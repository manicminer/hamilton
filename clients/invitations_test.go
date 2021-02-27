package clients_test

import (
	"fmt"
	"testing"

	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/clients"
	"github.com/manicminer/hamilton/clients/internal"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/models"
)

type InvitationsClientTest struct {
	connection   *internal.Connection
	client       *clients.InvitationsClient
	randomString string
}

func TestInvitationsClient(t *testing.T) {
	rs := internal.RandomString()
	c := InvitationsClientTest{
		connection:   internal.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: rs,
	}
	c.client = clients.NewInvitationsClient(c.connection.AuthConfig.TenantID)
	c.client.BaseClient.Authorizer = c.connection.Authorizer

	testInvitationsClient_Create(t, c, models.Invitation{
		InvitedUserEmailAddress: utils.StringPtr(fmt.Sprintf("test-user-%s@test.com", c.randomString)),
		InviteRedirectURL:       utils.StringPtr(fmt.Sprintf("https://myapp-%s.contoso.com", c.randomString)),
	})
}

func testInvitationsClient_Create(t *testing.T, c InvitationsClientTest, i models.Invitation) (invitation *models.Invitation) {
	invitation, status, err := c.client.Create(c.connection.Context, i)
	if err != nil {
		t.Fatalf("InvitationsClient.Create(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("InvitationsClient.Create(): invalid status: %d", status)
	}
	if invitation == nil {
		t.Fatal("InvitationsClient.Create(): invitation was nil")
	}
	if invitation.ID == nil {
		t.Fatal("InvitationsClient.Create(): invitation.ID was nil")
	}
	return
}
