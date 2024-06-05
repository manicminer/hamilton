package msgraph_test

import (
	"testing"

	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/msgraph"
)

func testPrivilegedAccessGroupAssignmentScheduleInstancesClient_List(t *testing.T, c *test.Test, query odata.Query) (instances *[]msgraph.PrivilegedAccessGroupAssignmentScheduleInstance) {
	instances, status, err := c.PrivilegedAccessGroupAssignmentScheduleInstancesClient.List(c.Context, query)
	if err != nil {
		t.Fatalf("PrivilegedAccessGroupAssignmentScheduleInstancesClient.List(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("PrivilegedAccessGroupAssignmentScheduleInstancesClient.List(): invalid status: %d", status)
	}
	if instances == nil {
		t.Fatal("PrivilegedAccessGroupAssignmentScheduleInstancesClient.List(): PrivilegedAccessGroupAssignmentSchedule was nil")
	}
	if len(*instances) == 0 {
		t.Fatal("PrivilegedAccessGroupAssignmentScheduleInstancesClient.List(): Returned zero results")
	}
	return
}

func testPrivilegedAccessGroupAssignmentScheduleInstancesClient_Get(t *testing.T, c *test.Test, id string) (request *msgraph.PrivilegedAccessGroupAssignmentScheduleInstance) {
	request, status, err := c.PrivilegedAccessGroupAssignmentScheduleInstancesClient.Get(c.Context, id)
	if err != nil {
		t.Fatalf("PrivilegedAccessGroupAssignmentScheduleInstancesClient.Get(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("PrivilegedAccessGroupAssignmentScheduleInstancesClient.Get(): invalid status: %d", status)
	}
	if request == nil {
		t.Fatal("PrivilegedAccessGroupAssignmentScheduleInstancesClient.Get(): PrivilegedAccessGroupAssignmentSchedule was nil")
	}
	return
}
