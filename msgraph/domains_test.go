package msgraph_test

import (
	"testing"

	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/msgraph"
	"github.com/manicminer/hamilton/odata"
)

type DomainsClientTest struct {
	connection   *test.Connection
	client       *msgraph.DomainsClient
	randomString string
}

func TestDomainsClient(t *testing.T) {
	c := DomainsClientTest{
		connection:   test.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: test.RandomString(),
	}
	c.client = msgraph.NewDomainsClient(c.connection.AuthConfig.TenantID)
	c.client.BaseClient.Authorizer = c.connection.Authorizer
	c.client.BaseClient.Endpoint = c.connection.AuthConfig.Environment.MsGraph.Endpoint

	domains := testDomainsClient_List(t, c)
	testDomainsClient_Get(t, c, *(*domains)[0].ID)
}

func testDomainsClient_List(t *testing.T, c DomainsClientTest) (domains *[]msgraph.Domain) {
	domains, _, err := c.client.List(c.connection.Context, odata.Query{})
	if err != nil {
		t.Fatalf("DomainsClient.List(): %v", err)
	}
	if domains == nil {
		t.Fatal("DomainsClient.List(): domains was nil")
	}
	return
}

func testDomainsClient_Get(t *testing.T, c DomainsClientTest, id string) (domain *msgraph.Domain) {
	domain, status, err := c.client.Get(c.connection.Context, id, odata.Query{})
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
