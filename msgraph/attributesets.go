package msgraph

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// AttributeSetsClient performs operations on Attribute Sets.
type AttributeSetsClient struct {
	BaseClient Client
}

// NewAttributeSetsClient returns a new AttributeSetsClient.
func NewAttributeSetsClient() *AttributeSetsClient {
	return &AttributeSetsClient{
		BaseClient: NewClient(Version10),
	}
}

// List returns a list of AttributeSet
func (c *AttributeSetsClient) List(ctx context.Context, query odata.Query) (*[]AttributeSet, int, error) {
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		DisablePaging:    query.Top > 0,
		OData:            query,
		ValidStatusCodes: []int{http.StatusOK},
		Uri: Uri{
			Entity: "/directory/attributeSets",
		},
	})

	if err != nil {
		return nil, status, fmt.Errorf("AttributeSetsClient.BaseClient.Get(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var data struct {
		AttributeSets []AttributeSet `json:"value"`
	}

	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &data.AttributeSets, status, nil
}

// Get retrieves an AttributeSet
func (c *AttributeSetsClient) Get(ctx context.Context, id string, query odata.Query) (*AttributeSet, int, error) {
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		OData:                  query,
		ValidStatusCodes:       []int{http.StatusOK},
		Uri: Uri{
			Entity: fmt.Sprintf("/directory/attributeSets/%s", id),
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("AttributeSetsClient.BaseClient.Get(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var attributeSet AttributeSet
	if err := json.Unmarshal(respBody, &attributeSet); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &attributeSet, status, nil
}

// Create creates an AttributeSet
func (c *AttributeSetsClient) Create(ctx context.Context, attributeSet AttributeSet) (*AttributeSet, int, error) {
	var status int

	body, err := json.Marshal(attributeSet)
	if err != nil {
		return nil, status, fmt.Errorf("json.Marshal(): %v", err)
	}

	resp, status, _, err := c.BaseClient.Post(ctx, PostHttpRequestInput{
		Body:                   body,
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		OData: odata.Query{
			Metadata: odata.MetadataFull,
		},
		ValidStatusCodes: []int{http.StatusCreated},
		Uri: Uri{
			Entity: "/directory/attributeSets",
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("AttributeSetsClient.BaseClient.Post(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var newAttributeSet AttributeSet
	if err := json.Unmarshal(respBody, &newAttributeSet); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &newAttributeSet, status, nil
}

// Update updates an AttributeSet
func (c *AttributeSetsClient) Update(ctx context.Context, attributeSet AttributeSet) (int, error) {
	var status int

	body, err := json.Marshal(attributeSet)
	if err != nil {
		return status, fmt.Errorf("json.Marshal(): %v", err)
	}

	_, status, _, err = c.BaseClient.Patch(ctx, PatchHttpRequestInput{
		Body:                   body,
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		ValidStatusCodes: []int{http.StatusNoContent},
		Uri: Uri{
			Entity: fmt.Sprintf("/directory/attributeSets/%s", *attributeSet.ID),
		},
	})
	if err != nil {
		return status, fmt.Errorf("AttributeSetsClient.BaseClient.Patch(): %v", err)
	}

	return status, nil
}
