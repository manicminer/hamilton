package msgraph_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/msgraph"
)

func TestB2CUserFlowClient(t *testing.T) {
	c := test.NewTest(t)
	defer c.CancelFunc()

	userflow := testB2CUserFlowClient_Create(t, c, msgraph.B2CUserFlow{
		ID:                  utils.StringPtr(fmt.Sprintf("test_b2c_user_flow_%s", c.RandomString)),
		UserFlowType:        utils.StringPtr("signUp"),
		UserFlowTypeVersion: utils.Float32Ptr(1.0),
	})
	testB2CUserFlowClient_Get(t, c, *userflow.ID)
	testB2CUserFlowClient_Update(t, c, msgraph.B2CUserFlow{
		ID:                 userflow.ID,
		DefaultLanguageTag: utils.StringPtr("en"),
	})
	testB2CUserFlowClient_List(t, c)
	testB2CUserFlowClient_Delete(t, c, *userflow.ID)
}

func testB2CUserFlowClient_Create(t *testing.T, c *test.Test, u msgraph.B2CUserFlow) *msgraph.B2CUserFlow {
	userflow, status, err := c.B2CUserFlowClient.Create(c.Context, u)
	if err != nil {
		t.Fatalf("B2CUserFlowclient.Create(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("B2CUserFlowClient.Create(): invalid status: %d", status)
	}
	if userflow == nil {
		t.Fatal("B2CUserFlowClient.Create(): userflow was nil")
	}
	if userflow.ID == nil {
		t.Fatal("B2CUserFlowClient.Create(): userflow.ID was nil")
	}
	return userflow
}

func testB2CUserFlowClient_Get(t *testing.T, c *test.Test, id string) *msgraph.B2CUserFlow {
	userflow, status, err := c.B2CUserFlowClient.Get(c.Context, id, odata.Query{})
	if err != nil {
		t.Fatalf("B2CUserFlowClient.Get(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("B2CUserFlowClient.Get(): invalid status: %d", status)
	}
	if userflow == nil {
		t.Fatal("B2CUserFlowClient.Get(): userflow was nil")
	}
	return userflow
}

func testB2CUserFlowClient_List(t *testing.T, c *test.Test) *[]msgraph.B2CUserFlow {
	userflows, _, err := c.B2CUserFlowClient.List(c.Context, odata.Query{Top: 10})
	if err != nil {
		t.Fatalf("B2CUserFlowClient.List(): %v", err)
	}
	if userflows == nil {
		t.Fatal("B2CUserFlowClient.List(): userflows was nil")
	}
	return userflows
}

func testB2CUserFlowClient_Update(t *testing.T, c *test.Test, u msgraph.B2CUserFlow) {
	status, err := c.B2CUserFlowClient.Update(c.Context, u)
	if err != nil {
		t.Fatalf("B2CUserFlowClient.Update(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("B2CUserFlowClient.Update(): invalid status: %d", status)
	}
}

func testB2CUserFlowClient_Delete(t *testing.T, c *test.Test, id string) {
	status, err := c.B2CUserFlowClient.Delete(c.Context, id)
	if err != nil {
		t.Fatalf("B2CUserFlowClient.Delete(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("B2CUserFlowClient.Delete(): invalid status: %d", status)
	}
}
