package containerregistry

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// V2Check tells whether this Docker Registry instance supports Docker Registry HTTP API v2.
// Reference: https://docs.microsoft.com/en-us/rest/api/containerregistry/v2-support/check
// Parameters:
// ctx - the context
func (cr *ContainerRegistryClient) V2Check(ctx context.Context) (bool, error) {
	baseURL, err := cr.getBaseURL()
	if err != nil {
		return false, err
	}

	v2URL, err := url.Parse(fmt.Sprintf("%s/v2/", baseURL.String()))
	if err != nil {
		return false, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, v2URL.String(), http.NoBody)
	if err != nil {
		return false, err
	}

	err = cr.setAuthorizationHeader(ctx, req)
	if err != nil {
		return false, err
	}

	res, err := cr.httpClient.Do(req)
	if err != nil {
		return false, err
	}

	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return false, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return false, fmt.Errorf("received non-200 status code - %d: %s", res.StatusCode, string(resBytes))
	}

	return true, nil
}
