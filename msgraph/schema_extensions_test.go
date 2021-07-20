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

type MyExtensionProperties struct {
	Property1 *string `json:"property1,omitempty"`
	Property2 *bool   `json:"property2,omitempty"`
}

func (e *MyExtensionProperties) UnmarshalJSON(data []byte) error {
	type ep MyExtensionProperties
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

	g := GroupsClientTest{
		connection:   test.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: rs,
	}
	g.client = msgraph.NewGroupsClient(g.connection.AuthConfig.TenantID)
	g.client.BaseClient.Authorizer = g.connection.Authorizer

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

	targetTypes := []msgraph.ExtensionSchemaTargetType{
		msgraph.ExtensionSchemaTargetTypeGroup,
		msgraph.ExtensionSchemaTargetTypeUser,
	}
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

	// Replication seems to be a problem with schema extensions, no viable workaround as yet
	time.Sleep(10 * time.Second)

	testSchemaExtensionsGroup(t, g, schema)
	testSchemaExtensionsUser(t, u, schema)

	testSchemaExtensionsClient_Delete(t, c, *schema.ID)
}

func testSchemaExtensionsGroup(t *testing.T, g GroupsClientTest, schema *msgraph.SchemaExtension) {
	// First create a group having schema extension data expressed using msgraph.SchemaExtensionMap
	group := testGroupsClient_Create(t, g, msgraph.Group{
		DisplayName:     utils.StringPtr("test-group"),
		GroupTypes:      []msgraph.GroupType{msgraph.GroupTypeUnified},
		MailEnabled:     utils.BoolPtr(true),
		MailNickname:    utils.StringPtr(fmt.Sprintf("test-365-group-%s", g.randomString)),
		SecurityEnabled: utils.BoolPtr(true),
		SchemaExtensions: &[]msgraph.SchemaExtensionData{
			{
				ID: *schema.ID,
				Properties: &msgraph.SchemaExtensionMap{
					"property1": utils.StringPtr("my string value"),
					"property2": utils.BoolPtr(true),
				},
			},
		},
	})

	// Then retrieve the group and populate using msgraph.SchemaExtensionGraph
	group = testSchemaExtensionsGroup_Get(t, g, *group.ID, []msgraph.SchemaExtensionData{
		{
			ID:         *schema.ID,
			Properties: &msgraph.SchemaExtensionMap{},
		},
	})
	schemaExtensions := *group.SchemaExtensions
	m := *schemaExtensions[0].Properties.(*msgraph.SchemaExtensionMap)
	if val, ok := m["property1"].(string); !ok || val != "my string value" {
		t.Fatalf("Unexpected value for property1 returned: %+v", val)
	}
	if val, ok := m["property2"].(bool); !ok || !val {
		t.Fatalf("Unexpected value for property2 returned: %+v", val)
	}

	// Next, update the group with schema extension data expressed using MyExtensionProperties
	group.SchemaExtensions = &[]msgraph.SchemaExtensionData{
		{
			ID: *schema.ID,
			Properties: &MyExtensionProperties{
				Property1: utils.StringPtr("my stringy value"),
				Property2: utils.BoolPtr(true),
			},
		},
	}
	testGroupsClient_Update(t, g, *group)

	// Finally retrieve the group and populate using MyExtensionProperties
	group = testSchemaExtensionsGroup_Get(t, g, *group.ID, []msgraph.SchemaExtensionData{
		{
			ID:         *schema.ID,
			Properties: &MyExtensionProperties{},
		},
	})
	schemaExtensions = *group.SchemaExtensions
	if val := schemaExtensions[0].Properties.(*MyExtensionProperties).Property1; val == nil || *val != "my stringy value" {
		t.Fatalf("Unexpected value for Property1 returned: %+v", val)
	}
	if val := schemaExtensions[0].Properties.(*MyExtensionProperties).Property2; val == nil || !*val {
		t.Fatalf("Unexpected value for Property2 returned: %+v", val)
	}

	testGroupsClient_Delete(t, g, *group.ID)
}

