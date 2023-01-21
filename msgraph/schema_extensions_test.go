package msgraph_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/msgraph"
	"github.com/manicminer/hamilton/odata"
)

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
	c := test.NewTest(t)
	defer c.CancelFunc()

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

	testSchemaExtensionsGroup(t, c, schema)
	testSchemaExtensionsUser(t, c, schema)

	testSchemaExtensionsClient_Delete(t, c, *schema.ID)
}

func testSchemaExtensionsGroup(t *testing.T, c *test.Test, schema *msgraph.SchemaExtension) {
	// First create a group having schema extension data expressed using msgraph.SchemaExtensionMap
	group := testGroupsClient_Create(t, c, msgraph.Group{
		DisplayName:     utils.StringPtr("test-group"),
		GroupTypes:      &[]msgraph.GroupType{msgraph.GroupTypeUnified},
		MailEnabled:     utils.BoolPtr(true),
		MailNickname:    utils.StringPtr(fmt.Sprintf("test-365-group-%s", c.RandomString)),
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
	group = testSchemaExtensionsGroup_Get(t, c, *group.ID(), []msgraph.SchemaExtensionData{
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
	testGroupsClient_Update(t, c, *group)

	// Finally retrieve the group and populate using MyExtensionProperties
	group = testSchemaExtensionsGroup_Get(t, c, *group.ID(), []msgraph.SchemaExtensionData{
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

	testGroupsClient_Delete(t, c, *group.ID())
}

func testSchemaExtensionsGroup_Get(t *testing.T, c *test.Test, id string, schemaExtensions []msgraph.SchemaExtensionData) (group *msgraph.Group) {
	sel := []string{"id", "displayName"}
	for _, s := range schemaExtensions {
		sel = append(sel, s.ID)
	}
	group, status, err := c.GroupsClient.GetWithSchemaExtensions(c.Context, id, odata.Query{Select: sel}, &schemaExtensions)
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

func testSchemaExtensionsUser(t *testing.T, c *test.Test, schema *msgraph.SchemaExtension) {
	// First create a user having schema extension data expressed using msgraph.SchemaExtensionMap
	user := testUsersClient_Create(t, c, msgraph.User{
		AccountEnabled:    utils.BoolPtr(true),
		DisplayName:       utils.StringPtr("test-user"),
		MailNickname:      utils.StringPtr(fmt.Sprintf("test-user-%s", c.RandomString)),
		UserPrincipalName: utils.StringPtr(fmt.Sprintf("test-user-%s@%s", c.RandomString, c.Connections["default"].DomainName)),
		PasswordProfile: &msgraph.UserPasswordProfile{
			Password: utils.StringPtr(fmt.Sprintf("IrPa55w0rd%s", c.RandomString)),
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
	user = testSchemaExtensionsUser_Get(t, c, *user.ID(), []msgraph.SchemaExtensionData{
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
	testUsersClient_Update(t, c, *user)

	// Finally retrieve the user and populate using MyExtensionProperties
	user = testSchemaExtensionsUser_Get(t, c, *user.ID(), []msgraph.SchemaExtensionData{
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

	testUsersClient_Delete(t, c, *user.ID())
}

func testSchemaExtensionsUser_Get(t *testing.T, c *test.Test, id string, schemaExtensions []msgraph.SchemaExtensionData) (user *msgraph.User) {
	sel := []string{"id", "displayName"}
	for _, s := range schemaExtensions {
		sel = append(sel, s.ID)
	}
	user, status, err := c.UsersClient.GetWithSchemaExtensions(c.Context, id, odata.Query{Select: sel}, &schemaExtensions)
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

func testSchemaExtensionsClient_Create(t *testing.T, c *test.Test, s msgraph.SchemaExtension) (schema *msgraph.SchemaExtension) {
	schema, status, err := c.SchemaExtensionsClient.Create(c.Context, s)
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

func testSchemaExtensionsClient_Get(t *testing.T, c *test.Test, id string) (schema *msgraph.SchemaExtension) {
	schema, status, err := c.SchemaExtensionsClient.Get(c.Context, id, odata.Query{})
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

func testSchemaExtensionsClient_Update(t *testing.T, c *test.Test, s msgraph.SchemaExtension) {
	status, err := c.SchemaExtensionsClient.Update(c.Context, s)
	if err != nil {
		t.Fatalf("SchemaExtensionsClient.Update(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("SchemaExtensionsClient.Update(): invalid status: %d", status)
	}
}

func testSchemaExtensionsClient_Delete(t *testing.T, c *test.Test, id string) {
	status, err := c.SchemaExtensionsClient.Delete(c.Context, id)
	if err != nil {
		t.Fatalf("SchemaExtensionsClient.Delete(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("SchemaExtensionsClient.Delete(): invalid status: %d", status)
	}
}
