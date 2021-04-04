package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/manicminer/hamilton/base"
	"github.com/manicminer/hamilton/models"
)

// CountryNamedLocationClient performs operations on Country Named Locations.
type CountryNamedLocationClient struct {
	BaseClient base.Client
}

// NewCountryNamedLocationClient returns a new CountryNamedLocationClient.
func NewCountryNamedLocationClient(tenantId string) *CountryNamedLocationClient {
	return &CountryNamedLocationClient{
		BaseClient: base.NewClient(base.Version10, tenantId),
	}
}

// List returns a list of Country Named Locations.
func (c *CountryNamedLocationClient) List(ctx context.Context) (*[]models.CountryNamedLocation, int, error) {
	params := url.Values{}
	params.Add("$filter", "isof('microsoft.graph.countryNamedLocation')")

	resp, status, _, err := c.BaseClient.Get(ctx, base.GetHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		Uri: base.Uri{
			Entity:      "/identity/conditionalAccess/namedLocations",
			Params:      params,
			HasTenantId: true,
		},
	})

	if err != nil {
		return nil, status, err
	}

	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)

	var data struct {
		CountryNamedLocations []models.CountryNamedLocation `json:"value"`
	}

	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, status, err
	}

	return &data.CountryNamedLocations, status, nil
}

// Create creates a new Country Named Location.
func (c *CountryNamedLocationClient) Create(ctx context.Context, countryNamedLocation models.CountryNamedLocation) (*models.CountryNamedLocation, int, error) {
	var status int
	body, err := json.Marshal(countryNamedLocation)
	if err != nil {
		return nil, status, err
	}
	resp, status, _, err := c.BaseClient.Post(ctx, base.PostHttpRequestInput{
		Body:             body,
		ValidStatusCodes: []int{http.StatusCreated},
		Uri: base.Uri{
			Entity:      "/identity/conditionalAccess/namedLocations",
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, err
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	var newCountryNamedLocation models.CountryNamedLocation
	if err := json.Unmarshal(respBody, &newCountryNamedLocation); err != nil {
		return nil, status, err
	}
	return &newCountryNamedLocation, status, nil
}

// Get retrieves an Country Named Location.
func (c *CountryNamedLocationClient) Get(ctx context.Context, id string) (*models.CountryNamedLocation, int, error) {
	resp, status, _, err := c.BaseClient.Get(ctx, base.GetHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		Uri: base.Uri{
			Entity:      fmt.Sprintf("/identity/conditionalAccess/namedLocations/%s", id),
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, err
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	var countryNamedLocation models.CountryNamedLocation
	if err := json.Unmarshal(respBody, &countryNamedLocation); err != nil {
		return nil, status, err
	}
	return &countryNamedLocation, status, nil
}

// Update amends an existing Country Named Location.
func (c *CountryNamedLocationClient) Update(ctx context.Context, countryNamedLocation models.CountryNamedLocation) (int, error) {
	var status int
	body, err := json.Marshal(countryNamedLocation)
	if err != nil {
		return status, err
	}
	_, status, _, err = c.BaseClient.Patch(ctx, base.PatchHttpRequestInput{
		Body:             body,
		ValidStatusCodes: []int{http.StatusNoContent},
		Uri: base.Uri{
			Entity:      fmt.Sprintf("/identity/conditionalAccess/namedLocations/%s", *countryNamedLocation.ID),
			HasTenantId: true,
		},
	})
	if err != nil {
		return status, err
	}
	return status, nil
}

// Delete removes an Country Named Location.
func (c *CountryNamedLocationClient) Delete(ctx context.Context, id string) (int, error) {
	_, status, _, err := c.BaseClient.Delete(ctx, base.DeleteHttpRequestInput{
		ValidStatusCodes: []int{http.StatusNoContent},
		Uri: base.Uri{
			Entity:      fmt.Sprintf("/identity/conditionalAccess/namedLocations/%s", id),
			HasTenantId: true,
		},
	})
	if err != nil {
		return status, err
	}
	return status, nil
}
