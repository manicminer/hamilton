package msgraph_test

import (
	"fmt"
	"testing"

	"github.com/manicminer/hamilton/environments"
	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/msgraph"
	"github.com/manicminer/hamilton/odata"
)

func TestDelegatedPermissionGrantsClient(t *testing.T) {
	c := test.NewTest(t)
	defer c.CancelFunc()

	app := testApplicationsClient_Create(t, c, msgraph.Application{
		DisplayName: utils.StringPtr(fmt.Sprintf("test-delegatedPermissionGrants-%s", c.RandomString)),
	})

	sp := testServicePrincipalsClient_Create(t, c, msgraph.ServicePrincipal{
		AccountEnabled: utils.BoolPtr(true),
		AppId:          app.AppId,
		DisplayName:    app.DisplayName,
	})

	result := testServicePrincipalsClient_List(t, c, odata.Query{Filter: fmt.Sprintf("appId eq '%s'", environments.PublishedApis["MicrosoftGraph"])})
	if len(*result) == 0 {
		t.Fatalf("msgraph service principal not found")
	}

	user := testUsersClient_Create(t, c, msgraph.User{
		AccountEnabled:    utils.BoolPtr(true),
		DisplayName:       utils.StringPtr("test-user"),
		MailNickname:      utils.StringPtr(fmt.Sprintf("test-user-%s", c.RandomString)),
		UserPrincipalName: utils.StringPtr(fmt.Sprintf("test-user-%s@%s", c.RandomString, c.Connections["default"].DomainName)),
		PasswordProfile: &msgraph.UserPasswordProfile{
			Password: utils.StringPtr(fmt.Sprintf("IrPa55w0rd%s", c.RandomString)),
		},
	})

	grant := testDelegatedPermissionGrantsClient_Create(t, c, msgraph.DelegatedPermissionGrant{
		ClientId:    sp.ID(),
		ConsentType: utils.StringPtr(msgraph.DelegatedPermissionGrantConsentTypePrincipal),
		PrincipalId: user.ID(),
		ResourceId:  (*result)[0].ID(),
		Scopes:      &[]string{"openid", "User.Read"},
	})

	testDelegatedPermissionGrantsClient_Get(t, c, *grant.Id)

	testDelegatedPermissionGrantsClient_Update(t, c, *grant)

	testDelegatedPermissionGrantsClient_List(t, c, odata.Query{})
	testDelegatedPermissionGrantsClient_List(t, c, odata.Query{Filter: fmt.Sprintf("clientId eq '%s'", *sp.ID())})

	testDelegatedPermissionGrantsClient_Delete(t, c, *grant.Id)
	testUsersClient_Delete(t, c, *user.ID())
	testServicePrincipalsClient_Delete(t, c, *sp.ID())
	testApplicationsClient_Delete(t, c, *app.ID())
}

func testDelegatedPermissionGrantsClient_Create(t *testing.T, c *test.Test, sp msgraph.DelegatedPermissionGrant) (delegatedPermissionGrant *msgraph.DelegatedPermissionGrant) {
	delegatedPermissionGrant, status, err := c.DelegatedPermissionGrantsClient.Create(c.Context, sp)
	if err != nil {
		t.Fatalf("DelegatedPermissionGrantsClient.Create(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("DelegatedPermissionGrantsClient.Create(): invalid status: %d", status)
	}
	if delegatedPermissionGrant == nil {
		t.Fatal("DelegatedPermissionGrantsClient.Create(): delegatedPermissionGrant was nil")
	}
	if delegatedPermissionGrant.Id == nil {
		t.Fatal("DelegatedPermissionGrantsClient.Create(): delegatedPermissionGrant.ID was nil")
	}
	return
}

func testDelegatedPermissionGrantsClient_Update(t *testing.T, c *test.Test, sp msgraph.DelegatedPermissionGrant) (delegatedPermissionGrant *msgraph.DelegatedPermissionGrant) {
	status, err := c.DelegatedPermissionGrantsClient.Update(c.Context, sp)
	if err != nil {
		t.Fatalf("DelegatedPermissionGrantsClient.Update(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("DelegatedPermissionGrantsClient.Update(): invalid status: %d", status)
	}
	return
}

func testDelegatedPermissionGrantsClient_List(t *testing.T, c *test.Test, query odata.Query) (delegatedPermissionGrants *[]msgraph.DelegatedPermissionGrant) {
	query.Top = 10
	delegatedPermissionGrants, _, err := c.DelegatedPermissionGrantsClient.List(c.Context, query)
	if err != nil {
		t.Fatalf("DelegatedPermissionGrantsClient.List(): %v", err)
	}
	if delegatedPermissionGrants == nil {
		t.Fatal("DelegatedPermissionGrantsClient.List(): delegatedPermissionGrants was nil")
	}
	return
}

func testDelegatedPermissionGrantsClient_Get(t *testing.T, c *test.Test, id string) (delegatedPermissionGrant *msgraph.DelegatedPermissionGrant) {
	delegatedPermissionGrant, status, err := c.DelegatedPermissionGrantsClient.Get(c.Context, id, odata.Query{})
	if err != nil {
		t.Fatalf("DelegatedPermissionGrantsClient.Get(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("DelegatedPermissionGrantsClient.Get(): invalid status: %d", status)
	}
	if delegatedPermissionGrant == nil {
		t.Fatal("DelegatedPermissionGrantsClient.Get(): delegatedPermissionGrant was nil")
	}
	return
}

func testDelegatedPermissionGrantsClient_Delete(t *testing.T, c *test.Test, id string) {
	status, err := c.DelegatedPermissionGrantsClient.Delete(c.Context, id)
	if err != nil {
		t.Fatalf("DelegatedPermissionGrantsClient.Delete(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("DelegatedPermissionGrantsClient.Delete(): invalid status: %d", status)
	}
}
