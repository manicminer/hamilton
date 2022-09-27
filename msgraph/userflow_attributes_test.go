package msgraph_test

import (
	"testing"

	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/msgraph"
	"github.com/manicminer/hamilton/odata"
)

func TestUserFlowAttributesClient(t *testing.T) {
	c := test.NewTest(t)
	defer c.CancelFunc()

	userflowAttribute := testUserFlowAttributesClient_Create(t, c, msgraph.UserFlowAttribute{
		ID:                    utils.StringPtr("test attribute"),
		DisplayName:           utils.StringPtr("test attribute"),
		UserFlowAttributeType: utils.StringPtr("custom"),
		DataType:              utils.StringPtr("string"),
	})
	testUserFlowAttributesClient_Get(t, c, *userflowAttribute.ID)
	userflowAttribute.DisplayName = utils.StringPtr("updated test attribute")
	testUserFlowAttributesClient_Update(t, c, *userflowAttribute)
	testUserFlowAttributesClient_List(t, c)
	testGroupsClient_Delete(t, c, *userflowAttribute.ID)
}

func testUserFlowAttributesClient_Create(t *testing.T, c *test.Test, u msgraph.UserFlowAttribute) *msgraph.UserFlowAttribute {
	userflowAttribute, status, err := c.UserFlowAttributesClient.Create(c.Context, u)
	if err != nil {
		t.Fatalf("UserFlowAttributeclient.Create(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("UserFlowAttributesClient.Create(): invalid status: %d", status)
	}
	if userflowAttribute == nil {
		t.Fatal("UserFlowAttributesClient.Create(): userflowAttribute was nil")
	}
	if userflowAttribute.ID == nil {
		t.Fatal("UserFlowAttributesClient.Create(): userflowAttribute.ID was nil")
	}
	return userflowAttribute
}

func testUserFlowAttributesClient_Get(t *testing.T, c *test.Test, id string) *msgraph.UserFlowAttribute {
	userflowAttribute, status, err := c.UserFlowAttributesClient.Get(c.Context, id, odata.Query{})
	if err != nil {
		t.Fatalf("UserFlowAttributesClient.Get(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("UserFlowAttributesClient.Get(): invalid status: %d", status)
	}
	if userflowAttribute == nil {
		t.Fatal("UserFlowAttributesClient.Get(): userflowAttribute was nil")
	}
	return userflowAttribute
}

func testUserFlowAttributesClient_List(t *testing.T, c *test.Test) *[]msgraph.UserFlowAttribute {
	userflowAttributes, _, err := c.UserFlowAttributesClient.List(c.Context, odata.Query{Top: 10})
	if err != nil {
		t.Fatalf("UserFlowAttributesClient.List(): %v", err)
	}
	if userflowAttributes == nil {
		t.Fatal("UserFlowAttributesClient.List(): userflowAttributes was nil")
	}
	return userflowAttributes
}

func testUserFlowAttributesClient_Update(t *testing.T, c *test.Test, u msgraph.UserFlowAttribute) {
	status, err := c.UserFlowAttributesClient.Update(c.Context, u)
	if err != nil {
		t.Fatalf("UserFlowAttributesClient.Update(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("UserFlowAttributesClient.Update(): invalid status: %d", status)
	}
}
