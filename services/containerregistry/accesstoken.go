package containerregistry

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// ExchangeAccessToken exchanges an Azure Container Registry refresh token for an Azure Container Registry access token
func (c *ContainerRegistryClient) ExchangeAccessToken(ctx context.Context, refreshToken string, scopes AccessTokenScopes) (string, AccessTokenClaims, error) {
	if len(scopes) == 0 {
		return "", AccessTokenClaims{}, fmt.Errorf("at least one scope is required")
	}
	serviceURL, err := parseService(c.serverURL)
	if err != nil {
		return "", AccessTokenClaims{}, err
	}

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	data.Set("service", serviceURL.Hostname())
	for _, scope := range scopes {
		data.Add("scope", scope.String())
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s://%s/oauth2/token", serviceURL.Scheme, serviceURL.Host), strings.NewReader(data.Encode()))
	if err != nil {
		return "", AccessTokenClaims{}, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return "", AccessTokenClaims{}, err
	}

	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return "", AccessTokenClaims{}, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", AccessTokenClaims{}, fmt.Errorf("received non-200 status code - %d: %s", res.StatusCode, string(resBytes))
	}

	var resData struct {
		AccessToken string `json:"access_token"`
	}

	err = json.Unmarshal(resBytes, &resData)
	if err != nil {
		return "", AccessTokenClaims{}, err
	}

	atClaims, err := decodeAccessTokenWithoutValidation(resData.AccessToken)
	if err != nil {
		return "", AccessTokenClaims{}, err
	}

	return resData.AccessToken, atClaims, nil
}

// AccessTokenScope contains the information about what access the access token has to the Azure Container Registry
type AccessTokenScope struct {
	Type    string   `json:"Type"`
	Name    string   `json:"Name"`
	Actions []string `json:"Actions"`
}

// String outputs the string of the scope in a scope compatible way
func (scope AccessTokenScope) String() string {
	return fmt.Sprintf("%s:%s:%s", scope.Type, scope.Name, strings.Join(scope.Actions, ","))
}

// AccessTokenScopes is just an array of AccessTokenScope
type AccessTokenScopes []AccessTokenScope

// AccessTokenClaims contians the claims of the access token
// Please observe: These claims are in no way validated, only decoded
type AccessTokenClaims struct {
	JwtID          string            `json:"jti"`
	Subject        string            `json:"sub"`
	NotBefore      int64             `json:"nbf"`
	ExpirationTime int64             `json:"exp"`
	IssuedAt       int64             `json:"iat"`
	Issuer         string            `json:"iss"`
	Audience       string            `json:"aud"`
	Version        string            `json:"version"`
	Rid            string            `json:"rid"`
	GrantType      string            `json:"grant_type"`
	ApplicationID  string            `json:"appid"`
	Scopes         AccessTokenScopes `json:"access"`
	Roles          []string          `json:"roles"`
}

func decodeAccessTokenWithoutValidation(token string) (AccessTokenClaims, error) {
	parts := strings.SplitN(token, ".", 3)
	claimsBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return AccessTokenClaims{}, err
	}

	var atClaims AccessTokenClaims
	err = json.Unmarshal(claimsBytes, &atClaims)
	if err != nil {
		return AccessTokenClaims{}, err
	}

	return atClaims, nil
}
