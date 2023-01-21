package msgraph_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-uuid"

	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/msgraph"
	"github.com/manicminer/hamilton/odata"
)

func TestAppRoleAssignedToClient(t *testing.T) {
	c := test.NewTest(t)
	defer c.CancelFunc()

	// Scaffold the resource application
	testResourceAppRoleId, _ := uuid.GenerateUUID()
	resourceApp := testApplicationsClient_Create(t, c, msgraph.Application{
		DisplayName: utils.StringPtr(fmt.Sprintf("test-application-appRoleAssignedTo-%s", c.RandomString)),
		AppRoles: &[]msgraph.AppRole{
			{
				ID:          utils.StringPtr(testResourceAppRoleId),
				DisplayName: utils.StringPtr(fmt.Sprintf("test-resourceApp-role-%s", c.RandomString)),
				IsEnabled:   utils.BoolPtr(true),
				Description: utils.StringPtr(fmt.Sprintf("test-resourceApp-role-description-%s", c.RandomString)),
				Value:       utils.StringPtr(fmt.Sprintf("test-resourceApp-role-value-%s", c.RandomString)),
				AllowedMemberTypes: &[]msgraph.AppRoleAllowedMemberType{
					msgraph.AppRoleAllowedMemberTypeApplication,
					msgraph.AppRoleAllowedMemberTypeUser,
				},
			},
		},
	})
	resourceServicePrincipal := testServicePrincipalsClient_Create(t, c, msgraph.ServicePrincipal{
		AccountEnabled:            utils.BoolPtr(true),
		AppId:                     resourceApp.AppId,
		AppRoleAssignmentRequired: utils.BoolPtr(true),
		DisplayName:               resourceApp.DisplayName,
	})

	// Scaffold the group, user and service principal to assign roles to
	group := testGroupsClient_Create(t, c, msgraph.Group{
		DisplayName:     utils.StringPtr("test-group-appRoleAssignments"),
		MailEnabled:     utils.BoolPtr(false),
		MailNickname:    utils.StringPtr(fmt.Sprintf("test-group-%s", c.RandomString)),
		SecurityEnabled: utils.BoolPtr(true),
	})
	user := testUsersClient_Create(t, c, msgraph.User{
		AccountEnabled:    utils.BoolPtr(true),
		DisplayName:       utils.StringPtr("test-user-appRoleAssignments"),
		MailNickname:      utils.StringPtr(fmt.Sprintf("test-user-%s", c.RandomString)),
		UserPrincipalName: utils.StringPtr(fmt.Sprintf("test-user-%s@%s", c.RandomString, c.Connections["default"].DomainName)),
		PasswordProfile: &msgraph.UserPasswordProfile{
			Password: utils.StringPtr(fmt.Sprintf("IrPa55w0rd%s", c.RandomString)),
		},
	})
	app := testApplicationsClient_Create(t, c, msgraph.Application{
		DisplayName: utils.StringPtr(fmt.Sprintf("test-application-appRoleAssignments-2-%s", c.RandomString)),
		RequiredResourceAccess: &[]msgraph.RequiredResourceAccess{
			{
				ResourceAppId: resourceApp.AppId,
				ResourceAccess: &[]msgraph.ResourceAccess{
					{
						ID:   (*resourceApp.AppRoles)[0].ID,
						Type: msgraph.ResourceAccessTypeRole,
					},
				},
			},
		},
	})
	servicePrincipal := testServicePrincipalsClient_Create(t, c, msgraph.ServicePrincipal{
		AccountEnabled: utils.BoolPtr(true),
		AppId:          app.AppId,
		DisplayName:    app.DisplayName,
	})

	// assign app role to group
	groupAssignment := testAppRoleAssignedToClient_Assign(t, c, msgraph.AppRoleAssignment{
		AppRoleId:   (*resourceApp.AppRoles)[0].ID,
		PrincipalId: group.ID(),
		ResourceId:  resourceServicePrincipal.ID(),
	})

	// assign app role to user
	userAssignment := testAppRoleAssignedToClient_Assign(t, c, msgraph.AppRoleAssignment{
		AppRoleId:   (*resourceApp.AppRoles)[0].ID,
		PrincipalId: user.ID(),
		ResourceId:  resourceServicePrincipal.ID(),
	})

	// assign app role to service principal
	servicePrincipalAssignment := testAppRoleAssignedToClient_Assign(t, c, msgraph.AppRoleAssignment{
		AppRoleId:   (*resourceApp.AppRoles)[0].ID,
		PrincipalId: servicePrincipal.ID(),
		ResourceId:  resourceServicePrincipal.ID(),
	})

	// list app roles assigned to resource service principal
	testAppRoleAssignedToClient_List(t, c, *resourceServicePrincipal.ID())

	// remove the assigned app roles
	testAppRoleAssignedToClient_Remove(t, c, *resourceServicePrincipal.ID(), *groupAssignment.Id)
	testAppRoleAssignedToClient_Remove(t, c, *resourceServicePrincipal.ID(), *userAssignment.Id)
	testAppRoleAssignedToClient_Remove(t, c, *resourceServicePrincipal.ID(), *servicePrincipalAssignment.Id)

	// clean up
	testGroupsClient_Delete(t, c, *group.ID())
	testUsersClient_Delete(t, c, *user.ID())
	testServicePrincipalsClient_Delete(t, c, *servicePrincipal.ID())
	testServicePrincipalsClient_Delete(t, c, *resourceServicePrincipal.ID())
	testApplicationsClient_Delete(t, c, *resourceApp.ID())
}

