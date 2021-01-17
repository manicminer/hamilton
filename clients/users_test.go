package clients_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/manicminer/hamilton/clients"
	"github.com/manicminer/hamilton/clients/internal"
	"github.com/manicminer/hamilton/models"
)

type UsersClient struct {
	connection   *internal.Connection
	context      context.Context
	client       *clients.UsersClient
	randomString string
}

func TestUsersClient(t *testing.T) {
	c := UsersClient{
		context:      context.Background(),
		randomString: internal.RandomString(),
	}
	c.connection = internal.NewConnection()
	c.client = clients.NewUsersClient(c.connection.AuthConfig.TenantID)
	c.client.BaseClient.Authorizer = c.connection.Authorizer

	user := testUsersClient_Create(t, c, models.User{
		AccountEnabled:    internal.Bool(true),
		DisplayName:       internal.String("Test User"),
		MailNickname:      internal.String(fmt.Sprintf("testuser-%s", c.randomString)),
		UserPrincipalName: internal.String(fmt.Sprintf("testuser-%s@%s", c.randomString, c.connection.DomainName)),
		PasswordProfile: &models.UserPasswordProfile{
			Password: internal.String(fmt.Sprintf("IrPa55w0rd%s", c.randomString)),
		},
	})
	testUsersClient_List(t, c)
	testUsersClient_Delete(t, c, *user.ID)
}

func testUsersClient_Create(t *testing.T, c UsersClient, u models.User) (user *models.User) {
	user, status, err := c.client.Create(c.context, u)
	if err != nil {
		t.Errorf("UsersClient.Create(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("UsersClient.Create(): invalid status: %d", status)
	}
	if user == nil {
		t.Fatal("UsersClient.Create(): user was nil")
	}
	if user.ID == nil {
		t.Fatal("UsersClient.Create(): user.ID was nil")
	}
	return
}

func testUsersClient_List(t *testing.T, c UsersClient) (users *[]models.User) {
	users, _, err := c.client.List(c.context, "")
	if err != nil {
		t.Errorf("UsersClient.List(): %v", err)
	}
	if users == nil {
		t.Error("UsersClient.List(): users was nil")
	}
	return
}

func testUsersClient_Delete(t *testing.T, c UsersClient, id string) {
	status, err := c.client.Delete(c.context, id)
	if err != nil {
		t.Errorf("UsersClient.Delete(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Errorf("UsersClient.Delete(): invalid status: %d", status)
	}
}
