package msgraph_test

import (
	"testing"

	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/msgraph"
)

func TestCustomSecurityAttributeDefinitionClient(t *testing.T) {
	c := test.NewTest(t)
	defer c.CancelFunc()

	attributeSet, _, err := c.AttributeSetClient.Create(
		c.Context,
		msgraph.AttributeSet{
      Description: utils.StringPtr("custom_security_attributes test"),
			ID: utils.StringPtr(c.RandomString),
		},
	)
	if err != nil {
		t.Fatalf("AttributeSetClient.Create(): %v", err)
	}

	customSecurityAttributeDefinition := testCustomSecurityAttributeDefinitionClientCreate(
		t,
		c,
		msgraph.CustomSecurityAttributeDefinition{
			AttributeSet:            attributeSet.ID,
			Description:             utils.StringPtr("test description"),
			IsCollection:            utils.BoolPtr(false),
			IsSearchable:            utils.BoolPtr(false),
			Name:                    utils.StringPtr(c.RandomString),
			Status:                  utils.StringPtr("Available"),
			Type:                    utils.StringPtr("Boolean"),
			UsePreDefinedValuesOnly: utils.BoolPtr(false),
		},
	)

	testCustomSecurityAttributeDefinitionClientGet(t, c, *customSecurityAttributeDefinition.ID)
	testCustomSecurityAttributeDefinitionClientUpdate(
		t,
		c,
		msgraph.CustomSecurityAttributeDefinition{
      ID: customSecurityAttributeDefinition.ID,
			Description: utils.StringPtr("updated test description"),
		},
	)

	testCustomSecurityAttributeDefinitionClientList(t, c)
}

func testCustomSecurityAttributeDefinitionClientCreate(t *testing.T, c *test.Test, csad msgraph.CustomSecurityAttributeDefinition) *msgraph.CustomSecurityAttributeDefinition {

	customSecurityAttributeDefinition, status, err := c.CustomSecurityAttributeDefinitionClient.Create(c.Context, csad)
	if err != nil {
		t.Fatalf("CustomSecurityAttributeDefinitionClient.Create(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("CustomSecurityAttributeDefinitionClient.Create(): invalid status:%d", status)
	}
	if customSecurityAttributeDefinition == nil {
		t.Fatalf("CustomSecurityAttributeDefinition.Create(): customSecurityAttributeDefinition was nil")
	}
	if customSecurityAttributeDefinition.ID == nil {
		t.Fatalf("CustomSecurityAttributeDefinitionClient.Create(): customSecurityAttributeDefinition.ID was nil")
	}

	return customSecurityAttributeDefinition
}

func testCustomSecurityAttributeDefinitionClientGet(t *testing.T, c *test.Test, id string) *msgraph.CustomSecurityAttributeDefinition {
	customSecurityAttributeDefinition, status, err := c.CustomSecurityAttributeDefinitionClient.Get(c.Context, id, odata.Query{})
	if err != nil {
		t.Fatalf("CustomSecurityAttributeDefinitionClient.Get(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("CustomSecurityAttributeDefinition.Client.Get(): invalid status: %d", status)
	}
	if customSecurityAttributeDefinition == nil {
		t.Fatalf("CustomSecurityAttributeDefinitionClient.Get(): customSecurityAttributeDefinition was nil")
	}

	return customSecurityAttributeDefinition
}

func testCustomSecurityAttributeDefinitionClientList(t *testing.T, c *test.Test) *[]msgraph.CustomSecurityAttributeDefinition {
	customSecurityAttributeDefinitions, _, err := c.CustomSecurityAttributeDefinitionClient.List(
		c.Context,
		odata.Query{Top: 10},
	)
	if err != nil {
		t.Fatalf("CustomSecurityAttributeDefinitionClient.List(): %v", err)
	}
	if customSecurityAttributeDefinitions == nil {
		t.Fatalf("CustomSecurityAttributeDefinitionClient.List(): customSecurityAttributeDefinitions was nil")
	}

	return customSecurityAttributeDefinitions
}

func testCustomSecurityAttributeDefinitionClientUpdate(t *testing.T, c *test.Test, csad msgraph.CustomSecurityAttributeDefinition) {
	status, err := c.CustomSecurityAttributeDefinitionClient.Update(c.Context, csad)
	if err != nil {
		t.Fatalf("CustomSecurityAttributeDefinitionClient.Update(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("CustomSecurityAttributeDefinitionClient.Update(): invalid status: %d", status)
	}
}

