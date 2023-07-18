package msgraph

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type AccessPackageAssignmentApprovalClient struct {
	BaseClient Client
}

func NewAccessPackageAssignmentApprovalClient() *AccessPackageAssignmentApprovalClient {
	return &AccessPackageAssignmentApprovalClient{
		BaseClient: NewClient(Version10),
	}
}

// Get will get an approval
func (c *AccessPackageAssignmentApprovalClient) Get(ctx context.Context, id string) (*Approval, int, error) {
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		ValidStatusCodes:       []int{http.StatusOK},
		Uri: Uri{
			Entity: fmt.Sprintf("/identityGovernance/entitlementManagement/accessPackageAssignmentApprovals/%s", id),
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("AccessPackageAssignmentApprovalClient.BaseClient.Get(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var approval Approval
	if err := json.Unmarshal(respBody, &approval); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &approval, status, nil
}
