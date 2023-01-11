package msgraph

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

var actionAdminAssign string = "adminAssign"
var actionAdminRemove string = "adminRemove"

// RoleEligibilityScheduleRequestClient performs operations on RoleEligibilityScheduleRequests.
type RoleEligibilityScheduleRequestClient struct {
	BaseClient Client
}

// NewRoleEligibilityScheduleRequest returns a new RoleEligibilityScheduleRequestClient
func NewRoleEligibilityScheduleRequestClient() *RoleEligibilityScheduleRequestClient {
	return &RoleEligibilityScheduleRequestClient{
		BaseClient: NewClient(Version10),
	}
}

// Get retrieves a UnifiedRoleEligibilityScheduleRequest
func (c *RoleEligibilityScheduleRequestClient) Get(ctx context.Context, id string, query odata.Query) (*UnifiedRoleEligibilityScheduleRequest, int, error) {
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		OData:                  query,
		ValidStatusCodes:       []int{http.StatusOK},
		Uri: Uri{
			Entity: fmt.Sprintf("/roleManagement/directory/roleEligibilityScheduleRequests/%s", id),
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("RoleEligibilityScheduleRequestClient.BaseClient.Get(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var dirRole UnifiedRoleEligibilityScheduleRequest
	if err := json.Unmarshal(respBody, &dirRole); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &dirRole, status, nil
}

// Create creates a new UnifiedRoleEligibilityScheduleRequest.
func (c *RoleEligibilityScheduleRequestClient) Create(ctx context.Context, resr UnifiedRoleEligibilityScheduleRequest) (*UnifiedRoleEligibilityScheduleRequest, int, error) {
	var status int

	resr.Action = &actionAdminAssign
	body, err := json.Marshal(resr)
	if err != nil {
		return nil, status, fmt.Errorf("json.Marshal(): %v", err)
	}

	resp, status, _, err := c.BaseClient.Post(ctx, PostHttpRequestInput{
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		Body:                   body,
		ValidStatusCodes:       []int{http.StatusCreated},
		Uri: Uri{
			Entity: "/roleManagement/directory/roleEligibilityScheduleRequests",
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("RoleEligibilityScheduleRequestClient.BaseClient.Post(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var newRoleAssignment UnifiedRoleEligibilityScheduleRequest
	if err := json.Unmarshal(respBody, &newRoleAssignment); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &newRoleAssignment, status, nil
}

// Delete removes a UnifiedRoleEligibilityScheduleRequest.
func (c *RoleEligibilityScheduleRequestClient) Delete(ctx context.Context, id string) (int, error) {
	resr, status, err := c.Get(ctx, id, odata.Query{})
	if err != nil {
		return status, fmt.Errorf("RoleEligibilityScheduleRequestClient.Get(): %v", err)
	}
	resrBody := UnifiedRoleEligibilityScheduleRequest{
		Action:           &actionAdminRemove,
		RoleDefinitionId: resr.RoleDefinitionId,
		DirectoryScopeId: resr.DirectoryScopeId,
		PrincipalId:      resr.PrincipalId,
	}
	body, err := json.Marshal(resrBody)
	if err != nil {
		return status, fmt.Errorf("json.Marshal(): %v", err)
	}
	_, status, _, err = c.BaseClient.Post(ctx, PostHttpRequestInput{
		Body:             body,
		ValidStatusCodes: []int{http.StatusCreated},
		Uri: Uri{
			Entity: "/roleManagement/directory/roleEligibilityScheduleRequests",
		},
	})
	if err != nil {
		return status, fmt.Errorf("RoleAssignments.BaseClient.Get(): %v", err)
	}

	return status, nil
}
