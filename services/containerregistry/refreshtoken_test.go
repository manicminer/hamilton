package containerregistry

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/manicminer/hamilton/auth"
)

const (
	testFakeRefreshToken = "foo.eyJqdGkiOiIwMDAwMDAwMC0wMDAwLTAwMDAtMDAwMC0wMDAwMDAwMDAwMDAiLCJzdWIiOiJ6ZSNmb29AYmFyLmNvbSIsIm5iZiI6MTY0MzQ4OTgxNiwiZXhwIjoxNjQzNTAxNTE2LCJpYXQiOjE2NDM0ODk4MTYsImlzcyI6IkF6dXJlIENvbnRhaW5lciBSZWdpc3RyeSIsImF1ZCI6ImZvb2Jhci5henVyZWNyLmlvIiwidmVyc2lvbiI6IjEuMCIsInJpZCI6IjAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwIiwiZ3JhbnRfdHlwZSI6InJlZnJlc2hfdG9rZW4iLCJhcHBpZCI6IjAwMDAwMDAwLTAwMDAtMDAwMC0wMDAwLTAwMDAwMDAwMDAwMCIsInBlcm1pc3Npb25zIjp7IkFjdGlvbnMiOlsicmVhZCIsIndyaXRlIiwiZGVsZXRlIl0sIk5vdEFjdGlvbnMiOltdfSwicm9sZXMiOltdfQ.bar"
)

func testExchangeRefreshTokenSuccess(t *testing.T, authorizer auth.Authorizer, serverURL string, tenant string, httpClient *http.Client) {
	t.Helper()

	cr := NewContainerRegistryClient(authorizer, serverURL, tenant)
	cr.WithHttpClient(httpClient)
	ctx := context.Background()
	token, rtClaims, err := cr.ExchangeRefreshToken(ctx)
	if err != nil {
		t.Fatalf("Received unexpected error: %v", err)
	}

	if token != testFakeRefreshToken {
		t.Fatalf("Expected token %q, received: %s", testFakeRefreshToken, token)
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
	err := h.validateExchangeRefreshTokenRequest(t, r)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error())) //nolint
		return
	}

	response := struct {
		RefreshToken string `json:"refresh_token"`
	}{
		RefreshToken: testFakeRefreshToken,
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
	rtClaims, err := decodeRefreshTokenWithoutValidation(testFakeRefreshToken)
	if err != nil {
		t.Fatalf("received unexpected error: %v", err)
	}

	if rtClaims.JwtID != "00000000-0000-0000-0000-000000000000" {
		t.Fatalf("expected JwtID to be '00000000-0000-0000-0000-000000000000' but received: %s", rtClaims.JwtID)
	}
}
