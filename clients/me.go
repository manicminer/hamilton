package clients

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/base"
	"github.com/manicminer/hamilton/models"
)

type MeClient struct {
	BaseClient base.Client
}

func NewMeClient(authorizer auth.Authorizer, tenantId string) *MeClient {
	return &MeClient{
		BaseClient: base.NewClient(authorizer, base.DefaultEndpoint, tenantId, base.VersionBeta),
	}
}

func (c *MeClient) Get(ctx context.Context) (*models.Me, int, error) {
	var status int
	resp, status, err := c.BaseClient.Get(ctx, base.GetHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		Uri: base.Uri{
			Entity:      "/me",
			HasTenantId: false,
		},
	})
	if err != nil {
		return nil, status, err
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	var me models.Me
	if err := json.Unmarshal(respBody, &me); err != nil {
		return nil, status, err
	}
	return &me, status, nil
}

func (c *MeClient) GetProfile(ctx context.Context) (*models.Me, int, error) {
	var status int
	resp, status, err := c.BaseClient.Get(ctx, base.GetHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		Uri: base.Uri{
			Entity:      "/me/profile",
			HasTenantId: false,
		},
	})
	if err != nil {
		return nil, status, err
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	var me models.Me
	if err := json.Unmarshal(respBody, &me); err != nil {
		return nil, status, err
	}
	return &me, status, nil
}
