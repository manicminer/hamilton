package auth_test

import (
	"context"
	"encoding/base64"
	"os"
	"testing"

	"github.com/manicminer/hamilton/internal/test"

	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/environments"
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
	var pfx []byte
	if clientCertificate != "" {
		out := make([]byte, base64.StdEncoding.DecodedLen(len(clientCertificate)))
		n, err := base64.StdEncoding.Decode(out, []byte(clientCertificate))
		if err != nil {
			t.Fatalf("NewClientCertificateAuthorizer(): could not decode value of CLIENT_CERTIFICATE: %v", err)
		}
		pfx = out[:n]
	}
	auth, err := auth.NewClientCertificateAuthorizer(ctx, environments.Global, auth.MsGraph, auth.TokenVersion1, tenantId, clientId, pfx, clientCertificatePath, clientCertPassword)
	if err != nil {
		t.Fatalf("NewClientCertificateAuthorizer(): %v", err)
	}
	if auth == nil {
		t.Fatal("auth is nil, expected Authorizer")
	}
	token, err := auth.Token()
	if err != nil {
		t.Fatalf("auth.Token(): %v", err)
	}
	if token == nil {
		t.Fatalf("token was nil")
	}
	if token.AccessToken == "" {
		t.Fatal("token.AccessToken was empty")
	}
}

func TestClientCertificateAuthorizerV2(t *testing.T) {
	ctx := context.Background()
	var pfx []byte
	if clientCertificate != "" {
		out := make([]byte, base64.StdEncoding.DecodedLen(len(clientCertificate)))
		n, err := base64.StdEncoding.Decode(out, []byte(clientCertificate))
		if err != nil {
			t.Fatalf("NewClientCertificateAuthorizer(): could not decode value of CLIENT_CERTIFICATE: %v", err)
		}
		pfx = out[:n]
	}
	auth, err := auth.NewClientCertificateAuthorizer(ctx, environments.Global, auth.MsGraph, auth.TokenVersion2, tenantId, clientId, pfx, clientCertificate, clientCertPassword)
	if err != nil {
		t.Fatalf("NewClientCertificateAuthorizer(): %v", err)
	}
	if auth == nil {
		t.Fatal("auth is nil, expected Authorizer")
	}
	token, err := auth.Token()
	if err != nil {
		t.Fatalf("auth.Token(): %v", err)
	}
	if token == nil {
		t.Fatalf("token was nil")
	}
	if token.AccessToken == "" {
		t.Fatal("token.AccessToken was empty")
	}
}

func TestClientSecretAuthorizerV1(t *testing.T) {
	ctx := context.Background()
	auth, err := auth.NewClientSecretAuthorizer(ctx, environments.Global, auth.MsGraph, auth.TokenVersion1, tenantId, clientId, clientSecret)
	if err != nil {
		t.Fatalf("NewClientSecretAuthorizer(): %v", err)
	}
	if auth == nil {
		t.Fatal("auth is nil, expected Authorizer")
	}
	token, err := auth.Token()
	if err != nil {
		t.Fatalf("auth.Token(): %v", err)
	}
	if token == nil {
		t.Fatalf("token was nil")
	}
	if token.AccessToken == "" {
		t.Fatalf("token.AccessToken was empty")
	}
}

func TestClientSecretAuthorizerV2(t *testing.T) {
	ctx := context.Background()
	auth, err := auth.NewClientSecretAuthorizer(ctx, environments.Global, auth.MsGraph, auth.TokenVersion2, tenantId, clientId, clientSecret)
	if err != nil {
		t.Fatalf("NewClientSecretAuthorizer(): %v", err)
	}
	if auth == nil {
		t.Fatal("auth is nil, expected Authorizer")
	}
	token, err := auth.Token()
	if err != nil {
		t.Fatalf("auth.Token(): %v", err)
	}
	if token == nil {
		t.Fatalf("token was nil")
	}
	if token.AccessToken == "" {
		t.Fatalf("token.AccessToken was empty")
	}
}

func TestAzureCliAuthorizer(t *testing.T) {
	ctx := context.Background()
	auth, err := auth.NewAzureCliAuthorizer(ctx, auth.MsGraph, tenantId)
	if err != nil {
		t.Fatalf("NewAzureCliAuthorizer(): %v", err)
	}
	if auth == nil {
		t.Fatal("auth is nil, expected Authorizer")
	}
	token, err := auth.Token()
	if err != nil {
		t.Fatalf("auth.Token(): %v", err)
	}
	if token == nil {
		t.Fatalf("token was nil")
	}
	if token.AccessToken == "" {
		t.Fatalf("token.AccessToken was empty")
	}
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
