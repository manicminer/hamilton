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

type CountryNamedLocationClientTest struct {
	connection   *internal.Connection
	client       *clients.CountryNamedLocationClient
	randomString string
}

func TestCountryNamedLocationClient(t *testing.T) {
	rs := internal.RandomString()
	c := CountryNamedLocationClientTest{
		connection:   internal.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: rs,
	}
	c.client = clients.NewCountryNamedLocationClient(c.connection.AuthConfig.TenantID)
	c.client.BaseClient.Authorizer = c.connection.Authorizer

	newCountryNamedLocation := models.CountryNamedLocation{
		NamedLocation: &models.NamedLocation{
			ODataType:   utils.StringPtr("#microsoft.graph.countryNamedLocation"),
			DisplayName: utils.StringPtr("Test Country Named Location")},
		CountriesAndRegions: &[]string{"US", "GB"},
	}

	countryNamedLocation := testCountryNamedLocationClient_Create(t, c, newCountryNamedLocation)
	// Running get too quickly after create often results in the resource not being found
	time.Sleep(5 * time.Second)
	testCountryNamedLocationClient_Get(t, c, *countryNamedLocation.ID)

	countryNamedLocation.DisplayName = utils.StringPtr(fmt.Sprintf("test-updated-cnl-%s", c.randomString))
	testCountryNamedLocationClient_Update(t, c, *countryNamedLocation)

	testCountryNamedLocationClient_List(t, c)
	testCountryNamedLocationClient_Delete(t, c, *countryNamedLocation.ID)
}

func testCountryNamedLocationClient_Create(t *testing.T, c CountryNamedLocationClientTest, cnl models.CountryNamedLocation) (countryNamedLocation *models.CountryNamedLocation) {
	countryNamedLocation, status, err := c.client.Create(c.connection.Context, cnl)
	if err != nil {
		t.Fatalf("CountryNamedLocationClient.Create(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("CountryNamedLocationClient.Create(): invalid status: %d", status)
	}
	if countryNamedLocation == nil {
		t.Fatal("CountryNamedLocationClient.Create(): countryNamedLocation was nil")
	}
	if countryNamedLocation.ID == nil {
		t.Fatal("CountryNamedLocationClient.Create(): countryNamedLocation.ID was nil")
	}
	return
}

func testCountryNamedLocationClient_Get(t *testing.T, c CountryNamedLocationClientTest, id string) (countryNamedLocation *models.CountryNamedLocation) {
	countryNamedLocation, status, err := c.client.Get(c.connection.Context, id)
	if err != nil {
		t.Fatalf("CountryNamedLocationClient.Get(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("CountryNamedLocationClient.Get(): invalid status: %d", status)
	}
	if countryNamedLocation == nil {
		t.Fatal("CountryNamedLocationClient.Get(): countryNamedLocation was nil")
	}
	return
}

func testCountryNamedLocationClient_Update(t *testing.T, c CountryNamedLocationClientTest, cnl models.CountryNamedLocation) {
	status, err := c.client.Update(c.connection.Context, cnl)
	if err != nil {
		t.Fatalf("CountryNamedLocationClient.Update(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("CountryNamedLocationClient.Update(): invalid status: %d", status)
	}
}

func testCountryNamedLocationClient_List(t *testing.T, c CountryNamedLocationClientTest) (countryNamedLocations *[]models.CountryNamedLocation) {
	countryNamedLocations, _, err := c.client.List(c.connection.Context)
	if err != nil {
		t.Fatalf("CountryNamedLocationClient.List(): %v", err)
	}
	if countryNamedLocations == nil {
		t.Fatal("CountryNamedLocationClient.List(): countryNamedLocations was nil")
	}
	return
}

func testCountryNamedLocationClient_Delete(t *testing.T, c CountryNamedLocationClientTest, id string) {
	status, err := c.client.Delete(c.connection.Context, id)
	if err != nil {
		t.Fatalf("CountryNamedLocationClient.Delete(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("CountryNamedLocationClient.Delete(): invalid status: %d", status)
	}
}
