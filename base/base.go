package base

import (
	"bytes"
	"context"
	"fmt"
	"github.com/manicminer/hamilton/auth"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	DefaultEndpoint = "https://graph.microsoft.com"
	Version10       = "v1.0"
	VersionBeta     = "beta"
)

type GraphClient = *http.Client

type BaseClient struct {
	ApiVersion string
	Endpoint   string
	TenantId   string

	authorizer auth.Authorizer
	httpClient GraphClient
}

func NewBaseClient(authorizer auth.Authorizer, endpoint, tenantId, version string) BaseClient {
	return BaseClient{
		authorizer: authorizer,
		httpClient: http.DefaultClient,
		Endpoint:   endpoint,
		TenantId:   tenantId,
		ApiVersion: version,
	}
}

type DeleteHttpRequestInput struct {
	ValidStatusCodes []int
	Uri              string
}

func (c BaseClient) Delete(ctx context.Context, input DeleteHttpRequestInput) (*http.Response, error) {
	url := c.buildUri(input.Uri)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, http.NoBody)
	if err != nil {
		return nil, err
	}
	resp, err := c.performRequest(ctx, req, input.ValidStatusCodes)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

type GetHttpRequestInput struct {
	ValidStatusCodes []int
	Uri              string
}

func (c BaseClient) Get(ctx context.Context, input GetHttpRequestInput) (*http.Response, error) {
	url := c.buildUri(input.Uri)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, err
	}
	resp, err := c.performRequest(ctx, req, input.ValidStatusCodes)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

type PostHttpRequestInput struct {
	Body             []byte
	ValidStatusCodes []int
	Uri              string
}

func (c BaseClient) Post(ctx context.Context, input PostHttpRequestInput) (*http.Response, error) {
	url := c.buildUri(input.Uri)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(input.Body))
	if err != nil {
		return nil, err
	}
	resp, err := c.performRequest(ctx, req, input.ValidStatusCodes)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

type PutHttpRequestInput struct {
	Body             []byte
	ValidStatusCodes []int
	Uri              string
}

func (c BaseClient) Put(ctx context.Context, input PutHttpRequestInput) (*http.Response, error) {
	url := c.buildUri(input.Uri)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewBuffer(input.Body))
	if err != nil {
		return nil, err
	}
	resp, err := c.performRequest(ctx, req, input.ValidStatusCodes)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

type PatchHttpRequestInput struct {
	Body             []byte
	ValidStatusCodes []int
	Uri              string
}

func (c BaseClient) Patch(ctx context.Context, input PatchHttpRequestInput) (*http.Response, error) {
	url := c.buildUri(input.Uri)
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, bytes.NewBuffer(input.Body))
	if err != nil {
		return nil, err
	}
	resp, err := c.performRequest(ctx, req, input.ValidStatusCodes)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c BaseClient) buildUri(uri string) string {
	return fmt.Sprintf("%s/%s/%s/%s", c.Endpoint, c.ApiVersion, c.TenantId, strings.TrimLeft(uri, "/"))
}

func (c BaseClient) performRequest(_ context.Context, req *http.Request, validStatusCodes []int) (*http.Response, error) {
	token, err := c.authorizer.Token()
	if err != nil {
		return nil, err
	}
	token.SetAuthHeader(req)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if !containsStatusCode(validStatusCodes, resp.StatusCode) {
		defer resp.Body.Close()
		respBody, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status %d with response: %s", resp.StatusCode, string(respBody))
	}
	return resp, nil
}

func containsStatusCode(expected []int, actual int) bool {
	for _, v := range expected {
		if actual == v {
			return true
		}
	}

	return false
}
