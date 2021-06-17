package msgraph

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// SignInReports Client performs operations on Sign in reports.
type SignInReportsClient struct {
	BaseClient Client
}

// NewSignInLogsClient returns a new SignInReportsClient.
func NewSignInLogsClient(tenantId string) *SignInReportsClient {
	return &SignInReportsClient{
		BaseClient: NewClient(VersionBeta, tenantId),
	}
}

// List returns a list of Sign-in Reports, optionally filtered using OData.
func (c *SignInReportsClient) List(ctx context.Context, filter string) (*[]SignInReport, int, error) {
	params := url.Values{}
	if filter != "" {
		params.Add("$filter", filter)
	}
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		Uri: Uri{
			Entity:      "/auditLogs/signIns",
			Params:      params,
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
	var data struct {
		SignInLogs []SignInReport `json:"value"`
	}
	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}
	return &data.SignInLogs, status, nil
}

// Get retrieves a Sign-in Report.
func (c *SignInReportsClient) Get(ctx context.Context, id string) (*SignInReport, int, error) {
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		ValidStatusCodes:       []int{http.StatusOK},
		Uri: Uri{
			Entity:      fmt.Sprintf("/auditLogs/signIns/%s", id),
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
	var signInReport SignInReport
	if err := json.Unmarshal(respBody, &signInReport); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}
	return &signInReport, status, nil
}
