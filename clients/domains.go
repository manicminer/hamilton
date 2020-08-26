package clients

import (
	"context"
	"encoding/json"
	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/base"
	"github.com/manicminer/hamilton/models"
	"io/ioutil"
	"net/http"
)

type DomainsClient struct {
	BaseClient base.BaseClient
}

func NewDomainsClient(authorizer auth.Authorizer, tenantId string) *DomainsClient {
	return &DomainsClient{
		BaseClient: base.NewBaseClient(authorizer, base.DefaultEndpoint, tenantId, base.VersionBeta),
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


