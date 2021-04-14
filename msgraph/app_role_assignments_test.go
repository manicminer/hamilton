package msgraph_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-uuid"

	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/msgraph"
)

type AppRoleAssignmentsClientTest struct {
	connection   *test.Connection
	client       *msgraph.AppRoleAssignmentsClient
	randomString string
}

func TestAppRoleAssignmentsClient(t *testing.T) {
	rs := test.RandomString()
	// setup service principle test client
	servicePrinciplesClient := ServicePrincipalsClientTest{
		connection:   test.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: rs,
	}
	servicePrinciplesClient.client = msgraph.NewServicePrincipalsClient(servicePrinciplesClient.connection.AuthConfig.TenantID)
	servicePrinciplesClient.client.BaseClient.Authorizer = servicePrinciplesClient.connection.Authorizer

	// setup groups test client
	groupsClient := GroupsClientTest{
		connection:   test.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: rs,
	}
	groupsClient.client = msgraph.NewGroupsClient(groupsClient.connection.AuthConfig.TenantID)
	groupsClient.client.BaseClient.Authorizer = groupsClient.connection.Authorizer

	// setup app role assignments test client
	appRoleAssignClient := AppRoleAssignmentsClientTest{
		connection:   test.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: rs,
	}
	appRoleAssignClient.client = msgraph.NewAppRoleAssignmentsClient(appRoleAssignClient.connection.AuthConfig.TenantID)
	appRoleAssignClient.client.BaseClient.Authorizer = appRoleAssignClient.connection.Authorizer

	// setup applications test client
	appClient := ApplicationsClientTest{
		connection:   test.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: rs,
	}
	appClient.client = msgraph.NewApplicationsClient(appClient.connection.AuthConfig.TenantID)
	appClient.client.BaseClient.Authorizer = appClient.connection.Authorizer

	// create a new test group
	newGroup := msgraph.Group{
		DisplayName:     utils.StringPtr("Test Group"),
		MailEnabled:     utils.BoolPtr(false),
		MailNickname:    utils.StringPtr(fmt.Sprintf("test-group-%s", groupsClient.randomString)),
		SecurityEnabled: utils.BoolPtr(true),
	}
	group := testGroupsClient_Create(t, groupsClient, newGroup)

	// pre-generate uuid for a test app role
	testAppRoleId, _ := uuid.GenerateUUID()
	// create a new test application with a test app role
	app := testApplicationsClient_Create(t, appClient, msgraph.Application{
		DisplayName: utils.StringPtr(fmt.Sprintf("test-application-%s", appClient.randomString)),
		AppRoles: &[]msgraph.AppRole{
			{
				ID:                 utils.StringPtr(testAppRoleId),
				DisplayName:        utils.StringPtr(fmt.Sprintf("test-app-role-%s", appClient.randomString)),
				IsEnabled:          utils.BoolPtr(true),
				AllowedMemberTypes: &[]string{"User", "Application"},
				Description:        utils.StringPtr(fmt.Sprintf("test-app-role-description-%s", appClient.randomString)),
				Value:              utils.StringPtr(fmt.Sprintf("test-app-role-value-%s", appClient.randomString)),
			},
		},
	})

	// create a new test service principle
	sp := testServicePrincipalsClient_Create(t, servicePrinciplesClient, msgraph.ServicePrincipal{
		AccountEnabled: utils.BoolPtr(true),
		AppId:          app.AppId,
		// display name needs to match app's display name
		DisplayName: app.DisplayName,
	})

	// assign app role to the test group
	appRoleAssignment := testAppRoleAssignmentsClient_Assign(t, appRoleAssignClient, *group.ID, *sp.ID, testAppRoleId)

	// list app role assignments for a test group
	appRoleAssignments := testAppRoleAssignmentsClient_List(t, appRoleAssignClient, *group.ID)
	if len(*appRoleAssignments) == 0 {
		t.Fatal("expected at least one app role assignment assigned to the test group")
	}

	// removes app role assignment previously set to the test group
	testAppRoleAssignmentsClient_Remove(t, appRoleAssignClient, *group.ID, *appRoleAssignment.Id)

	// remove all test resources to clean up
	testGroupsClient_Delete(t, groupsClient, *group.ID)
	testServicePrincipalsClient_Delete(t, servicePrinciplesClient, *sp.ID)
	testApplicationsClient_Delete(t, appClient, *app.ID)
}

func testAppRoleAssignmentsClient_List(t *testing.T, c AppRoleAssignmentsClientTest, groupId string) (appRoleAssignments *[]msgraph.AppRoleAssignment) {
	appRoleAssignments, _, err := c.client.List(c.connection.Context, groupId)
	if err != nil {
		t.Fatalf("AppRoleAssignmentsClient.List(): %v", err)
	}
	if appRoleAssignments == nil {
		t.Fatal("AppRoleAssignmentsClient.List(): appRoleAssignments was nil")
	}
	return
}

func testAppRoleAssignmentsClient_Remove(t *testing.T, c AppRoleAssignmentsClientTest, groupId, appRoleAssignmentId string) {
	status, err := c.client.Remove(c.connection.Context, groupId, appRoleAssignmentId)
	if err != nil {
		t.Fatalf("AppRoleAssignmentsClient.Delete(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AppRoleAssignmentsClient.Delete(): invalid status: %d", status)
	}
}

func testAppRoleAssignmentsClient_Assign(t *testing.T, c AppRoleAssignmentsClientTest, groupId, resourceId, appRoleId string) (appRoleAssignment *msgraph.AppRoleAssignment) {
	appRoleAssignment, status, err := c.client.Assign(c.connection.Context, groupId, resourceId, appRoleId)
	if err != nil {
		t.Fatalf("AppRoleAssignmentsClient.Create(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AppRoleAssignmentsClient.Create(): invalid status: %d", status)
	}
	if appRoleAssignment == nil {
		t.Fatal("AppRoleAssignmentsClient.Create(): appRoleAssignment was nil")
	}
	if appRoleAssignment.Id == nil {
		t.Fatal("AppRoleAssignmentsClient.Create(): appRoleAssignment.Id was nil")
	}
	return
}
