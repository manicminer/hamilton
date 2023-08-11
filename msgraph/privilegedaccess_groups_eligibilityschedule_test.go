package msgraph_test

import (
	"fmt"
	"os"
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

	var pimGroupId string
	if v := os.Getenv("PIM_GROUP_ID"); v != "" {
		pimGroupId = v
	} else {
		pimGroupId = ""
	}

	now := time.Now()
	future := now.AddDate(0, 0, 7)
	end := now.AddDate(0, 2, 0)

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

	testPrivilegedAccessGroupEligibilityScheduleClient_RequestsList(t, c)

	reqOwner := testPrivilegedAccessGroupEligibilityScheduleClient_RequestsCreate(t, c, msgraph.PrivilegedAccessGroupEligibilityScheduleRequest{
		AccessId:      utils.StringPtr(msgraph.PrivilegedAccessGroupRelationshipOwner),
		Action:        utils.StringPtr(msgraph.PrivilegedAccessGroupActionAdminAssign),
		GroupId:       utils.StringPtr(pimGroupId),
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

	reqMemberUser := testPrivilegedAccessGroupEligibilityScheduleClient_RequestsCreate(t, c, msgraph.PrivilegedAccessGroupEligibilityScheduleRequest{
		AccessId:      utils.StringPtr(msgraph.PrivilegedAccessGroupRelationshipMember),
		Action:        utils.StringPtr(msgraph.PrivilegedAccessGroupActionAdminAssign),
		GroupId:       utils.StringPtr(pimGroupId),
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

	reqMemberGroup := testPrivilegedAccessGroupEligibilityScheduleClient_RequestsCreate(t, c, msgraph.PrivilegedAccessGroupEligibilityScheduleRequest{
		AccessId:      utils.StringPtr(msgraph.PrivilegedAccessGroupRelationshipMember),
		Action:        utils.StringPtr(msgraph.PrivilegedAccessGroupActionAdminAssign),
		GroupId:       utils.StringPtr(pimGroupId),
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

	testPrivilegedAccessGroupEligibilityScheduleClient_RequestsGet(t, c, *reqOwner.ID)
	testPrivilegedAccessGroupEligibilityScheduleClient_RequestsGet(t, c, *reqMemberUser.ID)
	testPrivilegedAccessGroupEligibilityScheduleClient_RequestsGet(t, c, *reqMemberGroup.ID)

	schedules := testPrivilegedAccessGroupEligibilityScheduleClient_List(t, c, odata.Query{
		Filter: fmt.Sprintf("groupId eq '%s'", pimGroupId),
	})
	for _, sch := range *schedules {
		testPrivilegedAccessGroupEligibilityScheduleClient_Get(t, c, *sch.ID)
	}

	instances := testPrivilegedAccessGroupEligibilityScheduleClient_InstancesList(t, c, odata.Query{
		Filter: fmt.Sprintf("groupId eq '%s'", pimGroupId),
	})
	for _, inst := range *instances {
		testPrivilegedAccessGroupEligibilityScheduleClient_InstancesGet(t, c, *inst.ID)
	}

	testPrivilegedAccessGroupEligibilityScheduleClient_RequestsCancel(t, c, *reqMemberUser.ID)
	testPrivilegedAccessGroupEligibilityScheduleClient_RequestsCancel(t, c, *reqMemberGroup.ID)
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

func testPrivilegedAccessGroupEligibilityScheduleClient_InstancesList(t *testing.T, c *test.Test, query odata.Query) (instances *[]msgraph.PrivilegedAccessGroupEligibilityScheduleInstance) {
	instances, status, err := c.PrivilegedAccessGroupEligibilityScheduleClient.InstancesList(c.Context, query)
	if err != nil {
		t.Fatalf("PrivilegedAccessGroupEligibilityScheduleClient.InstancesList(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("PrivilegedAccessGroupEligibilityScheduleClient.InstancesList(): invalid status: %d", status)
	}
	if instances == nil {
		t.Fatal("PrivilegedAccessGroupEligibilityScheduleClient.InstancesList(): PrivilegedAccessGroupEligibilitySchedule was nil")
	}
	if len(*instances) == 0 {
		t.Fatal("PrivilegedAccessGroupEligibilityScheduleClient.List(): Returned zero results")
	}
	return
}

func testPrivilegedAccessGroupEligibilityScheduleClient_InstancesGet(t *testing.T, c *test.Test, id string) (request *msgraph.PrivilegedAccessGroupEligibilityScheduleInstance) {
	request, status, err := c.PrivilegedAccessGroupEligibilityScheduleClient.InstancesGet(c.Context, id)
	if err != nil {
		t.Fatalf("PrivilegedAccessGroupEligibilityScheduleClient.InstancesGet(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("PrivilegedAccessGroupEligibilityScheduleClient.InstancesGet(): invalid status: %d", status)
	}
	if request == nil {
		t.Fatal("PrivilegedAccessGroupEligibilityScheduleClient.InstancesGet(): PrivilegedAccessGroupEligibilitySchedule was nil")
	}
	return
}

func testPrivilegedAccessGroupEligibilityScheduleClient_RequestsCreate(t *testing.T, c *test.Test, r msgraph.PrivilegedAccessGroupEligibilityScheduleRequest) (request *msgraph.PrivilegedAccessGroupEligibilityScheduleRequest) {
	request, status, err := c.PrivilegedAccessGroupEligibilityScheduleClient.RequestsCreate(c.Context, r)
	if err != nil {
		t.Fatalf("PrivilegedAccessGroupEligibilityScheduleClient.Create(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("PrivilegedAccessGroupEligibilityScheduleClient.Create(): invalid status: %d", status)
	}
	if request == nil {
		t.Fatal("PrivilegedAccessGroupEligibilityScheduleClient.Create(): PrivilegedAccessGroupEligibilityScheduleRequest was nil")
	}
	if request.ID == nil {
		t.Fatal("PrivilegedAccessGroupEligibilityScheduleClient.Create(): PrivilegedAccessGroupEligibilityScheduleRequest.ID was nil")
	}
	return
}

func testPrivilegedAccessGroupEligibilityScheduleClient_RequestsList(t *testing.T, c *test.Test) (requests *[]msgraph.PrivilegedAccessGroupEligibilityScheduleRequest) {
	requests, status, err := c.PrivilegedAccessGroupEligibilityScheduleClient.RequestsList(c.Context, odata.Query{})
	if err != nil {
		t.Fatalf("PrivilegedAccessGroupEligibilityScheduleClient.List(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("PrivilegedAccessGroupEligibilityScheduleClient.List(): invalid status: %d", status)
	}
	if requests == nil {
		t.Fatal("PrivilegedAccessGroupEligibilityScheduleClient.List(): PrivilegedAccessGroupEligibilityScheduleRequest was nil")
	}
	return
}

func testPrivilegedAccessGroupEligibilityScheduleClient_RequestsGet(t *testing.T, c *test.Test, id string) (request *msgraph.PrivilegedAccessGroupEligibilityScheduleRequest) {
	request, status, err := c.PrivilegedAccessGroupEligibilityScheduleClient.RequestsGet(c.Context, id)
	if err != nil {
		t.Fatalf("PrivilegedAccessGroupEligibilityScheduleClient.Get(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("PrivilegedAccessGroupEligibilityScheduleClient.Get(): invalid status: %d", status)
	}
	if request == nil {
		t.Fatal("PrivilegedAccessGroupEligibilityScheduleClient.Get(): PrivilegedAccessGroupEligibilityScheduleRequest was nil")
	}
	return
}

func testPrivilegedAccessGroupEligibilityScheduleClient_RequestsCancel(t *testing.T, c *test.Test, id string) {
	status, err := c.PrivilegedAccessGroupEligibilityScheduleClient.RequestsCancel(c.Context, id)
	if err != nil {
		t.Fatalf("PrivilegedAccessGroupEligibilityScheduleClient.Cancel(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("PrivilegedAccessGroupEligibilityScheduleClient.Cancel(): invalid status: %d", status)
	}
}
