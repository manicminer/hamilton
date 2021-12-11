package msgraph_test

import (
	"fmt"
	"testing"

	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/msgraph"
	"github.com/manicminer/hamilton/odata"
)

func TestAccessPackageCatalogClient(t *testing.T) {
	c := test.NewTest(t)
	defer c.CancelFunc()

	accessPackageCatalog := testAccessPackageCatalogClient_Create(t, c, msgraph.AccessPackageCatalog{
		DisplayName:         utils.StringPtr(fmt.Sprintf("test-catalog-%s", c.RandomString)),
		CatalogType:         msgraph.AccessPackageCatalogTypeUserManaged,
		State:               msgraph.AccessPackageCatalogStatePublished,
		Description:         utils.StringPtr("Test Access Catalog"),
		IsExternallyVisible: utils.BoolPtr(false),
	})

	testAccessPackageCatalogClient_Update(t, c, msgraph.AccessPackageCatalog{
		ID:          accessPackageCatalog.ID,
		DisplayName: utils.StringPtr(fmt.Sprintf("test-catalog-updated-%s", c.RandomString)),
		Description: utils.StringPtr("Test Access Catalog"),
	})

	testAccessPackageCatalogClient_List(t, c)
	testAccessPackageCatalogClient_Get(t, c, *accessPackageCatalog.ID)
	testAccessPackageCatalogClient_Delete(t, c, *accessPackageCatalog.ID)
}

func testAccessPackageCatalogClient_Create(t *testing.T, c *test.Test, a msgraph.AccessPackageCatalog) (accessPackageCatalog *msgraph.AccessPackageCatalog) {
	accessPackageCatalog, status, err := c.AccessPackageCatalogClient.Create(c.Context, a)
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

func testAccessPackageCatalogClient_Get(t *testing.T, c *test.Test, id string) (accessPackageCatalog *msgraph.AccessPackageCatalog) {
	accessPackageCatalog, status, err := c.AccessPackageCatalogClient.Get(c.Context, id, odata.Query{})
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

func testAccessPackageCatalogClient_Update(t *testing.T, c *test.Test, accessPackageCatalog msgraph.AccessPackageCatalog) {
	status, err := c.AccessPackageCatalogClient.Update(c.Context, accessPackageCatalog)
	if err != nil {
		t.Fatalf("AccessPackageCatalogClient.Update(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AccessPackageCatalogClient.Update(): invalid status: %d", status)
	}
}

func testAccessPackageCatalogClient_List(t *testing.T, c *test.Test) (accessPackageCatalogs *[]msgraph.AccessPackageCatalog) {
	accessPackageCatalogs, _, err := c.AccessPackageCatalogClient.List(c.Context, odata.Query{Top: 10})
	if err != nil {
		t.Fatalf("AccessPackageCatalogClient.List(): %v", err)
	}
	if accessPackageCatalogs == nil {
		t.Fatal("AccessPackageCatalogClient.List(): accessPackageCatalogs was nil")
	}
	return
}

func testAccessPackageCatalogClient_Delete(t *testing.T, c *test.Test, id string) {
	status, err := c.AccessPackageCatalogClient.Delete(c.Context, id)
	if err != nil {
		t.Fatalf("AccessPackageCatalogClient.Delete(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AccessPackageCatalogClient.Delete(): invalid status: %d", status)
	}
}
