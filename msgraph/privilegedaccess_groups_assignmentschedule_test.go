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

func TestPrivilegedAccessGroupAssignmentScheduleClient(t *testing.T) {
	c := test.NewTest(t)
	defer c.CancelFunc()

	var pimGroupId string
	if v := os.Getenv("PIM_GROUP_ID"); v != "" {
		pimGroupId = v
	} else {
		pimGroupId = "8647a803-6803-46d6-bad4-bec15c5989d6"
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

	testPrivilegedAccessGroupAssignmentScheduleClient_RequestsList(t, c)

	reqOwner := testPrivilegedAccessGroupAssignmentScheduleClient_RequestsCreate(t, c, msgraph.PrivilegedAccessGroupAssignmentScheduleRequest{
		AccessId:      utils.StringPtr(msgraph.PrivilegedAccessGroupRelationshipOwner),
		Action:        utils.StringPtr(msgraph.PrivilegedAccessGroupAssignmentActionAdminAssign),
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

	reqMemberUser := testPrivilegedAccessGroupAssignmentScheduleClient_RequestsCreate(t, c, msgraph.PrivilegedAccessGroupAssignmentScheduleRequest{
		AccessId:      utils.StringPtr(msgraph.PrivilegedAccessGroupRelationshipMember),
		Action:        utils.StringPtr(msgraph.PrivilegedAccessGroupAssignmentActionAdminAssign),
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

	reqMemberGroup := testPrivilegedAccessGroupAssignmentScheduleClient_RequestsCreate(t, c, msgraph.PrivilegedAccessGroupAssignmentScheduleRequest{
		AccessId:      utils.StringPtr(msgraph.PrivilegedAccessGroupRelationshipMember),
		Action:        utils.StringPtr(msgraph.PrivilegedAccessGroupAssignmentActionAdminAssign),
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

	testPrivilegedAccessGroupAssignmentScheduleClient_RequestsGet(t, c, *reqOwner.ID)
	testPrivilegedAccessGroupAssignmentScheduleClient_RequestsGet(t, c, *reqMemberUser.ID)
	testPrivilegedAccessGroupAssignmentScheduleClient_RequestsGet(t, c, *reqMemberGroup.ID)

	testPrivilegedAccessGroupAssignmentScheduleClient_RequestsCancel(t, c, *reqMemberUser.ID)
	testPrivilegedAccessGroupAssignmentScheduleClient_RequestsCancel(t, c, *reqMemberGroup.ID)

}

func testPrivilegedAccessGroupAssignmentScheduleClient_List(t *testing.T, c *test.Test, query odata.Query) (requests *[]msgraph.PrivilegedAccessGroupAssignmentSchedule) {
	requests, status, err := c.PrivilegedAccessGroupAssignmentScheduleClient.List(c.Context, query)
	if err != nil {
		t.Fatalf("PrivilegedAccessGroupAssignmentScheduleClient.List(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("PrivilegedAccessGroupAssignmentScheduleClient.List(): invalid status: %d", status)
	}
	if requests == nil {
		t.Fatal("PrivilegedAccessGroupAssignmentScheduleClient.List(): PrivilegedAccessGroupAssignmentSchedule was nil")
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

func testPrivilegedAccessGroupAssignmentScheduleClient_InstancesList(t *testing.T, c *test.Test, query odata.Query) (requests *[]msgraph.PrivilegedAccessGroupAssignmentScheduleInstance) {
	requests, status, err := c.PrivilegedAccessGroupAssignmentScheduleClient.InstancesList(c.Context, query)
	if err != nil {
		t.Fatalf("PrivilegedAccessGroupAssignmentScheduleClient.InstancesList(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("PrivilegedAccessGroupAssignmentScheduleClient.InstancesList(): invalid status: %d", status)
	}
	if requests == nil {
		t.Fatal("PrivilegedAccessGroupAssignmentScheduleClient.InstancesList(): PrivilegedAccessGroupAssignmentSchedule was nil")
	}
	return
}

func testPrivilegedAccessGroupAssignmentScheduleClient_InstancesGet(t *testing.T, c *test.Test, id string) (request *msgraph.PrivilegedAccessGroupAssignmentScheduleInstance) {
	request, status, err := c.PrivilegedAccessGroupAssignmentScheduleClient.InstancesGet(c.Context, id)
	if err != nil {
		t.Fatalf("PrivilegedAccessGroupAssignmentScheduleClient.InstancesGet(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("PrivilegedAccessGroupAssignmentScheduleClient.InstancesGet(): invalid status: %d", status)
	}
	if request == nil {
		t.Fatal("PrivilegedAccessGroupAssignmentScheduleClient.InstancesGet(): PrivilegedAccessGroupAssignmentSchedule was nil")
	}
	return
}

func testPrivilegedAccessGroupAssignmentScheduleClient_RequestsCreate(t *testing.T, c *test.Test, r msgraph.PrivilegedAccessGroupAssignmentScheduleRequest) (request *msgraph.PrivilegedAccessGroupAssignmentScheduleRequest) {
	request, status, err := c.PrivilegedAccessGroupAssignmentScheduleClient.RequestsCreate(c.Context, r)
	if err != nil {
		t.Fatalf("PrivilegedAccessGroupAssignmentScheduleClient.Create(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("PrivilegedAccessGroupAssignmentScheduleClient.Create(): invalid status: %d", status)
	}
	if request == nil {
		t.Fatal("PrivilegedAccessGroupAssignmentScheduleClient.Create(): PrivilegedAccessGroupAssignmentScheduleRequest was nil")
	}
	if request.ID == nil {
		t.Fatal("PrivilegedAccessGroupAssignmentScheduleClient.Create(): PrivilegedAccessGroupAssignmentScheduleRequest.ID was nil")
	}
	return
}

func testPrivilegedAccessGroupAssignmentScheduleClient_RequestsList(t *testing.T, c *test.Test) (requests *[]msgraph.PrivilegedAccessGroupAssignmentScheduleRequest) {
	requests, status, err := c.PrivilegedAccessGroupAssignmentScheduleClient.RequestsList(c.Context, odata.Query{})
	if err != nil {
		t.Fatalf("PrivilegedAccessGroupAssignmentScheduleClient.List(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("PrivilegedAccessGroupAssignmentScheduleClient.List(): invalid status: %d", status)
	}
	if requests == nil {
		t.Fatal("PrivilegedAccessGroupAssignmentScheduleClient.List(): PrivilegedAccessGroupAssignmentScheduleRequest was nil")
	}
	return
}

func testPrivilegedAccessGroupAssignmentScheduleClient_RequestsGet(t *testing.T, c *test.Test, id string) (request *msgraph.PrivilegedAccessGroupAssignmentScheduleRequest) {
	request, status, err := c.PrivilegedAccessGroupAssignmentScheduleClient.RequestsGet(c.Context, id)
	if err != nil {
		t.Fatalf("PrivilegedAccessGroupAssignmentScheduleClient.Get(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("PrivilegedAccessGroupAssignmentScheduleClient.Get(): invalid status: %d", status)
	}
	if request == nil {
		t.Fatal("PrivilegedAccessGroupAssignmentScheduleClient.Get(): PrivilegedAccessGroupAssignmentScheduleRequest was nil")
	}
	return
}

func testPrivilegedAccessGroupAssignmentScheduleClient_RequestsCancel(t *testing.T, c *test.Test, id string) {
	status, err := c.PrivilegedAccessGroupAssignmentScheduleClient.RequestsCancel(c.Context, id)
	if err != nil {
		t.Fatalf("PrivilegedAccessGroupAssignmentScheduleClient.Cancel(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("PrivilegedAccessGroupAssignmentScheduleClient.Cancel(): invalid status: %d", status)
	}
}
