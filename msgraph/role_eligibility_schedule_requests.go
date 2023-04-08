package msgraph

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// RoleEligibilityScheduleRequestsClient performs operations on RoleEligibilityScheduleRequests.
type RoleEligibilityScheduleRequestsClient struct {
	BaseClient Client
}

// NewRoleEligibilityScheduleRequestsClient returns a new RoleEligibilityScheduleRequestsClient
func NewRoleEligibilityScheduleRequestsClient() *RoleEligibilityScheduleRequestsClient {
	return &RoleEligibilityScheduleRequestsClient{
		BaseClient: NewClient(Version10),
	}
}

// List returns a list of RoleEligibilityScheduleRequestsClient
func (c *RoleEligibilityScheduleRequestsClient) List(ctx context.Context, query odata.Query) (*[]UnifiedRoleEligibilityScheduleRequest, int, error) {
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		OData:            query,
		ValidStatusCodes: []int{http.StatusOK},
		Uri: Uri{
			Entity: "/roleManagement/directory/roleEligibilityScheduleRequests",
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("RoleEligibilityScheduleRequestsClient.BaseClient.Get(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var data struct {
		RoleEligibilityScheduleRequest []UnifiedRoleEligibilityScheduleRequest `json:"value"`
	}
	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &data.RoleEligibilityScheduleRequest, status, nil
}

// Get retrieves a UnifiedRoleEligibilityScheduleRequest
func (c *RoleEligibilityScheduleRequestsClient) Get(ctx context.Context, id string, query odata.Query) (*UnifiedRoleEligibilityScheduleRequest, int, error) {
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		OData:                  query,
		ValidStatusCodes:       []int{http.StatusOK},
		Uri: Uri{
			Entity: fmt.Sprintf("/roleManagement/directory/roleEligibilityScheduleRequests/%s", id),
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("RoleEligibilityScheduleRequestsClient.BaseClient.Get(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var roleEligibilityScheduleRequest UnifiedRoleEligibilityScheduleRequest
	if err := json.Unmarshal(respBody, &roleEligibilityScheduleRequest); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &roleEligibilityScheduleRequest, status, nil
}

// Create creates a new UnifiedRoleEligibilityScheduleRequest.
func (c *RoleEligibilityScheduleRequestsClient) Create(ctx context.Context, roleEligibilityScheduleRequest UnifiedRoleEligibilityScheduleRequest) (*UnifiedRoleEligibilityScheduleRequest, int, error) {
	var status int

	body, err := json.Marshal(roleEligibilityScheduleRequest)
	if err != nil {
		return nil, status, fmt.Errorf("json.Marshal(): %v", err)
	}

	resp, status, _, err := c.BaseClient.Post(ctx, PostHttpRequestInput{
		Body:             body,
		ValidStatusCodes: []int{http.StatusCreated},
		Uri: Uri{
			Entity: "/roleManagement/directory/roleEligibilityScheduleRequests",
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("RoleEligibilityScheduleRequestsClient.BaseClient.Post(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var newRoleEligibilityScheduleRequest UnifiedRoleEligibilityScheduleRequest
	if err := json.Unmarshal(respBody, &newRoleEligibilityScheduleRequest); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &newRoleEligibilityScheduleRequest, status, nil
}
