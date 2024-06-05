package msgraph_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/msgraph"
)

func TestEntitlementRoleAssignmentsClient(t *testing.T) {
	c := test.NewTest(t)
	defer c.CancelFunc()

	roleDefinitions := testEntitlementRoleDefinitionsClient_List(t, c)

	roleDefinition := (*roleDefinitions)[0]

	user := testUsersClient_Create(t, c, msgraph.User{
		AccountEnabled:    utils.BoolPtr(true),
		DisplayName:       utils.StringPtr("test-user"),
		MailNickname:      utils.StringPtr(fmt.Sprintf("test-user-%s", c.RandomString)),
		UserPrincipalName: utils.StringPtr(fmt.Sprintf("test-user-%s@%s", c.RandomString, c.Connections["default"].DomainName)),
		PasswordProfile: &msgraph.UserPasswordProfile{
			Password: utils.StringPtr(fmt.Sprintf("IrPa55w0rd%s", c.RandomString)),
		},
	})

	accessPackageCatalog := testAccessPackageCatalogClient_Create(t, c, msgraph.AccessPackageCatalog{
		DisplayName:         utils.StringPtr(fmt.Sprintf("test-catalog-%s", c.RandomString)),
		CatalogType:         msgraph.AccessPackageCatalogTypeUserManaged,
		State:               msgraph.AccessPackageCatalogStatePublished,
		Description:         utils.StringPtr("Test Access Catalog"),
		IsExternallyVisible: utils.BoolPtr(false),
	})

	roleAssignment := testEntitlementRoleAssignmentsClient_Create(t, c, msgraph.UnifiedRoleAssignment{
		DirectoryScopeId: utils.StringPtr("/"),
		PrincipalId:      user.ID(),
		RoleDefinitionId: roleDefinition.ID(),
		AppScopeId:       utils.StringPtr("/AccessPackageCatalog/" + *accessPackageCatalog.ID),
	})

	testEntitlementRoleAssignmentsClient_Get(t, c, *roleAssignment.ID())
	testEntitlementRoleAssignmentsClient_List(t, c, odata.Query{Filter: fmt.Sprintf("roleDefinitionId eq '%s'", *roleDefinition.ID())})
	testEntitlementRoleAssignmentsClient_Delete(t, c, *roleAssignment.ID())
	testAccessPackageCatalogClient_Delete(t, c, *accessPackageCatalog.ID)
	testUsersClient_Delete(t, c, *user.ID())
	testUsersClient_DeletePermanently(t, c, *user.ID())
}

func testEntitlementRoleAssignmentsClient_Create(t *testing.T, c *test.Test, r msgraph.UnifiedRoleAssignment) (roleAssignment *msgraph.UnifiedRoleAssignment) {
	roleAssignment, status, err := c.EntitlementRoleAssignmentsClient.Create(c.Context, r)
	if err != nil {
		t.Fatalf("EntitlementRoleAssignmentsClient.Create(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("EntitlementRoleAssignmentsClient.Create(): invalid status: %d", status)
	}
	if roleAssignment == nil {
		t.Fatal("EntitlementRoleAssignmentsClient.Create(): roleAssignment was nil")
	}
	if roleAssignment.ID() == nil {
		t.Fatal("EntitlementRoleAssignmentsClient.Create(): roleAssignment.ID was nil")
	}
	return
}

func testEntitlementRoleAssignmentsClient_List(t *testing.T, c *test.Test, query odata.Query) (roleAssignments *[]msgraph.UnifiedRoleAssignment) {
	roleAssignments, _, err := c.EntitlementRoleAssignmentsClient.List(c.Context, query)
	if err != nil {
		t.Fatalf("EntitlementRoleAssignmentsClient.List(): %v", err)
	}
	if roleAssignments == nil {
		t.Fatal("EntitlementRoleAssignmentsClient.List(): roleAssignments was nil")
	}
	return
}

func testEntitlementRoleAssignmentsClient_Get(t *testing.T, c *test.Test, id string) (roleAssignment *msgraph.UnifiedRoleAssignment) {
	roleAssignment, status, err := c.EntitlementRoleAssignmentsClient.Get(c.Context, id, odata.Query{})
	if err != nil {
		t.Fatalf("EntitlementRoleAssignmentsClient.Get(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("EntitlementRoleAssignmentsClient.Get(): invalid status: %d", status)
	}
	if roleAssignment == nil {
		t.Fatal("EntitlementRoleAssignmentsClient.Get(): roleAssignment was nil")
	}
	return
}

func testEntitlementRoleAssignmentsClient_Delete(t *testing.T, c *test.Test, id string) {
	status, err := c.EntitlementRoleAssignmentsClient.Delete(c.Context, id)
	if err != nil {
		t.Fatalf("EntitlementRoleAssignmentsClient.Delete(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("EntitlementRoleAssignmentsClient.Delete(): invalid status: %d", status)
	}
}
