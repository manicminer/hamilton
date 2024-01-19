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
)

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
	resp, status, o, err := c.performRequest(req, input)
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
