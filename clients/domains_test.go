package clients_test

import (
	"testing"

	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/clients"
	"github.com/manicminer/hamilton/clients/internal"
	"github.com/manicminer/hamilton/models"
)

type DomainsClientTest struct {
	connection   *internal.Connection
	client       *clients.DomainsClient
	randomString string
}

func TestDomainsClient(t *testing.T) {
	c := DomainsClientTest{
		connection:   internal.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: internal.RandomString(),
	}
	c.client = clients.NewDomainsClient(c.connection.AuthConfig.TenantID)
	c.client.BaseClient.Authorizer = c.connection.Authorizer

	domains := testDomainsClient_List(t, c)
	testDomainsClient_Get(t, c, *(*domains)[0].ID)
}

func testDomainsClient_List(t *testing.T, c DomainsClientTest) (domains *[]models.Domain) {
	domains, _, err := c.client.List(c.connection.Context)
	if err != nil {
		t.Fatalf("DomainsClient.List(): %v", err)
	}
	if domains == nil {
		t.Fatal("DomainsClient.List(): domains was nil")
	}
	return
}

func testDomainsClient_Get(t *testing.T, c DomainsClientTest, id string) (domain *models.Domain) {
	domain, status, err := c.client.Get(c.connection.Context, id)
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
