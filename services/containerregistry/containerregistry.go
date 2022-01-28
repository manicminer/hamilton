package containerregistry

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/manicminer/hamilton/auth"
)

type ContainerRegistry struct {
	authorizer auth.Authorizer
	serverURL  string
	tenantID   string
}

func New(authorizer auth.Authorizer, serverURL string, tenantID string) *ContainerRegistry {
	return &ContainerRegistry{
		authorizer,
		serverURL,
		tenantID,
	}
}

func (c ContainerRegistry) ExchangeToken(ctx context.Context) (string, error) {
	token, err := c.authorizer.Token()
	if err != nil {
		return "", err
	}

	service, err := getService(c.serverURL)
	if err != nil {
		return "", err
	}

	data := url.Values{}
	data.Set("grant_type", "access_token")
	data.Set("service", service)
	data.Set("access_token", token.AccessToken)
	if len(c.tenantID) > 0 {
		data.Set("tenant", c.tenantID)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("https://%s/oauth2/exchange", service), strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 status code: %d", res.StatusCode)
	}

	var resData struct {
		RefreshToken string `json:"refresh_token"`
	}

	err = json.Unmarshal(resBytes, &resData)
	if err != nil {
		return "", err
	}

	return resData.RefreshToken, nil
}

func getService(serverURL string) (string, error) {
	scheme := "https://"
	if strings.HasPrefix(serverURL, "https://") {
		scheme = ""
	}

	serviceURL, err := url.Parse(fmt.Sprintf("%s%s", scheme, serverURL))
	if err != nil {
		return "", fmt.Errorf("failed to parse server URL - %w", err)
	}

	return serviceURL.Hostname(), nil
}
