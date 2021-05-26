package msgraph_test

import (
	"fmt"
	"testing"

	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/msgraph"
)

type ServicePrincipalsClientTest struct {
	connection   *test.Connection
	client       *msgraph.ServicePrincipalsClient
	randomString string
}

func TestServicePrincipalsClient(t *testing.T) {
	rs := test.RandomString()
	c := ServicePrincipalsClientTest{
		connection:   test.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: rs,
	}
	c.client = msgraph.NewServicePrincipalsClient(c.connection.AuthConfig.TenantID)
	c.client.BaseClient.Authorizer = c.connection.Authorizer

	a := ApplicationsClientTest{
		connection:   test.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: rs,
	}
	a.client = msgraph.NewApplicationsClient(c.connection.AuthConfig.TenantID)
	a.client.BaseClient.Authorizer = c.connection.Authorizer
	app := testApplicationsClient_Create(t, a, msgraph.Application{
		DisplayName: utils.StringPtr(fmt.Sprintf("test-serviceprincipal-%s", a.randomString)),
	})

	sp := testServicePrincipalsClient_Create(t, c, msgraph.ServicePrincipal{
		AccountEnabled: utils.BoolPtr(true),
		AppId:          app.AppId,
		DisplayName:    utils.StringPtr(fmt.Sprintf("test-serviceprincipal-%s", c.randomString)),
	})

	appChild := testApplicationsClient_Create(t, a, msgraph.Application{
		DisplayName: utils.StringPtr(fmt.Sprintf("test-serviceprincipal-child%s", a.randomString)),
	})
	spChild := testServicePrincipalsClient_Create(t, c, msgraph.ServicePrincipal{
		AccountEnabled: utils.BoolPtr(true),
		AppId:          appChild.AppId,
		DisplayName:    utils.StringPtr(fmt.Sprintf("test-serviceprincipal-child%s", c.randomString)),
	})

	spChild.AppendOwner(string(c.client.BaseClient.Endpoint), string(c.client.BaseClient.ApiVersion), *sp.ID)
	testServicePrincipalsClient_AddOwners(t, c, spChild)
	testServicePrincipalsClient_ListOwners(t, c, *spChild.ID, []string{*sp.ID})
	testServicePrincipalsClient_GetOwner(t, c, *spChild.ID, *sp.ID)
	testServicePrincipalsClient_Get(t, c, *sp.ID)
	sp.Tags = &([]string{"TestTag"})
	testServicePrincipalsClient_Update(t, c, *sp)
	pwd := testServicePrincipalsClient_AddPassword(t, c, sp)
	testServicePrincipalsClient_RemovePassword(t, c, sp, pwd)
	testServicePrincipalsClient_List(t, c)

	g := GroupsClientTest{
		connection:   test.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: rs,
	}
	g.client = msgraph.NewGroupsClient(g.connection.AuthConfig.TenantID)
	g.client.BaseClient.Authorizer = g.connection.Authorizer

	newGroupParent := msgraph.Group{
		DisplayName:     utils.StringPtr("Test Group Parent"),
		MailEnabled:     utils.BoolPtr(false),
		MailNickname:    utils.StringPtr(fmt.Sprintf("test-group-parent-%s", c.randomString)),
		SecurityEnabled: utils.BoolPtr(true),
	}
	newGroupChild := msgraph.Group{
		DisplayName:     utils.StringPtr("Test Group Child"),
		MailEnabled:     utils.BoolPtr(false),
		MailNickname:    utils.StringPtr(fmt.Sprintf("test-group-child-%s", c.randomString)),
		SecurityEnabled: utils.BoolPtr(true),
	}

	groupParent := testGroupsClient_Create(t, g, newGroupParent)
	groupChild := testGroupsClient_Create(t, g, newGroupChild)
	groupParent.AppendMember(g.client.BaseClient.Endpoint, g.client.BaseClient.ApiVersion, *groupChild.ID)
	testGroupsClient_AddMembers(t, g, groupParent)
	groupChild.AppendMember(g.client.BaseClient.Endpoint, g.client.BaseClient.ApiVersion, *sp.ID)
	testGroupsClient_AddMembers(t, g, groupChild)

	testServicePrincipalsClient_ListGroupMemberships(t, c, *sp.ID)
	testServicePrincipalsClient_ListOwnedObjects(t, c, *sp.ID)

	testServicePrincipalsClient_RemoveOwners(t, c, *spChild.ID, []string{*sp.ID})
	testGroupsClient_Delete(t, g, *groupParent.ID)
	testGroupsClient_Delete(t, g, *groupChild.ID)

	testServicePrincipalsClient_Delete(t, c, *sp.ID)
	testServicePrincipalsClient_Delete(t, c, *spChild.ID)

	testApplicationsClient_Delete(t, a, *app.ID)
	testApplicationsClient_Delete(t, a, *appChild.ID)
}

