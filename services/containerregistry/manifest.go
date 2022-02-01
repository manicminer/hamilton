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

// ManifestList list manifests of a repository
// Reference: https://docs.microsoft.com/en-us/rest/api/containerregistry/manifests/get-list
// Parameters:
// ctx - the context
// imageName - name of the image (including the namespace)
// last - parameter for the last item in previous query. Result set will include values lexically after last.
// maxItems - parameter for max number of item
// orderBy - orderby query parameter
//
// Parameters are pointers and if set to nil they won't be added to the query.
func (cr *ContainerRegistryClient) ManifestList(ctx context.Context, imageName string, last *string, maxItems *int, orderBy *string) (ManifestList, error) {
	baseURL, err := cr.getBaseURL()
	if err != nil {
		return ManifestList{}, err
	}

	manifestURL, err := url.Parse(fmt.Sprintf("%s/acr/v1/%s/_manifests", baseURL.String(), imageName))
	if err != nil {
		return ManifestList{}, err
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

	if orderBy != nil {
		query.Set("orderby", *orderBy)
		setQuery = true
	}

	if setQuery {
		manifestURL.RawQuery = query.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, manifestURL.String(), http.NoBody)
	if err != nil {
		return ManifestList{}, err
	}

	err = cr.setAuthorizationHeader(ctx, req)
	if err != nil {
		return ManifestList{}, err
	}

	res, err := cr.httpClient.Do(req)
	if err != nil {
		return ManifestList{}, err
	}

	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return ManifestList{}, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return ManifestList{}, fmt.Errorf("received non-200 status code - %d: %s", res.StatusCode, string(resBytes))
	}

	var resData ManifestList
	err = json.Unmarshal(resBytes, &resData)
	if err != nil {
		return ManifestList{}, err
	}

	return resData, nil
}

// ManifestGet get the manifest identified by name and reference where reference can be a tag or digest.
// Reference: https://docs.microsoft.com/en-us/rest/api/containerregistry/manifests/get
// Parameters:
// ctx - the context
// imageName - name of the image (including the namespace)
// reference - a tag or a digest, pointing to a specific image
func (cr *ContainerRegistryClient) ManifestGet(ctx context.Context, imageName string, reference string) (ManifestGetResponse, error) {
	baseURL, err := cr.getBaseURL()
	if err != nil {
		return ManifestGetResponse{}, err
	}

	manifestURL, err := url.Parse(fmt.Sprintf("%s/v2/%s/manifests/%s", baseURL.String(), imageName, reference))
	if err != nil {
		return ManifestGetResponse{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, manifestURL.String(), http.NoBody)
	if err != nil {
		return ManifestGetResponse{}, err
	}

	req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json")

	err = cr.setAuthorizationHeader(ctx, req)
	if err != nil {
		return ManifestGetResponse{}, err
	}

	res, err := cr.httpClient.Do(req)
	if err != nil {
		return ManifestGetResponse{}, err
	}

	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return ManifestGetResponse{}, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return ManifestGetResponse{}, fmt.Errorf("received non-200 status code - %d: %s", res.StatusCode, string(resBytes))
	}

	var resData ManifestGetResponse
	err = json.Unmarshal(resBytes, &resData)
	if err != nil {
		return ManifestGetResponse{}, err
	}

	return resData, nil
}

// ManifestGetAttributes get manifest attributes
// Reference: https://docs.microsoft.com/en-us/rest/api/containerregistry/manifests/get-attributes
// Parameters:
// ctx - the context
// imageName - name of the image (including the namespace)
// reference - a tag or a digest, pointing to a specific image
func (cr *ContainerRegistryClient) ManifestGetAttributes(ctx context.Context, imageName string, reference string) (ManifestAttributesResponse, error) {
	baseURL, err := cr.getBaseURL()
	if err != nil {
		return ManifestAttributesResponse{}, err
	}

	manifestURL, err := url.Parse(fmt.Sprintf("%s/acr/v1/%s/_manifests/%s", baseURL.String(), imageName, reference))
	if err != nil {
		return ManifestAttributesResponse{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, manifestURL.String(), http.NoBody)
	if err != nil {
		return ManifestAttributesResponse{}, err
	}

	err = cr.setAuthorizationHeader(ctx, req)
	if err != nil {
		return ManifestAttributesResponse{}, err
	}

	res, err := cr.httpClient.Do(req)
	if err != nil {
		return ManifestAttributesResponse{}, err
	}

	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return ManifestAttributesResponse{}, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return ManifestAttributesResponse{}, fmt.Errorf("received non-200 status code - %d: %s", res.StatusCode, string(resBytes))
	}

	var resData ManifestAttributesResponse
	err = json.Unmarshal(resBytes, &resData)
	if err != nil {
		return ManifestAttributesResponse{}, err
	}

	return resData, nil
}

