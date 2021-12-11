package msgraph_test

import (
	"fmt"
	"testing"

	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/msgraph"
	"github.com/manicminer/hamilton/odata"
)

func TestAccessPackageClient(t *testing.T) {
	c := test.NewTest(t)
	defer c.CancelFunc()

	// Create test catalog
	accessPackageCatalog := testapCatalog_Create(t, c)

	// Create AP
	accessPackage := testAccessPackageClient_Create(t, c, msgraph.AccessPackage{
		DisplayName:         utils.StringPtr(fmt.Sprintf("test-accesspackage-%s", c.RandomString)),
		Catalog:           &msgraph.AccessPackageCatalog{
			ID: accessPackageCatalog.ID,
		},
		Description:         utils.StringPtr("Test Access Package"),
		IsHidden:            utils.BoolPtr(false),
	})

	// Update test
	testAccessPackageClient_Update(t, c, msgraph.AccessPackage{
		ID:          accessPackage.ID,
		DisplayName: utils.StringPtr(fmt.Sprintf("test-accesspackage-updated-%s", c.RandomString)),
	})

	// Other operations
	testAccessPackageClient_List(t, c)
	testAccessPackageClient_Get(t, c, *accessPackage.ID)
	testAccessPackageClient_Delete(t, c, *accessPackage.ID)

	// Cleanup
	testapCatalog_Delete(t, c, accessPackageCatalog)
}

// AP

func testAccessPackageClient_Create(t *testing.T, c *test.Test, a msgraph.AccessPackage) (accessPackage *msgraph.AccessPackage) {
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

func testAccessPackageClient_Get(t *testing.T, c *test.Test, id string) (accessPackage *msgraph.AccessPackage) {
	accessPackage, status, err := c.AccessPackageClient.Get(c.Context, id, odata.Query{})
	if err != nil {
		t.Fatalf("AccessPackageClient.Get(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AccessPackageClient.Get(): invalid status: %d", status)
	}
	if accessPackage == nil {
		t.Fatal("AccessPackageClient.Get(): policy was nil")
	}
	return
}

func testAccessPackageClient_Update(t *testing.T, c *test.Test, accessPackage msgraph.AccessPackage) {
	status, err := c.AccessPackageClient.Update(c.Context, accessPackage)
	if err != nil {
		t.Fatalf("AccessPackageClient.Update(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AccessPackageClient.Update(): invalid status: %d", status)
	}
}

func testAccessPackageClient_List(t *testing.T, c *test.Test) (accessPackages *[]msgraph.AccessPackage) {
	accessPackages, _, err := c.AccessPackageClient.List(c.Context, odata.Query{Top: 10})
	if err != nil {
		t.Fatalf("AccessPackageClient.List(): %v", err)
	}
	if accessPackages == nil {
		t.Fatal("AccessPackageClient.List(): accessPackages was nil")
	}
	return
}

func testAccessPackageClient_Delete(t *testing.T, c *test.Test, id string) {
	status, err := c.AccessPackageClient.Delete(c.Context, id)
	if err != nil {
		t.Fatalf("AccessPackageClient.Delete(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AccessPackageClient.Delete(): invalid status: %d", status)
	}
}

// AP Catalog

func testapCatalog_Create(t *testing.T, c *test.Test) (accessPackageCatalog *msgraph.AccessPackageCatalog) {
	accessPackageCatalog, _, err := c.AccessPackageCatalogClient.Create(c.Context, msgraph.AccessPackageCatalog{
		DisplayName:         utils.StringPtr(fmt.Sprintf("test-catalog-%s", c.RandomString)),
		CatalogType:         msgraph.AccessPackageCatalogTypeUserManaged,
		State:       msgraph.AccessPackageCatalogStatePublished,
		Description:         utils.StringPtr("Test Access Catalog"),
		IsExternallyVisible: utils.BoolPtr(false),
	})

	if err != nil {
		t.Fatalf("AccessPackageCatalogClient.Create() - Could not create test AccessPackage catalog: %v", err)
	}
	return
}

func testapCatalog_Delete(t *testing.T, c *test.Test, accessPackageCatalog *msgraph.AccessPackageCatalog) {
	_, err := c.AccessPackageCatalogClient.Delete(c.Context, *accessPackageCatalog.ID)
	if err != nil {
		t.Fatalf("AccessPackageCatalogClient.Delete() - Could not delete test AccessPackage catalog")
	}
}
