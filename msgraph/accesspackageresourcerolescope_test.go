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

type AccessPackageResourceRoleScopeTest struct {
	connection                *test.Connection
	apClient                  *msgraph.AccessPackageClient        //apClient
	apCatalogClient           *msgraph.AccessPackageCatalogClient //Client for Catalog Test to associate as required
	apResourceRequestClient   *msgraph.AccessPackageResourceRequestClient
	apResourceClient          *msgraph.AccessPackageResourceClient
	groupsClient              *msgraph.GroupsClient
	apResourceRoleScopeClient *msgraph.AccessPackageResourceRoleScopeClient
	randomString              string
}

func TestAccessPackageResourceRoleScopeClient(t *testing.T) {
	c := AccessPackageResourceRoleScopeTest{
		connection:   test.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: test.RandomString(),
	}

	// Init clients
	c.apClient = msgraph.NewAccessPackageClient(c.connection.AuthConfig.TenantID)
	c.apClient.BaseClient.Authorizer = c.connection.Authorizer

	c.apCatalogClient = msgraph.NewAccessPackageCatalogClient(c.connection.AuthConfig.TenantID)
	c.apCatalogClient.BaseClient.Authorizer = c.connection.Authorizer

	c.groupsClient = msgraph.NewGroupsClient(c.connection.AuthConfig.TenantID)
	c.groupsClient.BaseClient.Authorizer = c.connection.Authorizer

	c.apResourceRequestClient = msgraph.NewAccessPackageResourceRequestClient(c.connection.AuthConfig.TenantID)
	c.apResourceRequestClient.BaseClient.Authorizer = c.connection.Authorizer

	c.apResourceClient = msgraph.NewAccessPackageResourceClient(c.connection.AuthConfig.TenantID)
	c.apResourceClient.BaseClient.Authorizer = c.connection.Authorizer

	c.apResourceRoleScopeClient = msgraph.NewAccessPackageResourceRoleScopeClient(c.connection.AuthConfig.TenantID)
	c.apResourceRoleScopeClient.BaseClient.Authorizer = c.connection.Authorizer

	// make groups
	token, err := c.connection.Authorizer.Token()
	if err != nil {
		t.Fatalf("could not acquire access token: %v", err)
	}
	claims, err := auth.ParseClaims(token)
	if err != nil {
		t.Fatalf("could not parse claims: %v", err)
	}

	o := DirectoryObjectsClientTest{
		connection:   test.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: test.RandomString(),
	}
	o.client = msgraph.NewDirectoryObjectsClient(c.connection.AuthConfig.TenantID)
	o.client.BaseClient.Authorizer = o.connection.Authorizer

	self := testDirectoryObjectsClient_Get(t, o, claims.ObjectId)

	aadGroup := testapresourcerolescopeGroup_Create(t, c, msgraph.Owners{*self})

	// Create test catalog
	accessPackageCatalog := testapresourcerolescopeCatalog_Create(t, c)

	// Create AP
	accessPackage := testapresourcerolescopeAP_Create(t, c, msgraph.AccessPackage{
		DisplayName:         utils.StringPtr(fmt.Sprintf("test-accesspackage-%s", c.randomString)),
		CatalogId:           accessPackageCatalog.ID,
		Description:         utils.StringPtr("Test Access Package"),
		IsHidden:            utils.BoolPtr(false),
		IsRoleScopesVisible: utils.BoolPtr(false),
	})

	// Create Resource Request AND Poll for ID

	pollForId := true

	accessPackageResourceRequest := testapresourcerolescopeResourceRequest_Create(t, c, msgraph.AccessPackageResourceRequest{
		CatalogId:   accessPackage.CatalogId,
		RequestType: utils.StringPtr("AdminAdd"),
		AccessPackageResource: &msgraph.AccessPackageResource{
			OriginId:     aadGroup.ID,
			OriginSystem: utils.StringPtr("AadGroup"),
			//ResourceType: utils.StringPtr("Security Group") - This is not mandatory for groups but is seen in sharepoint emails - Probably make mandatory downstream for sanity? Return on create method if we poll the resource regardless
		},
	}, pollForId)

	// Imagine can roll this into the AP Assignment as a block Needing the resource request

	accessPackageResourceRoleScope := testAccessPackageResourceRoleScopeClient_Create(t, c, msgraph.AccessPackageResourceRoleScope{
		AccessPackageId: accessPackage.ID,
		AccessPackageResourceRole: &msgraph.AccessPackageResourceRole{
			//ID: utils.StringPtr("405b1a50-7e4c-4f82-ae46-b9e8ec0eb1e0"),
			DisplayName:  utils.StringPtr("Member"),
			OriginId:     utils.StringPtr(fmt.Sprintf("Member_%s", *accessPackageResourceRequest.AccessPackageResource.OriginId)),
			OriginSystem: accessPackageResourceRequest.AccessPackageResource.OriginSystem,
			AccessPackageResource: &msgraph.AccessPackageResource{
				ID:           accessPackageResourceRequest.AccessPackageResource.ID,           //Requires poll for id
				ResourceType: accessPackageResourceRequest.AccessPackageResource.ResourceType, //requires poll for id
				OriginId:     accessPackageResourceRequest.AccessPackageResource.OriginId,
				//OriginSystem: utils.StringPtr("AadGroup"),
			},
		},
		AccessPackageResourceScope: &msgraph.AccessPackageResourceScope{
			//ID: "12fec290-4dc6-4e82-88d6-7b5af842e0a3",
			OriginSystem: accessPackageResourceRequest.AccessPackageResource.OriginSystem,
			OriginId:     accessPackageResourceRequest.AccessPackageResource.OriginId,
		},
	})
	//Get tests
	testAccessPackageResourceRoleScopeClient_Get(t, c, *accessPackageResourceRoleScope.AccessPackageId, *accessPackageResourceRoleScope.ID)
	testapresourcerolescopeResource_Get(t, c, *accessPackageResourceRequest.CatalogId, *accessPackageResourceRequest.AccessPackageResource.OriginId)
	// List tests
	testAccessPackageResourceRoleScopeClient_List(t, c, *accessPackageResourceRoleScope.AccessPackageId)

	// Force replacement scenario

	testapresourcerolescopeAP_Delete(t, c, *accessPackage.ID)
	accessPackage = testapresourcerolescopeAP_Create(t, c, *accessPackage) //New ID

	testAccessPackageResourceRoleScopeClient_Create(t, c, msgraph.AccessPackageResourceRoleScope{
		AccessPackageId: accessPackage.ID,
		AccessPackageResourceRole: &msgraph.AccessPackageResourceRole{
			//ID: utils.StringPtr("405b1a50-7e4c-4f82-ae46-b9e8ec0eb1e0"),
			DisplayName:  utils.StringPtr("Owner"),
			OriginId:     utils.StringPtr(fmt.Sprintf("Owner_%s", *accessPackageResourceRequest.AccessPackageResource.OriginId)),
			OriginSystem: accessPackageResourceRequest.AccessPackageResource.OriginSystem,
			AccessPackageResource: &msgraph.AccessPackageResource{
				ID:           accessPackageResourceRequest.AccessPackageResource.ID,           //Requires poll for id
				ResourceType: accessPackageResourceRequest.AccessPackageResource.ResourceType, //requires poll for id
				OriginId:     accessPackageResourceRequest.AccessPackageResource.OriginId,
				//OriginSystem: utils.StringPtr("AadGroup"),
			},
		},
		AccessPackageResourceScope: &msgraph.AccessPackageResourceScope{
			//ID: "12fec290-4dc6-4e82-88d6-7b5af842e0a3",
			OriginSystem: accessPackageResourceRequest.AccessPackageResource.OriginSystem,
			OriginId:     accessPackageResourceRequest.AccessPackageResource.OriginId,
		},
	})

	//testAccessPackageResourceRoleScopeClient_Delete(t, c, *accessPackageResourceRoleScope)

	//testAccessPackageResourceRequestClient_Delete(t, c, accessPackageResourceRequest)
	// Cleanup
	testapresourcerolescopeAP_Delete(t, c, *accessPackage.ID)
	testapresourcerolescopeResourceRequest_Delete(t, c, accessPackageResourceRequest)
	testapresourcerolescopeCatalog_Delete(t, c, *accessPackageCatalog.ID)
	testapresourcerolescopeGroup_Delete(t, c, aadGroup)

}

