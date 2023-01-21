package msgraph_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/msgraph"
	"github.com/manicminer/hamilton/odata"
)

func TestAuthenticationMethodsClient(t *testing.T) {
	c := test.NewTest(t)
	defer c.CancelFunc()

	user := testUsersClient_Create(t, c, msgraph.User{
		AccountEnabled:    utils.BoolPtr(true),
		DisplayName:       utils.StringPtr("test-user-authenticationmethods"),
		MailNickname:      utils.StringPtr(fmt.Sprintf("test-user-authenticationmethods-%s", c.RandomString)),
		UserPrincipalName: utils.StringPtr(fmt.Sprintf("test-user-authenticationmethods-%s@%s", c.RandomString, c.Connections["default"].DomainName)),
		PasswordProfile: &msgraph.UserPasswordProfile{
			Password: utils.StringPtr(fmt.Sprintf("IrPa55w0rd%s", c.RandomString)),
		},
	})

	_ = testAuthMethods_List(t, c, *user.ID())
	_ = testAuthMethods_ListFido2Methods(t, c, *user.ID())
	_ = testAuthMethods_ListMicrosoftAuthenticatorMethods(t, c, *user.ID())
	_ = testAuthMethods_ListWindowsHelloMethods(t, c, *user.ID())
	tempAccessPass := testAuthMethods_CreateTemporaryAccessPassMethod(t, c, *user.ID())
	_ = testAuthMethods_GetTemporaryAccessPassMethod(t, c, *user.ID(), *tempAccessPass.ID)
	_ = testAuthMethods_ListTemporaryAccessPassMethods(t, c, *user.ID())
	testAuthMethods_DeleteTemporaryAccessPassMethod(t, c, *user.ID(), *tempAccessPass.ID)
	phoneAuthMethod := testAuthMethods_CreatePhoneMethod(t, c, *user.ID())
	_ = testAuthMethods_GetPhoneMethod(t, c, *user.ID(), *phoneAuthMethod.ID)
	_ = testAuthMethods_ListPhoneMethods(t, c, *user.ID())
	phoneAuthMethod.PhoneNumber = utils.StringPtr("+44 07777777778")
	testAuthMethods_UpdatePhoneMethod(t, c, *user.ID(), *phoneAuthMethod)
	testAuthMethods_EnablePhoneSMS(t, c, *user.ID(), *phoneAuthMethod.ID)
	testAuthMethods_DisablePhoneSMS(t, c, *user.ID(), *phoneAuthMethod.ID)
	testAuthMethods_DeletePhoneMethod(t, c, *user.ID(), *phoneAuthMethod.ID)
	emailAuthMethod := testAuthMethods_CreateEmailMethod(t, c, *user.ID())
	_ = testAuthMethods_GetEmailMethod(t, c, *user.ID(), *emailAuthMethod.ID)
	_ = testAuthMethods_ListEmailMethods(t, c, *user.ID())
	emailAuthMethod.EmailAddress = utils.StringPtr("test-user-authenticationmethods@contoso.com")
	testAuthMethods_UpdateEmailMethod(t, c, *user.ID(), *emailAuthMethod)
	testAuthMethods_DeleteEmailMethod(t, c, *user.ID(), *emailAuthMethod.ID)
	_ = testAuthMethods_ListPasswordMethods(t, c, *user.ID())
	testUsersClient_Delete(t, c, *user.ID())
}

