package hamilton

import (
	"context"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"golang.org/x/oauth2/microsoft"
)

type Authorizer interface {
	Token() (*oauth2.Token, error)
}

func NewClientSecretAuthorizer(ctx context.Context, clientId, clientSecret, tenantId string) oauth2.TokenSource {
	conf := clientcredentials.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		Scopes:       []string{"https://graph.microsoft.com/.default"},
		TokenURL:     microsoft.AzureADEndpoint(tenantId).TokenURL,
	}
	return conf.TokenSource(ctx)
}