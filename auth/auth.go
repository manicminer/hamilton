package auth

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"strings"

	"golang.org/x/crypto/pkcs12"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	microsoft2 "github.com/manicminer/hamilton/auth/internal/microsoft"
	"github.com/manicminer/hamilton/environments"
)

type Config struct {
	// Specifies the national cloud environment to use
	Environment environments.Environment

	// Azure Active Directory tenant to connect to, should be a valid UUID
	TenantID string

	// Client ID for the application used to authenticate the connection
	ClientID string

	// Enables authentication using Azure CLI
	EnableAzureCliToken bool

	// Enables authentication using managed service identity. Not yet supported.
	// TODO: NOT YET SUPPORTED
	EnableMsiAuth bool

	// Specifies a custom MSI endpoint to connect to
	MsiEndpoint string

	// Enables client certificate authentication using client assertions
	EnableClientCertAuth bool

	// Specifies the path to a client certificate bundle in PFX format
	ClientCertPath string

	// Specifies the encryption password to unlock a client certificate
	ClientCertPassword string

	// Enables client secret authentication using client credentials
	EnableClientSecretAuth bool

	// Specifies the password to authenticate with using client secret authentication
	ClientSecret string
}

// Authorizer is anything that can return an access token for authorizing API connections
type Authorizer interface {
	Token() (*oauth2.Token, error)
}

// NewAuthorizer returns a suitable Authorizer depending on what is defined in the Config
// Authorizers are selected for authentication methods in the following preferential order:
// - Client certificate authentication
// - Client secret authentication
// - Azure CLI authentication
//
// Whether one of these is returned depends on whether it is enabled in the Config, and whether sufficient
// configuration fields are set to enable that authentication method.
//
// For client certificate authentication, specify TenantID, ClientID and ClientCertPath.
// For client secret authentication, specify TenantID, ClientID and ClientSecret.
// Azure CLI authentication (if enabled) is used as a fallback mechanism.
func (c *Config) NewAuthorizer(ctx context.Context) (Authorizer, error) {
	if c.EnableClientCertAuth && strings.TrimSpace(c.TenantID) != "" && strings.TrimSpace(c.ClientID) != "" && strings.TrimSpace(c.ClientCertPath) != "" {
		a, err := NewClientCertificateAuthorizer(ctx, c.Environment, c.TenantID, c.ClientID, c.ClientCertPath, c.ClientCertPassword)
		if err != nil {
			return nil, fmt.Errorf("could not configure ClientCertificate Authorizer: %s", err)
		}
		if a != nil {
			return a, nil
		}
	}

	if c.EnableClientSecretAuth && strings.TrimSpace(c.TenantID) != "" && strings.TrimSpace(c.ClientID) != "" && strings.TrimSpace(c.ClientSecret) != "" {
		a, err := NewClientSecretAuthorizer(ctx, c.Environment, c.TenantID, c.ClientID, c.ClientSecret)
		if err != nil {
			return nil, fmt.Errorf("could not configure ClientCertificate Authorizer: %s", err)
		}
		if a != nil {
			return a, nil
		}
	}

	if c.EnableAzureCliToken {
		a, err := NewAzureCliAuthorizer(ctx, c.TenantID)
		if err != nil {
			return nil, fmt.Errorf("could not configure AzureCli Authorizer: %s", err)
		}
		if a != nil {
			return a, nil
		}
	}

	return nil, fmt.Errorf("no Authorizer could be configured, please check your configuration")
}

// NewAzureCliAuthorizer returns an Authorizer which authenticates using the Azure CLI.
func NewAzureCliAuthorizer(ctx context.Context, tenantId string) (Authorizer, error) {
	conf, err := NewAzureCliConfig(tenantId)
	if err != nil {
		return nil, err
	}
	return conf.TokenSource(ctx), nil
}

// NewClientCertificateAuthorizer returns an authorizer which uses client certificate authentication.
func NewClientCertificateAuthorizer(ctx context.Context, environment environments.Environment, tenantId, clientId, pfxPath, pfxPass string) (Authorizer, error) {
	pfx, err := ioutil.ReadFile(pfxPath)
	if err != nil {
		return nil, fmt.Errorf("could not read pkcs12 store at %q: %s", pfxPath, err)
	}

	key, cert, err := pkcs12.Decode(pfx, pfxPass)
	if err != nil {
		return nil, fmt.Errorf("could not decode pkcs12 credential store: %s", err)
	}

	priv, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("unsupported non-rsa key was found in pkcs12 store %q", pfxPath)
	}

	conf := microsoft2.Config{
		ClientID:    clientId,
		PrivateKey:  x509.MarshalPKCS1PrivateKey(priv),
		Certificate: cert.Raw,
		Scopes:      []string{fmt.Sprintf("%s/.default", environment.MsGraphEndpoint)},
		TokenURL:    environments.AzureAD(environment.AzureADEndpoint, tenantId).TokenURL,
	}
	return conf.TokenSource(ctx), nil
}

// NewClientSecretAuthorizer returns an authorizer which uses client secret authentication.
func NewClientSecretAuthorizer(ctx context.Context, environment environments.Environment, tenantId, clientId, clientSecret string) (Authorizer, error) {
	conf := clientcredentials.Config{
		AuthStyle:    oauth2.AuthStyleInParams,
		ClientID:     clientId,
		ClientSecret: clientSecret,
		Scopes:       []string{fmt.Sprintf("%s/.default", environment.MsGraphEndpoint)},
		TokenURL:     environments.AzureAD(environment.AzureADEndpoint, tenantId).TokenURL,
	}
	return conf.TokenSource(ctx), nil
}
