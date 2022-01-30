package containerregistry

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/manicminer/hamilton/auth"
)

// ContainerRegistryClient handles communication with Azure Container Registry
type ContainerRegistryClient struct {
	authorizer auth.Authorizer
	httpClient *http.Client
	serverURL  string
	tenantID   string
}

// NewContainerRegistryClient returns a ContainerRegistryClient
func NewContainerRegistryClient(authorizer auth.Authorizer, serverURL string, tenantID string) *ContainerRegistryClient {
	httpClient := http.DefaultClient
	return &ContainerRegistryClient{
		authorizer,
		httpClient,
		serverURL,
		tenantID,
	}
}

// WithHttpClient replaces what http client is used for communication to Azure Container Registry
func (c *ContainerRegistryClient) WithHttpClient(httpClient *http.Client) {
	c.httpClient = httpClient
}

func (c *ContainerRegistryClient) getBaseURL() (*url.URL, error) {
	return parseService(c.serverURL)
}

func parseService(serverURL string) (*url.URL, error) {
	scheme := "https://"
	if strings.HasPrefix(serverURL, "https://") {
		scheme = ""
	}

	serviceURL, err := url.Parse(fmt.Sprintf("%s%s", scheme, serverURL))
	if err != nil {
		return nil, fmt.Errorf("failed to parse server URL - %w", err)
	}

	return serviceURL, nil
}