func testServicePrincipalsClient_Create(t *testing.T, c ServicePrincipalsClientTest, sp msgraph.ServicePrincipal) (servicePrincipal *msgraph.ServicePrincipal) {
	servicePrincipal, status, err := c.client.Create(c.connection.Context, sp)
	if err != nil {
		t.Fatalf("ServicePrincipalsClient.Create(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ServicePrincipalsClient.Create(): invalid status: %d", status)
	}
	if servicePrincipal == nil {
		t.Fatal("ServicePrincipalsClient.Create(): servicePrincipal was nil")
	}
	if servicePrincipal.ID == nil {
		t.Fatal("ServicePrincipalsClient.Create(): servicePrincipal.ID was nil")
	}
	return
}

func testServicePrincipalsClient_Update(t *testing.T, c ServicePrincipalsClientTest, sp msgraph.ServicePrincipal) (servicePrincipal *msgraph.ServicePrincipal) {
	status, err := c.client.Update(c.connection.Context, sp)
	if err != nil {
		t.Fatalf("ServicePrincipalsClient.Update(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ServicePrincipalsClient.Update(): invalid status: %d", status)
	}
	return
}

func testServicePrincipalsClient_List(t *testing.T, c ServicePrincipalsClientTest) (servicePrincipals *[]msgraph.ServicePrincipal) {
	servicePrincipals, _, err := c.client.List(c.connection.Context, "")
	if err != nil {
		t.Fatalf("ServicePrincipalsClient.List(): %v", err)
	}
	if servicePrincipals == nil {
		t.Fatal("ServicePrincipalsClient.List(): servicePrincipals was nil")
	}
	return
}

func testServicePrincipalsClient_Get(t *testing.T, c ServicePrincipalsClientTest, id string) (servicePrincipal *msgraph.ServicePrincipal) {
	servicePrincipal, status, err := c.client.Get(c.connection.Context, id)
	if err != nil {
		t.Fatalf("ServicePrincipalsClient.Get(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ServicePrincipalsClient.Get(): invalid status: %d", status)
	}
	if servicePrincipal == nil {
		t.Fatal("ServicePrincipalsClient.Get(): servicePrincipal was nil")
	}
	return
}

func testServicePrincipalsClient_Delete(t *testing.T, c ServicePrincipalsClientTest, id string) {
	status, err := c.client.Delete(c.connection.Context, id)
	if err != nil {
		t.Fatalf("ServicePrincipalsClient.Delete(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ServicePrincipalsClient.Delete(): invalid status: %d", status)
	}
}

func testServicePrincipalsClient_ListGroupMemberships(t *testing.T, c ServicePrincipalsClientTest, id string) (groups *[]msgraph.Group) {
	groups, _, err := c.client.ListGroupMemberships(c.connection.Context, id, "")
	if err != nil {
		t.Fatalf("ServicePrincipalsClient.ListGroupMemberships(): %v", err)
	}

	if groups == nil {
		t.Fatal("ServicePrincipalsClient.ListGroupMemberships(): groups was nil")
	}

	if len(*groups) != 2 {
		t.Fatalf("ServicePrincipalsClient.ListGroupMemberships(): expected groups length 2. was: %d", len(*groups))
	}

	return
}