func TestGroupsAppRoleAssignmentsClient(t *testing.T) {
	c := test.NewTest(t)
	defer c.CancelFunc()

	// create a new test group
	newGroup := msgraph.Group{
		DisplayName:     utils.StringPtr("test-group-appRoleAssignments"),
		MailEnabled:     utils.BoolPtr(false),
		MailNickname:    utils.StringPtr(fmt.Sprintf("test-group-%s", c.RandomString)),
		SecurityEnabled: utils.BoolPtr(true),
	}
	group := testGroupsClient_Create(t, c, newGroup)

	// pre-generate uuid for a test resourceApp role
	testResourceAppRoleId, _ := uuid.GenerateUUID()
	// create a new test application with a test resourceApp role
	resourceApp := testApplicationsClient_Create(t, c, msgraph.Application{
		DisplayName: utils.StringPtr(fmt.Sprintf("test-application-appRoleAssignments-%s", c.RandomString)),
		AppRoles: &[]msgraph.AppRole{
			{
				ID:          utils.StringPtr(testResourceAppRoleId),
				DisplayName: utils.StringPtr(fmt.Sprintf("test-resourceApp-role-%s", c.RandomString)),
				IsEnabled:   utils.BoolPtr(true),
				Description: utils.StringPtr(fmt.Sprintf("test-resourceApp-role-description-%s", c.RandomString)),
				Value:       utils.StringPtr(fmt.Sprintf("test-resourceApp-role-value-%s", c.RandomString)),
				AllowedMemberTypes: &[]msgraph.AppRoleAllowedMemberType{
					msgraph.AppRoleAllowedMemberTypeUser,
					msgraph.AppRoleAllowedMemberTypeApplication,
				},
			},
		},
	})

	// create a new test resource (API) service principal which has defined the resourceApp role (the application permission)
	resourceServicePrincipal := testServicePrincipalsClient_Create(t, c, msgraph.ServicePrincipal{
		AccountEnabled: utils.BoolPtr(true),
		AppId:          resourceApp.AppId,
		// display name needs to match resourceApp's display name
		DisplayName: resourceApp.DisplayName,
	})

	// assign resourceApp role to the test group
	appRoleAssignment := testGroupsAppRoleAssignmentsClient_Assign(t, c, *group.ID(), *resourceServicePrincipal.ID(), testResourceAppRoleId)

	// list resourceApp role assignments for a test group
	appRoleAssignments := testGroupsAppRoleAssignmentsClient_List(t, c, *group.ID())
	if len(*appRoleAssignments) == 0 {
		t.Fatal("expected at least one resourceApp role assignment assigned to the test group")
	}

	// removes resourceApp role assignment previously set to the test group
	testGroupsAppRoleAssignmentsClient_Remove(t, c, *group.ID(), *appRoleAssignment.Id)

	// remove all test resources to clean up
	testGroupsClient_Delete(t, c, *group.ID())
	testServicePrincipalsClient_Delete(t, c, *resourceServicePrincipal.ID())
	testApplicationsClient_Delete(t, c, *resourceApp.ID())
}