func testSchemaExtensionsGroup_Get(t *testing.T, g GroupsClientTest, id string, schemaExtensions []msgraph.SchemaExtensionData) (group *msgraph.Group) {
	sel := []string{"id", "displayName"}
	for _, s := range schemaExtensions {
		sel = append(sel, s.ID)
	}
	group, status, err := g.client.GetWithSchemaExtensions(g.connection.Context, id, odata.Query{Select: sel}, &schemaExtensions)
	if err != nil {
		t.Fatalf("GroupsClient.Get(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("GroupsClient.Get(): invalid status: %d", status)
	}
	if group == nil {
		t.Fatal("GroupsClient.Get(): group was nil")
	}
	if group.SchemaExtensions == nil {
		t.Fatal("GroupsClient.Get(): group.SchemaExtensions was nil")
	}
	if len(*group.SchemaExtensions) != len(schemaExtensions) {
		t.Fatalf("GroupsClient.Get(): unexpected length of group.SchemaExtensions, was %d, expected %d", len(*group.SchemaExtensions), len(schemaExtensions))
	}
	return
}

func testSchemaExtensionsUser(t *testing.T, u UsersClientTest, schema *msgraph.SchemaExtension) {
	// First create a user having schema extension data expressed using msgraph.SchemaExtensionMap
	user := testUsersClient_Create(t, u, msgraph.User{
		AccountEnabled:    utils.BoolPtr(true),
		DisplayName:       utils.StringPtr("test-user"),
		MailNickname:      utils.StringPtr(fmt.Sprintf("test-user-%s", u.randomString)),
		UserPrincipalName: utils.StringPtr(fmt.Sprintf("test-user-%s@%s", u.randomString, u.connection.DomainName)),
		PasswordProfile: &msgraph.UserPasswordProfile{
			Password: utils.StringPtr(fmt.Sprintf("IrPa55w0rd%s", u.randomString)),
		},
		SchemaExtensions: &[]msgraph.SchemaExtensionData{
			{
				ID: *schema.ID,
				Properties: &msgraph.SchemaExtensionMap{
					"property1": utils.StringPtr("my string value"),
					"property2": utils.BoolPtr(true),
				},
			},
		},
	})

	// Then retrieve the user and populate using msgraph.SchemaExtensionGraph
	user = testSchemaExtensionsUser_Get(t, u, *user.ID, []msgraph.SchemaExtensionData{
		{
			ID:         *schema.ID,
			Properties: &msgraph.SchemaExtensionMap{},
		},
	})
	schemaExtensions := *user.SchemaExtensions
	m := *schemaExtensions[0].Properties.(*msgraph.SchemaExtensionMap)
	if val, ok := m["property1"].(string); !ok || val != "my string value" {
		t.Fatalf("Unexpected value for property1 returned: %+v", val)
	}
	if val, ok := m["property2"].(bool); !ok || !val {
		t.Fatalf("Unexpected value for property2 returned: %+v", val)
	}

	// Next, update the user with schema extension data expressed using MyExtensionProperties
	user.SchemaExtensions = &[]msgraph.SchemaExtensionData{
		{
			ID: *schema.ID,
			Properties: &MyExtensionProperties{
				Property1: utils.StringPtr("my stringy value"),
				Property2: utils.BoolPtr(true),
			},
		},
	}
	testUsersClient_Update(t, u, *user)

	// Finally retrieve the user and populate using MyExtensionProperties
	user = testSchemaExtensionsUser_Get(t, u, *user.ID, []msgraph.SchemaExtensionData{
		{
			ID:         *schema.ID,
			Properties: &MyExtensionProperties{},
		},
	})
	schemaExtensions = *user.SchemaExtensions
	if val := schemaExtensions[0].Properties.(*MyExtensionProperties).Property1; val == nil || *val != "my stringy value" {
		t.Fatalf("Unexpected value for Property1 returned: %+v", val)
	}
	if val := schemaExtensions[0].Properties.(*MyExtensionProperties).Property2; val == nil || !*val {
		t.Fatalf("Unexpected value for Property2 returned: %+v", val)
	}

	testUsersClient_Delete(t, u, *user.ID)
}

func testSchemaExtensionsUser_Get(t *testing.T, u UsersClientTest, id string, schemaExtensions []msgraph.SchemaExtensionData) (user *msgraph.User) {
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
