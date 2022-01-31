package containerregistry

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/environments"
	"golang.org/x/oauth2"
)

func TestContainerRegistryClient(t *testing.T) {
	fa := testNewFakeAuthorizer(t)
	h := testNewACRHandler(t)
	httpServer := httptest.NewTLSServer(h.handler(t))
	h.serverURL = httpServer.URL
	h.expectedTenant = ""
	testExchangeRefreshTokenSuccess(t, fa, httpServer.URL, "", httpServer.Client())
	testExchangeAccessTokenSuccess(t, fa, httpServer.URL, "", httpServer.Client())

	h.expectedTenant = "ze-tenant"
	testExchangeRefreshTokenSuccess(t, fa, httpServer.URL, "ze-tenant", httpServer.Client())
	testExchangeAccessTokenSuccess(t, fa, httpServer.URL, "ze-tenant", httpServer.Client())

	h.expectedTenant = "ze-tenant"
	h.fakeError = fmt.Errorf("ze-fake-error")
	testExchangeRefreshTokenFailure(t, fa, httpServer.URL, "ze-tenant", httpServer.Client(), "ze-fake-error")
	testExchangeAccessTokenFailure(t, fa, httpServer.URL, "ze-tenant", httpServer.Client(), "ze-fake-error")
}

func TestContainerRegistryE2E(t *testing.T) {
	containerRegistryName := os.Getenv("CONTAINER_REGISTRY_NAME")
	if containerRegistryName == "" {
		t.Skip("environment variable CONTAINER_REGISTRY_NAME not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	serverURL := fmt.Sprintf("%s.azurecr.io", containerRegistryName)

	cr := testContainerRegistryClientE2E(t, ctx, serverURL)
	testCatalogClientE2E(t, ctx, cr)
}

func testContainerRegistryClientE2E(t *testing.T, ctx context.Context, serverURL string) *ContainerRegistryClient {
	t.Helper()

	authorizer := testNewAuthorizer(t, ctx)
	cr := NewContainerRegistryClient(authorizer, serverURL, "")
	refreshToken, rtClaims, err := cr.ExchangeRefreshToken(ctx)
	if err != nil {
		t.Fatalf("received unexpected error: %v", err)
	}

	if refreshToken == "" {
		t.Fatalf("refreshToken is empty")
	}

	if rtClaims.Issuer != "Azure Container Registry" {
		t.Fatalf("refresh token claim 'iss' (Issuer) expected to be 'Azure Container Registry', but received: %s", rtClaims.Issuer)
	}

	atScopes := AccessTokenScopes{
		{
			Type:    "registry",
			Name:    "catalog",
			Actions: []string{"*"},
		},
	}

	accessToken, atClaims, err := cr.ExchangeAccessToken(ctx, refreshToken, atScopes)
	if err != nil {
		t.Fatalf("received unexpected error: %v", err)
	}

	if accessToken == "" {
		t.Fatalf("accessToken is empty")
	}

	if atClaims.Issuer != "Azure Container Registry" {
		t.Fatalf("access token claim 'iss' (Issuer) expected to be 'Azure Container Registry', but received: %s", atClaims.Issuer)
	}

	return cr
}

func testCatalogClientE2E(t *testing.T, ctx context.Context, cr *ContainerRegistryClient) {
	t.Helper()

	catalogClient := NewCatalogClient(cr)
	repositories, err := catalogClient.List(ctx, nil, nil)
	if err != nil {
		t.Fatalf("received unexpected error: %v", err)
	}

	if len(repositories) < 2 {
		t.Fatalf("expected at least two repositories")
	}

	toStringPtr := func(v string) *string { return &v }
	toIntPtr := func(v int) *int { return &v }
	repositoriesLimit, err := catalogClient.List(ctx, toStringPtr(repositories[0]), toIntPtr(1))
	if err != nil {
		t.Fatalf("received unexpected error: %v", err)
	}

	if len(repositoriesLimit) != 1 {
		t.Fatalf("expected repositoriesLimit to contain exactly one repository")
	}

	if repositoriesLimit[0] != repositories[1] {
		t.Fatalf("expected %q (repositoriesLimit[0]) to be %q (repositories[1])", repositoriesLimit[0], repositories[1])
	}

	imageName := repositories[0]
	_, err = catalogClient.GetAttributes(ctx, imageName)
	if err != nil {
		t.Fatalf("received unexpected error: %v", err)
	}

	toBoolPtr := func(v bool) *bool { return &v }
	err = catalogClient.UpdateAttributes(ctx, imageName, RepositoryChangeableAttributes{ListEnabled: toBoolPtr(true)})
	if err != nil {
		t.Fatalf("received unexpected error: %v", err)
	}

	_, err = catalogClient.Delete(ctx, "image-that-does-not-exist")
	if err == nil {
		t.Fatal("expected error when running Delete()")
	}

	if !strings.Contains(err.Error(), "repository name not known to registry") {
		t.Fatalf("expected error of Delete() to contain 'repository name not known to registry' but received: %v", err)
	}
}

type testFakeAuthorizer struct {
	t *testing.T
}

func testNewFakeAuthorizer(t *testing.T) *testFakeAuthorizer {
	t.Helper()

	return &testFakeAuthorizer{
		t,
	}
}

func (fa *testFakeAuthorizer) Token() (*oauth2.Token, error) {
	fa.t.Helper()

	return &oauth2.Token{
		AccessToken: "foo",
	}, nil
}

func (fa *testFakeAuthorizer) AuxiliaryTokens() ([]*oauth2.Token, error) {
	fa.t.Helper()

	return nil, nil
}

type testACRHandler struct {
	serverURL      string
	expectedTenant string
	fakeError      error
}

func testNewACRHandler(t *testing.T) *testACRHandler {
	t.Helper()

	return &testACRHandler{}
}

func (h *testACRHandler) handler(t *testing.T) http.HandlerFunc {
	t.Helper()

	return func(w http.ResponseWriter, r *http.Request) {
		h.router(t, w, r)
	}
}

func (h *testACRHandler) router(t *testing.T, w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	switch path {
	case "/oauth2/exchange":
		h.refreshTokenHandler(t, w, r)
	case "/oauth2/token":
		h.accessTokenHandler(t, w, r)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("unknown path: %s", path)))
		return
	}
}

