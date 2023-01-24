package msgraph_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/msgraph"
)

func TestAccessPackageAssignmentRequestClient(t *testing.T) {
	c := test.NewTest(t)
	defer c.CancelFunc()

	// Create test Catalog
	accessPackageCatalog := testAccessPackageCatalog_Create(t, c)

	// Create AP
	accessPackage := testAccessPackage_Create(t, c, accessPackageCatalog)

	currentTimePlusDay := time.Now().AddDate(0, 0, 1)

	user := testUsersClient_Create(t, c, msgraph.User{
		AccountEnabled:    utils.BoolPtr(true),
		DisplayName:       utils.StringPtr("test-user"),
		MailNickname:      utils.StringPtr(fmt.Sprintf("test-user-%s", c.RandomString)),
		UserPrincipalName: utils.StringPtr(fmt.Sprintf("test-user-%s@%s", c.RandomString, c.Connections["default"].DomainName)),
		PasswordProfile: &msgraph.UserPasswordProfile{
			Password: utils.StringPtr(fmt.Sprintf("IrPa55w0rd%s", c.RandomString)),
		},
	})

	// Create Assignment Policy
	accessPackageAssignmentPolicy := testAccessPackageAssignmentPolicyClient_Create(t, c, msgraph.AccessPackageAssignmentPolicy{
		AccessPackageId: accessPackage.ID,
		AccessReviewSettings: &msgraph.AssignmentReviewSettings{
			AccessReviewTimeoutBehavior:     msgraph.AccessReviewTimeoutBehaviorTypeRemoveAccess,
			IsEnabled:                       utils.BoolPtr(true),
			StartDateTime:                   &currentTimePlusDay,
			DurationInDays:                  utils.Int32Ptr(5),
			RecurrenceType:                  msgraph.AccessReviewRecurranceTypeMonthly,
			ReviewerType:                    msgraph.AccessReviewReviewerTypeSelf,
			IsAccessRecommendationEnabled:   utils.BoolPtr(true),
			IsApprovalJustificationRequired: utils.BoolPtr(true),
		},
		DisplayName: utils.StringPtr(fmt.Sprintf("Test-AP-Policy-Assignment-%s", c.RandomString)),
		Description: utils.StringPtr("Test AP Policy Assignment Description"),
		//AccessReviewSettings: utils.BoolPtr()
		RequestorSettings: &msgraph.RequestorSettings{
			ScopeType:      msgraph.RequestorSettingsScopeTypeNoSubjects,
			AcceptRequests: utils.BoolPtr(true),
		},
		RequestApprovalSettings: &msgraph.ApprovalSettings{
			IsApprovalRequired:               utils.BoolPtr(false),
			IsApprovalRequiredForExtension:   utils.BoolPtr(false),
			IsRequestorJustificationRequired: utils.BoolPtr(false),
			ApprovalMode:                     msgraph.ApprovalModeNoApproval,
		},
		Questions: &[]msgraph.AccessPackageQuestion{},
	})

	accessPackageGet := testAccessPackageClient_Get(t, c, *accessPackage.ID)
	userGet := testUsersClient_Get(t, c, *user.Id)
	policyGetID := testAccessPackageAssignmentPolicyClient_Get(t, c, *accessPackageAssignmentPolicy.ID)

	ap := testAccessPackageAssignmentRequestClient_Create(t, c, msgraph.AccessPackageAssignmentRequest{
		RequestType: utils.StringPtr(msgraph.AccessPacakgeRequestTypeAdminAdd),
		AccessPackageAssignment: &msgraph.AccessPackageAssignment{
			TargetID:            userGet.Id,
			AssignementPolicyID: policyGetID.ID,
			AccessPackageID:     accessPackageGet.ID,
		},
	})

	updatedAPRequest := testAccessPackageAssignmentRequestClient_Get(t, c, *ap.ID)
	// Can only delete a request if it is in specific states
	switch updatedAPRequest.State {
	case utils.StringPtr(msgraph.AccessPackageRequestStateDenied):
		testAccessPacakgeAssignmentRequestClient_Delete(t, c, *updatedAPRequest.ID)
	case utils.StringPtr(msgraph.AccessPackageRequestStateCanceled):
		testAccessPacakgeAssignmentRequestClient_Delete(t, c, *updatedAPRequest.ID)
	case utils.StringPtr(msgraph.AccessPackageRequestStateDelivered):
		testAccessPacakgeAssignmentRequestClient_Delete(t, c, *updatedAPRequest.ID)
	}

	//Cleanup
	testUser_Delete(t, c, user)
	testAccessPackageAssignmentPolicyClient_Delete(t, c, *accessPackageAssignmentPolicy.ID)
	testAccessPackage_Delete(t, c, *accessPackage.ID)
	testAccessPackageCatalog_Delete(t, c, accessPackageCatalog)

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

func testAccessPackageAssignmentRequestClient_Get(t *testing.T, c *test.Test, id string) (request *msgraph.AccessPackageAssignmentRequest) {
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
	return request

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
