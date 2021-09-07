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

type DirectoryObjectsClientTest struct {
	connection   *test.Connection
	client       *msgraph.DirectoryObjectsClient
	randomString string
}

func TestDirectoryObjectsClient(t *testing.T) {
	rs := test.RandomString()
	c := DirectoryObjectsClientTest{
		connection:   test.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: rs,
	}
	c.client = msgraph.NewDirectoryObjectsClient(c.connection.AuthConfig.TenantID)
	c.client.BaseClient.Authorizer = c.connection.Authorizer

	g := GroupsClientTest{
		connection:   test.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: rs,
	}
	g.client = msgraph.NewGroupsClient(c.connection.AuthConfig.TenantID)
	g.client.BaseClient.Authorizer = c.connection.Authorizer

	u := UsersClientTest{
		connection:   test.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: rs,
	}
	u.client = msgraph.NewUsersClient(c.connection.AuthConfig.TenantID)
	u.client.BaseClient.Authorizer = c.connection.Authorizer

	user := testUsersClient_Create(t, u, msgraph.User{
		AccountEnabled:    utils.BoolPtr(true),
		DisplayName:       utils.StringPtr("test-user"),
		MailNickname:      utils.StringPtr(fmt.Sprintf("test-user-directoryobject-%s", c.randomString)),
		UserPrincipalName: utils.StringPtr(fmt.Sprintf("test-user-directoryobject-%s@%s", c.randomString, c.connection.DomainName)),
		PasswordProfile: &msgraph.UserPasswordProfile{
			Password: utils.StringPtr(fmt.Sprintf("IrPa55w0rd%s", c.randomString)),
		},
	})

	newGroup1 := msgraph.Group{
		DisplayName:     utils.StringPtr("test-group-directoryobject-member"),
		MailEnabled:     utils.BoolPtr(false),
		MailNickname:    utils.StringPtr(fmt.Sprintf("test-group-directoryobject-member-%s", c.randomString)),
		SecurityEnabled: utils.BoolPtr(true),
	}
	newGroup1.Members = &msgraph.Members{user.DirectoryObject}
	group1 := testGroupsClient_Create(t, g, newGroup1)

	newGroup2 := msgraph.Group{
		DisplayName:     utils.StringPtr("test-group-directoryobject"),
		MailEnabled:     utils.BoolPtr(false),
		MailNickname:    utils.StringPtr(fmt.Sprintf("test-group-directoryobject-%s", c.randomString)),
		SecurityEnabled: utils.BoolPtr(true),
		Members: &msgraph.Members{
			group1.DirectoryObject,
			user.DirectoryObject,
		},
	}
	group2 := testGroupsClient_Create(t, g, newGroup2)

	testDirectoryObjectsClient_Get(t, c, *user.ID)
	testDirectoryObjectsClient_Get(t, c, *group1.ID)
	testDirectoryObjectsClient_GetMemberGroups(t, c, *user.ID, true, []string{*group1.ID, *group2.ID})
	testDirectoryObjectsClient_GetMemberObjects(t, c, *group1.ID, true, []string{*group2.ID})
	testDirectoryObjectsClient_GetByIds(t, c, []string{*group1.ID, *group2.ID, *user.ID}, []string{odata.ShortTypeGroup})
	testDirectoryObjectsClient_Delete(t, c, *group1.ID)
}

