package clients_test

import (
	"fmt"
	"testing"

	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/clients"
	"github.com/manicminer/hamilton/clients/internal"
	"github.com/manicminer/hamilton/models"
)

type ServicePrincipalsClientTest struct {
	connection   *internal.Connection
	client       *clients.ServicePrincipalsClient
	randomString string
}

func TestServicePrincipalsClient(t *testing.T) {
	rs := internal.RandomString()
	c := ServicePrincipalsClientTest{
		connection:   internal.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: rs,
	}
	c.client = clients.NewServicePrincipalsClient(c.connection.AuthConfig.TenantID)
	c.client.BaseClient.Authorizer = c.connection.Authorizer

	a := ApplicationsClientTest{
		connection:   internal.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: rs,
	}
	a.client = clients.NewApplicationsClient(c.connection.AuthConfig.TenantID)
	a.client.BaseClient.Authorizer = c.connection.Authorizer
	app := testApplicationsClient_Create(t, a, models.Application{
		DisplayName: internal.String(fmt.Sprintf("test-serviceprincipal-%s", a.randomString)),
	})

	sp := testServicePrincipalsClient_Create(t, c, models.ServicePrincipal{
		AccountEnabled: internal.Bool(true),
		AppId:          app.AppId,
		DisplayName:    internal.String(fmt.Sprintf("test-serviceprincipal-%s", c.randomString)),
	})
	testServicePrincipalsClient_Get(t, c, *sp.ID)
	sp.Tags = &([]string{"TestTag"})
	testServicePrincipalsClient_Update(t, c, *sp)
	testServicePrincipalsClient_List(t, c)

	g := GroupsClientTest{
		connection:   internal.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: rs,
	}
	g.client = clients.NewGroupsClient(g.connection.AuthConfig.TenantID)
	g.client.BaseClient.Authorizer = g.connection.Authorizer

	newGroupParent := models.Group{
		DisplayName:     internal.String("Test Group Parent"),
		MailEnabled:     internal.Bool(false),
		MailNickname:    internal.String(fmt.Sprintf("test-group-parent-%s", c.randomString)),
		SecurityEnabled: internal.Bool(true),
	}
	newGroupChild := models.Group{
		DisplayName:     internal.String("Test Group Child"),
		MailEnabled:     internal.Bool(false),
		MailNickname:    internal.String(fmt.Sprintf("test-group-child-%s", c.randomString)),
		SecurityEnabled: internal.Bool(true),
	}

	groupParent := testGroupsClient_Create(t, g, newGroupParent)
	groupChild := testGroupsClient_Create(t, g, newGroupChild)
	groupParent.AppendMember(g.client.BaseClient.Endpoint, g.client.BaseClient.ApiVersion, *groupChild.ID)
	testGroupsClient_AddMembers(t, g, groupParent)
	groupChild.AppendMember(g.client.BaseClient.Endpoint, g.client.BaseClient.ApiVersion, *sp.ID)
	testGroupsClient_AddMembers(t, g, groupChild)

	testServicePrincipalsClient_ListGroupMemberships(t, c, *sp.ID)
	testGroupsClient_Delete(t, g, *groupParent.ID)
	testGroupsClient_Delete(t, g, *groupChild.ID)

	testServicePrincipalsClient_Delete(t, c, *sp.ID)

	testApplicationsClient_Delete(t, a, *app.ID)
}

func testServicePrincipalsClient_Create(t *testing.T, c ServicePrincipalsClientTest, sp models.ServicePrincipal) (servicePrincipal *models.ServicePrincipal) {
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

func testServicePrincipalsClient_Update(t *testing.T, c ServicePrincipalsClientTest, sp models.ServicePrincipal) (servicePrincipal *models.ServicePrincipal) {
	status, err := c.client.Update(c.connection.Context, sp)
	if err != nil {
		t.Fatalf("ServicePrincipalsClient.Update(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ServicePrincipalsClient.Update(): invalid status: %d", status)
	}
	return
}

func testServicePrincipalsClient_List(t *testing.T, c ServicePrincipalsClientTest) (servicePrincipals *[]models.ServicePrincipal) {
	servicePrincipals, _, err := c.client.List(c.connection.Context, "")
	if err != nil {
		t.Fatalf("ServicePrincipalsClient.List(): %v", err)
	}
	if servicePrincipals == nil {
		t.Fatal("ServicePrincipalsClient.List(): servicePrincipals was nil")
	}
	return
}

func testServicePrincipalsClient_Get(t *testing.T, c ServicePrincipalsClientTest, id string) (servicePrincipal *models.ServicePrincipal) {
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

func testServicePrincipalsClient_ListGroupMemberships(t *testing.T, c ServicePrincipalsClientTest, id string) (groups *[]models.Group) {
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
