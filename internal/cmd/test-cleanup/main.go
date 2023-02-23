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

	b2cTenantId = os.Getenv("B2C_TENANT_ID")
)

var (
	ctx           context.Context
	authorizer    auth.Authorizer
	b2cAuthorizer auth.Authorizer
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

	b2cCreds := auth.Credentials{
		Environment:               *env,
		TenantID:                  b2cTenantId,
		ClientID:                  clientId,
		ClientCertificateData:     utils.Base64DecodeCertificate(clientCertificate),
		ClientCertificatePath:     clientCertificatePath,
		ClientCertificatePassword: clientCertPassword,
		ClientSecret:              clientSecret,
		EnableAuthenticatingUsingClientCertificate: true,
		EnableAuthenticatingUsingClientSecret:      true,
	}

	b2cAuthorizer, err = auth.NewAuthorizerFromCredentials(ctx, b2cCreds, env.MicrosoftGraph)
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	log.Println("Starting test cleanup...")
	cleanupAccessPackages()
	cleanupAccessPackageAssignmentPolicies()
	cleanupAccessPackageAssignmentRequests()
	cleanupAccessPackageCatalogs()
	cleanupAdministrativeUnits()
	cleanupB2CUserFlows()
	cleanupClaimsMappingPolicies()
	cleanupConditionalAccessPolicies()
	cleanupConnectedOrganizations()
	cleanupNamedLocations()
	cleanupServicePrincipals()
	cleanupApplications()
	cleanupGroups()
	cleanupUsers()
	cleanupSchemaExtensions()
	cleanupRoleDefinitions()
	log.Println("Finished test cleanup")
}
