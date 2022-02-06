package containerregistry

import (
	"context"
	"fmt"
	"time"

	"github.com/manicminer/hamilton/auth"
	"golang.org/x/oauth2"
)

// NewAuthorizer returns an auth.Authorizer based on the ContainerRegistryClient
func (cr *ContainerRegistryClient) NewAuthorizer() (auth.Authorizer, error) {
	return cr, nil
}

// Token returns an *oauth2.Token or an error
func (cr *ContainerRegistryClient) Token() (*oauth2.Token, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	return cr.getOAuth2Token(ctx)
}

// AuxiliaryTokens returns a slice of *oauth2.Tokens or an error.
// Not supported with ContainerRegistryClient and will always return an error.
func (cr *ContainerRegistryClient) AuxiliaryTokens() ([]*oauth2.Token, error) {
	return nil, fmt.Errorf("AuxiliaryTokens not supported with ContainerRegistryClient")
}

func (cr *ContainerRegistryClient) getOAuth2Token(ctx context.Context) (*oauth2.Token, error) {
	_, err := cr.getAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	return &oauth2.Token{
		AccessToken:  cr.accessToken,
		RefreshToken: cr.refreshToken,
		TokenType:    "Bearer",
		Expiry:       time.Unix(cr.accessTokenClaims.ExpirationTime, 0),
	}, nil
}
