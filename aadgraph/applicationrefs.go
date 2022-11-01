package aadgraph

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/manicminer/hamilton/environments"
)

// ApplicationRefsClient performs operations on Applications.
type ApplicationRefsClient struct {
	BaseClient Client
}

// NewApplicationRefsClient returns a new ApplicationRefsClient
func NewApplicationRefsClient(tenantId string) *ApplicationRefsClient {
	return &ApplicationRefsClient{
		BaseClient: NewClient(Version20, tenantId),
	}
}

// Get retrieves an Application manifest.
func (c *ApplicationRefsClient) Get(ctx context.Context, id environments.ApiAppId) (*ApplicationRef, int, error) {
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		Uri: Uri{
			Entity: fmt.Sprintf("/applicationRefs/%s", id),
		},
	})
	if err != nil {
		return nil, status, err
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	var appRef ApplicationRef
	if err := json.Unmarshal(respBody, &appRef); err != nil {
		return nil, status, err
	}
	return &appRef, status, nil
}
