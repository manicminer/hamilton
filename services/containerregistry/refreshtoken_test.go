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

func testNewFakeRefreshToken(t *testing.T) string {
	t.Helper()

	claims := RefreshTokenClaims{
		JwtID:          "00000000-0000-0000-0000-000000000000",
		Subject:        "ze#foo@bar.com",
		NotBefore:      time.Now().Unix(),
		ExpirationTime: time.Now().Add(time.Hour).Unix(),
		IssuedAt:       time.Now().Unix(),
		Issuer:         "Azure Container Registry",
		Audience:       "foobar.azurecr.io",
		Version:        "1.0",
		Rid:            "00000000000000000000000000000000",
		GrantType:      "refresh_token",
		ApplicationID:  "00000000-0000-0000-0000-000000000000",
		Permissions:    RefreshTokenClaimsPermissions{},
		Roles:          []string{},
	}

	claimsBytes, err := json.Marshal(claims)
	if err != nil {
		t.Fatalf("received unexpected error: %v", err)
	}

	b64Claims := base64.RawURLEncoding.EncodeToString(claimsBytes)

	return fmt.Sprintf("foo.%s.bar", b64Claims)
}

func testExchangeRefreshTokenSuccess(t *testing.T, authorizer auth.Authorizer, serverURL string, tenant string, httpClient *http.Client) {
	t.Helper()

	cr := NewContainerRegistryClient(authorizer, serverURL, tenant)
	cr.WithHttpClient(httpClient)
	ctx := context.Background()
	token, rtClaims, err := cr.ExchangeRefreshToken(ctx)
	if err != nil {
		t.Fatalf("Received unexpected error: %v", err)
	}

	if token == "" {
		t.Fatalf("received empty token")
	}

	if rtClaims.JwtID != "00000000-0000-0000-0000-000000000000" {
		t.Fatalf("expected JwtID to be '00000000-0000-0000-0000-000000000000' but received: %s", rtClaims.JwtID)
	}
}

func testExchangeRefreshTokenFailure(t *testing.T, authorizer auth.Authorizer, serverURL string, tenant string, httpClient *http.Client, errContains string) {
	t.Helper()

	cr := NewContainerRegistryClient(authorizer, serverURL, tenant)
	cr.WithHttpClient(httpClient)
	ctx := context.Background()
	_, _, err := cr.ExchangeRefreshToken(ctx)
	if err == nil {
		t.Fatalf("Expected to receive error but didn't")
	}

	if !strings.Contains(err.Error(), errContains) {
		t.Fatalf("Expected to receive error containing %q but received: %v", errContains, err)
	}
}

func (h *testACRHandler) refreshTokenHandler(t *testing.T, w http.ResponseWriter, r *http.Request) {
	t.Helper()

	err := h.validateExchangeRefreshTokenRequest(t, r)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error())) //nolint
		return
	}

	response := struct {
		RefreshToken string `json:"refresh_token"`
	}{
		RefreshToken: testNewFakeRefreshToken(t),
	}

	json.NewEncoder(w).Encode(response) //nolint
}

func (h *testACRHandler) validateExchangeRefreshTokenRequest(t *testing.T, r *http.Request) error {
	t.Helper()

	if r.Method != http.MethodPost {
		return fmt.Errorf("expected method to be POST, received: %s", r.Method)
	}

	path := r.URL.Path
	if path != "/oauth2/exchange" {
		return fmt.Errorf("expected path '/oauth2/exchange', received path: %s", path)
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
	if grantType != "access_token" {
		return fmt.Errorf("expected req body grant_type to be 'access_token', received: %s", grantType)
	}

	expectedService, err := url.Parse(h.serverURL)
	if err != nil {
		return fmt.Errorf("received unexpected error parsing serverURL: %v", err)
	}
	service := reqData.Get("service")
	if service != expectedService.Hostname() {
		return fmt.Errorf("expected req body service to be %q, received: %s", expectedService.Hostname(), service)
	}

	accessToken := reqData.Get("access_token")
	if accessToken != "foo" {
		return fmt.Errorf("expected req body access_token to be 'foo', received: %s", accessToken)
	}

	tenant := reqData.Get("tenant")
	if tenant != h.expectedTenant {
		return fmt.Errorf("expected req body tenant to be %q, received: %s", h.expectedTenant, tenant)
	}

	if h.fakeError != nil {
		return h.fakeError
	}

	return nil
}

func TestDecodeRefreshTokenWithoutValidation(t *testing.T) {
	rtClaims, err := decodeRefreshTokenWithoutValidation(testNewFakeRefreshToken(t))
	if err != nil {
		t.Fatalf("received unexpected error: %v", err)
	}

	if rtClaims.JwtID != "00000000-0000-0000-0000-000000000000" {
		t.Fatalf("expected JwtID to be '00000000-0000-0000-0000-000000000000' but received: %s", rtClaims.JwtID)
	}
}
