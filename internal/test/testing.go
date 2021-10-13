package test

import (
	"context"
	"log"
	"os"

	"github.com/manicminer/hamilton/internal/utils"

	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/environments"
)

var (
	tenantId              = os.Getenv("TENANT_ID")
	tenantDomain          = os.Getenv("TENANT_DOMAIN")
	clientId              = os.Getenv("CLIENT_ID")
	clientCertificate     = os.Getenv("CLIENT_CERTIFICATE")
	clientCertificatePath = os.Getenv("CLIENT_CERTIFICATE_PATH")
	clientCertPassword    = os.Getenv("CLIENT_CERTIFICATE_PASSWORD")
	clientSecret          = os.Getenv("CLIENT_SECRET")
	environment           = os.Getenv("AZURE_ENVIRONMENT")
)

type Connection struct {
	AuthConfig *auth.Config
	Authorizer auth.Authorizer
	Context    context.Context
	DomainName string
}

// NewConnection configures and returns a Connection for use in tests.
func NewConnection(api auth.Api, tokenVersion auth.TokenVersion) *Connection {
	env := environments.Global
	switch environment {
	case "usgovernment", "usgovernmentl4":
		env = environments.USGovernmentL4
	case "dod", "usgovernmentl5":
		env = environments.USGovernmentL5
	case "canary":
		env = environments.Canary
	case "china":
		env = environments.China
	case "germany":
		env = environments.Germany
	}
	t := Connection{
		AuthConfig: &auth.Config{
			Environment:            env,
			Version:                tokenVersion,
			TenantID:               tenantId,
			ClientID:               clientId,
			ClientCertData:         utils.Base64DecodeCertificate(clientCertificate),
			ClientCertPath:         clientCertificatePath,
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
