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

func TestPrivilegedAccessGroupEligibilityScheduleClient(t *testing.T) {
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

	testPrivilegedAccessGroupEligibilityScheduleRequestsClient_List(t, c)

	reqOwner := testPrivilegedAccessGroupEligibilityScheduleRequestsClient_Create(t, c, msgraph.PrivilegedAccessGroupEligibilityScheduleRequest{
		AccessId:      msgraph.PrivilegedAccessGroupRelationshipOwner,
		Action:        msgraph.PrivilegedAccessGroupActionAdminAssign,
		GroupId:       pimGroup.ID(),
		PrincipalId:   userOwner.ID(),
		Justification: utils.StringPtr("Hamilton Testing"),
		ScheduleInfo: &msgraph.RequestSchedule{
			StartDateTime: &now,
			Expiration: &msgraph.ExpirationPattern{
				EndDateTime: &end,
				Type:        utils.StringPtr(msgraph.ExpirationPatternTypeAfterDateTime),
			},
		},
	})

	reqMemberUser := testPrivilegedAccessGroupEligibilityScheduleRequestsClient_Create(t, c, msgraph.PrivilegedAccessGroupEligibilityScheduleRequest{
		AccessId:      msgraph.PrivilegedAccessGroupRelationshipMember,
		Action:        msgraph.PrivilegedAccessGroupActionAdminAssign,
		GroupId:       pimGroup.ID(),
		PrincipalId:   userMember.ID(),
		Justification: utils.StringPtr("Hamilton Testing"),
		ScheduleInfo: &msgraph.RequestSchedule{
			StartDateTime: &future,
			Expiration: &msgraph.ExpirationPattern{
				EndDateTime: &end,
				Type:        utils.StringPtr(msgraph.ExpirationPatternTypeAfterDateTime),
			},
		},
	})

	reqMemberGroup := testPrivilegedAccessGroupEligibilityScheduleRequestsClient_Create(t, c, msgraph.PrivilegedAccessGroupEligibilityScheduleRequest{
		AccessId:      msgraph.PrivilegedAccessGroupRelationshipMember,
		Action:        msgraph.PrivilegedAccessGroupActionAdminAssign,
		GroupId:       pimGroup.ID(),
		PrincipalId:   groupMember.ID(),
		Justification: utils.StringPtr("Hamilton Testing"),
		ScheduleInfo: &msgraph.RequestSchedule{
			StartDateTime: &future,
			Expiration: &msgraph.ExpirationPattern{
				EndDateTime: &end,
				Type:        utils.StringPtr(msgraph.ExpirationPatternTypeAfterDateTime),
			},
		},
	})

	testPrivilegedAccessGroupEligibilityScheduleRequestsClient_Get(t, c, *reqOwner.ID)
	testPrivilegedAccessGroupEligibilityScheduleRequestsClient_Get(t, c, *reqMemberUser.ID)
	testPrivilegedAccessGroupEligibilityScheduleRequestsClient_Get(t, c, *reqMemberGroup.ID)

	schedules := testPrivilegedAccessGroupEligibilityScheduleClient_List(t, c, odata.Query{
		Filter: fmt.Sprintf("groupId eq '%s'", *pimGroup.ID()),
	})
	for _, sch := range *schedules {
		testPrivilegedAccessGroupEligibilityScheduleClient_Get(t, c, *sch.ID)
	}

	instances := testPrivilegedAccessGroupEligibilityScheduleInstancesClient_List(t, c, odata.Query{
		Filter: fmt.Sprintf("groupId eq '%s'", *pimGroup.ID()),
	})
	for _, inst := range *instances {
		testPrivilegedAccessGroupEligibilityScheduleInstancesClient_Get(t, c, *inst.ID)
	}

	testPrivilegedAccessGroupEligibilityScheduleRequestsClient_Cancel(t, c, *reqMemberUser.ID)
	testPrivilegedAccessGroupEligibilityScheduleRequestsClient_Cancel(t, c, *reqMemberGroup.ID)
}

func testPrivilegedAccessGroupEligibilityScheduleClient_List(t *testing.T, c *test.Test, query odata.Query) (schedules *[]msgraph.PrivilegedAccessGroupEligibilitySchedule) {
	schedules, status, err := c.PrivilegedAccessGroupEligibilityScheduleClient.List(c.Context, query)
	if err != nil {
		t.Fatalf("PrivilegedAccessGroupEligibilityScheduleClient.List(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("PrivilegedAccessGroupEligibilityScheduleClient.List(): invalid status: %d", status)
	}
	if schedules == nil {
		t.Fatal("PrivilegedAccessGroupEligibilityScheduleClient.List(): PrivilegedAccessGroupEligibilitySchedule was nil")
	}
	if len(*schedules) == 0 {
		t.Fatal("PrivilegedAccessGroupEligibilityScheduleClient.List(): Returned zero results")
	}
	return
}

func testPrivilegedAccessGroupEligibilityScheduleClient_Get(t *testing.T, c *test.Test, id string) (request *msgraph.PrivilegedAccessGroupEligibilitySchedule) {
	request, status, err := c.PrivilegedAccessGroupEligibilityScheduleClient.Get(c.Context, id)
	if err != nil {
		t.Fatalf("PrivilegedAccessGroupEligibilityScheduleClient.Get(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("PrivilegedAccessGroupEligibilityScheduleClient.Get(): invalid status: %d", status)
	}
	if request == nil {
		t.Fatal("PrivilegedAccessGroupEligibilityScheduleClient.Get(): PrivilegedAccessGroupEligibilitySchedule was nil")
	}
	return
}
