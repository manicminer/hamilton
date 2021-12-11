package msgraph_test

import (
	"fmt"
	"testing"

	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/msgraph"
	"github.com/manicminer/hamilton/odata"
)

func TestAccessPackageResourceRequestClient(t *testing.T) {
	c := test.NewTest(t)
	defer c.CancelFunc()

	self := testDirectoryObjectsClient_Get(t, c, c.Claims.ObjectId)

	// Create group
	aadGroup := testAccessPackageResourceRequestGroup_Create(t, c, msgraph.Owners{*self})

	// Create test catalog
	accessPackageCatalog := testAccessPackageResourceRequestCatalog_Create(t, c)

	// Create access package
	accessPackage := testAccessPackageResourceRequestAP_Create(t, c, msgraph.AccessPackage{
		DisplayName: utils.StringPtr(fmt.Sprintf("test-accesspackage-%s", c.RandomString)),
		Catalog: &msgraph.AccessPackageCatalog{
			ID: accessPackageCatalog.ID,
		},
		Description: utils.StringPtr("Test Access Package"),
		IsHidden:    utils.BoolPtr(false),
	})

	// Create resource request and poll for ID
	accessPackageResourceRequest := testAccessPackageResourceRequestClient_Create(t, c, msgraph.AccessPackageResourceRequest{
		CatalogId:   accessPackage.Catalog.ID,
		RequestType: utils.StringPtr("AdminAdd"),
		AccessPackageResource: &msgraph.AccessPackageResource{
			OriginId:     aadGroup.ID,
			OriginSystem: msgraph.AccessPackageResourceOriginSystemAadGroup,
		},
	}, true)

	// Resource client
	testAccessPackageResourceRequestResource_Get(t, c, *accessPackageResourceRequest.CatalogId, *accessPackageResourceRequest.AccessPackageResource.OriginId)
	testAccessPackageResourceRequestResource_List(t, c, *accessPackageResourceRequest.CatalogId)

	// Requests client
	testAccessPackageResourceRequestClient_List(t, c)
	testAccessPackageResourceRequestClient_Get(t, c, *accessPackageResourceRequest.ID)
	testAccessPackageResourceRequestClient_Delete(t, c, accessPackageResourceRequest)

	// Cleanup
	testAccessPackageResourceRequestAP_Delete(t, c, *accessPackage.ID)
	testAccessPackageResourceRequestCatalog_Delete(t, c, *accessPackageCatalog.ID)
	testAccessPackageResourceRequestGroup_Delete(t, c, aadGroup)
}

// AccessPackageResourceRequest

func testAccessPackageResourceRequestClient_Create(t *testing.T, c *test.Test, a msgraph.AccessPackageResourceRequest, pollForId bool) (accessPackageResourceRequest *msgraph.AccessPackageResourceRequest) {
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
		t.Fatal("AccessPackageResourceRequestClient.Create(): acccessPackageResourceRequest.ID was nil")
	}
	return
}

func testAccessPackageResourceRequestClient_Get(t *testing.T, c *test.Test, id string) (accessPackageResourceRequest *msgraph.AccessPackageResourceRequest) {
	accessPackageResourceRequest, status, err := c.AccessPackageResourceRequestClient.Get(c.Context, id)
	if err != nil {
		t.Fatalf("AccessPackageResourceRequestClient.Get(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AccessPackageResourceRequestClient.Get(): invalid status: %d", status)
	}
	if accessPackageResourceRequest == nil {
		t.Fatal("AccessPackageResourceRequestClient.Get(): policy was nil")
	}
	return
}

func testAccessPackageResourceRequestClient_List(t *testing.T, c *test.Test) (accessPackageResourceRequests *[]msgraph.AccessPackageResourceRequest) {
	accessPackageResourceRequests, _, err := c.AccessPackageResourceRequestClient.List(c.Context, odata.Query{Top: 10})
	if err != nil {
		t.Fatalf("AccessPackageResourceRequestClient.List(): %v", err)
	}
	if accessPackageResourceRequests == nil {
		t.Fatal("AccessPackageResourceRequestClient.List(): accessPackageResourceRequests was nil")
	}
	return
}

