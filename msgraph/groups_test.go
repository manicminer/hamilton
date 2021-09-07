package msgraph_test

import (
	"fmt"
	"testing"

	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/msgraph"
	"github.com/manicminer/hamilton/odata"
)

type GroupsClientTest struct {
	connection   *test.Connection
	client       *msgraph.GroupsClient
	randomString string
}

func TestGroupsClient(t *testing.T) {
	rs := test.RandomString()
	c := GroupsClientTest{
		connection:   test.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: rs,
	}
	c.client = msgraph.NewGroupsClient(c.connection.AuthConfig.TenantID)
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
		connection:   test.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: rs,
	}
	u.client = msgraph.NewUsersClient(c.connection.AuthConfig.TenantID)
	u.client.BaseClient.Authorizer = c.connection.Authorizer

	o := DirectoryObjectsClientTest{
		connection:   test.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: rs,
	}
	o.client = msgraph.NewDirectoryObjectsClient(c.connection.AuthConfig.TenantID)
	o.client.BaseClient.Authorizer = o.connection.Authorizer

	self := testDirectoryObjectsClient_Get(t, o, claims.ObjectId)

	newGroup := msgraph.Group{
		DisplayName:     utils.StringPtr("test-group"),
		MailEnabled:     utils.BoolPtr(false),
		MailNickname:    utils.StringPtr(fmt.Sprintf("test-group-%s", c.randomString)),
		SecurityEnabled: utils.BoolPtr(true),
		Owners:          &msgraph.Owners{*self},
		Members:         &msgraph.Members{*self},
	}
	group := testGroupsClient_Create(t, c, newGroup)
	testGroupsClient_Get(t, c, *group.ID)

	owners := testGroupsClient_ListOwners(t, c, *group.ID)
	testGroupsClient_GetOwner(t, c, *group.ID, (*owners)[0])

	members := testGroupsClient_ListMembers(t, c, *group.ID)
	testGroupsClient_GetMember(t, c, *group.ID, (*members)[0])

	group.DisplayName = utils.StringPtr(fmt.Sprintf("test-updated-group-%s", c.randomString))
	testGroupsClient_Update(t, c, *group)

	user := testUsersClient_Create(t, u, msgraph.User{
		AccountEnabled:    utils.BoolPtr(true),
		DisplayName:       utils.StringPtr("test-user"),
		MailNickname:      utils.StringPtr(fmt.Sprintf("test-user-%s", c.randomString)),
		UserPrincipalName: utils.StringPtr(fmt.Sprintf("test-user-%s@%s", c.randomString, c.connection.DomainName)),
		PasswordProfile: &msgraph.UserPasswordProfile{
			Password: utils.StringPtr(fmt.Sprintf("IrPa55w0rd%s", c.randomString)),
		},
	})

	group.Owners = &msgraph.Owners{user.DirectoryObject}
	testGroupsClient_AddOwners(t, c, group)
	testGroupsClient_RemoveOwners(t, c, *group.ID, &([]string{claims.ObjectId}))

	group.Members = &msgraph.Members{user.DirectoryObject}
	testGroupsClient_AddMembers(t, c, group)
	testGroupsClient_RemoveMembers(t, c, *group.ID, &([]string{claims.ObjectId}))

	testGroupsClient_List(t, c)
	testGroupsClient_Delete(t, c, *group.ID)
	testUsersClient_Delete(t, u, *user.ID)

	newGroup365 := msgraph.Group{
		DisplayName:     utils.StringPtr("test-group-365"),
		GroupTypes:      []msgraph.GroupType{msgraph.GroupTypeUnified},
		MailEnabled:     utils.BoolPtr(true),
		MailNickname:    utils.StringPtr(fmt.Sprintf("test-365-group-%s", c.randomString)),
		SecurityEnabled: utils.BoolPtr(true),
	}
	group365 := testGroupsClient_Create(t, c, newGroup365)
	testGroupsClient_Delete(t, c, *group365.ID)
	testGroupsClient_GetDeleted(t, c, *group365.ID)
	testGroupsClient_RestoreDeleted(t, c, *group365.ID)
	testGroupsClient_Delete(t, c, *group365.ID)
	testGroupsClient_ListDeleted(t, c, *group365.ID)
	testGroupsClient_DeletePermanently(t, c, *group365.ID)
}

