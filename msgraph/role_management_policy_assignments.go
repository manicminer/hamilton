package msgraph

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type RoleManagementPolicyAssignmentsClient struct {
	BaseClient Client
}

func NewRoleManagementPolicyAssignmentsClient() *RoleManagementPolicyAssignmentsClient {
	return &RoleManagementPolicyAssignmentsClient{
		BaseClient: NewClient(VersionBeta),
	}
}

// List retrieves a list of Role Management Assignments for a policy
func (c *RoleManagementPolicyAssignmentsClient) List(ctx context.Context, policyId string, query odata.Query) (*[]UnifiedRoleManagementPolicyAssignment, int, error) {
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		OData:            query,
		ValidStatusCodes: []int{http.StatusOK},
		Uri: Uri{
			Entity: fmt.Sprintf("/policies/roleManagementPolicies/%s/rules", policyId),
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("RoleManagementPolicyAssignmentsClient.BaseClient.Get(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var data struct {
		UnifiedRoleManagementPolicyAssignment []UnifiedRoleManagementPolicyAssignment `json:"value"`
	}
	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &data.UnifiedRoleManagementPolicyAssignment, status, nil
}

// Get retrieves a UnifiedRoleManagementPolicyAssignment
func (c *RoleManagementPolicyAssignmentsClient) Get(ctx context.Context, id, policyId string, query odata.Query) (*UnifiedRoleManagementPolicyAssignment, int, error) {
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		OData:            query,
		ValidStatusCodes: []int{http.StatusOK},
		Uri: Uri{
			Entity: fmt.Sprintf("/policies/roleManagementPolicies/%s/rules/%s", policyId, id),
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("RoleManagementPolicyAssignmentsClient.BaseClient.Get(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var assignment UnifiedRoleManagementPolicyAssignment
	if err := json.Unmarshal(respBody, &assignment); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &assignment, status, nil
}
