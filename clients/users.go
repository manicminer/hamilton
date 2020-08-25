package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/base"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/manicminer/hamilton/models"
)

type UsersClient struct {
	BaseClient base.BaseClient
}

func NewUsersClient(authorizer auth.Authorizer, tenantId string) *UsersClient {
	return &UsersClient{
		BaseClient: base.NewBaseClient(authorizer, base.DefaultEndpoint, tenantId, base.VersionBeta),
	}
}

func (c *UsersClient) List(ctx context.Context, filter string) (*[]models.User, error) {
	params := url.Values{}
	if filter != "" {
		params.Add("$filter", filter)
	}
	resp, err := c.BaseClient.Get(ctx, base.GetHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		Uri:              fmt.Sprintf("/users?%s", params.Encode()),
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	var data struct {
		Users []models.User `json:"value"`
	}
	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, err
	}
	return &data.Users, nil
}

func (c *UsersClient) Create(ctx context.Context, user models.User) (*models.User, error) {
	body, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}
	resp, err := c.BaseClient.Post(ctx, base.PostHttpRequestInput{
		Body:             body,
		ValidStatusCodes: []int{http.StatusCreated},
		Uri:              "/users",
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	var newUser models.User
	if err := json.Unmarshal(respBody, &newUser); err != nil {
		return nil, err
	}
	return &newUser, nil
}

func (c *UsersClient) Get(ctx context.Context, id string) (*models.User, error) {
	resp, err := c.BaseClient.Get(ctx, base.GetHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		Uri:              fmt.Sprintf("/users/%s", id),
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	var user models.User
	if err := json.Unmarshal(respBody, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (c *UsersClient) Update(ctx context.Context, user models.User) error {
	body, err := json.Marshal(user)
	if err != nil {
		return err
	}
	_, err = c.BaseClient.Patch(ctx, base.PatchHttpRequestInput{
		Body:             body,
		ValidStatusCodes: []int{http.StatusNoContent},
		Uri:              fmt.Sprintf("/users/%s", *user.ID),
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *UsersClient) Delete(ctx context.Context, id string) error {
	_, err := c.BaseClient.Delete(ctx, base.DeleteHttpRequestInput{
		ValidStatusCodes: []int{http.StatusNoContent},
		Uri:              fmt.Sprintf("/users/%s", id),
	})
	if err != nil {
		return err
	}
	return nil
}
