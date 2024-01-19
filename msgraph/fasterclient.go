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

// OData is used to unmarshall OData metadata from an API response.
type OData struct {
	Context      *string     `json:"@odata.context"`
	MetadataEtag *string     `json:"@odata.metadataEtag"`
	Type         *odata.Type `json:"@odata.type"`
	Count        *int        `json:"@odata.count"`
	NextLink     *odata.Link `json:"@odata.nextLink"`
	Delta        *string     `json:"@odata.delta"`
	DeltaLink    *odata.Link `json:"@odata.deltaLink"`
	Id           *odata.Id   `json:"@odata.id"`
	EditLink     *odata.Link `json:"@odata.editLink"`
	Etag         *string     `json:"@odata.etag"`

	Error *odata.Error `json:"-"`

	Value interface{} `json:"value"`
}

func (o *OData) UnmarshalJSON(data []byte) error {
	// Unmarshal using a local type
	type od OData
	var o2 od
	if err := json.Unmarshal(data, &o2); err != nil {
		return err
	}
	*o = OData(o2)

	// Look for errors in the "error" and "odata.error" fields
	var e map[string]json.RawMessage
	if err := json.Unmarshal(data, &e); err != nil {
		return err
	}
	for _, k := range []string{"error", "odata.error"} {
		if v, ok := e[k]; ok {
			var e2 odata.Error
			if err := json.Unmarshal(v, &e2); err != nil {
				return err
			}
			o.Error = &e2
			break
		}
	}
	return nil
}

// FasterFromResponse parses an http.Response and returns an unmarshaled OData
// If no odata is present in the response, or the content type is invalid, returns nil
func FasterFromResponse(resp *http.Response) (*OData, error) {
	if resp == nil {
		return nil, nil
	}

	var o OData

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
func (c Client) fasterPerformRequest(req *http.Request, input FasterGetHttpRequestInput) (*http.Response, int, *OData, error) {
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
	var o *OData
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
func (c Client) FasterGet(ctx context.Context, input FasterGetHttpRequestInput) (*http.Response, int, *OData, error) {
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

// FasterConsistencyFailureFunc is a function that determines whether an HTTP request has failed due to eventual consistency and should be retried
type FasterConsistencyFailureFunc func(*http.Response, *OData) bool

// FasterValidStatusFunc is a function that tests whether an HTTP response is considered valid for the particular request.
type FasterValidStatusFunc func(*http.Response, *OData) bool

// FasterRetryOn404ConsistencyFailureFunc can be used to retry a request when a 404 response is received
func FasterRetryOn404ConsistencyFailureFunc(resp *http.Response, _ *OData) bool {
	return resp != nil && resp.StatusCode == http.StatusNotFound
}

// FasterGetHttpRequestInput configures a GET request.
type FasterGetHttpRequestInput struct {
	ConsistencyFailureFunc FasterConsistencyFailureFunc
	DisablePaging          bool
	OData                  odata.Query
	ValidStatusCodes       []int
	ValidStatusFunc        FasterValidStatusFunc
	Uri                    Uri
	rawUri                 string
}

// GetConsistencyFailureFunc returns a function used to evaluate whether a failed request is due to eventual consistency and should be retried.
func (i FasterGetHttpRequestInput) GetConsistencyFailureFunc() FasterConsistencyFailureFunc {
	return i.ConsistencyFailureFunc
}

// GetContentType returns the content type for the request, currently only application/json is supported
func (i FasterGetHttpRequestInput) GetContentType() string {
	return "application/json; charset=utf-8"
}

// GetOData returns the OData request metadata
func (i FasterGetHttpRequestInput) GetOData() odata.Query {
	return i.OData
}

// GetValidStatusCodes returns a []int of status codes considered valid for a GET request.
func (i FasterGetHttpRequestInput) GetValidStatusCodes() []int {
	return i.ValidStatusCodes
}

// GetValidStatusFunc returns a function used to evaluate whether the response to a GET request is considered valid.
func (i FasterGetHttpRequestInput) GetValidStatusFunc() FasterValidStatusFunc {
	return i.ValidStatusFunc
}
