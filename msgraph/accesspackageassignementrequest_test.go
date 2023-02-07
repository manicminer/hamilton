package msgraph_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/msgraph"
	"github.com/manicminer/hamilton/odata"
)

func TestAccessPackageAssignmentRequestClient(t *testing.T) {
	c := test.NewTest(t)
	defer c.CancelFunc()

	// Create test Catalog
	accessPackageCatalog := testAccessPackageCatalog_Create(t, c)

	// Create AP
	accessPackage := testAccessPackage_Create(t, c, accessPackageCatalog)

	currentTime := time.Now()

	approverUser := testUsersClient_Create(t, c, msgraph.User{
		AccountEnabled:    utils.BoolPtr(true),
		DisplayName:       utils.StringPtr("test-approver"),
		MailNickname:      utils.StringPtr(fmt.Sprintf("test-approver-%s", c.RandomString)),
		UserPrincipalName: utils.StringPtr(fmt.Sprintf("test-approver-%s@%s", c.RandomString, c.Connections["default"].DomainName)),
		PasswordProfile: &msgraph.UserPasswordProfile{
			Password: utils.StringPtr(fmt.Sprintf("IrPa55w0rd%s", c.RandomString)),
		},
	})

	user := testUsersClient_Create(t, c, msgraph.User{
		AccountEnabled:    utils.BoolPtr(true),
		DisplayName:       utils.StringPtr("test-user"),
		MailNickname:      utils.StringPtr(fmt.Sprintf("test-user-%s", c.RandomString)),
		UserPrincipalName: utils.StringPtr(fmt.Sprintf("test-user-%s@%s", c.RandomString, c.Connections["default"].DomainName)),
		PasswordProfile: &msgraph.UserPasswordProfile{
			Password: utils.StringPtr(fmt.Sprintf("IrPa55w0rd%s", c.RandomString)),
		},
	})

	user2 := testUsersClient_Create(t, c, msgraph.User{
		AccountEnabled:    utils.BoolPtr(true),
		DisplayName:       utils.StringPtr("test-user2"),
		MailNickname:      utils.StringPtr(fmt.Sprintf("test-user2-%s", c.RandomString)),
		UserPrincipalName: utils.StringPtr(fmt.Sprintf("test-user2-%s@%s", c.RandomString, c.Connections["default"].DomainName)),
		PasswordProfile: &msgraph.UserPasswordProfile{
			Password: utils.StringPtr(fmt.Sprintf("IrPa55w0rd%s", c.RandomString)),
		},
	})

	// Adding a sleep, original testing showed the policy could not find the user applied as an approver. Adding a sleep solved this issue.
	// Originally attempted a get of user after create but didn't help. Seems to be an eventual consistency issue
	time.Sleep(time.Second * 20)

	// Create Assignment Policy
	accessPackageAssignmentPolicy := testAccessPackageAssignmentPolicyClient_Create(t, c, msgraph.AccessPackageAssignmentPolicy{
		AccessPackageId: accessPackage.ID,
		AccessReviewSettings: &msgraph.AssignmentReviewSettings{
			AccessReviewTimeoutBehavior:     msgraph.AccessReviewTimeoutBehaviorTypeRemoveAccess,
			IsEnabled:                       utils.BoolPtr(true),
			StartDateTime:                   &currentTime,
			DurationInDays:                  utils.Int32Ptr(5),
			RecurrenceType:                  msgraph.AccessReviewRecurranceTypeMonthly,
			ReviewerType:                    msgraph.AccessReviewReviewerTypeSelf,
			IsAccessRecommendationEnabled:   utils.BoolPtr(true),
			IsApprovalJustificationRequired: utils.BoolPtr(true),
		},
		DisplayName: utils.StringPtr(fmt.Sprintf("Test-AP-Policy-Assignment-%s", c.RandomString)),
		Description: utils.StringPtr("Test AP Policy Assignment Description"),
		RequestorSettings: &msgraph.RequestorSettings{
			ScopeType:      msgraph.RequestorSettingsScopeTypeNoSubjects,
			AcceptRequests: utils.BoolPtr(true),
		},
		RequestApprovalSettings: &msgraph.ApprovalSettings{
			IsApprovalRequired:               utils.BoolPtr(true),
			IsApprovalRequiredForExtension:   utils.BoolPtr(false),
			IsRequestorJustificationRequired: utils.BoolPtr(false),
			ApprovalMode:                     msgraph.ApprovalModeSingleStage,
			ApprovalStages: &[]msgraph.ApprovalStage{
				{
					ApprovalStageTimeOutInDays:      utils.Int32Ptr(7),
					IsApproverJustificationRequired: utils.BoolPtr(false),
					IsEscalationEnabled:             utils.BoolPtr(false),
					PrimaryApprovers: &[]msgraph.UserSet{
						{
							ODataType:   utils.StringPtr(odata.TypeSingleUser),
							Description: utils.StringPtr("approver"),
							IsBackup:    utils.BoolPtr(false),
							ID:          approverUser.ID(),
						},
					},
				},
			},
		},
	})

	ap := testAccessPackageAssignmentRequestClient_Create(t, c, msgraph.AccessPackageAssignmentRequest{
		RequestType: utils.StringPtr(msgraph.AccessPackageRequestTypeAdminAdd),
		Assignment: &msgraph.AccessPackageAssignment{
			TargetID:            user.ID(),
			AssignementPolicyID: accessPackageAssignmentPolicy.ID,
			AccessPackageID:     accessPackage.ID,
		},
	})

	ap2 := testAccessPackageAssignmentRequestClient_Create(t, c, msgraph.AccessPackageAssignmentRequest{
		RequestType: utils.StringPtr(msgraph.AccessPackageRequestTypeAdminAdd),
		Assignment: &msgraph.AccessPackageAssignment{
			TargetID:            user2.ID(),
			AssignementPolicyID: accessPackageAssignmentPolicy.ID,
			AccessPackageID:     accessPackage.ID,
		},
	})

	apIDs := []string{*ap.ID, *ap2.ID}
	testAccessPackageAssignmentRequestClient_List(t, c, apIDs)

	testAccessPackageAssignmentRequestClient_Cancel(t, c, *ap.ID)
	testAccessPackageAssignmentRequestClient_Cancel(t, c, *ap2.ID)

	testAccessPackageAssignmentRequestClient_Get(t, c, *ap.ID)

	//Cleanup
	deleteWhenPossible(t, c, ap)
	deleteWhenPossible(t, c, ap2)
	testAccessPackageAssignmentPolicyClient_Delete(t, c, *accessPackageAssignmentPolicy.ID)
	testAccessPackage_Delete(t, c, *accessPackage.ID)
	testAccessPackageCatalog_Delete(t, c, accessPackageCatalog)
	testUser_Delete(t, c, user)
	testUser_Delete(t, c, user2)
	testUser_Delete(t, c, approverUser)
}

