package msgraph_test

import (
	"testing"

	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/msgraph"
)

func testPrivilegedAccessGroupEligibilityScheduleInstancesClient_List(t *testing.T, c *test.Test, query odata.Query) (instances *[]msgraph.PrivilegedAccessGroupEligibilityScheduleInstance) {
	instances, status, err := c.PrivilegedAccessGroupEligibilityScheduleInstancesClient.List(c.Context, query)
	if err != nil {
		t.Fatalf("PrivilegedAccessGroupEligibilityScheduleInstancesClient.InstancesList(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("PrivilegedAccessGroupEligibilityScheduleInstancesClient.InstancesList(): invalid status: %d", status)
	}
	if instances == nil {
		t.Fatal("PrivilegedAccessGroupEligibilityScheduleInstancesClient.InstancesList(): PrivilegedAccessGroupEligibilitySchedule was nil")
	}
	if len(*instances) == 0 {
		t.Fatal("PrivilegedAccessGroupEligibilityScheduleInstancesClient.List(): Returned zero results")
	}
	return
}

func testPrivilegedAccessGroupEligibilityScheduleInstancesClient_Get(t *testing.T, c *test.Test, id string) (request *msgraph.PrivilegedAccessGroupEligibilityScheduleInstance) {
	request, status, err := c.PrivilegedAccessGroupEligibilityScheduleInstancesClient.Get(c.Context, id)
	if err != nil {
		t.Fatalf("PrivilegedAccessGroupEligibilityScheduleInstancesClient.InstancesGet(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("PrivilegedAccessGroupEligibilityScheduleInstancesClient.InstancesGet(): invalid status: %d", status)
	}
	if request == nil {
		t.Fatal("PrivilegedAccessGroupEligibilityScheduleInstancesClient.InstancesGet(): PrivilegedAccessGroupEligibilitySchedule was nil")
	}
	return
}
