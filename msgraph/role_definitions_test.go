package msgraph_test

import (
	"testing"

	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/msgraph"
)

func TestRoleDefinitionsClient(t *testing.T) {
	c := test.NewTest(t)
	defer c.CancelFunc()

	roleDefinition := testRoleDefinitionsClient_Create(t, c, msgraph.UnifiedRoleDefinition{
		Description: msgraph.NullableString("my test role definition"),
		DisplayName: utils.StringPtr("test-pontificator"),
		IsEnabled:   utils.BoolPtr(true),
		RolePermissions: &[]msgraph.UnifiedRolePermission{
			{
				AllowedResourceActions: &[]string{
					"microsoft.directory/applications/allProperties/read",
					"microsoft.directory/applications/synchronization/standard/read",
					"microsoft.directory/groups/allProperties/read",
				},
				//Condition: utils.StringPtr("@Subject.objectId Any_of @Resource.owners"), // not yet supported by API
			},
		},
		Version: utils.StringPtr("1.5"),
	})

	roleDefinition.Description = msgraph.NullableString("for the tests")
	testRoleDefinitionsClient_Update(t, c, *roleDefinition)
	testRoleDefinitionsClient_Get(t, c, *roleDefinition.ID())

	roleDefinitions := testRoleDefinitionsClient_List(t, c)
	role := (*roleDefinitions)[0]
	testRoleDefinitionsClient_Get(t, c, *role.ID())

	testRoleDefinitionsClient_Delete(t, c, *roleDefinition.ID())
}

func testRoleDefinitionsClient_Create(t *testing.T, c *test.Test, r msgraph.UnifiedRoleDefinition) (roleDefinition *msgraph.UnifiedRoleDefinition) {
	roleDefinition, status, err := c.RoleDefinitionsClient.Create(c.Context, r)
	if err != nil {
		t.Fatalf("RoleDefinitionsClient.Create(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("RoleDefinitionsClient.Create(): invalid status: %d", status)
	}
	if roleDefinition == nil {
		t.Fatal("RoleDefinitionsClient.Create(): roleDefinition was nil")
	}
	if roleDefinition.ID() == nil {
		t.Fatal("RoleDefinitionsClient.Create(): roleDefinition.ID was nil")
	}
	return
}

func testRoleDefinitionsClient_Update(t *testing.T, c *test.Test, r msgraph.UnifiedRoleDefinition) {
	status, err := c.RoleDefinitionsClient.Update(c.Context, r)
	if err != nil {
		t.Fatalf("RoleDefinitionsClient.Update(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("RoleDefinitionsClient.Update(): invalid status: %d", status)
	}
}

func testRoleDefinitionsClient_List(t *testing.T, c *test.Test) (roleDefinitions *[]msgraph.UnifiedRoleDefinition) {
	roleDefinitions, _, err := c.RoleDefinitionsClient.List(c.Context, odata.Query{})
	if err != nil {
		t.Fatalf("RoleDefinitionsClient.List(): %v", err)
	}
	if roleDefinitions == nil {
		t.Fatal("RoleDefinitionsClient.List(): roleDefinitions was nil")
	}
	return
}

func testRoleDefinitionsClient_Get(t *testing.T, c *test.Test, id string) (roleDefinition *msgraph.UnifiedRoleDefinition) {
	roleDefinition, status, err := c.RoleDefinitionsClient.Get(c.Context, id, odata.Query{})
	if err != nil {
		t.Fatalf("RoleDefinitionsClient.Get(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("RoleDefinitionsClient.Get(): invalid status: %d", status)
	}
	if roleDefinition == nil {
		t.Fatal("RoleDefinitionsClient.Get(): roleDefinition was nil")
	}
	return
}

func testRoleDefinitionsClient_Delete(t *testing.T, c *test.Test, id string) {
	status, err := c.RoleDefinitionsClient.Delete(c.Context, id)
	if err != nil {
		t.Fatalf("RoleDefinitionsClient.Delete(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("RoleDefinitionsClient.Delete(): invalid status: %d", status)
	}
}
