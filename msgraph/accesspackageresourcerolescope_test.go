package msgraph_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/msgraph"
)

func TestAccessPackageResourceRoleScopeClient(t *testing.T) {
	c := test.NewTest(t)
	defer c.CancelFunc()

	self := testDirectoryObjectsClient_Get(t, c, c.Claims.ObjectId)

	// Create group
	aadGroup := testAccessPackageResourceRoleScopeGroup_Create(t, c, msgraph.Owners{*self})

	// Create test catalog
	accessPackageCatalog := testAccessPackageResourceRoleScopeCatalog_Create(t, c)

	// Create access package
	accessPackage := testAccessPackageResourceRoleScopeAP_Create(t, c, msgraph.AccessPackage{
		DisplayName: utils.StringPtr(fmt.Sprintf("test-accesspackage-%s", c.RandomString)),
		Catalog: &msgraph.AccessPackageCatalog{
			ID: accessPackageCatalog.ID,
		},
		Description: utils.StringPtr("Test Access Package"),
		IsHidden:    utils.BoolPtr(false),
	})

	// Create Resource Request and poll for ID
	accessPackageResourceRequest := testAccessPackageResourceRoleScopeResourceRequest_Create(t, c, msgraph.AccessPackageResourceRequest{
		CatalogId:   accessPackage.Catalog.ID,
		RequestType: utils.StringPtr("AdminAdd"),
		AccessPackageResource: &msgraph.AccessPackageResource{
			OriginId:     aadGroup.ID(),
			OriginSystem: msgraph.AccessPackageResourceOriginSystemAadGroup,
			//ResourceType: utils.StringPtr("Security Group") // This is not mandatory for groups but is seen in sharepoint emails
		},
	}, true)

	accessPackageResourceRoleScope := testAccessPackageResourceRoleScopeClient_Create(t, c, msgraph.AccessPackageResourceRoleScope{
		AccessPackageId: accessPackage.ID,
		AccessPackageResourceRole: &msgraph.AccessPackageResourceRole{
			//ID: utils.StringPtr("405b1a50-7e4c-4f82-ae46-b9e8ec0eb1e0"),
			DisplayName:  utils.StringPtr("Member"),
			OriginId:     utils.StringPtr(fmt.Sprintf("Member_%s", *accessPackageResourceRequest.AccessPackageResource.OriginId)),
			OriginSystem: accessPackageResourceRequest.AccessPackageResource.OriginSystem,
			AccessPackageResource: &msgraph.AccessPackageResource{
				ID:           accessPackageResourceRequest.AccessPackageResource.ID,           //Requires poll for ID
				ResourceType: accessPackageResourceRequest.AccessPackageResource.ResourceType, //requires poll for ID
				OriginId:     accessPackageResourceRequest.AccessPackageResource.OriginId,
				//OriginSystem: msgraph.AccessPackageResourceOriginSystemAadGroup,
			},
		},
		AccessPackageResourceScope: &msgraph.AccessPackageResourceScope{
			//ID: "12fec290-4dc6-4e82-88d6-7b5af842e0a3",
			OriginSystem: accessPackageResourceRequest.AccessPackageResource.OriginSystem,
			OriginId:     accessPackageResourceRequest.AccessPackageResource.OriginId,
		},
	})

	testAccessPackageResourceRoleScopeClient_Get(t, c, *accessPackageResourceRoleScope.AccessPackageId, *accessPackageResourceRoleScope.ID)
	testAccessPackageResourceRoleScopeResource_Get(t, c, *accessPackageResourceRequest.CatalogId, *accessPackageResourceRequest.AccessPackageResource.OriginId)
	testAccessPackageResourceRoleScopeClient_List(t, c, *accessPackageResourceRoleScope.AccessPackageId)
	testAccessPackageResourceRoleScopeClient_Delete(t, c, *accessPackageResourceRoleScope.AccessPackageId, *accessPackageResourceRoleScope.ID)

	// Force-replacement scenario
	testAccessPackageResourceRoleScopeAP_Delete(t, c, *accessPackage.ID)
	accessPackage = testAccessPackageResourceRoleScopeAP_Create(t, c, *accessPackage) //New ID

	testAccessPackageResourceRoleScopeClient_Create(t, c, msgraph.AccessPackageResourceRoleScope{
		AccessPackageId: accessPackage.ID,
		AccessPackageResourceRole: &msgraph.AccessPackageResourceRole{
			//ID: utils.StringPtr("405b1a50-7e4c-4f82-ae46-b9e8ec0eb1e0"),
			DisplayName:  utils.StringPtr("Owner"),
			OriginId:     utils.StringPtr(fmt.Sprintf("Owner_%s", *accessPackageResourceRequest.AccessPackageResource.OriginId)),
			OriginSystem: accessPackageResourceRequest.AccessPackageResource.OriginSystem,
			AccessPackageResource: &msgraph.AccessPackageResource{
				ID:           accessPackageResourceRequest.AccessPackageResource.ID,           //Requires poll for ID
				ResourceType: accessPackageResourceRequest.AccessPackageResource.ResourceType, //requires poll for ID
				OriginId:     accessPackageResourceRequest.AccessPackageResource.OriginId,
				//OriginSystem: msgraph.AccessPackageResourceOriginSystemAadGroup,
			},
		},
		AccessPackageResourceScope: &msgraph.AccessPackageResourceScope{
			//ID: "12fec290-4dc6-4e82-88d6-7b5af842e0a3",
			OriginSystem: accessPackageResourceRequest.AccessPackageResource.OriginSystem,
			OriginId:     accessPackageResourceRequest.AccessPackageResource.OriginId,
		},
	})

	//testAccessPackageResourceRequestClient_Delete(t, c, accessPackageResourceRequest)

	// Cleanup
	testAccessPackageResourceRoleScopeAP_Delete(t, c, *accessPackage.ID)
	testAccessPackageResourceRoleScopeResourceRequest_Delete(t, c, accessPackageResourceRequest)
	testAccessPackageResourceRoleScopeCatalog_Delete(t, c, *accessPackageCatalog.ID)
	testAccessPackageResourceRoleScopeGroup_Delete(t, c, aadGroup)
}

