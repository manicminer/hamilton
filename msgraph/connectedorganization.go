package msgraph

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/manicminer/hamilton/odata"
)

type ConnectedOrganizationClient struct {
	BaseClient Client
}

func NewConnectedOrganizationClient(tenantId string) *ConnectedOrganizationClient {
	return &ConnectedOrganizationClient{
		BaseClient: NewClient(Version10, tenantId),
	}
}

// List returns a list of ConnectedOrganization
func (c *ConnectedOrganizationClient) List(ctx context.Context, query odata.Query) (*[]ConnectedOrganization, int, error) {
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		DisablePaging:    query.Top > 0,
		OData:            query,
		ValidStatusCodes: []int{http.StatusOK},
		Uri: Uri{
			Entity:      "/identityGovernance/entitlementManagement/connectedOrganizations",
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("ConnectedOrganizationClient.BaseClient.Get(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var data struct {
		ConnectedOrganizations []ConnectedOrganization `json:"value"`
	}
	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &data.ConnectedOrganizations, status, nil
}

// Create creates a new ConnectedOrganization.
func (c *ConnectedOrganizationClient) Create(ctx context.Context, connectedOrganization ConnectedOrganization) (*ConnectedOrganization, int, error) {
	var status int
	body, err := json.Marshal(connectedOrganization)
	if err != nil {
		return nil, status, fmt.Errorf("json.Marshal(): %v", err)
	}

	resp, status, _, err := c.BaseClient.Post(ctx, PostHttpRequestInput{
		Body:             body,
		ValidStatusCodes: []int{http.StatusCreated},
		Uri: Uri{
			Entity:      "/identityGovernance/entitlementManagement/connectedOrganizations",
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("ConnectedOrganizationClient.BaseClient.Post(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var newConnectedOrganization ConnectedOrganization
	if err := json.Unmarshal(respBody, &newConnectedOrganization); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &newConnectedOrganization, status, nil
}

// Get retrieves a ConnectedOrganization.
func (c *ConnectedOrganizationClient) Get(ctx context.Context, id string, query odata.Query) (*ConnectedOrganization, int, error) {
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		OData:                  query,
		ValidStatusCodes:       []int{http.StatusOK},
		Uri: Uri{
			Entity:      fmt.Sprintf("/identityGovernance/entitlementManagement/connectedOrganizations/%s", id),
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("ConnectedOrganizationClient.BaseClient.Get(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var connectedOrganization ConnectedOrganization
	if err := json.Unmarshal(respBody, &connectedOrganization); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &connectedOrganization, status, nil
}

// Update amends an existing ConnectedOrganization.
func (c *ConnectedOrganizationClient) Update(ctx context.Context, connectedOrganization ConnectedOrganization) (int, error) {
	var status int

	if connectedOrganization.ID == nil {
		return status, errors.New("cannot update ConnectedOrganization with nil ID")
	}

	body, err := json.Marshal(connectedOrganization)
	if err != nil {
		return status, fmt.Errorf("json.Marshal(): %v", err)
	}

	_, status, _, err = c.BaseClient.Patch(ctx, PatchHttpRequestInput{
		Body:                   body,
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		ValidStatusCodes:       []int{http.StatusNoContent},
		Uri: Uri{
			Entity:      fmt.Sprintf("/identityGovernance/entitlementManagement/connectedOrganizations/%s", *connectedOrganization.ID),
			HasTenantId: true,
		},
	})
	if err != nil {
		return status, fmt.Errorf("ConnectedOrganizationClient.BaseClient.Patch(): %v", err)
	}

	return status, nil
}

// Delete removes a ConnectedOrganization.
func (c *ConnectedOrganizationClient) Delete(ctx context.Context, id string) (int, error) {
	_, status, _, err := c.BaseClient.Delete(ctx, DeleteHttpRequestInput{
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		ValidStatusCodes:       []int{http.StatusNoContent},
		Uri: Uri{
			Entity:      fmt.Sprintf("/identityGovernance/entitlementManagement/connectedOrganizations/%s", id),
			HasTenantId: true,
		},
	})
	if err != nil {
		return status, fmt.Errorf("ConnectedOrganizationClient.BaseClient.Delete(): %v", err)
	}

	return status, nil
}
