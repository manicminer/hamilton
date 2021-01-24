package clients_test

import (
	"fmt"
	"testing"

	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/clients"
	"github.com/manicminer/hamilton/clients/internal"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/models"
)

type GroupsClientTest struct {
	connection   *internal.Connection
	client       *clients.GroupsClient
	randomString string
}

func TestGroupsClient(t *testing.T) {
	rs := internal.RandomString()
	c := GroupsClientTest{
		connection:   internal.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: rs,
	}
	c.client = clients.NewGroupsClient(c.connection.AuthConfig.TenantID)
	c.client.BaseClient.Authorizer = c.connection.Authorizer

	token, err := c.connection.Authorizer.Token()
	if err != nil {
		t.Fatalf("could not acquire access token: %v", err)
	}
	claims, err := auth.ParseClaims(token)
	if err != nil {
		t.Fatalf("could not parse claims: %v", err)
	}

	u := UsersClientTest{
		connection:   internal.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: rs,
	}
	u.client = clients.NewUsersClient(c.connection.AuthConfig.TenantID)
	u.client.BaseClient.Authorizer = c.connection.Authorizer

	newGroup := models.Group{
		DisplayName:     utils.StringPtr("Test Group"),
		MailEnabled:     utils.BoolPtr(false),
		MailNickname:    utils.StringPtr(fmt.Sprintf("test-group-%s", c.randomString)),
		SecurityEnabled: utils.BoolPtr(true),
	}
	newGroup.AppendOwner(c.client.BaseClient.Endpoint, c.client.BaseClient.ApiVersion, claims.ObjectId)
	group := testGroupsClient_Create(t, c, newGroup)
	testGroupsClient_Get(t, c, *group.ID)
	group.DisplayName = utils.StringPtr(fmt.Sprintf("test-updated-group-%s", c.randomString))

	testGroupsClient_Update(t, c, *group)

	owners := testGroupsClient_ListOwners(t, c, *group.ID)
	testGroupsClient_GetOwner(t, c, *group.ID, (*owners)[0])

	user := testUsersClient_Create(t, u, models.User{
		AccountEnabled:    utils.BoolPtr(true),
		DisplayName:       utils.StringPtr("Test User"),
		MailNickname:      utils.StringPtr(fmt.Sprintf("test-user-%s", c.randomString)),
		UserPrincipalName: utils.StringPtr(fmt.Sprintf("test-user-%s@%s", c.randomString, c.connection.DomainName)),
		PasswordProfile: &models.UserPasswordProfile{
			Password: utils.StringPtr(fmt.Sprintf("IrPa55w0rd%s", c.randomString)),
		},
	})
	group.AppendOwner(c.client.BaseClient.Endpoint, c.client.BaseClient.ApiVersion, *user.ID)
	testGroupsClient_AddOwners(t, c, group)

	testGroupsClient_RemoveOwners(t, c, *group.ID, &([]string{claims.ObjectId}))

	group.AppendMember(c.client.BaseClient.Endpoint, c.client.BaseClient.ApiVersion, claims.ObjectId)
	testGroupsClient_AddMembers(t, c, group)
	members := testGroupsClient_ListMembers(t, c, *group.ID)
	testGroupsClient_GetMember(t, c, *group.ID, (*members)[0])
	testGroupsClient_RemoveMembers(t, c, *group.ID, &([]string{claims.ObjectId}))

	testGroupsClient_List(t, c)
	testGroupsClient_Delete(t, c, *group.ID)
	testUsersClient_Delete(t, u, *user.ID)
}