// AP Resource Scope

func testAccessPackageResourceRoleScopeClient_Create(t *testing.T, c AccessPackageResourceRoleScopeTest, a msgraph.AccessPackageResourceRoleScope) (accessPackageResourceRoleScope *msgraph.AccessPackageResourceRoleScope) {
	accessPackageResourceRoleScope, status, err := c.apResourceRoleScopeClient.Create(c.connection.Context, a)
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

func testAccessPackageResourceRoleScopeClient_Get(t *testing.T, c AccessPackageResourceRoleScopeTest, accessPackageId string, id string) (accessPackageResourceRoleScope *msgraph.AccessPackageResourceRoleScope) {
	accessPackageResourceRoleScope, status, err := c.apResourceRoleScopeClient.Get(c.connection.Context, accessPackageId, id)
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

func testAccessPackageResourceRoleScopeClient_List(t *testing.T, c AccessPackageResourceRoleScopeTest, accessPackageId string) (accessPackageResourceRoleScope *[]msgraph.AccessPackageResourceRoleScope) {
	accessPackageResourceRoleScopes, _, err := c.apResourceRoleScopeClient.List(c.connection.Context, odata.Query{}, accessPackageId)
	if err != nil {
		t.Fatalf("AccessPackageResourceRequestClient.List(): %v", err)
	}
	if accessPackageResourceRoleScopes == nil {
		t.Fatal("AccessPackageResourceRequestClient.List(): accessPackageResourceRequests was nil")
	}
	return
}

// AP Resource

func testapresourcerolescopeResourceRequest_Create(t *testing.T, c AccessPackageResourceRoleScopeTest, a msgraph.AccessPackageResourceRequest, pollForId bool) (accessPackageResourceRequest *msgraph.AccessPackageResourceRequest) {
	accessPackageResourceRequest, status, err := c.apResourceRequestClient.Create(c.connection.Context, a, pollForId)
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

func testapresourcerolescopeResourceRequest_Delete(t *testing.T, c AccessPackageResourceRoleScopeTest, accessPackageResourceRequest *msgraph.AccessPackageResourceRequest) {
	status, err := c.apResourceRequestClient.Delete(c.connection.Context, *accessPackageResourceRequest)
	if err != nil {
		t.Fatalf("AccessPackageResourceRequestClient.Delete(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AccessPackageResourceRequestClient.Delete(): invalid status: %d", status)
	}
}

// AP

func testapresourcerolescopeAP_Create(t *testing.T, c AccessPackageResourceRoleScopeTest, a msgraph.AccessPackage) (accessPackage *msgraph.AccessPackage) {
	accessPackage, status, err := c.apClient.Create(c.connection.Context, a)
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

func testapresourcerolescopeAP_Delete(t *testing.T, c AccessPackageResourceRoleScopeTest, id string) {
	status, err := c.apClient.Delete(c.connection.Context, id)
	if err != nil {
		t.Fatalf("AccessPackageClient.Delete(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AccessPackageClient.Delete(): invalid status: %d", status)
	}
}

// AP Catalog

func testapresourcerolescopeCatalog_Create(t *testing.T, c AccessPackageResourceRoleScopeTest) (accessPackageCatalog *msgraph.AccessPackageCatalog) {
	accessPackageCatalog, _, err := c.apCatalogClient.Create(c.connection.Context, msgraph.AccessPackageCatalog{
		DisplayName:         utils.StringPtr(fmt.Sprintf("test-catalog-%s", c.randomString)),
		CatalogType:         utils.StringPtr("UserManaged"),
		CatalogStatus:       utils.StringPtr("Published"),
		Description:         utils.StringPtr("Test Access Catalog"),
		IsExternallyVisible: utils.BoolPtr(false),
	})

	if err != nil {
		t.Fatalf("AccessPackageCatalogClient.Create() - Could not create test AccessPackage catalog: %v", err)
	}
	return
}

func testapresourcerolescopeCatalog_Delete(t *testing.T, c AccessPackageResourceRoleScopeTest, id string) {
	_, err := c.apCatalogClient.Delete(c.connection.Context, id)
	if err != nil {
		t.Fatalf("AccessPackageCatalogClient.Delete() - Could not delete test AccessPackage catalog")
	}
}

func testapresourcerolescopeGroup_Create(t *testing.T, c AccessPackageResourceRoleScopeTest, self msgraph.Owners) (group *msgraph.Group) {
	group, _, err := c.groupsClient.Create(c.connection.Context, msgraph.Group{
		DisplayName:     utils.StringPtr(fmt.Sprintf("%s-%s", "testapresourcerequest", c.randomString)),
		MailEnabled:     utils.BoolPtr(false),
		MailNickname:    utils.StringPtr(fmt.Sprintf("%s-%s", "testapresourcerequest", c.randomString)),
		SecurityEnabled: utils.BoolPtr(true),
		Owners:          &self,
	})

	if err != nil {
		t.Fatalf("GroupsClient.Create() - Could not create test group: %v", err)
	}
	return
}

func testapresourcerolescopeGroup_Delete(t *testing.T, c AccessPackageResourceRoleScopeTest, group *msgraph.Group) {
	_, err := c.groupsClient.Delete(c.connection.Context, *group.ID)
	if err != nil {
		t.Fatalf("GroupsClient.Delete() - Could not delete test group: %v", err)
	}
}

//APResource

func testapresourcerolescopeResource_Get(t *testing.T, c AccessPackageResourceRoleScopeTest, catalogId string, originId string) (accessPackageResource *msgraph.AccessPackageResource) {
	accessPackageResource, status, err := c.apResourceClient.Get(c.connection.Context, catalogId, originId)

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
