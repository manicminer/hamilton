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

type AccessPackageResourceRequestTest struct {
	connection      *test.Connection
	apClient        *msgraph.AccessPackageClient        //apClient
	apCatalogClient *msgraph.AccessPackageCatalogClient //Client for Catalog Test to associate as required
	apResourceRequestClient *msgraph.AccessPackageResourceRequestClient
	apResourceClient *msgraph.AccessPackageResourceClient
	groupsClient *msgraph.GroupsClient
	randomString    string
}

func TestAccessPackageResourceRequestClient(t *testing.T) {
	c := AccessPackageResourceRequestTest{
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

	aadGroup := testapresourcerequestGroup_Create(t, c, msgraph.Owners{*self})

	// Create test catalog
	accessPackageCatalog := testapresourcerequestCatalog_Create(t, c)

	// Create AP
	accessPackage := testapresourcerequestAP_Create(t, c, msgraph.AccessPackage{
		DisplayName:         utils.StringPtr(fmt.Sprintf("test-accesspackage-%s", c.randomString)),
		CatalogId:           accessPackageCatalog.ID,
		Description:         utils.StringPtr("Test Access Package"),
		IsHidden:            utils.BoolPtr(false),
		IsRoleScopesVisible: utils.BoolPtr(false),
	})

	// Create Resource Request AND Poll for ID

	pollForId := true

	accessPackageResourceRequest := testAccessPackageResourceRequestClient_Create(t, c, msgraph.AccessPackageResourceRequest{
		CatalogId: accessPackage.CatalogId,
		RequestType: utils.StringPtr("AdminAdd"),
		AccessPackageResource: &msgraph.AccessPackageResource{
			OriginId: aadGroup.ID,
			OriginSystem: utils.StringPtr("AadGroup"),
		},
	}, pollForId)


	// Figure out what the ResourceId is IF We haven't Polled for it below
	
	//accessPackageResource := testapresourcerequestResource_Get(t, c, *accessPackageResourceRequest.CatalogId, *accessPackageResourceRequest.AccessPackageResource.OriginId)
	//accessPackageResourceRequest.AccessPackageResource.ID = accessPackageResource.ID

	// There is no update tests here

	// Tests

	//Resource Client
	testapresourcerequestResource_Get(t, c, *accessPackageResourceRequest.CatalogId, *accessPackageResourceRequest.AccessPackageResource.OriginId) //Resource we created from Request Exists
	testapresourcerequestResource_List(t, c, *accessPackageResourceRequest.CatalogId)
	// Requests Client
	testAccessPackageResourceRequestClient_List(t, c)
	testAccessPackageResourceRequestClient_Get(t, c, *accessPackageResourceRequest.ID) //Req Exists
	
	testAccessPackageResourceRequestClient_Delete(t, c, accessPackageResourceRequest)
	// Cleanup
	testapresourcerequestAP_Delete(t, c, *accessPackage.ID)
	testapresourcerequestCatalog_Delete(t, c, *accessPackageCatalog.ID)
	testapresourcerequestGroup_Delete(t, c, aadGroup)

}

// AP Resource Request

func testAccessPackageResourceRequestClient_Create(t *testing.T, c AccessPackageResourceRequestTest, a msgraph.AccessPackageResourceRequest, pollForId bool) (accessPackageResourceRequest *msgraph.AccessPackageResourceRequest) {
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

func testAccessPackageResourceRequestClient_Get(t *testing.T, c AccessPackageResourceRequestTest, id string) (accessPackageResourceRequest *msgraph.AccessPackageResourceRequest) {
	accessPackageResourceRequest, status, err := c.apResourceRequestClient.Get(c.connection.Context, id,)
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

func testAccessPackageResourceRequestClient_List(t *testing.T, c AccessPackageResourceRequestTest) (accessPackageResourceRequests *[]msgraph.AccessPackageResourceRequest) {
	accessPackageResourceRequests, _, err := c.apResourceRequestClient.List(c.connection.Context, odata.Query{Top: 10})
	if err != nil {
		t.Fatalf("AccessPackageResourceRequestClient.List(): %v", err)
	}
	if accessPackageResourceRequests == nil {
		t.Fatal("AccessPackageResourceRequestClient.List(): accessPackageResourceRequests was nil")
	}
	return
}

func testAccessPackageResourceRequestClient_Delete(t *testing.T, c AccessPackageResourceRequestTest, accessPackageResourceRequest *msgraph.AccessPackageResourceRequest) {
	status, err := c.apResourceRequestClient.Delete(c.connection.Context, *accessPackageResourceRequest)
	if err != nil {
		t.Fatalf("AccessPackageResourceRequestClient.Delete(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AccessPackageResourceRequestClient.Delete(): invalid status: %d", status)
	}
}



// AP

func testapresourcerequestAP_Create(t *testing.T, c AccessPackageResourceRequestTest, a msgraph.AccessPackage) (accessPackage *msgraph.AccessPackage) {
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

func testapresourcerequestAP_Delete(t *testing.T, c AccessPackageResourceRequestTest, id string) {
	status, err := c.apClient.Delete(c.connection.Context, id)
	if err != nil {
		t.Fatalf("AccessPackageClient.Delete(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AccessPackageClient.Delete(): invalid status: %d", status)
	}
}

// AP Catalog

func testapresourcerequestCatalog_Create(t *testing.T, c AccessPackageResourceRequestTest) (accessPackageCatalog *msgraph.AccessPackageCatalog) {
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

func testapresourcerequestCatalog_Delete(t *testing.T, c AccessPackageResourceRequestTest, id string) {
	_, err := c.apCatalogClient.Delete(c.connection.Context, id)
	if err != nil {
		t.Fatalf("AccessPackageCatalogClient.Delete() - Could not delete test AccessPackage catalog")
	}
}

func testapresourcerequestGroup_Create(t *testing.T, c AccessPackageResourceRequestTest, self msgraph.Owners) (group *msgraph.Group) {
	group, _, err := c.groupsClient.Create(c.connection.Context, msgraph.Group{
		DisplayName:     utils.StringPtr(fmt.Sprintf("%s-%s", "testapresourcerequest", c.randomString)),
		MailEnabled:     utils.BoolPtr(false),
		MailNickname:    utils.StringPtr(fmt.Sprintf("%s-%s", "testapresourcerequest", c.randomString)),
		SecurityEnabled: utils.BoolPtr(true),
		Owners: &self,
	})

	if err != nil {
		t.Fatalf("GroupsClient.Create() - Could not create test group: %v", err)
	}
	return
}

func testapresourcerequestGroup_Delete(t *testing.T, c AccessPackageResourceRequestTest, group *msgraph.Group) {
	_, err := c.groupsClient.Delete(c.connection.Context, *group.ID)
	if err != nil {
		t.Fatalf("GroupsClient.Delete() - Could not delete test group: %v", err)
	}
}



//APResource


func testapresourcerequestResource_Get(t *testing.T, c AccessPackageResourceRequestTest, catalogId string, originId string) (accessPackageResource *msgraph.AccessPackageResource) {
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

func testapresourcerequestResource_List(t *testing.T, c AccessPackageResourceRequestTest, catalogId string) (accessPackageResources *[]msgraph.AccessPackageResource) {
	accessPackageResources, status, err := c.apResourceClient.List(c.connection.Context, odata.Query{Top: 10}, catalogId)

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