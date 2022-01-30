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

// ExchangeRefreshToken exchanges an Azure AD token for an Azure Container Registry refresh token
func (c *ContainerRegistryClient) ExchangeRefreshToken(ctx context.Context) (string, RefreshTokenClaims, error) {
	token, err := c.authorizer.Token()
	if err != nil {
		return "", RefreshTokenClaims{}, err
	}

	serviceURL, err := parseService(c.serverURL)
	if err != nil {
		return "", RefreshTokenClaims{}, err
	}

	data := url.Values{}
	data.Set("grant_type", "access_token")
	data.Set("service", serviceURL.Hostname())
	data.Set("access_token", token.AccessToken)
	if len(c.tenantID) > 0 {
		data.Set("tenant", c.tenantID)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s://%s/oauth2/exchange", serviceURL.Scheme, serviceURL.Host), strings.NewReader(data.Encode()))
	if err != nil {
		return "", RefreshTokenClaims{}, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return "", RefreshTokenClaims{}, err
	}

	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return "", RefreshTokenClaims{}, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", RefreshTokenClaims{}, fmt.Errorf("received non-200 status code - %d: %s", res.StatusCode, string(resBytes))
	}

	var resData struct {
		RefreshToken string `json:"refresh_token"`
	}

	err = json.Unmarshal(resBytes, &resData)
	if err != nil {
		return "", RefreshTokenClaims{}, err
	}

	rtClaims, err := decodeRefreshTokenWithoutValidation(resData.RefreshToken)
	if err != nil {
		return "", RefreshTokenClaims{}, err
	}

	return resData.RefreshToken, rtClaims, nil
}

// RefreshTokenClaimsPermissions contains the actions permitted for the refresh token
type RefreshTokenClaimsPermissions struct {
	Actions    []string `json:"Actions"`
	NotActions []string `json:"NotActions"`
}

// RefreshTokenClaims contians the claims of the refresh token
// Please observe: These claims are in no way validated, only decoded
type RefreshTokenClaims struct {
	JwtID          string                        `json:"jti"`
	Subject        string                        `json:"sub"`
	NotBefore      int64                         `json:"nbf"`
	ExpirationTime int64                         `json:"exp"`
	IssuedAt       int64                         `json:"iat"`
	Issuer         string                        `json:"iss"`
	Audience       string                        `json:"aud"`
	Version        string                        `json:"version"`
	Rid            string                        `json:"rid"`
	GrantType      string                        `json:"grant_type"`
	ApplicationID  string                        `json:"appid"`
	Permissions    RefreshTokenClaimsPermissions `json:"permissions"`
	Roles          []string                      `json:"roles"`
}

func decodeRefreshTokenWithoutValidation(token string) (RefreshTokenClaims, error) {
	parts := strings.SplitN(token, ".", 3)
	claimsBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return RefreshTokenClaims{}, err
	}

	var rtClaims RefreshTokenClaims
	err = json.Unmarshal(claimsBytes, &rtClaims)
	if err != nil {
		return RefreshTokenClaims{}, err
	}

	return rtClaims, nil
}
