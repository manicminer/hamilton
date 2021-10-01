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

type UsersClientTest struct {
	connection   *test.Connection
	client       *msgraph.UsersClient
	randomString string
}

func TestUsersClient(t *testing.T) {
	rs := test.RandomString()
	c := UsersClientTest{
		connection:   test.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: rs,
	}
	c.client = msgraph.NewUsersClient(c.connection.AuthConfig.TenantID)
	c.client.BaseClient.Authorizer = c.connection.Authorizer

	user := testUsersClient_Create(t, c, msgraph.User{
		AccountEnabled:    utils.BoolPtr(true),
		DisplayName:       utils.StringPtr("test-user"),
		MailNickname:      utils.StringPtr(fmt.Sprintf("test-user-%s", c.randomString)),
		UserPrincipalName: utils.StringPtr(fmt.Sprintf("test-user-%s@%s", c.randomString, c.connection.DomainName)),
		PasswordProfile: &msgraph.UserPasswordProfile{
			Password: utils.StringPtr(fmt.Sprintf("IrPa55w0rd%s", c.randomString)),
		},
	})
	testUsersClient_Get(t, c, *user.ID)
	user.DisplayName = utils.StringPtr(fmt.Sprintf("test-updated-user-%s", c.randomString))
	testUsersClient_Update(t, c, *user)
	testUsersClient_List(t, c)

	manager := testUsersClient_Create(t, c, msgraph.User{
		AccountEnabled:    utils.BoolPtr(true),
		DisplayName:       utils.StringPtr("test-user-manager"),
		MailNickname:      utils.StringPtr(fmt.Sprintf("test-user-manager-%s", c.randomString)),
		UserPrincipalName: utils.StringPtr(fmt.Sprintf("test-user-manager-%s@%s", c.randomString, c.connection.DomainName)),
		PasswordProfile: &msgraph.UserPasswordProfile{
			Password: utils.StringPtr(fmt.Sprintf("IrPa55w0rd%s", c.randomString)),
		},
	})
	testUsersClient_Get(t, c, *manager.ID)

	g := GroupsClientTest{
		connection:   test.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: rs,
	}
	g.client = msgraph.NewGroupsClient(g.connection.AuthConfig.TenantID)
	g.client.BaseClient.Authorizer = g.connection.Authorizer

	newGroupParent := msgraph.Group{
		DisplayName:     utils.StringPtr("test-group-parent-users"),
		MailEnabled:     utils.BoolPtr(false),
		MailNickname:    utils.StringPtr(fmt.Sprintf("test-group-parent-%s", c.randomString)),
		SecurityEnabled: utils.BoolPtr(true),
	}
	newGroupChild := msgraph.Group{
		DisplayName:     utils.StringPtr("test-group-child-users"),
		MailEnabled:     utils.BoolPtr(false),
		MailNickname:    utils.StringPtr(fmt.Sprintf("test-group-child-%s", c.randomString)),
		SecurityEnabled: utils.BoolPtr(true),
	}

	groupParent := testGroupsClient_Create(t, g, newGroupParent)
	groupChild := testGroupsClient_Create(t, g, newGroupChild)
	groupParent.Members = &msgraph.Members{groupChild.DirectoryObject}
	testGroupsClient_AddMembers(t, g, groupParent)
	groupChild.Members = &msgraph.Members{user.DirectoryObject}
	testGroupsClient_AddMembers(t, g, groupChild)

	testUsersClient_ListGroupMemberships(t, c, *user.ID)
	testGroupsClient_Delete(t, g, *groupParent.ID)
	testGroupsClient_Delete(t, g, *groupChild.ID)

	testUsersClient_AssignManager(t, c, *user.ID, *manager)
	testUsersClient_GetManager(t, c, *user.ID)
	testUsersClient_Delete(t, c, *manager.ID)

	testUsersClient_Delete(t, c, *user.ID)
	testUsersClient_ListDeleted(t, c, *user.ID)
	testUsersClient_GetDeleted(t, c, *user.ID)
	testUsersClient_RestoreDeleted(t, c, *user.ID)
	testUsersClient_Delete(t, c, *user.ID)
	testUsersClient_DeletePermanently(t, c, *user.ID)
}

