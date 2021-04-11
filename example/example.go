package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/msgraph"
	"github.com/manicminer/hamilton/environments"
)

var (
	tenantId           = os.Getenv("TENANT_ID")
	clientId           = os.Getenv("CLIENT_ID")
	clientSecret       = os.Getenv("CLIENT_SECRET")
)

func main() {
	ctx := context.Background()

	authConfig := &auth.Config{
		Environment:            environments.Global,
		TenantID:               tenantId,
		ClientID:               clientId,
		ClientSecret:           clientSecret,
		EnableClientSecretAuth: true,
	}

	authorizer, err := authConfig.NewAuthorizer(ctx, auth.MsGraph)
	if err != nil {
		log.Fatal(err)
	}

	client := msgraph.NewUsersClient(tenantId)
	client.BaseClient.Authorizer = authorizer

	users, _, err := client.List(ctx, "")
	if err != nil {
		log.Fatal(err)
	}
	if users == nil {
		log.Fatalln("bad API response, nil result received")
	}

	for _, user := range *users {
		fmt.Printf("%s: %s <%s>\n", *user.ID, *user.DisplayName, *user.UserPrincipalName)
	}
}
