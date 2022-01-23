package autorest_test

import (
	"os"
	"testing"

	msautorest "github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"

	"github.com/manicminer/hamilton/auth/autorest"
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
	environment           = os.Getenv("AZURE_ENVIRONMENT")
	msiEndpoint           = os.Getenv("MSI_ENDPOINT")
	msiToken              = os.Getenv("MSI_TOKEN")
)

func TestAutorestAuthorizerWrapper(t *testing.T) {
	env, err := environments.EnvironmentFromString(environment)
	if err != nil {
		t.Fatal(err)
	}

	// adal.ServicePrincipalToken.refreshInternal() doesn't support v2 tokens
	oauthConfig, err := adal.NewOAuthConfigWithAPIVersion(string(env.AzureADEndpoint), tenantId, utils.StringPtr("1.0"))
	if err != nil {
		t.Fatalf("adal.NewOAuthConfig(): %v", err)
	}

	spt, err := adal.NewServicePrincipalToken(*oauthConfig, clientId, clientSecret, string(env.MsGraph.Endpoint))
	if err != nil {
		t.Fatalf("adal.NewServicePrincipalToken(): %v", err)
	}

	auth, err := autorest.NewAutorestAuthorizerWrapper(msautorest.NewBearerAuthorizer(spt))
	if err != nil {
		t.Fatalf("NewAutorestAuthorizerWrapper(): %v", err)
	}
	if auth == nil {
		t.Fatal("auth is nil, expected Authorizer")
	}

	token, err := auth.Token()
	if err != nil {
		t.Fatalf("auth.Token(): %v", err)
	}
	if token == nil {
		t.Fatal("token was nil")
	}
	if token.AccessToken == "" {
		t.Fatal("token.AccessToken was empty")
	}
}
