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
}

// NewContainerRegistryClient returns a ContainerRegistryClient
func NewContainerRegistryClient(authorizer auth.Authorizer, serverURL string, tenantID string) *ContainerRegistryClient {
	httpClient := http.DefaultClient
	return &ContainerRegistryClient{
		authorizer: authorizer,
		httpClient: httpClient,
		serverURL:  serverURL,
		tenantID:   tenantID,
	}
}

// WithHttpClient replaces what http client is used for communication to Azure Container Registry
func (cr *ContainerRegistryClient) WithHttpClient(httpClient *http.Client) {
	cr.httpClient = httpClient
}

func (cr *ContainerRegistryClient) getAccessToken(ctx context.Context) (string, error) {
	cr.mu.Lock()
	at := cr.accessToken
	atClaims := cr.accessTokenClaims
	cr.mu.Unlock()

	atExpiry := time.Unix(atClaims.ExpirationTime, 0)
	atValid := at != "" && atExpiry.After(time.Now().Add(time.Minute))
	if atValid {
		return at, nil
	}

	cr.mu.Lock()
	rt := cr.refreshToken
	rtClaims := cr.refreshTokenClaims
	cr.mu.Unlock()

	rtExpiry := time.Unix(rtClaims.ExpirationTime, 0)
	rtValid := rt != "" && !rtExpiry.IsZero() && rtExpiry.After(time.Now().Add(time.Minute))
	if !rtValid {
		rt, rtClaims, err := cr.ExchangeRefreshToken(ctx)
		if err != nil {
			return "", err
		}

		cr.mu.Lock()
		cr.refreshToken = rt
		cr.refreshTokenClaims = rtClaims
		cr.mu.Unlock()
	}

	scopes := AccessTokenScopes{
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

	at, atClaims, err := cr.ExchangeAccessToken(ctx, cr.refreshToken, scopes)
	if err != nil {
		return "", nil
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
