package msgraph

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// DirectoryRoleTemplatesClient performs operations on DirectoryRoleTemplates.
type DirectoryRoleTemplatesClient struct {
	BaseClient Client
}

// NewDirectoryRoleTemplatesClient returns a new DirectoryRoleTemplatesClient
func NewDirectoryRoleTemplatesClient(tenantId string) *DirectoryRoleTemplatesClient {
	return &DirectoryRoleTemplatesClient{
		BaseClient: NewClient(Version10, tenantId),
	}
}

// List returns a list of DirectoryRoleTemplates, optionally filtered using OData.
func (c *DirectoryRoleTemplatesClient) List(ctx context.Context) (*[]DirectoryRoleTemplate, int, error) {

	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		Uri: Uri{
			Entity:      "/directoryRoleTemplates",
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, err
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	var data struct {
		DirectoryRoleTemplates []DirectoryRoleTemplate `json:"value"`
	}
	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, status, err
	}
	return &data.DirectoryRoleTemplates, status, nil
}

// Get retrieves an DirectoryRoleTemplates manifest.
func (c *DirectoryRoleTemplatesClient) Get(ctx context.Context, id string) (*DirectoryRoleTemplate, int, error) {
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		Uri: Uri{
			Entity:      fmt.Sprintf("/directoryRoleTemplates/%s", id),
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, err
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	var dirRoleTemplate DirectoryRoleTemplate
	if err := json.Unmarshal(respBody, &dirRoleTemplate); err != nil {
		return nil, status, err
	}
	return &dirRoleTemplate, status, nil
}