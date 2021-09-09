package auth_test

import (
	"context"
	"testing"

	"github.com/manicminer/hamilton/auth"
)

func TestParseClaims_azureCli(t *testing.T) {
	ctx := context.Background()
	token := testAzureCliAuthorizer(ctx, t)
	claims, err := auth.ParseClaims(token)
	if err != nil {
		t.Fatal(err)
	}
	checkClaims(t, claims)
}

func TestParseClaims_clientCertificate(t *testing.T) {
	ctx := context.Background()
	token := testClientCertificateAuthorizer(ctx, t, auth.TokenVersion2)
	claims, err := auth.ParseClaims(token)
	if err != nil {
		t.Fatal(err)
	}
	checkClaims(t, claims)
}

func TestParseClaims_clientSecret(t *testing.T) {
	ctx := context.Background()
	token := testClientSecretAuthorizer(ctx, t, auth.TokenVersion2)
	claims, err := auth.ParseClaims(token)
	if err != nil {
		t.Fatal(err)
	}
	checkClaims(t, claims)
}

func checkClaims(t *testing.T, claims auth.Claims) {
	if claims.AppId == "" {
		t.Fatal("claims.AppId was empty")
	}
	if claims.Audience == "" {
		t.Fatal("claims.Audience was empty")
	}
	if claims.Issuer == "" {
		t.Fatal("claims.Issuer was empty")
	}
	if len(claims.Roles) == 0 && claims.Scopes == "" {
		t.Fatal("claims.Roles and claims.Scopes were empty")
	}
	if claims.Subject == "" {
		t.Fatal("claims.Subject was empty")
	}
	if claims.TenantId == "" {
		t.Fatal("claims.TenantId was empty")
	}
}
