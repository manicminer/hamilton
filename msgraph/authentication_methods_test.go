package msgraph_test

import (
	"testing"

	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/msgraph"
	"github.com/manicminer/hamilton/odata"
)

type AuthenticationMethodsClientTest struct {
	connection   *test.Connection
	client       *msgraph.ApplicationsClient
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
		connection: test.NewConnection(auth.MsGraph, auth.TokenVersion2)
		randomString: rs,
	}

	u.client = msgraph.NewUsersClient(u.connection.AuthConfig.TenantID)
	u.client.BaseClient.Authorizer = u.connection.Authorizer

	user := testUsersClient_Create(t, u, msgraph.User{
		AccountEnabled: utils.BoolPtr(true),
		DisplayName: utils.StringPtr("test-user-authenticationmethods"),
		MailNickname:      utils.StringPtr(fmt.Sprintf("test-user-authenticationmethods-%s", c.randomString)),
		UserPrincipalName: utils.StringPtr(fmt.Sprintf("test-user-authenticationmethods-%s@%s", c.randomString, c.connection.DomainName)),
		PasswordProfile: &msgraph.UserPasswordProfile{
			Password: utils.StringPtr(fmt.Sprintf("IrPa55w0rd%s", c.randomString)),
		},
	})

	listAllAuthMethods := testAuthMethods_List(t,c,user)

	
}


func testAuthMethods_List(t *testing.T, c AuthenticationMethodsClientTest, u msgraph.User) (authMethods *[]msgraph.AuthenticationMethod ) {
	authMethods, status, err := c.client.List(c.connection.Context,u,odata.Query{})
	if status < 200 || status >= 300 {
		t.Fatalf("AuthenticationMethodsClientTest.List(): invalid status: %d", status)
	}

	if err != nil {
		t.Fatalf("AuthenticationMethodsClientTest.List(): %v", err)
	}

	if signInLogs == nil {
		t.Fatal("AuthenticationMethodsClientTest.List():logs was nil")
	}
	return
}