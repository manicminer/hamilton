package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/environments"
	"github.com/manicminer/hamilton/msgraph"
	"github.com/manicminer/hamilton/odata"
)

var (
	tenantId     = os.Getenv("TENANT_ID")
	clientId     = os.Getenv("CLIENT_ID")
	clientSecret = os.Getenv("CLIENT_SECRET")
)

func main() {
	ctx := context.Background()

	environment := environments.Global

	authConfig := &auth.Config{
		Environment:            environment,
		TenantID:               tenantId,
		ClientID:               clientId,
		ClientSecret:           clientSecret,
		EnableClientSecretAuth: true,
	}

	authorizer, err := authConfig.NewAuthorizer(ctx, environment.MsGraph)
	if err != nil {
		log.Fatal(err)
	}

	requestLogger := func(req *http.Request) (*http.Request, error) {
		if req != nil {
			dmp, err := httputil.DumpRequestOut(req, true)
			if err == nil {
				log.Printf("REQUEST: %s", dmp)
			}
		}
		return req, nil
	}

	responseLogger := func(req *http.Request, resp *http.Response) (*http.Response, error) {
		if resp != nil {
			dmp, err := httputil.DumpResponse(resp, true)
			if err == nil {
				log.Printf("RESPONSE: %s", dmp)
			}
		}
		return resp, nil
	}

	client := msgraph.NewUsersClient(tenantId)
	client.BaseClient.Authorizer = authorizer
	client.BaseClient.RequestMiddlewares = &[]msgraph.RequestMiddleware{requestLogger}
	client.BaseClient.ResponseMiddlewares = &[]msgraph.ResponseMiddleware{responseLogger}

	users, _, err := client.List(ctx, odata.Query{})
	if err != nil {
		log.Println(err)
		return
	}
	if users == nil {
		log.Println("bad API response, nil result received")
		return
	}
	for _, user := range *users {
		fmt.Printf("%s: %s <%s>\n", *user.ID, *user.DisplayName, *user.UserPrincipalName)
	}
}
