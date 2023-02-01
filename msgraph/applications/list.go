package applications

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/manicminer/hamilton/msgraph"
)

type ListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]msgraph.Application
}

type ListOptions struct {
	OData odata.Query
}

func (o ListOptions) ToHeaders() *client.Headers {
	return nil
}

func (o ListOptions) ToOData() *odata.Query {
	return &o.OData
}

func (o ListOptions) ToQuery() *client.QueryParams {
	return nil
}

// List returns a list of Applications, optionally queried using OData.
func (c *ApplicationsClient) List(ctx context.Context, options ListOptions) (*ListOperationResponse, error) {
	req, err := c.BaseClient.NewRequest(ctx, client.RequestOptions{
		ContentType: "application/json",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Path:          "/applications",
	})
	if err != nil {
		return nil, err
	}

	var resp *client.Response
	if options.OData.Top > 0 {
		resp, err = req.Execute(ctx)
	} else {
		resp, err = req.ExecutePaged(ctx)
	}

	result := ListOperationResponse{}
	if resp != nil {
		result.OData = resp.OData
		result.HttpResponse = resp.Response
	}
	if err != nil {
		return nil, fmt.Errorf("ApplicationsClient.List(): %v", err)
	}

	var value struct {
		Value *[]msgraph.Application `json:"value"`
	}
	if err = resp.Unmarshal(&value); err != nil {
		return nil, err
	}

	result.Model = value.Value
	return &result, nil
}
