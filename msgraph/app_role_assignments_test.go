package msgraph_test

import (
	"fmt"
	"testing"

	"github.com/manicminer/hamilton/odata"

	"github.com/hashicorp/go-uuid"

	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/msgraph"
)

type AppRoleAssignedToClientTest struct {
	connection   *test.Connection
	client       *msgraph.AppRoleAssignedToClient
	randomString string
}

type AppRoleAssignmentsClientTest struct {
	connection   *test.Connection
	client       *msgraph.AppRoleAssignmentsClient
	randomString string
}

func TestAppRoleAssignedToClient(t *testing.T) {
	rs := test.RandomString()

	c := AppRoleAssignedToClientTest{
		connection:   test.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: rs,
	}
	c.client = msgraph.NewAppRoleAssignedToClient(c.connection.AuthConfig.TenantID)
	c.client.BaseClient.Authorizer = c.connection.Authorizer
	c.client.BaseClient.Endpoint = c.connection.AuthConfig.Environment.MsGraph.Endpoint

	a := ApplicationsClientTest{
		connection:   test.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: rs,
	}
	a.client = msgraph.NewApplicationsClient(a.connection.AuthConfig.TenantID)
	a.client.BaseClient.Authorizer = a.connection.Authorizer
	a.client.BaseClient.Endpoint = a.connection.AuthConfig.Environment.MsGraph.Endpoint

	s := ServicePrincipalsClientTest{
		connection:   test.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: rs,
	}
	s.client = msgraph.NewServicePrincipalsClient(s.connection.AuthConfig.TenantID)
	s.client.BaseClient.Authorizer = s.connection.Authorizer
	s.client.BaseClient.Endpoint = s.connection.AuthConfig.Environment.MsGraph.Endpoint

	g := GroupsClientTest{
		connection:   test.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: rs,
	}
	g.client = msgraph.NewGroupsClient(g.connection.AuthConfig.TenantID)
	g.client.BaseClient.Authorizer = g.connection.Authorizer
	g.client.BaseClient.Endpoint = g.connection.AuthConfig.Environment.MsGraph.Endpoint

	u := UsersClientTest{
		connection:   test.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: rs,
	}
	u.client = msgraph.NewUsersClient(u.connection.AuthConfig.TenantID)
	u.client.BaseClient.Authorizer = u.connection.Authorizer
	u.client.BaseClient.Endpoint = u.connection.AuthConfig.Environment.MsGraph.Endpoint

	// Scaffold the resource application
	testResourceAppRoleId, _ := uuid.GenerateUUID()
	resourceApp := testApplicationsClient_Create(t, a, msgraph.Application{
		DisplayName: utils.StringPtr(fmt.Sprintf("test-application-appRoleAssignedTo-%s", a.randomString)),
		AppRoles: &[]msgraph.AppRole{
			{
				ID:          utils.StringPtr(testResourceAppRoleId),
				DisplayName: utils.StringPtr(fmt.Sprintf("test-resourceApp-role-%s", a.randomString)),
				IsEnabled:   utils.BoolPtr(true),
				Description: utils.StringPtr(fmt.Sprintf("test-resourceApp-role-description-%s", a.randomString)),
				Value:       utils.StringPtr(fmt.Sprintf("test-resourceApp-role-value-%s", a.randomString)),
				AllowedMemberTypes: &[]msgraph.AppRoleAllowedMemberType{
					msgraph.AppRoleAllowedMemberTypeApplication,
					msgraph.AppRoleAllowedMemberTypeUser,
				},
			},
		},
	})
	resourceServicePrincipal := testServicePrincipalsClient_Create(t, s, msgraph.ServicePrincipal{
		AccountEnabled:            utils.BoolPtr(true),
		AppId:                     resourceApp.AppId,
		AppRoleAssignmentRequired: utils.BoolPtr(true),
		DisplayName:               resourceApp.DisplayName,
	})

	// Scaffold the group, user and service principal to assign roles to
	group := testGroupsClient_Create(t, g, msgraph.Group{
		DisplayName:     utils.StringPtr("test-group-appRoleAssignments"),
		MailEnabled:     utils.BoolPtr(false),
		MailNickname:    utils.StringPtr(fmt.Sprintf("test-group-%s", g.randomString)),
		SecurityEnabled: utils.BoolPtr(true),
	})
	user := testUsersClient_Create(t, u, msgraph.User{
		AccountEnabled:    utils.BoolPtr(true),
		DisplayName:       utils.StringPtr("test-user-appRoleAssignments"),
		MailNickname:      utils.StringPtr(fmt.Sprintf("test-user-%s", u.randomString)),
		UserPrincipalName: utils.StringPtr(fmt.Sprintf("test-user-%s@%s", u.randomString, u.connection.DomainName)),
		PasswordProfile: &msgraph.UserPasswordProfile{
			Password: utils.StringPtr(fmt.Sprintf("IrPa55w0rd%s", u.randomString)),
		},
	})
	app := testApplicationsClient_Create(t, a, msgraph.Application{
		DisplayName: utils.StringPtr(fmt.Sprintf("test-application-appRoleAssignments-2-%s", a.randomString)),
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
	servicePrincipal := testServicePrincipalsClient_Create(t, s, msgraph.ServicePrincipal{
		AccountEnabled: utils.BoolPtr(true),
		AppId:          app.AppId,
		DisplayName:    app.DisplayName,
	})

	// assign app role to group
	groupAssignment := testAppRoleAssignedToClient_Assign(t, c, msgraph.AppRoleAssignment{
		AppRoleId:   (*resourceApp.AppRoles)[0].ID,
		PrincipalId: group.ID,
		ResourceId:  resourceServicePrincipal.ID,
	})

	// assign app role to user
	userAssignment := testAppRoleAssignedToClient_Assign(t, c, msgraph.AppRoleAssignment{
		AppRoleId:   (*resourceApp.AppRoles)[0].ID,
		PrincipalId: user.ID,
		ResourceId:  resourceServicePrincipal.ID,
	})

	// assign app role to service principal
	servicePrincipalAssignment := testAppRoleAssignedToClient_Assign(t, c, msgraph.AppRoleAssignment{
		AppRoleId:   (*resourceApp.AppRoles)[0].ID,
		PrincipalId: servicePrincipal.ID,
		ResourceId:  resourceServicePrincipal.ID,
	})

	// list app roles assigned to resource service principal
	testAppRoleAssignedToClient_List(t, c, *resourceServicePrincipal.ID)

	// remove the assigned app roles
	testAppRoleAssignedToClient_Remove(t, c, *resourceServicePrincipal.ID, *groupAssignment.Id)
	testAppRoleAssignedToClient_Remove(t, c, *resourceServicePrincipal.ID, *userAssignment.Id)
	testAppRoleAssignedToClient_Remove(t, c, *resourceServicePrincipal.ID, *servicePrincipalAssignment.Id)

	// clean up
	testGroupsClient_Delete(t, g, *group.ID)
	testUsersClient_Delete(t, u, *user.ID)
	testServicePrincipalsClient_Delete(t, s, *servicePrincipal.ID)
	testServicePrincipalsClient_Delete(t, s, *resourceServicePrincipal.ID)
	testApplicationsClient_Delete(t, a, *resourceApp.ID)
}

func TestGroupsAppRoleAssignmentsClient(t *testing.T) {
	rs := test.RandomString()
	// setup service principal test client
	servicePrincipalsClient := ServicePrincipalsClientTest{
		connection:   test.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: rs,
	}
	servicePrincipalsClient.client = msgraph.NewServicePrincipalsClient(servicePrincipalsClient.connection.AuthConfig.TenantID)
	servicePrincipalsClient.client.BaseClient.Authorizer = servicePrincipalsClient.connection.Authorizer
	servicePrincipalsClient.client.BaseClient.Endpoint = servicePrincipalsClient.connection.AuthConfig.Environment.MsGraph.Endpoint

	// setup groups test client
	groupsClient := GroupsClientTest{
		connection:   test.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: rs,
	}
	groupsClient.client = msgraph.NewGroupsClient(groupsClient.connection.AuthConfig.TenantID)
	groupsClient.client.BaseClient.Authorizer = groupsClient.connection.Authorizer
	groupsClient.client.BaseClient.Endpoint = groupsClient.connection.AuthConfig.Environment.MsGraph.Endpoint

	// setup resourceApp role assignments test client
	appRoleAssignClient := AppRoleAssignmentsClientTest{
		connection:   test.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: rs,
	}
	appRoleAssignClient.client = msgraph.NewGroupsAppRoleAssignmentsClient(appRoleAssignClient.connection.AuthConfig.TenantID)
	appRoleAssignClient.client.BaseClient.Authorizer = appRoleAssignClient.connection.Authorizer
	appRoleAssignClient.client.BaseClient.Endpoint = appRoleAssignClient.connection.AuthConfig.Environment.MsGraph.Endpoint

	// setup applications test client
	appClient := ApplicationsClientTest{
		connection:   test.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: rs,
	}
	appClient.client = msgraph.NewApplicationsClient(appClient.connection.AuthConfig.TenantID)
	appClient.client.BaseClient.Authorizer = appClient.connection.Authorizer
	appClient.client.BaseClient.Endpoint = appClient.connection.AuthConfig.Environment.MsGraph.Endpoint

	// create a new test group
	newGroup := msgraph.Group{
		DisplayName:     utils.StringPtr("test-group-appRoleAssignments"),
		MailEnabled:     utils.BoolPtr(false),
		MailNickname:    utils.StringPtr(fmt.Sprintf("test-group-%s", groupsClient.randomString)),
		SecurityEnabled: utils.BoolPtr(true),
	}
	group := testGroupsClient_Create(t, groupsClient, newGroup)

	// pre-generate uuid for a test resourceApp role
	testResourceAppRoleId, _ := uuid.GenerateUUID()
	// create a new test application with a test resourceApp role
	resourceApp := testApplicationsClient_Create(t, appClient, msgraph.Application{
		DisplayName: utils.StringPtr(fmt.Sprintf("test-application-appRoleAssignments-%s", appClient.randomString)),
		AppRoles: &[]msgraph.AppRole{
			{
				ID:          utils.StringPtr(testResourceAppRoleId),
				DisplayName: utils.StringPtr(fmt.Sprintf("test-resourceApp-role-%s", appClient.randomString)),
				IsEnabled:   utils.BoolPtr(true),
				Description: utils.StringPtr(fmt.Sprintf("test-resourceApp-role-description-%s", appClient.randomString)),
				Value:       utils.StringPtr(fmt.Sprintf("test-resourceApp-role-value-%s", appClient.randomString)),
				AllowedMemberTypes: &[]msgraph.AppRoleAllowedMemberType{
					msgraph.AppRoleAllowedMemberTypeUser,
					msgraph.AppRoleAllowedMemberTypeApplication,
				},
			},
		},
	})

	// create a new test resource (API) service principal which has defined the resourceApp role (the application permission)
	resourceServicePrincipal := testServicePrincipalsClient_Create(t, servicePrincipalsClient, msgraph.ServicePrincipal{
		AccountEnabled: utils.BoolPtr(true),
		AppId:          resourceApp.AppId,
		// display name needs to match resourceApp's display name
		DisplayName: resourceApp.DisplayName,
	})

	// assign resourceApp role to the test group
	appRoleAssignment := testAppRoleAssignmentsClient_Assign(t, appRoleAssignClient, *group.ID, *resourceServicePrincipal.ID, testResourceAppRoleId)

	// list resourceApp role assignments for a test group
	appRoleAssignments := testAppRoleAssignmentsClient_List(t, appRoleAssignClient, *group.ID)
	if len(*appRoleAssignments) == 0 {
		t.Fatal("expected at least one resourceApp role assignment assigned to the test group")
	}

	// removes resourceApp role assignment previously set to the test group
	testAppRoleAssignmentsClient_Remove(t, appRoleAssignClient, *group.ID, *appRoleAssignment.Id)

	// remove all test resources to clean up
	testGroupsClient_Delete(t, groupsClient, *group.ID)
	testServicePrincipalsClient_Delete(t, servicePrincipalsClient, *resourceServicePrincipal.ID)
	testApplicationsClient_Delete(t, appClient, *resourceApp.ID)
}

func TestUsersAppRoleAssignmentsClient(t *testing.T) {
	rs := test.RandomString()
	// setup service principal test client
	servicePrincipalsClient := ServicePrincipalsClientTest{
		connection:   test.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: rs,
	}
	servicePrincipalsClient.client = msgraph.NewServicePrincipalsClient(servicePrincipalsClient.connection.AuthConfig.TenantID)
	servicePrincipalsClient.client.BaseClient.Authorizer = servicePrincipalsClient.connection.Authorizer
	servicePrincipalsClient.client.BaseClient.Endpoint = servicePrincipalsClient.connection.AuthConfig.Environment.MsGraph.Endpoint

	// setup users test client
	usersClient := UsersClientTest{
		connection:   test.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: rs,
	}
	usersClient.client = msgraph.NewUsersClient(usersClient.connection.AuthConfig.TenantID)
	usersClient.client.BaseClient.Authorizer = usersClient.connection.Authorizer
	usersClient.client.BaseClient.Endpoint = usersClient.connection.AuthConfig.Environment.MsGraph.Endpoint

	// setup resourceApp role assignments test client
	appRoleAssignClient := AppRoleAssignmentsClientTest{
		connection:   test.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: rs,
	}
	appRoleAssignClient.client = msgraph.NewUsersAppRoleAssignmentsClient(appRoleAssignClient.connection.AuthConfig.TenantID)
	appRoleAssignClient.client.BaseClient.Authorizer = appRoleAssignClient.connection.Authorizer
	appRoleAssignClient.client.BaseClient.Endpoint = appRoleAssignClient.connection.AuthConfig.Environment.MsGraph.Endpoint

	// setup applications test client
	appClient := ApplicationsClientTest{
		connection:   test.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: rs,
	}
	appClient.client = msgraph.NewApplicationsClient(appClient.connection.AuthConfig.TenantID)
	appClient.client.BaseClient.Authorizer = appClient.connection.Authorizer
	appClient.client.BaseClient.Endpoint = appClient.connection.AuthConfig.Environment.MsGraph.Endpoint

	// create a new test user
	newUser := msgraph.User{
		AccountEnabled:    utils.BoolPtr(true),
		DisplayName:       utils.StringPtr("test-user-appRoleAssignments"),
		MailNickname:      utils.StringPtr(fmt.Sprintf("test-user-%s", usersClient.randomString)),
		UserPrincipalName: utils.StringPtr(fmt.Sprintf("test-user-%s@%s", usersClient.randomString, usersClient.connection.DomainName)),
		PasswordProfile: &msgraph.UserPasswordProfile{
			Password: utils.StringPtr(fmt.Sprintf("IrPa55w0rd%s", usersClient.randomString)),
		},
	}
	user := testUsersClient_Create(t, usersClient, newUser)

	// pre-generate uuid for a test resourceApp role
	testResourceAppRoleId, _ := uuid.GenerateUUID()
	// create a new test application with a test resourceApp role
	resourceApp := testApplicationsClient_Create(t, appClient, msgraph.Application{
		DisplayName: utils.StringPtr(fmt.Sprintf("test-application-appRoleAssignments-%s", appClient.randomString)),
		AppRoles: &[]msgraph.AppRole{
			{
				ID:          utils.StringPtr(testResourceAppRoleId),
				DisplayName: utils.StringPtr(fmt.Sprintf("test-resourceApp-role-%s", appClient.randomString)),
				IsEnabled:   utils.BoolPtr(true),
				Description: utils.StringPtr(fmt.Sprintf("test-resourceApp-role-description-%s", appClient.randomString)),
				Value:       utils.StringPtr(fmt.Sprintf("test-resourceApp-role-value-%s", appClient.randomString)),
				AllowedMemberTypes: &[]msgraph.AppRoleAllowedMemberType{
					msgraph.AppRoleAllowedMemberTypeUser,
					msgraph.AppRoleAllowedMemberTypeApplication,
				},
			},
		},
	})

	// create a new test resource (API) service principal which has defined the resourceApp role (the application permission)
	resourceServicePrincipal := testServicePrincipalsClient_Create(t, servicePrincipalsClient, msgraph.ServicePrincipal{
		AccountEnabled: utils.BoolPtr(true),
		AppId:          resourceApp.AppId,
		// display name needs to match resourceApp's display name
		DisplayName: resourceApp.DisplayName,
	})

	// assign resourceApp role to the test user
	appRoleAssignment := testAppRoleAssignmentsClient_Assign(t, appRoleAssignClient, *user.ID, *resourceServicePrincipal.ID, testResourceAppRoleId)

	// list resourceApp role assignments for a test user
	appRoleAssignments := testAppRoleAssignmentsClient_List(t, appRoleAssignClient, *user.ID)
	if len(*appRoleAssignments) == 0 {
		t.Fatal("expected at least one resourceApp role assignment assigned to the test user")
	}

	// removes resourceApp role assignment previously set to the test user
	testAppRoleAssignmentsClient_Remove(t, appRoleAssignClient, *user.ID, *appRoleAssignment.Id)

	// remove all test resources to clean up
	testUsersClient_Delete(t, usersClient, *user.ID)
	testServicePrincipalsClient_Delete(t, servicePrincipalsClient, *resourceServicePrincipal.ID)
	testApplicationsClient_Delete(t, appClient, *resourceApp.ID)
}

func TestServicePrincipalsAppRoleAssignmentsClient(t *testing.T) {
	rs := test.RandomString()
	// setup service principal test client
	servicePrincipalsClient := ServicePrincipalsClientTest{
		connection:   test.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: rs,
	}
	servicePrincipalsClient.client = msgraph.NewServicePrincipalsClient(servicePrincipalsClient.connection.AuthConfig.TenantID)
	servicePrincipalsClient.client.BaseClient.Authorizer = servicePrincipalsClient.connection.Authorizer
	servicePrincipalsClient.client.BaseClient.Endpoint = servicePrincipalsClient.connection.AuthConfig.Environment.MsGraph.Endpoint

	// setup resourceApp role assignments test client
	appRoleAssignClient := AppRoleAssignmentsClientTest{
		connection:   test.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: rs,
	}
	appRoleAssignClient.client = msgraph.NewServicePrincipalsAppRoleAssignmentsClient(appRoleAssignClient.connection.AuthConfig.TenantID)
	appRoleAssignClient.client.BaseClient.Authorizer = appRoleAssignClient.connection.Authorizer
	appRoleAssignClient.client.BaseClient.Endpoint = appRoleAssignClient.connection.AuthConfig.Environment.MsGraph.Endpoint

	// setup applications test client
	appClient := ApplicationsClientTest{
		connection:   test.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: rs,
	}
	appClient.client = msgraph.NewApplicationsClient(appClient.connection.AuthConfig.TenantID)
	appClient.client.BaseClient.Authorizer = appClient.connection.Authorizer
	appClient.client.BaseClient.Endpoint = appClient.connection.AuthConfig.Environment.MsGraph.Endpoint

	// pre-generate uuid for a test resourceApp role
	testResourceAppRoleId, _ := uuid.GenerateUUID()
	// create a new test application with a test resourceApp role
	resourceApp := testApplicationsClient_Create(t, appClient, msgraph.Application{
		DisplayName: utils.StringPtr(fmt.Sprintf("test-application-appRoleAssignments-%s", appClient.randomString)),
		AppRoles: &[]msgraph.AppRole{
			{
				ID:          utils.StringPtr(testResourceAppRoleId),
				DisplayName: utils.StringPtr(fmt.Sprintf("test-resourceApp-role-%s", appClient.randomString)),
				IsEnabled:   utils.BoolPtr(true),
				Description: utils.StringPtr(fmt.Sprintf("test-resourceApp-role-description-%s", appClient.randomString)),
				Value:       utils.StringPtr(fmt.Sprintf("test-resourceApp-role-value-%s", appClient.randomString)),
				AllowedMemberTypes: &[]msgraph.AppRoleAllowedMemberType{
					msgraph.AppRoleAllowedMemberTypeUser,
					msgraph.AppRoleAllowedMemberTypeApplication,
				},
			},
		},
	})

	// create a new test resource (API) service principal which has defined the resourceApp role (the application permission)
	resourceServicePrincipal := testServicePrincipalsClient_Create(t, servicePrincipalsClient, msgraph.ServicePrincipal{
		AccountEnabled: utils.BoolPtr(true),
		AppId:          resourceApp.AppId,
		// display name needs to match resourceApp's display name
		DisplayName: resourceApp.DisplayName,
	})

	// create a new test 2 application
	clientApp := testApplicationsClient_Create(t, appClient, msgraph.Application{
		DisplayName: utils.StringPtr(fmt.Sprintf("test-application-appRoleAssignments-2-%s", appClient.randomString)),
		AppRoles: &[]msgraph.AppRole{
			{
				ID:          utils.StringPtr(testResourceAppRoleId),
				DisplayName: utils.StringPtr(fmt.Sprintf("test-2-resourceApp-role-%s", appClient.randomString)),
				IsEnabled:   utils.BoolPtr(true),
				Description: utils.StringPtr(fmt.Sprintf("test-2-resourceApp-role-description-%s", appClient.randomString)),
				Value:       utils.StringPtr(fmt.Sprintf("test-2-resourceApp-role-value-%s", appClient.randomString)),
				AllowedMemberTypes: &[]msgraph.AppRoleAllowedMemberType{
					msgraph.AppRoleAllowedMemberTypeUser,
					msgraph.AppRoleAllowedMemberTypeApplication,
				},
			},
		},
	})
	// create a new test client service principal
	clientServicePrincipal := testServicePrincipalsClient_Create(t, servicePrincipalsClient, msgraph.ServicePrincipal{
		AccountEnabled: utils.BoolPtr(true),
		AppId:          clientApp.AppId,
		// display name needs to match clientApp's display name
		DisplayName: clientApp.DisplayName,
	})

	// assign resourceApp role to the test client service principal
	appRoleAssignment := testAppRoleAssignmentsClient_Assign(t, appRoleAssignClient, *clientServicePrincipal.ID, *resourceServicePrincipal.ID, testResourceAppRoleId)

	// list resourceApp role assignments for a test client service principal
	appRoleAssignments := testAppRoleAssignmentsClient_List(t, appRoleAssignClient, *clientServicePrincipal.ID)
	if len(*appRoleAssignments) == 0 {
		t.Fatal("expected at least one resourceApp role assignment assigned to the test client service principal")
	}

	// removes resourceApp role assignment previously set to the test client service principal
	testAppRoleAssignmentsClient_Remove(t, appRoleAssignClient, *clientServicePrincipal.ID, *appRoleAssignment.Id)

	// remove all test resources to clean up
	testServicePrincipalsClient_Delete(t, servicePrincipalsClient, *clientServicePrincipal.ID)
	testServicePrincipalsClient_Delete(t, servicePrincipalsClient, *resourceServicePrincipal.ID)
	testApplicationsClient_Delete(t, appClient, *resourceApp.ID)
}

func testAppRoleAssignmentsClient_List(t *testing.T, c AppRoleAssignmentsClientTest, id string) (appRoleAssignments *[]msgraph.AppRoleAssignment) {
	appRoleAssignments, _, err := c.client.List(c.connection.Context, id)
	if err != nil {
		t.Fatalf("AppRoleAssignmentsClient.List(): %v", err)
	}
	if appRoleAssignments == nil {
		t.Fatal("AppRoleAssignmentsClient.List(): appRoleAssignments was nil")
	}
	return
}

func testAppRoleAssignmentsClient_Remove(t *testing.T, c AppRoleAssignmentsClientTest, id, appRoleAssignmentId string) {
	status, err := c.client.Remove(c.connection.Context, id, appRoleAssignmentId)
	if err != nil {
		t.Fatalf("AppRoleAssignmentsClient.Remove(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AppRoleAssignmentsClient.Remove(): invalid status: %d", status)
	}
}

func testAppRoleAssignmentsClient_Assign(t *testing.T, c AppRoleAssignmentsClientTest, principalId, resourceServicePrincipalId, appRoleId string) (appRoleAssignment *msgraph.AppRoleAssignment) {
	appRoleAssignment, status, err := c.client.Assign(c.connection.Context, principalId, resourceServicePrincipalId, appRoleId)
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

func testAppRoleAssignedToClient_List(t *testing.T, c AppRoleAssignedToClientTest, resourceAppId string) (appRoleAssignments *[]msgraph.AppRoleAssignment) {
	appRoleAssignments, status, err := c.client.List(c.connection.Context, resourceAppId, odata.Query{})
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

func testAppRoleAssignedToClient_Assign(t *testing.T, c AppRoleAssignedToClientTest, appRoleAssignment msgraph.AppRoleAssignment) (newAppRoleAssignment *msgraph.AppRoleAssignment) {
	newAppRoleAssignment, status, err := c.client.Assign(c.connection.Context, appRoleAssignment)
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

func testAppRoleAssignedToClient_Remove(t *testing.T, c AppRoleAssignedToClientTest, resourceAppId, appRoleAssignmentId string) {
	status, err := c.client.Remove(c.connection.Context, resourceAppId, appRoleAssignmentId)
	if err != nil {
		t.Fatalf("AppRoleAssignedToClient.Remove(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AppRoleAssignedToClient.Remove(): invalid status: %d", status)
	}
}