func testGroupsClient_Create(t *testing.T, c GroupsClientTest, g models.Group) (group *models.Group) {
	group, status, err := c.client.Create(c.connection.Context, g)
	if err != nil {
		t.Fatalf("GroupsClient.Create(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("GroupsClient.Create(): invalid status: %d", status)
	}
	if group == nil {
		t.Fatal("GroupsClient.Create(): group was nil")
	}
	if group.ID == nil {
		t.Fatal("GroupsClient.Create(): group.ID was nil")
	}
	return
}

func testGroupsClient_Update(t *testing.T, c GroupsClientTest, g models.Group) {
	status, err := c.client.Update(c.connection.Context, g)
	if err != nil {
		t.Fatalf("GroupsClient.Update(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("GroupsClient.Update(): invalid status: %d", status)
	}
}

func testGroupsClient_List(t *testing.T, c GroupsClientTest) (groups *[]models.Group) {
	groups, _, err := c.client.List(c.connection.Context, "")
	if err != nil {
		t.Fatalf("GroupsClient.List(): %v", err)
	}
	if groups == nil {
		t.Fatal("GroupsClient.List(): groups was nil")
	}
	return
}

func testGroupsClient_Get(t *testing.T, c GroupsClientTest, id string) (group *models.Group) {
	group, status, err := c.client.Get(c.connection.Context, id)
	if err != nil {
		t.Fatalf("GroupsClient.Get(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("GroupsClient.Get(): invalid status: %d", status)
	}
	if group == nil {
		t.Fatal("GroupsClient.Get(): group was nil")
	}
	return
}

func testGroupsClient_Delete(t *testing.T, c GroupsClientTest, id string) {
	status, err := c.client.Delete(c.connection.Context, id)
	if err != nil {
		t.Fatalf("GroupsClient.Delete(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("GroupsClient.Delete(): invalid status: %d", status)
	}
}

func testGroupsClient_ListOwners(t *testing.T, c GroupsClientTest, id string) (owners *[]string) {
	owners, status, err := c.client.ListOwners(c.connection.Context, id)
	if err != nil {
		t.Fatalf("GroupsClient.ListOwners(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("GroupsClient.ListOwners(): invalid status: %d", status)
	}
	if owners == nil {
		t.Fatal("GroupsClient.ListOwners(): owners was nil")
	}
	if len(*owners) == 0 {
		t.Fatal("GroupsClient.ListOwners(): owners was empty")
	}
	return
}

func testGroupsClient_GetOwner(t *testing.T, c GroupsClientTest, groupId string, ownerId string) (owner *string) {
	owner, status, err := c.client.GetOwner(c.connection.Context, groupId, ownerId)
	if err != nil {
		t.Fatalf("GroupsClient.GetOwner(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("GroupsClient.GetOwner(): invalid status: %d", status)
	}
	if owner == nil {
		t.Fatal("GroupsClient.GetOwner(): owner was nil")
	}
	return
}

func testGroupsClient_AddOwners(t *testing.T, c GroupsClientTest, g *models.Group) {
	status, err := c.client.AddOwners(c.connection.Context, g)
	if err != nil {
		t.Fatalf("GroupsClient.AddOwners(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("GroupsClient.AddOwners(): invalid status: %d", status)
	}
}

func testGroupsClient_RemoveOwners(t *testing.T, c GroupsClientTest, groupId string, ownerIds *[]string) {
	status, err := c.client.RemoveOwners(c.connection.Context, groupId, ownerIds)
	if err != nil {
		t.Fatalf("GroupsClient.RemoveOwners(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("GroupsClient.RemoveOwners(): invalid status: %d", status)
	}
}

func testGroupsClient_ListMembers(t *testing.T, c GroupsClientTest, id string) (members *[]string) {
	members, status, err := c.client.ListMembers(c.connection.Context, id)
	if err != nil {
		t.Fatalf("GroupsClient.ListMembers(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("GroupsClient.ListMembers(): invalid status: %d", status)
	}
	if members == nil {
		t.Fatal("GroupsClient.ListMembers(): members was nil")
	}
	if len(*members) == 0 {
		t.Fatal("GroupsClient.ListMembers(): members was empty")
	}
	return
}

func testGroupsClient_GetMember(t *testing.T, c GroupsClientTest, groupId string, memberId string) (member *string) {
	member, status, err := c.client.GetMember(c.connection.Context, groupId, memberId)
	if err != nil {
		t.Fatalf("GroupsClient.GetMember(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("GroupsClient.GetMember(): invalid status: %d", status)
	}
	if member == nil {
		t.Fatal("GroupsClient.GetMember(): member was nil")
	}
	return
}

func testGroupsClient_AddMembers(t *testing.T, c GroupsClientTest, g *models.Group) {
	status, err := c.client.AddMembers(c.connection.Context, g)
	if err != nil {
		t.Fatalf("GroupsClient.AddMembers(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("GroupsClient.AddMembers(): invalid status: %d", status)
	}
}

func testGroupsClient_RemoveMembers(t *testing.T, c GroupsClientTest, groupId string, memberIds *[]string) {
	status, err := c.client.RemoveMembers(c.connection.Context, groupId, memberIds)
	if err != nil {
		t.Fatalf("GroupsClient.RemoveMembers(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("GroupsClient.RemoveMembers(): invalid status: %d", status)
	}
}
