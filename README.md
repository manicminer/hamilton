# Hamilton is a Go SDK for Microsoft Graph

This is a working Go client for the [Microsoft Graph API][ms-graph-docs]. It is actively maintained and has growing support for services and objects in Azure Active Directory.

## Example Usage

```go
package example

import (
	"context"
	"fmt"
	"log"

	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/clients"
)

const (
	tenantId = "00000000-0000-0000-0000-000000000000"
	clientId = "11111111-1111-1111-1111-111111111111"
	clientSecret = "My$3cR3t"
)

func main() {
	ctx := context.Background()

	authConfig := &auth.Config{
		TenantID:               tenantId,
		ClientID:               clientId,
		ClientSecret:           clientSecret,
		EnableClientSecretAuth: true,
	}

	authorizer, err := authConfig.NewAuthorizer(ctx)
	if err != nil {
		log.Fatal(err)
	}

	client := clients.NewUsersClient(tenantId)
	client.BaseClient.Authorizer = authorizer

	users, _, err := client.List(ctx, "")
	if err != nil {
		log.Fatal(err)
	}
	if users == nil {
		log.Fatalln("bad API response, nil result received")
	}

	for _, user := range *users {
		fmt.Printf("%s: %s <%s>\n", *user.ID, *user.DisplayName, *user.Mail)
	}
}
```

[ms-graph-docs]: https://docs.microsoft.com/en-us/graph/overview