func testAuthMethods_List(t *testing.T, c *test.Test, userID string) (authMethods *[]msgraph.AuthenticationMethod) {
	authMethods, status, err := c.AuthenticationMethodsClient.List(c.Context, userID, odata.Query{})
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

func testAuthMethods_ListFido2Methods(t *testing.T, c *test.Test, userID string) (fido2Methods *[]msgraph.Fido2AuthenticationMethod) {
	fido2Methods, status, err := c.AuthenticationMethodsClient.ListFido2Methods(c.Context, userID, odata.Query{})
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

func testAuthMethods_ListMicrosoftAuthenticatorMethods(t *testing.T, c *test.Test, userID string) (msAuthMethods *[]msgraph.MicrosoftAuthenticatorAuthenticationMethod) {
	msAuthMethods, status, err := c.AuthenticationMethodsClient.ListMicrosoftAuthenticatorMethods(c.Context, userID, odata.Query{})
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

func testAuthMethods_ListWindowsHelloMethods(t *testing.T, c *test.Test, userID string) (windowsHelloMethods *[]msgraph.WindowsHelloForBusinessAuthenticationMethod) {
	windowsHelloMethods, status, err := c.AuthenticationMethodsClient.ListWindowsHelloMethods(c.Context, userID, odata.Query{})
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

func testAuthMethods_CreateTemporaryAccessPassMethod(t *testing.T, c *test.Test, userID string) (tempAccessPass *msgraph.TemporaryAccessPassAuthenticationMethod) {
	startPassTime := time.Now().UTC()
	startPassTime.AddDate(0, 0, 1)
	tempPass := msgraph.TemporaryAccessPassAuthenticationMethod{
		StartDateTime:     &startPassTime,
		LifetimeInMinutes: utils.Int32Ptr(60),
		IsUsableOnce:      utils.BoolPtr(true),
	}

	tempAccessPass, status, err := c.AuthenticationMethodsClient.CreateTemporaryAccessPassMethod(c.Context, userID, tempPass)
	if status < 200 || status >= 300 {
		t.Fatalf("AuthenticationMethodsClientTest.CreateTemporaryAccessPassMethod(): invalid status: %d", status)
	}

	if err != nil {
		t.Fatalf("AuthenticationMethodsClientTest.CreateTemporaryAccessPassMethod(): %v", err)
	}

	if tempAccessPass == nil {
		t.Fatal("AuthenticationMethodsClientTest.CreateTemporaryAccessPassMethod():logs was nil")
	}
	return
}

func testAuthMethods_GetTemporaryAccessPassMethod(t *testing.T, c *test.Test, userID, ID string) (tempAccessPass *msgraph.TemporaryAccessPassAuthenticationMethod) {
	tempAccessPass, status, err := c.AuthenticationMethodsClient.GetTemporaryAccessPassMethod(c.Context, userID, ID, odata.Query{})
	if status < 200 || status >= 300 {
		t.Fatalf("AuthenticationMethodsClientTest.GetTemporaryAccessPassMethod(): invalid status: %d", status)
	}

	if err != nil {
		t.Fatalf("AuthenticationMethodsClientTest.GetTemporaryAccessPassMethod(): %v", err)
	}

	if tempAccessPass == nil {
		t.Fatal("AuthenticationMethodsClientTest.GetTemporaryAccessPassMethod():logs was nil")
	}
	return
}

func testAuthMethods_ListTemporaryAccessPassMethods(t *testing.T, c *test.Test, userID string) (tempAccessPasses *[]msgraph.TemporaryAccessPassAuthenticationMethod) {
	tempAccessPasses, status, err := c.AuthenticationMethodsClient.ListTemporaryAccessPassMethods(c.Context, userID, odata.Query{})
	if status < 200 || status >= 300 {
		t.Fatalf("AuthenticationMethodsClientTest.ListTemporaryAccessPassMethod(): invalid status: %d", status)
	}

	if err != nil {
		t.Fatalf("AuthenticationMethodsClientTest.ListTemporaryAccessPassMethod(): %v", err)
	}

	if tempAccessPasses == nil {
		t.Fatal("AuthenticationMethodsClientTest.ListTemporaryAccessPassMethod():logs was nil")
	}
	return
}

func testAuthMethods_DeleteTemporaryAccessPassMethod(t *testing.T, c *test.Test, userID, ID string) {

	status, err := c.AuthenticationMethodsClient.DeleteTemporaryAccessPassMethod(c.Context, userID, ID)
	if status < 200 || status >= 300 {
		t.Fatalf("AuthenticationMethodsClientTest.DeleteTemporaryAccessPassMethod(): invalid status: %d", status)
	}

	if err != nil {
		t.Fatalf("AuthenticationMethodsClientTest.DeleteTemporaryAccessPassMethod(): %v", err)
	}
}

func testAuthMethods_CreatePhoneMethod(t *testing.T, c *test.Test, userID string) (phoneAuthMethod *msgraph.PhoneAuthenticationMethod) {
	phoneType := msgraph.AuthenticationPhoneTypeMobile
	phoneAuth := msgraph.PhoneAuthenticationMethod{
		PhoneNumber: utils.StringPtr("+44 07777777777"),
		PhoneType:   &phoneType,
	}
	phoneAuthMethod, status, err := c.AuthenticationMethodsClient.CreatePhoneMethod(c.Context, userID, phoneAuth)
	if status < 200 || status >= 300 {
		t.Fatalf("AuthenticationMethodsClientTest.CreatePhoneMethod(): invalid status: %d", status)
	}

	if err != nil {
		t.Fatalf("AuthenticationMethodsClientTest.CreatePhoneMethod(): %v", err)
	}

	if phoneAuthMethod == nil {
		t.Fatal("AuthenticationMethodsClientTest.CreatePhoneMethod():logs was nil")
	}
	return
}

func testAuthMethods_GetPhoneMethod(t *testing.T, c *test.Test, userID, ID string) (phoneAuthMethod *msgraph.PhoneAuthenticationMethod) {
	phoneAuthMethod, status, err := c.AuthenticationMethodsClient.GetPhoneMethod(c.Context, userID, ID, odata.Query{})
	if status < 200 || status >= 300 {
		t.Fatalf("AuthenticationMethodsClientTest.GetPhoneMethod(): invalid status: %d", status)
	}

	if err != nil {
		t.Fatalf("AuthenticationMethodsClientTest.GetPhoneMethod(): %v", err)
	}

	if phoneAuthMethod == nil {
		t.Fatal("AuthenticationMethodsClientTest.GetPhoneMethod():logs was nil")
	}
	return
}

func testAuthMethods_ListPhoneMethods(t *testing.T, c *test.Test, userID string) (phoneAuthMethods *[]msgraph.PhoneAuthenticationMethod) {
	phoneAuthMethods, status, err := c.AuthenticationMethodsClient.ListPhoneMethods(c.Context, userID, odata.Query{})
	if status < 200 || status >= 300 {
		t.Fatalf("AuthenticationMethodsClientTest.ListPhoneMethods(): invalid status: %d", status)
	}

	if err != nil {
		t.Fatalf("AuthenticationMethodsClientTest.ListPhoneMethods(): %v", err)
	}

	if phoneAuthMethods == nil {
		t.Fatal("AuthenticationMethodsClientTest.ListPhoneMethods():logs was nil")
	}
	return
}

func testAuthMethods_UpdatePhoneMethod(t *testing.T, c *test.Test, userID string, phone msgraph.PhoneAuthenticationMethod) {
	status, err := c.AuthenticationMethodsClient.UpdatePhoneMethod(c.Context, userID, phone)
	if status < 200 || status >= 300 {
		t.Fatalf("AuthenticationMethodsClientTest.UpdatePhoneMethod(): invalid status: %d", status)
	}

	if err != nil {
		t.Fatalf("AuthenticationMethodsClientTest.UpdatePhoneMethod(): %v", err)
	}
}

func testAuthMethods_EnablePhoneSMS(t *testing.T, c *test.Test, userID, ID string) {
	status, err := c.AuthenticationMethodsClient.EnablePhoneSMS(c.Context, userID, ID)
	if status < 200 || status >= 300 {
		t.Fatalf("AuthenticationMethodsClientTest.EnablePhoneSMS(): invalid status: %d", status)
	}

	if err != nil {
		t.Fatalf("AuthenticationMethodsClientTest.EnablePhoneSMS(): %v", err)
	}
}

func testAuthMethods_DisablePhoneSMS(t *testing.T, c *test.Test, userID, ID string) {
	status, err := c.AuthenticationMethodsClient.DisablePhoneSMS(c.Context, userID, ID)
	if status < 200 || status >= 300 {
		t.Fatalf("AuthenticationMethodsClientTest.DisablePhoneSMS(): invalid status: %d", status)
	}

	if err != nil {
		t.Fatalf("AuthenticationMethodsClientTest.DisablePhoneSMS(): %v", err)
	}
}

func testAuthMethods_DeletePhoneMethod(t *testing.T, c *test.Test, userID, ID string) {
	status, err := c.AuthenticationMethodsClient.DeletePhoneMethod(c.Context, userID, ID)
	if status < 200 || status >= 300 {
		t.Fatalf("AuthenticationMethodsClientTest.DeletePhoneMethod(): invalid status: %d", status)
	}

	if err != nil {
		t.Fatalf("AuthenticationMethodsClientTest.DeletePhoneMethod(): %v", err)
	}
}

func testAuthMethods_CreateEmailMethod(t *testing.T, c *test.Test, userID string) (emailMethod *msgraph.EmailAuthenticationMethod) {
	email := msgraph.EmailAuthenticationMethod{
		EmailAddress: utils.StringPtr("test-user-authenticationmethods@testdomain.com"),
	}
	emailMethod, status, err := c.AuthenticationMethodsClient.CreateEmailMethod(c.Context, userID, email)
	if status < 200 || status >= 300 {
		t.Fatalf("AuthenticationMethodsClientTest.CreateEmailMethod(): invalid status: %d", status)
	}

	if err != nil {
		t.Fatalf("AuthenticationMethodsClientTest.CreateEmailMethod(): %v", err)
	}

	if emailMethod == nil {
		t.Fatal("AuthenticationMethodsClientTest.CreateEmailMethod():logs was nil")
	}
	return
}

func testAuthMethods_GetEmailMethod(t *testing.T, c *test.Test, userID, ID string) (emailMethod *msgraph.EmailAuthenticationMethod) {
	emailMethod, status, err := c.AuthenticationMethodsClient.GetEmailMethod(c.Context, userID, ID, odata.Query{})
	if status < 200 || status >= 300 {
		t.Fatalf("AuthenticationMethodsClientTest.GetEmailMethod(): invalid status: %d", status)
	}

	if err != nil {
		t.Fatalf("AuthenticationMethodsClientTest.GetEmailMethod(): %v", err)
	}

	if emailMethod == nil {
		t.Fatal("AuthenticationMethodsClientTest.GetEmailMethod():logs was nil")
	}
	return
}

func testAuthMethods_ListEmailMethods(t *testing.T, c *test.Test, userID string) (emailMethods *[]msgraph.EmailAuthenticationMethod) {
	emailMethods, status, err := c.AuthenticationMethodsClient.ListEmailMethods(c.Context, userID, odata.Query{})
	if status < 200 || status >= 300 {
		t.Fatalf("AuthenticationMethodsClientTest.ListEmailMethods(): invalid status: %d", status)
	}

	if err != nil {
		t.Fatalf("AuthenticationMethodsClientTest.ListEmailMethods(): %v", err)
	}

	if emailMethods == nil {
		t.Fatal("AuthenticationMethodsClientTest.ListEmailMethods():logs was nil")
	}
	return
}

func testAuthMethods_UpdateEmailMethod(t *testing.T, c *test.Test, userID string, email msgraph.EmailAuthenticationMethod) {
	status, err := c.AuthenticationMethodsClient.UpdateEmailMethod(c.Context, userID, email)
	if status < 200 || status >= 300 {
		t.Fatalf("AuthenticationMethodsClientTest.UpdateEmailMethod(): invalid status: %d", status)
	}

	if err != nil {
		t.Fatalf("AuthenticationMethodsClientTest.UpdateEmailMethod(): %v", err)
	}
}

func testAuthMethods_DeleteEmailMethod(t *testing.T, c *test.Test, userID, ID string) {
	status, err := c.AuthenticationMethodsClient.DeleteEmailMethod(c.Context, userID, ID)
	if status < 200 || status >= 300 {
		t.Fatalf("AuthenticationMethodsClientTest.DeleteEmailMethod(): invalid status: %d", status)
	}

	if err != nil {
		t.Fatalf("AuthenticationMethodsClientTest.DeleteEmailMethod(): %v", err)
	}
}

func testAuthMethods_ListPasswordMethods(t *testing.T, c *test.Test, userID string) (passwordMethods *[]msgraph.
	PasswordAuthenticationMethod) {
	passwordMethods, status, err := c.AuthenticationMethodsClient.ListPasswordMethods(c.Context, userID, odata.Query{})
	if status < 200 || status >= 300 {
		t.Fatalf("AuthenticationMethodsClientTest.ListPasswordMethods(): invalid status: %d", status)
	}

	if err != nil {
		t.Fatalf("AuthenticationMethodsClientTest.ListPasswordMethods(): %v", err)
	}

	if passwordMethods == nil {
		t.Fatal("AuthenticationMethodsClientTest.ListPasswordMethods():logs was nil")
	}
	return
}
