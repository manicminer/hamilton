package msgraph_test

import (
	"fmt"
	"testing"

	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/msgraph"
)

type ApplicationsClientTest struct {
	connection   *test.Connection
	client       *msgraph.ApplicationsClient
	randomString string
}

func TestApplicationsClient(t *testing.T) {
	c := ApplicationsClientTest{
		connection:   test.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: test.RandomString(),
	}
	c.client = msgraph.NewApplicationsClient(c.connection.AuthConfig.TenantID)
	c.client.BaseClient.Authorizer = c.connection.Authorizer

	app := testApplicationsClient_Create(t, c, msgraph.Application{
		DisplayName: utils.StringPtr(fmt.Sprintf("test-application-%s", c.randomString)),
	})
	testApplicationsClient_Get(t, c, *app.ID)
	app.DisplayName = utils.StringPtr(fmt.Sprintf("test-app-updated-%s", c.randomString))
	testApplicationsClient_Update(t, c, *app)
	owners := testApplicationsClient_ListOwners(t, c, *app.ID)
	testApplicationsClient_GetOwner(t, c, *app.ID, (*owners)[0])
	testApplicationsClient_RemoveOwners(t, c, *app.ID, owners)
	app.AppendOwner(c.client.BaseClient.Endpoint, c.client.BaseClient.ApiVersion, (*owners)[0])
	testApplicationsClient_AddOwners(t, c, app)
	pwd := testApplicationsClient_AddPassword(t, c, app)
	testApplicationsClient_RemovePassword(t, c, app, pwd)
	testApplicationsClient_List(t, c)
	testApplicationsClient_Delete(t, c, *app.ID)
}

func TestApplicationsClient_groupMembershipClaims(t *testing.T) {
	c := ApplicationsClientTest{
		connection:   test.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: test.RandomString(),
	}
	c.client = msgraph.NewApplicationsClient(c.connection.AuthConfig.TenantID)
	c.client.BaseClient.Authorizer = c.connection.Authorizer

	app := testApplicationsClient_Create(t, c, msgraph.Application{
		DisplayName:           utils.StringPtr(fmt.Sprintf("test-application-%s", c.randomString)),
		GroupMembershipClaims: &[]msgraph.GroupMembershipClaim{"SecurityGroup", "ApplicationGroup"},
	})
	testApplicationsClient_Delete(t, c, *app.ID)
}

func testApplicationsClient_Create(t *testing.T, c ApplicationsClientTest, a msgraph.Application) (application *msgraph.Application) {
	application, status, err := c.client.Create(c.connection.Context, a)
	if err != nil {
		t.Fatalf("ApplicationsClient.Create(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ApplicationsClient.Create(): invalid status: %d", status)
	}
	if application == nil {
		t.Fatal("ApplicationsClient.Create(): application was nil")
	}
	if application.ID == nil {
		t.Fatal("ApplicationsClient.Create(): application.ID was nil")
	}
	return
}

func testApplicationsClient_Update(t *testing.T, c ApplicationsClientTest, a msgraph.Application) {
	status, err := c.client.Update(c.connection.Context, a)
	if err != nil {
		t.Fatalf("ApplicationsClient.Update(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ApplicationsClient.Update(): invalid status: %d", status)
	}
}

func testApplicationsClient_List(t *testing.T, c ApplicationsClientTest) (applications *[]msgraph.Application) {
	applications, _, err := c.client.List(c.connection.Context, "")
	if err != nil {
		t.Fatalf("ApplicationsClient.List(): %v", err)
	}
	if applications == nil {
		t.Fatal("ApplicationsClient.List(): applications was nil")
	}
	return
}

func testApplicationsClient_Get(t *testing.T, c ApplicationsClientTest, id string) (application *msgraph.Application) {
	application, status, err := c.client.Get(c.connection.Context, id)
	if err != nil {
		t.Fatalf("ApplicationsClient.Get(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ApplicationsClient.Get(): invalid status: %d", status)
	}
	if application == nil {
		t.Fatal("ApplicationsClient.Get(): application was nil")
	}
	return
}

func testApplicationsClient_Delete(t *testing.T, c ApplicationsClientTest, id string) {
	status, err := c.client.Delete(c.connection.Context, id)
	if err != nil {
		t.Fatalf("ApplicationsClient.Delete(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ApplicationsClient.Delete(): invalid status: %d", status)
	}
}

func testApplicationsClient_ListOwners(t *testing.T, c ApplicationsClientTest, id string) (owners *[]string) {
	owners, status, err := c.client.ListOwners(c.connection.Context, id)
	if err != nil {
		t.Fatalf("ApplicationsClient.ListOwners(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ApplicationsClient.ListOwners(): invalid status: %d", status)
	}
	if owners == nil {
		t.Fatal("ApplicationsClient.ListOwners(): owners was nil")
	}
	if len(*owners) == 0 {
		t.Fatal("ApplicationsClient.ListOwners(): owners was empty")
	}
	return
}

func testApplicationsClient_GetOwner(t *testing.T, c ApplicationsClientTest, appId string, ownerId string) (owner *string) {
	owner, status, err := c.client.GetOwner(c.connection.Context, appId, ownerId)
	if err != nil {
		t.Fatalf("ApplicationsClient.GetOwner(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ApplicationsClient.GetOwner(): invalid status: %d", status)
	}
	if owner == nil {
		t.Fatal("ApplicationsClient.GetOwner(): owner was nil")
	}
	return
}

func testApplicationsClient_AddOwners(t *testing.T, c ApplicationsClientTest, a *msgraph.Application) {
	status, err := c.client.AddOwners(c.connection.Context, a)
	if err != nil {
		t.Fatalf("ApplicationsClient.AddOwners(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ApplicationsClient.AddOwners(): invalid status: %d", status)
	}
}

func testApplicationsClient_RemoveOwners(t *testing.T, c ApplicationsClientTest, appId string, ownerIds *[]string) {
	status, err := c.client.RemoveOwners(c.connection.Context, appId, ownerIds)
	if err != nil {
		t.Fatalf("ApplicationsClient.RemoveOwners(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ApplicationsClient.RemoveOwners(): invalid status: %d", status)
	}
}

func testApplicationsClient_AddPassword(t *testing.T, c ApplicationsClientTest, a *msgraph.Application) *msgraph.PasswordCredential {
	pwd := msgraph.PasswordCredential{
		DisplayName: utils.StringPtr("test password"),
	}
	newPwd, status, err := c.client.AddPassword(c.connection.Context, *a.ID, pwd)
	if err != nil {
		t.Fatalf("ApplicationsClient.AddPassword(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ApplicationsClient.AddPassword(): invalid status: %d", status)
	}
	if newPwd.SecretText == nil || len(*newPwd.SecretText) == 0 {
		t.Fatalf("ApplicationsClient.AddPassword(): nil or empty secretText returned by API")
	}
	return newPwd
}

func testApplicationsClient_RemovePassword(t *testing.T, c ApplicationsClientTest, a *msgraph.Application, p *msgraph.PasswordCredential) {
	status, err := c.client.RemovePassword(c.connection.Context, *a.ID, *p.KeyId)
	if err != nil {
		t.Fatalf("ApplicationsClient.RemovePassword(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ApplicationsClient.RemovePassword(): invalid status: %d", status)
	}
}
