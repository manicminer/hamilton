package clients_test

import (
	"fmt"
	"testing"

	"github.com/manicminer/hamilton/clients"
	"github.com/manicminer/hamilton/clients/internal"
	"github.com/manicminer/hamilton/models"
)

type UsersClientTest struct {
	connection   *internal.Connection
	client       *clients.UsersClient
	randomString string
}

func TestUsersClient(t *testing.T) {
	rs := internal.RandomString()
	c := UsersClientTest{
		connection:   internal.NewConnection(),
		randomString: rs,
	}
	c.client = clients.NewUsersClient(c.connection.AuthConfig.TenantID)
	c.client.BaseClient.Authorizer = c.connection.Authorizer

	user := testUsersClient_Create(t, c, models.User{
		AccountEnabled:    internal.Bool(true),
		DisplayName:       internal.String("Test User"),
		MailNickname:      internal.String(fmt.Sprintf("test-user-%s", c.randomString)),
		UserPrincipalName: internal.String(fmt.Sprintf("test-user-%s@%s", c.randomString, c.connection.DomainName)),
		PasswordProfile: &models.UserPasswordProfile{
			Password: internal.String(fmt.Sprintf("IrPa55w0rd%s", c.randomString)),
		},
	})
	testUsersClient_Get(t, c, *user.ID)
	user.DisplayName = internal.String(fmt.Sprintf("test-updated-user-%s", c.randomString))
	testUsersClient_Update(t, c, *user)
	testUsersClient_List(t, c)

	g := GroupsClientTest{
		connection:   internal.NewConnection(),
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
	groupChild.AppendMember(g.client.BaseClient.Endpoint, g.client.BaseClient.ApiVersion, *user.ID)
	testGroupsClient_AddMembers(t, g, groupChild)

	testUsersClient_GetMemberGroups(t, c, *user.ID)
	testGroupsClient_Delete(t, g, *groupParent.ID)
	testGroupsClient_Delete(t, g, *groupChild.ID)

	testUsersClient_Delete(t, c, *user.ID)
}

func testUsersClient_Create(t *testing.T, c UsersClientTest, u models.User) (user *models.User) {
	user, status, err := c.client.Create(c.connection.Context, u)
	if err != nil {
		t.Fatalf("UsersClient.Create(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("UsersClient.Create(): invalid status: %d", status)
	}
	if user == nil {
		t.Fatal("UsersClient.Create(): user was nil")
	}
	if user.ID == nil {
		t.Fatal("UsersClient.Create(): user.ID was nil")
	}
	return
}

func testUsersClient_Update(t *testing.T, c UsersClientTest, u models.User) {
	status, err := c.client.Update(c.connection.Context, u)
	if err != nil {
		t.Fatalf("UsersClient.Update(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("UsersClient.Update(): invalid status: %d", status)
	}
}

func testUsersClient_List(t *testing.T, c UsersClientTest) (users *[]models.User) {
	users, _, err := c.client.List(c.connection.Context, "")
	if err != nil {
		t.Fatalf("UsersClient.List(): %v", err)
	}
	if users == nil {
		t.Fatal("UsersClient.List(): users was nil")
	}
	return
}

func testUsersClient_Get(t *testing.T, c UsersClientTest, id string) (user *models.User) {
	user, status, err := c.client.Get(c.connection.Context, id)
	if err != nil {
		t.Fatalf("UsersClient.Delete(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("UsersClient.Delete(): invalid status: %d", status)
	}
	if user == nil {
		t.Fatal("UsersClient.Get(): user was nil")
	}
	return
}

func testUsersClient_Delete(t *testing.T, c UsersClientTest, id string) {
	status, err := c.client.Delete(c.connection.Context, id)
	if err != nil {
		t.Fatalf("UsersClient.Delete(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("UsersClient.Delete(): invalid status: %d", status)
	}
}

func testUsersClient_GetMemberGroups(t *testing.T, c UsersClientTest, id string) (groups *[]models.Group) {
	groups, _, err := c.client.GetMemberGroups(c.connection.Context, id, "")
	if err != nil {
		t.Fatalf("UsersClient.GetMemberGroups(): %v", err)
	}

	if groups == nil {
		t.Fatal("UsersClient.GetMemberGroups(): groups was nil")
	}

	if len(*groups) != 2 {
		t.Fatalf("UsersClient.GetMemberGroups(): expected groups length 1. was: %d", len(*groups))
	}

	return
}
