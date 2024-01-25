package msgraph

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
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

	Value json.RawMessage `json:"value"`

	InternalError1 json.RawMessage `json:"error"`
	InternalError2 json.RawMessage `json:"odata.error"`
}

func (o *OData) UnmarshalJSON(data []byte) error {
	// Unmarshal using a local type
	type od OData
	var o2 od
	if err := json.Unmarshal(data, &o2); err != nil {
		return err
	}
	*o = OData(o2)

	// Look for errors in the "error" and "odata.error" fields, unmarshal separately if any
	if o.InternalError1 != nil {
		var e odata.Error
		if err := json.Unmarshal(o.InternalError1, &e); err != nil {
			return err
		}
		o.Error = &e
	}
	if o.InternalError2 != nil {
		var e odata.Error
		if err := json.Unmarshal(o.InternalError2, &e); err != nil {
			return err
		}
		o.Error = &e
	}

	return nil
}

// FasterFromResponse parses an http.Response and returns:
// - unmarshaled OData (if no odata is present in the response, or the content type is invalid, returns nil)
// - unmarshaled result (as pointer to the specified resultType, normally a slice of structs)
// - link to the next page (if present)
func FasterFromResponse(resp *http.Response, resultType reflect.Type) (*OData, interface{}, *odata.Link, error) {
	if resp == nil {
		return nil, nil, nil, nil
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
			return nil, nil, nil, fmt.Errorf("could not read response body: %s", err)
		}

		// Unmarshal odata
		if err := json.Unmarshal(respBody, &o); err != nil {
			return nil, nil, nil, err
		}

		// Unmarshal result
		result := reflect.New(resultType).Interface()
		if err := json.Unmarshal(o.Value, result); err != nil {
			return nil, nil, nil, err
		}

		return &o, result, o.NextLink, nil
	}

	return nil, nil, nil, nil
}

// fasterPerformRequest is used by the package to send an HTTP request to the API.
func (c Client) fasterPerformRequest(req *http.Request, input FasterGetHttpRequestInput, resultType reflect.Type) (int, interface{}, *odata.Link, error) {
	var status int

	query := input.GetOData()
	req.Header = query.AppendHeaders(req.Header)
	req.Header.Add("Content-Type", input.GetContentType())

	if c.Authorizer != nil {
		token, err := c.Authorizer.Token(req.Context(), req)
		if err != nil {
			return status, nil, nil, err
		}
		token.SetAuthHeader(req)
	}

	if c.UserAgent != "" {
		req.Header.Add("User-Agent", c.UserAgent)
	}

	var resp *http.Response
	var o *OData
	var result interface{}
	var nextLink *odata.Link
	var err error

	var reqBody []byte
	if req.Body != nil {
		reqBody, err = io.ReadAll(req.Body)
		if err != nil {
			return status, nil, nil, fmt.Errorf("reading request body: %v", err)
		}
	}

	c.RetryableClient.CheckRetry = func(ctx context.Context, resp *http.Response, err error) (bool, error) {
		if resp != nil && !c.DisableRetries {
			if resp.StatusCode == http.StatusFailedDependency {
				return true, nil
			}

			o, result, nextLink, err = FasterFromResponse(resp, resultType)
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
				return status, nil, nil, err
			}
			req = r
		}
	}

	resp, err = c.HttpClient.Do(req)
	if err != nil {
		return status, nil, nil, err
	}
	if resp == nil {
		return status, nil, nil, fmt.Errorf("nil response received")
	}

	if c.ResponseMiddlewares != nil {
		for _, m := range *c.ResponseMiddlewares {
			r, err := m(req, resp)
			if err != nil {
				return status, nil, nil, err
			}
			resp = r
		}
	}

	status = resp.StatusCode
	if !containsStatusCode(input.GetValidStatusCodes(), status) {
		f := input.GetValidStatusFunc()
		if f != nil && f(resp, o) {
			return status, result, nextLink, nil
		}

		var errText string
		switch {
		case o != nil && o.Error != nil && o.Error.String() != "":
			errText = fmt.Sprintf("OData error: %s", o.Error)
		default:
			defer resp.Body.Close()
			respBody, err := io.ReadAll(resp.Body)
			if err != nil {
				return status, result, nil, fmt.Errorf("unexpected status %d, could not read response body", status)
			}
			if len(respBody) == 0 {
				return status, result, nil, fmt.Errorf("unexpected status %d received with no body", status)
			}
			errText = fmt.Sprintf("response: %s", respBody)
		}
		return status, result, nil, fmt.Errorf("unexpected status %d with %s", status, errText)
	}

	return status, result, nextLink, nil
}

// FasterGet performs a GET request.
func (c Client) FasterGet(ctx context.Context, input FasterGetHttpRequestInput, result interface{}) (int, error) {
	var status int

	// Check for a raw uri, else build one from the Uri field
	url := input.rawUri
	if url == "" {
		// Append odata query parameters
		input.Uri.Params = input.OData.AppendValues(input.Uri.Params)

		var err error
		url, err = c.buildUri(input.Uri)
		if err != nil {
			return status, fmt.Errorf("unable to make request: %v", err)
		}
	}

	// Build a new request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return status, err
	}

	// Perform the request
	status, partial, nextLink, err := c.fasterPerformRequest(req, input, reflect.TypeOf(result).Elem())
	if err != nil {
		return status, err
	}

	// Append the partial result to the result
	if reflect.ValueOf(partial).IsValid() && !reflect.ValueOf(partial).IsZero() {
		// equivalent of *result = append(*result, partial)
		reflect.ValueOf(result).Elem().Set(reflect.AppendSlice(reflect.ValueOf(result).Elem(), reflect.ValueOf(partial).Elem()))
	}

	if input.DisablePaging || nextLink == nil || partial == nil {
		// No more pages, reassign result and return
		return status, nil
	}

	// Get the next page, recursively
	nextInput := input
	nextInput.rawUri = string(*nextLink)
	return c.FasterGet(ctx, nextInput, result)
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
