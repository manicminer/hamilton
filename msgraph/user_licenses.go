package msgraph

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type AssignLicenseResourceClient struct {
	BaseClient Client
}

func NewAssignLicenseResourceClient() *AssignLicenseResourceClient {
	return &AssignLicenseResourceClient{
		BaseClient: NewClient(Version10),
	}
}

// List returns a list of Users, optionally queried using OData.
func (c *AssignLicenseResourceClient) ListUserLicenses(ctx context.Context, query odata.Query, id string) (*[]AddLicenses, int, error) {
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		DisablePaging:    query.Top > 0,
		OData:            query,
		ValidStatusCodes: []int{http.StatusOK},
		Uri: Uri{
			Entity: fmt.Sprintf("/users/%s/licenseDetails", id),
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("AssignLicenseResourceClient.BaseClient.Get(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var data struct {
		Licenses []AddLicenses `json:"value"`
	}
	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &data.Licenses, status, nil
}

// Create creates a new User.
// Assign can accept userPrincipalName but List cannot. Restricting input to id.
func (c *AssignLicenseResourceClient) Create(ctx context.Context, id string, license AssignLicenseRequest) (*User, int, error) {
	var status int

	body, err := json.Marshal(license)
	if err != nil {
		return nil, status, fmt.Errorf("json.Marshal(): %v", err)
	}

	consistencyFunc := func(resp *http.Response, o *odata.OData) bool {
		if resp != nil && resp.StatusCode == http.StatusBadRequest && o != nil && o.Error != nil {
			return o.Error.Match(odata.ErrorPropertyValuesAreInvalid)
		}
		return false
	}

	resp, status, _, err := c.BaseClient.Post(ctx, PostHttpRequestInput{
		Body:                   body,
		ConsistencyFailureFunc: consistencyFunc,
		OData: odata.Query{
			Metadata: odata.MetadataFull,
		},
		ValidStatusCodes: []int{http.StatusCreated},
		Uri: Uri{
			Entity: fmt.Sprintf("/users/%s/assignLicense", id),
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("AssignLicenseResourceClient.BaseClient.Post(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var user User
	if err := json.Unmarshal(respBody, &license); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &user, status, nil
}

// GetDeleted retrieves a deleted User.
func (c *AssignLicenseResourceClient) GetDeleted(ctx context.Context, id string, query odata.Query) (*User, int, error) {
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		OData:                  query,
		ValidStatusCodes:       []int{http.StatusOK},
		Uri: Uri{
			Entity: fmt.Sprintf("/directory/deletedItems/%s", id),
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("AssignLicenseResourceClient.BaseClient.Get(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var user User
	if err := json.Unmarshal(respBody, &user); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &user, status, nil
}

// Delete removes a User.
func (c *AssignLicenseResourceClient) Delete(ctx context.Context, id string) (int, error) {
	_, status, _, err := c.BaseClient.Delete(ctx, DeleteHttpRequestInput{
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		ValidStatusCodes:       []int{http.StatusNoContent},
		Uri: Uri{
			Entity: fmt.Sprintf("/users/%s", id),
		},
	})
	if err != nil {
		return status, fmt.Errorf("AssignLicenseResourceClient.BaseClient.Delete(): %v", err)
	}

	return status, nil
}

// ListDeleted retrieves a list of recently deleted users, optionally queried using OData.
func (c *AssignLicenseResourceClient) ListDeleted(ctx context.Context, query odata.Query) (*[]User, int, error) {
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		DisablePaging:    query.Top > 0,
		OData:            query,
		ValidStatusCodes: []int{http.StatusOK},
		Uri: Uri{
			Entity: "/directory/deleteditems/microsoft.graph.user",
		},
	})
	if err != nil {
		return nil, status, err
	}

	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	var data struct {
		DeletedUsers []User `json:"value"`
	}
	if err = json.Unmarshal(respBody, &data); err != nil {
		return nil, status, err
	}

	return &data.DeletedUsers, status, nil
}
