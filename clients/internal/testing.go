package internal

import (
	"context"
	"log"
	"os"

	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/environments"
)

var (
	tenantId          = os.Getenv("TENANT_ID")
	tenantDomain      = os.Getenv("TENANT_DOMAIN")
	clientId          = os.Getenv("CLIENT_ID")
	//clientCertificate = os.Getenv("CLIENT_CERTIFICATE")
	clientSecret      = os.Getenv("CLIENT_SECRET")
)

type Connection struct {
	AuthConfig *auth.Config
	Authorizer auth.Authorizer
	Context    context.Context
	DomainName string
}

func NewConnection() *Connection {
	t := Connection{
		AuthConfig: &auth.Config{
			Environment:            environments.Global,
			TenantID:               tenantId,
			ClientID:               clientId,
			ClientSecret:           clientSecret,
			EnableClientCertAuth:   true,
			EnableClientSecretAuth: true,
			EnableAzureCliToken:    true,
		},
		Context:    context.Background(),
		DomainName: tenantDomain,
	}

	var err error
	t.Authorizer, err = t.AuthConfig.NewAuthorizer(t.Context)
	if err != nil {
		log.Fatal(err)
	}

	return &t
}

func Bool(b bool) *bool {
	return &b
}

func String(s string) *string {
	return &s
}
