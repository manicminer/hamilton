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

// IPNamedLocationClient performs operations on IP Named Locations.
type IPNamedLocationClient struct {
	BaseClient base.Client
}

// NewIPNamedLocationClient returns a new IPNamedLocationClient.
func NewIPNamedLocationClient(tenantId string) *IPNamedLocationClient {
	return &IPNamedLocationClient{
		BaseClient: base.NewClient(base.Version10, tenantId),
	}
}

// List returns a list of IP Named Locations.
func (c *IPNamedLocationClient) List(ctx context.Context) (*[]models.IPNamedLocation, int, error) {
	params := url.Values{}
	params.Add("$filter", "isof('microsoft.graph.ipNamedLocation')")

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
		IPNamedLocations []models.IPNamedLocation `json:"value"`
	}

	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, status, err
	}

	return &data.IPNamedLocations, status, nil
}

// Create creates a new IP Named Location.
func (c *IPNamedLocationClient) Create(ctx context.Context, ipNamedLocation models.IPNamedLocation) (*models.IPNamedLocation, int, error) {
	var status int
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

// Get retrieves an IP Named Location.
func (c *IPNamedLocationClient) Get(ctx context.Context, id string) (*models.IPNamedLocation, int, error) {
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

// Update amends an existing IP Named Location.
func (c *IPNamedLocationClient) Update(ctx context.Context, ipNamedLocation models.IPNamedLocation) (int, error) {
	var status int
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

// Delete removes an IP Named Location.
func (c *IPNamedLocationClient) Delete(ctx context.Context, id string) (int, error) {
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