func TestUsersAppRoleAssignmentsClient(t *testing.T) {
	c := test.NewTest(t)
	defer c.CancelFunc()

	// create a new test user
	newUser := msgraph.User{
		AccountEnabled:    utils.BoolPtr(true),
		DisplayName:       utils.StringPtr("test-user-appRoleAssignments"),
		MailNickname:      utils.StringPtr(fmt.Sprintf("test-user-%s", c.RandomString)),
		UserPrincipalName: utils.StringPtr(fmt.Sprintf("test-user-%s@%s", c.RandomString, c.Connections["default"].DomainName)),
		PasswordProfile: &msgraph.UserPasswordProfile{
			Password: utils.StringPtr(fmt.Sprintf("IrPa55w0rd%s", c.RandomString)),
		},
	}
	user := testUsersClient_Create(t, c, newUser)

	// pre-generate uuid for a test resourceApp role
	testResourceAppRoleId, _ := uuid.GenerateUUID()
	// create a new test application with a test resourceApp role
	resourceApp := testApplicationsClient_Create(t, c, msgraph.Application{
		DisplayName: utils.StringPtr(fmt.Sprintf("test-application-appRoleAssignments-%s", c.RandomString)),
		AppRoles: &[]msgraph.AppRole{
			{
				ID:          utils.StringPtr(testResourceAppRoleId),
				DisplayName: utils.StringPtr(fmt.Sprintf("test-resourceApp-role-%s", c.RandomString)),
				IsEnabled:   utils.BoolPtr(true),
				Description: utils.StringPtr(fmt.Sprintf("test-resourceApp-role-description-%s", c.RandomString)),
				Value:       utils.StringPtr(fmt.Sprintf("test-resourceApp-role-value-%s", c.RandomString)),
				AllowedMemberTypes: &[]msgraph.AppRoleAllowedMemberType{
					msgraph.AppRoleAllowedMemberTypeUser,
					msgraph.AppRoleAllowedMemberTypeApplication,
				},
			},
		},
	})

	// create a new test resource (API) service principal which has defined the resourceApp role (the application permission)
	resourceServicePrincipal := testServicePrincipalsClient_Create(t, c, msgraph.ServicePrincipal{
		AccountEnabled: utils.BoolPtr(true),
		AppId:          resourceApp.AppId,
		// display name needs to match resourceApp's display name
		DisplayName: resourceApp.DisplayName,
	})

	// assign resourceApp role to the test user
	appRoleAssignment := testUsersAppRoleAssignmentsClient_Assign(t, c, *user.ID(), *resourceServicePrincipal.ID(), testResourceAppRoleId)

	// list resourceApp role assignments for a test user
	appRoleAssignments := testUsersAppRoleAssignmentsClient_List(t, c, *user.ID())
	if len(*appRoleAssignments) == 0 {
		t.Fatal("expected at least one resourceApp role assignment assigned to the test user")
	}

	// removes resourceApp role assignment previously set to the test user
	testUsersAppRoleAssignmentsClient_Remove(t, c, *user.ID(), *appRoleAssignment.Id)

	// remove all test resources to clean up
	testUsersClient_Delete(t, c, *user.ID())
	testServicePrincipalsClient_Delete(t, c, *resourceServicePrincipal.ID())
	testApplicationsClient_Delete(t, c, *resourceApp.ID())
}

