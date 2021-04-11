package msgraph_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/msgraph"
)

var (
	enterpriseAppId = os.Getenv("ENTERPRISE_APP_ID") // Enterprise Apps cannot be created using Graph. Therefore, to simplify testing pass the Enterprise App ID via an Env variable.
)

type ConditionalAccessPolicyTest struct {
	connection   *test.Connection
	policyClient *msgraph.ConditionalAccessPolicyClient
	groupClient  *msgraph.GroupsClient
	userClient   *msgraph.UsersClient
	randomString string
}

func TestConditionalAccessPolicyClient(t *testing.T) {
	c := ConditionalAccessPolicyTest{
		connection:   test.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: test.RandomString(),
	}

	if enterpriseAppId == "" {
		t.Fatalf("ConditionalAccessPolicyClient.Create(): ENTERPRISE_APP_ID is not set")
	}

	c.policyClient = msgraph.NewConditionalAccessPolicyClient(c.connection.AuthConfig.TenantID)
	c.policyClient.BaseClient.Authorizer = c.connection.Authorizer

	c.groupClient = msgraph.NewGroupsClient(c.connection.AuthConfig.TenantID)
	c.groupClient.BaseClient.Authorizer = c.connection.Authorizer

	c.userClient = msgraph.NewUsersClient(c.connection.AuthConfig.TenantID)
	c.userClient.BaseClient.Authorizer = c.connection.Authorizer

	testIncGroup := testGroup_Create(t, c, "inc")
	testExcGroup := testGroup_Create(t, c, "exc")
	testUser := testUser_Create(t, c)

	// act
	policy := testConditionalAccessPolicysClient_Create(t, c, msgraph.ConditionalAccessPolicy{
		DisplayName: utils.StringPtr(fmt.Sprintf("test-policy-%s", c.randomString)),
		State:       utils.StringPtr("enabled"),
		Conditions: &msgraph.ConditionalAccessConditionSet{
			ClientAppTypes: &[]string{"mobileAppsAndDesktopClients", "browser"},
			Applications: &msgraph.ConditionalAccessApplications{
				IncludeApplications: &[]string{enterpriseAppId},
			},
			Users: &msgraph.ConditionalAccessUsers{
				IncludeUsers:  &[]string{"All"},
				ExcludeUsers:  &[]string{*testUser.ID, "GuestsOrExternalUsers"},
				IncludeGroups: &[]string{*testIncGroup.ID},
				ExcludeGroups: &[]string{*testExcGroup.ID},
			},
			Locations: &msgraph.ConditionalAccessLocations{
				IncludeLocations: &[]string{"All"},
				ExcludeLocations: &[]string{"AllTrusted"},
			},
		},
		GrantControls: &msgraph.ConditionalAccessGrantControls{
			Operator:        utils.StringPtr("OR"),
			BuiltInControls: &[]string{"block"},
		},
	})

	initCreatedDate := policy.CreatedDateTime
	policy.DisplayName = utils.StringPtr(fmt.Sprintf("test-policy-updated-%s", c.randomString))
	testConditionalAccessPolicysClient_Update(t, c, *policy)

	if policy.CreatedDateTime != initCreatedDate {
		t.Fatalf("ConditionalAccessPolicyClient.Create(): unintended mutation on CreatedDateTime")
	}

	testConditionalAccessPolicysClient_List(t, c)
	testConditionalAccessPolicysClient_Get(t, c, *policy.ID)
	testConditionalAccessPolicysClient_Delete(t, c, *policy.ID)

	// cleanup
	testGroup_Delete(t, c, testIncGroup)
	testGroup_Delete(t, c, testExcGroup)
	testUser_Delete(t, c, testUser)
}

func testConditionalAccessPolicysClient_Create(t *testing.T, c ConditionalAccessPolicyTest, a msgraph.ConditionalAccessPolicy) (conditionalAccessPolicy *msgraph.ConditionalAccessPolicy) {
	conditionalAccessPolicy, status, err := c.policyClient.Create(c.connection.Context, a)
	if err != nil {
		t.Fatalf("ConditionalAccessPolicyClient.Create(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ConditionalAccessPolicyClient.Create(): invalid status: %d", status)
	}
	if conditionalAccessPolicy == nil {
		t.Fatal("ConditionalAccessPolicyClient.Create(): conditionalAccessPolicy was nil")
	}
	if conditionalAccessPolicy.ID == nil {
		t.Fatal("ConditionalAccessPolicyClient.Create(): conditionalAccessPolicy.ID was nil")
	}
	return
}

