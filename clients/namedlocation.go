package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/manicminer/hamilton/base"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/models"
)

// NamedLocationClient performs operations on Named Locations.
type NamedLocationClient struct {
	BaseClient base.Client
}

// NewNamedLocationClient returns a new NamedLocationClient.
func NewNamedLocationClient(tenantId string) *NamedLocationClient {
	return &NamedLocationClient{
		BaseClient: base.NewClient(base.Version10, tenantId),
	}
}

// List returns a list of Named Locations, optionally filtered using OData.
func (c *NamedLocationClient) List(ctx context.Context, filter string) (*[]models.NamedLocation, int, error) {
	params := url.Values{}
	if filter != "" {
		params.Add("$filter", filter)
	}

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
		NamedLocations []models.NamedLocation `json:"value"`
	}

	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, status, err
	}

	return &data.NamedLocations, status, nil
}

// Delete removes a Named Location.
func (c *NamedLocationClient) Delete(ctx context.Context, id string) (int, error) {
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

// CreateIP creates a new IP Named Location.
func (c *NamedLocationClient) CreateIP(ctx context.Context, ipNamedLocation models.IPNamedLocation) (*models.IPNamedLocation, int, error) {
	var status int

	ipNamedLocation.ODataType = utils.StringPtr("#microsoft.graph.ipNamedLocation")
	// This API does not handle PATCH on some properties
	// i.e. cannot update values such as CreatedDateTime
	ipNamedLocation.CreatedDateTime = nil
	ipNamedLocation.ModifiedDateTime = nil

	body, err := json.Marshal(ipNamedLocation)
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
	var newIPNamedLocation models.IPNamedLocation
	if err := json.Unmarshal(respBody, &newIPNamedLocation); err != nil {
		return nil, status, err
	}
	return &newIPNamedLocation, status, nil
}

// CreateCountry creates a new Country Named Location.
func (c *NamedLocationClient) CreateCountry(ctx context.Context, countryNamedLocation models.CountryNamedLocation) (*models.CountryNamedLocation, int, error) {
	var status int

	countryNamedLocation.ODataType = utils.StringPtr("#microsoft.graph.countryNamedLocation")
	// This API does not handle PATCH on some properties
	// i.e. cannot update values such as CreatedDateTime
	countryNamedLocation.CreatedDateTime = nil
	countryNamedLocation.ModifiedDateTime = nil

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

// GetIP retrieves an IP Named Location.
func (c *NamedLocationClient) GetIP(ctx context.Context, id string) (*models.IPNamedLocation, int, error) {
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
	var ipNamedLocation models.IPNamedLocation
	if err := json.Unmarshal(respBody, &ipNamedLocation); err != nil {
		return nil, status, err
	}
	return &ipNamedLocation, status, nil
}

// GetCountry retrieves an Country Named Location.
func (c *NamedLocationClient) GetCountry(ctx context.Context, id string) (*models.CountryNamedLocation, int, error) {
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

// UpdateIP amends an existing IP Named Location.
func (c *NamedLocationClient) UpdateIP(ctx context.Context, ipNamedLocation models.IPNamedLocation) (int, error) {
	var status int

	// This API does not handle PATCH on some properties
	// i.e. cannot update values such as CreatedDateTime
	ipNamedLocation.CreatedDateTime = nil
	ipNamedLocation.ModifiedDateTime = nil

	body, err := json.Marshal(ipNamedLocation)
	if err != nil {
		return status, err
	}
	_, status, _, err = c.BaseClient.Patch(ctx, base.PatchHttpRequestInput{
		Body:             body,
		ValidStatusCodes: []int{http.StatusNoContent},
		Uri: base.Uri{
			Entity:      fmt.Sprintf("/identity/conditionalAccess/namedLocations/%s", *ipNamedLocation.ID),
			HasTenantId: true,
		},
	})
	if err != nil {
		return status, err
	}
	return status, nil
}

// UpdateCountry amends an existing Country Named Location.
func (c *NamedLocationClient) UpdateCountry(ctx context.Context, countryNamedLocation models.CountryNamedLocation) (int, error) {
	var status int

	// This API does not handle PATCH on some properties
	// i.e. cannot update values such as CreatedDateTime
	countryNamedLocation.CreatedDateTime = nil
	countryNamedLocation.ModifiedDateTime = nil

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
