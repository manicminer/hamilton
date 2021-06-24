# Hamilton is a Go SDK for Microsoft Graph

This is a working Go client for the [Microsoft Graph API][ms-graph-docs]. It is actively maintained and has growing
support for services and objects in Azure Active Directory.

## Documentation

See [pkg.go.dev](https://pkg.go.dev/github.com/manicminer/hamilton).

## Example Usage

```go
package main

import (
	"context"
	"fmt"
	"log"
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
```

## Contributing

Contributions are welcomed! Please note that clients must have tests that cover all methods where feasible.

Please raise a pull request on GitHub to submit contributions. Bug reports and feature requests are happily received.

## Testing

Testing requires an Azure AD tenant and real credentials. You can authenticate with any supported method for the client
tests, and the auth tests are split by authentication method.

Note that each client generally has a single test that exercises all methods. This is to help ensure that test objects
are cleaned up where possible. Where tests fail, often objects will be left behind and should be cleaned up manually.

It's recommended to use an isolated tenant for testing and _not_ a production tenant.

To run all the tests:
```shell
$ make test
```

[ms-graph-docs]: https://docs.microsoft.com/en-us/graph/overview
