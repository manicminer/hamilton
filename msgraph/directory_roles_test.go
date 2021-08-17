package msgraph_test

import (
	"fmt"
	"testing"

	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/msgraph"
)

type DirectoryRolesClientTest struct {
	connection   *test.Connection
	client       *msgraph.DirectoryRolesClient
	randomString string
}

func TestDirectoryRolesClient(t *testing.T) {
	rs := test.RandomString()
	// set up groups test client
	groupsClient := GroupsClientTest{
		connection:   test.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: rs,
	}
	groupsClient.client = msgraph.NewGroupsClient(groupsClient.connection.AuthConfig.TenantID)
	groupsClient.client.BaseClient.Authorizer = groupsClient.connection.Authorizer

	// set up directory roles test client
	dirRolesClient := DirectoryRolesClientTest{
		connection:   test.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: rs,
	}
	dirRolesClient.client = msgraph.NewDirectoryRolesClient(dirRolesClient.connection.AuthConfig.TenantID)
	dirRolesClient.client.BaseClient.Authorizer = dirRolesClient.connection.Authorizer

	// list directory roles; usually at least few directory roles are activated within a tenant
	directoryRoles := testDirectoryRolesClient_List(t, dirRolesClient)
	directoryRole := (*directoryRoles)[0]
	testDirectoryRolesClient_Get(t, dirRolesClient, *directoryRole.ID)

	// create a new test group which can be later assigned as a member of the previously listed directory role
	newGroup := msgraph.Group{
		DisplayName:     utils.StringPtr("test-group-directoryRoles"),
		MailEnabled:     utils.BoolPtr(false),
		MailNickname:    utils.StringPtr(fmt.Sprintf("test-group-%s", groupsClient.randomString)),
		SecurityEnabled: utils.BoolPtr(true),
		// required attribute to set if you plan to assign a directory role to an object (e.g. user, group etc)
		// can be only set on a group creation using MS Graph Beta API
		IsAssignableToRole: utils.BoolPtr(true),
	}
	group := testGroupsClient_Create(t, groupsClient, newGroup)

	// add the test group as a member of directory role
	directoryRole.Members = &msgraph.Members{group.DirectoryObject}
	testDirectoryRolesClient_AddMembers(t, dirRolesClient, &directoryRole)

	// list members of the directory role; then remove the added group member to clean up
	testDirectoryRolesClient_ListMembers(t, dirRolesClient, *directoryRole.ID)
	testDirectoryRolesClient_GetMember(t, dirRolesClient, *directoryRole.ID, *group.ID)
	testDirectoryRolesClient_RemoveMembers(t, dirRolesClient, *directoryRole.ID, &[]string{*group.ID})

	// remove the test group to clean up
	testGroupsClient_Delete(t, groupsClient, *group.ID)
}

func testDirectoryRolesClient_List(t *testing.T, c DirectoryRolesClientTest) (directoryRoles *[]msgraph.DirectoryRole) {
	directoryRoles, _, err := c.client.List(c.connection.Context)
	if err != nil {
		t.Fatalf("DirectoryRolesClient.List(): %v", err)
	}
	if directoryRoles == nil {
		t.Fatal("DirectoryRolesClient.List(): directoryRoles was nil")
	}
	return
}

func testDirectoryRolesClient_Get(t *testing.T, c DirectoryRolesClientTest, id string) (directoryRole *msgraph.DirectoryRole) {
	directoryRole, status, err := c.client.Get(c.connection.Context, id)
	if err != nil {
		t.Fatalf("DirectoryRolesClient.Get(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("DirectoryRolesClient.Get(): invalid status: %d", status)
	}
	if directoryRole == nil {
		t.Fatal("DirectoryRolesClient.Get(): directoryRole was nil")
	}
	return
}

func testDirectoryRolesClient_ListMembers(t *testing.T, c DirectoryRolesClientTest, id string) (members *[]string) {
	members, status, err := c.client.ListMembers(c.connection.Context, id)
	if err != nil {
		t.Fatalf("DirectoryRolesClient.ListMembers(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("DirectoryRolesClient.ListMembers(): invalid status: %d", status)
	}
	if members == nil {
		t.Fatal("DirectoryRolesClient.ListMembers(): members was nil")
	}
	if len(*members) == 0 {
		t.Fatal("DirectoryRolesClient.ListMembers(): members was empty")
	}
	return
}

func testDirectoryRolesClient_GetMember(t *testing.T, c DirectoryRolesClientTest, dirRoleId string, memberId string) (member *string) {
	member, status, err := c.client.GetMember(c.connection.Context, dirRoleId, memberId)
	if err != nil {
		t.Fatalf("DirectoryRolesClient.GetMember(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("DirectoryRolesClient.GetMember(): invalid status: %d", status)
	}
	if member == nil {
		t.Fatal("DirectoryRolesClient.GetMember(): member was nil")
	}
	return
}

func testDirectoryRolesClient_RemoveMembers(t *testing.T, c DirectoryRolesClientTest, dirRoleId string, memberIds *[]string) {
	status, err := c.client.RemoveMembers(c.connection.Context, dirRoleId, memberIds)
	if err != nil {
		t.Fatalf("DirectoryRolesClient.RemoveMembers(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("DirectoryRolesClient.RemoveMembers(): invalid status: %d", status)
	}
}

func testDirectoryRolesClient_AddMembers(t *testing.T, c DirectoryRolesClientTest, dirRole *msgraph.DirectoryRole) {
	status, err := c.client.AddMembers(c.connection.Context, dirRole)
	if err != nil {
		t.Fatalf("DirectoryRolesClient.AddMembers(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("DirectoryRolesClient.AddMembers(): invalid status: %d", status)
	}
}
