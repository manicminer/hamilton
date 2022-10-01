package msgraph

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/manicminer/hamilton/odata"
)

// B2CUserFlowClient performs operations on B2CUserFlow.
type B2CUserFlowClient struct {
	BaseClient Client
}

// NewB2CUserFlowClient returns a new B2CUserFlowClient.
func NewB2CUserFlowClient(tenantId string) *B2CUserFlowClient {
	return &B2CUserFlowClient{
		BaseClient: NewClient(VersionBeta, tenantId),
	}
}

// List returns a list of B2C UserFlows, optionally queried using OData.
func (c *B2CUserFlowClient) List(ctx context.Context, query odata.Query) (*[]B2CUserFlow, int, error) {
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		OData:            query,
		ValidStatusCodes: []int{http.StatusOK},
		Uri: Uri{
			Entity:      "/identity/b2cUserFlows",
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("B2CUserFlowClient.BaseClient.Get(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var data struct {
		UserFlows []B2CUserFlow `json:"value"`
	}
	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &data.UserFlows, status, nil
}

// Create creates a new B2CUserFlow.
func (c *B2CUserFlowClient) Create(ctx context.Context, userflow B2CUserFlow) (*B2CUserFlow, int, error) {
	var status int

	body, err := json.Marshal(userflow)
	if err != nil {
		return nil, status, fmt.Errorf("json.Marshal(): %v", err)
	}

	resp, status, _, err := c.BaseClient.Post(ctx, PostHttpRequestInput{
		Body: body,
		OData: odata.Query{
			Metadata: odata.MetadataFull,
		},
		ValidStatusCodes: []int{http.StatusCreated},
		Uri: Uri{
			Entity:      "/identity/b2cUserFlows",
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("B2CUserFlowClient.BaseClient.Post(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var newUserFlow B2CUserFlow
	if err := json.Unmarshal(respBody, &newUserFlow); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &newUserFlow, status, nil
}

// Get returns an existing B2CUserFlow.
func (c *B2CUserFlowClient) Get(ctx context.Context, id string, query odata.Query) (*B2CUserFlow, int, error) {
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		OData:                  query,
		ValidStatusCodes:       []int{http.StatusOK},
		Uri: Uri{
			Entity:      fmt.Sprintf("/identity/b2cUserFlows/%s", id),
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("B2CUserFlowClient.BaseClient.Get(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var userflow B2CUserFlow
	if err := json.Unmarshal(respBody, &userflow); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &userflow, status, nil
}

// Update amends an existing B2CUserFlow.
func (c *B2CUserFlowClient) Update(ctx context.Context, userflow B2CUserFlow) (int, error) {
	var status int
	if userflow.ID == nil {
		return status, fmt.Errorf("cannot update userflow with nil ID")
	}

	userflowID := *userflow.ID
	userflow.ID = nil

	body, err := json.Marshal(userflow)
	if err != nil {
		return status, fmt.Errorf("json.Marshal(): %v", err)
	}

	_, status, _, err = c.BaseClient.Patch(ctx, PatchHttpRequestInput{
		Body:                   body,
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		ValidStatusCodes: []int{
			http.StatusOK,
			http.StatusNoContent,
		},
		Uri: Uri{
			Entity:      fmt.Sprintf("/identity/b2cUserFlows//%s", userflowID),
			HasTenantId: true,
		},
	})
	if err != nil {
		return status, fmt.Errorf("B2CUserFlowClient.BaseClient.Patch(): %v", err)
	}

	return status, nil
}

// Delete removes a B2CUserFlow.
func (c *B2CUserFlowClient) Delete(ctx context.Context, id string) (int, error) {
	_, status, _, err := c.BaseClient.Delete(ctx, DeleteHttpRequestInput{
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		ValidStatusCodes:       []int{http.StatusNoContent},
		Uri: Uri{
			Entity:      fmt.Sprintf("/identity/b2cUserFlows/%s", id),
			HasTenantId: true,
		},
	})
	if err != nil {
		return status, fmt.Errorf("B2CUserFlowClient.BaseClient.Delete(): %v", err)
	}

	return status, nil
}

// AssignAttribute assigns the user flow attribute to user flow.
func (c *B2CUserFlowClient) AssignAttribute(ctx context.Context, userflowId string, assignment UserFlowAttributeAssignment) (*UserFlowAttributeAssignment, int, error) {
	var status int

	if assignment.UserAttribute == nil {
		return nil, status, fmt.Errorf("UserAttribute cannot be nil")
	}
	if assignment.UserAttribute.ID == nil {
		return nil, status, fmt.Errorf("UserAttribute.ID cannot be nil")
	}
	body, err := json.Marshal(assignment)
	if err != nil {
		return nil, status, fmt.Errorf("json.Marshal(): %v", err)
	}
	resp, status, _, err := c.BaseClient.Post(ctx, PostHttpRequestInput{
		Body: body,
		OData: odata.Query{
			Metadata: odata.MetadataFull,
		},
		ValidStatusCodes: []int{http.StatusCreated},
		Uri: Uri{
			Entity: fmt.Sprintf("/identity/b2cUserFlows/%s/userAttributeAssignments", userflowId),
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("UserFlowAttributeAssignmentsClient.BaseClient.Post(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var newAttrAssignment UserFlowAttributeAssignment
	if err := json.Unmarshal(respBody, &newAttrAssignment); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &newAttrAssignment, status, nil
}

// ListAssignedAttributes returns all the assigned attributes for the user flow.
func (c *B2CUserFlowClient) ListAssignedAttributes(ctx context.Context, id string, query odata.Query) (*[]UserFlowAttributeAssignment, int, error) {
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		OData:            query,
		ValidStatusCodes: []int{http.StatusOK},
		Uri: Uri{
			Entity: fmt.Sprintf("/identity/b2cUserFlows/%s/userAttributeAssignments", id)},
	})
	if err != nil {
		return nil, status, fmt.Errorf("UserFlowAttributeAssignmentsClient.BaseClient.Get(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var data struct {
		UserFlowAttributeAssignments []UserFlowAttributeAssignment `json:"value"`
	}
	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &data.UserFlowAttributeAssignments, status, nil
}

// GetAssignedAttribute returns the assigned attribute.
func (c *B2CUserFlowClient) GetAssignedAttribute(ctx context.Context, userflowId, assignmentId string) (*UserFlowAttributeAssignment, int, error) {
	var status int

	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		OData: odata.Query{
			Metadata: odata.MetadataFull,
		},
		ValidStatusCodes: []int{http.StatusCreated},
		Uri: Uri{
			Entity: fmt.Sprintf("/identity/b2cUserFlows/%s/userAttributeAssignments/%s", userflowId, assignmentId),
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("UserFlowAttributeAssignmentsClient.BaseClient.Get(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var newAttrAssignment UserFlowAttributeAssignment
	if err := json.Unmarshal(respBody, &newAttrAssignment); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &newAttrAssignment, status, nil
}

// RemoveAttributeAssignment removes the assigned attribute from the user flow.
func (c *B2CUserFlowClient) RemoveAttributeAssignment(ctx context.Context, userflowId, assignmentId string) (int, error) {
	var status int

	_, status, _, err := c.BaseClient.Delete(ctx, DeleteHttpRequestInput{
		OData: odata.Query{
			Metadata: odata.MetadataFull,
		},
		ValidStatusCodes: []int{http.StatusCreated},
		Uri: Uri{
			Entity: fmt.Sprintf("/identity/b2cUserFlows/%s/userAttributeAssignments/%s", userflowId, assignmentId),
		},
	})
	if err != nil {
		return status, fmt.Errorf("UserFlowAttributeAssignmentsClient.BaseClient.Delete(): %v", err)
	}

	return status, nil
}

// UpdateAttributeAssignment updates the user flow attribute assignment.
func (c *B2CUserFlowClient) UpdateAttributeAssignment(ctx context.Context, userflowId string, assignment UserFlowAttributeAssignment) (*UserFlowAttributeAssignment, int, error) {
	var status int

	if assignment.ID == nil {
		return nil, status, fmt.Errorf("Assignment ID cannot be nil")
	}
	// The API doesn't allow updating the user attribute in an assignment. The user should create
	// a new assignment if they wish to change the attribute.
	if assignment.UserAttribute != nil {
		return nil, status, fmt.Errorf("UserAttribute should be nil")
	}
	body, err := json.Marshal(assignment)
	if err != nil {
		return nil, status, fmt.Errorf("json.Marshal(): %v", err)
	}
	id := *assignment.ID
	resp, status, _, err := c.BaseClient.Patch(ctx, PatchHttpRequestInput{
		Body: body,
		OData: odata.Query{
			Metadata: odata.MetadataFull,
		},
		ValidStatusCodes: []int{http.StatusCreated},
		Uri: Uri{
			Entity: fmt.Sprintf("/identity/b2cUserFlows/%s/userAttributeAssignments/%s", userflowId, id),
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("UserFlowAttributeAssignmentsClient.BaseClient.Patch(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var newAttrAssignment UserFlowAttributeAssignment
	if err := json.Unmarshal(respBody, &newAttrAssignment); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &newAttrAssignment, status, nil
}
