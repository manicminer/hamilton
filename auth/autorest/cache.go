package autorest

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Azure/go-autorest/autorest"
	"github.com/manicminer/hamilton/auth"
	"golang.org/x/oauth2"
)

// CachedAuthorizer caches a token until it expires, then acquires a new token from Source
type CachedAutorestAuthorizer struct {
	cache *auth.CachedAuthorizer
}

// Token returns the current token if it's still valid, else will acquire a new token
func (c *CachedAutorestAuthorizer) Token() (*oauth2.Token, error) {
	return c.cache.Token()
}

// AuxiliaryTokens returns additional tokens for auxiliary tenant IDs, for use in multi-tenant scenarios
func (c *CachedAutorestAuthorizer) AuxiliaryTokens() ([]*oauth2.Token, error) {
	return c.cache.AuxiliaryTokens()
}

// WithAuthorization implements the autorest.Authorizer interface
func (c *CachedAutorestAuthorizer) WithAuthorization() autorest.PrepareDecorator {
	return func(p autorest.Preparer) autorest.Preparer {
		return autorest.PreparerFunc(func(req *http.Request) (*http.Request, error) {
			var err error
			req, err = p.Prepare(req)
			if err == nil {
				token, err := c.Token()
				if err != nil {
					return nil, err
				}

				req, err = autorest.Prepare(req, autorest.WithHeader("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken)))
				if err != nil {
					return req, err
				}

				auxTokens, err := c.AuxiliaryTokens()
				if err != nil {
					return req, err
				}

				auxTokenList := make([]string, 0)
				for _, a := range auxTokens {
					if a != nil && a.AccessToken != "" {
						auxTokenList = append(auxTokenList, fmt.Sprintf("%s %s", a.TokenType, a.AccessToken))
					}
				}

				return autorest.Prepare(req, autorest.WithHeader("x-ms-authorization-auxiliary", strings.Join(auxTokenList, ", ")))
			}

			return req, err
		})
	}
}

// BearerAuthorizerCallback is a helper that returns an *autorest.BearerAuthorizerCallback for use in data plane API clients in the Azure SDK
func (c *CachedAutorestAuthorizer) BearerAuthorizerCallback() *autorest.BearerAuthorizerCallback {
	return autorest.NewBearerAuthorizerCallback(nil, func(_, resource string) (*autorest.BearerAuthorizer, error) {
		token, err := c.Token()
		if err != nil {
			return nil, fmt.Errorf("obtaining token: %v", err)
		}

		return autorest.NewBearerAuthorizer(&servicePrincipalTokenWrapper{
			tokenType:  "Bearer",
			tokenValue: token.AccessToken,
		}), nil
	})
}

// NewCachedAutorestAuthorizer returns an Authorizer that caches an access token for the duration of its validity.
// If the cached token expires, a new one is acquired and cached.
func NewCachedAutorestAuthorizer(cache *auth.CachedAuthorizer) auth.Authorizer {
	return &CachedAutorestAuthorizer{cache}
}
