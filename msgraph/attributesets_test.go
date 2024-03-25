package msgraph_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/msgraph"
)

func TestAttributeSetsClient(t *testing.T) {
	c := test.NewTest(t)
	defer c.CancelFunc()

	randomString := test.RandomString()

	attributeSet := testAttributeSetsClient_Create(t, c, msgraph.AttributeSet{
		ID:                  utils.StringPtr(fmt.Sprintf("testing%s", randomString)),
		Description:         utils.StringPtr("This is a test attribute set"),
		MaxAttributesPerSet: utils.Int32Ptr(20),
	})

	testAttributeSetsClient_Get(t, c, *attributeSet.ID)

	attributeSet.Description = utils.StringPtr("This is an updated test attribute set")
	attributeSet.MaxAttributesPerSet = utils.Int32Ptr(25)

	testAttributeSetsClient_Update(t, c, *attributeSet)

	testAttributeSetsClient_List(t, c)
}

func testAttributeSetsClient_Create(t *testing.T, c *test.Test, as msgraph.AttributeSet) (attributeSet *msgraph.AttributeSet) {
	attributeSet, status, err := c.AttributeSetsClient.Create(c.Context, as)
	if err != nil {
		t.Fatalf("AttributeSetsClient.Create(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AttributeSetsClient.Create(): invalid status: %d", status)
	}
	if attributeSet == nil {
		t.Fatal("AttributeSetsClient.Create(): attributeSet was nil")
	}
	if attributeSet.ID == nil {
		t.Fatal("AttributeSetsClient.Create(): attributeSet.ID was nil")
	}
	return
}

func testAttributeSetsClient_Update(t *testing.T, c *test.Test, as msgraph.AttributeSet) {
	status, err := c.AttributeSetsClient.Update(c.Context, as)
	if err != nil {
		t.Fatalf("AttributeSetsClient.Update(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AttributeSetsClient.Update(): invalid status: %d", status)
	}
}

func testAttributeSetsClient_Get(t *testing.T, c *test.Test, id string) (attributeSet *msgraph.AttributeSet) {
	attributeSet, status, err := c.AttributeSetsClient.Get(c.Context, id, odata.Query{})
	if err != nil {
		t.Fatalf("AttributeSetsClient.Get(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AttributeSetsClient.Get(): invalid status: %d", status)
	}
	if attributeSet == nil {
		t.Fatal("AttributeSetsClient.Get(): attributeSet was nil")
	}
	return
}

func testAttributeSetsClient_List(t *testing.T, c *test.Test) (attributeSets *[]msgraph.AttributeSet) {
	attributeSets, _, err := c.AttributeSetsClient.List(c.Context, odata.Query{})
	if err != nil {
		t.Fatalf("AttributeSetsClient.List(): %v", err)
	}
	if attributeSets == nil {
		t.Fatal("AttributeSetsClient.List(): attributeSet was nil")
	}
	return
}