func testDirectoryObjectsClient_Get(t *testing.T, c DirectoryObjectsClientTest, id string) (directoryObject *msgraph.DirectoryObject) {
	directoryObject, status, err := c.client.Get(c.connection.Context, id, odata.Query{})
	if err != nil {
		t.Fatalf("DirectoryObjectsClient.Get(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("DirectoryObjectsClient.Get(): invalid status: %d", status)
	}
	if directoryObject == nil {
		t.Fatal("DirectoryObjectsClient.Get(): directoryObject was nil")
	}
	if directoryObject.ID == nil {
		t.Fatal("DirectoryObjectsClient.Get(): directoryObject ID was nil")
	}
	return
}

func testDirectoryObjectsClient_GetByIds(t *testing.T, c DirectoryObjectsClientTest, ids []string, types []odata.ShortType) (directoryObjects *[]msgraph.DirectoryObject) {
	directoryObjects, status, err := c.client.GetByIds(c.connection.Context, ids, types)
	if err != nil {
		t.Fatalf("DirectoryObjectsClient.GetByIds(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("DirectoryObjectsClient.GetByIds(): invalid status: %d", status)
	}
	if directoryObjects == nil {
		t.Fatal("DirectoryObjectsClient.GetByIds(): directoryObject was nil")
	}
	if len(*directoryObjects) == 0 {
		t.Fatal("DirectoryObjectsClient.GetByIds(): directoryObjects was empty")
	}
	return
}

func testDirectoryObjectsClient_GetMemberGroups(t *testing.T, c DirectoryObjectsClientTest, id string, securityEnabledOnly bool, expected []string) (directoryObjects *[]msgraph.DirectoryObject) {
	directoryObjects, status, err := c.client.GetMemberGroups(c.connection.Context, id, securityEnabledOnly)
	if err != nil {
		t.Fatalf("DirectoryObjectsClient.GetMemberGroups(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("DirectoryObjectsClient.GetMemberGroups(): invalid status: %d", status)
	}
	if directoryObjects == nil {
		t.Fatal("DirectoryObjectsClient.GetMemberGroups(): directoryObjects was nil")
	}
	if len(*directoryObjects) == 0 {
		t.Fatal("DirectoryObjectsClient.GetMemberGroups(): directoryObjects was empty")
	}

	expectedCount := len(expected)
	if len(*directoryObjects) < expectedCount {
		t.Fatalf("DirectoryObjectsClient.GetMemberGroups(): expected at least %d result. has: %d", expectedCount, len(*directoryObjects))
	}
	var actualCount int
	for _, e := range expected {
		for _, o := range *directoryObjects {
			if o.ID != nil && e == *o.ID {
				actualCount++
			}
		}
	}
	if actualCount < expectedCount {
		t.Fatalf("DirectoryObjectsClient.GetMemberGroups(): expected %d matching objects. found: %d", expectedCount, actualCount)
	}
	return
}

func testDirectoryObjectsClient_GetMemberObjects(t *testing.T, c DirectoryObjectsClientTest, id string, securityEnabledOnly bool, expected []string) (directoryObjects *[]msgraph.DirectoryObject) {
	directoryObjects, status, err := c.client.GetMemberObjects(c.connection.Context, id, securityEnabledOnly)
	if err != nil {
		t.Fatalf("DirectoryObjectsClient.GetMemberObjects(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("DirectoryObjectsClient.GetMemberObjects(): invalid status: %d", status)
	}
	if directoryObjects == nil {
		t.Fatal("DirectoryObjectsClient.GetMemberObjects(): directoryObjects was nil")
	}
	if len(*directoryObjects) == 0 {
		t.Fatal("DirectoryObjectsClient.GetMemberObjects(): directoryObjects was empty")
	}

	expectedCount := len(expected)
	if len(*directoryObjects) < expectedCount {
		t.Fatalf("DirectoryObjectsClient.GetMemberObjects(): expected at least %d result. has: %d", expectedCount, len(*directoryObjects))
	}
	var actualCount int
	for _, e := range expected {
		for _, o := range *directoryObjects {
			if o.ID != nil && e == *o.ID {
				actualCount++
			}
		}
	}
	if actualCount < expectedCount {
		t.Fatalf("DirectoryObjectsClient.GetMemberObjects(): expected %d matching objects. found: %d", expectedCount, actualCount)
	}
	return
}

func testDirectoryObjectsClient_Delete(t *testing.T, c DirectoryObjectsClientTest, id string) {
	status, err := c.client.Delete(c.connection.Context, id)
	if err != nil {
		t.Fatalf("DirectoryObjectsClient.Delete(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("DirectoryObjectsClient.Delete(): invalid status: %d", status)
	}
}
