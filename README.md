# Hamilton is a Go SDK for Microsoft Graph

This is a working Go client for the [Microsoft Graph API][ms-graph-docs]. It is actively maintained and has growing
support for services and objects in Azure Active Directory.

## Documentation

See [pkg.go.dev](https://pkg.go.dev/github.com/manicminer/hamilton).

## Features

- Automatic retries for failed requests and handling of eventual consistency on writes due to propagation delays
- Automatic paging of results
- Native model structs for marshaling and unmarshaling
- Support for national clouds including US Government (L4 and L5) and China
- Support for both the v1.0 and beta API endpoints
- Ability to inject middleware functions for logging etc
- OData parsing in API responses and support for OData queries such as filters, sorting, searching, expand and select
- Authentication now uses [github.com/hashicorp/go-azure-sdk/sdk/auth](https://github.com/hashicorp/go-azure-sdk/tree/main/sdk/auth)

## Getting Started

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/go-azure-sdk/sdk/auth"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/manicminer/hamilton/msgraph"
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

	client := msgraph.NewUsersClient()
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
```

## Configure retry limit for all failed requests

```go
client := msgraph.NewUsersClient(tenantId)
client.BaseClient.Authorizer = authorizer
client.BaseClient.RetryableClient.RetryMax = 8
```

## Disable eventual consistency handling

_Note: this does **not** disable auto-retries for failed requests (e.g. HTTP 429 or 500 responses)_

```go
client := msgraph.NewUsersClient(tenantId)
client.BaseClient.Authorizer = authorizer
client.BaseClient.DisableRetries = true
```

## Log requests and responses

```go
requestLogger := func(req *http.Request) (*http.Request, error) {
	if req != nil {
		if dump, err := httputil.DumpRequestOut(req, true); err == nil {
			log.Printf("%s\n", dump)
		}
	}
	return req, nil
}

responseLogger := func(req *http.Request, resp *http.Response) (*http.Response, error) {
	if resp != nil {
		if dump, err := httputil.DumpResponse(resp, true); err == nil {
			log.Printf("%s\n", dump)
		}
	}
	return resp, nil
}

client := msgraph.NewUsersClient(tenantId)
client.BaseClient.Authorizer = authorizer
client.BaseClient.DisableRetries = true
client.BaseClient.RequestMiddlewares = &[]msgraph.RequestMiddleware{requestLogger}
client.BaseClient.ResponseMiddlewares = &[]msgraph.ResponseMiddleware{responseLogger}
```

## Contributing

Contributions are welcomed! Please note that clients must have tests that cover all methods where feasible.

Please raise a pull request [on GitHub][gh-project] to submit contributions. Bug reports and feature requests are happily received.

## Testing

Testing requires an Azure AD tenant and real credentials. Note that some tests require an Azure AD Premium P2 license and/or an Office 365 license.
You can authenticate with any supported method for the client tests, and the auth tests are split by authentication method.

> ℹ️ You can sign up for the [Microsoft 365 Developer Program](https://developer.microsoft.com/en-us/microsoft-365/dev-program) which offers a Microsoft 365 E5 subscription for 25 users, at no cost for development purposes.

Note that each client generally has a single test that exercises all methods. This is to help ensure that test objects
are cleaned up where possible. Where tests fail, often objects will be left behind and should be cleaned up separately.
The [test-cleanup](https://github.com/manicminer/hamilton/tree/main/internal/cmd/test-cleanup) command can be used to
delete leftover test objects in the event of test failure.

It's recommended to use an isolated tenant for testing and _not_ a production tenant.

To run all the tests:
```shell
$ make test
```

[gh-project]: https://github.com/manicminer/hamilton
[ms-graph-docs]: https://docs.microsoft.com/en-us/graph/overview
