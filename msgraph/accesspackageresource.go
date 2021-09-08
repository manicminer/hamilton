package msgraph

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"github.com/manicminer/hamilton/odata"
)

type AccessPackageResourceClient struct {
	BaseClient Client
}

func NewAccessPackageResourceClient(tenantId string) *AccessPackageResourceClient {
	return &AccessPackageResourceClient{
		BaseClient: NewClient(VersionBeta, tenantId),
	}
}

// List Method
func (c *AccessPackageResourceClient) List(ctx context.Context, query odata.Query, catalogId string) (*[]AccessPackageResource, int, error) {
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		DisablePaging:    query.Top > 0,
		ValidStatusCodes: []int{http.StatusOK},
		Uri: Uri{
			Entity:      fmt.Sprintf("/identityGovernance/entitlementManagement/accessPackageCatalogs/%s/accessPackageResources", catalogId),
			Params:      query.Values(),
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("AccessPackageResourceClient.BaseClient.Get(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var data struct {
		AccessPackageResources []AccessPackageResource `json:"value"`
	}
	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &data.AccessPackageResources, status, nil
}

//Pseudo Get Method via OData

func (c *AccessPackageResourceClient) Get(ctx context.Context, catalogId string, originId string) (*AccessPackageResource, int, error) {
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		ValidStatusCodes: []int{http.StatusOK},
		Uri: Uri{
			Entity:      fmt.Sprintf("/identityGovernance/entitlementManagement/accessPackageCatalogs/%s/accessPackageResources", catalogId),
			Params:      odata.Query{Filter: fmt.Sprintf("startswith(originId,'%s')", originId)}.Values(),
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("AccessPackageResourceClient.BaseClient.Get(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}
	var data struct {
	AccessPackageResources []AccessPackageResource `json:"value"`
	}

	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	accessPackageResource := data.AccessPackageResources[0]

	return &accessPackageResource, status, nil
}