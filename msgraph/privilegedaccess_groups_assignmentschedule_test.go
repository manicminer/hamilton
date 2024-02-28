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

func TestPrivilegedAccessGroupAssignmentScheduleClient(t *testing.T) {
	c := test.NewTest(t)
	defer c.CancelFunc()

	now := time.Now()
	future := now.AddDate(0, 0, 7)
	end := now.AddDate(0, 2, 0)

	pimGroup := testGroupsClient_Create(t, c, msgraph.Group{
		DisplayName:     utils.StringPtr("test-pim-group"),
		MailEnabled:     utils.BoolPtr(false),
		MailNickname:    utils.StringPtr(fmt.Sprintf("test-pim-group-%s", c.RandomString)),
		SecurityEnabled: utils.BoolPtr(true),
	})
	defer testGroupsClient_Delete(t, c, *pimGroup.ID())

	groupMember := testGroupsClient_Create(t, c, msgraph.Group{
		DisplayName:     utils.StringPtr("test-group-member"),
		MailEnabled:     utils.BoolPtr(false),
		MailNickname:    utils.StringPtr(fmt.Sprintf("test-group-member-%s", c.RandomString)),
		SecurityEnabled: utils.BoolPtr(true),
	})
	defer testGroupsClient_Delete(t, c, *groupMember.ID())

	userMember := testUsersClient_Create(t, c, msgraph.User{
		AccountEnabled:    utils.BoolPtr(true),
		DisplayName:       utils.StringPtr("test-user-groupmember"),
		MailNickname:      utils.StringPtr(fmt.Sprintf("test-user-groupmember-%s", c.RandomString)),
		UserPrincipalName: utils.StringPtr(fmt.Sprintf("test-user-groupmember-%s@%s", c.RandomString, c.Connections["default"].DomainName)),
		PasswordProfile: &msgraph.UserPasswordProfile{
			Password: utils.StringPtr(fmt.Sprintf("IrPa55w0rd%s", c.RandomString)),
		},
	})
	defer testUsersClient_Delete(t, c, *userMember.ID())

	userOwner := testUsersClient_Create(t, c, msgraph.User{
		AccountEnabled:    utils.BoolPtr(true),
		DisplayName:       utils.StringPtr("test-user-groupowner"),
		MailNickname:      utils.StringPtr(fmt.Sprintf("test-user-groupowner-%s", c.RandomString)),
		UserPrincipalName: utils.StringPtr(fmt.Sprintf("test-user-groupowner-%s@%s", c.RandomString, c.Connections["default"].DomainName)),
		PasswordProfile: &msgraph.UserPasswordProfile{
			Password: utils.StringPtr(fmt.Sprintf("IrPa55w0rd%s", c.RandomString)),
		},
	})
	defer testUsersClient_Delete(t, c, *userOwner.ID())

	testPrivilegedAccessGroupAssignmentScheduleRequestsClient_List(t, c)

	reqOwner := testPrivilegedAccessGroupAssignmentScheduleRequestsClient_Create(t, c, msgraph.PrivilegedAccessGroupAssignmentScheduleRequest{
		AccessId:      msgraph.PrivilegedAccessGroupRelationshipOwner,
		Action:        msgraph.PrivilegedAccessGroupActionAdminAssign,
		GroupId:       pimGroup.ID(),
		PrincipalId:   userOwner.ID(),
		Justification: utils.StringPtr("Hamilton Testing"),
		ScheduleInfo: &msgraph.RequestSchedule{
			StartDateTime: &now,
			Expiration: &msgraph.ExpirationPattern{
				EndDateTime: &end,
				Type:        msgraph.ExpirationPatternTypeAfterDateTime,
			},
		},
	})

	reqMemberUser := testPrivilegedAccessGroupAssignmentScheduleRequestsClient_Create(t, c, msgraph.PrivilegedAccessGroupAssignmentScheduleRequest{
		AccessId:      msgraph.PrivilegedAccessGroupRelationshipMember,
		Action:        msgraph.PrivilegedAccessGroupActionAdminAssign,
		GroupId:       pimGroup.ID(),
		PrincipalId:   userMember.ID(),
		Justification: utils.StringPtr("Hamilton Testing"),
		ScheduleInfo: &msgraph.RequestSchedule{
			StartDateTime: &future,
			Expiration: &msgraph.ExpirationPattern{
				EndDateTime: &end,
				Type:        msgraph.ExpirationPatternTypeAfterDateTime,
			},
		},
	})

	reqMemberGroup := testPrivilegedAccessGroupAssignmentScheduleRequestsClient_Create(t, c, msgraph.PrivilegedAccessGroupAssignmentScheduleRequest{
		AccessId:      msgraph.PrivilegedAccessGroupRelationshipMember,
		Action:        msgraph.PrivilegedAccessGroupActionAdminAssign,
		GroupId:       pimGroup.ID(),
		PrincipalId:   groupMember.ID(),
		Justification: utils.StringPtr("Hamilton Testing"),
		ScheduleInfo: &msgraph.RequestSchedule{
			StartDateTime: &future,
			Expiration: &msgraph.ExpirationPattern{
				EndDateTime: &end,
				Type:        msgraph.ExpirationPatternTypeAfterDateTime,
			},
		},
	})

	testPrivilegedAccessGroupAssignmentScheduleRequestsClient_Get(t, c, *reqOwner.ID)
	testPrivilegedAccessGroupAssignmentScheduleRequestsClient_Get(t, c, *reqMemberUser.ID)
	testPrivilegedAccessGroupAssignmentScheduleRequestsClient_Get(t, c, *reqMemberGroup.ID)

	schedules := testPrivilegedAccessGroupAssignmentScheduleClient_List(t, c, odata.Query{
		Filter: fmt.Sprintf("groupId eq '%s'", *pimGroup.ID()),
	})
	for _, sch := range *schedules {
		testPrivilegedAccessGroupAssignmentScheduleClient_Get(t, c, *sch.ID)
	}

	instances := testPrivilegedAccessGroupAssignmentScheduleInstancesClient_List(t, c, odata.Query{
		Filter: fmt.Sprintf("groupId eq '%s'", *pimGroup.ID()),
	})
	for _, inst := range *instances {
		testPrivilegedAccessGroupAssignmentScheduleInstancesClient_Get(t, c, *inst.ID)
	}

	testPrivilegedAccessGroupAssignmentScheduleRequestsClient_Cancel(t, c, *reqMemberUser.ID)
	testPrivilegedAccessGroupAssignmentScheduleRequestsClient_Cancel(t, c, *reqMemberGroup.ID)
}

