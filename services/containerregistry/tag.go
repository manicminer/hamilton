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

// TagList lists tags of a repository
// Reference: https://docs.microsoft.com/en-us/rest/api/containerregistry/tag/get-list
// Parameters:
// ctx - the context
// imageName - name of the image (including the namespace)
// digest - filter by digest
// last - parameter for the last item in previous query. Result set will include values lexically after last.
// maxItems - parameter for max number of item
// orderBy - orderby query parameter
//
// Parameters are pointers and if set to nil they won't be added to the query.
func (cr *ContainerRegistryClient) TagList(ctx context.Context, imageName string, digest *string, last *string, maxItems *int, orderBy *string) (TagList, error) {
	baseURL, err := cr.getBaseURL()
	if err != nil {
		return TagList{}, err
	}

	tagURL, err := url.Parse(fmt.Sprintf("%s/acr/v1/%s/_tags", baseURL.String(), imageName))
	if err != nil {
		return TagList{}, err
	}

	query := url.Values{}
	setQuery := false
	if digest != nil {
		query.Set("digest", *digest)
		setQuery = true
	}

	if last != nil {
		query.Set("last", *last)
		setQuery = true
	}

	if maxItems != nil {
		query.Set("n", fmt.Sprintf("%d", *maxItems))
		setQuery = true
	}

	if orderBy != nil {
		query.Set("orderby", *orderBy)
		setQuery = true
	}

	if setQuery {
		tagURL.RawQuery = query.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, tagURL.String(), http.NoBody)
	if err != nil {
		return TagList{}, err
	}

	err = cr.setAuthorizationHeader(ctx, req)
	if err != nil {
		return TagList{}, err
	}

	res, err := cr.httpClient.Do(req)
	if err != nil {
		return TagList{}, err
	}

	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return TagList{}, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return TagList{}, fmt.Errorf("received non-200 status code - %d: %s", res.StatusCode, string(resBytes))
	}

	var resData TagList
	err = json.Unmarshal(resBytes, &resData)
	if err != nil {
		return TagList{}, err
	}

	return resData, nil
}

// TagGetAttributes get tag attributes by tag
// Reference: https://docs.microsoft.com/en-us/rest/api/containerregistry/tag/get-attributes
// Parameters:
// ctx - the context
// imageName - name of the image (including the namespace)
// reference - tag name
func (cr *ContainerRegistryClient) TagGetAttributes(ctx context.Context, imageName string, reference string) (TagAttributesResponse, error) {
	baseURL, err := cr.getBaseURL()
	if err != nil {
		return TagAttributesResponse{}, err
	}

	tagURL, err := url.Parse(fmt.Sprintf("%s/acr/v1/%s/_tags/%s", baseURL.String(), imageName, reference))
	if err != nil {
		return TagAttributesResponse{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, tagURL.String(), http.NoBody)
	if err != nil {
		return TagAttributesResponse{}, err
	}

	err = cr.setAuthorizationHeader(ctx, req)
	if err != nil {
		return TagAttributesResponse{}, err
	}

	res, err := cr.httpClient.Do(req)
	if err != nil {
		return TagAttributesResponse{}, err
	}

	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return TagAttributesResponse{}, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return TagAttributesResponse{}, fmt.Errorf("received non-200 status code - %d: %s", res.StatusCode, string(resBytes))
	}

	var resData TagAttributesResponse
	err = json.Unmarshal(resBytes, &resData)
	if err != nil {
		return TagAttributesResponse{}, err
	}

	return resData, nil
}

// TagUpdateAttributes update tag attributes
// Reference: https://docs.microsoft.com/en-us/rest/api/containerregistry/tag/update-attributes
// Parameters:
// ctx - the context
// imageName - name of the image (including the namespace)
// reference - tag name
// attributes - the attributes (that are non-nil) that should be updated
func (cr *ContainerRegistryClient) TagUpdateAttributes(ctx context.Context, imageName string, reference string, attributes TagChangeableAttributes) error {
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

	tagURL, err := url.Parse(fmt.Sprintf("%s/acr/v1/%s/_tags/%s", baseURL.String(), imageName, reference))
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, tagURL.String(), bytes.NewReader(bodyBytes))
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

// TagDelete delete tag
// Reference: https://docs.microsoft.com/en-us/rest/api/containerregistry/tag/delete
// Parameters:
// ctx - the context
// imageName - name of the image (including the namespace)
// reference - tag name
func (cr *ContainerRegistryClient) TagDelete(ctx context.Context, imageName string, reference string) error {
	baseURL, err := cr.getBaseURL()
	if err != nil {
		return err
	}

	tagURL, err := url.Parse(fmt.Sprintf("%s/acr/v1/%s/_tags/%s", baseURL.String(), imageName, reference))
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, tagURL.String(), http.NoBody)
	if err != nil {
		return err
	}

	err = cr.setAuthorizationHeader(ctx, req)
	if err != nil {
		return err
	}

	res, err := cr.httpClient.Do(req)
	if err != nil {
		return err
	}

	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusAccepted {
		return fmt.Errorf("received non-202 status code - %d: %s", res.StatusCode, string(resBytes))
	}

	return nil
}

type TagChangeableAttributes struct {
	DeleteEnabled *bool `json:"deleteEnabled,omitempty"`
	WriteEnabled  *bool `json:"writeEnabled,omitempty"`
	ReadEnabled   *bool `json:"readEnabled,omitempty"`
	ListEnabled   *bool `json:"listEnabled,omitempty"`
}

type TagChangeableAttributesResponse struct {
	DeleteEnabled bool `json:"deleteEnabled"`
	WriteEnabled  bool `json:"writeEnabled"`
	ReadEnabled   bool `json:"readEnabled"`
	ListEnabled   bool `json:"listEnabled"`
}

type Tag struct {
	Name                 string                          `json:"name"`
	Digest               string                          `json:"digest"`
	CreatedTime          time.Time                       `json:"createdTime"`
	LastUpdateTime       time.Time                       `json:"lastUpdateTime"`
	Signed               bool                            `json:"signed"`
	ChangeableAttributes TagChangeableAttributesResponse `json:"changeableAttributes"`
}

type TagList struct {
	Registry  string `json:"registry"`
	ImageName string `json:"imageName"`
	Tags      []Tag  `json:"tags"`
}

type TagAttributesResponse struct {
	Registry  string `json:"registry"`
	ImageName string `json:"imageName"`
	Tag       Tag    `json:"tag"`
}
