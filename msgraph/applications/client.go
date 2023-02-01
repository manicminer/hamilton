package applications

import (
	"github.com/hashicorp/go-azure-sdk/sdk/client/msgraph"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// ApplicationsClient performs operations on Applications.
type ApplicationsClient struct {
	BaseClient *msgraph.Client
}

// NewApplicationsClient returns a new ApplicationsClient
func NewApplicationsClient(api environments.Api, apiVersion msgraph.ApiVersion, tenantId string) (*ApplicationsClient, error) {
	client, err := msgraph.NewMsGraphClient(api, apiVersion, tenantId)
	if err != nil {
		return nil, err
	}
	return &ApplicationsClient{
		BaseClient: client,
	}, nil
}
