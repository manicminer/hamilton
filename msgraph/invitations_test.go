package msgraph_test

import (
	"fmt"
	"testing"

	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/msgraph"
)

type InvitationsClientTest struct {
	connection   *test.Connection
	client       *msgraph.InvitationsClient
	randomString string
}

func TestInvitationsClient(t *testing.T) {
	rs := test.RandomString()
	c := InvitationsClientTest{
		connection:   test.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: rs,
	}
	c.client = msgraph.NewInvitationsClient(c.connection.AuthConfig.TenantID)
	c.client.BaseClient.Authorizer = c.connection.Authorizer

	testInvitationsClient_Create(t, c, msgraph.Invitation{
		InvitedUserDisplayName:  utils.StringPtr("test-user-invited"),
		InvitedUserEmailAddress: utils.StringPtr(fmt.Sprintf("test-user-%s@test.com", c.randomString)),
		InviteRedirectURL:       utils.StringPtr(fmt.Sprintf("https://myapp-%s.contoso.com", c.randomString)),
	})
}

func testInvitationsClient_Create(t *testing.T, c InvitationsClientTest, i msgraph.Invitation) (invitation *msgraph.Invitation) {
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
