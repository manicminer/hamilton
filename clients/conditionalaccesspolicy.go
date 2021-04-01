package clients

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/manicminer/hamilton/base"
	"github.com/manicminer/hamilton/models"
)

// ConditionalAccessPolicyClient performs operations on ConditionalAccessPolicy.
type ConditionalAccessPolicyClient struct {
	BaseClient base.Client
}

// NewConditionalAccessPolicyClient returns a new ConditionalAccessPolicyClient
func NewConditionalAccessPolicyClient(tenantId string) *ConditionalAccessPolicyClient {
	return &ConditionalAccessPolicyClient{
		BaseClient: base.NewClient(base.VersionBeta, tenantId),
	}
}

// List returns a list of ConditionalAccessPolicys, optionally filtered using OData.
func (c *ConditionalAccessPolicyClient) List(ctx context.Context, filter string) (*[]models.ConditionalAccessPolicy, int, error) {
	params := url.Values{}
	if filter != "" {
		params.Add("$filter", filter)
	}
	resp, status, _, err := c.BaseClient.Get(ctx, base.GetHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		Uri: base.Uri{
			Entity:      "/identity/conditionalAccess/policies",
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
		ConditionalAccessPolicys []models.ConditionalAccessPolicy `json:"value"`
	}
	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, status, err
	}
	return &data.ConditionalAccessPolicys, status, nil
}

// Create creates a new ConditionalAccessPolicy.
func (c *ConditionalAccessPolicyClient) Create(ctx context.Context, conditionalAccessPolicy models.ConditionalAccessPolicy) (*models.ConditionalAccessPolicy, int, error) {
	var status int
	body, err := json.Marshal(conditionalAccessPolicy)
	if err != nil {
		return nil, status, err
	}
	resp, status, _, err := c.BaseClient.Post(ctx, base.PostHttpRequestInput{
		Body:             body,
		ValidStatusCodes: []int{http.StatusCreated},
		Uri: base.Uri{
			Entity:      "/identity/conditionalAccess/policies",
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, err
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	var newConditionalAccessPolicy models.ConditionalAccessPolicy
	if err := json.Unmarshal(respBody, &newConditionalAccessPolicy); err != nil {
		return nil, status, err
	}
	return &newConditionalAccessPolicy, status, nil
}

// Get retrieves an ConditionalAccessPolicy.
func (c *ConditionalAccessPolicyClient) Get(ctx context.Context, id string) (*models.ConditionalAccessPolicy, int, error) {
	resp, status, _, err := c.BaseClient.Get(ctx, base.GetHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		Uri: base.Uri{
			Entity:      fmt.Sprintf("/identity/conditionalAccess/policies/%s", id),
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, err
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	var conditionalAccessPolicy models.ConditionalAccessPolicy
	if err := json.Unmarshal(respBody, &conditionalAccessPolicy); err != nil {
		return nil, status, err
	}
	return &conditionalAccessPolicy, status, nil
}

// Update amends an existing ConditionalAccessPolicy.
func (c *ConditionalAccessPolicyClient) Update(ctx context.Context, conditionalAccessPolicy models.ConditionalAccessPolicy) (int, error) {
	var status int
	if conditionalAccessPolicy.ID == nil {
		return status, errors.New("cannot update conditionalAccessPolicy with nil ID")
	}
	body, err := json.Marshal(conditionalAccessPolicy)
	if err != nil {
		return status, err
	}
	_, status, _, err = c.BaseClient.Patch(ctx, base.PatchHttpRequestInput{
		Body:             body,
		ValidStatusCodes: []int{http.StatusNoContent},
		Uri: base.Uri{
			Entity:      fmt.Sprintf("/identity/conditionalAccess/policies/%s", *conditionalAccessPolicy.ID),
			HasTenantId: true,
		},
	})
	if err != nil {
		return status, err
	}
	return status, nil
}

// Delete removes a ConditionalAccessPolicy.
func (c *ConditionalAccessPolicyClient) Delete(ctx context.Context, id string) (int, error) {
	_, status, _, err := c.BaseClient.Delete(ctx, base.DeleteHttpRequestInput{
		ValidStatusCodes: []int{http.StatusNoContent},
		Uri: base.Uri{
			Entity:      fmt.Sprintf("/identity/conditionalAccess/policies/%s", id),
			HasTenantId: true,
		},
	})
	if err != nil {
		return status, err
	}
	return status, nil
}
