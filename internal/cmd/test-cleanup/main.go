package main

import (
	"context"
	"log"
	"os"

	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/environments"
	"github.com/manicminer/hamilton/internal/utils"
)

var (
	tenantId              = os.Getenv("TENANT_ID")
	clientId              = os.Getenv("CLIENT_ID")
	clientCertificate     = os.Getenv("CLIENT_CERTIFICATE")
	clientCertificatePath = os.Getenv("CLIENT_CERTIFICATE_PATH")
	clientCertPassword    = os.Getenv("CLIENT_CERTIFICATE_PASSWORD")
	clientSecret          = os.Getenv("CLIENT_SECRET")
)

var (
	ctx        context.Context
	authorizer auth.Authorizer
)

const displayNamePrefix = "test-"

func init() {
	ctx = context.Background()

	authConfig := &auth.Config{
		Environment:            environments.Global,
		TenantID:               tenantId,
		ClientID:               clientId,
		ClientCertData:         utils.Base64DecodeCertificate(clientCertificate),
		ClientCertPath:         clientCertificatePath,
		ClientCertPassword:     clientCertPassword,
		ClientSecret:           clientSecret,
		EnableClientCertAuth:   true,
		EnableClientSecretAuth: true,
	}

	var err error
	authorizer, err = authConfig.NewAuthorizer(ctx, environments.Global.MsGraph)
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	log.Println("Starting test cleanup...")
	cleanupConditionalAccessPolicies()
	cleanupNamedLocations()
	cleanupServicePrincipals()
	cleanupApplications()
	cleanupGroups()
	cleanupUsers()
	cleanupSchemaExtensions()
	log.Println("Finished test cleanup")
}
