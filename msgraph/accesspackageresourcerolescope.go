package msgraph

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/odata"
)

type AccessPackageResourceRoleScopeClient struct {
	BaseClient Client
}

func NewAccessPackageResourceRoleScopeClient(tenantId string) *AccessPackageResourceRoleScopeClient {
	return &AccessPackageResourceRoleScopeClient{
		BaseClient: NewClient(VersionBeta, tenantId),
	}
}

// List returns a list of AccessPackageResourceRoleScope(s)
func (c *AccessPackageResourceRoleScopeClient) List(ctx context.Context, query odata.Query, accessPackageId string) (*[]AccessPackageResourceRoleScope, int, error) {
	//Append Query with Expand required for endpoint
	query.Expand = odata.Expand{
		Relationship: "accessPackageResourceRoleScopes", Select: []string{"accessPackageResourceRole", "accessPackageResourceScope"},
	}

	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		//DisablePaging:    query.Top > 0, Query is not a collection
		ValidStatusCodes: []int{http.StatusOK},
		Uri: Uri{
			Entity:      fmt.Sprintf("/identityGovernance/entitlementManagement/accessPackages/%s", accessPackageId),
			Params:      query.Values(),
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("AccessPackageResourceRoleScopeClient.BaseClient.Get(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var data struct {
		AccessPackageResourceRoleScopes []AccessPackageResourceRoleScope `json:"value"`
	}
	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &data.AccessPackageResourceRoleScopes, status, nil
}

// Create creates a new AccessPackageResourceRoleScope.
func (c *AccessPackageResourceRoleScopeClient) Create(ctx context.Context, accessPackageResourceRoleScope AccessPackageResourceRoleScope) (*AccessPackageResourceRoleScope, int, error) {
	var status int
	body, err := json.Marshal(accessPackageResourceRoleScope)
	if err != nil {
		return nil, status, fmt.Errorf("json.Marshal(): %v", err)
	}

	resp, status, _, err := c.BaseClient.Post(ctx, PostHttpRequestInput{
		Body:             body,
		ValidStatusCodes: []int{http.StatusCreated},
		Uri: Uri{
			Entity:      fmt.Sprintf("/identityGovernance/entitlementManagement/accessPackages/%s/accessPackageResourceRoleScopes", *accessPackageResourceRoleScope.AccessPackageId),
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("AccessPackageResourceRoleScopeClient.BaseClient.Post(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var newAccessPackageResourceRoleScope AccessPackageResourceRoleScope
	if err := json.Unmarshal(respBody, &newAccessPackageResourceRoleScope); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	accessPackageResourceRoleScope.ID = newAccessPackageResourceRoleScope.ID //Keep the rest of the information just set the ID is it does not return whole request
	ids := strings.Split(*newAccessPackageResourceRoleScope.ID, "_")
	// We can derive the IDs of the AccessPackageResourceRole & AccessPackageResourceScope out of the combined _ ID returned back
	accessPackageResourceRoleScope.AccessPackageResourceRole.ID = utils.StringPtr(ids[0])
	accessPackageResourceRoleScope.AccessPackageResourceScope.ID = utils.StringPtr(ids[1])

	return &accessPackageResourceRoleScope, status, nil
}

// Get retrieves a AccessPackageResourceRoleScope.
func (c *AccessPackageResourceRoleScopeClient) Get(ctx context.Context, accessPackageId string, id string) (*AccessPackageResourceRoleScope, int, error) {
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		Uri: Uri{
			Entity: fmt.Sprintf("/identityGovernance/entitlementManagement/accessPackages/%s", accessPackageId),
			Params: odata.Query{
				Expand: odata.Expand{
					Relationship: "accessPackageResourceRoleScopes",
					Select:       []string{"accessPackageResourceRole", "accessPackageResourceScope"},
				},
				//Filter: fmt.Sprintf("startswith(originId,'%s')", id),
			}.Values(), //The Resource we made a request to add
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("AccessPackageResourceRoleScopeClient.BaseClient.Get(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var data struct {
		AccessPackageResourceRoleScopes []AccessPackageResourceRoleScope `json:"accessPackageResourceRoleScopes"`
	}
	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	var accessPackageResourceRoleScope AccessPackageResourceRoleScope
	// There is only a select and expand method on this endpoint - Instead foreach the list till we find it

	for _, roleScope := range data.AccessPackageResourceRoleScopes {
		if *roleScope.ID == id {
			accessPackageResourceRoleScope = roleScope
			accessPackageResourceRoleScope.AccessPackageId = &accessPackageId //Should probably also pass this back
		}
	}

	if accessPackageResourceRoleScope.ID == nil { //Throw error if we cant find it
		return nil, status, fmt.Errorf("AccessPackageResourceRoleScopeClient.BaseClient.Get(): Could not find accessPackageResourceRoleScope ID%v", err)
	}

	return &accessPackageResourceRoleScope, status, nil
}

// No Update Method

// No Delete Method
// Manage this downstream, intergrate into a accessPackage block and trigger a force replacement and recreate the AP & Role scope resources
