package msgraph_test

import (
	"fmt"
	"testing"

	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/msgraph"
	"github.com/manicminer/hamilton/odata"
)

type AuthenticationMethodsClientTest struct {
	connection   *test.Connection
	client       *msgraph.AuthenticationMethodsClient
	randomString string
}

func TestAuthenticationMethodsClient(t *testing.T) {
	rs := test.RandomString()
	c := AuthenticationMethodsClientTest{
		connection:   test.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: rs,
	}

	c.client = msgraph.NewAuthenticationMethodsClient(c.connection.AuthConfig.TenantID)
	c.client.BaseClient.Authorizer = c.connection.Authorizer

	u := UsersClientTest{
		connection:   test.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: rs,
	}

	u.client = msgraph.NewUsersClient(u.connection.AuthConfig.TenantID)
	u.client.BaseClient.Authorizer = u.connection.Authorizer

	user := testUsersClient_Create(t, u, msgraph.User{
		AccountEnabled:    utils.BoolPtr(true),
		DisplayName:       utils.StringPtr("test-user-authenticationmethods"),
		MailNickname:      utils.StringPtr(fmt.Sprintf("test-user-authenticationmethods-%s", c.randomString)),
		UserPrincipalName: utils.StringPtr(fmt.Sprintf("test-user-authenticationmethods-%s@%s", c.randomString, c.connection.DomainName)),
		PasswordProfile: &msgraph.UserPasswordProfile{
			Password: utils.StringPtr(fmt.Sprintf("IrPa55w0rd%s", c.randomString)),
		},
	})

	_ = testAuthMethods_List(t, c, *user.ID)
	_ = testAuthMethods_ListFido2Methods(t, c, *user.ID)
	_ = testAuthMethods_ListMicrosoftAuthenticatorMethods(t, c, *user.ID)
	_ = testAuthMethods_ListWindowsHelloMethods(t, c, *user.ID)

}

func testAuthMethods_List(t *testing.T, c AuthenticationMethodsClientTest, userID string) (authMethods *[]msgraph.AuthenticationMethod) {
	authMethods, status, err := c.client.List(c.connection.Context, userID, odata.Query{})
	if status < 200 || status >= 300 {
		t.Fatalf("AuthenticationMethodsClientTest.List(): invalid status: %d", status)
	}

	if err != nil {
		t.Fatalf("AuthenticationMethodsClientTest.List(): %v", err)
	}

	if authMethods == nil {
		t.Fatal("AuthenticationMethodsClientTest.List():logs was nil")
	}
	return
}

func testAuthMethods_ListFido2Methods(t *testing.T, c AuthenticationMethodsClientTest, userID string) (fido2Methods *[]msgraph.Fido2AuthenticationMethod) {
	fido2Methods, status, err := c.client.ListFido2Methods(c.connection.Context, userID, odata.Query{})
	if status < 200 || status >= 300 {
		t.Fatalf("AuthenticationMethodsClientTest.ListFido2Methods(): invalid status: %d", status)
	}

	if err != nil {
		t.Fatalf("AuthenticationMethodsClientTest.ListFido2Methods(): %v", err)
	}

	if fido2Methods == nil {
		t.Fatal("AuthenticationMethodsClientTest.ListFido2Methods():logs was nil")
	}
	return
}

func testAuthMethods_ListMicrosoftAuthenticatorMethods(t *testing.T, c AuthenticationMethodsClientTest, userID string) (msAuthMethods *[]msgraph.MicrosoftAuthenticatorAuthenticationMethod) {
	msAuthMethods, status, err := c.client.ListMicrosoftAuthenticatorMethods(c.connection.Context, userID, odata.Query{})
	if status < 200 || status >= 300 {
		t.Fatalf("AuthenticationMethodsClientTest.ListMicrosoftAuthenticatorMethods(): invalid status: %d", status)
	}

	if err != nil {
		t.Fatalf("AuthenticationMethodsClientTest.ListMicrosoftAuthenticatorMethods(): %v", err)
	}

	if msAuthMethods == nil {
		t.Fatal("AuthenticationMethodsClientTest.ListMicrosoftAuthenticatorMethods():logs was nil")
	}
	return
}

func testAuthMethods_ListWindowsHelloMethods(t *testing.T, c AuthenticationMethodsClientTest, userID string) (windowsHelloMethods *[]msgraph.WindowsHelloForBusinessAuthenticationMethod) {
	windowsHelloMethods, status, err := c.client.ListWindowsHelloMethods(c.connection.Context, userID, odata.Query{})
	if status < 200 || status >= 300 {
		t.Fatalf("AuthenticationMethodsClientTest.ListWindowsHelloMethods(): invalid status: %d", status)
	}

	if err != nil {
		t.Fatalf("AuthenticationMethodsClientTest.ListWindowsHelloMethods(): %v", err)
	}

	if windowsHelloMethods == nil {
		t.Fatal("AuthenticationMethodsClientTest.ListWindowsHelloMethods():logs was nil")
	}
	return
}
