package auth_test

import (
	"context"
	"os"
	"testing"

	"golang.org/x/oauth2"

	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/environments"
	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/internal/utils"
)

var (
	tenantId              = os.Getenv("TENANT_ID")
	clientId              = os.Getenv("CLIENT_ID")
	clientCertificate     = os.Getenv("CLIENT_CERTIFICATE")
	clientCertificatePath = os.Getenv("CLIENT_CERTIFICATE_PATH")
	clientCertPassword    = os.Getenv("CLIENT_CERTIFICATE_PASSWORD")
	clientSecret          = os.Getenv("CLIENT_SECRET")
	msiEndpoint           = os.Getenv("MSI_ENDPOINT")
	msiToken              = os.Getenv("MSI_TOKEN")
)

func TestClientCertificateAuthorizerV1(t *testing.T) {
	ctx := context.Background()
	testClientCertificateAuthorizer(ctx, t, auth.TokenVersion1)
}

func TestClientCertificateAuthorizerV2(t *testing.T) {
	ctx := context.Background()
	testClientCertificateAuthorizer(ctx, t, auth.TokenVersion2)
}

func testClientCertificateAuthorizer(ctx context.Context, t *testing.T, tokenVersion auth.TokenVersion) (token *oauth2.Token) {
	pfx := utils.Base64DecodeCertificate(clientCertificate)
	auth, err := auth.NewClientCertificateAuthorizer(ctx, environments.Global, auth.MsGraph, tokenVersion, tenantId, clientId, pfx, clientCertificatePath, clientCertPassword)
	if err != nil {
		t.Fatalf("NewClientCertificateAuthorizer(): %v", err)
	}
	if auth == nil {
		t.Fatal("auth is nil, expected Authorizer")
	}
	token, err = auth.Token()
	if err != nil {
		t.Fatalf("auth.Token(): %v", err)
	}
	if token == nil {
		t.Fatalf("token was nil")
	}
	if token.AccessToken == "" {
		t.Fatal("token.AccessToken was empty")
	}
	return
}

func TestClientSecretAuthorizerV1(t *testing.T) {
	ctx := context.Background()
	testClientSecretAuthorizer(ctx, t, auth.TokenVersion1)
}

func TestClientSecretAuthorizerV2(t *testing.T) {
	ctx := context.Background()
	testClientSecretAuthorizer(ctx, t, auth.TokenVersion2)
}

func testClientSecretAuthorizer(ctx context.Context, t *testing.T, tokenVersion auth.TokenVersion) (token *oauth2.Token) {
	auth, err := auth.NewClientSecretAuthorizer(ctx, environments.Global, auth.MsGraph, tokenVersion, tenantId, clientId, clientSecret)
	if err != nil {
		t.Fatalf("NewClientSecretAuthorizer(): %v", err)
	}
	if auth == nil {
		t.Fatal("auth is nil, expected Authorizer")
	}
	token, err = auth.Token()
	if err != nil {
		t.Fatalf("auth.Token(): %v", err)
	}
	if token == nil {
		t.Fatalf("token was nil")
	}
	if token.AccessToken == "" {
		t.Fatalf("token.AccessToken was empty")
	}
	return
}

func TestAzureCliAuthorizer(t *testing.T) {
	ctx := context.Background()
	testAzureCliAuthorizer(ctx, t)
}

func testAzureCliAuthorizer(ctx context.Context, t *testing.T) (token *oauth2.Token) {
	auth, err := auth.NewAzureCliAuthorizer(ctx, auth.MsGraph, tenantId)
	if err != nil {
		t.Fatalf("NewAzureCliAuthorizer(): %v", err)
	}
	if auth == nil {
		t.Fatal("auth is nil, expected Authorizer")
	}
	token, err = auth.Token()
	if err != nil {
		t.Fatalf("auth.Token(): %v", err)
	}
	if token == nil {
		t.Fatalf("token was nil")
	}
	if token.AccessToken == "" {
		t.Fatalf("token.AccessToken was empty")
	}
	return
}

func TestMsiAuthorizer(t *testing.T) {
	ctx := context.Background()
	if msiToken != "" {
		msiEndpoint = "http://localhost:8080/metadata/identity/oauth2/token"
		done := test.MsiStubServer(ctx, 8080, msiToken)
		defer func() {
			done <- true
		}()
	}
	auth, err := auth.NewMsiAuthorizer(ctx, environments.Global, auth.MsGraph, msiEndpoint)
	if err != nil {
		t.Fatalf("NewMsiAuthorizer(): %v", err)
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
