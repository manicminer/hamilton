package msgraph

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/manicminer/hamilton/odata"
)

// AuthenticationClient performs operations on the Authentications methods endpoint under Identity and Sign-in
type AuthenticationClient struct {
	BaseClient Client
}

// NewAuthenticationClient returns a new ApplicationsClient
func NewAuthenticationClient(tenantId string) *AuthenticationClient {
	return &AuthenticationClient{
		BaseClient: NewClient(Version10, tenantId),
	}
}

//List all authentication methods
func (c *AuthenticationClient) List(ctx context.Context, query odata.Query) (*[]AuthenticationMethod, int, error) {
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		DisablePaging:    query.Top > 0,
		ValidStatusCodes: []int{http.StatusOK},
		Uri: Uri{
			Entity:      "/applications",
			Params:      query.Values(),
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("ApplicationsClient.BaseClient.Get(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var data struct {
		AuthMethods []AuthenticationMethod `json:"value"`
	}
	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &data.AuthMethods, status, nil
}
