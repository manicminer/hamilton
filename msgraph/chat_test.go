package msgraph_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/msgraph"
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
		MailNickname:      utils.StringPtr(fmt.Sprintf("test-user1-%s", c.RandomString)),
		UserPrincipalName: utils.StringPtr(fmt.Sprintf("test-user1-%s@%s", c.RandomString, c.Connections["default"].DomainName)),
		PasswordProfile: &msgraph.UserPasswordProfile{
			Password: utils.StringPtr(fmt.Sprintf("IrPa55w0rd%s", c.RandomString)),
		},
	})
	defer testUsersClient_Delete(t, c, *user1.ID())

	user2 := testUsersClient_Create(t, c, msgraph.User{
		AccountEnabled:    utils.BoolPtr(true),
		DisplayName:       utils.StringPtr("test-user2"),
		MailNickname:      utils.StringPtr(fmt.Sprintf("test-user2-%s", c.RandomString)),
		UserPrincipalName: utils.StringPtr(fmt.Sprintf("test-user2-%s@%s", c.RandomString, c.Connections["default"].DomainName)),
		PasswordProfile: &msgraph.UserPasswordProfile{
			Password: utils.StringPtr(fmt.Sprintf("IrPa55w0rd%s", c.RandomString)),
		},
	})
	defer testUsersClient_Delete(t, c, *user2.ID())

	// Check that a group chat and a OneOnOne chat can be created
	newChat := msgraph.Chat{
		Topic:    utils.StringPtr(fmt.Sprintf("test-chat-%s", c.RandomString)),
		ChatType: utils.StringPtr(msgraph.ChatTypeGroup),
		Members: &[]msgraph.ConversationMember{
			{
				ODataType: utils.StringPtr(msgraph.TypeConversationMember),
				User:      utils.StringPtr(fmt.Sprintf("https://graph.microsoft.com/v1.0/users('%s')", *user1.Id)),
				Roles:     &[]string{"owner"},
			},
			{
				ODataType: utils.StringPtr(msgraph.TypeConversationMember),
				User:      utils.StringPtr(fmt.Sprintf("https://graph.microsoft.com/v1.0/users('%s')", *user2.Id)),
				Roles:     &[]string{"owner"},
			},
			{
				ODataType: utils.StringPtr(msgraph.TypeConversationMember),
				User:      utils.StringPtr(fmt.Sprintf("https://graph.microsoft.com/v1.0/users('%s')", *self.Id)),
				Roles:     &[]string{"owner"},
			},
		},
	}

	chat := testChatClient_Create(t, c, newChat)
	defer testChatClient_Delete(t, c, *chat.ID)
	testChatClient_Get(t, c, *chat.ID)
	testChatClient_List(t, c, *self.Id)

	chat.Topic = utils.StringPtr(fmt.Sprintf("test-chat-archived-%s", c.RandomString))
	chat.Viewpoint.IsHidden = utils.BoolPtr(true)
	testChatClient_Update(t, c, *chat)
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

func testChatClient_List(t *testing.T, c *test.Test, userID string) (chats *[]msgraph.Chat) {
	chats, _, err := c.ChatClient.List(c.Context, userID, odata.Query{Top: 10})
	if err != nil {
		t.Fatalf("ChatClient.List(): %v", err)
	}
	if chats == nil {
		t.Fatal("ChatClient.List(): chats was nil")
	}
	return
}

func testChatClient_Update(t *testing.T, c *test.Test, chat msgraph.Chat) (updatedChat *msgraph.Chat) {
	chat.Topic = utils.StringPtr(fmt.Sprintf("test-chat-%s", c.RandomString))
	status, err := c.ChatClient.Update(c.Context, chat)
	if err != nil {
		t.Fatalf("ChatClient.Update(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ChatClient.Update(): invalid status: %d", status)
	}
	return
}

func testChatClient_Delete(t *testing.T, c *test.Test, chatId string) {
	status, err := c.ChatClient.Delete(c.Context, chatId)
	if err != nil {
		t.Fatalf("ChatClient.Delete(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ChatClient.Delete(): invalid status: %d", status)
	}
	return
}