func testUsersClient_Create(t *testing.T, c UsersClientTest, u msgraph.User) (user *msgraph.User) {
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

func testUsersClient_Update(t *testing.T, c UsersClientTest, u msgraph.User) {
	status, err := c.client.Update(c.connection.Context, u)
	if err != nil {
		t.Fatalf("UsersClient.Update(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("UsersClient.Update(): invalid status: %d", status)
	}
}

func testUsersClient_List(t *testing.T, c UsersClientTest) (users *[]msgraph.User) {
	users, _, err := c.client.List(c.connection.Context, odata.Query{Top: 10, Expand: odata.Expand{Relationship: "memberOf"}})
	if err != nil {
		t.Fatalf("UsersClient.List(): %v", err)
	}
	if users == nil {
		t.Fatal("UsersClient.List(): users was nil")
	}
	return
}

func testUsersClient_Get(t *testing.T, c UsersClientTest, id string) (user *msgraph.User) {
	user, status, err := c.client.Get(c.connection.Context, id, odata.Query{})
	if err != nil {
		t.Fatalf("UsersClient.Get(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("UsersClient.Get(): invalid status: %d", status)
	}
	if user == nil {
		t.Fatal("UsersClient.Get(): user was nil")
	}
	return
}

func testUsersClient_GetDeleted(t *testing.T, c UsersClientTest, id string) (user *msgraph.User) {
	user, status, err := c.client.GetDeleted(c.connection.Context, id, odata.Query{})
	if err != nil {
		t.Fatalf("UsersClient.GetDeleted(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("UsersClient.GetDeleted(): invalid status: %d", status)
	}
	if user == nil {
		t.Fatal("UsersClient.GetDeleted(): user was nil")
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

func testUsersClient_DeletePermanently(t *testing.T, c UsersClientTest, id string) {
	status, err := c.client.DeletePermanently(c.connection.Context, id)
	if err != nil {
		t.Fatalf("UsersClient.DeletePermanently(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("UsersClient.DeletePermanently(): invalid status: %d", status)
	}
}

func testUsersClient_ListGroupMemberships(t *testing.T, c UsersClientTest, id string) (groups *[]msgraph.Group) {
	groups, _, err := c.client.ListGroupMemberships(c.connection.Context, id, odata.Query{})
	if err != nil {
		t.Fatalf("UsersClient.ListGroupMemberships(): %v", err)
	}

	if groups == nil {
		t.Fatal("UsersClient.ListGroupMemberships(): groups was nil")
	}

	if len(*groups) != 2 {
		t.Fatalf("UsersClient.ListGroupMemberships(): expected groups length 2. was: %d", len(*groups))
	}

	return
}

func testUsersClient_ListDeleted(t *testing.T, c UsersClientTest, expectedId string) (deletedUsers *[]msgraph.User) {
	deletedUsers, status, err := c.client.ListDeleted(c.connection.Context, odata.Query{
		Filter: fmt.Sprintf("id eq '%s'", expectedId),
		Top:    10,
	})
	if err != nil {
		t.Fatalf("UsersClient.ListDeleted(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("UsersClient.ListDeleted(): invalid status: %d", status)
	}
	if deletedUsers == nil {
		t.Fatal("UsersClient.ListDeleted(): deletedUsers was nil")
	}
	if len(*deletedUsers) == 0 {
		t.Fatal("UsersClient.ListDeleted(): expected at least 1 deleted user. was: 0")
	}
	found := false
	for _, user := range *deletedUsers {
		if user.ID != nil && *user.ID == expectedId {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("UsersClient.ListDeleted(): expected userId %q in result", expectedId)
	}
	return
}

func testUsersClient_RestoreDeleted(t *testing.T, c UsersClientTest, id string) {
	user, status, err := c.client.RestoreDeleted(c.connection.Context, id)
	if err != nil {
		t.Fatalf("UsersClient.RestoreDeleted(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("UsersClient.RestoreDeleted(): invalid status: %d", status)
	}
	if user == nil {
		t.Fatal("UsersClient.RestoreDeleted(): user was nil")
	}
	if user.ID == nil {
		t.Fatal("UsersClient.RestoreDeleted(): user.ID was nil")
	}
	if *user.ID != id {
		t.Fatal("UsersClient.RestoreDeleted(): user ids do not match")
	}
}

func testUsersClient_AssignManager(t *testing.T, c UsersClientTest, id string, manager msgraph.User) {
	status, err := c.client.AssignManager(c.connection.Context, id, manager)
	if err != nil {
		t.Fatalf("UsersClient.AssignManager(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("UsersClient.AssignManager(): invalid status: %d", status)
	}
}

func testUsersClient_GetManager(t *testing.T, c UsersClientTest, id string) (user *msgraph.User) {
	user, status, err := c.client.GetManager(c.connection.Context, id)
	if err != nil {
		t.Fatalf("UsersClient.GetManager(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("UsersClient.GetManager(): invalid status: %d", status)
	}
	if user == nil {
		t.Fatal("UsersClient.GetManager(): user was nil")
	}
	return
}
