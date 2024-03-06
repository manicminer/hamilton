package msgraph_test

import (
	"fmt"
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

	user := testUsersClient_Create(t, c, msgraph.User{
		AccountEnabled:    utils.BoolPtr(true),
		DisplayName:       utils.StringPtr("test-user"),
		MailNickname:      utils.StringPtr(fmt.Sprintf("test-user-%s", c.RandomString)),
		UserPrincipalName: utils.StringPtr(fmt.Sprintf("test-user-%s@%s", c.RandomString, c.Connections["default"].DomainName)),
		PasswordProfile: &msgraph.UserPasswordProfile{
			Password: utils.StringPtr(fmt.Sprintf("IrPa55w0rd%s", c.RandomString)),
		},
	})

	directoryRoles := testDirectoryRolesClient_List(t, c)
	directoryRole := (*directoryRoles)[0]

	now := time.Now()

	roleEligibilityScheduleRequest := testRoleEligibilityScheduleRequestsClient_Create(t, c, msgraph.UnifiedRoleEligibilityScheduleRequest{
		Action:           utils.StringPtr(msgraph.UnifiedRoleScheduleRequestActionAdminAssign),
		RoleDefinitionId: directoryRole.RoleTemplateId,
		PrincipalId:      user.ID(),
		DirectoryScopeId: utils.StringPtr("/"),
		Justification:    utils.StringPtr("Test eligible"),
		ScheduleInfo: &msgraph.RequestSchedule{
			StartDateTime: &now,
			Expiration: &msgraph.ExpirationPattern{
				Type: msgraph.ExpirationPatternTypeNoExpiration,
			},
		},
	})

	testRoleEligibilityScheduleRequestsClient_Get(t, c, *roleEligibilityScheduleRequest.ID)
	testListReturnsID(t, c, *testRoleEligibilityScheduleRequestsClient_List(t, c), *roleEligibilityScheduleRequest.ID)
	roleEligibilityScheduleRequest.Action = utils.StringPtr(msgraph.UnifiedRoleScheduleRequestActionAdminRemove)
	testRoleEligibilityScheduleRequestsClient_Create(t, c, *roleEligibilityScheduleRequest)
	testUsersClient_Delete(t, c, *user.ID())
	testUsersClient_DeletePermanently(t, c, *user.ID())
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
	return
}

func testRoleEligibilityScheduleRequestsClient_List(t *testing.T, c *test.Test) (roleEligibilityScheduleRequests *[]msgraph.UnifiedRoleEligibilityScheduleRequest) {
	roleEligibilityScheduleRequests, status, err := c.RoleEligibilityScheduleRequestsClient.List(c.Context)
	if err != nil {
		t.Fatalf("RoleEligibilityScheduleRequestsClient.List(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("RoleEligibilityScheduleRequestsClient.List(): invalid status: %d", status)
	}
	if roleEligibilityScheduleRequests == nil {
		t.Fatal("RoleEligibilityScheduleRequestsClient.List(): roleEligibilityScheduleRequests was nil")
	}
	return
}

func testListReturnsID(t *testing.T, c *test.Test, roleEligibilityScheduleRequests []msgraph.UnifiedRoleEligibilityScheduleRequest, id string) {
	for _, r := range roleEligibilityScheduleRequests {
		if *r.ID == id {
			return
		}
	}
	t.Fatalf("RoleEligibilityScheduleRequestsClient.List(): didn't return roleEligibilityScheduleRequest with id %s", id)
}
