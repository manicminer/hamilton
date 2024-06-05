package msgraph

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type ChatClient struct {
	BaseClient Client
}

func NewChatClient() *ChatClient {
	return &ChatClient{
		BaseClient: NewClient(Version10),
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
			Entity: "/chats",
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
			Entity: fmt.Sprintf("/chats/%s", id),
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
			Entity: fmt.Sprintf("/users/%s/chats", userID),
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
			Entity: fmt.Sprintf("/chats/%s", *chat.ID),
		},
	})
	if err != nil {
		return status, fmt.Errorf("ChatsClient.BaseClient.Patch(): %v", err)
	}

	return status, nil
}
