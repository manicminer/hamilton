package msgraph_test

import (
	"fmt"
	"testing"

	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/msgraph"
)

func TestAccessPackageResourceRoleClient(t *testing.T) {
	c := test.NewTest(t)
	defer c.CancelFunc()

	self := testDirectoryObjectsClient_Get(t, c, c.Claims.ObjectId)

	// Create group
	aadGroup := testAccessPackageResourceRoleGroup_Create(t, c, msgraph.Owners{*self})

	// Create test catalog
	accessPackageCatalog := testAccessPackageResourceRoleCatalog_Create(t, c)

	// Create access package
	accessPackage := testAccessPackageResourceRoleAP_Create(t, c, msgraph.AccessPackage{
		DisplayName: utils.StringPtr(fmt.Sprintf("test-accesspackage-%s", c.RandomString)),
		Catalog: &msgraph.AccessPackageCatalog{
			ID: accessPackageCatalog.ID,
		},
		Description: utils.StringPtr("Test Access Package"),
		IsHidden:    utils.BoolPtr(false),
	})

	// Create Resource Request and poll for ID
	accessPackageResourceRequest := testAccessPackageResourceRoleResourceRequest_Create(t, c, msgraph.AccessPackageResourceRequest{
		CatalogId:   accessPackage.Catalog.ID,
		RequestType: utils.StringPtr("AdminAdd"),
		AccessPackageResource: &msgraph.AccessPackageResource{
			OriginId:     aadGroup.ID(),
			OriginSystem: msgraph.AccessPackageResourceOriginSystemAadGroup,
			//ResourceType: utils.StringPtr("Security Group") // This is not mandatory for groups but is seen in sharepoint emails
		},
	}, true)

	// Try to get roles for group we added to Catalog
	testAccessPackageResourceRoleClient_Get(t, c, *accessPackage.Catalog.ID, msgraph.AccessPackageResourceOriginSystemAadGroup, *accessPackageResourceRequest.AccessPackageResource.ID)

	// Cleanup
	testAccessPackageResourceRoleAP_Delete(t, c, *accessPackage.ID)
	testAccessPackageResourceRoleResourceRequest_Delete(t, c, accessPackageResourceRequest)
	testAccessPackageResourceRoleCatalog_Delete(t, c, *accessPackageCatalog.ID)
	testAccessPackageResourceRoleGroup_Delete(t, c, aadGroup)
}

// AccessPackageResourceRole
func testAccessPackageResourceRoleClient_Get(t *testing.T, c *test.Test, catalogId string, originSystem msgraph.AccessPackageResourceOriginSystem, accessPackageResourceId string) (accessPackageResourceRoleScope *msgraph.AccessPackageResourceRoleScope) {
	accessPackageResourceRole, status, err := c.AccessPackageResourceRoleClient.Get(c.Context, catalogId, originSystem, accessPackageResourceId)
	if err != nil {
		t.Fatalf("AccessPackageResourceRequestClient.Get(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AccessPackageResourceRequestClient.Get(): invalid status: %d", status)
	}
	if accessPackageResourceRole == nil {
		t.Fatal("AccessPackageResourceRequestClient.Get(): policy was nil")
	}
	return
}

// AccessPackageResourceRequest
func testAccessPackageResourceRoleResourceRequest_Create(t *testing.T, c *test.Test, a msgraph.AccessPackageResourceRequest, pollForId bool) (accessPackageResourceRequest *msgraph.AccessPackageResourceRequest) {
	accessPackageResourceRequest, status, err := c.AccessPackageResourceRequestClient.Create(c.Context, a, pollForId)
	if err != nil {
		t.Fatalf("AccessPackageResourceRequestClient.Create(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AccessPackageResourceRequestClient.Create(): invalid status: %d", status)
	}
	if accessPackageResourceRequest == nil {
		t.Fatal("AccessPackageResourceRequestClient.Create(): accessPackageResourceRequest was nil")
	}
	if accessPackageResourceRequest.ID == nil {
		t.Fatal("AccessPackageResourceRequestClient.Create(): accessPackageResourceRequest.ID was nil")
	}
	return
}