func testServicePrincipalsClient_AddPassword(t *testing.T, c ServicePrincipalsClientTest, a *msgraph.ServicePrincipal) *msgraph.PasswordCredential {
	pwd := msgraph.PasswordCredential{
		DisplayName: utils.StringPtr("test password"),
	}
	newPwd, status, err := c.client.AddPassword(c.connection.Context, *a.ID, pwd)
	if err != nil {
		t.Fatalf("ServicePrincipalsClient.AddPassword(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ServicePrincipalsClient.AddPassword(): invalid status: %d", status)
	}
	if newPwd.SecretText == nil || len(*newPwd.SecretText) == 0 {
		t.Fatalf("ServicePrincipalsClient.AddPassword(): nil or empty secretText returned by API")
	}
	if *newPwd.DisplayName != *pwd.DisplayName {
		t.Fatalf("ServicePrincipalsClient.AddPassword(): password names do not match")
	}
	return newPwd
}

func testServicePrincipalsClient_RemovePassword(t *testing.T, c ServicePrincipalsClientTest, a *msgraph.ServicePrincipal, p *msgraph.PasswordCredential) {
	status, err := c.client.RemovePassword(c.connection.Context, *a.ID, *p.KeyId)
	if err != nil {
		t.Fatalf("ServicePrincipalsClient.RemovePassword(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ServicePrincipalsClient.RemovePassword(): invalid status: %d", status)
	}
}

func testServicePrincipalsClient_AddOwners(t *testing.T, c ServicePrincipalsClientTest, sp *msgraph.ServicePrincipal) {
	status, err := c.client.AddOwners(c.connection.Context, sp)
	if err != nil {
		t.Fatalf("ServicePrincipalsClient.AddOwners(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ServicePrincipalsClient.AddOwners(): invalid status: %d", status)
	}
}

func testServicePrincipalsClient_ListOwnedObjects(t *testing.T, c ServicePrincipalsClientTest, id string) (ownedObjects *[]string) {
	ownedObjects, _, err := c.client.ListOwnedObjects(c.connection.Context, id)
	if err != nil {
		t.Fatalf("ServicePrincipalsClient.ListOwnedObjects(): %v", err)
	}

	if ownedObjects == nil {
		t.Fatal("ServicePrincipalsClient.ListOwnedObjects(): ownedObjects was nil")
	}

	if len(*ownedObjects) != 1 {
		t.Fatalf("ServicePrincipalsClient.ListOwnedObjects(): expected ownedObjects length 1. was: %d", len(*ownedObjects))
	}
	return
}

func testServicePrincipalsClient_ListOwners(t *testing.T, c ServicePrincipalsClientTest, id string, expected []string) (owners *[]string) {
	owners, status, err := c.client.ListOwners(c.connection.Context, id)
	if err != nil {
		t.Fatalf("ServicePrincipalsClient.ListOwners(): %v", err)
	}

	if status < 200 || status >= 300 {
		t.Fatalf("ServicePrincipalsClient.ListOwners(): invalid status: %d", status)
	}

	ownersExpected := len(expected)

	if len(*owners) < ownersExpected {
		t.Fatalf("ServicePrincipalsClient.ListOwners(): expected at least %d owner. has: %d", ownersExpected, len(*owners))
	}

	var ownersFound int

	for _, e := range expected {
		for _, o := range *owners {
			if e == o {
				ownersFound++
				continue
			}
		}
	}

	if ownersFound < ownersExpected {
		t.Fatalf("ServicePrincipalsClient.ListOwners(): expected %d matching owners. found: %d", ownersExpected, ownersFound)
	}
	return
}

func testServicePrincipalsClient_GetOwner(t *testing.T, c ServicePrincipalsClientTest, spId, ownerId string) (owner *string) {
	owner, status, err := c.client.GetOwner(c.connection.Context, spId, ownerId)
	if err != nil {
		t.Fatalf("ServicePrincipalsClient.GetOwner(): %v", err)
	}

	if status < 200 || status >= 300 {
		t.Fatalf("ServicePrincipalsClient.GetOwner(): invalid status: %d", status)
	}

	if owner == nil {
		t.Fatalf("ServicePrincipalsClient.GetOwner(): owner was nil")
	}
	return
}

func testServicePrincipalsClient_RemoveOwners(t *testing.T, c ServicePrincipalsClientTest, spId string, ownerIds []string) {
	_, err := c.client.RemoveOwners(c.connection.Context, spId, &ownerIds)
	if err != nil {
		t.Fatalf("ServicePrincipalsClient.RemoveOwners(): %v", err)
	}
}
