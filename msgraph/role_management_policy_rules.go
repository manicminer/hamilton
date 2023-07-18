package msgraph

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type RoleManagementPolicyRulesClient struct {
	BaseClient Client
}

func NewRoleManagementPolicyRulesClient() *RoleManagementPolicyRulesClient {
	return &RoleManagementPolicyRulesClient{
		BaseClient: NewClient(VersionBeta),
	}
}

// List retrieves a list of Role Management Rules for a policy
func (c *RoleManagementPolicyRulesClient) List(ctx context.Context, policyId string, query odata.Query) (*[]UnifiedRoleManagementPolicyRule, int, error) {
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		OData:            query,
		ValidStatusCodes: []int{http.StatusOK},
		Uri: Uri{
			Entity: fmt.Sprintf("/policies/roleManagementPolicies/%s/rules", policyId),
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("RoleManagementPolicyRulesClient.BaseClient.Get(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var data struct {
		UnifiedRoleManagementPolicyRule []UnifiedRoleManagementPolicyRule `json:"value"`
	}
	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &data.UnifiedRoleManagementPolicyRule, status, nil
}

// Get retrieves a UnifiedRoleManagementPolicyRule
func (c *RoleManagementPolicyRulesClient) Get(ctx context.Context, id, policyId string, query odata.Query) (*UnifiedRoleManagementPolicyRule, int, error) {
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		OData:            query,
		ValidStatusCodes: []int{http.StatusOK},
		Uri: Uri{
			Entity: fmt.Sprintf("/policies/roleManagementPolicies/%s/rules/%s", policyId, id),
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("RoleManagementPolicyRulesClient.BaseClient.Get(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var rule UnifiedRoleManagementPolicyRule
	if err := json.Unmarshal(respBody, &rule); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &rule, status, nil
}

// Update amends an existing UnifiedRoleManagementPolicyRule.
func (c *RoleManagementPolicyRulesClient) Update(ctx context.Context, rule UnifiedRoleManagementPolicyRule, policyId string) (int, error) {
	var status int

	body, err := json.Marshal(rule)
	if err != nil {
		return status, fmt.Errorf("json.Marshal(): %v", err)
	}

	_, status, _, err = c.BaseClient.Patch(ctx, PatchHttpRequestInput{
		Body:                   body,
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		ValidStatusCodes:       []int{http.StatusNoContent},
		Uri: Uri{
			Entity: fmt.Sprintf("/policies/roleManagementPolicies/%s/rules/%s", policyId, *rule.ID),
		},
	})
	if err != nil {
		return status, fmt.Errorf("RoleManagementPolicyRulesClient.BaseClient.Patch(): %v", err)
	}

	return status, nil
}
