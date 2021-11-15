package msgraph_test

import (
	"testing"

	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/msgraph"
	"github.com/manicminer/hamilton/odata"
)

func TestDomainsClient(t *testing.T) {
	c := test.NewTest(t)
	defer c.CancelFunc()

	domains := testDomainsClient_List(t, c)
	testDomainsClient_Get(t, c, *(*domains)[0].ID)
}

func testDomainsClient_List(t *testing.T, c *test.Test) (domains *[]msgraph.Domain) {
	domains, _, err := c.DomainsClient.List(c.Context, odata.Query{})
	if err != nil {
		t.Fatalf("DomainsClient.List(): %v", err)
	}
	if domains == nil {
		t.Fatal("DomainsClient.List(): domains was nil")
	}
	return
}

func testDomainsClient_Get(t *testing.T, c *test.Test, id string) (domain *msgraph.Domain) {
	domain, status, err := c.DomainsClient.Get(c.Context, id, odata.Query{})
	if err != nil {
		t.Fatalf("DomainsClient.Get(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("DomainsClient.Get(): invalid status: %d", status)
	}
	if domain == nil {
		t.Fatal("DomainsClient.Get(): domain was nil")
	}
	return
}
