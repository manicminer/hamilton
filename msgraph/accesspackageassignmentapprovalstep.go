package msgraph

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type AccessPackageAssignmentApprovalStepClient struct {
	BaseClient Client
}

func NewAccessPackageAssignmentApprovalStepClient() *AccessPackageAssignmentApprovalStepClient {
	return &AccessPackageAssignmentApprovalStepClient{
		BaseClient: NewClient(Version10),
	}
}

// List returns a list of Approval Steps for an approval
func (c *AccessPackageAssignmentApprovalStepClient) List(ctx context.Context, approvalId string, query odata.Query) (*[]ApprovalStep, int, error) {
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		DisablePaging:    query.Top > 0,
		OData:            query,
		ValidStatusCodes: []int{http.StatusOK},
		Uri: Uri{
			Entity: fmt.Sprintf("/identityGovernance/entitlementManagement/accessPackageAssignmentApprovals/%s/steps", approvalId),
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("AccessPackageAssignmentPolicyClient.BaseClient.Get(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var data struct {
		ApprovalStep []ApprovalStep `json:"value"`
	}
	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &data.ApprovalStep, status, nil
}

// Get will get an approval step
func (c *AccessPackageAssignmentApprovalStepClient) Get(ctx context.Context, id, approvalId string) (*ApprovalStep, int, error) {
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		ValidStatusCodes:       []int{http.StatusOK},
		Uri: Uri{
			Entity: fmt.Sprintf("/identityGovernance/entitlementManagement/accessPackageAssignmentApprovals/%s/steps/%s", approvalId, id),
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("AccessPackageAssignmentApprovalStepClient.BaseClient.Get(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var step ApprovalStep
	if err := json.Unmarshal(respBody, &step); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &step, status, nil
}

// Update amends an existing AccessPackageAssignmentPolicy.
func (c *AccessPackageAssignmentApprovalStepClient) Update(ctx context.Context, step ApprovalStep, approvalId string) (int, error) {
	var status int

	if step.ID == nil {
		return status, errors.New("cannot update ApprovalStep with nil ID")
	}

	body, err := json.Marshal(step)
	if err != nil {
		return status, fmt.Errorf("json.Marshal(): %v", err)
	}

	_, status, _, err = c.BaseClient.Put(ctx, PutHttpRequestInput{ //This is usually a patch but this endpoint uses PUT
		Body:                   body,
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		ValidStatusCodes:       []int{http.StatusOK},
		Uri: Uri{
			Entity: fmt.Sprintf("/identityGovernance/entitlementManagement/accessPackageAssignmentApprovals/%s/steps/%s", approvalId, *step.ID),
		},
	})
	if err != nil {
		return status, fmt.Errorf("AccessPackageAssignmentApprovalStepClient.BaseClient.Put(): %v", err)
	}

	return status, nil
}
