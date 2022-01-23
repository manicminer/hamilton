package containerregistry

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/docker/docker-credential-helpers/credentials"
	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/environments"
)

type ACRHelper struct{}

var _ credentials.Helper = (*ACRHelper)(nil)

func NewACRHelper() credentials.Helper {
	return &ACRHelper{}
}

func (ACRHelper) Add(_ *credentials.Credentials) error {
	return fmt.Errorf("method Add() is not implemented")
}

func (ACRHelper) Delete(_ string) error {
	return fmt.Errorf("method Delete() not implemented")
}

func (self ACRHelper) Get(serverURL string) (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	authConfig := self.newConfig()
	authorizer, err := self.newAuthorizer(ctx, authConfig)
	if err != nil {
		return "", "", err
	}

	token, err := authorizer.Token()
	if err != nil {
		return "", "", err
	}

	acrToken, err := self.getACRToken(ctx, token.AccessToken, serverURL, authConfig.TenantID)
	if err != nil {
		return "", "", err
	}

	return "<token>", acrToken, nil
}

func (self ACRHelper) List() (map[string]string, error) {
	return nil, fmt.Errorf("method List() is not implemented")
}

func (self ACRHelper) newConfig() *auth.Config {
	auxTenants := strings.Split(os.Getenv("AZURE_AUXILIARY_TENANT_IDS"), ";")
	for i := range auxTenants {
		auxTenants[i] = strings.TrimSpace(auxTenants[i])
	}

	return &auth.Config{
		Environment:        environments.Global,
		TenantID:           strings.TrimSpace(os.Getenv("AZURE_TENANT_ID")),
		ClientID:           strings.TrimSpace(os.Getenv("AZURE_CLIENT_ID")),
		ClientSecret:       strings.TrimSpace(os.Getenv("AZURE_CLIENT_SECRET")),
		ClientCertPath:     os.Getenv("AZURE_CERTIFICATE_PATH"),
		ClientCertPassword: os.Getenv("AZURE_CERTIFICATE_PASSWORD"),
		AuxiliaryTenantIDs: auxTenants,
		Version:            auth.TokenVersion2,
	}
}

func (self ACRHelper) newAuthorizer(ctx context.Context, authCfg *auth.Config) (auth.Authorizer, error) {
	a, err := auth.NewClientCertificateAuthorizer(ctx, authCfg.Environment, authCfg.Environment.ResourceManager, authCfg.Version, authCfg.TenantID, authCfg.AuxiliaryTenantIDs, authCfg.ClientID, authCfg.ClientCertData, authCfg.ClientCertPath, authCfg.ClientCertPassword)
	if err == nil {
		return a, nil
	}

	if authCfg.TenantID != "" && authCfg.ClientID != "" && authCfg.ClientSecret != "" {
		a, err = auth.NewClientSecretAuthorizer(ctx, authCfg.Environment, authCfg.Environment.ResourceManager, authCfg.Version, authCfg.TenantID, authCfg.AuxiliaryTenantIDs, authCfg.ClientID, authCfg.ClientSecret)
		if err == nil {
			return a, nil
		}
	}

	a, err = auth.NewMsiAuthorizer(ctx, authCfg.Environment.ResourceManager, authCfg.MsiEndpoint, authCfg.ClientID)
	if err == nil {
		return a, nil
	}

	a, err = auth.NewAzureCliAuthorizer(ctx, authCfg.Environment.ResourceManager, authCfg.TenantID)
	if err == nil {
		return a, nil
	}

	return nil, fmt.Errorf("no valid authorizer could be found")
}

func (self ACRHelper) getACRToken(ctx context.Context, accessToken string, serverURL string, tenantID string) (string, error) {
	service, err := self.getService(serverURL)
	if err != nil {
		return "", err
	}

	data := url.Values{}
	data.Set("grant_type", "access_token")
	data.Set("service", service)
	data.Set("access_token", accessToken)
	if len(tenantID) > 0 {
		data.Set("tenant", tenantID)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("https://%s/oauth2/exchange", service), strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 status code: %d", res.StatusCode)
	}

	var resData struct {
		RefreshToken string `json:"refresh_token"`
	}

	err = json.Unmarshal(resBytes, &resData)

	return resData.RefreshToken, nil
}

func (self ACRHelper) getService(serverURL string) (string, error) {
	scheme := "https://"
	if strings.HasPrefix(serverURL, "https://") {
		scheme = ""
	}

	serviceURL, err := url.Parse(fmt.Sprintf("%s%s", scheme, serverURL))
	if err != nil {
		return "", fmt.Errorf("failed to parse server URL - %w", err)
	}

	return serviceURL.Hostname(), nil
}
