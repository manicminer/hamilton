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

func TestRoleEligibilityScheduleRequestClient(t *testing.T) {
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

	roleDefinition := testRoleDefinitionsClient_Create(t, c, msgraph.UnifiedRoleDefinition{
		Description: msgraph.NullableString("testing role eligibility schedule request"),
		DisplayName: utils.StringPtr("test-eligible"),
		IsEnabled:   utils.BoolPtr(true),
		RolePermissions: &[]msgraph.UnifiedRolePermission{
			{
				AllowedResourceActions: &[]string{
					"microsoft.directory/groups/allProperties/read",
				},
			},
		},
		Version: utils.StringPtr("1.5"),
	})

	now := time.Now()

	roleEligibilityScheduleRequest := testRoleEligibilityScheduleRequestClient_Create(t, c, msgraph.UnifiedRoleEligibilityScheduleRequest{
		Action:           utils.StringPtr(msgraph.UnifiedRoleScheduleRequestActionAdminAssign),
		RoleDefinitionId: roleDefinition.ID(),
		PrincipalId:      user.ID(),
		DirectoryScopeId: utils.StringPtr("/"),
		Justification:    utils.StringPtr("Test eligible"),
		ScheduleInfo: &msgraph.RequestSchedule{
			StartDateTime: &now,
			Expiration: &msgraph.ExpirationPattern{
				Type: utils.StringPtr(msgraph.ExpirationPatternTypeNoExpiration),
			},
		},
	})

	testRoleEligibilityScheduleRequestClient_Get(t, c, *roleEligibilityScheduleRequest.ID)
	testListReturnsID(t, c, *testRoleEligibilityScheduleRequestClient_List(t, c), *roleEligibilityScheduleRequest.ID)
	roleEligibilityScheduleRequest.Action = utils.StringPtr(msgraph.UnifiedRoleScheduleRequestActionAdminRemove)
	testRoleEligibilityScheduleRequestClient_Create(t, c, *roleEligibilityScheduleRequest)
	testUsersClient_Delete(t, c, *user.ID())
	testUsersClient_DeletePermanently(t, c, *user.ID())
	testRoleDefinitionsClient_Delete(t, c, *roleDefinition.ID())
}

func testRoleEligibilityScheduleRequestClient_Create(t *testing.T, c *test.Test, r msgraph.UnifiedRoleEligibilityScheduleRequest) (roleEligibilityScheduleRequest *msgraph.UnifiedRoleEligibilityScheduleRequest) {
	roleEligibilityScheduleRequest, status, err := c.RoleEligibilityScheduleRequestClient.Create(c.Context, r)
	if err != nil {
		t.Fatalf("RoleEligibilityScheduleRequestClient.Create(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("RoleEligibilityScheduleRequestClient.Create(): invalid status: %d", status)
	}
	if roleEligibilityScheduleRequest == nil {
		t.Fatal("RoleEligibilityScheduleRequestClient.Create(): roleEligibilityScheduleRequest was nil")
	}
	if roleEligibilityScheduleRequest.ID == nil {
		t.Fatal("RoleEligibilityScheduleRequestClient.Create(): roleEligibilityScheduleRequest.ID was nil")
	}
	return
}

func testRoleEligibilityScheduleRequestClient_Get(t *testing.T, c *test.Test, id string) (roleEligibilityScheduleRequest *msgraph.UnifiedRoleEligibilityScheduleRequest) {
	roleEligibilityScheduleRequest, status, err := c.RoleEligibilityScheduleRequestClient.Get(c.Context, id, odata.Query{})
	if err != nil {
		t.Fatalf("RoleEligibilityScheduleRequestClient.Get(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("RoleEligibilityScheduleRequestClient.Get(): invalid status: %d", status)
	}
	if roleEligibilityScheduleRequest == nil {
		t.Fatal("RoleEligibilityScheduleRequestClient.Get(): roleEligibilityScheduleRequest was nil")
	}
	return
}

func testRoleEligibilityScheduleRequestClient_List(t *testing.T, c *test.Test) (roleEligibilityScheduleRequests *[]msgraph.UnifiedRoleEligibilityScheduleRequest) {
	roleEligibilityScheduleRequests, status, err := c.RoleEligibilityScheduleRequestClient.List(c.Context)
	if err != nil {
		t.Fatalf("RoleEligibilityScheduleRequestClient.List(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("RoleEligibilityScheduleRequestClient.List(): invalid status: %d", status)
	}
	if roleEligibilityScheduleRequests == nil {
		t.Fatal("RoleEligibilityScheduleRequestClient.List(): roleEligibilityScheduleRequests was nil")
	}
	return
}

func testListReturnsID(t *testing.T, c *test.Test, roleEligibilityScheduleRequests []msgraph.UnifiedRoleEligibilityScheduleRequest, id string) {
	for _, r := range roleEligibilityScheduleRequests {
		if *r.ID == id {
			return
		}
	}
	t.Fatalf("RoleEligibilityScheduleRequestClient.List(): didn't return roleEligibilityScheduleRequest with id %s", id)
}
