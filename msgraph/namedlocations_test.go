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

type NamedLocationsClientTest struct {
	connection   *test.Connection
	client       *msgraph.NamedLocationsClient
	randomString string
}

func TestNamedLocationsClient(t *testing.T) {
	c := NamedLocationsClientTest{
		connection:   test.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: test.RandomString(),
	}
	c.client = msgraph.NewNamedLocationsClient(c.connection.AuthConfig.TenantID)
	c.client.BaseClient.Authorizer = c.connection.Authorizer
	c.client.BaseClient.Endpoint = c.connection.AuthConfig.Environment.MsGraph.Endpoint

	newIPNamedLocation := msgraph.IPNamedLocation{
		BaseNamedLocation: &msgraph.BaseNamedLocation{
			DisplayName: utils.StringPtr("test-ip-named-location"),
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
			DisplayName: utils.StringPtr("test-country-named-location"),
		},
		CountriesAndRegions: &[]string{"US", "GB"},
	}

	ipNamedLocation := testNamedLocationsClient_CreateIP(t, c, newIPNamedLocation)
	countryNamedLocation := testNamedLocationsClient_CreateCountry(t, c, newCountryNamedLocation)

	ipNamedLocation.DisplayName = utils.StringPtr(fmt.Sprintf("test-updated-ipnl-%s", c.randomString))
	ipNamedLocation.CreatedDateTime = nil
	ipNamedLocation.ModifiedDateTime = nil
	countryNamedLocation.DisplayName = utils.StringPtr(fmt.Sprintf("test-updated-cnl-%s", c.randomString))
	countryNamedLocation.CreatedDateTime = nil
	countryNamedLocation.ModifiedDateTime = nil

	testNamedLocationsClient_UpdateIP(t, c, *ipNamedLocation)
	testNamedLocationsClient_UpdateCountry(t, c, *countryNamedLocation)

	testNamedLocationsClient_List(t, c)
	// Running get after the update to give the API a chance to catch up
	testNamedLocationsClient_GetIP(t, c, *ipNamedLocation.ID)
	testNamedLocationsClient_GetCountry(t, c, *countryNamedLocation.ID)
	testNamedLocationsClient_Get(t, c, *ipNamedLocation.ID)
	testNamedLocationsClient_Get(t, c, *countryNamedLocation.ID)

	testNamedLocationsClient_Delete(t, c, *ipNamedLocation.ID)
	testNamedLocationsClient_Delete(t, c, *countryNamedLocation.ID)
}

func testNamedLocationsClient_CreateIP(t *testing.T, c NamedLocationsClientTest, ipnl msgraph.IPNamedLocation) (ipNamedLocation *msgraph.IPNamedLocation) {
	ipNamedLocation, status, err := c.client.CreateIP(c.connection.Context, ipnl)
	if err != nil {
		t.Fatalf("NamedLocationsClient.CreateIP(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("NamedLocationsClient.CreateIP(): invalid status: %d", status)
	}
	if ipNamedLocation == nil {
		t.Fatal("NamedLocationsClient.CreateIP(): ipNamedLocation was nil")
	}
	if ipNamedLocation.ID == nil {
		t.Fatal("NamedLocationsClient.CreateIP(): ipNamedLocation.ID was nil")
	}
	return
}

func testNamedLocationsClient_CreateCountry(t *testing.T, c NamedLocationsClientTest, cnl msgraph.CountryNamedLocation) (countryNamedLocation *msgraph.CountryNamedLocation) {
	countryNamedLocation, status, err := c.client.CreateCountry(c.connection.Context, cnl)
	if err != nil {
		t.Fatalf("NamedLocationsClient.CreateCountry(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("NamedLocationsClient.CreateCountry(): invalid status: %d", status)
	}
	if countryNamedLocation == nil {
		t.Fatal("NamedLocationsClient.CreateCountry(): countryNamedLocation was nil")
	}
	if countryNamedLocation.ID == nil {
		t.Fatal("NamedLocationsClient.CreateCountry(): countryNamedLocation.ID was nil")
	}
	return
}

func testNamedLocationsClient_GetIP(t *testing.T, c NamedLocationsClientTest, id string) (ipNamedLocation *msgraph.IPNamedLocation) {
	ipNamedLocation, status, err := c.client.GetIP(c.connection.Context, id, odata.Query{})
	if err != nil {
		t.Fatalf("IPNamedLocationsClient.GetIP(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("IPNamedLocationsClient.GetIP(): invalid status: %d", status)
	}
	if ipNamedLocation == nil {
		t.Fatal("IPNamedLocationsClient.GetIP(): ipNamedLocation was nil")
	}
	return
}

func testNamedLocationsClient_GetCountry(t *testing.T, c NamedLocationsClientTest, id string) (countryNamedLocation *msgraph.CountryNamedLocation) {
	countryNamedLocation, status, err := c.client.GetCountry(c.connection.Context, id, odata.Query{})
	if err != nil {
		t.Fatalf("NamedLocationsClient.GetCountry(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("NamedLocationsClient.GetCountry(): invalid status: %d", status)
	}
	if countryNamedLocation == nil {
		t.Fatal("NamedLocationsClient.GetCountry(): countryNamedLocation was nil")
	}
	return
}

func testNamedLocationsClient_Get(t *testing.T, c NamedLocationsClientTest, id string) (namedLocation *msgraph.NamedLocation) {
	namedLocation, status, err := c.client.Get(c.connection.Context, id, odata.Query{})
	if err != nil {
		t.Fatalf("NamedLocationsClient.Get(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("NamedLocationsClient.Get(): invalid status: %d", status)
	}
	if namedLocation == nil {
		t.Fatal("NamedLocationsClient.Get(): NamedLocation was nil")
	}
	_, ok1 := (*namedLocation).(msgraph.CountryNamedLocation)
	_, ok2 := (*namedLocation).(msgraph.IPNamedLocation)
	if !ok1 && !ok2 {
		t.Fatal("NamedLocationsClient.Get(): a NamedLocation was returned that did not match a known model")
	}
	return
}

func testNamedLocationsClient_UpdateIP(t *testing.T, c NamedLocationsClientTest, ipnl msgraph.IPNamedLocation) {
	status, err := c.client.UpdateIP(c.connection.Context, ipnl)
	if err != nil {
		t.Fatalf("NamedLocationsClient.UpdateIP(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("NamedLocationsClient.UpdateIP(): invalid status: %d", status)
	}
}

func testNamedLocationsClient_UpdateCountry(t *testing.T, c NamedLocationsClientTest, cnl msgraph.CountryNamedLocation) {
	status, err := c.client.UpdateCountry(c.connection.Context, cnl)
	if err != nil {
		t.Fatalf("NamedLocationsClient.UpdateCountry(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("NamedLocationsClient.UpdateCountry(): invalid status: %d", status)
	}
}

func testNamedLocationsClient_List(t *testing.T, c NamedLocationsClientTest) (namedLocations *[]msgraph.NamedLocation) {
	namedLocations, _, err := c.client.List(c.connection.Context, odata.Query{})
	if err != nil {
		t.Fatalf("NamedLocationsClient.List(): %v", err)
	}
	if namedLocations == nil {
		t.Fatal("NamedLocationsClient.List(): ipNamedLocations was nil")
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

func testNamedLocationsClient_Delete(t *testing.T, c NamedLocationsClientTest, id string) {
	status, err := c.client.Delete(c.connection.Context, id)
	if err != nil {
		t.Fatalf("NamedLocationsClient.Delete(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("NamedLocationsClient.Delete(): invalid status: %d", status)
	}
}
