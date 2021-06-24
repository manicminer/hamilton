package msgraph_test

import (
	"fmt"
	"testing"

	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/msgraph"
)

type SchemaExtensionsClientTest struct {
	connection   *test.Connection
	client       *msgraph.SchemaExtensionsClient
	randomstring string
}

func TestSchemaExtensionsClient(t *testing.T) {
	c := SchemaExtensionsClientTest{
		connection:   test.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomstring: test.RandomString(),
	}

	property1 := msgraph.ExtensionSchemaProperty{
		Name: utils.StringPtr("property1"),
		Type: msgraph.ExtensionSchemaPropertyDataString,
	}

	property2 := msgraph.ExtensionSchemaProperty{
		Name: utils.StringPtr("property2"),
		Type: msgraph.ExtensionSchemaPropertyDataBoolean,
	}

	targetTypes := []msgraph.ExtensionSchemaTargetType{msgraph.ExtensionSchemaTargetTypeUser}
	schemaExtension := msgraph.SchemaExtension{
		Description: utils.StringPtr("This is a description"),
		ID:          utils.StringPtr(fmt.Sprintf("schemaid%s", c.randomstring)),
		TargetTypes: &targetTypes,
		Properties:  &[]msgraph.ExtensionSchemaProperty{property1},
	}

	c.client = msgraph.NewSchemaExtensionsClient(c.connection.AuthConfig.TenantID)
	c.client.BaseClient.Authorizer = c.connection.Authorizer
	schema := testSchemaExtensionsClient_Create(t, c, schemaExtension)
	testSchemaExtensionsClient_Get(t, c, *schema.ID)

	updateExtension := msgraph.SchemaExtension{
		Properties: &[]msgraph.ExtensionSchemaProperty{property1, property2},
		ID:         schema.ID,
	}

	testSchemaExtensionsClient_Update(t, c, updateExtension)
	testSchemaExtensionsClient_Delete(t, c, *schema.ID)
}

func testSchemaExtensionsClient_Create(t *testing.T, c SchemaExtensionsClientTest, s msgraph.SchemaExtension) (schema *msgraph.SchemaExtension) {
	schema, status, err := c.client.Create(c.connection.Context, s)
	if err != nil {
		t.Fatalf("ApplicationsClient.Create(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("SchemaExtensionsClient.Create(): invalid status: %d", status)
	}
	if schema == nil {
		t.Fatal("SchemaExtensionsClient.Create(): schema was nil")
	}
	if schema.ID == nil {
		t.Fatal("SchemaExtensionsClient.Create(): schema.ID was nil")
	}
	return
}

func testSchemaExtensionsClient_Get(t *testing.T, c SchemaExtensionsClientTest, id string) (schema *msgraph.SchemaExtension) {
	schema, status, err := c.client.Get(c.connection.Context, id)
	if err != nil {
		t.Fatalf("SchemaExtensionsClient.Get(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("SchemaExtensionsClient.Get(): invalid status: %d", status)
	}
	if schema == nil {
		t.Fatal("SchemaExtensionsClient.Get(): schema was nil")
	}
	return
}

func testSchemaExtensionsClient_Update(t *testing.T, c SchemaExtensionsClientTest, s msgraph.SchemaExtension) {
	status, err := c.client.Update(c.connection.Context, s)
	if err != nil {
		t.Fatalf("SchemaExtensionsClient.Update(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("SchemaExtensionsClient.Update(): invalid status: %d", status)
	}
}

func testSchemaExtensionsClient_Delete(t *testing.T, c SchemaExtensionsClientTest, id string) {
	status, err := c.client.Delete(c.connection.Context, id)
	if err != nil {
		t.Fatalf("SchemaExtensionsClient.Delete(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("SchemaExtensionsClient.Delete(): invalid status: %d", status)
	}
}
