package clients

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"

	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/base"
	"github.com/manicminer/hamilton/models"
)

type ApplicationsClient struct {
	BaseClient base.Client
}

func NewApplicationsClient(authorizer auth.Authorizer, tenantId string) *ApplicationsClient {
	return &ApplicationsClient{
		BaseClient: base.NewClient(authorizer, base.DefaultEndpoint, tenantId, base.VersionBeta),
	}
}

func (c *ApplicationsClient) List(ctx context.Context, filter string) (*[]models.Application, int, error) {
	params := url.Values{}
	if filter != "" {
		params.Add("$filter", filter)
	}
	resp, status, err := c.BaseClient.Get(ctx, base.GetHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		Uri: base.Uri{
			Entity:      "/applications",
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
		Applications []models.Application `json:"value"`
	}
	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, status, err
	}
	return &data.Applications, status, nil
}

func (c *ApplicationsClient) Create(ctx context.Context, application models.Application) (*models.Application, int, error) {
	var status int
	body, err := json.Marshal(application)
	if err != nil {
		return nil, status, err
	}
	resp, status, err := c.BaseClient.Post(ctx, base.PostHttpRequestInput{
		Body:             body,
		ValidStatusCodes: []int{http.StatusCreated},
		Uri: base.Uri{
			Entity:      "/applications",
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, err
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	var newApplication models.Application
	if err := json.Unmarshal(respBody, &newApplication); err != nil {
		return nil, status, err
	}
	return &newApplication, status, nil
}

func (c *ApplicationsClient) Get(ctx context.Context, id string) (*models.Application, int, error) {
	resp, status, err := c.BaseClient.Get(ctx, base.GetHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		Uri: base.Uri{
			Entity:      fmt.Sprintf("/applications/%s", id),
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, err
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	var application models.Application
	if err := json.Unmarshal(respBody, &application); err != nil {
		return nil, status, err
	}
	return &application, status, nil
}

func (c *ApplicationsClient) Update(ctx context.Context, application models.Application) (int, error) {
	var status int
	if application.ID == nil {
		return status, errors.New("cannot update application with nil ID")
	}
	body, err := json.Marshal(application)
	if err != nil {
		return status, err
	}
	_, status, err = c.BaseClient.Patch(ctx, base.PatchHttpRequestInput{
		Body:             body,
		ValidStatusCodes: []int{http.StatusNoContent},
		Uri: base.Uri{
			Entity:      fmt.Sprintf("/applications/%s", *application.ID),
			HasTenantId: true,
		},
	})
	if err != nil {
		return status, err
	}
	return status, nil
}

func (c *ApplicationsClient) Delete(ctx context.Context, id string) (int, error) {
	_, status, err := c.BaseClient.Delete(ctx, base.DeleteHttpRequestInput{
		ValidStatusCodes: []int{http.StatusNoContent},
		Uri: base.Uri{
			Entity:      fmt.Sprintf("/applications/%s", id),
			HasTenantId: true,
		},
	})
	if err != nil {
		return status, err
	}
	return status, nil
}

func (c *ApplicationsClient) ListOwners(ctx context.Context, id string) (*[]string, int, error) {
	resp, status, err := c.BaseClient.Get(ctx, base.GetHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		Uri: base.Uri{
			Entity:      fmt.Sprintf("/applications/%s/owners?$select=id", id),
			HasTenantId: true,
		},

	})
	if err != nil {
		return nil, status, err
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	var data struct {
		Owners []struct {
			Type string `json:"@odata.type"`
			Id   string `json:"id"`
		} `json:"value"`
	}
	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, status, err
	}
	ret := make([]string, len(data.Owners))
	for i, v := range data.Owners {
		ret[i] = v.Id
	}
	return &ret, status, nil
}

func (c *ApplicationsClient) GetOwner(ctx context.Context, applicationId, ownerId string) (*string, int, error) {
	resp, status, err := c.BaseClient.Get(ctx, base.GetHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		Uri: base.Uri{
			Entity:      fmt.Sprintf("/applications/%s/owners/%s/$ref?$select=id,url", applicationId, ownerId),
			HasTenantId: true,
		},

	})
	if err != nil {
		return nil, status, err
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	var data struct {
		Context string `json:"@odata.context"`
		Type    string `json:"@odata.type"`
		Id      string `json:"id"`
		Url     string `json:"url"`
	}
	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, status, err
	}
	return &data.Id, status, nil
}

func (c *ApplicationsClient) AddOwners(ctx context.Context, application *models.Application) (int, error) {
	var status int
	if application.ID == nil {
		return status, errors.New("cannot update application with nil ID")
	}
	if application.Owners == nil {
		return status, errors.New("cannot update application with nil Owners")
	}
	for _, owner := range *application.Owners {
		data := struct {
			Owner string `json:"@odata.id"`
		}{
			Owner: owner,
		}
		body, err := json.Marshal(data)
		if err != nil {
			return status, err
		}
		_, status, err = c.BaseClient.Post(ctx, base.PostHttpRequestInput{
			Body:             body,
			ValidStatusCodes: []int{http.StatusNoContent},
			Uri: base.Uri{
				Entity:      fmt.Sprintf("/applications/%s/owners/$ref", *application.ID),
				HasTenantId: true,
			},
		})
		if err != nil {
			return status, err
		}
	}
	return status, nil
}

func (c *ApplicationsClient) RemoveOwners(ctx context.Context, id string, ownerIds *[]string) (int, error) {
	var status int
	if ownerIds == nil {
		return status, errors.New("cannot remove, nil ownerIds")
	}
	for _, ownerId := range *ownerIds {
		// check for ownership before attempting deletion
		if _, status, err := c.GetOwner(ctx, id, ownerId); err != nil {
			if status == http.StatusNotFound {
				continue
			}
			return status, err
		}

		// despite the above check, sometimes owners are just gone
		checkOwnerGone := func(resp *http.Response) bool {
			if resp.StatusCode == http.StatusBadRequest {
				defer resp.Body.Close()
				respBody, _ := ioutil.ReadAll(resp.Body)
				var apiError models.Error
				if err := json.Unmarshal(respBody, &apiError); err != nil {
					return false
				}
				re := regexp.MustCompile("One or more removed object references do not exist")
				if re.MatchString(apiError.Error.Message) {
					return true
				}
			}
			return false
		}

		_, status, err := c.BaseClient.Delete(ctx, base.DeleteHttpRequestInput{
			ValidStatusCodes: []int{http.StatusNoContent},
			ValidStatusFunc:  checkOwnerGone,
			Uri: base.Uri{
				Entity:      fmt.Sprintf("/applications/%s/owners/%s/$ref", id, ownerId),
				HasTenantId: true,
			},
		})
		if err != nil {
			return status, err
		}
	}
	return status, nil
}
