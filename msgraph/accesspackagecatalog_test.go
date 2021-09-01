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

type AccessPackageCatalogTest struct {
	connection   *test.Connection
	apCatalogClient *msgraph.AccessPackageCatalogClient
	randomString string
}

func TestAccessPackageCatalogClient(t *testing.T) {
	c := AccessPackageCatalogTest{
		connection:   test.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: test.RandomString(),
	}

	c.apCatalogClient = msgraph.NewAccessPackageCatalogClient(c.connection.AuthConfig.TenantID)
	c.apCatalogClient.BaseClient.Authorizer = c.connection.Authorizer

	// act
	accessPackageCatalog := testAccessPackageCatalogClient_Create(t, c, msgraph.AccessPackageCatalog{
		DisplayName: utils.StringPtr(fmt.Sprintf("test-catalog-%s", c.randomString)),
		CatalogType:       utils.StringPtr("UserManaged"),
		CatalogStatus:     utils.StringPtr("Published"),
		Description: utils.StringPtr("Test Access Catalog"),
		IsExternallyVisible: utils.BoolPtr(false),
		
	})

	updateAccessPackageCatalog := msgraph.AccessPackageCatalog{
		ID:          accessPackageCatalog.ID,
		DisplayName: utils.StringPtr(fmt.Sprintf("test-catalog-updated-%s", c.randomString)),
		Description: utils.StringPtr("Test Access Catalog"),
	}
	testAccessPackageCatalogClient_Update(t, c, updateAccessPackageCatalog)

	testAccessPackageCatalogClient_List(t, c)
	testAccessPackageCatalogClient_Get(t, c, *accessPackageCatalog.ID)
    testAccessPackageCatalogClient_Delete(t, c, *accessPackageCatalog.ID)

}

func testAccessPackageCatalogClient_Create(t *testing.T, c AccessPackageCatalogTest, a msgraph.AccessPackageCatalog) (accessPackageCatalog *msgraph.AccessPackageCatalog) {
	accessPackageCatalog, status, err := c.apCatalogClient.Create(c.connection.Context, a)
	if err != nil {
		t.Fatalf("AccessPackageCatalogClient.Create(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AccessPackageCatalogClient.Create(): invalid status: %d", status)
	}
	if accessPackageCatalog == nil {
		t.Fatal("AccessPackageCatalogClient.Create(): accessPackageCatalog was nil")
	}
	if accessPackageCatalog.ID == nil {
		t.Fatal("AccessPackageCatalogClient.Create(): acccessPackageCatalog.ID was nil")
	}
	return
}

func testAccessPackageCatalogClient_Get(t *testing.T, c AccessPackageCatalogTest, id string) (accessPackageCatalog *msgraph.AccessPackageCatalog) {
	accessPackageCatalog, status, err := c.apCatalogClient.Get(c.connection.Context, id, odata.Query{})
	if err != nil {
		t.Fatalf("AccessPackageCatalogClient.Get(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AccessPackageCatalogClient.Get(): invalid status: %d", status)
	}
	if accessPackageCatalog == nil {
		t.Fatal("AccessPackageCatalogClient.Get(): policy was nil")
	}
	return
}

func testAccessPackageCatalogClient_Update(t *testing.T, c AccessPackageCatalogTest, accessPackageCatalog msgraph.AccessPackageCatalog) {
	status, err := c.apCatalogClient.Update(c.connection.Context, accessPackageCatalog)
	if err != nil {
		t.Fatalf("AccessPackageCatalogClient.Update(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AccessPackageCatalogClient.Update(): invalid status: %d", status)
	}
}

func testAccessPackageCatalogClient_List(t *testing.T, c AccessPackageCatalogTest) (accessPackageCatalogs *[]msgraph.AccessPackageCatalog) {
	accessPackageCatalogs, _, err := c.apCatalogClient.List(c.connection.Context, odata.Query{Top: 10})
	if err != nil {
		t.Fatalf("AccessPackageCatalogClient.List(): %v", err)
	}
	if accessPackageCatalogs == nil {
		t.Fatal("AccessPackageCatalogClient.List(): accessPackageCatalogs was nil")
	}
	return
}

func testAccessPackageCatalogClient_Delete(t *testing.T, c AccessPackageCatalogTest, id string) {
	status, err := c.apCatalogClient.Delete(c.connection.Context, id)
	if err != nil {
		t.Fatalf("AccessPackageCatalogClient.Delete(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AccessPackageCatalogClient.Delete(): invalid status: %d", status)
	}
}