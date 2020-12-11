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
	BaseClient base.Client
}

func NewGroupsClient(authorizer auth.Authorizer, tenantId string) *GroupsClient {
	return &GroupsClient{
		BaseClient: base.NewClient(authorizer, base.DefaultEndpoint, tenantId, base.VersionBeta),
	}
}

func (c *GroupsClient) List(ctx context.Context, filter string) (*[]models.Group, int, error) {
	params := url.Values{}
	if filter != "" {
		params.Add("$filter", filter)
	}
	resp, status, err := c.BaseClient.Get(ctx, base.GetHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		Uri:              fmt.Sprintf("/groups?%s", params.Encode()),
	})
	if err != nil {
		return nil, status, err
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	var data struct {
		Groups []models.Group `json:"value"`
	}
	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, status, err
	}
	return &data.Groups, status, nil
}

func (c *GroupsClient) Create(ctx context.Context, group models.Group) (*models.Group, int, error) {
	var status int
	body, err := json.Marshal(group)
	if err != nil {
		return nil, status, err
	}
	resp, status, err := c.BaseClient.Post(ctx, base.PostHttpRequestInput{
		Body:             body,
		ValidStatusCodes: []int{http.StatusCreated},
		Uri:              "/groups",
	})
	if err != nil {
		return nil, status, err
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	var newGroup models.Group
	if err := json.Unmarshal(respBody, &newGroup); err != nil {
		return nil, status, err
	}
	return &newGroup, status, nil
}

func (c *GroupsClient) Get(ctx context.Context, id string) (*models.Group, int, error) {
	resp, status, err := c.BaseClient.Get(ctx, base.GetHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		Uri:              fmt.Sprintf("/groups/%s", id),
	})
	if err != nil {
		return nil, status, err
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	var group models.Group
	if err := json.Unmarshal(respBody, &group); err != nil {
		return nil, status, err
	}
	return &group, status, nil
}

func (c *GroupsClient) Update(ctx context.Context, group models.Group) (int, error) {
	var status int
	body, err := json.Marshal(group)
	if err != nil {
		return status, err
	}
	_, status, err = c.BaseClient.Patch(ctx, base.PatchHttpRequestInput{
		Body:             body,
		ValidStatusCodes: []int{http.StatusNoContent},
		Uri:              fmt.Sprintf("/groups/%s", *group.ID),
	})
	if err != nil {
		return status, err
	}
	return status, nil
}

func (c *GroupsClient) Delete(ctx context.Context, id string) (int, error) {
	_, status, err := c.BaseClient.Delete(ctx, base.DeleteHttpRequestInput{
		ValidStatusCodes: []int{http.StatusNoContent},
		Uri:              fmt.Sprintf("/groups/%s", id),
	})
	if err != nil {
		return status, err
	}
	return status, nil
}

func (c *GroupsClient) ListMembers(ctx context.Context, id string) (*[]string, int, error) {
	resp, status, err := c.BaseClient.Get(ctx, base.GetHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		Uri:              fmt.Sprintf("/groups/%s/members?$select=id", id),
	})
	if err != nil {
		return nil, status, err
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
		return nil, status, err
	}
	ret := make([]string, len(data.Members))
	for i, v := range data.Members {
		ret[i] = v.Id
	}
	return &ret, status, nil
}

func (c *GroupsClient) GetMember(ctx context.Context, groupId, memberId string) (*string, int, error) {
	resp, status, err := c.BaseClient.Get(ctx, base.GetHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		Uri:              fmt.Sprintf("/groups/%s/members/%s/$ref?$select=id,url", groupId, memberId),
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

func (c *GroupsClient) AddMembers(ctx context.Context, group *models.Group) (int, error) {
	var status int
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
			return status, err
		}
		_, status, err = c.BaseClient.Patch(ctx, base.PatchHttpRequestInput{
			Body:             body,
			ValidStatusCodes: []int{http.StatusNoContent},
			Uri:              fmt.Sprintf("/groups/%s", *group.ID),
		})
		if err != nil {
			return status, err
		}
	}
	return status, nil
}

func (c *GroupsClient) RemoveMembers(ctx context.Context, id string, memberIds *[]string) (int, error) {
	var status int
	for _, memberId := range *memberIds {
		// check for membership before attempting deletion
		if _, status, err := c.GetMember(ctx, id, memberId); err != nil {
			if status == http.StatusNotFound {
				continue
			}
			return status, err
		}
		_, status, err := c.BaseClient.Delete(ctx, base.DeleteHttpRequestInput{
			ValidStatusCodes: []int{http.StatusNoContent},
			Uri:              fmt.Sprintf("/groups/%s/members/%s/$ref", id, memberId),
		})
		if err != nil {
			return status, err
		}
	}
	return status, nil
}

func (c *GroupsClient) ListOwners(ctx context.Context, id string) (*[]string, int, error) {
	resp, status, err := c.BaseClient.Get(ctx, base.GetHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		Uri:              fmt.Sprintf("/groups/%s/owners?$select=id", id),
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

func (c *GroupsClient) GetOwner(ctx context.Context, groupId, ownerId string) (*string, int, error) {
	resp, status, err := c.BaseClient.Get(ctx, base.GetHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		Uri:              fmt.Sprintf("/groups/%s/owners/%s/$ref?$select=id,url", groupId, ownerId),
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

func (c *GroupsClient) AddOwners(ctx context.Context, group *models.Group) (int, error) {
	var status int
	for _, owner := range *group.Owners {
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
			Uri:              fmt.Sprintf("/groups/%s/owners/$ref", *group.ID),
		})
		if err != nil {
			return status, err
		}
	}
	return status, nil
}

func (c *GroupsClient) RemoveOwners(ctx context.Context, id string, ownerIds *[]string) (int, error) {
	var status int
	for _, ownerId := range *ownerIds {
		// check for ownership before attempting deletion
		if _, status, err := c.GetOwner(ctx, id, ownerId); err != nil {
			if status == http.StatusNotFound {
				continue
			}
			return status, err
		}
		_, status, err := c.BaseClient.Delete(ctx, base.DeleteHttpRequestInput{
			ValidStatusCodes: []int{http.StatusNoContent},
			Uri:              fmt.Sprintf("/groups/%s/owners/%s/$ref", id, ownerId),
		})
		if err != nil {
			return status, err
		}
	}
	return status, nil
}
