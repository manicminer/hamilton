package msgraph_test

import (
	"fmt"
	"testing"

	"github.com/manicminer/hamilton/environments"
	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/msgraph"
	"github.com/manicminer/hamilton/odata"
)

func TestConditionalAccessPolicyClient(t *testing.T) {
	c := test.NewTest(t)
	defer c.CancelFunc()

	testAppId := environments.PublishedApis["Office365ExchangeOnline"]
	testIncGroup := testGroup_Create(t, c, "test-conditionalAccessPolicy-inc")
	testExcGroup := testGroup_Create(t, c, "test-conditionalAccessPolicy-exc")
	testUser := testUser_Create(t, c)

	policy := testConditionalAccessPolicysClient_Create(t, c, msgraph.ConditionalAccessPolicy{
		DisplayName: utils.StringPtr(fmt.Sprintf("test-policy-%s", c.RandomString)),
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
		DisplayName: utils.StringPtr(fmt.Sprintf("test-policy-updated-%s", c.RandomString)),
	}
	testConditionalAccessPolicysClient_Update(t, c, updatePolicy)

	testConditionalAccessPolicysClient_List(t, c)
	testConditionalAccessPolicysClient_Get(t, c, *policy.ID)
	testConditionalAccessPolicysClient_Delete(t, c, *policy.ID)

	testGroup_Delete(t, c, testIncGroup)
	testGroup_Delete(t, c, testExcGroup)
	testUser_Delete(t, c, testUser)
}

func testConditionalAccessPolicysClient_Create(t *testing.T, c *test.Test, a msgraph.ConditionalAccessPolicy) (conditionalAccessPolicy *msgraph.ConditionalAccessPolicy) {
	conditionalAccessPolicy, status, err := c.ConditionalAccessPoliciesClient.Create(c.Context, a)
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

func testConditionalAccessPolicysClient_Get(t *testing.T, c *test.Test, id string) (policy *msgraph.ConditionalAccessPolicy) {
	policy, status, err := c.ConditionalAccessPoliciesClient.Get(c.Context, id, odata.Query{})
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

func testConditionalAccessPolicysClient_Update(t *testing.T, c *test.Test, policy msgraph.ConditionalAccessPolicy) {
	status, err := c.ConditionalAccessPoliciesClient.Update(c.Context, policy)
	if err != nil {
		t.Fatalf("ConditionalAccessPolicyClient.Update(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ConditionalAccessPolicyClient.Update(): invalid status: %d", status)
	}
}

func testConditionalAccessPolicysClient_List(t *testing.T, c *test.Test) (policies *[]msgraph.ConditionalAccessPolicy) {
	policies, _, err := c.ConditionalAccessPoliciesClient.List(c.Context, odata.Query{Top: 10})
	if err != nil {
		t.Fatalf("ConditionalAccessPolicyClient.List(): %v", err)
	}
	if policies == nil {
		t.Fatal("ConditionalAccessPolicyClient.List(): policies was nil")
	}
	return
}

func testConditionalAccessPolicysClient_Delete(t *testing.T, c *test.Test, id string) {
	status, err := c.ConditionalAccessPoliciesClient.Delete(c.Context, id)
	if err != nil {
		t.Fatalf("ConditionalAccessPolicyClient.Delete(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ConditionalAccessPolicyClient.Delete(): invalid status: %d", status)
	}
}

func testGroup_Create(t *testing.T, c *test.Test, prefix string) (group *msgraph.Group) {
	group, _, err := c.GroupsClient.Create(c.Context, msgraph.Group{
		DisplayName:     utils.StringPtr(fmt.Sprintf("%s-%s", prefix, c.RandomString)),
		MailEnabled:     utils.BoolPtr(false),
		MailNickname:    utils.StringPtr(fmt.Sprintf("%s-%s", prefix, c.RandomString)),
		SecurityEnabled: utils.BoolPtr(true),
	})

	if err != nil {
		t.Fatalf("GroupsClient.Create() - Could not create test group: %v", err)
	}
	return
}

func testGroup_Delete(t *testing.T, c *test.Test, group *msgraph.Group) {
	_, err := c.GroupsClient.Delete(c.Context, *group.ID)
	if err != nil {
		t.Fatalf("GroupsClient.Delete() - Could not delete test group: %v", err)
	}
}

func testUser_Create(t *testing.T, c *test.Test) (user *msgraph.User) {
	user, _, err := c.UsersClient.Create(c.Context, msgraph.User{
		AccountEnabled:    utils.BoolPtr(true),
		DisplayName:       utils.StringPtr("test-user-conditionalAccessPolicy"),
		MailNickname:      utils.StringPtr(fmt.Sprintf("test-user-%s", c.RandomString)),
		UserPrincipalName: utils.StringPtr(fmt.Sprintf("test-user-%s@%s", c.RandomString, c.Connections["default"].DomainName)),
		PasswordProfile: &msgraph.UserPasswordProfile{
			Password: utils.StringPtr(fmt.Sprintf("IrPa55w0rd%s", c.RandomString)),
		},
	})

	if err != nil {
		t.Fatalf("UsersClient.Create() - Could not create test user: %v", err)
	}
	return
}

func testUser_Delete(t *testing.T, c *test.Test, user *msgraph.User) {
	_, err := c.UsersClient.Delete(c.Context, *user.ID)
	if err != nil {
		t.Fatalf("UsersClient.Delete() - Could not delete test user: %v", err)
	}
}
