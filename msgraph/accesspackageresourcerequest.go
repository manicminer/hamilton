package msgraph

import (
	"context"
	"encoding/json"
	"github.com/manicminer/hamilton/internal/utils"
	"fmt"
	"io"
	"net/http"
	"log"

	"github.com/manicminer/hamilton/odata"
)

type AccessPackageResourceRequestClient struct {
	BaseClient Client
}

func NewAccessPackageResourceRequestClient(tenantId string) *AccessPackageResourceRequestClient {
	return &AccessPackageResourceRequestClient{
		BaseClient: NewClient(VersionBeta, tenantId),
	}
}

// List returns a list of AccessPackageResourceRequest
func (c *AccessPackageResourceRequestClient) List(ctx context.Context, query odata.Query) (*[]AccessPackageResourceRequest, int, error) {
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		DisablePaging:    query.Top > 0,
		ValidStatusCodes: []int{http.StatusOK},
		Uri: Uri{
			Entity:      "/identityGovernance/entitlementManagement/accessPackageResourceRequests",
			Params:      query.Values(),
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("AccessPackageResourceRequestClient.BaseClient.Get(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var data struct {
		AccessPackageResourceRequests []AccessPackageResourceRequest `json:"value"`
	}
	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &data.AccessPackageResourceRequests, status, nil
}

// Create creates a new AccessPackageResourceRequest.
func (c *AccessPackageResourceRequestClient) Create(ctx context.Context, accessPackageResourceRequest AccessPackageResourceRequest) (*AccessPackageResourceRequest, int, error) {
	var status int
	body, err := json.Marshal(accessPackageResourceRequest)
	jsonOutput := string(body[:])
	log.Printf("Json: %s", jsonOutput)
	if err != nil {
		return nil, status, fmt.Errorf("json.Marshal(): %v", err)
	}

	 ResourceDoesNotExist := func(resp *http.Response, o *odata.OData) bool {
	 	return o != nil && o.Error != nil && o.Error.Match(odata.ErrorResourceDoesNotExist)
	 }

	resp, status, _, err := c.BaseClient.Post(ctx, PostHttpRequestInput{
		Body:             body,
		ConsistencyFailureFunc: ResourceDoesNotExist,  // There appears to be potential backend lag when creating new groups, AP then adding resources. For stability wait here for 15 seconds before continuing
		ValidStatusCodes: []int{http.StatusCreated},
		Uri: Uri{
			Entity:      "/identityGovernance/entitlementManagement/accessPackageResourceRequests",
			HasTenantId: true,
		},
	})

	if err != nil {
		return nil, status, fmt.Errorf("AccessPackageResourceRequestClient.BaseClient.Post(): %v ", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var newAccessPackageResourceRequest AccessPackageResourceRequest
	if err := json.Unmarshal(respBody, &newAccessPackageResourceRequest); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	// The endpoint does not actually return the AccessPackageResources created which makes implementation impossible, workaround make sure to tag responses that 201 back with what we've just posted

	newAccessPackageResourceRequest.AccessPackageResource = accessPackageResourceRequest.AccessPackageResource

	return &newAccessPackageResourceRequest, status, nil
}

//Get retrieves a AccessPackageResourceRequest.
//Pseudo Get via Odata - Note this downstream
func (c *AccessPackageResourceRequestClient) Get(ctx context.Context, id string) (*AccessPackageResourceRequest, int, error) {
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		ValidStatusCodes:       []int{http.StatusOK},
		Uri: Uri{
			Entity:      "/identityGovernance/entitlementManagement/accessPackageResourceRequests",
			Params:      odata.Query{Filter: fmt.Sprintf("startswith(id,'%s')", id)}.Values(),
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("AccessPackageResourceRequestClient.BaseClient.Get(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var accessPackageResourceRequest AccessPackageResourceRequest
	if err := json.Unmarshal(respBody, &accessPackageResourceRequest); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &accessPackageResourceRequest, status, nil
}

// Update amends an existing AccessPackageResourceRequest.

//There is no endpoint for this, implement force replacement in downstream providers

// Delete removes a AccessPackageResourceRequest - Pseudo Delete
// There is no Delete Endpoint - Instead we create a request to delete a resource assignment. This is weird behavior but can be implemented, specific calls are needed
// Implement downstream via Resources endpoint: See how this is preformed in Tests
//https://docs.microsoft.com/en-us/graph/api/accesspackageresourcerequest-post?view=graph-rest-beta&tabs=http#example-5-create-an-accesspackageresourcerequest-for-removing-a-resource
func (c *AccessPackageResourceRequestClient) Delete(ctx context.Context, accessPackageResourceRequest AccessPackageResourceRequest, accessPackageResource AccessPackageResource) (int, error) {

	var status int

	// Deletion Request based off resource that can be found out via List endpoint on APResources. Pass in APRequest Resources Origin ID w/ Catalog To find the resource ID
	// See tests for examples
	newaccessPackageResourceRequest := AccessPackageResourceRequest{
		CatalogId: accessPackageResourceRequest.CatalogId,
		RequestType: utils.StringPtr("AdminRemove"),
		AccessPackageResource: &AccessPackageResource{
			ID: accessPackageResource.ID,
		},
	}
	
	body, err := json.Marshal(newaccessPackageResourceRequest)
	if err != nil {
		return status, fmt.Errorf("json.Marshal(): %v", err)
	}

	resp, status, _, err := c.BaseClient.Post(ctx, PostHttpRequestInput{
		Body:             body,
		ValidStatusCodes: []int{http.StatusCreated},
		Uri: Uri{
			Entity:      "/identityGovernance/entitlementManagement/accessPackageResourceRequests",
			HasTenantId: true,
		},
	})
	if err != nil {
		return status, fmt.Errorf("AccessPackageResourceRequestClient.BaseClient.Post(): %v", err)
	}

	defer resp.Body.Close()
	_, err2 := io.ReadAll(resp.Body)
	if err2 != nil {
		return status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	return status, nil //Expecting 201 Back
}
