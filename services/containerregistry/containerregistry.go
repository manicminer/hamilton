package containerregistry

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/manicminer/hamilton/auth"
)

var (
	defaultExpiryDelta       = 10 * time.Second
	defaultAccessTokenScopes = AccessTokenScopes{
		{
			Type:    "registry",
			Name:    "catalog",
			Actions: []string{"*"},
		},
		{
			Type:    "repository",
			Name:    "*",
			Actions: []string{"*"},
		},
	}
)

// ContainerRegistryClient handles communication with Azure Container Registry
type ContainerRegistryClient struct {
	mu                 sync.Mutex
	authorizer         auth.Authorizer
	httpClient         *http.Client
	serverURL          string
	tenantID           string
	refreshToken       string
	refreshTokenClaims RefreshTokenClaims
	accessToken        string
	accessTokenClaims  AccessTokenClaims
	accessTokenScopes  AccessTokenScopes
	expiryDelta        time.Duration
}

// NewContainerRegistryClient returns a ContainerRegistryClient
func NewContainerRegistryClient(authorizer auth.Authorizer, serverURL string, tenantID string) *ContainerRegistryClient {

	return &ContainerRegistryClient{
		authorizer:        authorizer,
		httpClient:        http.DefaultClient,
		serverURL:         serverURL,
		tenantID:          tenantID,
		accessTokenScopes: defaultAccessTokenScopes,
		expiryDelta:       defaultExpiryDelta,
	}
}

// WithHttpClient replaces what http client is used for communication to Azure Container Registry
func (cr *ContainerRegistryClient) WithHttpClient(httpClient *http.Client) {
	cr.httpClient = httpClient
}

// WithAccessTokenScopes replaces the default access token scopes with the ones provided.
// This is used with the ContainerRegistryClient when communicating with the different APIs.
func (cr *ContainerRegistryClient) WithAccessTokenScopes(accessTokenScopes AccessTokenScopes) {
	cr.accessTokenScopes = accessTokenScopes
}

// WithAccessTokenScopes replaces the default access token scopes with the ones provided.
func (cr *ContainerRegistryClient) WithExpiryDelta(expiryDelta time.Duration) {
	cr.expiryDelta = expiryDelta
}

func (cr *ContainerRegistryClient) getAccessToken(ctx context.Context) (string, error) {
	cr.mu.Lock()
	at := cr.accessToken
	atClaims := cr.accessTokenClaims
	cr.mu.Unlock()

	if isTokenValid(at, time.Now(), atClaims.ExpirationTime, cr.expiryDelta) {
		return at, nil
	}

	cr.mu.Lock()
	rt := cr.refreshToken
	rtClaims := cr.refreshTokenClaims
	cr.mu.Unlock()

	if !isTokenValid(rt, time.Now(), rtClaims.ExpirationTime, cr.expiryDelta) {
		rt, rtClaims, err := cr.ExchangeRefreshToken(ctx)
		if err != nil {
			return "", err
		}

		cr.mu.Lock()
		cr.refreshToken = rt
		cr.refreshTokenClaims = rtClaims
		cr.mu.Unlock()
	}

	at, atClaims, err := cr.ExchangeAccessToken(ctx, cr.refreshToken, cr.accessTokenScopes)
	if err != nil {
		return "", err
	}

	cr.mu.Lock()
	cr.accessToken = at
	cr.accessTokenClaims = atClaims
	cr.mu.Unlock()

	return at, nil
}

func (cr *ContainerRegistryClient) setAuthorizationHeader(ctx context.Context, req *http.Request) error {
	at, err := cr.getAccessToken(ctx)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", at))

	return nil
}

func (cr *ContainerRegistryClient) getBaseURL() (*url.URL, error) {
	return parseService(cr.serverURL)
}

func parseService(serverURL string) (*url.URL, error) {
	scheme := "https://"
	if strings.HasPrefix(serverURL, "https://") {
		scheme = ""
	}

	serviceURL, err := url.Parse(fmt.Sprintf("%s%s", scheme, serverURL))
	if err != nil {
		return nil, fmt.Errorf("failed to parse server URL - %w", err)
	}

	return serviceURL, nil
}

func isTokenValid(token string, currentTime time.Time, expirationTime int64, expiryDelta time.Duration) bool {
	if token == "" {
		return false
	}

	expiry := time.Unix(expirationTime, 0)
	return expiry.Round(0).After(currentTime.Add(expiryDelta))
}