func testAccessPackageResourceRequestClient_Delete(t *testing.T, c *test.Test, accessPackageResourceRequest *msgraph.AccessPackageResourceRequest) {
	status, err := c.AccessPackageResourceRequestClient.Delete(c.Context, *accessPackageResourceRequest)
	if err != nil {
		t.Fatalf("AccessPackageResourceRequestClient.Delete(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AccessPackageResourceRequestClient.Delete(): invalid status: %d", status)
	}
}

// AccessPackage

func testAccessPackageResourceRequestAP_Create(t *testing.T, c *test.Test, a msgraph.AccessPackage) (accessPackage *msgraph.AccessPackage) {
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
		t.Fatal("AccessPackageClient.Create(): acccessPackage.ID was nil")
	}
	return
}

func testAccessPackageResourceRequestAP_Delete(t *testing.T, c *test.Test, id string) {
	status, err := c.AccessPackageClient.Delete(c.Context, id)
	if err != nil {
		t.Fatalf("AccessPackageClient.Delete(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AccessPackageClient.Delete(): invalid status: %d", status)
	}
}

// AccessPackageCatalog

func testAccessPackageResourceRequestCatalog_Create(t *testing.T, c *test.Test) (accessPackageCatalog *msgraph.AccessPackageCatalog) {
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

func testAccessPackageResourceRequestCatalog_Delete(t *testing.T, c *test.Test, id string) {
	_, err := c.AccessPackageCatalogClient.Delete(c.Context, id)
	if err != nil {
		t.Fatalf("AccessPackageCatalogClient.Delete() - Could not delete test AccessPackage catalog")
	}
}

func testAccessPackageResourceRequestGroup_Create(t *testing.T, c *test.Test, owners msgraph.Owners) (group *msgraph.Group) {
	group, _, err := c.GroupsClient.Create(c.Context, msgraph.Group{
		DisplayName:     utils.StringPtr(fmt.Sprintf("%s-%s", "testapresourcerequest", c.RandomString)),
		MailEnabled:     utils.BoolPtr(false),
		MailNickname:    utils.StringPtr(fmt.Sprintf("%s-%s", "testapresourcerequest", c.RandomString)),
		SecurityEnabled: utils.BoolPtr(true),
		Owners:          &owners,
	})

	if err != nil {
		t.Fatalf("GroupsClient.Create() - Could not create test group: %v", err)
	}
	return
}

func testAccessPackageResourceRequestGroup_Delete(t *testing.T, c *test.Test, group *msgraph.Group) {
	_, err := c.GroupsClient.Delete(c.Context, *group.ID)
	if err != nil {
		t.Fatalf("GroupsClient.Delete() - Could not delete test group: %v", err)
	}
}

// AccessPackageResource

func testAccessPackageResourceRequestResource_Get(t *testing.T, c *test.Test, catalogId string, originId string) (accessPackageResource *msgraph.AccessPackageResource) {
	accessPackageResource, status, err := c.AccessPackageResourceClient.Get(c.Context, catalogId, originId)

	if err != nil {
		t.Fatalf("AccessPackageCatalogClient.Get(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AccessPackageCatalogClient.Get(): invalid status: %d", status)
	}
	if accessPackageResource == nil {
		t.Fatal("AccessPackageCatalogClient.Get(): policy was nil")
	}

	return
}

func testAccessPackageResourceRequestResource_List(t *testing.T, c *test.Test, catalogId string) (accessPackageResources *[]msgraph.AccessPackageResource) {
	accessPackageResources, status, err := c.AccessPackageResourceClient.List(c.Context, catalogId, odata.Query{Top: 10})

	if err != nil {
		t.Fatalf("AccessPackageCatalogClient.Get(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AccessPackageCatalogClient.Get(): invalid status: %d", status)
	}
	if accessPackageResources == nil {
		t.Fatal("AccessPackageCatalogClient.Get(): policys was nil")
	}

	return
}
