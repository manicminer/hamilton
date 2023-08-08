package msgraph

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type PrivilegedAccessGroupClient struct {
	BaseClient Client
}

func NewPrivilegedAccessGroupClient() *PrivilegedAccessGroupClient {
	return &PrivilegedAccessGroupClient{
		BaseClient: NewClient(VersionBeta),
	}
}

// Register enables a group for PIM.
// Technically this is part of the deprecated v2 PIM API, however there is no equivilent in v3
// and automatic enablement doesn't work for service principals with assigned permissions.
func (c *PrivilegedAccessGroupClient) Register(ctx context.Context, groupId string) (int, error) {
	// We get an unknown 401 error if the group has not fully provisioned before
	// attempting to enable PIM.
	consistencyFunc := func(resp *http.Response, o *odata.OData) bool {
		return o != nil && o.Error != nil && resp.StatusCode == http.StatusUnauthorized
	}

	_, status, _, err := c.BaseClient.Post(ctx, PostHttpRequestInput{
		ConsistencyFailureFunc: consistencyFunc,
		Body:                   []byte(fmt.Sprintf("{\"externalId\": \"%s\"}", groupId)),
		ValidStatusCodes:       []int{http.StatusOK},
		Uri: Uri{
			Entity: "/privilegedAccess/aadGroups/resources/register",
		},
	})
	if err != nil {
		return status, fmt.Errorf("PrivilegedAccessGroupClient.BaseClient.Post(): %v", err)
	}

	return status, nil
}
