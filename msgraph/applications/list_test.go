package applications

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-sdk/sdk/auth"
	"github.com/hashicorp/go-azure-sdk/sdk/client/msgraph"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

func TestList(t *testing.T) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Minute))
	defer cancel()
	env := environments.AzurePublic()
	authorizer, err := auth.NewAuthorizerFromCredentials(ctx, auth.Credentials{
		Environment:                       *env,
		EnableAuthenticatingUsingAzureCLI: true,
	}, env.MicrosoftGraph)
	if err != nil {
		t.Fatal(err)
	}
	client, err := NewApplicationsClient(env.MicrosoftGraph, msgraph.Version10, os.Getenv("TENANT_ID"))
	client.BaseClient.Client.Authorizer = authorizer

	resp, err := client.List(ctx, ListOptions{OData: odata.Query{Top: 10}})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%#v\n", resp)

}
