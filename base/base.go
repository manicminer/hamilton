package base

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/manicminer/hamilton/auth"
)

const (
	DefaultEndpoint = "https://graph.microsoft.com"
	Version10       = "v1.0"
	VersionBeta     = "beta"
)

// ValidStatusFunc is a function that tests whether an HTTP response is considered valid for the particular request.
type ValidStatusFunc func(response *http.Response) bool

// HttpRequestInput is any type that can validate the response to an HTTP request.
type HttpRequestInput interface {
	GetValidStatusCodes() []int
	GetValidStatusFunc() ValidStatusFunc
}

// Uri represents a Microsoft Graph endpoint.
type Uri struct {
	Entity      string
	Params      url.Values
	HasTenantId bool
}

// GraphClient is any suitable HTTP client.
type GraphClient = *http.Client

// Client is a base client to be used by clients for specific entities.
// It can send GET, POST, PUT, PATCH and DELETE requests to Microsoft Graph and is API version and tenant aware.
type Client struct {
	// ApiVersion is the Microsoft Graph API version to use.
	ApiVersion string

	// Endpoint is the base endpoint for Microsoft Graph, usually "https://graph.microsoft.com".
	Endpoint string

	// TenantId is the tenant ID to use in requests.
	TenantId string

	// UserAgent is the HTTP user agent string to send in requests.
	UserAgent string

	// Authorizer is anything that can provide an access token with which to authorize requests.
	Authorizer auth.Authorizer

	httpClient GraphClient
}

// NewClient returns a new Client configured with the specified endpoint, tenant ID and API version.
func NewClient(endpoint, tenantId, apiVersion string) Client {
	return Client{
		httpClient: http.DefaultClient,
		Endpoint:   endpoint,
		TenantId:   tenantId,
		ApiVersion: apiVersion,
	}
}

// buildUri is used by the package to build a complete URI string for API requests.
func (c Client) buildUri(uri Uri) (string, error) {
	url, err := url.Parse(c.Endpoint)
	if err != nil {
		return "", err
	}
	url.Path = "/" + c.ApiVersion
	if uri.HasTenantId {
		url.Path = fmt.Sprintf("%s/%s", url.Path, c.TenantId)
	}
	url.Path = fmt.Sprintf("%s/%s", url.Path, strings.TrimLeft(uri.Entity, "/"))
	if uri.Params != nil {
		url.RawQuery = uri.Params.Encode()
	}
	return url.String(), nil
}

// performRequest is used by the package to send an HTTP request to the API.
func (c Client) performRequest(req *http.Request, input HttpRequestInput) (*http.Response, int, error) {
	var status int

	if c.Authorizer != nil {
		token, err := c.Authorizer.Token()
		if err != nil {
			return nil, status, err
		}
		token.SetAuthHeader(req)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	if c.UserAgent != "" {
		req.Header.Add("User-Agent", c.UserAgent)
	}

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

// containsStatusCode determines whether the returned status code is in the []int of expected status codes.
func containsStatusCode(expected []int, actual int) bool {
	for _, v := range expected {
		if actual == v {
			return true
		}
	}

	return false
}