func testAccessPackageResourceRoleResourceRequest_Delete(t *testing.T, c *test.Test, accessPackageResourceRequest *msgraph.AccessPackageResourceRequest) {
	status, err := c.AccessPackageResourceRequestClient.Delete(c.Context, *accessPackageResourceRequest)
	if err != nil {
		t.Fatalf("AccessPackageResourceRequestClient.Delete(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AccessPackageResourceRequestClient.Delete(): invalid status: %d", status)
	}
}

// AccessPackage
func testAccessPackageResourceRoleAP_Create(t *testing.T, c *test.Test, a msgraph.AccessPackage) (accessPackage *msgraph.AccessPackage) {
	accessPackage, status, err := c.AccessPackageClient.Create(c.Context, a)
	if err != nil {
		t.Fatalf("AccessPackageClient.Create(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AccessPackageClient.Create(): invalid status: %d", status)
	}
	if accessPackage == nil {
		t.Fatal("AccessPackageClient.Create(): accessPackage was nil")
	}
	if accessPackage.ID == nil {
		t.Fatal("AccessPackageClient.Create(): accessPackage.ID was nil")
	}
	return
}

func testAccessPackageResourceRoleAP_Delete(t *testing.T, c *test.Test, id string) {
	status, err := c.AccessPackageClient.Delete(c.Context, id)
	if err != nil {
		t.Fatalf("AccessPackageClient.Delete(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AccessPackageClient.Delete(): invalid status: %d", status)
	}
}

// AccessPackageCatalog
func testAccessPackageResourceRoleCatalog_Create(t *testing.T, c *test.Test) (accessPackageCatalog *msgraph.AccessPackageCatalog) {
	accessPackageCatalog, _, err := c.AccessPackageCatalogClient.Create(c.Context, msgraph.AccessPackageCatalog{
		DisplayName:         utils.StringPtr(fmt.Sprintf("test-catalog-%s", c.RandomString)),
		CatalogType:         msgraph.AccessPackageCatalogTypeUserManaged,
		State:               msgraph.AccessPackageCatalogStatePublished,
		Description:         utils.StringPtr("Test Access Catalog"),
		IsExternallyVisible: utils.BoolPtr(false),
	})

	if err != nil {
		t.Fatalf("AccessPackageCatalogClient.Create() - Could not create test AccessPackage catalog: %v", err)
	}
	return
}

func testAccessPackageResourceRoleCatalog_Delete(t *testing.T, c *test.Test, id string) {
	_, err := c.AccessPackageCatalogClient.Delete(c.Context, id)
	if err != nil {
		t.Fatalf("AccessPackageCatalogClient.Delete() - Could not delete test AccessPackage catalog")
	}
}

func testAccessPackageResourceRoleGroup_Create(t *testing.T, c *test.Test, self msgraph.Owners) (group *msgraph.Group) {
	group, _, err := c.GroupsClient.Create(c.Context, msgraph.Group{
		DisplayName:     utils.StringPtr(fmt.Sprintf("%s-%s", "testapresourcerequest", c.RandomString)),
		MailEnabled:     utils.BoolPtr(false),
		MailNickname:    utils.StringPtr(fmt.Sprintf("%s-%s", "testapresourcerequest", c.RandomString)),
		SecurityEnabled: utils.BoolPtr(true),
		Owners:          &self,
	})

	if err != nil {
		t.Fatalf("GroupsClient.Create() - Could not create test group: %v", err)
	}
	return
}

func testAccessPackageResourceRoleGroup_Delete(t *testing.T, c *test.Test, group *msgraph.Group) {
	_, err := c.GroupsClient.Delete(c.Context, *group.ID())
	if err != nil {
		t.Fatalf("GroupsClient.Delete() - Could not delete test group: %v", err)
	}
}
