package msgraph

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/manicminer/hamilton/odata"
	"io"
	"net/http"
)

type ChatClient struct {
	BaseClient Client
}

func NewChatClient(tenantId string) *ChatClient {
	return &ChatClient{
		BaseClient: NewClient(VersionBeta, tenantId),
	}
}

// Create creates a new chat.
func (c *ChatClient) Create(ctx context.Context, chat Chat) (*Chat, int, error) {
	var status int

	body, err := json.Marshal(chat)
	if err != nil {
		return nil, status, fmt.Errorf("json.Marshal(): %v", err)
	}

	resp, status, _, err := c.BaseClient.Post(ctx, PostHttpRequestInput{
		Body: body,
		OData: odata.Query{
			Metadata: odata.MetadataFull,
		},
		ValidStatusCodes: []int{http.StatusCreated},
		Uri: Uri{
			Entity:      "/chats",
			HasTenantId: true,
		},
	})

	if err != nil {
		return nil, status, fmt.Errorf("ChatsClient.BaseClient.Post(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var newChat Chat
	if err := json.Unmarshal(respBody, &newChat); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &newChat, status, nil
}

// Get retrieves a chat.
func (c *ChatClient) Get(ctx context.Context, id string, query odata.Query) (*Chat, int, error) {
	query.Metadata = odata.MetadataFull

	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		OData:                  query,
		ValidStatusCodes:       []int{http.StatusOK},
		Uri: Uri{
			Entity:      fmt.Sprintf("/chats/%s", id),
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("ChatsClient.BaseClient.Get(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var chat Chat
	if err := json.Unmarshal(respBody, &chat); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &chat, status, nil
}

// List returns a list of chats as Chat objects.
// To return just a lost of IDs then place the query to be Odata.Query{Select: "id"}.
func (c *ChatClient) List(ctx context.Context, userID string, query odata.Query) (*[]Chat, int, error) {
	var status int

	query.Metadata = odata.MetadataFull

	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		OData:                  query,
		ValidStatusCodes:       []int{http.StatusOK},
		Uri: Uri{
			Entity:      fmt.Sprintf("/users/%s/chats", userID),
			HasTenantId: false,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("ChatsClient.BaseClient.Get(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var chatList struct {
		Value []Chat `json:"value"`
	}
	if err := json.Unmarshal(respBody, &chatList); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &chatList.Value, status, nil

}

// Update updates a chat.
func (c *ChatClient) Update(ctx context.Context, chat Chat) (int, error) {
	var status int

	body, err := json.Marshal(chat)
	if err != nil {
		return status, fmt.Errorf("json.Marshal(): %v", err)
	}

	_, status, _, err = c.BaseClient.Patch(ctx, PatchHttpRequestInput{
		Body:                   body,
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		ValidStatusCodes:       []int{http.StatusNoContent},
		Uri: Uri{
			Entity:      fmt.Sprintf("/chats/%s", *chat.ID),
			HasTenantId: true,
		},
	})
	if err != nil {
		return status, fmt.Errorf("ChatsClient.BaseClient.Patch(): %v", err)
	}

	return status, nil
}

// ListMembers returns a list of conversationMembers in a chat.
// Note that a conversationMember can also be a group or team and not just a user.
func (c *ChatClient) ListMembers(ctx context.Context, id string, query odata.Query) (*[]ConversationMember, int, error) {
	query.Metadata = odata.MetadataFull

	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		OData:                  query,
		ValidStatusCodes:       []int{http.StatusOK},
		Uri: Uri{
			Entity:      fmt.Sprintf("/chats/%s/members", id),
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("ChatsClient.BaseClient.Get(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var chatMember []ConversationMember
	if err := json.Unmarshal(respBody, &chatMember); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &chatMember, status, nil
}

// GetMember returns a single conversationMember in a chat.
// Note that a conversationMember can also be a group or team and not just a user.
func (c *ChatClient) GetMember(ctx context.Context, chatId, memberId string) (*ConversationMember, int, error) {

	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		ValidStatusCodes:       []int{http.StatusOK},
		Uri: Uri{
			Entity:      fmt.Sprintf("/chats/%s/members/%s", chatId, memberId),
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("ChatsClient.BaseClient.Get(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var chatMember ConversationMember
	if err := json.Unmarshal(respBody, &chatMember); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &chatMember, status, nil

}

// AddMember adds a member to a chat.
// There is currently no AddMembers Function - that abstraction is required outside hamilton.
func (c *ChatClient) AddMember(ctx context.Context, chatId string, conversationMember *ConversationMember) (int, error) {
	var status int

	if conversationMember == nil {
		return status, fmt.Errorf("conversationMember is nil || No conversationMember provided")
	}

	// Populate the body
	body, err := json.Marshal(conversationMember)
	if err != nil {
		return status, fmt.Errorf("json.Marshal(): %v", err)
	}

	_, status, _, err = c.BaseClient.Post(ctx, PostHttpRequestInput{
		Body:                   body,
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		ValidStatusCodes:       []int{http.StatusCreated},
		Uri: Uri{
			Entity:      fmt.Sprintf("/chats/%s/members", chatId),
			HasTenantId: false,
		},
	})
	if err != nil {
		return status, fmt.Errorf("ChatsClient.BaseClient.Post(): %v", err)
	}

	return status, nil
}

// RemoveMember removes a member from a chat.
func (c *ChatClient) RemoveMember(ctx context.Context, chatId, memberId string) (int, error) {
	var status int

	_, status, _, err := c.BaseClient.Delete(ctx, DeleteHttpRequestInput{
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		ValidStatusCodes:       []int{http.StatusNoContent},
		Uri: Uri{
			Entity:      fmt.Sprintf("/chats/%s/members/%s", chatId, memberId),
			HasTenantId: true,
		},
	})
	if err != nil {
		return status, fmt.Errorf("ChatsClient.BaseClient.Delete(): %v", err)
	}

	return status, nil
}

// LIST MESSAGES IN CHAT
// GET MESSAGES FROM ALL CHATS FOR USER
// GET MESSAGE IN CHAT
// SEND MESSAGE IN CHAT
// GET CHAT BETWEEN USER AND APP
// MARK CHAT AS READ
// MARK CHAT AS UNREAD
// HIDE CHAT
// UNHIDE CHAT
// LIST PINNED MESSAGES
// PIN MESSAGE
// UNPIN MESSAGE