func TestServicePrincipalsAppRoleAssignmentsClient(t *testing.T) {
	c := test.NewTest(t)
	defer c.CancelFunc()

	// pre-generate uuid for a test resourceApp role
	testResourceAppRoleId, _ := uuid.GenerateUUID()
	// create a new test application with a test resourceApp role
	resourceApp := testApplicationsClient_Create(t, c, msgraph.Application{
		DisplayName: utils.StringPtr(fmt.Sprintf("test-application-appRoleAssignments-%s", c.RandomString)),
		AppRoles: &[]msgraph.AppRole{
			{
				ID:          utils.StringPtr(testResourceAppRoleId),
				DisplayName: utils.StringPtr(fmt.Sprintf("test-resourceApp-role-%s", c.RandomString)),
				IsEnabled:   utils.BoolPtr(true),
				Description: utils.StringPtr(fmt.Sprintf("test-resourceApp-role-description-%s", c.RandomString)),
				Value:       utils.StringPtr(fmt.Sprintf("test-resourceApp-role-value-%s", c.RandomString)),
				AllowedMemberTypes: &[]msgraph.AppRoleAllowedMemberType{
					msgraph.AppRoleAllowedMemberTypeUser,
					msgraph.AppRoleAllowedMemberTypeApplication,
				},
			},
		},
	})

	// create a new test resource (API) service principal which has defined the resourceApp role (the application permission)
	resourceServicePrincipal := testServicePrincipalsClient_Create(t, c, msgraph.ServicePrincipal{
		AccountEnabled: utils.BoolPtr(true),
		AppId:          resourceApp.AppId,
		// display name needs to match resourceApp's display name
		DisplayName: resourceApp.DisplayName,
	})

	// create a new test 2 application
	clientApp := testApplicationsClient_Create(t, c, msgraph.Application{
		DisplayName: utils.StringPtr(fmt.Sprintf("test-application-appRoleAssignments-2-%s", c.RandomString)),
		AppRoles: &[]msgraph.AppRole{
			{
				ID:          utils.StringPtr(testResourceAppRoleId),
				DisplayName: utils.StringPtr(fmt.Sprintf("test-2-resourceApp-role-%s", c.RandomString)),
				IsEnabled:   utils.BoolPtr(true),
				Description: utils.StringPtr(fmt.Sprintf("test-2-resourceApp-role-description-%s", c.RandomString)),
				Value:       utils.StringPtr(fmt.Sprintf("test-2-resourceApp-role-value-%s", c.RandomString)),
				AllowedMemberTypes: &[]msgraph.AppRoleAllowedMemberType{
					msgraph.AppRoleAllowedMemberTypeUser,
					msgraph.AppRoleAllowedMemberTypeApplication,
				},
			},
		},
	})
	// create a new test client service principal
	clientServicePrincipal := testServicePrincipalsClient_Create(t, c, msgraph.ServicePrincipal{
		AccountEnabled: utils.BoolPtr(true),
		AppId:          clientApp.AppId,
		// display name needs to match clientApp's display name
		DisplayName: clientApp.DisplayName,
	})

	// assign resourceApp role to the test client service principal
	appRoleAssignment := testServicePrincipalsAppRoleAssignmentsClient_Assign(t, c, *clientServicePrincipal.ID(), *resourceServicePrincipal.ID(), testResourceAppRoleId)

	// list resourceApp role assignments for a test client service principal
	appRoleAssignments := testServicePrincipalsAppRoleAssignmentsClient_List(t, c, *clientServicePrincipal.ID())
	if len(*appRoleAssignments) == 0 {
		t.Fatal("expected at least one resourceApp role assignment assigned to the test client service principal")
	}

	// removes resourceApp role assignment previously set to the test client service principal
	testServicePrincipalsAppRoleAssignmentsClient_Remove(t, c, *clientServicePrincipal.ID(), *appRoleAssignment.Id)

	// remove all test resources to clean up
	testServicePrincipalsClient_Delete(t, c, *clientServicePrincipal.ID())
	testServicePrincipalsClient_Delete(t, c, *resourceServicePrincipal.ID())
	testApplicationsClient_Delete(t, c, *resourceApp.ID())
}

