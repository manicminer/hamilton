package base

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/manicminer/hamilton/auth"
)

const (
	DefaultEndpoint = "https://graph.microsoft.com"
	Version10       = "v1.0"
	VersionBeta     = "beta"
)

type ValidStatusFunc func(response *http.Response) bool

type HttpRequestInput interface {
	GetValidStatusCodes() []int
	GetValidStatusFunc() ValidStatusFunc
}

type GraphClient = *http.Client

type Client struct {
	ApiVersion string
	Endpoint   string
	TenantId   string

	authorizer auth.Authorizer
	httpClient GraphClient
}

func NewClient(authorizer auth.Authorizer, endpoint, tenantId, version string) Client {
	return Client{
		authorizer: authorizer,
		httpClient: http.DefaultClient,
		Endpoint:   endpoint,
		TenantId:   tenantId,
		ApiVersion: version,
	}
}

func (c Client) buildUri(uri string) string {
	return fmt.Sprintf("%s/%s/%s/%s", c.Endpoint, c.ApiVersion, c.TenantId, strings.TrimLeft(uri, "/"))
}

func (c Client) performRequest(_ context.Context, req *http.Request, input HttpRequestInput) (*http.Response, int, error) {
	var status int

	token, err := c.authorizer.Token()
	if err != nil {
		return nil, status, err
	}

	token.SetAuthHeader(req)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, status, err
	}

	status = resp.StatusCode
	if !containsStatusCode(input.GetValidStatusCodes(), status) {
		f := input.GetValidStatusFunc()
		if f != nil && f(resp) {
			return resp, status, nil
		}

		defer resp.Body.Close()
		respBody, _ := ioutil.ReadAll(resp.Body)
		return nil, status, fmt.Errorf("unexpected status %d with response: %s", resp.StatusCode, string(respBody))
	}

	return resp, status, nil
}

func containsStatusCode(expected []int, actual int) bool {
	for _, v := range expected {
		if actual == v {
			return true
		}
	}

	return false
}
