package msgraph_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/msgraph"
)

func TestNamedLocationsClient(t *testing.T) {
	c := test.NewTest(t)
	defer c.CancelFunc()

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

	ipNamedLocation.DisplayName = utils.StringPtr(fmt.Sprintf("test-updated-ipnl-%s", c.RandomString))
	ipNamedLocation.CreatedDateTime = nil
	ipNamedLocation.ModifiedDateTime = nil
	countryNamedLocation.DisplayName = utils.StringPtr(fmt.Sprintf("test-updated-cnl-%s", c.RandomString))
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

	ipNamedLocation.IsTrusted = utils.BoolPtr(false)
	testNamedLocationsClient_UpdateIP(t, c, *ipNamedLocation)

	countryNamedLocation.IncludeUnknownCountriesAndRegions = utils.BoolPtr(true)
	testNamedLocationsClient_UpdateCountry(t, c, *countryNamedLocation)

	testNamedLocationsClient_Delete(t, c, *ipNamedLocation.ID)
	testNamedLocationsClient_Delete(t, c, *countryNamedLocation.ID)
}

func testNamedLocationsClient_CreateIP(t *testing.T, c *test.Test, ipnl msgraph.IPNamedLocation) (ipNamedLocation *msgraph.IPNamedLocation) {
	ipNamedLocation, status, err := c.NamedLocationsClient.CreateIP(c.Context, ipnl)
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

func testNamedLocationsClient_CreateCountry(t *testing.T, c *test.Test, cnl msgraph.CountryNamedLocation) (countryNamedLocation *msgraph.CountryNamedLocation) {
	countryNamedLocation, status, err := c.NamedLocationsClient.CreateCountry(c.Context, cnl)
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

func testNamedLocationsClient_GetIP(t *testing.T, c *test.Test, id string) (ipNamedLocation *msgraph.IPNamedLocation) {
	ipNamedLocation, status, err := c.NamedLocationsClient.GetIP(c.Context, id, odata.Query{})
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

func testNamedLocationsClient_GetCountry(t *testing.T, c *test.Test, id string) (countryNamedLocation *msgraph.CountryNamedLocation) {
	countryNamedLocation, status, err := c.NamedLocationsClient.GetCountry(c.Context, id, odata.Query{})
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

func testNamedLocationsClient_Get(t *testing.T, c *test.Test, id string) (namedLocation *msgraph.NamedLocation) {
	namedLocation, status, err := c.NamedLocationsClient.Get(c.Context, id, odata.Query{})
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

func testNamedLocationsClient_UpdateIP(t *testing.T, c *test.Test, ipnl msgraph.IPNamedLocation) {
	status, err := c.NamedLocationsClient.UpdateIP(c.Context, ipnl)
	if err != nil {
		t.Fatalf("NamedLocationsClient.UpdateIP(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("NamedLocationsClient.UpdateIP(): invalid status: %d", status)
	}
}

func testNamedLocationsClient_UpdateCountry(t *testing.T, c *test.Test, cnl msgraph.CountryNamedLocation) {
	status, err := c.NamedLocationsClient.UpdateCountry(c.Context, cnl)
	if err != nil {
		t.Fatalf("NamedLocationsClient.UpdateCountry(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("NamedLocationsClient.UpdateCountry(): invalid status: %d", status)
	}
}

func testNamedLocationsClient_List(t *testing.T, c *test.Test) (namedLocations *[]msgraph.NamedLocation) {
	namedLocations, _, err := c.NamedLocationsClient.List(c.Context, odata.Query{})
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

func testNamedLocationsClient_Delete(t *testing.T, c *test.Test, id string) {
	status, err := c.NamedLocationsClient.Delete(c.Context, id)
	if err != nil {
		t.Fatalf("NamedLocationsClient.Delete(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("NamedLocationsClient.Delete(): invalid status: %d", status)
	}
}
