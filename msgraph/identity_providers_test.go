package msgraph_test

import (
	"strings"
	"testing"

	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/msgraph"
)

func TestIdentityProvidersClient(t *testing.T) {
	c := test.NewTest(t)
	defer c.CancelFunc()

	providers := testIdentityProvidersClient_List(t, c)
	for _, provider := range *providers {
		if strings.EqualFold(*provider.Type, "Google") {
			testIdentityProvidersClient_Delete(t, c, *provider.ID)
		}
	}

	testIdentityProvidersClient_ListAvailableProviderTypes(t, c)

	identityProvider := testIdentityProvidersClient_Create(t, c, msgraph.IdentityProvider{
		ODataType:    utils.StringPtr(odata.TypeSocialIdentityProvider),
		Name:         utils.StringPtr("Login with Google"),
		Type:         utils.StringPtr("Google"),
		ClientId:     utils.StringPtr("56433757-cadd-4135-8431-2c9e3fd68ae8"),
		ClientSecret: utils.StringPtr("000000000000"),
	})

	testIdentityProvidersClient_Get(t, c, *identityProvider.ID)

	patchIdentityProvider := &msgraph.IdentityProvider{}
	patchIdentityProvider.ODataType = identityProvider.ODataType
	patchIdentityProvider.ID = identityProvider.ID
	patchIdentityProvider.ClientSecret = utils.StringPtr("1111111111111")
	testIdentityProvidersClient_Update(t, c, *patchIdentityProvider)

	testIdentityProvidersClient_List(t, c)

	testIdentityProvidersClient_Delete(t, c, *identityProvider.ID)
}

func testIdentityProvidersClient_Create(t *testing.T, c *test.Test, p msgraph.IdentityProvider) (provider *msgraph.IdentityProvider) {
	provider, status, err := c.IdentityProvidersClient.Create(c.Context, p)
	if err != nil {
		t.Fatalf("IdentityProvidersClient.Create(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("IdentityProvidersClient.Create(): invalid status: %d", status)
	}
	if provider == nil {
		t.Fatal("IdentityProvidersClient.Create(): provider was nil")
	}
	if provider.ID == nil {
		t.Fatal("IdentityProvidersClient.Create(): provider.ID was nil")
	}
	return
}

func testIdentityProvidersClient_Update(t *testing.T, c *test.Test, p msgraph.IdentityProvider) {
	status, err := c.IdentityProvidersClient.Update(c.Context, p)
	if err != nil {
		t.Fatalf("IdentityProvidersClient.Update(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("IdentityProvidersClient.Update(): invalid status: %d", status)
	}
}

func testIdentityProvidersClient_List(t *testing.T, c *test.Test) (providers *[]msgraph.IdentityProvider) {
	providers, _, err := c.IdentityProvidersClient.List(c.Context)
	if err != nil {
		t.Fatalf("IdentityProvidersClient.List(): %v", err)
	}
	if providers == nil {
		t.Fatal("IdentityProvidersClient.List(): providers was nil")
	}
	return
}

func testIdentityProvidersClient_Get(t *testing.T, c *test.Test, id string) (provider *msgraph.IdentityProvider) {
	provider, status, err := c.IdentityProvidersClient.Get(c.Context, id)
	if err != nil {
		t.Fatalf("IdentityProvidersClient.Get(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("IdentityProvidersClient.Get(): invalid status: %d", status)
	}
	if provider == nil {
		t.Fatal("IdentityProvidersClient.Get(): provider was nil")
	}
	if provider.ID == nil {
		t.Fatal("IdentityProvidersClient.Get(): provider.ID was nil")
	}
	return
}

func testIdentityProvidersClient_Delete(t *testing.T, c *test.Test, id string) {
	status, err := c.IdentityProvidersClient.Delete(c.Context, id)
	if err != nil {
		t.Fatalf("IdentityProvidersClient.Delete(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("IdentityProvidersClient.Delete(): invalid status: %d", status)
	}
}

func testIdentityProvidersClient_ListAvailableProviderTypes(t *testing.T, c *test.Test) {
	availableIdentityProviders, _, err := c.IdentityProvidersClient.ListAvailableProviderTypes(c.Context)
	if err != nil {
		t.Fatalf("IdentityProvidersClient.ListAvailableProviderTypes(): %v", err)
	}

	if availableIdentityProviders == nil {
		t.Fatal("IdentityProvidersClient.ListAvailableProviderTypes(): availableIdentityProviders was nil")
	}

	if len(*availableIdentityProviders) == 0 {
		t.Fatal("IdentityProvidersClient.ListAvailableProviderTypes(): expected availableIdentityProviders at least one available provider")
	}
}
