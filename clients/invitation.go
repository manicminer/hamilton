package clients

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/manicminer/hamilton/base"
	"github.com/manicminer/hamilton/models"
)

// InvitationsClient performs operations on Invitations.
type InvitationsClient struct {
	BaseClient base.Client
}

// NewInvitationsClient returns a new InvitationsClient.
func NewInvitationsClient(tenantId string) *InvitationsClient {
	return &InvitationsClient{
		BaseClient: base.NewClient(base.VersionBeta, tenantId),
	}
}

// Create creates a new Invitation.
func (c *InvitationsClient) Create(ctx context.Context, invitation models.Invitation) (*models.Invitation, int, error) {
	var status int
	body, err := json.Marshal(invitation)
	if err != nil {
		return nil, status, err
	}
	resp, status, _, err := c.BaseClient.Post(ctx, base.PostHttpRequestInput{
		Body:             body,
		ValidStatusCodes: []int{http.StatusCreated},
		Uri: base.Uri{
			Entity:      "/invitations",
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, err
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	var newInvitation models.Invitation
	if err := json.Unmarshal(respBody, &newInvitation); err != nil {
		return nil, status, err
	}
	return &newInvitation, status, nil
}
