package msgraph

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/hashicorp/go-retryablehttp"
)

// FasterFromResponse parses an http.Response and returns an unmarshaled OData
// If no odata is present in the response, or the content type is invalid, returns nil
func FasterFromResponse(resp *http.Response) (*odata.OData, error) {
	if resp == nil {
		return nil, nil
	}

	var o odata.OData

	// Check for json content before looking for odata metadata
	contentType := strings.ToLower(resp.Header.Get("Content-Type"))
	if strings.HasPrefix(contentType, "application/json") {
		// Read the response body and close it
		respBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		// Always reassign the response body
		resp.Body = io.NopCloser(bytes.NewBuffer(respBody))

		if err != nil {
			return nil, fmt.Errorf("could not read response body: %s", err)
		}

		// Unmarshal odata
		if err := json.Unmarshal(respBody, &o); err != nil {
			return nil, err
		}

		return &o, nil
	}

	return nil, nil
}

// fasterPerformRequest is used by the package to send an HTTP request to the API.
func (c Client) fasterPerformRequest(req *http.Request, input HttpRequestInput) (*http.Response, int, *odata.OData, error) {
	var status int

	query := input.GetOData()
	req.Header = query.AppendHeaders(req.Header)
	req.Header.Add("Content-Type", input.GetContentType())

	if c.Authorizer != nil {
		token, err := c.Authorizer.Token(req.Context(), req)
		if err != nil {
			return nil, status, nil, err
		}
		token.SetAuthHeader(req)
	}

	if c.UserAgent != "" {
		req.Header.Add("User-Agent", c.UserAgent)
	}

	var resp *http.Response
	var o *odata.OData
	var err error

	var reqBody []byte
	if req.Body != nil {
		reqBody, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, status, nil, fmt.Errorf("reading request body: %v", err)
		}
	}

	c.RetryableClient.CheckRetry = func(ctx context.Context, resp *http.Response, err error) (bool, error) {
		if resp != nil && !c.DisableRetries {
			if resp.StatusCode == http.StatusFailedDependency {
				return true, nil
			}

			o, err = FasterFromResponse(resp)
			if err != nil {
				return false, err
			}

			f := input.GetConsistencyFailureFunc()
			if f != nil && f(resp, o) {
				return true, nil
			}
		}
		return retryablehttp.DefaultRetryPolicy(ctx, resp, err)
	}

	req.Body = io.NopCloser(bytes.NewBuffer(reqBody))

	if c.RequestMiddlewares != nil {
		for _, m := range *c.RequestMiddlewares {
			r, err := m(req)
			if err != nil {
				return nil, status, nil, err
			}
			req = r
		}
	}

	resp, err = c.HttpClient.Do(req)
	if err != nil {
		return nil, status, nil, err
	}

	if c.ResponseMiddlewares != nil {
		for _, m := range *c.ResponseMiddlewares {
			r, err := m(req, resp)
			if err != nil {
				return nil, status, nil, err
			}
			resp = r
		}
	}

	o, err = FasterFromResponse(resp)
	if err != nil {
		return nil, status, o, err
	}
	if resp == nil {
		return resp, status, o, fmt.Errorf("nil response received")
	}

	status = resp.StatusCode
	if !containsStatusCode(input.GetValidStatusCodes(), status) {
		f := input.GetValidStatusFunc()
		if f != nil && f(resp, o) {
			return resp, status, o, nil
		}

		var errText string
		switch {
		case o != nil && o.Error != nil && o.Error.String() != "":
			errText = fmt.Sprintf("OData error: %s", o.Error)
		default:
			defer resp.Body.Close()
			respBody, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, status, o, fmt.Errorf("unexpected status %d, could not read response body", status)
			}
			if len(respBody) == 0 {
				return nil, status, o, fmt.Errorf("unexpected status %d received with no body", status)
			}
			errText = fmt.Sprintf("response: %s", respBody)
		}
		return nil, status, o, fmt.Errorf("unexpected status %d with %s", status, errText)
	}

	return resp, status, o, nil
}

// FasterGet performs a GET request.
func (c Client) FasterGet(ctx context.Context, input GetHttpRequestInput) (*http.Response, int, *odata.OData, error) {
	var status int

	// Check for a raw uri, else build one from the Uri field
	url := input.rawUri
	if url == "" {
		// Append odata query parameters
		input.Uri.Params = input.OData.AppendValues(input.Uri.Params)

		var err error
		url, err = c.buildUri(input.Uri)
		if err != nil {
			return nil, status, nil, fmt.Errorf("unable to make request: %v", err)
		}
	}

	// Build a new request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, status, nil, err
	}

	// Perform the request
	resp, status, o, err := c.fasterPerformRequest(req, input)
	if err != nil {
		return nil, status, o, err
	}

	// Check for json content before handling pagination
	contentType := strings.ToLower(resp.Header.Get("Content-Type"))
	if strings.HasPrefix(contentType, "application/json") {
		// Read the response body and close it
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, status, o, fmt.Errorf("could not parse response body")
		}
		resp.Body.Close()

		// Unmarshall firstOdata
		var firstOdata odata.OData
		if err := json.Unmarshal(respBody, &firstOdata); err != nil {
			return nil, status, o, err
		}

		firstValue, ok := firstOdata.Value.([]interface{})
		if input.DisablePaging || firstOdata.NextLink == nil || firstValue == nil || !ok {
			// No more pages, reassign response body and return
			resp.Body = io.NopCloser(bytes.NewBuffer(respBody))
			return resp, status, o, nil
		}

		// Get the next page, recursively
		nextInput := input
		nextInput.rawUri = string(*firstOdata.NextLink)
		nextResp, status, o, err := c.FasterGet(ctx, nextInput)
		if err != nil {
			return resp, status, o, err
		}

		// Read the next page response body and close it
		nextRespBody, err := io.ReadAll(nextResp.Body)
		if err != nil {
			return nil, status, o, fmt.Errorf("could not parse response body")
		}
		nextResp.Body.Close()

		// Unmarshall firstOdata from the next page
		var nextOdata odata.OData
		if err := json.Unmarshal(nextRespBody, &nextOdata); err != nil {
			return resp, status, o, err
		}

		if nextValue, ok := nextOdata.Value.([]interface{}); ok {
			// Next page has results, append to current page
			value := append(firstValue, nextValue...)
			nextOdata.Value = &value
		}

		// Marshal the entire result, along with fields from the final page
		newJson, err := json.Marshal(nextOdata)
		if err != nil {
			return resp, status, o, err
		}

		// Reassign the response body
		resp.Body = io.NopCloser(bytes.NewBuffer(newJson))
	}

	return resp, status, o, nil
}