func deleteWhenPossible(t *testing.T, c *test.Test, ap *msgraph.AccessPackageAssignmentRequest) {
	// Can only delete a request if it is in specific states
	switch ap.State {
	case utils.StringPtr(msgraph.AccessPackageRequestStateDenied):
		testAccessPacakgeAssignmentRequestClient_Delete(t, c, *ap.ID)
	case utils.StringPtr(msgraph.AccessPackageRequestStateCanceled):
		testAccessPacakgeAssignmentRequestClient_Delete(t, c, *ap.ID)
	case utils.StringPtr(msgraph.AccessPackageRequestStateDelivered):
		testAccessPacakgeAssignmentRequestClient_Delete(t, c, *ap.ID)
	}
}

func testAccessPackageAssignmentRequestClient_Create(t *testing.T, c *test.Test, ar msgraph.AccessPackageAssignmentRequest) (request *msgraph.AccessPackageAssignmentRequest) {
	request, status, err := c.AccessPackageAssignmentRequestClient.Create(c.Context, ar)
	if err != nil {
		t.Fatalf("AccessPackageAssignementRequestClient.Create(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AccessPackageAssignementRequestClient.Create(): invalid status: %d", status)
	}
	if request == nil {
		t.Fatal("AccessPackageAssignementRequestClient.Create(): AccessPackageAssignmentRequest was nil")
	}
	if request.ID == nil {
		t.Fatal("AccessPackageAssignementRequestClient.Create(): AccessPackageAssignmentRequest.ID was nil")
	}
	return request
}

func testAccessPackageAssignmentRequestClient_Get(t *testing.T, c *test.Test, id string) {
	request, status, err := c.AccessPackageAssignmentRequestClient.Get(c.Context, id)
	if err != nil {
		t.Fatalf("AccessPackageAssignementRequestClient.Get(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AccessPackageAssignementRequestClient.Get(): invalid status: %d", status)
	}
	if request == nil {
		t.Fatal("AccessPackageAssignementRequestClient.Get(): AccessPackageAssignmentRequest was nil")
	}
	if request.ID == nil {
		t.Fatal("AccessPackageAssignementRequestClient.Get(): AccessPackageAssignmentRequest.ID was nil")
	}
}

func testAccessPackageAssignmentRequestClient_Cancel(t *testing.T, c *test.Test, id string) {
	status, err := c.AccessPackageAssignmentRequestClient.Cancel(c.Context, id)
	if err != nil {
		t.Fatalf("AccessPackageAssignmentRequestClient.Cancel(): %v", err)
	}
	if status != 200 {
		t.Fatalf("AccessPackageAssignmentRequestClient.Cancel(): invalid status: %d", status)
	}
}

func testAccessPackageAssignmentRequestClient_List(t *testing.T, c *test.Test, IDs []string) {
	requests, status, err := c.AccessPackageAssignmentRequestClient.List(c.Context, odata.Query{})
	if err != nil {
		t.Fatalf("AccessPackageAssignmentRequestClient.List(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AccessPackageAssignementRequestClient.List(): invalid status: %d", status)
	}
	count := len(IDs)

	for _, request := range *requests {
		if strContains(IDs, *request.ID) {
			count--
		}
	}
	if count != 0 {
		t.Fatalf("AccessPackageAssignementRequestClient.List(): failed to find AP requests that should have existed")
	}
}

func testAccessPacakgeAssignmentRequestClient_Delete(t *testing.T, c *test.Test, id string) {
	status, err := c.AccessPackageAssignmentRequestClient.Delete(c.Context, id)
	if err != nil {
		t.Fatalf("AccessPackageAssignmentRequestClient.Delete(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AccessPackageAssignmentRequestClient.Delete(): invalid status: %d", status)
	}
}

func strContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
