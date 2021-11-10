package msgraph_test

import (
	"fmt"
	"testing"

	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/msgraph"
	"github.com/manicminer/hamilton/odata"
)

func TestDirectoryObjectsClient(t *testing.T) {
	c := test.NewTest()

	user := testUsersClient_Create(t, c, msgraph.User{
		AccountEnabled:    utils.BoolPtr(true),
		DisplayName:       utils.StringPtr("test-user"),
		MailNickname:      utils.StringPtr(fmt.Sprintf("test-user-directoryobject-%s", c.RandomString)),
		UserPrincipalName: utils.StringPtr(fmt.Sprintf("test-user-directoryobject-%s@%s", c.RandomString, c.Connection.DomainName)),
		PasswordProfile: &msgraph.UserPasswordProfile{
			Password: utils.StringPtr(fmt.Sprintf("IrPa55w0rd%s", c.RandomString)),
		},
	})

	newGroup1 := msgraph.Group{
		DisplayName:     utils.StringPtr("test-group-directoryobject-member"),
		MailEnabled:     utils.BoolPtr(false),
		MailNickname:    utils.StringPtr(fmt.Sprintf("test-group-directoryobject-member-%s", c.RandomString)),
		SecurityEnabled: utils.BoolPtr(true),
	}
	newGroup1.Members = &msgraph.Members{user.DirectoryObject}
	group1 := testGroupsClient_Create(t, c, newGroup1)

	newGroup2 := msgraph.Group{
		DisplayName:     utils.StringPtr("test-group-directoryobject"),
		MailEnabled:     utils.BoolPtr(false),
		MailNickname:    utils.StringPtr(fmt.Sprintf("test-group-directoryobject-%s", c.RandomString)),
		SecurityEnabled: utils.BoolPtr(true),
		Members: &msgraph.Members{
			group1.DirectoryObject,
			user.DirectoryObject,
		},
	}
	group2 := testGroupsClient_Create(t, c, newGroup2)

	testDirectoryObjectsClient_Get(t, c, *user.ID)
	testDirectoryObjectsClient_Get(t, c, *group1.ID)
	testDirectoryObjectsClient_GetMemberGroups(t, c, *user.ID, true, []string{*group1.ID, *group2.ID})
	testDirectoryObjectsClient_GetMemberObjects(t, c, *group1.ID, true, []string{*group2.ID})
	testDirectoryObjectsClient_GetByIds(t, c, []string{*group1.ID, *group2.ID, *user.ID}, []string{odata.ShortTypeGroup})
	testDirectoryObjectsClient_Delete(t, c, *group1.ID)
}

func testDirectoryObjectsClient_Get(t *testing.T, c *test.Test, id string) (directoryObject *msgraph.DirectoryObject) {
	directoryObject, status, err := c.DirectoryObjectsClient.Get(c.Connection.Context, id, odata.Query{})
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

func testDirectoryObjectsClient_GetByIds(t *testing.T, c *test.Test, ids []string, types []odata.ShortType) (directoryObjects *[]msgraph.DirectoryObject) {
	directoryObjects, status, err := c.DirectoryObjectsClient.GetByIds(c.Connection.Context, ids, types)
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

func testDirectoryObjectsClient_GetMemberGroups(t *testing.T, c *test.Test, id string, securityEnabledOnly bool, expected []string) (directoryObjects *[]msgraph.DirectoryObject) {
	directoryObjects, status, err := c.DirectoryObjectsClient.GetMemberGroups(c.Connection.Context, id, securityEnabledOnly)
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

func testDirectoryObjectsClient_GetMemberObjects(t *testing.T, c *test.Test, id string, securityEnabledOnly bool, expected []string) (directoryObjects *[]msgraph.DirectoryObject) {
	directoryObjects, status, err := c.DirectoryObjectsClient.GetMemberObjects(c.Connection.Context, id, securityEnabledOnly)
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

func testDirectoryObjectsClient_Delete(t *testing.T, c *test.Test, id string) {
	status, err := c.DirectoryObjectsClient.Delete(c.Connection.Context, id)
	if err != nil {
		t.Fatalf("DirectoryObjectsClient.Delete(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("DirectoryObjectsClient.Delete(): invalid status: %d", status)
	}
}