// ManifestUpdateAttributes update attributes of a manifest
// Reference: https://docs.microsoft.com/en-us/rest/api/containerregistry/manifests/update-attributes
// Parameters:
// ctx - the context
// imageName - name of the image (including the namespace)
// reference - a tag or a digest, pointing to a specific image
// attributes - the attributes (that are non-nil) that should be updated
func (cr *ContainerRegistryClient) ManifestUpdateAttributes(ctx context.Context, imageName string, reference string, attributes ManifestChangeableAttributes) error {
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

	manifestURL, err := url.Parse(fmt.Sprintf("%s/acr/v1/%s/_manifests/%s", baseURL.String(), imageName, reference))
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, manifestURL.String(), bytes.NewReader(bodyBytes))
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

// ManifestDelete delete the manifest identified by name and reference. Note that a manifest can only be deleted by digest.
// Reference: https://docs.microsoft.com/en-us/rest/api/containerregistry/tag/delete
// Parameters:
// ctx - the context
// imageName - name of the image (including the namespace)
// digest - a digest, pointing to a specific image
func (cr *ContainerRegistryClient) ManifestDelete(ctx context.Context, imageName string, digest string) error {
	baseURL, err := cr.getBaseURL()
	if err != nil {
		return err
	}

	manifestURL, err := url.Parse(fmt.Sprintf("%s/v2/%s/manifests/%s", baseURL.String(), imageName, digest))
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, manifestURL.String(), http.NoBody)
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

type ManifestChangeableAttributes struct {
	DeleteEnabled     *bool   `json:"deleteEnabled,omitempty"`
	ListEnabled       *bool   `json:"listEnabled,omitempty"`
	QuarantineDetails *string `json:"quarantineDetails,omitempty"`
	QuarantineState   *string `json:"quarantineState,omitempty"`
	ReadEnabled       *bool   `json:"readEnabled,omitempty"`
	WriteEnabled      *bool   `json:"writeEnabled,omitempty"`
}

type ManifestChangeableAttributesResponse struct {
	DeleteEnabled bool `json:"deleteEnabled"`
	WriteEnabled  bool `json:"writeEnabled"`
	ReadEnabled   bool `json:"readEnabled"`
	ListEnabled   bool `json:"listEnabled"`
}

type Manifest struct {
	Digest               string                               `json:"digest"`
	ImageSize            int64                                `json:"imageSize"`
	CreatedTime          time.Time                            `json:"createdTime"`
	LastUpdateTime       time.Time                            `json:"lastUpdateTime"`
	Architecture         string                               `json:"architecture"`
	Os                   string                               `json:"os"`
	MediaType            string                               `json:"mediaType"`
	ConfigMediaType      string                               `json:"configMediaType"`
	Tags                 []string                             `json:"tags"`
	ChangeableAttributes ManifestChangeableAttributesResponse `json:"changeableAttributes"`
}

type ManifestList struct {
	Registry  string     `json:"registry"`
	ImageName string     `json:"imageName"`
	Manifests []Manifest `json:"manifests"`
}

type ManifestAttributesResponse struct {
	Registry  string   `json:"registry"`
	ImageName string   `json:"imageName"`
	Manifest  Manifest `json:"manifest"`
}

type ManifestConfigResponse struct {
	MediaType string `json:"mediaType"`
	Size      int64  `json:"size"`
	Digest    string `json:"digest"`
}

type ManifestLayerResponse struct {
	MediaType string `json:"mediaType"`
	Size      int64  `json:"size"`
	Digest    string `json:"digest"`
}

type ManifestGetResponse struct {
	SchemaVersion int                     `json:"schemaVersion"`
	MediaType     string                  `json:"mediaType"`
	Config        ManifestConfigResponse  `json:"config"`
	Layers        []ManifestLayerResponse `json:"layers"`
}
