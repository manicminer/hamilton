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

func TestPrivilegedAccessGroupAssignmentScheduleRequestClient(t *testing.T) {
	c := test.NewTest(t)
	defer c.CancelFunc()

	self := testDirectoryObjectsClient_Get(t, c, c.Claims.ObjectId)
	now := time.Now()
	end := now.AddDate(0, 3, 0)

	pimGroup := testGroupsClient_Create(t, c, msgraph.Group{
		DisplayName:     utils.StringPtr("test-group"),
		MailEnabled:     utils.BoolPtr(false),
		MailNickname:    utils.StringPtr(fmt.Sprintf("test-group-%s", c.RandomString)),
		SecurityEnabled: utils.BoolPtr(true),
	})

	reqOwner := testPrivilegedAccessGroupAssignmentScheduleRequestClient_Create(t, c, msgraph.PrivilegedAccessGroupAssignmentScheduleRequest{
		AccessId:      utils.StringPtr(msgraph.PrivilegedAccessGroupRelationshipOwner),
		Action:        utils.StringPtr(msgraph.PrivilegedAccessGroupAssignmentActionAdminAssign),
		GroupId:       pimGroup.ID(),
		PrincipalId:   self.ID(),
		Justification: utils.StringPtr("Hamilton Testing"),
		ScheduleInfo: &msgraph.RequestSchedule{
			StartDateTime: &now,
			Expiration: &msgraph.ExpirationPattern{
				EndDateTime: &end,
				Type:        utils.StringPtr(msgraph.ExpirationPatternTypeAfterDateTime),
			},
		},
	})

	reqMember := testPrivilegedAccessGroupAssignmentScheduleRequestClient_Create(t, c, msgraph.PrivilegedAccessGroupAssignmentScheduleRequest{
		AccessId:      utils.StringPtr(msgraph.PrivilegedAccessGroupRelationshipMember),
		Action:        utils.StringPtr(msgraph.PrivilegedAccessGroupAssignmentActionAdminAssign),
		GroupId:       pimGroup.ID(),
		PrincipalId:   self.ID(),
		Justification: utils.StringPtr("Hamilton Testing"),
		ScheduleInfo: &msgraph.RequestSchedule{
			StartDateTime: &now,
			Expiration: &msgraph.ExpirationPattern{
				EndDateTime: &end,
				Type:        utils.StringPtr(msgraph.ExpirationPatternTypeAfterDateTime),
			},
		},
	})

	testPrivilegedAccessGroupAssignmentScheduleRequestClient_List(t, c)
	testPrivilegedAccessGroupAssignmentScheduleRequestClient_Get(t, c, *reqOwner.ID)
	testPrivilegedAccessGroupAssignmentScheduleRequestClient_Get(t, c, *reqMember.ID)
	testPrivilegedAccessGroupAssignmentScheduleRequestClient_Cancel(t, c, *reqOwner.ID)
	testPrivilegedAccessGroupAssignmentScheduleRequestClient_Cancel(t, c, *reqMember.ID)

	testGroupsClient_Delete(t, c, *pimGroup.ID())
}

func testPrivilegedAccessGroupAssignmentScheduleRequestClient_Create(t *testing.T, c *test.Test, r msgraph.PrivilegedAccessGroupAssignmentScheduleRequest) (request *msgraph.PrivilegedAccessGroupAssignmentScheduleRequest) {
	request, status, err := c.PrivilegedAccessGroupAssignmentScheduleRequestClient.Create(c.Context, r)
	if err != nil {
		t.Fatalf("PrivilegedAccessGroupAssignmentScheduleRequestClient.Create(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("PrivilegedAccessGroupAssignmentScheduleRequestClient.Create(): invalid status: %d", status)
	}
	if request == nil {
		t.Fatal("PrivilegedAccessGroupAssignmentScheduleRequestClient.Create(): PrivilegedAccessGroupAssignmentScheduleRequest was nil")
	}
	if request.ID == nil {
		t.Fatal("PrivilegedAccessGroupAssignmentScheduleRequestClient.Create(): PrivilegedAccessGroupAssignmentScheduleRequest.ID was nil")
	}
	return
}

func testPrivilegedAccessGroupAssignmentScheduleRequestClient_List(t *testing.T, c *test.Test) (requests *[]msgraph.PrivilegedAccessGroupAssignmentScheduleRequest) {
	requests, status, err := c.PrivilegedAccessGroupAssignmentScheduleRequestClient.List(c.Context, odata.Query{})
	if err != nil {
		t.Fatalf("PrivilegedAccessGroupAssignmentScheduleRequestClient.List(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("PrivilegedAccessGroupAssignmentScheduleRequestClient.List(): invalid status: %d", status)
	}
	if requests == nil {
		t.Fatal("PrivilegedAccessGroupAssignmentScheduleRequestClient.List(): PrivilegedAccessGroupAssignmentScheduleRequest was nil")
	}
	return
}

func testPrivilegedAccessGroupAssignmentScheduleRequestClient_Get(t *testing.T, c *test.Test, id string) (request *msgraph.PrivilegedAccessGroupAssignmentScheduleRequest) {
	request, status, err := c.PrivilegedAccessGroupAssignmentScheduleRequestClient.Get(c.Context, id)
	if err != nil {
		t.Fatalf("PrivilegedAccessGroupAssignmentScheduleRequestClient.Get(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("PrivilegedAccessGroupAssignmentScheduleRequestClient.Get(): invalid status: %d", status)
	}
	if request == nil {
		t.Fatal("PrivilegedAccessGroupAssignmentScheduleRequestClient.Get(): PrivilegedAccessGroupAssignmentScheduleRequest was nil")
	}
	return
}

func testPrivilegedAccessGroupAssignmentScheduleRequestClient_Cancel(t *testing.T, c *test.Test, id string) {
	status, err := c.PrivilegedAccessGroupAssignmentScheduleRequestClient.Cancel(c.Context, id)
	if err != nil {
		t.Fatalf("PrivilegedAccessGroupAssignmentScheduleRequestClient.Cancel(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("PrivilegedAccessGroupAssignmentScheduleRequestClient.Cancel(): invalid status: %d", status)
	}
}