func testNewAuthorizer(t *testing.T, ctx context.Context) auth.Authorizer {
	t.Helper()

	authCfg := testNewAuthConfig(t)
	gitHubTokenURL := strings.TrimSpace(os.Getenv("ACTIONS_ID_TOKEN_REQUEST_URL"))
	gitHubToken := strings.TrimSpace(os.Getenv("ACTIONS_ID_TOKEN_REQUEST_TOKEN"))

	a, err := auth.NewClientCertificateAuthorizer(ctx, authCfg.Environment, authCfg.Environment.ResourceManager, authCfg.Version, authCfg.TenantID, authCfg.AuxiliaryTenantIDs, authCfg.ClientID, authCfg.ClientCertData, authCfg.ClientCertPath, authCfg.ClientCertPassword)
	if err == nil {
		return a
	}

	if authCfg.TenantID != "" && authCfg.ClientID != "" && authCfg.ClientSecret != "" {
		a, err = auth.NewClientSecretAuthorizer(ctx, authCfg.Environment, authCfg.Environment.ResourceManager, authCfg.Version, authCfg.TenantID, authCfg.AuxiliaryTenantIDs, authCfg.ClientID, authCfg.ClientSecret)
		if err == nil {
			return a
		}
	}

	a, err = auth.NewMsiAuthorizer(ctx, authCfg.Environment.ResourceManager, authCfg.MsiEndpoint, authCfg.ClientID)
	if err == nil {
		_, err := a.Token()
		if err == nil {
			return a
		}
	}

	if gitHubToken != "" || gitHubTokenURL != "" {
		a, err = auth.NewGitHubOIDCAuthorizer(ctx, authCfg.Environment, authCfg.Environment.ResourceManager, authCfg.TenantID, authCfg.AuxiliaryTenantIDs, authCfg.ClientID, gitHubTokenURL, gitHubToken)
		if err == nil {
			return a
		}
	}

	a, err = auth.NewAzureCliAuthorizer(ctx, authCfg.Environment.ResourceManager, authCfg.TenantID)
	if err == nil {
		return a
	}

	t.Fatalf("no valid authorizer could be found")
	return nil
}

func testNewAuthConfig(t *testing.T) *auth.Config {
	t.Helper()

	auxTenants := strings.Split(os.Getenv("AZURE_AUXILIARY_TENANT_IDS"), ";")
	for i := range auxTenants {
		auxTenants[i] = strings.TrimSpace(auxTenants[i])
	}

	return &auth.Config{
		Environment:        environments.Global,
		TenantID:           strings.TrimSpace(os.Getenv("AZURE_TENANT_ID")),
		ClientID:           strings.TrimSpace(os.Getenv("AZURE_CLIENT_ID")),
		ClientSecret:       strings.TrimSpace(os.Getenv("AZURE_CLIENT_SECRET")),
		ClientCertPath:     os.Getenv("AZURE_CERTIFICATE_PATH"),
		ClientCertPassword: os.Getenv("AZURE_CERTIFICATE_PASSWORD"),
		AuxiliaryTenantIDs: auxTenants,
		Version:            auth.TokenVersion2,
	}
}
