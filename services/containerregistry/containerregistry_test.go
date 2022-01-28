package containerregistry

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/manicminer/hamilton/auth"
	"golang.org/x/oauth2"
)

func TestContainerRegistryClient(t *testing.T) {
	fa := testNewFakeAuthorizer(t)
	h := testNewACRHandler(t)
	httpServer := httptest.NewTLSServer(h.handler(t))
	h.serverURL = httpServer.URL
	h.expectedTenant = ""
	testContainerRegistryClient(t, fa, httpServer.URL, "", httpServer.Client())

	h.expectedTenant = "ze-tenant"
	testContainerRegistryClient(t, fa, httpServer.URL, "ze-tenant", httpServer.Client())

	h.expectedTenant = "ze-tenant"
	h.fakeError = fmt.Errorf("ze-fake-error")
	testContainerRegistryClientFailure(t, fa, httpServer.URL, "ze-tenant", httpServer.Client(), "ze-fake-error")
}

func testContainerRegistryClient(t *testing.T, authorizer auth.Authorizer, serverURL string, tenant string, httpClient *http.Client) {
	t.Helper()

	cr := NewContainerRegistryClient(authorizer, serverURL, tenant)
	cr.WithHttpClient(httpClient)
	ctx := context.Background()
	token, err := cr.ExchangeToken(ctx)
	if err != nil {
		t.Fatalf("Received unexpected error: %v", err)
	}

	if token != "foobar" {
		t.Fatalf("Expected token 'foobar', received: %s", token)
	}
}

func testContainerRegistryClientFailure(t *testing.T, authorizer auth.Authorizer, serverURL string, tenant string, httpClient *http.Client, errContains string) {
	t.Helper()

	cr := NewContainerRegistryClient(authorizer, serverURL, tenant)
	cr.WithHttpClient(httpClient)
	ctx := context.Background()
	_, err := cr.ExchangeToken(ctx)
	if err == nil {
		t.Fatalf("Expected to receive error but didn't")
	}

	if !strings.Contains(err.Error(), errContains) {
		t.Fatalf("Expected to receive error containing %q but received: %v", errContains, err)
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
	response       struct {
		RefreshToken string `json:"refresh_token"`
	}
}

func testNewACRHandler(t *testing.T) *testACRHandler {
	t.Helper()

	return &testACRHandler{}
}

func (h *testACRHandler) handler(t *testing.T) http.HandlerFunc {
	t.Helper()

	return func(w http.ResponseWriter, r *http.Request) {
		accessToken, err := h.validateRequest(t, r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		h.response.RefreshToken = fmt.Sprintf("%sbar", accessToken)

		json.NewEncoder(w).Encode(h.response)
	}
}

func (h *testACRHandler) validateRequest(t *testing.T, r *http.Request) (string, error) {
	t.Helper()

	path := r.URL.Path
	if path != "/oauth2/exchange" {
		return "", fmt.Errorf("expected path '/oauth2/exchange', received path: %s", path)
	}

	query := r.URL.Query()
	if len(query) > 0 {
		return "", fmt.Errorf("expected empty query, received: %s", query)
	}

	contentType := r.Header.Get("Content-Type")
	if contentType != "application/x-www-form-urlencoded" {
		return "", fmt.Errorf("expected Content-Type 'application/x-www-form-urlencoded', received: %s", contentType)
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return "", fmt.Errorf("received unexpected error reading body: %v", err)
	}
	defer r.Body.Close()

	reqData, err := url.ParseQuery(string(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("received unexpected error parsing bodyBytes: %v", err)
	}

	grantType := reqData.Get("grant_type")
	if grantType != "access_token" {
		return "", fmt.Errorf("expected req body grant_type to be 'access_token', received: %s", grantType)
	}

	expectedService, err := url.Parse(h.serverURL)
	if err != nil {
		return "", fmt.Errorf("received unexpected error parsing serverURL: %v", err)
	}
	service := reqData.Get("service")
	if service != expectedService.Hostname() {
		return "", fmt.Errorf("expected req body service to be %q, received: %s", expectedService.Hostname(), service)
	}

	accessToken := reqData.Get("access_token")
	if accessToken != "foo" {
		return "", fmt.Errorf("expected req body access_token to be 'foo', received: %s", accessToken)
	}

	tenant := reqData.Get("tenant")
	if tenant != h.expectedTenant {
		return "", fmt.Errorf("expected req body tenant to be %q, received: %s", h.expectedTenant, tenant)
	}

	if h.fakeError != nil {
		return "", h.fakeError
	}

	return accessToken, nil
}
