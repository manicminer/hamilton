package clients_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/clients"
	"github.com/manicminer/hamilton/clients/internal"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/models"
)

type ConditionalAccessPolicyTest struct {
	connection      *internal.Connection
	policyClient    *clients.ConditionalAccessPolicyClient
	groupClient     *clients.GroupsClient
	userClient      *clients.UsersClient
	randomString    string
	enterpriseAppId string
}

func TestConditionalAccessPolicyClient(t *testing.T) {
	c := ConditionalAccessPolicyTest{
		connection:      internal.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString:    internal.RandomString(),
		enterpriseAppId: os.Getenv("ENTERPRISE_APP_ID"), // Enterprise Apps cannot be created using Graph. Therefore, to simplify testing pass the Enterprise App ID via an Env variable.
	}

	if c.enterpriseAppId == "" {
		t.Fatalf("ConditionalAccessPolicyClient.Create(): ENTERPRISE_APP_ID is not set")
	}

	c.policyClient = clients.NewConditionalAccessPolicyClient(c.connection.AuthConfig.TenantID)
	c.policyClient.BaseClient.Authorizer = c.connection.Authorizer

	c.groupClient = clients.NewGroupsClient(c.connection.AuthConfig.TenantID)
	c.groupClient.BaseClient.Authorizer = c.connection.Authorizer

	c.userClient = clients.NewUsersClient(c.connection.AuthConfig.TenantID)
	c.userClient.BaseClient.Authorizer = c.connection.Authorizer

	testIncGroup := testGroup_Create(t, c, "inc")
	testExcGroup := testGroup_Create(t, c, "exc")
	testUser := testUser_Create(t, c)

	// act
	policy := testConditionalAccessPolicysClient_Create(t, c, models.ConditionalAccessPolicy{
		DisplayName: utils.StringPtr(fmt.Sprintf("test-policy-%s", c.randomString)),
		State:       utils.StringPtr("enabled"),
		Conditions: &models.ConditionalAccessConditionSet{
			ClientAppTypes: &[]string{"mobileAppsAndDesktopClients", "browser"},
			Applications: &models.ConditionalAccessApplications{
				IncludeApplications: &[]string{*&c.enterpriseAppId},
			},
			Users: &models.ConditionalAccessUsers{
				IncludeUsers:  &[]string{"All"},
				ExcludeUsers:  &[]string{*testUser.ID, "GuestsOrExternalUsers"},
				IncludeGroups: &[]string{*testIncGroup.ID},
				ExcludeGroups: &[]string{*testExcGroup.ID},
			},
			Locations: &models.ConditionalAccessLocations{
				IncludeLocations: &[]string{"All"},
				ExcludeLocations: &[]string{"AllTrusted"},
			},
		},
		GrantControls: &models.ConditionalAccessGrantControls{
			Operator:        utils.StringPtr("OR"),
			BuiltInControls: &[]string{"block"},
		},
	})

	policy.DisplayName = utils.StringPtr(fmt.Sprintf("test-policy-updated-%s", c.randomString))
	testConditionalAccessPolicysClient_Update(t, c, *policy)

	testConditionalAccessPolicysClient_List(t, c)
	testConditionalAccessPolicysClient_Get(t, c, *policy.ID)
	testConditionalAccessPolicysClient_Delete(t, c, *policy.ID)

	// cleanup
	testGroup_Delete(t, c, testIncGroup)
	testGroup_Delete(t, c, testExcGroup)
	testUser_Delete(t, c, testUser)
}

func testConditionalAccessPolicysClient_Create(t *testing.T, c ConditionalAccessPolicyTest, a models.ConditionalAccessPolicy) (conditionalAccessPolicy *models.ConditionalAccessPolicy) {
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

func testConditionalAccessPolicysClient_Get(t *testing.T, c ConditionalAccessPolicyTest, id string) (policy *models.ConditionalAccessPolicy) {
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

func testConditionalAccessPolicysClient_Update(t *testing.T, c ConditionalAccessPolicyTest, policy models.ConditionalAccessPolicy) {
	status, err := c.policyClient.Update(c.connection.Context, policy)
	if err != nil {
		t.Fatalf("ConditionalAccessPolicyClient.Update(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ConditionalAccessPolicyClient.Update(): invalid status: %d", status)
	}
}

func testConditionalAccessPolicysClient_List(t *testing.T, c ConditionalAccessPolicyTest) (policies *[]models.ConditionalAccessPolicy) {
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

func testGroup_Create(t *testing.T, c ConditionalAccessPolicyTest, prefix string) (group *models.Group) {
	group, _, err := c.groupClient.Create(c.connection.Context, models.Group{
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

func testGroup_Delete(t *testing.T, c ConditionalAccessPolicyTest, group *models.Group) {
	_, err := c.groupClient.Delete(c.connection.Context, *group.ID)
	if err != nil {
		t.Fatalf("ConditionalAccessPolicyClient.Create() - Could not delete test group: %v", err)
	}
	return
}

func testUser_Create(t *testing.T, c ConditionalAccessPolicyTest) (user *models.User) {
	user, _, err := c.userClient.Create(c.connection.Context, models.User{
		AccountEnabled:    utils.BoolPtr(true),
		DisplayName:       utils.StringPtr("Test User"),
		MailNickname:      utils.StringPtr(fmt.Sprintf("test-user-%s", c.randomString)),
		UserPrincipalName: utils.StringPtr(fmt.Sprintf("test-user-%s@%s", c.randomString, c.connection.DomainName)),
		PasswordProfile: &models.UserPasswordProfile{
			Password: utils.StringPtr(fmt.Sprintf("IrPa55w0rd%s", c.randomString)),
		},
	})

	if err != nil {
		t.Fatalf("ConditionalAccessPolicyClient.Create() - Could not create test user: %v", err)
	}
	return
}

func testUser_Delete(t *testing.T, c ConditionalAccessPolicyTest, user *models.User) {
	_, err := c.userClient.Delete(c.connection.Context, *user.ID)
	if err != nil {
		t.Fatalf("ConditionalAccessPolicyClient.Create() - Could not delete test user: %v", err)
	}
	return
}
