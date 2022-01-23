package autorest

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/manicminer/hamilton/auth"
)

// NewAutorestAuthorizerWrapper returns an Authorizer that sources tokens from a supplied autorest.BearerAuthorizer
func NewAutorestAuthorizerWrapper(autorestAuthorizer autorest.Authorizer) (auth.Authorizer, error) {
	return &AutorestAuthorizerWrapper{authorizer: autorestAuthorizer}, nil
}
