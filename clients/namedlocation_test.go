package clients_test

import (
	"fmt"
	"testing"

	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/clients"
	"github.com/manicminer/hamilton/clients/internal"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/models"
)

type NamedLocationClientTest struct {
	connection   *internal.Connection
	client       *clients.NamedLocationClient
	randomString string
}

func TestNamedLocationClient(t *testing.T) {
	rs := internal.RandomString()
	c := NamedLocationClientTest{
		connection:   internal.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: rs,
	}
	c.client = clients.NewNamedLocationClient(c.connection.AuthConfig.TenantID)
	c.client.BaseClient.Authorizer = c.connection.Authorizer

	newIPNamedLocation := models.IPNamedLocation{
		BaseNamedLocation: &models.BaseNamedLocation{
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

	newCountryNamedLocation := models.CountryNamedLocation{
		BaseNamedLocation: &models.BaseNamedLocation{
			DisplayName: utils.StringPtr("Test Country Named Location")},
		CountriesAndRegions: &[]string{"US", "GB"},
	}

	ipNamedLocation := testNamedLocationClient_CreateIP(t, c, newIPNamedLocation)
	countryNamedLocation := testNamedLocationClient_CreateCountry(t, c, newCountryNamedLocation)

	ipNamedLocation.DisplayName = utils.StringPtr(fmt.Sprintf("test-updated-ipnl-%s", c.randomString))
	countryNamedLocation.DisplayName = utils.StringPtr(fmt.Sprintf("test-updated-cnl-%s", c.randomString))

	testNamedLocationClient_UpdateIP(t, c, *ipNamedLocation)
	testNamedLocationClient_UpdateCountry(t, c, *countryNamedLocation)

	namedLocationSlice := testNamedLocationClient_List(t, c, "")
	if namedLocationSlice != nil {
		for _, l := range *namedLocationSlice {
			t.Logf("%+v ", l)
			if _, ok := l.(models.CountryNamedLocation); ok {
				t.Logf("is a CountryNamedLocation")
			} else if _, ok := l.(models.IPNamedLocation); ok {
				t.Logf("is an IPNamedLocation")
			} else {
				t.Logf("Did not match a type")
			}
		}
	}
	// Running get after the update to give the API a chance to catch up
	testNamedLocationClient_GetIP(t, c, *ipNamedLocation.ID)
	testNamedLocationClient_GetCountry(t, c, *countryNamedLocation.ID)

	testNamedLocationClient_Delete(t, c, *ipNamedLocation.ID)
	testNamedLocationClient_Delete(t, c, *countryNamedLocation.ID)
}

func testNamedLocationClient_CreateIP(t *testing.T, c NamedLocationClientTest, ipnl models.IPNamedLocation) (ipNamedLocation *models.IPNamedLocation) {
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

func testNamedLocationClient_CreateCountry(t *testing.T, c NamedLocationClientTest, cnl models.CountryNamedLocation) (countryNamedLocation *models.CountryNamedLocation) {
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

func testNamedLocationClient_GetIP(t *testing.T, c NamedLocationClientTest, id string) (ipNamedLocation *models.IPNamedLocation) {
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

func testNamedLocationClient_GetCountry(t *testing.T, c NamedLocationClientTest, id string) (countryNamedLocation *models.CountryNamedLocation) {
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

func testNamedLocationClient_UpdateIP(t *testing.T, c NamedLocationClientTest, ipnl models.IPNamedLocation) {
	status, err := c.client.UpdateIP(c.connection.Context, ipnl)
	if err != nil {
		t.Fatalf("NamedLocationClient.UpdateIP(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("NamedLocationClient.UpdateIP(): invalid status: %d", status)
	}
}

func testNamedLocationClient_UpdateCountry(t *testing.T, c NamedLocationClientTest, cnl models.CountryNamedLocation) {
	status, err := c.client.UpdateCountry(c.connection.Context, cnl)
	if err != nil {
		t.Fatalf("NamedLocationClient.UpdateCountry(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("NamedLocationClient.UpdateCountry(): invalid status: %d", status)
	}
}

func testNamedLocationClient_List(t *testing.T, c NamedLocationClientTest, f string) (namedLocations *[]models.NamedLocation) {
	namedLocations, _, err := c.client.List(c.connection.Context, f)
	if err != nil {
		t.Fatalf("NamedLocationClient.List(): %v", err)
	}
	if namedLocations == nil {
		t.Fatal("NamedLocationClient.List(): ipNamedLocations was nil")
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
