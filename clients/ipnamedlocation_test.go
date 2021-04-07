package clients_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/clients"
	"github.com/manicminer/hamilton/clients/internal"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/models"
)

type IPNamedLocationClientTest struct {
	connection   *internal.Connection
	client       *clients.IPNamedLocationClient
	randomString string
}

func TestIPNamedLocationClient(t *testing.T) {
	rs := internal.RandomString()
	c := IPNamedLocationClientTest{
		connection:   internal.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: rs,
	}
	c.client = clients.NewIPNamedLocationClient(c.connection.AuthConfig.TenantID)
	c.client.BaseClient.Authorizer = c.connection.Authorizer

	newIPNamedLocation := models.IPNamedLocation{
		NamedLocation: &models.NamedLocation{
			ODataType:   utils.StringPtr("#microsoft.graph.ipNamedLocation"),
			DisplayName: utils.StringPtr("Test IP Named Location")},
		IPRanges: &[]models.IPNamedLocationIPRange{
			{
				CIDRAddress: utils.StringPtr("8.8.8.8/32"),
			},
			{
				CIDRAddress: utils.StringPtr("2001:4860:4860::8888/128"),
			},
		},
		IsTrusted: utils.BoolPtr(true),
	}

	ipNamedLocation := testIPNamedLocationClient_Create(t, c, newIPNamedLocation)
	// Running get too quickly after create often results in the resource not being found
	time.Sleep(5 * time.Second)
	testIPNamedLocationClient_Get(t, c, *ipNamedLocation.ID)

	ipNamedLocation.DisplayName = utils.StringPtr(fmt.Sprintf("test-updated-ipnl-%s", c.randomString))
	testIPNamedLocationClient_Update(t, c, *ipNamedLocation)

	testIPNamedLocationClient_List(t, c)
	testIPNamedLocationClient_Delete(t, c, *ipNamedLocation.ID)
}

func testIPNamedLocationClient_Create(t *testing.T, c IPNamedLocationClientTest, ipnl models.IPNamedLocation) (ipNamedLocation *models.IPNamedLocation) {
	ipNamedLocation, status, err := c.client.Create(c.connection.Context, ipnl)
	if err != nil {
		t.Fatalf("IPNamedLocationClient.Create(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("IPNamedLocationClient.Create(): invalid status: %d", status)
	}
	if ipNamedLocation == nil {
		t.Fatal("IPNamedLocationClient.Create(): ipNamedLocation was nil")
	}
	if ipNamedLocation.ID == nil {
		t.Fatal("IPNamedLocationClient.Create(): ipNamedLocation.ID was nil")
	}
	return
}

func testIPNamedLocationClient_Get(t *testing.T, c IPNamedLocationClientTest, id string) (ipNamedLocation *models.IPNamedLocation) {
	ipNamedLocation, status, err := c.client.Get(c.connection.Context, id)
	if err != nil {
		t.Fatalf("IPNamedLocationClient.Get(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("IPNamedLocationClient.Get(): invalid status: %d", status)
	}
	if ipNamedLocation == nil {
		t.Fatal("IPNamedLocationClient.Get(): ipNamedLocation was nil")
	}
	return
}

func testIPNamedLocationClient_Update(t *testing.T, c IPNamedLocationClientTest, ipnl models.IPNamedLocation) {
	status, err := c.client.Update(c.connection.Context, ipnl)
	if err != nil {
		t.Fatalf("IPNamedLocationClient.Update(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("IPNamedLocationClient.Update(): invalid status: %d", status)
	}
}

func testIPNamedLocationClient_List(t *testing.T, c IPNamedLocationClientTest) (ipNamedLocations *[]models.IPNamedLocation) {
	ipNamedLocations, _, err := c.client.List(c.connection.Context)
	if err != nil {
		t.Fatalf("IPNamedLocationClient.List(): %v", err)
	}
	if ipNamedLocations == nil {
		t.Fatal("IPNamedLocationClient.List(): ipNamedLocations was nil")
	}
	return
}

func testIPNamedLocationClient_Delete(t *testing.T, c IPNamedLocationClientTest, id string) {
	status, err := c.client.Delete(c.connection.Context, id)
	if err != nil {
		t.Fatalf("IPNamedLocationClient.Delete(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("IPNamedLocationClient.Delete(): invalid status: %d", status)
	}
}
