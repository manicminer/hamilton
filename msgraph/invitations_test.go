package msgraph_test

import (
	"fmt"
	"testing"

	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/msgraph"
)

func TestInvitationsClient(t *testing.T) {
	c := test.NewTest()

	testInvitationsClient_Create(t, c, msgraph.Invitation{
		InvitedUserDisplayName:  utils.StringPtr("test-user-invited"),
		InvitedUserEmailAddress: utils.StringPtr(fmt.Sprintf("test-user-%s@test.com", c.RandomString)),
		InviteRedirectURL:       utils.StringPtr(fmt.Sprintf("https://myapp-%s.contoso.com", c.RandomString)),
	})
}

func testInvitationsClient_Create(t *testing.T, c *test.Test, i msgraph.Invitation) (invitation *msgraph.Invitation) {
	invitation, status, err := c.InvitationsClient.Create(c.Connection.Context, i)
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