func testPrivilegedAccessGroupAssignmentScheduleClient_List(t *testing.T, c *test.Test, query odata.Query) (schedules *[]msgraph.PrivilegedAccessGroupAssignmentSchedule) {
	schedules, status, err := c.PrivilegedAccessGroupAssignmentScheduleClient.List(c.Context, query)
	if err != nil {
		t.Fatalf("PrivilegedAccessGroupAssignmentScheduleClient.List(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("PrivilegedAccessGroupAssignmentScheduleClient.List(): invalid status: %d", status)
	}
	if schedules == nil {
		t.Fatal("PrivilegedAccessGroupAssignmentScheduleClient.List(): PrivilegedAccessGroupAssignmentSchedule was nil")
	}
	if len(*schedules) == 0 {
		t.Fatal("PrivilegedAccessGroupAssignmentScheduleClient.List(): Returned zero results")
	}
	return
}

func testPrivilegedAccessGroupAssignmentScheduleClient_Get(t *testing.T, c *test.Test, id string) (request *msgraph.PrivilegedAccessGroupAssignmentSchedule) {
	request, status, err := c.PrivilegedAccessGroupAssignmentScheduleClient.Get(c.Context, id)
	if err != nil {
		t.Fatalf("PrivilegedAccessGroupAssignmentScheduleClient.Get(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("PrivilegedAccessGroupAssignmentScheduleClient.Get(): invalid status: %d", status)
	}
	if request == nil {
		t.Fatal("PrivilegedAccessGroupAssignmentScheduleClient.Get(): PrivilegedAccessGroupAssignmentSchedule was nil")
	}
	return
}
