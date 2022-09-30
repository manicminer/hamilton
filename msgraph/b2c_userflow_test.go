package msgraph_test

import (
	"testing"

	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/msgraph"
	"github.com/manicminer/hamilton/odata"
)

func TestB2CUserFlowClient(t *testing.T) {
	c := test.NewTest(t)
	defer c.CancelFunc()

	userflow := testB2CUserFlowClient_Create(t, c, msgraph.B2CUserFlow{
		ID:                  utils.StringPtr("test b2c user flow"),
		UserFlowType:        utils.StringPtr("signup"),
		UserFlowTypeVersion: utils.Float32Ptr(3.0),
	})
	testB2CUserFlowClient_Get(t, c, *userflow.ID)
	userflow.DefaultLanguageTag = utils.StringPtr("en")
	testB2CUserFlowClient_Update(t, c, *userflow)
	testB2CUserFlowClient_List(t, c)
	testGroupsClient_Delete(t, c, *userflow.ID)

	attr := testB2CUserFlowClient_CreateAttribute(t, c)
	testB2CUserFlowClient_AssignAttribute(t, c, *userflow.ID, &msgraph.UserFlowAttributeAssignment{
		UserInputType: utils.StringPtr(msgraph.UserInpuTypeTextBox),
		UserAttribute: attr,
		DisplayName:   utils.StringPtr("test assignment"),
	})

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

func testB2CUserFlowClient_CreateAttribute(t *testing.T, c *test.Test) *msgraph.UserFlowAttribute {
	attr := msgraph.UserFlowAttribute{
		DisplayName:           utils.StringPtr("testattr"),
		UserFlowAttributeType: utils.StringPtr("custom"),
		DataType:              utils.StringPtr(msgraph.UserflowAttributeDataTypeString),
		Description:           utils.StringPtr("test attr description"),
	}
	resp, _, err := c.UserFlowAttributesClient.Create(c.Context, attr)
	if err != nil {
		t.Fatalf("failed to create user flow attribute. err: %s", err)
	}
	return resp
}

func testB2CUserFlowClient_AssignAttribute(t *testing.T, c *test.Test, id string, u *msgraph.UserFlowAttributeAssignment) *msgraph.UserFlowAttributeAssignment {
	resp, _, err := c.B2CUserFlowClient.AssignAttribute(c.Context, id, *u)
	if err != nil {
		t.Fatalf("failed to assign user flow attribute. err: %s", err)
	}
	return resp
}
