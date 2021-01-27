package auth

import (
	"context"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/jws" //nolint:staticcheck
)

type ClientCredentialsType int

const (
	ClientCredentialsAssertionType ClientCredentialsType = iota
	ClientCredentialsSecretType
)

// ClientCredentialsConfig is the configuration for using client credentials flow.
//
// For more information see:
// https://docs.microsoft.com/en-us/azure/active-directory/develop/v2-oauth2-client-creds-grant-flow#get-a-token
// https://docs.microsoft.com/en-us/azure/active-directory/develop/active-directory-certificate-credentials
type ClientCredentialsConfig struct {
	// ClientID is the application's ID.
	ClientID string

	// ClientSecret is the application's secret.
	ClientSecret string

	// PrivateKey contains the contents of an RSA private key or the
	// contents of a PEM file that contains a private key. The provided
	// private key is used to sign JWT assertions.
	// PEM containers with a passphrase are not supported.
	// Use the following command to convert a PKCS 12 file into a PEM.
	//
	//    $ openssl pkcs12 -in key.p12 -out key.pem -nodes
	//
	PrivateKey []byte

	// Certificate contains the (optionally PEM encoded) X509 certificate registered
	// for the application with which you are authenticating.
	Certificate []byte

	// Resource specifies an API resource for which to request access (used for v1 tokens)
	Resource string

	// Scopes specifies a list of requested permission scopes (used for v2 tokens)
	Scopes []string

	// TokenURL is the clientCredentialsToken endpoint. Typically you can use the AzureADEndpoint
	// function to obtain this value, but it may change for non-public clouds.
	TokenURL string

	// Expires optionally specifies how long the clientCredentialsToken is valid for.
	Expires time.Duration

	// Audience optionally specifies the intended audience of the
	// request.  If empty, the value of TokenURL is used as the
	// intended audience.
	Audience string
}

// TokenSource provides a source for obtaining access tokens using clientAssertionAuthorizer or clientSecretAuthorizer.
func (c *ClientCredentialsConfig) TokenSource(ctx context.Context, authType ClientCredentialsType) (source Authorizer) {
	switch authType {
	case ClientCredentialsAssertionType:
		source = CachedAuthorizer(clientAssertionAuthorizer{ctx, c})
	case ClientCredentialsSecretType:
		source = CachedAuthorizer(clientSecretAuthorizer{ctx, c})
	}
	return
}

func clientCredentialsToken(resp *http.Response) (*oauth2.Token, error) {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("oauth2: cannot fetch clientCredentialsToken: %v", err)
	}

	if c := resp.StatusCode; c < 200 || c > 299 {
		return nil, &oauth2.RetrieveError{
			Response: resp,
			Body:     body,
		}
	}

	// clientCredentialsToken response can arrive with numeric values as integers or strings :(
	var tokenRes struct {
		AccessToken string      `json:"access_token"`
		TokenType   string      `json:"token_type"`
		IDToken     string      `json:"id_token"`
		Resource    string      `json:"resource"`
		Scope       string      `json:"scope"`
		ExpiresIn   interface{} `json:"expires_in"` // relative seconds from now
		ExpiresOn   interface{} `json:"expires_on"` // timestamp
	}
	if err := json.Unmarshal(body, &tokenRes); err != nil {
		return nil, fmt.Errorf("oauth2: cannot fetch clientCredentialsToken: %v", err)
	}

	token := &oauth2.Token{
		AccessToken: tokenRes.AccessToken,
		TokenType:   tokenRes.TokenType,
	}
	var secs time.Duration
	if exp, ok := tokenRes.ExpiresIn.(string); ok && exp != "" {
		if v, err := strconv.Atoi(exp); err == nil {
			secs = time.Duration(v)
		}
	} else if exp, ok := tokenRes.ExpiresIn.(int64); ok {
		secs = time.Duration(exp)
	} else if exp, ok := tokenRes.ExpiresIn.(float64); ok {
		secs = time.Duration(exp)
	}
	if secs > 0 {
		token.Expiry = time.Now().Add(secs * time.Second)
	}

	return token, nil
}

type clientAssertionAuthorizer struct {
	ctx  context.Context
	conf *ClientCredentialsConfig
}

func (a clientAssertionAuthorizer) Token() (*oauth2.Token, error) {
	crt := a.conf.Certificate
	if der, _ := pem.Decode(a.conf.Certificate); der != nil {
		crt = der.Bytes
	}
	cert, err := x509.ParseCertificate(crt)
	if err != nil {
		return nil, fmt.Errorf("oauth2: cannot parse certificate: %v", err)
	}
	s := sha1.Sum(cert.Raw)
	fp := base64.URLEncoding.EncodeToString(s[:])
	h := jws.Header{
		Algorithm: "RS256",
		Typ:       "JWT",
		KeyID:     fp,
	}

	claimSet := &jws.ClaimSet{
		Iss: a.conf.ClientID,
		Sub: a.conf.ClientID,
		Aud: a.conf.TokenURL,
	}
	if t := a.conf.Expires; t > 0 {
		claimSet.Exp = time.Now().Add(t).Unix()
	}
	if aud := a.conf.Audience; aud != "" {
		claimSet.Aud = aud
	}

	pk, err := parseKey(a.conf.PrivateKey)
	if err != nil {
		return nil, err
	}

	payload, err := jws.Encode(&h, claimSet, pk)
	if err != nil {
		return nil, err
	}

	hc := oauth2.NewClient(a.ctx, nil)
	v := url.Values{
		"client_assertion":      {payload},
		"client_assertion_type": {"urn:ietf:params:oauth:client-assertion-type:jwt-bearer"},
		"client_id":             {a.conf.ClientID},
		"grant_type":            {"client_credentials"},
	}
	if a.conf.Resource != "" {
		v["resource"] = []string{a.conf.Resource}
	} else {
		v["scope"] = []string{strings.Join(a.conf.Scopes, " ")}
	}
	resp, err := hc.PostForm(a.conf.TokenURL, v)
	if err != nil {
		return nil, fmt.Errorf("oauth2: cannot fetch clientCredentialsToken: %v", err)
	}

	return clientCredentialsToken(resp)
}

// parseKey returns an rsa.PrivateKey containing the provided binary key data.
// If the provided key is PEM encoded, it is decoded first.
func parseKey(key []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(key)
	if block != nil {
		key = block.Bytes
	}
	parsedKey, err := x509.ParsePKCS8PrivateKey(key)
	if err != nil {
		parsedKey, err = x509.ParsePKCS1PrivateKey(key)
		if err != nil {
			return nil, fmt.Errorf("private key should be a PEM or plain PKCS1 or PKCS8; parse error: %v", err)
		}
	}
	parsed, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("private key is invalid")
	}
	return parsed, nil
}

type clientSecretAuthorizer struct {
	ctx  context.Context
	conf *ClientCredentialsConfig
}

func (a clientSecretAuthorizer) Token() (*oauth2.Token, error) {
	hc := oauth2.NewClient(a.ctx, nil)
	v := url.Values{
		"client_id":     {a.conf.ClientID},
		"client_secret": {a.conf.ClientSecret},
		"grant_type":    {"client_credentials"},
	}
	if a.conf.Resource != "" {
		v["resource"] = []string{a.conf.Resource}
	} else {
		v["scope"] = []string{strings.Join(a.conf.Scopes, " ")}
	}
	resp, err := hc.PostForm(a.conf.TokenURL, v)
	if err != nil {
		return nil, fmt.Errorf("oauth2: cannot fetch clientCredentialsToken: %v", err)
	}

	return clientCredentialsToken(resp)
}
