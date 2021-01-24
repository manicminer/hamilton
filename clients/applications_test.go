package clients_test

import (
	"fmt"
	"testing"

	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/clients"
	"github.com/manicminer/hamilton/clients/internal"
	"github.com/manicminer/hamilton/models"
)

type ApplicationsClientTest struct {
	connection   *internal.Connection
	client       *clients.ApplicationsClient
	randomString string
}

func TestApplicationsClient(t *testing.T) {
	c := ApplicationsClientTest{
		connection:   internal.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: internal.RandomString(),
	}
	c.client = clients.NewApplicationsClient(c.connection.AuthConfig.TenantID)
	c.client.BaseClient.Authorizer = c.connection.Authorizer

	app := testApplicationsClient_Create(t, c, models.Application{
		DisplayName: internal.String(fmt.Sprintf("test-application-%s", c.randomString)),
	})
	testApplicationsClient_Get(t, c, *app.ID)
	app.DisplayName = internal.String(fmt.Sprintf("test-app-updated-%s", c.randomString))
	testApplicationsClient_Update(t, c, *app)
	owners := testApplicationsClient_ListOwners(t, c, *app.ID)
	testApplicationsClient_GetOwner(t, c, *app.ID, (*owners)[0])
	testApplicationsClient_RemoveOwners(t, c, *app.ID, owners)
	app.AppendOwner(c.client.BaseClient.Endpoint, c.client.BaseClient.ApiVersion, (*owners)[0])
	testApplicationsClient_AddOwners(t, c, app)
	testApplicationsClient_List(t, c)
	testApplicationsClient_Delete(t, c, *app.ID)
}

func testApplicationsClient_Create(t *testing.T, c ApplicationsClientTest, a models.Application) (application *models.Application) {
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

func testApplicationsClient_Update(t *testing.T, c ApplicationsClientTest, a models.Application) {
	status, err := c.client.Update(c.connection.Context, a)
	if err != nil {
		t.Fatalf("ApplicationsClient.Update(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ApplicationsClient.Update(): invalid status: %d", status)
	}
}

func testApplicationsClient_List(t *testing.T, c ApplicationsClientTest) (applications *[]models.Application) {
	applications, _, err := c.client.List(c.connection.Context, "")
	if err != nil {
		t.Fatalf("ApplicationsClient.List(): %v", err)
	}
	if applications == nil {
		t.Fatal("ApplicationsClient.List(): applications was nil")
	}
	return
}

func testApplicationsClient_Get(t *testing.T, c ApplicationsClientTest, id string) (application *models.Application) {
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

func testApplicationsClient_AddOwners(t *testing.T, c ApplicationsClientTest, a *models.Application) {
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
