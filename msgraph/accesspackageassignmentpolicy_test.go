package msgraph_test

import (
	"fmt"
	"testing"

	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/msgraph"
	"github.com/manicminer/hamilton/odata"
)

type AccessPackageAssignmentPolicyTest struct {
	connection               *test.Connection
	apClient                 *msgraph.AccessPackageClient        //apClient
	apCatalogClient          *msgraph.AccessPackageCatalogClient //Client for Catalog Test to associate as required
	apAssignmentPolicyClient *msgraph.AccessPackageAssignmentPolicyClient
	randomString             string
}

func TestAccessPackageAssignmentPolicyClient(t *testing.T) {
	c := AccessPackageAssignmentPolicyTest{
		connection:   test.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: test.RandomString(),
	}

	// Init clients
	c.apClient = msgraph.NewAccessPackageClient(c.connection.AuthConfig.TenantID)
	c.apClient.BaseClient.Authorizer = c.connection.Authorizer

	c.apCatalogClient = msgraph.NewAccessPackageCatalogClient(c.connection.AuthConfig.TenantID)
	c.apCatalogClient.BaseClient.Authorizer = c.connection.Authorizer

	c.apAssignmentPolicyClient = msgraph.NewAccessPackageAssignmentPolicyClient(c.connection.AuthConfig.TenantID)
	c.apAssignmentPolicyClient.BaseClient.Authorizer = c.connection.Authorizer

	// Create test catalog
	accessPackageCatalog := testAccessPackageCatalog_Create(t, c)

	// Create AP
	accessPackage := testAccessPackage_Create(t, c, accessPackageCatalog)

	// Create Assignment Policy
	accessPackageAssignmentPolicy := testAccessPackageAssignmentPolicyClient_Create(t, c, msgraph.AccessPackageAssignmentPolicy{
		AccessPackageId: accessPackage.ID,
		DisplayName:     utils.StringPtr(fmt.Sprintf("Test-AP-Policy-Assignment-%s", c.randomString)),
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
		DisplayName:     utils.StringPtr(fmt.Sprintf("Test-AP-Policy-Assignment-Updated-%s", c.randomString)),
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

func testAccessPackageAssignmentPolicyClient_Create(t *testing.T, c AccessPackageAssignmentPolicyTest, a msgraph.AccessPackageAssignmentPolicy) (accessPackageAssignmentPolicy *msgraph.AccessPackageAssignmentPolicy) {
	accessPackageAssignmentPolicy, status, err := c.apAssignmentPolicyClient.Create(c.connection.Context, a)
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

func testAccessPackageAssignmentPolicyClient_Get(t *testing.T, c AccessPackageAssignmentPolicyTest, id string) (accessPackageAssignmentPolicy *msgraph.AccessPackageAssignmentPolicy) {
	accessPackageAssignmentPolicy, status, err := c.apAssignmentPolicyClient.Get(c.connection.Context, id, odata.Query{})
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

func testAccessPackageAssignmentPolicyClient_Update(t *testing.T, c AccessPackageAssignmentPolicyTest, accessPackageAssignmentPolicy msgraph.AccessPackageAssignmentPolicy) {
	status, err := c.apAssignmentPolicyClient.Update(c.connection.Context, accessPackageAssignmentPolicy)
	if err != nil {
		t.Fatalf("AccessPackageAssignmentPolicyClient.Update(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AccessPackageAssignmentPolicyClient.Update(): invalid status: %d", status)
	}
}

func testAccessPackageAssignmentPolicyClient_List(t *testing.T, c AccessPackageAssignmentPolicyTest) (accessPackageAssignmentPolicys *[]msgraph.AccessPackageAssignmentPolicy) {
	accessPackageAssignmentPolicys, _, err := c.apAssignmentPolicyClient.List(c.connection.Context, odata.Query{Top: 10})
	if err != nil {
		t.Fatalf("AccessPackageAssignmentPolicyClient.List(): %v", err)
	}
	if accessPackageAssignmentPolicys == nil {
		t.Fatal("AccessPackageAssignmentPolicyClient.List(): accessPackageAssignmentPolicys was nil")
	}
	return
}

func testAccessPackageAssignmentPolicyClient_Delete(t *testing.T, c AccessPackageAssignmentPolicyTest, id string) {
	status, err := c.apAssignmentPolicyClient.Delete(c.connection.Context, id)
	if err != nil {
		t.Fatalf("AccessPackageAssignmentPolicyClient.Delete(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AccessPackageAssignmentPolicyClient.Delete(): invalid status: %d", status)
	}
}

// AccessPackage

func testAccessPackage_Create(t *testing.T, c AccessPackageAssignmentPolicyTest, accessPackageCatalog *msgraph.AccessPackageCatalog) (accessPackage *msgraph.AccessPackage) {
	accessPackage, _, err := c.apClient.Create(c.connection.Context, msgraph.AccessPackage{
		DisplayName:         utils.StringPtr(fmt.Sprintf("test-accesspackage-%s", c.randomString)),
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

func testAccessPackage_Delete(t *testing.T, c AccessPackageAssignmentPolicyTest, id string) {
	_, err := c.apClient.Delete(c.connection.Context, id)
	if err != nil {
		t.Fatalf("AccessPackageClient.Delete() - Could not delete test AccessPackage catalog")
	}
}

// AccessPackageCatalog

func testAccessPackageCatalog_Create(t *testing.T, c AccessPackageAssignmentPolicyTest) (accessPackageCatalog *msgraph.AccessPackageCatalog) {
	accessPackageCatalog, _, err := c.apCatalogClient.Create(c.connection.Context, msgraph.AccessPackageCatalog{
		DisplayName:         utils.StringPtr(fmt.Sprintf("test-catalog-%s", c.randomString)),
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

func testAccessPackageCatalog_Delete(t *testing.T, c AccessPackageAssignmentPolicyTest, accessPackageCatalog *msgraph.AccessPackageCatalog) {
	_, err := c.apCatalogClient.Delete(c.connection.Context, *accessPackageCatalog.ID)
	if err != nil {
		t.Fatalf("AccessPackageCatalogClient.Delete() - Could not delete test AccessPackage catalog")
	}
}
