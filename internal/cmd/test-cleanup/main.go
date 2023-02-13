package main

import (
	"context"
	"log"
	"os"

	"github.com/hashicorp/go-azure-sdk/sdk/auth"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
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
	env := environments.AzurePublic()

	creds := auth.Credentials{
		Environment:               *env,
		TenantID:                  tenantId,
		ClientID:                  clientId,
		ClientCertificateData:     utils.Base64DecodeCertificate(clientCertificate),
		ClientCertificatePath:     clientCertificatePath,
		ClientCertificatePassword: clientCertPassword,
		ClientSecret:              clientSecret,
		EnableAuthenticatingUsingClientCertificate: true,
		EnableAuthenticatingUsingClientSecret:      true,
	}

	var err error
	authorizer, err = auth.NewAuthorizerFromCredentials(ctx, creds, env.MicrosoftGraph)
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	log.Println("Starting test cleanup...")
	cleanupAdministrativeUnits()
	cleanupConditionalAccessPolicies()
	cleanupNamedLocations()
	cleanupServicePrincipals()
	cleanupApplications()
	cleanupGroups()
	cleanupUsers()
	cleanupSchemaExtensions()
	log.Println("Finished test cleanup")
}
