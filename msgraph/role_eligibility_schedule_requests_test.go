package msgraph_test

import (
	"testing"
	"time"

	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/msgraph"
)

func TestRoleEligibilityScheduleRequestsClient(t *testing.T) {
	c := test.NewTest(t)
	defer c.CancelFunc()

	testUser := testUser_Create(t, c)

	roleEligibilityScheduleRequest := testRoleEligibilityScheduleRequestsClient_Create(t, c, msgraph.UnifiedRoleEligibilityScheduleRequest{
		Action:           utils.StringPtr(msgraph.UnifiedRoleScheduleRequestActionsAdminAssign),
		DirectoryScopeId: utils.StringPtr("/"),
		PrincipalId:      testUser.ID(),
		// ffd52fa5-98dc-465c-991d-fc073eb59f8f is the template id of the "Attribute Assignment Reader" built-in role
		RoleDefinitionId: utils.StringPtr("ffd52fa5-98dc-465c-991d-fc073eb59f8f"),
		ScheduleInfo: &msgraph.RequestSchedule{
			Expiration: &msgraph.ExpirationPattern{
				Type: utils.StringPtr(msgraph.ExpirationPatternTypeNoExpiration),
			},
		},
	})
	testRoleEligibilityScheduleRequestsClient_List(t, c)
	testRoleEligibilityScheduleRequestsClient_Get(t, c, *roleEligibilityScheduleRequest.ID)

	// Processing of the request action takes a little time internally, so it is necessary to sleep for another action
	time.Sleep(20 * time.Second)

	// Remove the role assignment from UI
	testRoleEligibilityScheduleRequestsClient_Create(t, c, msgraph.UnifiedRoleEligibilityScheduleRequest{
		Action:           utils.StringPtr(msgraph.UnifiedRoleScheduleRequestActionsAdminRemove),
		DirectoryScopeId: utils.StringPtr("/"),
		PrincipalId:      testUser.ID(),
		// ffd52fa5-98dc-465c-991d-fc073eb59f8f is the template id of the "Attribute Assignment Reader" built-in role
		RoleDefinitionId: utils.StringPtr("ffd52fa5-98dc-465c-991d-fc073eb59f8f"),
	})

	testUser_Delete(t, c, testUser)
}

func testRoleEligibilityScheduleRequestsClient_List(t *testing.T, c *test.Test) (roleEligibilityScheduleRequests *[]msgraph.UnifiedRoleEligibilityScheduleRequest) {
	roleEligibilityScheduleRequests, _, err := c.RoleEligibilityScheduleRequestsClient.List(c.Context, odata.Query{Top: 10})
	if err != nil {
		t.Fatalf("RoleEligibilityScheduleRequestsClient.List(): %v", err)
	}
	if roleEligibilityScheduleRequests == nil {
		t.Fatal("RoleEligibilityScheduleRequestsClient.List(): roleEligibilityScheduleRequests was nil")
	}
	return
}

func testRoleEligibilityScheduleRequestsClient_Get(t *testing.T, c *test.Test, id string) (roleEligibilityScheduleRequest *msgraph.UnifiedRoleEligibilityScheduleRequest) {
	roleEligibilityScheduleRequest, status, err := c.RoleEligibilityScheduleRequestsClient.Get(c.Context, id, odata.Query{})
	if err != nil {
		t.Fatalf("RoleEligibilityScheduleRequestsClient.Get(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("RoleEligibilityScheduleRequestsClient.Get(): invalid status: %d", status)
	}
	if roleEligibilityScheduleRequest == nil {
		t.Fatal("RoleEligibilityScheduleRequestsClient.Get(): roleEligibilityScheduleRequest was nil")
	}
	if *roleEligibilityScheduleRequest.ID != id {
		t.Fatal("RoleEligibilityScheduleRequestsClient.Get(): roleEligibilityScheduleRequest.ID was different")
	}
	return
}

func testRoleEligibilityScheduleRequestsClient_Create(t *testing.T, c *test.Test, r msgraph.UnifiedRoleEligibilityScheduleRequest) (roleEligibilityScheduleRequest *msgraph.UnifiedRoleEligibilityScheduleRequest) {
	roleEligibilityScheduleRequest, status, err := c.RoleEligibilityScheduleRequestsClient.Create(c.Context, r)
	if err != nil {
		t.Fatalf("RoleEligibilityScheduleRequestsClient.Create(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("RoleEligibilityScheduleRequestsClient.Create(): invalid status: %d", status)
	}
	if roleEligibilityScheduleRequest == nil {
		t.Fatal("RoleEligibilityScheduleRequestsClient.Create(): roleEligibilityScheduleRequest was nil")
	}
	if roleEligibilityScheduleRequest.ID == nil {
		t.Fatal("RoleEligibilityScheduleRequestsClient.Create(): roleEligibilityScheduleRequest.ID was nil")
	}
	return
}
