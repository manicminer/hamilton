package containerregistry

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type CatalogClient struct {
	mu                 sync.Mutex
	cr                 *ContainerRegistryClient
	refreshToken       string
	refreshTokenClaims RefreshTokenClaims
	accessToken        string
	accessTokenClaims  AccessTokenClaims
}

func NewCatalogClient(cr *ContainerRegistryClient) *CatalogClient {
	return &CatalogClient{
		cr: cr,
	}
}

// List repositories
// Reference: https://docs.microsoft.com/en-us/rest/api/containerregistry/repository/get-list
// Parameters:
// last - parameter for the last item in previous query. Result set will include values lexically after last.
// maxItems - parameter for max number of item
//
// Parameters are pointers and if set to nil they won't be added to the query.
func (c *CatalogClient) List(ctx context.Context, last *string, maxItems *int) ([]string, error) {
	baseURL, err := c.cr.getBaseURL()
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

	err = c.setAuthorizationHeader(ctx, req)
	if err != nil {
		return nil, err
	}

	res, err := c.cr.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	var resData struct {
		Repositories []string `json:"repositories"`
	}

	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 status code - %d: %s", res.StatusCode, string(resBytes))
	}

	err = json.Unmarshal(resBytes, &resData)
	if err != nil {
		return nil, err
	}

	return resData.Repositories, nil
}

// GetAttributes returns repository attributes for a specific image
// Reference: https://docs.microsoft.com/en-us/rest/api/containerregistry/repository/get-attributes
// Parameters:
// ctx - the context
// imageName - name of the image (including the namespace)
func (c *CatalogClient) GetAttributes(ctx context.Context, imageName string) (RepositoryAttributesResponse, error) {
	baseURL, err := c.cr.getBaseURL()
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

	err = c.setAuthorizationHeader(ctx, req)
	if err != nil {
		return RepositoryAttributesResponse{}, err
	}

	res, err := c.cr.httpClient.Do(req)
	if err != nil {
		return RepositoryAttributesResponse{}, err
	}

	var resData RepositoryAttributesResponse

	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return RepositoryAttributesResponse{}, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return RepositoryAttributesResponse{}, fmt.Errorf("received non-200 status code - %d: %s", res.StatusCode, string(resBytes))
	}

	err = json.Unmarshal(resBytes, &resData)
	if err != nil {
		return RepositoryAttributesResponse{}, err
	}

	return resData, nil
}

// UpdateAttributes updates the attribute identified by name where reference is the name of the repository
// Parameters:
// ctx - the context
// imageName - name of the image (including the namespace)
// attributes - the attributes (that are non-nil) that should be updated
func (c *CatalogClient) UpdateAttributes(ctx context.Context, imageName string, attributes RepositoryChangeableAttributes) error {
	bodyBytes, err := json.Marshal(attributes)
	if err != nil {
		return err
	}

	if string(bodyBytes) == "{}" {
		return fmt.Errorf("no attributes set")
	}

	baseURL, err := c.cr.getBaseURL()
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

	err = c.setAuthorizationHeader(ctx, req)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := c.cr.httpClient.Do(req)
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

// Delete the repository identified by name
// Parameters:
// ctx - the context
// imageName - name of the image (including the namespace)
func (c *CatalogClient) Delete(ctx context.Context, imageName string) (RepositoryDeleteResponse, error) {
	return RepositoryDeleteResponse{}, fmt.Errorf("Delete not implemented yet")
}

func (c *CatalogClient) setAuthorizationHeader(ctx context.Context, req *http.Request) error {
	at, err := c.getAccessToken(ctx)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", at))

	return nil
}

func (c *CatalogClient) getAccessToken(ctx context.Context) (string, error) {
	c.mu.Lock()
	at := c.accessToken
	atClaims := c.accessTokenClaims
	c.mu.Unlock()

	atExpiry := time.Unix(atClaims.ExpirationTime, 0)
	if at != "" && atExpiry.After(time.Now().Add(time.Minute)) {
		return at, nil
	}

	c.mu.Lock()
	rt := c.refreshToken
	rtClaims := c.refreshTokenClaims
	c.mu.Unlock()

	rtExpiry := time.Unix(rtClaims.ExpirationTime, 0)
	if rt != "" && !rtExpiry.IsZero() && rtExpiry.After(time.Now().Add(time.Minute)) {
		scopes := AccessTokenScopes{
			{
				Type:    "registry",
				Name:    "catalog",
				Actions: []string{"*"},
			},
			{
				Type:    "repository",
				Name:    "*",
				Actions: []string{"*"},
			},
		}

		at, atClaims, err := c.cr.ExchangeAccessToken(ctx, c.refreshToken, scopes)
		if err != nil {
			return "", nil
		}

		c.mu.Lock()
		c.accessToken = at
		c.accessTokenClaims = atClaims
		c.mu.Unlock()

		return c.getAccessToken(ctx)
	}

	rt, rtClaims, err := c.cr.ExchangeRefreshToken(ctx)
	if err != nil {
		return "", err
	}

	c.mu.Lock()
	c.refreshToken = rt
	c.refreshTokenClaims = rtClaims
	c.mu.Unlock()

	return c.getAccessToken(ctx)
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
