package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/base"
	"github.com/manicminer/hamilton/models"
)

type UsersClient struct {
	BaseClient base.Client
}

func NewUsersClient(authorizer auth.Authorizer, tenantId string) *UsersClient {
	return &UsersClient{
		BaseClient: base.NewClient(authorizer, base.DefaultEndpoint, tenantId, base.VersionBeta),
	}
}

func (c *UsersClient) List(ctx context.Context, filter string) (*[]models.User, int, error) {
	params := url.Values{}
	if filter != "" {
		params.Add("$filter", filter)
	}
	resp, status, err := c.BaseClient.Get(ctx, base.GetHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		Uri: base.Uri{
			Entity:      "/users",
			Params:      params,
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, err
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	var data struct {
		Users []models.User `json:"value"`
	}
	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, status, err
	}
	return &data.Users, status, nil
}

func (c *UsersClient) Create(ctx context.Context, user models.User) (*models.User, int, error) {
	var status int
	body, err := json.Marshal(user)
	if err != nil {
		return nil, status, err
	}
	resp, status, err := c.BaseClient.Post(ctx, base.PostHttpRequestInput{
		Body:             body,
		ValidStatusCodes: []int{http.StatusCreated},
		Uri: base.Uri{
			Entity:      "/users",
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, err
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	var newUser models.User
	if err := json.Unmarshal(respBody, &newUser); err != nil {
		return nil, status, err
	}
	return &newUser, status, nil
}

func (c *UsersClient) Get(ctx context.Context, id string) (*models.User, int, error) {
	resp, status, err := c.BaseClient.Get(ctx, base.GetHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		Uri: base.Uri{
			Entity:      fmt.Sprintf("/users/%s", id),
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, err
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	var user models.User
	if err := json.Unmarshal(respBody, &user); err != nil {
		return nil, status, err
	}
	return &user, status, nil
}

func (c *UsersClient) Update(ctx context.Context, user models.User) (int, error) {
	var status int
	body, err := json.Marshal(user)
	if err != nil {
		return status, err
	}
	_, status, err = c.BaseClient.Patch(ctx, base.PatchHttpRequestInput{
		Body:             body,
		ValidStatusCodes: []int{http.StatusNoContent},
		Uri: base.Uri{
			Entity:      fmt.Sprintf("/users/%s", *user.ID),
			HasTenantId: true,
		},
	})
	if err != nil {
		return status, err
	}
	return status, nil
}

func (c *UsersClient) Delete(ctx context.Context, id string) (int, error) {
	_, status, err := c.BaseClient.Delete(ctx, base.DeleteHttpRequestInput{
		ValidStatusCodes: []int{http.StatusNoContent},
		Uri: base.Uri{
			Entity:      fmt.Sprintf("/users/%s", id),
			HasTenantId: true,
		},
	})
	if err != nil {
		return status, err
	}
	return status, nil
}
