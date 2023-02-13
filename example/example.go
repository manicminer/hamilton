package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/go-azure-sdk/sdk/auth"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
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
	env := environments.AzurePublic()

	credentials := auth.Credentials{
		Environment:  *env,
		TenantID:     tenantId,
		ClientID:     clientId,
		ClientSecret: clientSecret,

		EnableAuthenticatingUsingClientSecret: true,
	}

	authorizer, err := auth.NewAuthorizerFromCredentials(ctx, credentials, env.MicrosoftGraph)
	if err != nil {
		log.Fatal(err)
	}

	client := msgraph.NewUsersClient(tenantId)
	client.BaseClient.Authorizer = authorizer

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
		fmt.Printf("%s: %s <%s>\n", *user.ID(), *user.DisplayName, *user.UserPrincipalName)
	}
}
