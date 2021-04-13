package msgraph_test

import (
	"testing"

	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/msgraph"
)

type DirectoryRoleTemplatesClientTest struct {
	connection   *test.Connection
	client       *msgraph.DirectoryRoleTemplatesClient
	randomString string
}

func TestDirectoryRoleTemplatesClient(t *testing.T) {
	c := DirectoryRoleTemplatesClientTest{
		connection:   test.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: test.RandomString(),
	}
	c.client = msgraph.NewDirectoryRoleTemplatesClient(c.connection.AuthConfig.TenantID)
	c.client.BaseClient.Authorizer = c.connection.Authorizer

	directoryRoleTemplates := testDirectoryRoleTemplatesClient_List(t, c)
	testDirectoryRoleTemplatesClient_Get(t, c, *(*directoryRoleTemplates)[0].ID)
}

func testDirectoryRoleTemplatesClient_List(t *testing.T, c DirectoryRoleTemplatesClientTest) (directoryRoleTemplates *[]msgraph.DirectoryRoleTemplate) {
	directoryRoleTemplates, _, err := c.client.List(c.connection.Context)
	if err != nil {
		t.Fatalf("DirectoryRoleTemplatesClient.List(): %v", err)
	}
	if directoryRoleTemplates == nil {
		t.Fatal("DirectoryRoleTemplatesClient.List(): directoryRoleTemplates was nil")
	}
	return
}

func testDirectoryRoleTemplatesClient_Get(t *testing.T, c DirectoryRoleTemplatesClientTest, id string) (directoryRoleTemplate *msgraph.DirectoryRoleTemplate) {
	directoryRoleTemplate, status, err := c.client.Get(c.connection.Context, id)
	if err != nil {
		t.Fatalf("DirectoryRoleTemplatesClient.Get(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("DirectoryRoleTemplatesClient.Get(): invalid status: %d", status)
	}
	if directoryRoleTemplate == nil {
		t.Fatal("DirectoryRoleTemplatesClient.Get(): directoryRoleTemplate was nil")
	}
	return
}
