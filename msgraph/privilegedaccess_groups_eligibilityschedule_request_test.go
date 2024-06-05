package msgraph_test

import (
	"testing"

	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/msgraph"
)

func testPrivilegedAccessGroupEligibilityScheduleRequestsClient_Create(t *testing.T, c *test.Test, r msgraph.PrivilegedAccessGroupEligibilityScheduleRequest) (request *msgraph.PrivilegedAccessGroupEligibilityScheduleRequest) {
	request, status, err := c.PrivilegedAccessGroupEligibilityScheduleRequestsClient.Create(c.Context, r)
	if err != nil {
		t.Fatalf("PrivilegedAccessGroupEligibilityScheduleRequestsClient.Create(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("PrivilegedAccessGroupEligibilityScheduleRequestsClient.Create(): invalid status: %d", status)
	}
	if request == nil {
		t.Fatal("PrivilegedAccessGroupEligibilityScheduleRequestsClient.Create(): PrivilegedAccessGroupEligibilityScheduleRequest was nil")
	}
	if request.ID == nil {
		t.Fatal("PrivilegedAccessGroupEligibilityScheduleRequestsClient.Create(): PrivilegedAccessGroupEligibilityScheduleRequest.ID was nil")
	}
	return
}

func testPrivilegedAccessGroupEligibilityScheduleRequestsClient_List(t *testing.T, c *test.Test) (requests *[]msgraph.PrivilegedAccessGroupEligibilityScheduleRequest) {
	requests, status, err := c.PrivilegedAccessGroupEligibilityScheduleRequestsClient.List(c.Context, odata.Query{})
	if err != nil {
		t.Fatalf("PrivilegedAccessGroupEligibilityScheduleRequestsClient.List(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("PrivilegedAccessGroupEligibilityScheduleRequestsClient.List(): invalid status: %d", status)
	}
	if requests == nil {
		t.Fatal("PrivilegedAccessGroupEligibilityScheduleRequestsClient.List(): PrivilegedAccessGroupEligibilityScheduleRequest was nil")
	}
	return
}

func testPrivilegedAccessGroupEligibilityScheduleRequestsClient_Get(t *testing.T, c *test.Test, id string) (request *msgraph.PrivilegedAccessGroupEligibilityScheduleRequest) {
	request, status, err := c.PrivilegedAccessGroupEligibilityScheduleRequestsClient.Get(c.Context, id)
	if err != nil {
		t.Fatalf("PrivilegedAccessGroupEligibilityScheduleRequestsClient.Get(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("PrivilegedAccessGroupEligibilityScheduleRequestsClient.Get(): invalid status: %d", status)
	}
	if request == nil {
		t.Fatal("PrivilegedAccessGroupEligibilityScheduleRequestsClient.Get(): PrivilegedAccessGroupEligibilityScheduleRequest was nil")
	}
	return
}

func testPrivilegedAccessGroupEligibilityScheduleRequestsClient_Cancel(t *testing.T, c *test.Test, id string) {
	status, err := c.PrivilegedAccessGroupEligibilityScheduleRequestsClient.Cancel(c.Context, id)
	if err != nil {
		t.Fatalf("PrivilegedAccessGroupEligibilityScheduleRequestsClient.Cancel(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("PrivilegedAccessGroupEligibilityScheduleRequestsClient.Cancel(): invalid status: %d", status)
	}
}
