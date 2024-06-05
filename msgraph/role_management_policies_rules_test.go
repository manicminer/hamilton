package msgraph_test

import (
	"testing"

	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/msgraph"
)

func testRoleManagementPolicyRuleClient_List(t *testing.T, c *test.Test, policyId string) (rules *[]msgraph.UnifiedRoleManagementPolicyRule) {
	rules, status, err := c.RoleManagementPolicyRuleClient.List(c.Context, policyId)
	if err != nil {
		t.Fatalf("RoleManagementPolicyRuleClient.List(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("RoleManagementPolicyRuleClient.List(): invalid status: %d", status)
	}
	if rules == nil {
		t.Fatal("RoleManagementPolicyRuleClient.List(): response was nil")
	}
	return
}

func testRoleManagementPolicyRuleClient_Get(t *testing.T, c *test.Test, policyId string, ruleId string) (rule *msgraph.UnifiedRoleManagementPolicyRule) {
	rule, status, err := c.RoleManagementPolicyRuleClient.Get(c.Context, policyId, ruleId)
	if err != nil {
		t.Fatalf("RoleManagementPolicyRuleClient.Get(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("RoleManagementPolicyRuleClient.Get(): invalid status: %d", status)
	}
	if rule == nil {
		t.Fatal("RoleManagementPolicyRuleClient.Get(): response was nil")
	}
	return
}

func testRoleManagementPolicyRuleClient_Update(t *testing.T, c *test.Test, policyId string, rule msgraph.UnifiedRoleManagementPolicyRule) {
	status, err := c.RoleManagementPolicyRuleClient.Update(c.Context, policyId, rule)
	if err != nil {
		t.Fatalf("RoleManagementPolicyRuleClient.Update(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("RoleManagementPolicyRuleClient.Update(): invalid status: %d", status)
	}
}
