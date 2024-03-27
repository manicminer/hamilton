package msgraph_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/msgraph"
)

func TestRoleManagementPolicyClient(t *testing.T) {
	c := test.NewTest(t)
	defer c.CancelFunc()

	testRoleManagementPolicyClient_DirectoryRoles(t, c)
	testRoleManagementPolicyClient_GroupRoles(t, c)
}

func testRoleManagementPolicyClient_DirectoryRoles(t *testing.T, c *test.Test) {
	policies := testRoleManagementPolicyClient_List(t, c, odata.Query{
		Filter: "scopeId eq '/' and scopeType eq 'DirectoryRole'",
	})
	policy := testRoleManagementPolicyClient_Get(t, c, *(*policies)[0].ID)

	assignments := testRoleManagementPolicyAssignmentClient_List(t, c, odata.Query{
		Filter: "scopeId eq '/' and scopeType eq 'DirectoryRole'",
	})
	testRoleManagementPolicyAssignmentClient_Get(t, c, *(*assignments)[0].ID)

	rules := testRoleManagementPolicyRuleClient_List(t, c, *policy.ID)
	for _, rule := range *rules {
		if *rule.ID == "Expiration_Admin_Eligibility" {
			rule.MaximumDuration = utils.StringPtr("P180D")
		}
	}
	testRoleManagementPolicyClient_Update(t, c, msgraph.UnifiedRoleManagementPolicy{
		ID:    policy.ID,
		Rules: rules,
	})

	rule := testRoleManagementPolicyRuleClient_Get(t, c, *policy.ID, "Expiration_Admin_Eligibility")
	testRoleManagementPolicyRuleClient_Update(t, c, *policy.ID, msgraph.UnifiedRoleManagementPolicyRule{
		ID:              rule.ID,
		ODataType:       rule.ODataType,
		MaximumDuration: utils.StringPtr("PT0S"), // Revert to default for future tests
	})
}

func testRoleManagementPolicyClient_GroupRoles(t *testing.T, c *test.Test) {
	pimGroup := testGroupsClient_Create(t, c, msgraph.Group{
		DisplayName:     utils.StringPtr("test-pim-group"),
		MailEnabled:     utils.BoolPtr(false),
		MailNickname:    utils.StringPtr(fmt.Sprintf("test-pim-group-%s", c.RandomString)),
		SecurityEnabled: utils.BoolPtr(true),
	})
	defer testGroupsClient_Delete(t, c, *pimGroup.ID())

	policies := testRoleManagementPolicyClient_List(t, c, odata.Query{
		Filter: fmt.Sprintf("scopeId eq '%s' and scopeType eq 'Group'", *pimGroup.ID()),
	})
	policy := testRoleManagementPolicyClient_Get(t, c, *(*policies)[0].ID)

	assignments := testRoleManagementPolicyAssignmentClient_List(t, c, odata.Query{
		Filter: fmt.Sprintf("scopeId eq '%s' and scopeType eq 'Group'", *pimGroup.ID()),
	})
	testRoleManagementPolicyAssignmentClient_Get(t, c, *(*assignments)[0].ID)

	rules := testRoleManagementPolicyRuleClient_List(t, c, *policy.ID)
	for _, rule := range *rules {
		if *rule.ID == "Expiration_Admin_Eligibility" {
			rule.MaximumDuration = utils.StringPtr("P730D")
		}
	}
	testRoleManagementPolicyClient_Update(t, c, msgraph.UnifiedRoleManagementPolicy{
		ID:    policy.ID,
		Rules: rules,
	})

	testRoleManagementPolicyRuleClient_List(t, c, *policy.ID)
	rule := testRoleManagementPolicyRuleClient_Get(t, c, *policy.ID, "Expiration_Admin_Eligibility")
	testRoleManagementPolicyRuleClient_Update(t, c, *policy.ID, msgraph.UnifiedRoleManagementPolicyRule{
		ID:              rule.ID,
		ODataType:       rule.ODataType,
		MaximumDuration: utils.StringPtr("P365D"),
	})
}

func testRoleManagementPolicyClient_List(t *testing.T, c *test.Test, query odata.Query) (policies *[]msgraph.UnifiedRoleManagementPolicy) {
	policies, status, err := c.RoleManagementPolicyClient.List(c.Context, query)
	if err != nil {
		t.Fatalf("RoleManagementPolicyClient.List(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("RoleManagementPolicyClient.List(): invalid status: %d", status)
	}
	if policies == nil {
		t.Fatal("RoleManagementPolicyClient.List(): response was nil")
	}
	return
}

func testRoleManagementPolicyClient_Get(t *testing.T, c *test.Test, id string) (policy *msgraph.UnifiedRoleManagementPolicy) {
	policy, status, err := c.RoleManagementPolicyClient.Get(c.Context, id)
	if err != nil {
		t.Fatalf("RoleManagementPolicyClient.Get(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("RoleManagementPolicyClient.Get(): invalid status: %d", status)
	}
	if policy == nil {
		t.Fatal("RoleManagementPolicyClient.Get(): response was nil")
	}
	return
}

func testRoleManagementPolicyClient_Update(t *testing.T, c *test.Test, policy msgraph.UnifiedRoleManagementPolicy) {
	status, err := c.RoleManagementPolicyClient.Update(c.Context, policy)
	if err != nil {
		t.Fatalf("RoleManagementPolicyClient.Update(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("RoleManagementPolicyClient.Update(): invalid status: %d", status)
	}
}
