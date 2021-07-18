package msgraph_test

import (
	"fmt"
	"testing"

	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/environments"
	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/msgraph"
	"github.com/manicminer/hamilton/odata"
)

type ConditionalAccessPolicyTest struct {
	connection   *test.Connection
	policyClient *msgraph.ConditionalAccessPolicyClient
	groupsClient *msgraph.GroupsClient
	usersClient  *msgraph.UsersClient
	randomString string
}

func TestConditionalAccessPolicyClient(t *testing.T) {
	c := ConditionalAccessPolicyTest{
		connection:   test.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: test.RandomString(),
	}

	c.policyClient = msgraph.NewConditionalAccessPolicyClient(c.connection.AuthConfig.TenantID)
	c.policyClient.BaseClient.Authorizer = c.connection.Authorizer

	c.groupsClient = msgraph.NewGroupsClient(c.connection.AuthConfig.TenantID)
	c.groupsClient.BaseClient.Authorizer = c.connection.Authorizer

	c.usersClient = msgraph.NewUsersClient(c.connection.AuthConfig.TenantID)
	c.usersClient.BaseClient.Authorizer = c.connection.Authorizer

	testAppId := string(environments.PublishedApis["Office365ExchangeOnline"])
	testIncGroup := testGroup_Create(t, c, "test-conditionalAccessPolicy-inc")
	testExcGroup := testGroup_Create(t, c, "test-conditionalAccessPolicy-exc")
	testUser := testUser_Create(t, c)

	// act
	policy := testConditionalAccessPolicysClient_Create(t, c, msgraph.ConditionalAccessPolicy{
		DisplayName: utils.StringPtr(fmt.Sprintf("test-policy-%s", c.randomString)),
		State:       utils.StringPtr("enabled"),
		Conditions: &msgraph.ConditionalAccessConditionSet{
			ClientAppTypes: &[]string{"mobileAppsAndDesktopClients", "browser"},
			Applications: &msgraph.ConditionalAccessApplications{
				IncludeApplications: &[]string{testAppId},
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

	updatePolicy := msgraph.ConditionalAccessPolicy{
		ID:          policy.ID,
		DisplayName: utils.StringPtr(fmt.Sprintf("test-policy-updated-%s", c.randomString)),
	}
	testConditionalAccessPolicysClient_Update(t, c, updatePolicy)

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
	policy, status, err := c.policyClient.Get(c.connection.Context, id, odata.Query{})
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
	policies, _, err := c.policyClient.List(c.connection.Context, odata.Query{Top: 10})
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
	group, _, err := c.groupsClient.Create(c.connection.Context, msgraph.Group{
		DisplayName:     utils.StringPtr(fmt.Sprintf("%s-%s", prefix, c.randomString)),
		MailEnabled:     utils.BoolPtr(false),
		MailNickname:    utils.StringPtr(fmt.Sprintf("%s-%s", prefix, c.randomString)),
		SecurityEnabled: utils.BoolPtr(true),
	})

	if err != nil {
		t.Fatalf("GroupsClient.Create() - Could not create test group: %v", err)
	}
	return
}

func testGroup_Delete(t *testing.T, c ConditionalAccessPolicyTest, group *msgraph.Group) {
	_, err := c.groupsClient.Delete(c.connection.Context, *group.ID)
	if err != nil {
		t.Fatalf("GroupsClient.Delete() - Could not delete test group: %v", err)
	}
}

func testUser_Create(t *testing.T, c ConditionalAccessPolicyTest) (user *msgraph.User) {
	user, _, err := c.usersClient.Create(c.connection.Context, msgraph.User{
		AccountEnabled:    utils.BoolPtr(true),
		DisplayName:       utils.StringPtr("test-user-conditionalAccessPolicy"),
		MailNickname:      utils.StringPtr(fmt.Sprintf("test-user-%s", c.randomString)),
		UserPrincipalName: utils.StringPtr(fmt.Sprintf("test-user-%s@%s", c.randomString, c.connection.DomainName)),
		PasswordProfile: &msgraph.UserPasswordProfile{
			Password: utils.StringPtr(fmt.Sprintf("IrPa55w0rd%s", c.randomString)),
		},
	})

	if err != nil {
		t.Fatalf("UsersClient.Create() - Could not create test user: %v", err)
	}
	return
}

func testUser_Delete(t *testing.T, c ConditionalAccessPolicyTest, user *msgraph.User) {
	_, err := c.usersClient.Delete(c.connection.Context, *user.ID)
	if err != nil {
		t.Fatalf("UsersClient.Delete() - Could not delete test user: %v", err)
	}
}