func testConditionalAccessPolicysClient_Get(t *testing.T, c ConditionalAccessPolicyTest, id string) (policy *msgraph.ConditionalAccessPolicy) {
	policy, status, err := c.policyClient.Get(c.connection.Context, id)
	if err != nil {
		t.Fatalf("ConditionalAccessPolicyClient.Get(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ConditionalAccessPolicyClient.Get(): invalid status: %d", status)
	}
	if policy == nil {
		t.Fatal("ConditionalAccessPolicyClient.Get(): policy was nil")
	}
	return
}

func testConditionalAccessPolicysClient_Update(t *testing.T, c ConditionalAccessPolicyTest, policy msgraph.ConditionalAccessPolicy) {
	status, err := c.policyClient.Update(c.connection.Context, policy)
	if err != nil {
		t.Fatalf("ConditionalAccessPolicyClient.Update(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ConditionalAccessPolicyClient.Update(): invalid status: %d", status)
	}
}

func testConditionalAccessPolicysClient_List(t *testing.T, c ConditionalAccessPolicyTest) (policies *[]msgraph.ConditionalAccessPolicy) {
	policies, _, err := c.policyClient.List(c.connection.Context, "")
	if err != nil {
		t.Fatalf("ConditionalAccessPolicyClient.List(): %v", err)
	}
	if policies == nil {
		t.Fatal("ConditionalAccessPolicyClient.List(): policies was nil")
	}
	return
}

func testConditionalAccessPolicysClient_Delete(t *testing.T, c ConditionalAccessPolicyTest, id string) {
	status, err := c.policyClient.Delete(c.connection.Context, id)
	if err != nil {
		t.Fatalf("ConditionalAccessPolicyClient.Delete(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ConditionalAccessPolicyClient.Delete(): invalid status: %d", status)
	}
}

func testGroup_Create(t *testing.T, c ConditionalAccessPolicyTest, prefix string) (group *msgraph.Group) {
	group, _, err := c.groupClient.Create(c.connection.Context, msgraph.Group{
		DisplayName:     utils.StringPtr(fmt.Sprintf("%s-test-group-%s", prefix, c.randomString)),
		MailEnabled:     utils.BoolPtr(false),
		MailNickname:    utils.StringPtr(fmt.Sprintf("%s-test-group-%s", prefix, c.randomString)),
		SecurityEnabled: utils.BoolPtr(true),
	})

	if err != nil {
		t.Fatalf("ConditionalAccessPolicyClient.Create() - Could not create test group: %v", err)
	}
	return
}

func testGroup_Delete(t *testing.T, c ConditionalAccessPolicyTest, group *msgraph.Group) {
	_, err := c.groupClient.Delete(c.connection.Context, *group.ID)
	if err != nil {
		t.Fatalf("ConditionalAccessPolicyClient.Create() - Could not delete test group: %v", err)
	}
}

func testUser_Create(t *testing.T, c ConditionalAccessPolicyTest) (user *msgraph.User) {
	user, _, err := c.userClient.Create(c.connection.Context, msgraph.User{
		AccountEnabled:    utils.BoolPtr(true),
		DisplayName:       utils.StringPtr("Test User"),
		MailNickname:      utils.StringPtr(fmt.Sprintf("test-user-%s", c.randomString)),
		UserPrincipalName: utils.StringPtr(fmt.Sprintf("test-user-%s@%s", c.randomString, c.connection.DomainName)),
		PasswordProfile: &msgraph.UserPasswordProfile{
			Password: utils.StringPtr(fmt.Sprintf("IrPa55w0rd%s", c.randomString)),
		},
	})

	if err != nil {
		t.Fatalf("ConditionalAccessPolicyClient.Create() - Could not create test user: %v", err)
	}
	return
}

func testUser_Delete(t *testing.T, c ConditionalAccessPolicyTest, user *msgraph.User) {
	_, err := c.userClient.Delete(c.connection.Context, *user.ID)
	if err != nil {
		t.Fatalf("ConditionalAccessPolicyClient.Create() - Could not delete test user: %v", err)
	}
}