func testGroupsClient_Create(t *testing.T, c GroupsClientTest, g msgraph.Group) (group *msgraph.Group) {
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

func testGroupsClient_Update(t *testing.T, c GroupsClientTest, g msgraph.Group) {
	status, err := c.client.Update(c.connection.Context, g)
	if err != nil {
		t.Fatalf("GroupsClient.Update(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("GroupsClient.Update(): invalid status: %d", status)
	}
}

func testGroupsClient_List(t *testing.T, c GroupsClientTest) (groups *[]msgraph.Group) {
	groups, _, err := c.client.List(c.connection.Context, odata.Query{Top: 10})
	if err != nil {
		t.Fatalf("GroupsClient.List(): %v", err)
	}
	if groups == nil {
		t.Fatal("GroupsClient.List(): groups was nil")
	}
	return
}

func testGroupsClient_Get(t *testing.T, c GroupsClientTest, id string) (group *msgraph.Group) {
	group, status, err := c.client.Get(c.connection.Context, id, odata.Query{})
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

func testGroupsClient_DeletePermanently(t *testing.T, c GroupsClientTest, id string) {
	status, err := c.client.DeletePermanently(c.connection.Context, id)
	if err != nil {
		t.Fatalf("GroupsClient.DeletePermanently(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("GroupsClient.DeletePermanently(): invalid status: %d", status)
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

func testGroupsClient_AddOwners(t *testing.T, c GroupsClientTest, g *msgraph.Group) {
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

func testGroupsClient_AddMembers(t *testing.T, c GroupsClientTest, g *msgraph.Group) {
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

func testGroupsClient_GetDeleted(t *testing.T, c GroupsClientTest, id string) (group *msgraph.Group) {
	group, status, err := c.client.GetDeleted(c.connection.Context, id, odata.Query{})
	if err != nil {
		t.Fatalf("GroupsClient.GetDeleted(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("GroupsClient.GetDeleted(): invalid status: %d", status)
	}
	if group == nil {
		t.Fatal("GroupsClient.GetDeleted(): group was nil")
	}
	return
}

func testGroupsClient_ListDeleted(t *testing.T, c GroupsClientTest, expectedId string) (deletedGroups *[]msgraph.Group) {
	deletedGroups, status, err := c.client.ListDeleted(c.connection.Context, odata.Query{
		Filter: fmt.Sprintf("id eq '%s'", expectedId),
		Top:    10,
	})
	if err != nil {
		t.Fatalf("GroupsClient.ListDeleted(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("GroupsClient.ListDeleted(): invalid status: %d", status)
	}
	if deletedGroups == nil {
		t.Fatal("GroupsClient.ListDeleted(): deletedGroups was nil")
	}
	if len(*deletedGroups) == 0 {
		t.Fatal("GroupsClient.ListDeleted(): expected at least 1 deleted group, was: 0")
	}
	found := false
	for _, group := range *deletedGroups {
		if group.ID != nil && *group.ID == expectedId {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("GroupsClient.ListDeleted(): expected group ID %q in result", expectedId)
	}
	return
}

func testGroupsClient_RestoreDeleted(t *testing.T, c GroupsClientTest, id string) {
	group, status, err := c.client.RestoreDeleted(c.connection.Context, id)
	if err != nil {
		t.Fatalf("GroupsClient.RestoreDeleted(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("GroupsClient.RestoreDeleted(): invalid status: %d", status)
	}
	if group == nil {
		t.Fatal("GroupsClient.RestoreDeleted(): group was nil")
	}
	if group.ID == nil {
		t.Fatal("GroupsClient.RestoreDeleted(): group.ID was nil")
	}
	if *group.ID != id {
		t.Fatal("GroupsClient.RestoreDeleted(): group IDs do not match")
	}
}
