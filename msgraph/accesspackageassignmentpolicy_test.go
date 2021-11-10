package msgraph_test

import (
	"fmt"
	"testing"

	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/msgraph"
	"github.com/manicminer/hamilton/odata"
)

func TestAccessPackageAssignmentPolicyClient(t *testing.T) {
	c := test.NewTest()

	// Create test catalog
	accessPackageCatalog := testAccessPackageCatalog_Create(t, c)

	// Create AP
	accessPackage := testAccessPackage_Create(t, c, accessPackageCatalog)

	// Create Assignment Policy
	accessPackageAssignmentPolicy := testAccessPackageAssignmentPolicyClient_Create(t, c, msgraph.AccessPackageAssignmentPolicy{
		AccessPackageId: accessPackage.ID,
		DisplayName:     utils.StringPtr(fmt.Sprintf("Test-AP-Policy-Assignment-%s", c.RandomString)),
		Description:     utils.StringPtr("Test AP Policy Assignment Description"),
		//AccessReviewSettings: utils.BoolPtr()
		RequestorSettings: &msgraph.RequestorSettings{
			ScopeType:      msgraph.RequestorSettingsScopeTypeNoSubjects,
			AcceptRequests: utils.BoolPtr(true),
			//AllowedRequestors: &msgraph.UserSet{}
		},
		RequestApprovalSettings: &msgraph.ApprovalSettings{
			IsApprovalRequired:               utils.BoolPtr(false),
			IsApprovalRequiredForExtension:   utils.BoolPtr(false),
			IsRequestorJustificationRequired: utils.BoolPtr(false),
			ApprovalMode:                     msgraph.ApprovalModeNoApproval,
			//ApprovalStages: &msgraph.ApprovalStages{},
		},
	})

	testAccessPackageAssignmentPolicyClient_Get(t, c, *accessPackageAssignmentPolicy.ID)

	// Update test https://docs.microsoft.com/en-us/graph/api/accesspackageassignmentpolicy-update?view=graph-rest-beta
	newAccessPackageAssignmentPolicy := msgraph.AccessPackageAssignmentPolicy{
		ID:              accessPackageAssignmentPolicy.ID,
		AccessPackageId: accessPackageAssignmentPolicy.AccessPackageId, // Both the ID and AccessPackageId MUST Be specified. API complains vaguely as just "the Id"
		DisplayName:     utils.StringPtr(fmt.Sprintf("Test-AP-Policy-Assignment-Updated-%s", c.RandomString)),
		Description:     utils.StringPtr("Test AP Policy Assignment Description Updated"),
	}

	testAccessPackageAssignmentPolicyClient_Update(t, c, newAccessPackageAssignmentPolicy)
	testAccessPackageAssignmentPolicyClient_List(t, c)
	testAccessPackageAssignmentPolicyClient_Get(t, c, *newAccessPackageAssignmentPolicy.ID)

	//Cleanup
	testAccessPackageAssignmentPolicyClient_Delete(t, c, *newAccessPackageAssignmentPolicy.ID)
	testAccessPackage_Delete(t, c, *accessPackage.ID)
	testAccessPackageCatalog_Delete(t, c, accessPackageCatalog)
}

// AccessPackageAssignmentPolicy

func testAccessPackageAssignmentPolicyClient_Create(t *testing.T, c *test.Test, a msgraph.AccessPackageAssignmentPolicy) (accessPackageAssignmentPolicy *msgraph.AccessPackageAssignmentPolicy) {
	accessPackageAssignmentPolicy, status, err := c.AccessPackageAssignmentPolicyClient.Create(c.Connection.Context, a)
	if err != nil {
		t.Fatalf("AccessPackageAssignmentPolicyClient.Create(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AccessPackageAssignmentPolicyClient.Create(): invalid status: %d", status)
	}
	if accessPackageAssignmentPolicy == nil {
		t.Fatal("AccessPackageAssignmentPolicyClient.Create(): accessPackageAssignmentPolicy was nil")
	}
	if accessPackageAssignmentPolicy.ID == nil {
		t.Fatal("AccessPackageAssignmentPolicyClient.Create(): acccessPackageAssignmentPolicy.ID was nil")
	}
	return
}

