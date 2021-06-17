package msgraph

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// DirectoryAuditReportsClient performs operations on directory Audit reports.
type DirectoryAuditReportsClient struct {
	BaseClient Client
}

// NewDirectoryAuditReportsClient returns a new DirectoryAuditReportsClient.
func NewDirectoryAuditReportsClient(tenantId string) *DirectoryAuditReportsClient {
	return &DirectoryAuditReportsClient{
		BaseClient: NewClient(VersionBeta, tenantId),
	}
}

// List returns a list of Directory audit report logs, optionally filtered using OData.
func (c *DirectoryAuditReportsClient) List(ctx context.Context, filter string) (*[]AuditLog, int, error) {
	params := url.Values{}
	if filter != "" {
		params.Add("$filter", filter)
	}
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		Uri: Uri{
			Entity:      "/auditLogs/directoryAudits",
			Params:      params,
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("DirectoryAuditReportsClient.BaseClient.Get(): %v", err)
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("ioutil.ReadAll(): %v", err)
	}
	var data struct {
		DirectoryAuditReports []AuditLog `json:"value"`
	}
	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}
	return &data.DirectoryAuditReports, status, nil
}

// Get retrieves a Directory audit report.
func (c *DirectoryAuditReportsClient) Get(ctx context.Context, id string) (*AuditLog, int, error) {
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		ValidStatusCodes:       []int{http.StatusOK},
		Uri: Uri{
			Entity:      fmt.Sprintf("/auditLogs/directoryAudits/%s", id),
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("SignInLogsClient.BaseClient.Get(): %v", err)
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("ioutil.ReadAll(): %v", err)
	}
	var directoryAuditReport AuditLog
	if err := json.Unmarshal(respBody, &directoryAuditReport); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}
	return &directoryAuditReport, status, nil
}
