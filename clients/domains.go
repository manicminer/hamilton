package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/base"
	"github.com/manicminer/hamilton/models"
)

type DomainsClient struct {
	BaseClient base.BaseClient
}

func NewDomainsClient(authorizer auth.Authorizer, tenantId string) *DomainsClient {
	return &DomainsClient{
		BaseClient: base.NewBaseClient(authorizer, base.DefaultEndpoint, tenantId, base.Version10),
	}
}

func (c *DomainsClient) List(ctx context.Context) (*[]models.Domain, error) {
	resp, err := c.BaseClient.Get(ctx, base.GetHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		Uri:              "/domains",
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	var data struct {
		Domains []models.Domain `json:"value"`
	}
	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, err
	}
	return &data.Domains, nil
}

func (c *DomainsClient) Get(ctx context.Context, id string) (*models.Domain, error) {
	resp, err := c.BaseClient.Get(ctx, base.GetHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		Uri:              fmt.Sprintf("/domains/%s", id),
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	var domain models.Domain
	if err := json.Unmarshal(respBody, &domain); err != nil {
		return nil, err
	}
	return &domain, nil
}
