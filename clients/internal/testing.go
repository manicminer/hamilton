package internal

import (
	"context"
	"log"
	"os"

	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/environments"
)

var (
	tenantId           = os.Getenv("TENANT_ID")
	tenantDomain       = os.Getenv("TENANT_DOMAIN")
	clientId           = os.Getenv("CLIENT_ID")
	clientCertificate  = os.Getenv("CLIENT_CERTIFICATE")
	clientCertPassword = os.Getenv("CLIENT_CERTIFICATE_PASSWORD")
	clientSecret       = os.Getenv("CLIENT_SECRET")
)

type Connection struct {
	AuthConfig *auth.Config
	Authorizer auth.Authorizer
	Context    context.Context
	DomainName string
}

// NewConnection configures and returns a Connection for use in tests.
func NewConnection(api auth.Api, tokenVersion auth.TokenVersion) *Connection {
	t := Connection{
		AuthConfig: &auth.Config{
			Environment:            environments.Global,
			Version:                tokenVersion,
			TenantID:               tenantId,
			ClientID:               clientId,
			ClientCertPath:         clientCertificate,
			ClientCertPassword:     clientCertPassword,
			ClientSecret:           clientSecret,
			EnableClientCertAuth:   true,
			EnableClientSecretAuth: true,
			EnableAzureCliToken:    true,
		},
		Context:    context.Background(),
		DomainName: tenantDomain,
	}

	var err error
	t.Authorizer, err = t.AuthConfig.NewAuthorizer(t.Context, api)
	if err != nil {
		log.Fatal(err)
	}

	return &t
}
