package auth_test

import (
	"context"
	"os"
	"testing"

	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/environments"
)

var (
	tenantId           = os.Getenv("TENANT_ID")
	clientId           = os.Getenv("CLIENT_ID")
	clientCertificate  = os.Getenv("CLIENT_CERTIFICATE")
	clientCertPassword = os.Getenv("CLIENT_CERTIFICATE_PASSWORD")
	clientSecret       = os.Getenv("CLIENT_SECRET")
)

func TestClientCertificateAuthorizerV1(t *testing.T) {
	ctx := context.Background()
	auth, err := auth.NewClientCertificateAuthorizer(ctx, environments.Global, auth.MsGraph, auth.TokenVersion1, tenantId, clientId, clientCertificate, clientCertPassword)
	if err != nil {
		t.Errorf("NewClientCertificateAuthorizer(): %v", err)
	}
	if auth == nil {
		t.Error("auth is nil, expected Authorizer")
	}
	token, err := auth.Token()
	if err != nil {
		t.Errorf("auth.Token(): %v", err)
	}
	if token.AccessToken == "" {
		t.Error("token.AccessToken was empty")
	}
}

func TestClientCertificateAuthorizerV2(t *testing.T) {
	ctx := context.Background()
	auth, err := auth.NewClientCertificateAuthorizer(ctx, environments.Global, auth.MsGraph, auth.TokenVersion2, tenantId, clientId, clientCertificate, clientCertPassword)
	if err != nil {
		t.Errorf("NewClientCertificateAuthorizer(): %v", err)
	}
	if auth == nil {
		t.Error("auth is nil, expected Authorizer")
	}
	token, err := auth.Token()
	if err != nil {
		t.Errorf("auth.Token(): %v", err)
	}
	if token.AccessToken == "" {
		t.Error("token.AccessToken was empty")
	}
}

func TestClientSecretAuthorizerV1(t *testing.T) {
	ctx := context.Background()
	auth, err := auth.NewClientSecretAuthorizer(ctx, environments.Global, auth.MsGraph, auth.TokenVersion1, tenantId, clientId, clientSecret)
	if err != nil {
		t.Errorf("NewClientSecretAuthorizer(): %v", err)
	}
	if auth == nil {
		t.Error("auth is nil, expected Authorizer")
	}
	token, err := auth.Token()
	if err != nil {
		t.Errorf("auth.Token(): %v", err)
	}
	if token.AccessToken == "" {
		t.Errorf("token.AccessToken was empty")
	}
}

func TestClientSecretAuthorizerV2(t *testing.T) {
	ctx := context.Background()
	auth, err := auth.NewClientSecretAuthorizer(ctx, environments.Global, auth.MsGraph, auth.TokenVersion2, tenantId, clientId, clientSecret)
	if err != nil {
		t.Errorf("NewClientSecretAuthorizer(): %v", err)
	}
	if auth == nil {
		t.Error("auth is nil, expected Authorizer")
	}
	token, err := auth.Token()
	if err != nil {
		t.Errorf("auth.Token(): %v", err)
	}
	if token.AccessToken == "" {
		t.Errorf("token.AccessToken was empty")
	}
}

func TestAzureCliAuthorizer(t *testing.T) {
	ctx := context.Background()
	auth, err := auth.NewAzureCliAuthorizer(ctx, auth.MsGraph, tenantId)
	if err != nil {
		t.Errorf("NewAzureCliAuthorizer(): %v", err)
	}
	if auth == nil {
		t.Error("auth is nil, expected Authorizer")
	}
	token, err := auth.Token()
	if err != nil {
		t.Errorf("auth.Token(): %v", err)
	}
	if token.AccessToken == "" {
		t.Errorf("token.AccessToken was empty")
	}
}
