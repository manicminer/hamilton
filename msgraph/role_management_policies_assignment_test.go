package msgraph_test

import (
	"testing"

	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/msgraph"
)

func testRoleManagementPolicyAssignmentClient_List(t *testing.T, c *test.Test, query odata.Query) (rules *[]msgraph.UnifiedRoleManagementPolicyAssignment) {
	rules, status, err := c.RoleManagementPolicyAssignmentClient.List(c.Context, query)
	if err != nil {
		t.Fatalf("RoleManagementPolicyAssignmentClient.List(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("RoleManagementPolicyAssignmentClient.List(): invalid status: %d", status)
	}
	if rules == nil {
		t.Fatal("RoleManagementPolicyAssignmentClient.List(): response was nil")
	}
	return
}

func testRoleManagementPolicyAssignmentClient_Get(t *testing.T, c *test.Test, assignmentId string) (rule *msgraph.UnifiedRoleManagementPolicyAssignment) {
	rule, status, err := c.RoleManagementPolicyAssignmentClient.Get(c.Context, assignmentId)
	if err != nil {
		t.Fatalf("RoleManagementPolicyAssignmentClient.Get(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("RoleManagementPolicyAssignmentClient.Get(): invalid status: %d", status)
	}
	if rule == nil {
		t.Fatal("RoleManagementPolicyAssignmentClient.Get(): response was nil")
	}
	return
}
