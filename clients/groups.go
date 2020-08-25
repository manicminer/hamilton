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

type GroupsClient struct {
	BaseClient base.BaseClient
}

func NewGroupsClient(authorizer auth.Authorizer, tenantId string) *GroupsClient {
	return &GroupsClient{
		BaseClient: base.NewBaseClient(authorizer, base.DefaultEndpoint, tenantId, base.VersionBeta),
	}
}

func (c *GroupsClient) List(ctx context.Context, filter string) (*[]models.Group, error) {
	params := url.Values{}
	if filter != "" {
		params.Add("$filter", filter)
	}
	resp, err := c.BaseClient.Get(ctx, base.GetHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		Uri:              fmt.Sprintf("/groups?%s", params.Encode()),
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	var data struct {
		Groups []models.Group `json:"value"`
	}
	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, err
	}
	return &data.Groups, nil
}

func (c *GroupsClient) Create(ctx context.Context, group models.Group) (*models.Group, error) {
	body, err := json.Marshal(group)
	if err != nil {
		return nil, err
	}
	resp, err := c.BaseClient.Post(ctx, base.PostHttpRequestInput{
		Body:             body,
		ValidStatusCodes: []int{http.StatusCreated},
		Uri:              "/groups",
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	var newGroup models.Group
	if err := json.Unmarshal(respBody, &newGroup); err != nil {
		return nil, err
	}
	return &newGroup, nil
}

func (c *GroupsClient) Get(ctx context.Context, id string) (*models.Group, error) {
	resp, err := c.BaseClient.Get(ctx, base.GetHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		Uri:              fmt.Sprintf("/groups/%s", id),
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	var group models.Group
	if err := json.Unmarshal(respBody, &group); err != nil {
		return nil, err
	}
	return &group, nil
}

func (c *GroupsClient) Update(ctx context.Context, group models.Group) error {
	body, err := json.Marshal(group)
	if err != nil {
		return err
	}
	_, err = c.BaseClient.Patch(ctx, base.PatchHttpRequestInput{
		Body:             body,
		ValidStatusCodes: []int{http.StatusNoContent},
		Uri:              fmt.Sprintf("/groups/%s", *group.ID),
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *GroupsClient) Delete(ctx context.Context, id string) error {
	_, err := c.BaseClient.Delete(ctx, base.DeleteHttpRequestInput{
		ValidStatusCodes: []int{http.StatusNoContent},
		Uri:              fmt.Sprintf("/groups/%s", id),
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *GroupsClient) ListMembers(ctx context.Context, id string) (*[]string, error) {
	resp, err := c.BaseClient.Get(ctx, base.GetHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		Uri:              fmt.Sprintf("/groups/%s/members?$select=id", id),
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	var data struct {
		Members []struct {
			Type string `json:"@odata.type"`
			Id   string `json:"id"`
		} `json:"value"`
	}
	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, err
	}
	ret := make([]string, len(data.Members))
	for i, v := range data.Members {
		ret[i] = v.Id
	}
	return &ret, nil
}

func (c *GroupsClient) GetMember(ctx context.Context, groupId, memberId string) (*string, error) {
	resp, err := c.BaseClient.Get(ctx, base.GetHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		Uri:              fmt.Sprintf("/groups/%s/members/%s/$ref?$select=id,url", groupId, memberId),
	})
	if err != nil {
		return nil, err
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
		return nil, err
	}
	return &data.Id, nil
}

func (c *GroupsClient) AddMembers(ctx context.Context, group *models.Group) error {
	// Patching group members support up to 20 members per request
	var memberChunks [][]string
	members := *group.Members
	max := len(members)
	// Chunk into slices of 20 for batching
	for i := 0; i < max; i += 20 {
		end := i + 20
		if end > max {
			end = max
		}
		memberChunks = append(memberChunks, members[i:end])
	}
	for _, members := range memberChunks {
		data := struct {
			Members []string `json:"members@odata.bind"`
		}{
			Members: members,
		}
		body, err := json.Marshal(data)
		if err != nil {
			return err
		}
		_, err = c.BaseClient.Patch(ctx, base.PatchHttpRequestInput{
			Body:             body,
			ValidStatusCodes: []int{http.StatusNoContent},
			Uri:              fmt.Sprintf("/groups/%s", *group.ID),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *GroupsClient) RemoveMembers(ctx context.Context, id string, memberIds *[]string) error {
	for _, memberId := range *memberIds {
		// check for membership before attempting deletion
		if _, err := c.GetMember(ctx, id, memberId); err != nil {
			continue
		}
		_, err := c.BaseClient.Delete(ctx, base.DeleteHttpRequestInput{
			ValidStatusCodes: []int{http.StatusNoContent},
			Uri:              fmt.Sprintf("/groups/%s/members/%s/$ref", id, memberId),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *GroupsClient) ListOwners(ctx context.Context, id string) (*[]string, error) {
	resp, err := c.BaseClient.Get(ctx, base.GetHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		Uri:              fmt.Sprintf("/groups/%s/owners?$select=id", id),
	})
	if err != nil {
		return nil, err
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
		return nil, err
	}
	ret := make([]string, len(data.Owners))
	for i, v := range data.Owners {
		ret[i] = v.Id
	}
	return &ret, nil
}

func (c *GroupsClient) GetOwner(ctx context.Context, groupId, ownerId string) (*string, error) {
	resp, err := c.BaseClient.Get(ctx, base.GetHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		Uri:              fmt.Sprintf("/groups/%s/owners/%s/$ref?$select=id,url", groupId, ownerId),
	})
	if err != nil {
		return nil, err
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
		return nil, err
	}
	return &data.Id, nil
}

func (c *GroupsClient) AddOwners(ctx context.Context, group *models.Group) error {
	for _, owner := range *group.Owners {
		data := struct {
			Owner string `json:"@odata.id"`
		}{
			Owner: owner,
		}
		body, err := json.Marshal(data)
		if err != nil {
			return err
		}
		_, err = c.BaseClient.Post(ctx, base.PostHttpRequestInput{
			Body:             body,
			ValidStatusCodes: []int{http.StatusNoContent},
			Uri:              fmt.Sprintf("/groups/%s/owners/$ref", *group.ID),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *GroupsClient) RemoveOwners(ctx context.Context, id string, ownerIds *[]string) error {
	for _, ownerId := range *ownerIds {
		// check for ownership before attempting deletion
		if _, err := c.GetMember(ctx, id, ownerId); err != nil {
			continue
		}
		_, err := c.BaseClient.Delete(ctx, base.DeleteHttpRequestInput{
			ValidStatusCodes: []int{http.StatusNoContent},
			Uri:              fmt.Sprintf("/groups/%s/owners/%s/$ref", id, ownerId),
		})
		if err != nil {
			return err
		}
	}
	return nil
}
