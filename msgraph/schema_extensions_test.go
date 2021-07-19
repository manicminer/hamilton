package msgraph_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/msgraph"
	"github.com/manicminer/hamilton/odata"
)

type SchemaExtensionsClientTest struct {
	connection   *test.Connection
	client       *msgraph.SchemaExtensionsClient
	randomString string
}

type ExtensionProperties struct {
	Property1 *string `json:"property1,omitempty"`
	Property2 *bool   `json:"property2,omitempty"`
}

func (e *ExtensionProperties) UnmarshalJSON(data []byte) error {
	type ep ExtensionProperties
	e2 := (*ep)(e)
	return json.Unmarshal(data, e2)
}

func TestSchemaExtensionsClient(t *testing.T) {
	rs := test.RandomString()
	c := SchemaExtensionsClientTest{
		connection:   test.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: rs,
	}
	c.client = msgraph.NewSchemaExtensionsClient(c.connection.AuthConfig.TenantID)
	c.client.BaseClient.Authorizer = c.connection.Authorizer

	u := UsersClientTest{
		connection:   test.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: rs,
	}
	u.client = msgraph.NewUsersClient(u.connection.AuthConfig.TenantID)
	u.client.BaseClient.Authorizer = u.connection.Authorizer

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
		ID:          utils.StringPtr("testschema"),
		TargetTypes: &targetTypes,
		Properties:  &[]msgraph.ExtensionSchemaProperty{property1},
		Status:      msgraph.SchemaExtensionStatusInDevelopment,
	}

	schema := testSchemaExtensionsClient_Create(t, c, schemaExtension)
	testSchemaExtensionsClient_Get(t, c, *schema.ID)

	updateExtension := msgraph.SchemaExtension{
		Properties: &[]msgraph.ExtensionSchemaProperty{property1, property2},
		ID:         schema.ID,
	}

	testSchemaExtensionsClient_Update(t, c, updateExtension)

	time.Sleep(10 * time.Second)

	user := testUsersClient_Create(t, u, msgraph.User{
		AccountEnabled:    utils.BoolPtr(true),
		DisplayName:       utils.StringPtr("test-user"),
		MailNickname:      utils.StringPtr(fmt.Sprintf("test-user-%s", u.randomString)),
		UserPrincipalName: utils.StringPtr(fmt.Sprintf("test-user-%s@%s", u.randomString, u.connection.DomainName)),
		PasswordProfile: &msgraph.UserPasswordProfile{
			Password: utils.StringPtr(fmt.Sprintf("IrPa55w0rd%s", c.randomString)),
		},
		SchemaExtensions: &[]msgraph.SchemaExtensionMap{
			{
				ID: *schema.ID,
				Properties: &ExtensionProperties{
					Property1: utils.StringPtr("my string value"),
					Property2: utils.BoolPtr(true),
				},
			},
		},
	})

	user = testSchemaExtensionsUser_Get(t, u, *user.ID, []msgraph.SchemaExtensionMap{
		{
			ID:         *schema.ID,
			Properties: &ExtensionProperties{},
		},
	})
	schemaExtensions := *user.SchemaExtensions
	if val := schemaExtensions[0].Properties.(*ExtensionProperties).Property1; val == nil || *val != "my string value" {
		t.Fatalf("Unexpected value for Property1 returned: %+v", val)
	}
	if val := schemaExtensions[0].Properties.(*ExtensionProperties).Property2; val == nil || !*val {
		t.Fatalf("Unexpected value for Property2 returned: %+v", val)
	}
	testUsersClient_Delete(t, u, *user.ID)

	testSchemaExtensionsClient_Delete(t, c, *schema.ID)
}

func testSchemaExtensionsUser_Get(t *testing.T, u UsersClientTest, id string, schemaExtensions []msgraph.SchemaExtensionMap) (user *msgraph.User) {
	sel := []string{"id", "displayName"}
	for _, s := range schemaExtensions {
		sel = append(sel, s.ID)
	}
	user, status, err := u.client.GetWithSchemaExtensions(u.connection.Context, id, odata.Query{Select: sel}, &schemaExtensions)
	if err != nil {
		t.Fatalf("UsersClient.Get(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("UsersClient.Get(): invalid status: %d", status)
	}
	if user == nil {
		t.Fatal("UsersClient.Get(): user was nil")
	}
	if user.SchemaExtensions == nil {
		t.Fatal("UsersClient.Get(): user.SchemaExtensions was nil")
	}
	if len(*user.SchemaExtensions) != len(schemaExtensions) {
		t.Fatalf("UsersClient.Get(): unexpected length of user.SchemaExtensions, was %d, expected %d", len(*user.SchemaExtensions), len(schemaExtensions))
	}
	return
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
	schema, status, err := c.client.Get(c.connection.Context, id, odata.Query{})
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