func testAccessPackageAssignmentPolicyClient_Get(t *testing.T, c *test.Test, id string) (accessPackageAssignmentPolicy *msgraph.AccessPackageAssignmentPolicy) {
	accessPackageAssignmentPolicy, status, err := c.AccessPackageAssignmentPolicyClient.Get(c.Connection.Context, id, odata.Query{})
	if err != nil {
		t.Fatalf("AccessPackageAssignmentPolicyClient.Get(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AccessPackageAssignmentPolicyClient.Get(): invalid status: %d", status)
	}
	if accessPackageAssignmentPolicy == nil {
		t.Fatal("AccessPackageAssignmentPolicyClient.Get(): policy was nil")
	}
	return
}

func testAccessPackageAssignmentPolicyClient_Update(t *testing.T, c *test.Test, accessPackageAssignmentPolicy msgraph.AccessPackageAssignmentPolicy) {
	status, err := c.AccessPackageAssignmentPolicyClient.Update(c.Connection.Context, accessPackageAssignmentPolicy)
	if err != nil {
		t.Fatalf("AccessPackageAssignmentPolicyClient.Update(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AccessPackageAssignmentPolicyClient.Update(): invalid status: %d", status)
	}
}

func testAccessPackageAssignmentPolicyClient_List(t *testing.T, c *test.Test) (accessPackageAssignmentPolicys *[]msgraph.AccessPackageAssignmentPolicy) {
	accessPackageAssignmentPolicys, _, err := c.AccessPackageAssignmentPolicyClient.List(c.Connection.Context, odata.Query{Top: 10})
	if err != nil {
		t.Fatalf("AccessPackageAssignmentPolicyClient.List(): %v", err)
	}
	if accessPackageAssignmentPolicys == nil {
		t.Fatal("AccessPackageAssignmentPolicyClient.List(): accessPackageAssignmentPolicys was nil")
	}
	return
}

func testAccessPackageAssignmentPolicyClient_Delete(t *testing.T, c *test.Test, id string) {
	status, err := c.AccessPackageAssignmentPolicyClient.Delete(c.Connection.Context, id)
	if err != nil {
		t.Fatalf("AccessPackageAssignmentPolicyClient.Delete(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AccessPackageAssignmentPolicyClient.Delete(): invalid status: %d", status)
	}
}

// AccessPackage

func testAccessPackage_Create(t *testing.T, c *test.Test, accessPackageCatalog *msgraph.AccessPackageCatalog) (accessPackage *msgraph.AccessPackage) {
	accessPackage, _, err := c.AccessPackageClient.Create(c.Connection.Context, msgraph.AccessPackage{
		DisplayName:         utils.StringPtr(fmt.Sprintf("test-accesspackage-%s", c.RandomString)),
		CatalogId:           accessPackageCatalog.ID,
		Description:         utils.StringPtr("Test Access Package"),
		IsHidden:            utils.BoolPtr(false),
		IsRoleScopesVisible: utils.BoolPtr(false),
	})

	if err != nil {
		t.Fatalf("AccessPackageClient.Create() - Could not create test AccessPackage catalog: %v", err)
	}
	return
}

func testAccessPackage_Delete(t *testing.T, c *test.Test, id string) {
	_, err := c.AccessPackageClient.Delete(c.Connection.Context, id)
	if err != nil {
		t.Fatalf("AccessPackageClient.Delete() - Could not delete test AccessPackage catalog")
	}
}

// AccessPackageCatalog

func testAccessPackageCatalog_Create(t *testing.T, c *test.Test) (accessPackageCatalog *msgraph.AccessPackageCatalog) {
	accessPackageCatalog, _, err := c.AccessPackageCatalogClient.Create(c.Connection.Context, msgraph.AccessPackageCatalog{
		DisplayName:         utils.StringPtr(fmt.Sprintf("test-catalog-%s", c.RandomString)),
		CatalogType:         msgraph.AccessPackageCatalogTypeUserManaged,
		CatalogStatus:       msgraph.AccessPackageCatalogStatusPublished,
		Description:         utils.StringPtr("Test Access Catalog"),
		IsExternallyVisible: utils.BoolPtr(false),
	})

	if err != nil {
		t.Fatalf("AccessPackageCatalogClient.Create() - Could not create test AccessPackage catalog: %v", err)
	}
	return
}

func testAccessPackageCatalog_Delete(t *testing.T, c *test.Test, accessPackageCatalog *msgraph.AccessPackageCatalog) {
	_, err := c.AccessPackageCatalogClient.Delete(c.Connection.Context, *accessPackageCatalog.ID)
	if err != nil {
		t.Fatalf("AccessPackageCatalogClient.Delete() - Could not delete test AccessPackage catalog")
	}
}
