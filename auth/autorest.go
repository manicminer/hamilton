package auth

import (
	"context"
	"fmt"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"golang.org/x/oauth2"
)

// AutorestAuthorizerWrapper is an Authorizer which sources tokens from an autorest.Authorizer
// Currently supports only:
// - autorest.BearerAuthorizer
// - autorest.MultiTenantBearerAuthorizer
type AutorestAuthorizerWrapper struct {
	authorizer interface{}
}

func (a *AutorestAuthorizerWrapper) tokenProviders() (tokenProviders []adal.OAuthTokenProvider, err error) {
	if authorizer, ok := a.authorizer.(*autorest.BearerAuthorizer); ok && authorizer != nil {
		tokenProviders = append(tokenProviders, authorizer.TokenProvider())
	} else if authorizer, ok := a.authorizer.(*autorest.MultiTenantBearerAuthorizer); ok && authorizer != nil {
		if multiTokenProvider := authorizer.TokenProvider(); multiTokenProvider != nil {
			if m, ok := multiTokenProvider.(*adal.MultiTenantServicePrincipalToken); ok && m != nil {
				tokenProviders = append(tokenProviders, m.PrimaryToken)
				for _, aux := range m.AuxiliaryTokens {
					tokenProviders = append(tokenProviders, aux)
				}
			}
		}
	}

	for _, tokenProvider := range tokenProviders {
		if refresher, ok := tokenProvider.(adal.Refresher); ok {
			if err = refresher.EnsureFresh(); err != nil {
				return
			}
		} else if refresher, ok := tokenProvider.(adal.RefresherWithContext); ok {
			if err = refresher.EnsureFreshWithContext(context.Background()); err != nil {
				return
			}
		}
	}

	return
}

// Token returns an access token using an autorest.BearerAuthorizer struct
func (a *AutorestAuthorizerWrapper) Token() (*oauth2.Token, error) {
	tokenProviders, err := a.tokenProviders()
	if err != nil {
		return nil, err
	}
	if len(tokenProviders) == 0 {
		return nil, fmt.Errorf("no token providers returned")
	}

	var adalToken adal.Token
	if spToken, ok := tokenProviders[0].(*adal.ServicePrincipalToken); ok && spToken != nil {
		adalToken = spToken.Token()
	}
	if adalToken.AccessToken == "" {
		return nil, fmt.Errorf("could not obtain access token from token provider")
	}

	return &oauth2.Token{
		AccessToken:  adalToken.AccessToken,
		TokenType:    adalToken.Type,
		RefreshToken: adalToken.RefreshToken,
		Expiry:       adalToken.Expires(),
	}, nil
}

// AuxiliaryTokens returns additional tokens for auxiliary tenant IDs, sourced from an
// autorest.MultiTenantBearerAuthorizer, for use in multi-tenant scenarios
func (a *AutorestAuthorizerWrapper) AuxiliaryTokens() ([]*oauth2.Token, error) {
	tokenProviders, err := a.tokenProviders()
	if err != nil {
		return nil, err
	}

	var auxTokens []*oauth2.Token
	for i := 1; i < len(tokenProviders); i++ {
		var adalToken adal.Token

		if spToken, ok := tokenProviders[i].(*adal.ServicePrincipalToken); ok && spToken != nil {
			adalToken = spToken.Token()
		}

		if adalToken.AccessToken == "" {
			return nil, fmt.Errorf("could not obtain access token from token providers")
		}

		auxTokens = append(auxTokens, &oauth2.Token{
			AccessToken:  adalToken.AccessToken,
			TokenType:    adalToken.Type,
			RefreshToken: adalToken.RefreshToken,
			Expiry:       adalToken.Expires(),
		})
	}

	return auxTokens, nil
}
