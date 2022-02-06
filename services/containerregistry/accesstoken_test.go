package containerregistry

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/manicminer/hamilton/auth"
)

func testNewFakeAccessToken(t *testing.T) string {
	t.Helper()

	claims := AccessTokenClaims{
		JwtID:          "00000000-0000-0000-0000-000000000000",
		Subject:        "ze#foo@bar.com",
		NotBefore:      time.Now().Unix(),
		ExpirationTime: time.Now().Add(time.Hour).Unix(),
		IssuedAt:       time.Now().Unix(),
		Issuer:         "Azure Container Registry",
		Audience:       "foobar.azurecr.io",
		Version:        "1.0",
		Rid:            "00000000000000000000000000000000",
		GrantType:      "access_token",
		ApplicationID:  "00000000-0000-0000-0000-000000000000",
		Scopes:         AccessTokenScopes{},
		Roles:          []string{},
	}

	claimsBytes, err := json.Marshal(claims)
	if err != nil {
		t.Fatalf("received unexpected error: %v", err)
	}

	b64Claims := base64.RawURLEncoding.EncodeToString(claimsBytes)

	return fmt.Sprintf("foo.%s.bar", b64Claims)
}

func testExchangeAccessTokenSuccess(t *testing.T, authorizer auth.Authorizer, serverURL string, tenant string, httpClient *http.Client) {
	t.Helper()

	cr := NewContainerRegistryClient(authorizer, serverURL, tenant)
	cr.WithHttpClient(httpClient)
	ctx := context.Background()

	atScopes := AccessTokenScopes{
		{
			Type:    "registry",
			Name:    "catalog",
			Actions: []string{"*"},
		},
	}

	token, atClaims, err := cr.ExchangeAccessToken(ctx, "foobar", atScopes)
	if err != nil {
		t.Fatalf("Received unexpected error: %v", err)
	}

	if token == "" {
		t.Fatalf("received empty token")
	}

	if atClaims.JwtID != "00000000-0000-0000-0000-000000000000" {
		t.Fatalf("expected JwtID to be '00000000-0000-0000-0000-000000000000' but received: %s", atClaims.JwtID)
	}
}

func testExchangeAccessTokenFailure(t *testing.T, authorizer auth.Authorizer, serverURL string, tenant string, httpClient *http.Client, errContains string) {
	t.Helper()

	cr := NewContainerRegistryClient(authorizer, serverURL, tenant)
	cr.WithHttpClient(httpClient)
	ctx := context.Background()

	atScopes := AccessTokenScopes{
		{
			Type:    "registry",
			Name:    "catalog",
			Actions: []string{"*"},
		},
	}

	_, _, err := cr.ExchangeAccessToken(ctx, "foobar", atScopes)
	if err == nil {
		t.Fatalf("Expected to receive error but didn't")
	}

	if !strings.Contains(err.Error(), errContains) {
		t.Fatalf("Expected to receive error containing %q but received: %v", errContains, err)
	}
}

func (h *testACRHandler) accessTokenHandler(t *testing.T, w http.ResponseWriter, r *http.Request) {
	t.Helper()

	err := h.validateExchangeAccessTokenRequest(t, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error())) //nolint
		return
	}

	response := struct {
		AccessToken string `json:"access_token"`
	}{
		AccessToken: testNewFakeAccessToken(t),
	}

	json.NewEncoder(w).Encode(response) //nolint
}

func (h *testACRHandler) validateExchangeAccessTokenRequest(t *testing.T, r *http.Request) error {
	t.Helper()

	path := r.URL.Path
	if path != "/oauth2/token" {
		return fmt.Errorf("expected path '/oauth2/token', received path: %s", path)
	}

	query := r.URL.Query()
	if len(query) > 0 {
		return fmt.Errorf("expected empty query, received: %s", query)
	}

	contentType := r.Header.Get("Content-Type")
	if contentType != "application/x-www-form-urlencoded" {
		return fmt.Errorf("expected Content-Type 'application/x-www-form-urlencoded', received: %s", contentType)
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("received unexpected error reading body: %v", err)
	}
	defer r.Body.Close()

	reqData, err := url.ParseQuery(string(bodyBytes))
	if err != nil {
		return fmt.Errorf("received unexpected error parsing bodyBytes: %v", err)
	}

	grantType := reqData.Get("grant_type")
	if grantType != "refresh_token" {
		return fmt.Errorf("expected req body grant_type to be 'refresh_token', received: %s", grantType)
	}

	expectedService, err := url.Parse(h.serverURL)
	if err != nil {
		return fmt.Errorf("received unexpected error parsing serverURL: %v", err)
	}
	service := reqData.Get("service")
	if service != expectedService.Hostname() {
		return fmt.Errorf("expected req body service to be %q, received: %s", expectedService.Hostname(), service)
	}

	refreshToken := reqData.Get("refresh_token")
	if refreshToken == "" {
		return fmt.Errorf("received empty refresh_token")
	}

	scope := reqData.Get("scope")
	if scope != "registry:catalog:*" {
		return fmt.Errorf("expected req body scope to be 'registry:catalog:*', received: %s", scope)
	}

	if h.fakeError != nil {
		return h.fakeError
	}

	return nil
}

func TestDecodeAccessTokenWithoutValidation(t *testing.T) {
	atClaims, err := decodeAccessTokenWithoutValidation(testNewFakeAccessToken(t))
	if err != nil {
		t.Fatalf("received unexpected error: %v", err)
	}

	if atClaims.JwtID != "00000000-0000-0000-0000-000000000000" {
		t.Fatalf("expected JwtID to be '00000000-0000-0000-0000-000000000000' but received: %s", atClaims.JwtID)
	}
}
