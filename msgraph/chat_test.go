package msgraph_test

import (
	"fmt"
	"github.com/manicminer/hamilton/internal/utils"
	"testing"

	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/msgraph"
	"github.com/manicminer/hamilton/odata"
)

func TestChatClient(t *testing.T) {
	c := test.NewTest(t)
	defer c.CancelFunc()

	self := testDirectoryObjectsClient_Get(t, c, c.Claims.ObjectId)

	// To create a chat two users need to be assigned to the chat.
	// An owner needs to be assigned
	user1 := testUsersClient_Create(t, c, msgraph.User{
		AccountEnabled:    utils.BoolPtr(true),
		DisplayName:       utils.StringPtr("test-user1"),
		MailNickname:      utils.StringPtr(fmt.Sprintf("test-user-%s", c.RandomString)),
		UserPrincipalName: utils.StringPtr(fmt.Sprintf("test-user-%s@%s", c.RandomString, c.Connections["default"].DomainName)),
		PasswordProfile: &msgraph.UserPasswordProfile{
			Password: utils.StringPtr(fmt.Sprintf("IrPa55w0rd%s", c.RandomString)),
		},
	})
	user2 := testUsersClient_Create(t, c, msgraph.User{
		AccountEnabled:    utils.BoolPtr(true),
		DisplayName:       utils.StringPtr("test-user2"),
		MailNickname:      utils.StringPtr(fmt.Sprintf("test-user-%s", c.RandomString)),
		UserPrincipalName: utils.StringPtr(fmt.Sprintf("test-user-%s@%s", c.RandomString, c.Connections["default"].DomainName)),
		PasswordProfile: &msgraph.UserPasswordProfile{
			Password: utils.StringPtr(fmt.Sprintf("IrPa55w0rd%s", c.RandomString)),
		},
	})

	// Check that a group chat and a OneOnOne chat can be created
	newChat := msgraph.Chat{
		Topic:    utils.StringPtr(fmt.Sprintf("test-chat-%s", c.RandomString)),
		ChatType: utils.StringPtr(msgraph.ChatTypeGroup),
		Members: &[]msgraph.ConversationMember{
			{
				ID:    user1.Id,
				Roles: &[]string{"owner"},
			},
			{
				ID:    user2.Id,
				Roles: &[]string{"owner"},
			},
			{
				ID:    self.Id,
				Roles: &[]string{"owner"},
			},
		},
	}

	chat := testChatClient_Create(t, c, newChat)
	testChatClient_Get(t, c, *chat.ID)
	testChatClient_list(t, c, *self.Id)

}

func testChatClient_Create(t *testing.T, c *test.Test, newChat msgraph.Chat) (chat *msgraph.Chat) {
	chat, status, err := c.ChatClient.Create(c.Context, newChat)
	if err != nil {
		t.Fatalf("ChatClient.Create(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ChatClient.Create(): invalid status: %d", status)
	}
	if chat == nil {
		t.Fatal("ChatClient.Create(): chat was nil")
	}
	return
}

func testChatClient_Get(t *testing.T, c *test.Test, id string) (chat *msgraph.Chat) {
	chat, status, err := c.ChatClient.Get(c.Context, id, odata.Query{})
	if err != nil {
		t.Fatalf("ChatClient.Get(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ChatClient.Get(): invalid status: %d", status)
	}
	if chat == nil {
		t.Fatal("ChatClient.Get(): chat was nil")
	}
	return
}

func testChatClient_list(t *testing.T, c *test.Test, userID string) (chats *[]msgraph.Chat) {
	chats, _, err := c.ChatClient.List(c.Context, userID, odata.Query{Top: 10})
	if err != nil {
		t.Fatalf("ChatClient.List(): %v", err)
	}
	if chats == nil {
		t.Fatal("ChatClient.List(): chats was nil")
	}
	return
}