func testGroupsAppRoleAssignmentsClient_List(t *testing.T, c *test.Test, id string) (appRoleAssignments *[]msgraph.AppRoleAssignment) {
	appRoleAssignments, _, err := c.GroupsAppRoleAssignmentsClient.List(c.Context, id, odata.Query{})
	if err != nil {
		t.Fatalf("AppRoleAssignmentsClient.List(): %v", err)
	}
	if appRoleAssignments == nil {
		t.Fatal("AppRoleAssignmentsClient.List(): appRoleAssignments was nil")
	}
	return
}

func testServicePrincipalsAppRoleAssignmentsClient_List(t *testing.T, c *test.Test, id string) (appRoleAssignments *[]msgraph.AppRoleAssignment) {
	appRoleAssignments, _, err := c.ServicePrincipalsAppRoleAssignmentsClient.List(c.Context, id, odata.Query{})
	if err != nil {
		t.Fatalf("AppRoleAssignmentsClient.List(): %v", err)
	}
	if appRoleAssignments == nil {
		t.Fatal("AppRoleAssignmentsClient.List(): appRoleAssignments was nil")
	}
	return
}

func testUsersAppRoleAssignmentsClient_List(t *testing.T, c *test.Test, id string) (appRoleAssignments *[]msgraph.AppRoleAssignment) {
	appRoleAssignments, _, err := c.UsersAppRoleAssignmentsClient.List(c.Context, id, odata.Query{})
	if err != nil {
		t.Fatalf("AppRoleAssignmentsClient.List(): %v", err)
	}
	if appRoleAssignments == nil {
		t.Fatal("AppRoleAssignmentsClient.List(): appRoleAssignments was nil")
	}
	return
}

func testGroupsAppRoleAssignmentsClient_Remove(t *testing.T, c *test.Test, id, appRoleAssignmentId string) {
	status, err := c.GroupsAppRoleAssignmentsClient.Remove(c.Context, id, appRoleAssignmentId)
	if err != nil {
		t.Fatalf("AppRoleAssignmentsClient.Remove(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AppRoleAssignmentsClient.Remove(): invalid status: %d", status)
	}
}

func testServicePrincipalsAppRoleAssignmentsClient_Remove(t *testing.T, c *test.Test, id, appRoleAssignmentId string) {
	status, err := c.ServicePrincipalsAppRoleAssignmentsClient.Remove(c.Context, id, appRoleAssignmentId)
	if err != nil {
		t.Fatalf("AppRoleAssignmentsClient.Remove(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AppRoleAssignmentsClient.Remove(): invalid status: %d", status)
	}
}

func testUsersAppRoleAssignmentsClient_Remove(t *testing.T, c *test.Test, id, appRoleAssignmentId string) {
	status, err := c.UsersAppRoleAssignmentsClient.Remove(c.Context, id, appRoleAssignmentId)
	if err != nil {
		t.Fatalf("AppRoleAssignmentsClient.Remove(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AppRoleAssignmentsClient.Remove(): invalid status: %d", status)
	}
}

func testGroupsAppRoleAssignmentsClient_Assign(t *testing.T, c *test.Test, principalId, resourceServicePrincipalId, appRoleId string) (appRoleAssignment *msgraph.AppRoleAssignment) {
	appRoleAssignment, status, err := c.GroupsAppRoleAssignmentsClient.Assign(c.Context, principalId, resourceServicePrincipalId, appRoleId)
	if err != nil {
		t.Fatalf("AppRoleAssignmentsClient.Assign(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AppRoleAssignmentsClient.Assign(): invalid status: %d", status)
	}
	if appRoleAssignment == nil {
		t.Fatal("AppRoleAssignmentsClient.Assign(): appRoleAssignment was nil")
	}
	if appRoleAssignment.Id == nil {
		t.Fatal("AppRoleAssignmentsClient.Assign(): appRoleAssignment.Id was nil")
	}
	return
}

