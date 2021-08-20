package msgraph

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/manicminer/hamilton/odata"
)

// AuthenticationClient performs operations on the Authentications methods endpoint under Identity and Sign-in
type AuthenticationMethodsClient struct {
	BaseClient Client
}

// NewAuthenticationClient returns a new ApplicationsClient
func NewAuthenticationMethodsClient(tenantId string) *AuthenticationMethodsClient {
	return &AuthenticationMethodsClient{
		BaseClient: NewClient(Version10, tenantId),
	}
}

//List all authentication methods
func (c *AuthenticationMethodsClient) List(ctx context.Context, userID string, query odata.Query) (*[]AuthenticationMethod, int, error) {
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		DisablePaging:    query.Top > 0,
		ValidStatusCodes: []int{http.StatusOK},
		Uri: Uri{
			Entity:      fmt.Sprintf("/users/%s/authentication/methods", userID),
			Params:      query.Values(),
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("ApplicationsClient.BaseClient.Get(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var data struct {
		AuthenticationMethods *[]json.RawMessage `json:"value"`
	}

	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	//The graph API returns a mixture of types, this loop matches up the result to the appropriate model

	var ret []AuthenticationMethod

	if data.AuthenticationMethods == nil {
		return &ret, status, nil
	}

	for _, authMethod := range *data.AuthenticationMethods {
		var o odata.OData
		if err := json.Unmarshal(authMethod, &o); err != nil {
			return nil, status, fmt.Errorf("json.Unmarshall(): %v", err)
		}

		if o.Type == nil {
			continue
		}
		switch *o.Type {
		case odata.TypeFido2AuthenticationMethod:
			var auth Fido2AuthenticationMethod
			if err := json.Unmarshal(authMethod, &auth); err != nil {
				return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
			}
			ret = append(ret, auth)
		case odata.TypeMicrosoftAuthenticatorAuthenticationMethod:
			var auth MicrosoftAuthenticatorAuthenticationMethod
			if err := json.Unmarshal(authMethod, &auth); err != nil {
				return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
			}
			ret = append(ret, auth)
		case odata.TypeWindowsHelloForBusinessAuthenticationMethod:
			var auth WindowsHelloForBusinessAuthenticationMethod
			if err := json.Unmarshal(authMethod, &auth); err != nil {
				return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
			}
			ret = append(ret, auth)
		case odata.TypeTemporaryAccessPassAuthenticationMethod:
			var auth TemporaryAccessPassAuthenticationMethod
			if err := json.Unmarshal(authMethod, &auth); err != nil {
				return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
			}
			ret = append(ret, auth)
		case odata.TypePhoneAuthenticationMethod:
			var auth PhoneAuthenticationMethod
			if err := json.Unmarshal(authMethod, &auth); err != nil {
				return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
			}
			ret = append(ret, auth)
		case odata.TypeEmailAuthenticationMethod:
			var auth EmailAuthenticationMethod
			if err := json.Unmarshal(authMethod, &auth); err != nil {
				return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
			}
			ret = append(ret, auth)
		case odata.TypePasswordAuthenticationMethod:
			var auth PasswordAuthenticationMethod
			if err := json.Unmarshal(authMethod, &auth); err != nil {
				return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
			}
			ret = append(ret, auth)
		}
	}

	return &ret, status, nil
}

func (c *AuthenticationMethodsClient) ListFido2Methods(ctx context.Context, userID string, query odata.Query) (*[]Fido2AuthenticationMethod, int, error) {
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		DisablePaging:    query.Top > 0,
		ValidStatusCodes: []int{http.StatusOK},
		Uri: Uri{
			Entity:      fmt.Sprintf("/users/%s/authentication/fido2Methods", userID),
			Params:      query.Values(),
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("AuthenticationMethodsClient.BaseClient.Get(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var data struct {
		Fido2Methods []Fido2AuthenticationMethod `json:"value"`
	}
	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &data.Fido2Methods, status, nil
}

func (c *AuthenticationMethodsClient) GetFido2Method(ctx context.Context, userID, id string, query odata.Query) (*Fido2AuthenticationMethod, int, error) {
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		ValidStatusCodes:       []int{http.StatusOK},
		Uri: Uri{
			Entity:      fmt.Sprintf("/users/%s/authentication/fido2Methods/%s", userID, id),
			Params:      query.Values(),
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("AuthenticationMethodsClient.BaseClient.Get(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var fido2Method Fido2AuthenticationMethod
	if err := json.Unmarshal(respBody, &fido2Method); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &fido2Method, status, nil
}

func (c *AuthenticationMethodsClient) DeleteFido2Method(ctx context.Context, userID, id string) (int, error) {
	_, status, _, err := c.BaseClient.Delete(ctx, DeleteHttpRequestInput{
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		ValidStatusCodes:       []int{http.StatusNoContent},
		Uri: Uri{
			Entity:      fmt.Sprintf("/users/%s/authentication/fido2Methods/%s", userID, id),
			HasTenantId: true,
		},
	})
	if err != nil {
		return status, fmt.Errorf("AuthenticationMethodsClient.BaseClient.Delete(): %v", err)
	}

	return status, nil
}

func (c *AuthenticationMethodsClient) ListMSAuthenticatorMethods(ctx context.Context, userID string, query odata.Query) (*[]MicrosoftAuthenticatorAuthenticationMethod, int, error) {
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		DisablePaging:    query.Top > 0,
		ValidStatusCodes: []int{http.StatusOK},
		Uri: Uri{
			Entity:      fmt.Sprintf("/users/%s/authentication/microsoftAuthenticatorMethods", userID),
			Params:      query.Values(),
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("AuthenticationMethodsClient.BaseClient.Get(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var data struct {
		MicrosoftAuthenticatorMethods []MicrosoftAuthenticatorAuthenticationMethod `json:"value"`
	}
	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &data.MicrosoftAuthenticatorMethods, status, nil
}

func (c *AuthenticationMethodsClient) GetMSAuthenticatorMethod(ctx context.Context, userID, id string, query odata.Query) (*MicrosoftAuthenticatorAuthenticationMethod, int, error) {
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		ValidStatusCodes:       []int{http.StatusOK},
		Uri: Uri{
			Entity:      fmt.Sprintf("/users/%s/authentication/microsoftAuthenticatorMethods/%s", userID, id),
			Params:      query.Values(),
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("AuthenticationMethodsClient.BaseClient.Get(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var microsoftAuthenticatorMethod MicrosoftAuthenticatorAuthenticationMethod
	if err := json.Unmarshal(respBody, &microsoftAuthenticatorMethod); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &microsoftAuthenticatorMethod, status, nil
}

func (c *AuthenticationMethodsClient) DeleteMSAuthenticatorMethod(ctx context.Context, userID, id string) (int, error) {
	_, status, _, err := c.BaseClient.Delete(ctx, DeleteHttpRequestInput{
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		ValidStatusCodes:       []int{http.StatusNoContent},
		Uri: Uri{
			Entity:      fmt.Sprintf("/users/%s/authentication/microsoftAuthenticatorMethods/%s", userID, id),
			HasTenantId: true,
		},
	})
	if err != nil {
		return status, fmt.Errorf("AuthenticationMethodsClient.BaseClient.Delete(): %v", err)
	}

	return status, nil
}

func (c *AuthenticationMethodsClient) ListWindowsHelloMethods(ctx context.Context, userID string, query odata.Query) (*[]WindowsHelloForBusinessAuthenticationMethod, int, error) {
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		DisablePaging:    query.Top > 0,
		ValidStatusCodes: []int{http.StatusOK},
		Uri: Uri{
			Entity:      fmt.Sprintf("/users/%s/authentication/windowsHelloForBusinessMethods", userID),
			Params:      query.Values(),
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("AuthenticationMethodsClient.BaseClient.Get(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var data struct {
		WindowsHelloForBusinessMethods []WindowsHelloForBusinessAuthenticationMethod `json:"value"`
	}
	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &data.WindowsHelloForBusinessMethods, status, nil
}

func (c *AuthenticationMethodsClient) GetWindowsHelloMethod(ctx context.Context, userID, id string, query odata.Query) (*WindowsHelloForBusinessAuthenticationMethod, int, error) {
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		ValidStatusCodes:       []int{http.StatusOK},
		Uri: Uri{
			Entity:      fmt.Sprintf("/users/%s/authentication/windowsHelloForBusinessMethods/%s", userID, id),
			Params:      query.Values(),
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("AuthenticationMethodsClient.BaseClient.Get(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var windowsHelloForBusinessMethod WindowsHelloForBusinessAuthenticationMethod
	if err := json.Unmarshal(respBody, &windowsHelloForBusinessMethod); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &windowsHelloForBusinessMethod, status, nil
}

func (c *AuthenticationMethodsClient) DeleteWindowsHelloMethod(ctx context.Context, userID, id string) (int, error) {
	_, status, _, err := c.BaseClient.Delete(ctx, DeleteHttpRequestInput{
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		ValidStatusCodes:       []int{http.StatusNoContent},
		Uri: Uri{
			Entity:      fmt.Sprintf("/users/%s/authentication/windowsHelloForBusinessMethods/%s", userID, id),
			HasTenantId: true,
		},
	})
	if err != nil {
		return status, fmt.Errorf("AuthenticationMethodsClient.BaseClient.Delete(): %v", err)
	}

	return status, nil
}

func (c *AuthenticationMethodsClient) ListTempAccessPassMethods(ctx context.Context, userID string, query odata.Query) (*[]TemporaryAccessPassAuthenticationMethod, int, error) {
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		DisablePaging:    query.Top > 0,
		ValidStatusCodes: []int{http.StatusOK},
		Uri: Uri{
			Entity:      fmt.Sprintf("/users/%s/authentication/temporaryAccessPassMethods", userID),
			Params:      query.Values(),
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("AuthenticationMethodsClient.BaseClient.Get(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var data struct {
		TempAccessPassMethods []TemporaryAccessPassAuthenticationMethod `json:"value"`
	}
	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &data.TempAccessPassMethods, status, nil
}

func (c *AuthenticationMethodsClient) GetTempAccessPassMethod(ctx context.Context, userID, id string, query odata.Query) (*TemporaryAccessPassAuthenticationMethod, int, error) {
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		ValidStatusCodes:       []int{http.StatusOK},
		Uri: Uri{
			Entity:      fmt.Sprintf("/users/%s/authentication/temporaryAccessPassMethods/%s", userID, id),
			Params:      query.Values(),
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("AuthenticationMethodsClient.BaseClient.Get(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var temporaryAccessPassMethod TemporaryAccessPassAuthenticationMethod
	if err := json.Unmarshal(respBody, &temporaryAccessPassMethod); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &temporaryAccessPassMethod, status, nil
}

func (c *AuthenticationMethodsClient) CreateTempAccessPassMethod(ctx context.Context, userID string, accessPass TemporaryAccessPassAuthenticationMethod) (*TemporaryAccessPassAuthenticationMethod, int, error) {
	var status int

	body, err := json.Marshal(accessPass)
	if err != nil {
		return nil, status, fmt.Errorf("json.Marshal(): %v", err)
	}

	resp, status, _, err := c.BaseClient.Post(ctx, PostHttpRequestInput{
		Body:             body,
		ValidStatusCodes: []int{http.StatusCreated},
		Uri: Uri{
			Entity:      fmt.Sprintf("/users/%s/authentication/temporaryAccessPassMethods", userID),
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("AuthenticationMethodsClient.BaseClient.Post(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var newTempAccessPassAuthMethod TemporaryAccessPassAuthenticationMethod
	if err := json.Unmarshal(respBody, &newTempAccessPassAuthMethod); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &newTempAccessPassAuthMethod, status, nil
}

func (c *AuthenticationMethodsClient) DeleteTempAccessPassMethod(ctx context.Context, userID, id string) (int, error) {
	_, status, _, err := c.BaseClient.Delete(ctx, DeleteHttpRequestInput{
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		ValidStatusCodes:       []int{http.StatusNoContent},
		Uri: Uri{
			Entity:      fmt.Sprintf("/users/%s/authentication/temporaryAccessPassMethods/%s", userID, id),
			HasTenantId: true,
		},
	})
	if err != nil {
		return status, fmt.Errorf("AuthenticationMethodsClient.BaseClient.Delete(): %v", err)
	}

	return status, nil
}
