package msgraph

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// UsersClient performs operations on Users.
type SignInLogsClient struct {
	BaseClient Client
}

// NewUsersClient returns a new UsersClient.
func NewSignInLogsClient(tenantId string) *SignInLogsClient {
	return &SignInLogsClient{
		BaseClient: NewClient(VersionBeta, tenantId),
	}
}

// List returns a list of Sign-in Logs, optionally filtered using OData.
func (c *SignInLogsClient) List(ctx context.Context, filter string) (*[]SignInLog, int, error) {
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
		SignInLogs []SignInLog `json:"value"`
	}
	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}
	return &data.SignInLogs, status, nil
}

// Get retrieves a SignInLog.
func (c *SignInLogsClient) Get(ctx context.Context, id string) (*SignInLog, int, error) {
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
	var signInLog SignInLog
	if err := json.Unmarshal(respBody, &signInLog); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}
	return &signInLog, status, nil
}
