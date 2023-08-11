package msgraph_test

import (
	"testing"

	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/msgraph"
)

func testPrivilegedAccessGroupEligibilityScheduleRequestsClient_Create(t *testing.T, c *test.Test, r msgraph.PrivilegedAccessGroupEligibilityScheduleRequest) (request *msgraph.PrivilegedAccessGroupEligibilityScheduleRequest) {
	request, status, err := c.PrivilegedAccessGroupEligibilityScheduleRequestClient.Create(c.Context, r)
	if err != nil {
		t.Fatalf("PrivilegedAccessGroupEligibilityScheduleRequestClient.Create(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("PrivilegedAccessGroupEligibilityScheduleRequestClient.Create(): invalid status: %d", status)
	}
	if request == nil {
		t.Fatal("PrivilegedAccessGroupEligibilityScheduleRequestClient.Create(): PrivilegedAccessGroupEligibilityScheduleRequest was nil")
	}
	if request.ID == nil {
		t.Fatal("PrivilegedAccessGroupEligibilityScheduleRequestClient.Create(): PrivilegedAccessGroupEligibilityScheduleRequest.ID was nil")
	}
	return
}

func testPrivilegedAccessGroupEligibilityScheduleRequestsClient_List(t *testing.T, c *test.Test) (requests *[]msgraph.PrivilegedAccessGroupEligibilityScheduleRequest) {
	requests, status, err := c.PrivilegedAccessGroupEligibilityScheduleRequestClient.List(c.Context, odata.Query{})
	if err != nil {
		t.Fatalf("PrivilegedAccessGroupEligibilityScheduleRequestClient.List(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("PrivilegedAccessGroupEligibilityScheduleRequestClient.List(): invalid status: %d", status)
	}
	if requests == nil {
		t.Fatal("PrivilegedAccessGroupEligibilityScheduleRequestClient.List(): PrivilegedAccessGroupEligibilityScheduleRequest was nil")
	}
	return
}

func testPrivilegedAccessGroupEligibilityScheduleRequestsClient_Get(t *testing.T, c *test.Test, id string) (request *msgraph.PrivilegedAccessGroupEligibilityScheduleRequest) {
	request, status, err := c.PrivilegedAccessGroupEligibilityScheduleRequestClient.Get(c.Context, id)
	if err != nil {
		t.Fatalf("PrivilegedAccessGroupEligibilityScheduleRequestClient.Get(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("PrivilegedAccessGroupEligibilityScheduleRequestClient.Get(): invalid status: %d", status)
	}
	if request == nil {
		t.Fatal("PrivilegedAccessGroupEligibilityScheduleRequestClient.Get(): PrivilegedAccessGroupEligibilityScheduleRequest was nil")
	}
	return
}

func testPrivilegedAccessGroupEligibilityScheduleRequestsClient_Cancel(t *testing.T, c *test.Test, id string) {
	status, err := c.PrivilegedAccessGroupEligibilityScheduleRequestClient.Cancel(c.Context, id)
	if err != nil {
		t.Fatalf("PrivilegedAccessGroupEligibilityScheduleRequestClient.Cancel(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("PrivilegedAccessGroupEligibilityScheduleRequestClient.Cancel(): invalid status: %d", status)
	}
}
