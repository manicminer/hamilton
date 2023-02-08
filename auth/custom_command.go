package auth

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"text/template"

	"golang.org/x/oauth2"

	"github.com/manicminer/hamilton/environments"
)

// CustomCommandConfig configures an CustomCommandAuthorizer.
type CustomCommandConfig struct {
	endpoint environments.ApiEndpoint

	// tenantID is the required tenant ID for the primary token
	tenantID string

	// auxiliaryTenantIDs is an optional list of tenant IDs for which to obtain additional tokens
	auxiliaryTenantIDs []string

	// tokenType is either this or "Bearer" by default
	tokenType string

	// primaryCommand is the rendered command for the primary tenant
	primaryCommand []string

	// auxiliaryCommands are a list of rendered commands for each auxiliary tenant
	auxiliaryCommands [][]string
}

// NewCustomCommandConfig pre-processes the input command and returns a CustomCommandConfig.
// The command is the exec form of command used to retrieve the access token.
// Each command argument will be rendered as a Go template with an input object that has following fields:
// - .Endpoint: The token endpoint
// - .TenantID: The tenant ID. For auxiliary tokens, it is set as one of each auxiliary token.
//
// E.g. []string{"az", "account", "get-access-token", "--resource={{.Endpoint}}"}"}
func NewCustomCommandConfig(api environments.Api, tenantId string, auxiliaryTenantIDs []string, tokenType string, command []string) (*CustomCommandConfig, error) {
	config := CustomCommandConfig{
		endpoint:           api.Endpoint,
		tenantID:           tenantId,
		auxiliaryTenantIDs: auxiliaryTenantIDs,
	}

	if tokenType != "" {
		config.tokenType = tokenType
	}

	buildCommand := func(endpoint environments.ApiEndpoint, tenantID string, rawCommand []string) ([]string, error) {
		if len(rawCommand) == 0 {
			return nil, fmt.Errorf("missing command")
		}
		command := make([]string, len(rawCommand))
		copy(command, rawCommand)

		for i, arg := range rawCommand {
			// Not render the command to run
			if i == 0 {
				continue
			}

			inputObj := struct {
				Endpoint environments.ApiEndpoint
				TenantID string
			}{
				Endpoint: endpoint,
				TenantID: tenantID,
			}
			tpl, err := template.New("arg").Parse(arg)
			if err != nil {
				return nil, fmt.Errorf("format of the %d-th argument is not valid: %v", i, err)
			}
			var buf bytes.Buffer
			if err := tpl.Execute(&buf, inputObj); err != nil {
				return nil, fmt.Errorf("format of the %d-th argument is not valid: %v", i, err)
			}
			command[i] = buf.String()
		}
		return command, nil
	}

	pcommand, err := buildCommand(config.endpoint, config.tenantID, command)
	if err != nil {
		return nil, err
	}
	config.primaryCommand = pcommand

	for _, tenantId := range config.auxiliaryTenantIDs {
		pcommand, err := buildCommand(config.endpoint, tenantId, command)
		if err != nil {
			return nil, err
		}
		config.auxiliaryCommands = append(config.auxiliaryCommands, pcommand)
	}

	return &config, nil
}

// TokenSource provides a source for obtaining access tokens using CustomCommandAuthorizer.
func (c *CustomCommandConfig) TokenSource(ctx context.Context) Authorizer {
	// Cache access tokens internally to avoid unnecessary command invocations
	return NewCachedAuthorizer(&CustomCommandAuthorizer{
		ctx:  ctx,
		conf: c,
	})
}

// CustomCommandAuthorizer is an Authorizer which supports the Azure CLI.
type CustomCommandAuthorizer struct {
	ctx  context.Context
	conf *CustomCommandConfig
}

// Token returns an access token using the command as an authentication mechanism.
func (a *CustomCommandAuthorizer) Token() (*oauth2.Token, error) {
	token, err := runCustomCommand(a.conf.primaryCommand)
	if err != nil {
		return nil, err
	}
	return &oauth2.Token{
		AccessToken: token,
		TokenType:   a.conf.tokenType,
	}, nil
}

// AuxiliaryTokens returns additional tokens for auxiliary tenant IDs, for use in multi-tenant scenarios
func (a *CustomCommandAuthorizer) AuxiliaryTokens() ([]*oauth2.Token, error) {
	tokens := make([]*oauth2.Token, 0)
	for _, command := range a.conf.auxiliaryCommands {
		token, err := runCustomCommand(command)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, &oauth2.Token{
			AccessToken: token,
			TokenType:   a.conf.tokenType,
		})
	}

	return tokens, nil
}

// runCustomCommand executes the custom command and return the output access token with any quote/emptyspace trimmed.
func runCustomCommand(command []string) (string, error) {
	var stderr bytes.Buffer
	var stdout bytes.Buffer

	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	if err := cmd.Start(); err != nil {
		err := fmt.Errorf("launching custom command: %+v", err)
		if stdErrStr := stderr.String(); stdErrStr != "" {
			err = fmt.Errorf("%s: %s", err, strings.TrimSpace(stdErrStr))
		}
		return "", err
	}

	if err := cmd.Wait(); err != nil {
		err := fmt.Errorf("running custom command: %+v", err)
		if stdErrStr := stderr.String(); stdErrStr != "" {
			err = fmt.Errorf("%s: %s", err, strings.TrimSpace(stdErrStr))
		}
		return "", err
	}

	return strings.Trim(strings.TrimSpace(stdout.String()), `"`), nil
}
