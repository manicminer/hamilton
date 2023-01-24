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
			// Reviewers: &[]msgraph.UserSet{
			// 	{
			// 		ODataType:    utils.StringPtr(odata.TypeRequestorManager),
			// 		IsBackup:     utils.BoolPtr(false),
			// 		ManagerLevel: utils.Int32Ptr(1),
			// 	},
			// 	{
			// 		ODataType:    utils.StringPtr(odata.TypeSingleUser),
			// 		IsBackup:     utils.BoolPtr(true),
			// 		ID:           utils.StringPtr(""),
			// 	},
			// },
		},
		DisplayName: utils.StringPtr(fmt.Sprintf("Test-AP-Policy-Assignment-%s", c.RandomString)),
		Description: utils.StringPtr("Test AP Policy Assignment Description"),
		//AccessReviewSettings: utils.BoolPtr()
		RequestorSettings: &msgraph.RequestorSettings{
			//ScopeType:      msgraph.RequestorSettingsScopeTypeSpecificDirectorySubjects,
			ScopeType:      msgraph.RequestorSettingsScopeTypeNoSubjects,
			AcceptRequests: utils.BoolPtr(true),
			// AllowedRequestors: &[]msgraph.UserSet{
			// 		{
			// 			ODataType: utils.StringPtr(odata.TypeGroupMembers),
			// 			IsBackup: utils.BoolPtr(false),
			// 			ID: utils.StringPtr("xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"),
			// 			Description: utils.StringPtr("Sample users group"),
			// 		},
			// },
		},
		RequestApprovalSettings: &msgraph.ApprovalSettings{
			IsApprovalRequired:               utils.BoolPtr(false),
			IsApprovalRequiredForExtension:   utils.BoolPtr(false),
			IsRequestorJustificationRequired: utils.BoolPtr(false),
			ApprovalMode:                     msgraph.ApprovalModeNoApproval,
			//ApprovalStages: &msgraph.ApprovalStages{},
		},
		Questions: &[]msgraph.AccessPackageQuestion{
			{
				ODataType:  utils.StringPtr(odata.TypeAccessPackageTextInputQuestion),
				IsRequired: utils.BoolPtr(false),
				Sequence:   utils.Int32Ptr(1),
				Text: &msgraph.AccessPackageLocalizedContent{
					DefaultText: utils.StringPtr("Test"),
					LocalizedTexts: &[]msgraph.AccessPackageLocalizedTexts{
						{
							Text:         utils.StringPtr("abc"),
							LanguageCode: utils.StringPtr("en"),
						},
					},
				},
			},
			{
				ODataType:  utils.StringPtr(odata.TypeAccessPackageMultipleChoiceQuestion),
				IsRequired: utils.BoolPtr(false),
				Sequence:   utils.Int32Ptr(2),
				Text: &msgraph.AccessPackageLocalizedContent{
					DefaultText: utils.StringPtr("Test"),
					LocalizedTexts: &[]msgraph.AccessPackageLocalizedTexts{
						{
							Text:         utils.StringPtr("abc 2"),
							LanguageCode: utils.StringPtr("gb"),
						},
					},
				},
				Choices: &[]msgraph.AccessPackageMultipleChoiceQuestions{
					// Choice 1 containing a list of languages
					{
						ActualValue: utils.StringPtr("CHOICE1"),
						DisplayValue: &msgraph.AccessPackageLocalizedContent{
							DefaultText: utils.StringPtr("One"),
							LocalizedTexts: &[]msgraph.AccessPackageLocalizedTexts{
								{
									Text:         utils.StringPtr("Choice 1"),
									LanguageCode: utils.StringPtr("gb"),
								},
							},
						},
					},
					// Choice 2 containing a list of languages, etc.
					{
						ActualValue: utils.StringPtr("CHOICE2"),
						DisplayValue: &msgraph.AccessPackageLocalizedContent{
							DefaultText: utils.StringPtr("Two"),
							LocalizedTexts: &[]msgraph.AccessPackageLocalizedTexts{
								{
									Text:         utils.StringPtr("Choice 2"),
									LanguageCode: utils.StringPtr("gb"),
								},
								{
									Text:         utils.StringPtr("Zwei"),
									LanguageCode: utils.StringPtr("de"),
								},
							},
						},
					},
				},
			},
		},
	})

	// accessPackageAssignementRequest := msgraph.AccessPackageAssignmentRequest{
	// 	RequestType: utils.StringPtr(msgraph.AccessPacakgeRequestTypeAdminAdd),
	// 	AccessPackageAssignment: &msgraph.AccessPackageAssignment{
	// 		TargetID:            user.ID(),
	// 		AssignementPolicyID: accessPackageAssignmentPolicy.ID,
	// 		AccessPackageID:     accessPackage.ID,
	// 	},
	// }

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
	testAccessPacakgeAssignmentRequestClient_Delete(t, c, *ap.ID)

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

func testAccessPacakgeAssignmentRequestClient_Delete(t *testing.T, c *test.Test, id string) {
	status, err := c.AccessPackageAssignmentRequestClient.Delete(c.Context, id)
	if err != nil {
		t.Fatalf("AccessPackageAssignmentRequestClient.Delete(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AccessPackageAssignmentRequestClient.Delete(): invalid status: %d", status)
	}
}