func testServicePrincipalsAppRoleAssignmentsClient_Assign(t *testing.T, c *test.Test, principalId, resourceServicePrincipalId, appRoleId string) (appRoleAssignment *msgraph.AppRoleAssignment) {
	appRoleAssignment, status, err := c.ServicePrincipalsAppRoleAssignmentsClient.Assign(c.Context, principalId, resourceServicePrincipalId, appRoleId)
	if err != nil {
		t.Fatalf("AppRoleAssignmentsClient.Assign(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AppRoleAssignmentsClient.Assign(): invalid status: %d", status)
	}
	if appRoleAssignment == nil {
		t.Fatal("AppRoleAssignmentsClient.Assign(): appRoleAssignment was nil")
	}
	if appRoleAssignment.Id == nil {
		t.Fatal("AppRoleAssignmentsClient.Assign(): appRoleAssignment.Id was nil")
	}
	return
}

func testUsersAppRoleAssignmentsClient_Assign(t *testing.T, c *test.Test, principalId, resourceServicePrincipalId, appRoleId string) (appRoleAssignment *msgraph.AppRoleAssignment) {
	appRoleAssignment, status, err := c.UsersAppRoleAssignmentsClient.Assign(c.Context, principalId, resourceServicePrincipalId, appRoleId)
	if err != nil {
		t.Fatalf("AppRoleAssignmentsClient.Assign(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AppRoleAssignmentsClient.Assign(): invalid status: %d", status)
	}
	if appRoleAssignment == nil {
		t.Fatal("AppRoleAssignmentsClient.Assign(): appRoleAssignment was nil")
	}
	if appRoleAssignment.Id == nil {
		t.Fatal("AppRoleAssignmentsClient.Assign(): appRoleAssignment.Id was nil")
	}
	return
}

func testAppRoleAssignedToClient_List(t *testing.T, c *test.Test, resourceAppId string) (appRoleAssignments *[]msgraph.AppRoleAssignment) {
	appRoleAssignments, status, err := c.AppRoleAssignedToClient.List(c.Context, resourceAppId, odata.Query{})
	if err != nil {
		t.Fatalf("AppRoleAssignedToClient.List(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AppRoleAssignedToClient.List(): invalid status: %d", status)
	}
	if appRoleAssignments == nil {
		t.Fatal("AppRoleAssignedToClient.List(): appRoleAssignments was nil")
	}
	if len(*appRoleAssignments) == 0 {
		t.Fatal("AppRoleAssignedToClient.List(): appRoleAssignments was empty")
	}
	return
}

func testAppRoleAssignedToClient_Assign(t *testing.T, c *test.Test, appRoleAssignment msgraph.AppRoleAssignment) (newAppRoleAssignment *msgraph.AppRoleAssignment) {
	newAppRoleAssignment, status, err := c.AppRoleAssignedToClient.Assign(c.Context, appRoleAssignment)
	if err != nil {
		t.Fatalf("AppRoleAssignedToClient.Assign(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AppRoleAssignedToClient.Assign(): invalid status: %d", status)
	}
	if newAppRoleAssignment == nil {
		t.Fatal("AppRoleAssignedToClient.Assign(): newAppRoleAssignment was nil")
	}
	if newAppRoleAssignment.Id == nil {
		t.Fatal("AppRoleAssignedToClient.Assign(): newAppRoleAssignment.Id was nil")
	}
	return
}

func testAppRoleAssignedToClient_Remove(t *testing.T, c *test.Test, resourceAppId, appRoleAssignmentId string) {
	status, err := c.AppRoleAssignedToClient.Remove(c.Context, resourceAppId, appRoleAssignmentId)
	if err != nil {
		t.Fatalf("AppRoleAssignedToClient.Remove(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AppRoleAssignedToClient.Remove(): invalid status: %d", status)
	}
}