// AccessPackageResourceScope

func testAccessPackageResourceRoleScopeClient_Create(t *testing.T, c *test.Test, a msgraph.AccessPackageResourceRoleScope) (accessPackageResourceRoleScope *msgraph.AccessPackageResourceRoleScope) {
	accessPackageResourceRoleScope, status, err := c.AccessPackageResourceRoleScopeClient.Create(c.Context, a)
	if err != nil {
		t.Fatalf("AccessPackageResourceRequestClient.Create(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AccessPackageResourceRequestClient.Create(): invalid status: %d", status)
	}
	if accessPackageResourceRoleScope == nil {
		t.Fatal("AccessPackageResourceRequestClient.Create(): accessPackageResourceRequest was nil")
	}
	if accessPackageResourceRoleScope.ID == nil {
		t.Fatal("AccessPackageResourceRequestClient.Create(): accessPackageResourceRoleScope.ID was nil")
	}
	return
}

func testAccessPackageResourceRoleScopeClient_Get(t *testing.T, c *test.Test, accessPackageId string, id string) (accessPackageResourceRoleScope *msgraph.AccessPackageResourceRoleScope) {
	accessPackageResourceRoleScope, status, err := c.AccessPackageResourceRoleScopeClient.Get(c.Context, accessPackageId, id)
	if err != nil {
		t.Fatalf("AccessPackageResourceRequestClient.Get(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AccessPackageResourceRequestClient.Get(): invalid status: %d", status)
	}
	if accessPackageResourceRoleScope == nil {
		t.Fatal("AccessPackageResourceRequestClient.Get(): policy was nil")
	}
	return
}

func testAccessPackageResourceRoleScopeClient_List(t *testing.T, c *test.Test, accessPackageId string) (accessPackageResourceRoleScope *[]msgraph.AccessPackageResourceRoleScope) {
	accessPackageResourceRoleScopes, _, err := c.AccessPackageResourceRoleScopeClient.List(c.Context, odata.Query{}, accessPackageId)
	if err != nil {
		t.Fatalf("AccessPackageResourceRequestClient.List(): %v", err)
	}
	if accessPackageResourceRoleScopes == nil {
		t.Fatal("AccessPackageResourceRequestClient.List(): accessPackageResourceRequests was nil")
	}
	return
}

func testAccessPackageResourceRoleScopeClient_Delete(t *testing.T, c *test.Test, accessPackageId string, id string) {
	status, err := c.AccessPackageResourceRoleScopeClient.Delete(c.Context, accessPackageId, id)
	if err != nil {
		t.Fatalf("AccessPackageResourceRequestClient.Delete(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AccessPackageResourceRequestClient.Delete(): invalid status: %d", status)
	}
}

// AccessPackageResourceRequest

func testAccessPackageResourceRoleScopeResourceRequest_Create(t *testing.T, c *test.Test, a msgraph.AccessPackageResourceRequest, pollForId bool) (accessPackageResourceRequest *msgraph.AccessPackageResourceRequest) {
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

func testAccessPackageResourceRoleScopeResourceRequest_Delete(t *testing.T, c *test.Test, accessPackageResourceRequest *msgraph.AccessPackageResourceRequest) {
	status, err := c.AccessPackageResourceRequestClient.Delete(c.Context, *accessPackageResourceRequest)
	if err != nil {
		t.Fatalf("AccessPackageResourceRequestClient.Delete(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AccessPackageResourceRequestClient.Delete(): invalid status: %d", status)
	}
}

// AccessPackage

func testAccessPackageResourceRoleScopeAP_Create(t *testing.T, c *test.Test, a msgraph.AccessPackage) (accessPackage *msgraph.AccessPackage) {
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

func testAccessPackageResourceRoleScopeAP_Delete(t *testing.T, c *test.Test, id string) {
	status, err := c.AccessPackageClient.Delete(c.Context, id)
	if err != nil {
		t.Fatalf("AccessPackageClient.Delete(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AccessPackageClient.Delete(): invalid status: %d", status)
	}
}

// AccessPackageCatalog

func testAccessPackageResourceRoleScopeCatalog_Create(t *testing.T, c *test.Test) (accessPackageCatalog *msgraph.AccessPackageCatalog) {
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

func testAccessPackageResourceRoleScopeCatalog_Delete(t *testing.T, c *test.Test, id string) {
	_, err := c.AccessPackageCatalogClient.Delete(c.Context, id)
	if err != nil {
		t.Fatalf("AccessPackageCatalogClient.Delete() - Could not delete test AccessPackage catalog")
	}
}

func testAccessPackageResourceRoleScopeGroup_Create(t *testing.T, c *test.Test, self msgraph.Owners) (group *msgraph.Group) {
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

func testAccessPackageResourceRoleScopeGroup_Delete(t *testing.T, c *test.Test, group *msgraph.Group) {
	_, err := c.GroupsClient.Delete(c.Context, *group.ID())
	if err != nil {
		t.Fatalf("GroupsClient.Delete() - Could not delete test group: %v", err)
	}
}

// AccessPackageResource

func testAccessPackageResourceRoleScopeResource_Get(t *testing.T, c *test.Test, catalogId string, originId string) (accessPackageResource *msgraph.AccessPackageResource) {
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
