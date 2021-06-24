package msgraph_test

import (
	"strings"
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
	rs := test.RandomString()
	// set up directory role templates test client
	dirRoleTemplatesClient := DirectoryRoleTemplatesClientTest{
		connection:   test.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: rs,
	}
	dirRoleTemplatesClient.client = msgraph.NewDirectoryRoleTemplatesClient(dirRoleTemplatesClient.connection.AuthConfig.TenantID)
	dirRoleTemplatesClient.client.BaseClient.Authorizer = dirRoleTemplatesClient.connection.Authorizer

	// set up directory roles test client
	dirRolesClient := DirectoryRolesClientTest{
		connection:   test.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: rs,
	}
	dirRolesClient.client = msgraph.NewDirectoryRolesClient(dirRolesClient.connection.AuthConfig.TenantID)
	dirRolesClient.client.BaseClient.Authorizer = dirRolesClient.connection.Authorizer

	// list all directory roles available in the tenant
	directoryRoleTemplates := testDirectoryRoleTemplatesClient_List(t, dirRoleTemplatesClient)
	testDirectoryRoleTemplatesClient_Get(t, dirRoleTemplatesClient, *(*directoryRoleTemplates)[0].ID)

	// activate a directory role in the tenant using role template id if not already activated
	// https://docs.microsoft.com/en-us/azure/active-directory/roles/permissions-reference
	var globalReaderRoleId string
	for _, template := range *directoryRoleTemplates {
		if strings.EqualFold(*template.DisplayName, "Global Reader") {
			globalReaderRoleId = *template.ID
		}
	}
	testDirectoryRolesClient_Activate(t, dirRolesClient, globalReaderRoleId)
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

func testDirectoryRolesClient_Activate(t *testing.T, c DirectoryRolesClientTest, roleTemplateId string) (directoryRole *msgraph.DirectoryRole) {
	// list all activated directory roles in the tenant
	directoryRoles, _, err := c.client.List(c.connection.Context)
	if err != nil {
		t.Fatalf("DirectoryRolesClient.List(): %v", err)
	}
	if directoryRoles == nil {
		t.Fatal("DirectoryRolesClient.List(): directoryRoles was nil")
	}

	// helper function to find activate directory role by role template id
	// api does not support retrieving directory role by role template id; it does not support the OData Query Parameters
	findDirRoleByRoleTemplateId := func(directoryRoles []msgraph.DirectoryRole, roleTemplatedId string) *msgraph.DirectoryRole {
		for _, dirRole := range directoryRoles {
			if dirRole.RoleTemplateId != nil && (*dirRole.RoleTemplateId) == roleTemplateId {
				return &dirRole
			}
		}
		return nil
	}

	// attempt to activate directory role if not already present in the directory
	if dirRole := findDirRoleByRoleTemplateId(*directoryRoles, roleTemplateId); dirRole == nil {
		directoryRole, status, err := c.client.Activate(c.connection.Context, roleTemplateId)
		if err != nil {
			t.Fatalf("DirectoryRolesClient.Activate(): %v", err)
		}
		if status < 200 || status >= 300 {
			t.Fatalf("DirectoryRolesClient.Activate(): invalid status: %d", status)
		}
		if directoryRole == nil {
			t.Fatal("DirectoryRolesClient.Activate(): directoryRole was nil")
		}
	}

	// attempt to activate directory role a second time to test the API error handling
	directoryRole, status, err := c.client.Activate(c.connection.Context, roleTemplateId)
	if err != nil {
		t.Fatalf("DirectoryRolesClient.Activate() [attempt 2]: %v", err)
	}
	if (status < 200 || status >= 300) && (status < 400 || status >= 500) {
		t.Fatalf("DirectoryRolesClient.Activate() [attempt 2]: invalid status: %d", status)
	}
	if directoryRole == nil {
		t.Fatal("DirectoryRolesClient.Activate() [attempt 2]: directoryRole was nil")
	}
	return
}
