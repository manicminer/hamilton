package containerregistry

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strconv"
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

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	serverURL := fmt.Sprintf("%s.azurecr.io", containerRegistryName)

	cr := testContainerRegistryClientE2E(t, ctx, serverURL)
	imageName := testCatalogE2E(t, ctx, cr)
	testTagE2E(t, ctx, cr, imageName)
	testManifestE2E(t, ctx, cr, imageName)
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

func testCatalogE2E(t *testing.T, ctx context.Context, cr *ContainerRegistryClient) string {
	t.Helper()

	repositories, err := cr.CatalogList(ctx, nil, nil)
	if err != nil {
		t.Fatalf("received unexpected error: %v", err)
	}

	if len(repositories) < 2 {
		t.Fatalf("expected at least two repositories")
	}

	toStringPtr := func(v string) *string { return &v }
	toIntPtr := func(v int) *int { return &v }
	repositoriesLimit, err := cr.CatalogList(ctx, toStringPtr(repositories[0]), toIntPtr(1))
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
	attributes, err := cr.CatalogGetAttributes(ctx, imageName)
	if err != nil {
		t.Fatalf("received unexpected error: %v", err)
	}

	if attributes.ImageName != imageName {
		t.Errorf("expected attributes.ImageName to be %q but received: %s", imageName, attributes.ImageName)
	}

	toBoolPtr := func(v bool) *bool { return &v }
	err = cr.CatalogUpdateAttributes(ctx, imageName, RepositoryChangeableAttributes{ListEnabled: toBoolPtr(true)})
	if err != nil {
		t.Fatalf("received unexpected error: %v", err)
	}

	_, err = cr.CatalogDelete(ctx, "image-that-does-not-exist")
	if err == nil {
		t.Fatal("expected error when running Delete()")
	}

	if !strings.Contains(err.Error(), "repository name not known to registry") {
		t.Fatalf("expected error of Delete() to contain 'repository name not known to registry' but received: %v", err)
	}

	return repositories[0]
}

func testTagE2E(t *testing.T, ctx context.Context, cr *ContainerRegistryClient, imageName string) {
	t.Helper()

	tagList, err := cr.TagList(ctx, imageName, nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("received unexpected error: %v", err)
	}

	if len(tagList.Tags) == 0 {
		t.Errorf("expected at least one tag from TagList")
	}

	reference := tagList.Tags[0].Name
	attributes, err := cr.TagGetAttributes(ctx, imageName, reference)
	if err != nil {
		t.Fatalf("received unexpected error: %v", err)
	}

	if attributes.ImageName != imageName {
		t.Errorf("expected attributes.ImageName to be %q but received: %s", imageName, attributes.ImageName)
	}

	toBoolPtr := func(v bool) *bool { return &v }
	err = cr.TagUpdateAttributes(ctx, imageName, reference, TagChangeableAttributes{ListEnabled: toBoolPtr(true)})
	if err != nil {
		t.Fatalf("received unexpected error: %v", err)
	}

	err = cr.TagDelete(ctx, imageName, "this-tag-does-not-exist")
	if err == nil {
		t.Fatal("expected error when running Delete()")
	}

	if !strings.Contains(err.Error(), "the specified tag does not exist") {
		t.Fatalf("expected error of Delete() to contain 'the specified tag does not exist' but received: %v", err)
	}
}

func testManifestE2E(t *testing.T, ctx context.Context, cr *ContainerRegistryClient, imageName string) {
	t.Helper()

	manifestList, err := cr.ManifestList(ctx, imageName, nil, nil, nil)
	if err != nil {
		t.Fatalf("received unexpected error: %v", err)
	}

	if len(manifestList.Manifests) == 0 {
		t.Errorf("expected at least one tag from ManifestList")
	}

	digest := manifestList.Manifests[0].Digest
	manifest, err := cr.ManifestGet(ctx, imageName, digest)
	if err != nil {
		t.Fatalf("received unexpected error: %v", err)
	}

	if len(manifest.Layers) == 0 {
		t.Errorf("expected at least one layer from manifest")
	}

	attributes, err := cr.ManifestGetAttributes(ctx, imageName, digest)
	if err != nil {
		t.Fatalf("received unexpected error: %v", err)
	}

	if attributes.ImageName != imageName {
		t.Errorf("expected attributes.ImageName to be %q but received: %s", imageName, attributes.ImageName)
	}

	toBoolPtr := func(v bool) *bool { return &v }
	err = cr.ManifestUpdateAttributes(ctx, imageName, digest, ManifestChangeableAttributes{ListEnabled: toBoolPtr(true)})
	if err != nil {
		t.Fatalf("received unexpected error: %v", err)
	}

	err = cr.ManifestDelete(ctx, imageName, "sha256:0000000000000000000000000000000000000000000000000000000000000000")
	if err == nil {
		t.Fatal("expected error when running Delete()")
	}

	if !strings.Contains(err.Error(), "manifest sha256:0000000000000000000000000000000000000000000000000000000000000000 is not found") {
		t.Fatalf("expected error of Delete() to contain 'manifest sha256:0000000000000000000000000000000000000000000000000000000000000000 is not found' but received: %v", err)
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
	t.Helper()

	path := r.URL.Path
	switch {
	case testMatchPath(t, path, "/oauth2/exchange"):
		h.refreshTokenHandler(t, w, r)
	case testMatchPath(t, path, "/oauth2/token"):
		h.accessTokenHandler(t, w, r)
	case testMatchPath(t, path, "/acr/v1/.*/_tags.*"):
		h.tagHandler(t, w, r)
	case testMatchPath(t, path, "/acr/v1/.*"):
		h.catalogHandler(t, w, r)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("unknown path: %s", path))) //nolint
		return
	}
}

func testMatchPath(t *testing.T, path, pattern string, vars ...interface{}) bool {
	t.Helper()

	regex := regexp.MustCompile("^" + pattern + "$")
	matches := regex.FindStringSubmatch(path)
	if len(matches) <= 0 {
		return false
	}

	for i, match := range matches[1:] {
		switch p := vars[i].(type) {
		case *string:
			*p = match
		case *int:
			n, err := strconv.Atoi(match)
			if err != nil {
				return false
			}
			*p = n
		default:
			panic("vars must be *string or *int")
		}
	}
	return true
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
