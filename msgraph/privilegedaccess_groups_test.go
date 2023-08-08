package msgraph_test

import (
	"fmt"
	"testing"

	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/msgraph"
)

func TestPrivilegedAccessGroupClient(t *testing.T) {
	c := test.NewTest(t)
	defer c.CancelFunc()

	pimGroup := testGroupsClient_Create(t, c, msgraph.Group{
		DisplayName:     utils.StringPtr("test-group"),
		MailEnabled:     utils.BoolPtr(false),
		MailNickname:    utils.StringPtr(fmt.Sprintf("test-group-%s", c.RandomString)),
		SecurityEnabled: utils.BoolPtr(true),
	})
	defer testGroupsClient_Delete(t, c, *pimGroup.ID())

	testPrivilegedAccessGroupClient_Register(t, c, *pimGroup.ID())
}

func testPrivilegedAccessGroupClient_Register(t *testing.T, c *test.Test, groupId string) {
	status, err := c.PrivilegedAccessGroupClient.Register(c.Context, groupId)
	if err != nil {
		t.Fatalf("PrivilegedAccessGroupClient.Register(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("PrivilegedAccessGroupClient.Register(): invalid status: %d", status)
	}
}
