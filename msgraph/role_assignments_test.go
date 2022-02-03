package msgraph_test

import (
	"fmt"
	"testing"

	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/msgraph"
	"github.com/manicminer/hamilton/odata"
)

func TestRoleAssignmentsClient(t *testing.T) {
	c := test.NewTest(t)
	defer c.CancelFunc()

	roleDefinition := testRoleDefinitionsClient_Create(t, c, msgraph.UnifiedRoleDefinition{
		Description: msgraph.NullableString("testing custom role assignment"),
		DisplayName: utils.StringPtr("Test Assignor"),
		IsEnabled:   utils.BoolPtr(true),
		RolePermissions: &[]msgraph.UnifiedRolePermission{
			{
				AllowedResourceActions: &[]string{
					"microsoft.directory/groups/allProperties/read",
				},
			},
		},
		Version: utils.StringPtr("1.5"),
	})

	user := testUsersClient_Create(t, c, msgraph.User{
		AccountEnabled:    utils.BoolPtr(true),
		DisplayName:       utils.StringPtr("test-user"),
		MailNickname:      utils.StringPtr(fmt.Sprintf("test-user-%s", c.RandomString)),
		UserPrincipalName: utils.StringPtr(fmt.Sprintf("test-user-%s@%s", c.RandomString, c.Connection.DomainName)),
		PasswordProfile: &msgraph.UserPasswordProfile{
			Password: utils.StringPtr(fmt.Sprintf("IrPa55w0rd%s", c.RandomString)),
		},
	})

	roleAssignment := testRoleAssignmentsClient_Create(t, c, msgraph.UnifiedRoleAssignment{
		DirectoryScopeId: utils.StringPtr("/"),
		PrincipalId:      user.ID,
		RoleDefinitionId: roleDefinition.ID,
	})

	testRoleAssignmentsClient_Get(t, c, *roleAssignment.ID)
	testRoleAssignmentsClient_List(t, c, odata.Query{Filter: fmt.Sprintf("roleDefinitionId eq '%s'", *roleDefinition.ID)})
	testRoleAssignmentsClient_Delete(t, c, *roleAssignment.ID)
	testRoleDefinitionsClient_Delete(t, c, *roleDefinition.ID)
	testUsersClient_Delete(t, c, *user.ID)
	testUsersClient_DeletePermanently(t, c, *user.ID)
}

func testRoleAssignmentsClient_Create(t *testing.T, c *test.Test, r msgraph.UnifiedRoleAssignment) (roleAssignment *msgraph.UnifiedRoleAssignment) {
	roleAssignment, status, err := c.RoleAssignmentsClient.Create(c.Context, r)
	if err != nil {
		t.Fatalf("RoleAssignmentsClient.Create(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("RoleAssignmentsClient.Create(): invalid status: %d", status)
	}
	if roleAssignment == nil {
		t.Fatal("RoleAssignmentsClient.Create(): roleAssignment was nil")
	}
	if roleAssignment.ID == nil {
		t.Fatal("RoleAssignmentsClient.Create(): roleAssignment.ID was nil")
	}
	return
}

func testRoleAssignmentsClient_List(t *testing.T, c *test.Test, query odata.Query) (roleAssignments *[]msgraph.UnifiedRoleAssignment) {
	roleAssignments, _, err := c.RoleAssignmentsClient.List(c.Context, query)
	if err != nil {
		t.Fatalf("RoleAssignmentsClient.List(): %v", err)
	}
	if roleAssignments == nil {
		t.Fatal("RoleAssignmentsClient.List(): roleAssignments was nil")
	}
	return
}

func testRoleAssignmentsClient_Get(t *testing.T, c *test.Test, id string) (roleAssignment *msgraph.UnifiedRoleAssignment) {
	roleAssignment, status, err := c.RoleAssignmentsClient.Get(c.Context, id, odata.Query{})
	if err != nil {
		t.Fatalf("RoleAssignmentsClient.Get(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("RoleAssignmentsClient.Get(): invalid status: %d", status)
	}
	if roleAssignment == nil {
		t.Fatal("RoleAssignmentsClient.Get(): roleAssignment was nil")
	}
	return
}

func testRoleAssignmentsClient_Delete(t *testing.T, c *test.Test, id string) {
	status, err := c.RoleAssignmentsClient.Delete(c.Context, id)
	if err != nil {
		t.Fatalf("RoleAssignmentsClient.Delete(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("RoleAssignmentsClient.Delete(): invalid status: %d", status)
	}
}
