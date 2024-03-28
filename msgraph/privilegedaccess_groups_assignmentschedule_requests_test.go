package msgraph_test

import (
	"testing"

	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/msgraph"
)

func testPrivilegedAccessGroupAssignmentScheduleRequestsClient_Create(t *testing.T, c *test.Test, r msgraph.PrivilegedAccessGroupAssignmentScheduleRequest) (request *msgraph.PrivilegedAccessGroupAssignmentScheduleRequest) {
	request, status, err := c.PrivilegedAccessGroupAssignmentScheduleRequestsClient.Create(c.Context, r)
	if err != nil {
		t.Fatalf("PrivilegedAccessGroupAssignmentScheduleRequestsClient.Create(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("PrivilegedAccessGroupAssignmentScheduleRequestsClient.Create(): invalid status: %d", status)
	}
	if request == nil {
		t.Fatal("PrivilegedAccessGroupAssignmentScheduleRequestsClient.Create(): PrivilegedAccessGroupAssignmentScheduleRequest was nil")
	}
	if request.ID == nil {
		t.Fatal("PrivilegedAccessGroupAssignmentScheduleRequestsClient.Create(): PrivilegedAccessGroupAssignmentScheduleRequest.ID was nil")
	}
	return
}

func testPrivilegedAccessGroupAssignmentScheduleRequestsClient_List(t *testing.T, c *test.Test) (requests *[]msgraph.PrivilegedAccessGroupAssignmentScheduleRequest) {
	requests, status, err := c.PrivilegedAccessGroupAssignmentScheduleRequestsClient.List(c.Context, odata.Query{})
	if err != nil {
		t.Fatalf("PrivilegedAccessGroupAssignmentScheduleRequestsClient.List(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("PrivilegedAccessGroupAssignmentScheduleRequestsClient.List(): invalid status: %d", status)
	}
	if requests == nil {
		t.Fatal("PrivilegedAccessGroupAssignmentScheduleRequestsClient.List(): PrivilegedAccessGroupAssignmentScheduleRequest was nil")
	}
	return
}

func testPrivilegedAccessGroupAssignmentScheduleRequestsClient_Get(t *testing.T, c *test.Test, id string) (request *msgraph.PrivilegedAccessGroupAssignmentScheduleRequest) {
	request, status, err := c.PrivilegedAccessGroupAssignmentScheduleRequestsClient.Get(c.Context, id)
	if err != nil {
		t.Fatalf("PrivilegedAccessGroupAssignmentScheduleRequestsClient.Get(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("PrivilegedAccessGroupAssignmentScheduleRequestsClient.Get(): invalid status: %d", status)
	}
	if request == nil {
		t.Fatal("PrivilegedAccessGroupAssignmentScheduleRequestsClient.Get(): PrivilegedAccessGroupAssignmentScheduleRequest was nil")
	}
	return
}

func testPrivilegedAccessGroupAssignmentScheduleRequestsClient_Cancel(t *testing.T, c *test.Test, id string) {
	status, err := c.PrivilegedAccessGroupAssignmentScheduleRequestsClient.Cancel(c.Context, id)
	if err != nil {
		t.Fatalf("PrivilegedAccessGroupAssignmentScheduleRequestsClient.Cancel(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("PrivilegedAccessGroupAssignmentScheduleRequestsClient.Cancel(): invalid status: %d", status)
	}
}
