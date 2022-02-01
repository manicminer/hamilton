package containerregistry

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// CatalogList lists repositories
// Reference: https://docs.microsoft.com/en-us/rest/api/containerregistry/repository/get-list
// Parameters:
// last - parameter for the last item in previous query. Result set will include values lexically after last.
// maxItems - parameter for max number of item
//
// Parameters are pointers and if set to nil they won't be added to the query.
func (cr *ContainerRegistryClient) CatalogList(ctx context.Context, last *string, maxItems *int) ([]string, error) {
	baseURL, err := cr.getBaseURL()
	if err != nil {
		return nil, err
	}

	catalogURL, err := url.Parse(fmt.Sprintf("%s/acr/v1/_catalog", baseURL.String()))
	if err != nil {
		return nil, err
	}

	query := url.Values{}
	setQuery := false
	if last != nil {
		query.Set("last", *last)
		setQuery = true
	}

	if maxItems != nil {
		query.Set("n", fmt.Sprintf("%d", *maxItems))
		setQuery = true
	}

	if setQuery {
		catalogURL.RawQuery = query.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, catalogURL.String(), http.NoBody)
	if err != nil {
		return nil, err
	}

	err = cr.setAuthorizationHeader(ctx, req)
	if err != nil {
		return nil, err
	}

	res, err := cr.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 status code - %d: %s", res.StatusCode, string(resBytes))
	}

	var resData struct {
		Repositories []string `json:"repositories"`
	}
	err = json.Unmarshal(resBytes, &resData)
	if err != nil {
		return nil, err
	}

	return resData.Repositories, nil
}

// CatalogGetAttributes returns repository attributes for a specific image
// Reference: https://docs.microsoft.com/en-us/rest/api/containerregistry/repository/get-attributes
// Parameters:
// ctx - the context
// imageName - name of the image (including the namespace)
func (cr *ContainerRegistryClient) CatalogGetAttributes(ctx context.Context, imageName string) (RepositoryAttributesResponse, error) {
	baseURL, err := cr.getBaseURL()
	if err != nil {
		return RepositoryAttributesResponse{}, err
	}

	catalogURL, err := url.Parse(fmt.Sprintf("%s/acr/v1/%s", baseURL.String(), imageName))
	if err != nil {
		return RepositoryAttributesResponse{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, catalogURL.String(), http.NoBody)
	if err != nil {
		return RepositoryAttributesResponse{}, err
	}

	err = cr.setAuthorizationHeader(ctx, req)
	if err != nil {
		return RepositoryAttributesResponse{}, err
	}

	res, err := cr.httpClient.Do(req)
	if err != nil {
		return RepositoryAttributesResponse{}, err
	}

	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return RepositoryAttributesResponse{}, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return RepositoryAttributesResponse{}, fmt.Errorf("received non-200 status code - %d: %s", res.StatusCode, string(resBytes))
	}

	var resData RepositoryAttributesResponse
	err = json.Unmarshal(resBytes, &resData)
	if err != nil {
		return RepositoryAttributesResponse{}, err
	}

	return resData, nil
}

// CatalogUpdateAttributes updates the attribute identified by name where reference is the name of the repository
// Reference: https://docs.microsoft.com/en-us/rest/api/containerregistry/repository/update-attributes
// Parameters:
// ctx - the context
// imageName - name of the image (including the namespace)
// attributes - the attributes (that are non-nil) that should be updated
func (cr *ContainerRegistryClient) CatalogUpdateAttributes(ctx context.Context, imageName string, attributes RepositoryChangeableAttributes) error {
	bodyBytes, err := json.Marshal(attributes)
	if err != nil {
		return err
	}

	if string(bodyBytes) == "{}" {
		return fmt.Errorf("no attributes set")
	}

	baseURL, err := cr.getBaseURL()
	if err != nil {
		return err
	}

	catalogURL, err := url.Parse(fmt.Sprintf("%s/acr/v1/%s", baseURL.String(), imageName))
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, catalogURL.String(), bytes.NewReader(bodyBytes))
	if err != nil {
		return err
	}

	err = cr.setAuthorizationHeader(ctx, req)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := cr.httpClient.Do(req)
	if err != nil {
		return err
	}

	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 status code - %d: %s", res.StatusCode, string(resBytes))
	}

	return nil
}

// CatalogDelete the repository identified by name
// Reference: https://docs.microsoft.com/en-us/rest/api/containerregistry/repository/delete
// Parameters:
// ctx - the context
// imageName - name of the image (including the namespace)
func (cr *ContainerRegistryClient) CatalogDelete(ctx context.Context, imageName string) (RepositoryDeleteResponse, error) {
	baseURL, err := cr.getBaseURL()
	if err != nil {
		return RepositoryDeleteResponse{}, err
	}

	catalogURL, err := url.Parse(fmt.Sprintf("%s/acr/v1/%s", baseURL.String(), imageName))
	if err != nil {
		return RepositoryDeleteResponse{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, catalogURL.String(), http.NoBody)
	if err != nil {
		return RepositoryDeleteResponse{}, err
	}

	err = cr.setAuthorizationHeader(ctx, req)
	if err != nil {
		return RepositoryDeleteResponse{}, err
	}

	res, err := cr.httpClient.Do(req)
	if err != nil {
		return RepositoryDeleteResponse{}, err
	}

	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return RepositoryDeleteResponse{}, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusAccepted {
		return RepositoryDeleteResponse{}, fmt.Errorf("received non-202 status code - %d: %s", res.StatusCode, string(resBytes))
	}

	var resData RepositoryDeleteResponse
	err = json.Unmarshal(resBytes, &resData)
	if err != nil {
		return RepositoryDeleteResponse{}, err
	}

	return resData, nil
}

type RepositoryChangeableAttributes struct {
	DeleteEnabled   *bool `json:"deleteEnabled,omitempty"`
	WriteEnabled    *bool `json:"writeEnabled,omitempty"`
	ReadEnabled     *bool `json:"readEnabled,omitempty"`
	ListEnabled     *bool `json:"listEnabled,omitempty"`
	TeleportEnabled *bool `json:"teleportEnabled,omitempty"`
}

type RepositoryChangeableAttributesResponse struct {
	DeleteEnabled   bool `json:"deleteEnabled"`
	WriteEnabled    bool `json:"writeEnabled"`
	ReadEnabled     bool `json:"readEnabled"`
	ListEnabled     bool `json:"listEnabled"`
	TeleportEnabled bool `json:"teleportEnabled"`
}

type RepositoryAttributesResponse struct {
	Registry             string                                 `json:"registry"`
	ImageName            string                                 `json:"imageName"`
	CreatedTime          time.Time                              `json:"createdTime"`
	LastUpdateTime       time.Time                              `json:"lastUpdateTime"`
	ManifestCount        int                                    `json:"manifestCount"`
	TagCount             int                                    `json:"tagCount"`
	ChangeableAttributes RepositoryChangeableAttributesResponse `json:"changeableAttributes"`
}

type RepositoryDeleteResponse struct {
	ManifestsDeleted []string `json:"manifestsDeleted"`
	TagsDeleted      []string `json:"tagsDeleted"`
}
