package msgraph_test

import (
	"testing"

	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/msgraph"
)

func TestEntitlementRoleDefinitionsClient(t *testing.T) {
	c := test.NewTest(t)
	defer c.CancelFunc()

	roleDefinitions := testEntitlementRoleDefinitionsClient_List(t, c)

	role := (*roleDefinitions)[0]
	testEntitlementRoleDefinitionsClient_Get(t, c, *role.ID())
}

func testEntitlementRoleDefinitionsClient_List(t *testing.T, c *test.Test) (roleDefinitions *[]msgraph.UnifiedRoleDefinition) {
	roleDefinitions, _, err := c.EntitlementRoleDefinitionsClient.List(c.Context, odata.Query{})
	if err != nil {
		t.Fatalf("EntitlementRoleDefinitionsClient.List(): %v", err)
	}
	if roleDefinitions == nil {
		t.Fatal("EntitlementRoleDefinitionsClient.List(): roleDefinitions was nil")
	}
	return
}

func testEntitlementRoleDefinitionsClient_Get(t *testing.T, c *test.Test, id string) (roleDefinition *msgraph.UnifiedRoleDefinition) {
	roleDefinition, status, err := c.EntitlementRoleDefinitionsClient.Get(c.Context, id, odata.Query{})
	if err != nil {
		t.Fatalf("EntitlementRoleDefinitionsClient.Get(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("EntitlementRoleDefinitionsClient.Get(): invalid status: %d", status)
	}
	if roleDefinition == nil {
		t.Fatal("EntitlementRoleDefinitionsClient.Get(): roleDefinition was nil")
	}
	return
}
