package msgraph_test

import (
	"fmt"
	"testing"

	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/msgraph"
)

type NamedLocationClientTest struct {
	connection   *test.Connection
	client       *msgraph.NamedLocationClient
	randomString string
}

func TestNamedLocationClient(t *testing.T) {
	rs := test.RandomString()
	c := NamedLocationClientTest{
		connection:   test.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: rs,
	}
	c.client = msgraph.NewNamedLocationClient(c.connection.AuthConfig.TenantID)
	c.client.BaseClient.Authorizer = c.connection.Authorizer

	newIPNamedLocation := msgraph.IPNamedLocation{
		BaseNamedLocation: &msgraph.BaseNamedLocation{
			DisplayName: utils.StringPtr("Test IP Named Location"),
		},
		IPRanges: &[]msgraph.IPNamedLocationIPRange{
			{
				CIDRAddress: utils.StringPtr("8.8.8.8/32"),
			},
			{
				CIDRAddress: utils.StringPtr("2001:4860:4860::8888/128"),
			},
		},
		IsTrusted: utils.BoolPtr(true),
	}

	newCountryNamedLocation := msgraph.CountryNamedLocation{
		BaseNamedLocation: &msgraph.BaseNamedLocation{
			DisplayName: utils.StringPtr("Test Country Named Location"),
		},
		CountriesAndRegions: &[]string{"US", "GB"},
	}

	ipNamedLocation := testNamedLocationClient_CreateIP(t, c, newIPNamedLocation)
	countryNamedLocation := testNamedLocationClient_CreateCountry(t, c, newCountryNamedLocation)

	ipNamedLocation.DisplayName = utils.StringPtr(fmt.Sprintf("test-updated-ipnl-%s", c.randomString))
	ipNamedLocation.CreatedDateTime = nil
	ipNamedLocation.ModifiedDateTime = nil
	countryNamedLocation.DisplayName = utils.StringPtr(fmt.Sprintf("test-updated-cnl-%s", c.randomString))
	countryNamedLocation.CreatedDateTime = nil
	countryNamedLocation.ModifiedDateTime = nil

	testNamedLocationClient_UpdateIP(t, c, *ipNamedLocation)
	testNamedLocationClient_UpdateCountry(t, c, *countryNamedLocation)

	testNamedLocationClient_List(t, c, "")
	// Running get after the update to give the API a chance to catch up
	testNamedLocationClient_GetIP(t, c, *ipNamedLocation.ID)
	testNamedLocationClient_GetCountry(t, c, *countryNamedLocation.ID)

	testNamedLocationClient_Delete(t, c, *ipNamedLocation.ID)
	testNamedLocationClient_Delete(t, c, *countryNamedLocation.ID)
}

func testNamedLocationClient_CreateIP(t *testing.T, c NamedLocationClientTest, ipnl msgraph.IPNamedLocation) (ipNamedLocation *msgraph.IPNamedLocation) {
	ipNamedLocation, status, err := c.client.CreateIP(c.connection.Context, ipnl)
	if err != nil {
		t.Fatalf("NamedLocationClient.CreateIP(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("NamedLocationClient.CreateIP(): invalid status: %d", status)
	}
	if ipNamedLocation == nil {
		t.Fatal("NamedLocationClient.CreateIP(): ipNamedLocation was nil")
	}
	if ipNamedLocation.ID == nil {
		t.Fatal("NamedLocationClient.CreateIP(): ipNamedLocation.ID was nil")
	}
	return
}

func testNamedLocationClient_CreateCountry(t *testing.T, c NamedLocationClientTest, cnl msgraph.CountryNamedLocation) (countryNamedLocation *msgraph.CountryNamedLocation) {
	countryNamedLocation, status, err := c.client.CreateCountry(c.connection.Context, cnl)
	if err != nil {
		t.Fatalf("NamedLocationClient.CreateCountry(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("NamedLocationClient.CreateCountry(): invalid status: %d", status)
	}
	if countryNamedLocation == nil {
		t.Fatal("NamedLocationClient.CreateCountry(): countryNamedLocation was nil")
	}
	if countryNamedLocation.ID == nil {
		t.Fatal("NamedLocationClient.CreateCountry(): countryNamedLocation.ID was nil")
	}
	return
}

func testNamedLocationClient_GetIP(t *testing.T, c NamedLocationClientTest, id string) (ipNamedLocation *msgraph.IPNamedLocation) {
	ipNamedLocation, status, err := c.client.GetIP(c.connection.Context, id)
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

func testNamedLocationClient_GetCountry(t *testing.T, c NamedLocationClientTest, id string) (countryNamedLocation *msgraph.CountryNamedLocation) {
	countryNamedLocation, status, err := c.client.GetCountry(c.connection.Context, id)
	if err != nil {
		t.Fatalf("NamedLocationClient.GetCountry(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("NamedLocationClient.GetCountry(): invalid status: %d", status)
	}
	if countryNamedLocation == nil {
		t.Fatal("NamedLocationClient.GetCountry(): countryNamedLocation was nil")
	}
	return
}

func testNamedLocationClient_UpdateIP(t *testing.T, c NamedLocationClientTest, ipnl msgraph.IPNamedLocation) {
	status, err := c.client.UpdateIP(c.connection.Context, ipnl)
	if err != nil {
		t.Fatalf("NamedLocationClient.UpdateIP(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("NamedLocationClient.UpdateIP(): invalid status: %d", status)
	}
}

func testNamedLocationClient_UpdateCountry(t *testing.T, c NamedLocationClientTest, cnl msgraph.CountryNamedLocation) {
	status, err := c.client.UpdateCountry(c.connection.Context, cnl)
	if err != nil {
		t.Fatalf("NamedLocationClient.UpdateCountry(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("NamedLocationClient.UpdateCountry(): invalid status: %d", status)
	}
}

func testNamedLocationClient_List(t *testing.T, c NamedLocationClientTest, f string) (namedLocations *[]msgraph.NamedLocation) {
	namedLocations, _, err := c.client.List(c.connection.Context, f)
	if err != nil {
		t.Fatalf("NamedLocationClient.List(): %v", err)
	}
	if namedLocations == nil {
		t.Fatal("NamedLocationClient.List(): ipNamedLocations was nil")
	}
	for _, loc := range *namedLocations {
		_, ok1 := loc.(msgraph.CountryNamedLocation)
		_, ok2 := loc.(msgraph.IPNamedLocation)
		if !ok1 && !ok2 {
			t.Fatal("NamedLocationsClient.List(): a NamedLocation was returned that did not match a known model")
		}
	}
	return
}

func testNamedLocationClient_Delete(t *testing.T, c NamedLocationClientTest, id string) {
	status, err := c.client.Delete(c.connection.Context, id)
	if err != nil {
		t.Fatalf("NamedLocationClient.Delete(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("NamedLocationClient.Delete(): invalid status: %d", status)
	}
}
