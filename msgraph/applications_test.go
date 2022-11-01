package msgraph_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/msgraph"
	"github.com/manicminer/hamilton/odata"
)

func TestApplicationsClient(t *testing.T) {
	c := test.NewTest(t)
	defer c.CancelFunc()

	user := testUsersClient_Create(t, c, msgraph.User{
		AccountEnabled:    utils.BoolPtr(true),
		DisplayName:       utils.StringPtr("test-user-applicationowner"),
		MailNickname:      utils.StringPtr(fmt.Sprintf("test-user-applicationowner-%s", c.RandomString)),
		UserPrincipalName: utils.StringPtr(fmt.Sprintf("test-user-applicationowner-%s@%s", c.RandomString, c.Connections["default"].DomainName)),
		PasswordProfile: &msgraph.UserPasswordProfile{
			Password: utils.StringPtr(fmt.Sprintf("IrPa55w0rd%s", c.RandomString)),
		},
	})

	self := testDirectoryObjectsClient_Get(t, c, c.Claims.ObjectId)

	app := testApplicationsClient_Create(t, c, msgraph.Application{
		DisplayName: utils.StringPtr(fmt.Sprintf("test-application-%s", c.RandomString)),
		GroupMembershipClaims: &[]msgraph.GroupMembershipClaim{
			msgraph.GroupMembershipClaimApplicationGroup,
			msgraph.GroupMembershipClaimDirectoryRole,
			msgraph.GroupMembershipClaimSecurityGroup,
		},
		Owners: &msgraph.Owners{*self},
	})

	testApplicationsClient_Get(t, c, *app.ID)

	app.DisplayName = utils.StringPtr(fmt.Sprintf("test-app-updated-%s", c.RandomString))
	targetObject := []msgraph.ApplicationExtensionTargetObject{
		msgraph.ApplicationExtensionTargetObjectUser,
	}
	newExtension := msgraph.ApplicationExtension{
		DataType:      msgraph.ApplicationExtensionDataTypeString,
		Name:          utils.StringPtr("extName"),
		TargetObjects: &targetObject,
	}
	extensionId := testApplicationsClient_CreateExtension(t, c, newExtension, *app.ID)
	testApplicationsClient_ListExtension(t, c, *app.ID)
	testApplicationsClient_DeleteExtension(t, c, extensionId, *app.ID)

	testApplicationsClient_Update(t, c, *app)

	owners := testApplicationsClient_ListOwners(t, c, *app.ID)
	testApplicationsClient_GetOwner(t, c, *app.ID, (*owners)[0])
	testApplicationsClient_RemoveOwners(t, c, *app.ID, owners)
	app.Owners = &msgraph.Owners{user.DirectoryObject}
	testApplicationsClient_AddOwners(t, c, app)

	pwd := testApplicationsClient_AddPassword(t, c, app)
	testApplicationsClient_RemovePassword(t, c, app, pwd)

	testApplicationsClient_UploadLogo(t, c, app)

	credential := testApplicationsClient_CreateFederatedIdentityCredential(t, c, *app.ID, msgraph.FederatedIdentityCredential{
		Audiences:   &[]string{"api://AzureADTokenExchange"},
		Description: msgraph.NullableString("such testing many pull request"),
		Issuer:      utils.StringPtr("https://token.actions.githubusercontent.com"),
		Name:        utils.StringPtr(fmt.Sprintf("test-credential-%s", c.RandomString)),
		Subject:     utils.StringPtr("repo:manicminer-test/gha-test:pull-request"),
	})
	testApplicationsClient_GetFederatedIdentityCredential(t, c, *app.ID, *credential.ID)

	credential.Description = msgraph.NullableString("")
	testApplicationsClient_UpdateFederatedIdentityCredential(t, c, *app.ID, *credential)
	testApplicationsClient_ListFederatedIdentityCredentials(t, c, *app.ID)
	testApplicationsClient_DeleteFederatedIdentityCredential(t, c, *app.ID, *credential.ID)

	testApplicationsClient_List(t, c)
	testApplicationsClient_Delete(t, c, *app.ID)
	testApplicationsClient_ListDeleted(t, c, *app.ID)
	testApplicationsClient_GetDeleted(t, c, *app.ID)
	testApplicationsClient_RestoreDeleted(t, c, *app.ID)
	testApplicationsClient_Delete(t, c, *app.ID)
	testApplicationsClient_DeletePermanently(t, c, *app.ID)
}

func TestApplicationsClient_groupMembershipClaims(t *testing.T) {
	c := test.NewTest(t)
	defer c.CancelFunc()

	app := testApplicationsClient_Create(t, c, msgraph.Application{
		DisplayName:           utils.StringPtr(fmt.Sprintf("test-application-%s", c.RandomString)),
		GroupMembershipClaims: &[]msgraph.GroupMembershipClaim{"SecurityGroup", "ApplicationGroup"},
	})
	testApplicationsClient_Delete(t, c, *app.ID)
}

func testApplicationsClient_Create(t *testing.T, c *test.Test, a msgraph.Application) (application *msgraph.Application) {
	application, status, err := c.ApplicationsClient.Create(c.Context, a)
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

func testApplicationsClient_Update(t *testing.T, c *test.Test, a msgraph.Application) {
	status, err := c.ApplicationsClient.Update(c.Context, a)
	if err != nil {
		t.Fatalf("ApplicationsClient.Update(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ApplicationsClient.Update(): invalid status: %d", status)
	}
}

func testApplicationsClient_List(t *testing.T, c *test.Test) (applications *[]msgraph.Application) {
	applications, _, err := c.ApplicationsClient.List(c.Context, odata.Query{})
	if err != nil {
		t.Fatalf("ApplicationsClient.List(): %v", err)
	}
	if applications == nil {
		t.Fatal("ApplicationsClient.List(): applications was nil")
	}
	return
}

func testApplicationsClient_Get(t *testing.T, c *test.Test, id string) (application *msgraph.Application) {
	application, status, err := c.ApplicationsClient.Get(c.Context, id, odata.Query{})
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

func testApplicationsClient_CreateExtension(t *testing.T, c *test.Test, applicationExtension msgraph.ApplicationExtension, id string) string {
	extension, status, err := c.ApplicationsClient.CreateExtension(c.Context, applicationExtension, id)
	if err != nil {
		t.Fatalf("ApplicationsClient.CreateExtension(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ApplicationsClient.CreateExtension(): invalid status: %d", status)
	}
	if extension == nil {
		t.Fatal("ApplicationsClient.CreateExtension(): extension was nil")
	}
	if extension.Id == nil {
		t.Fatal("ApplicationsClient.CreateExtension(): extension.Id was nil")
	}
	return *extension.Id
}

func testApplicationsClient_ListExtension(t *testing.T, c *test.Test, id string) {
	extension, status, err := c.ApplicationsClient.ListExtensions(c.Context, id, odata.Query{})
	if err != nil {
		t.Fatalf("ApplicationsClient.ListExtensions(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ApplicationsClient.ListExtensions(): invalid status: %d", status)
	}
	if extension == nil {
		t.Fatal("ApplicationsClient.ListExtensions(): extension was nil")
	}
}

func testApplicationsClient_DeleteExtension(t *testing.T, c *test.Test, extensionId, id string) {
	status, err := c.ApplicationsClient.DeleteExtension(c.Context, id, extensionId)
	if err != nil {
		t.Fatalf("ApplicationsClient.DeleteExtension(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ApplicationsClient.DeleteExtension(): invalid status: %d", status)
	}
}

func testApplicationsClient_GetDeleted(t *testing.T, c *test.Test, id string) (application *msgraph.Application) {
	application, status, err := c.ApplicationsClient.GetDeleted(c.Context, id, odata.Query{})
	if err != nil {
		t.Fatalf("ApplicationsClient.GetDeleted(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ApplicationsClient.GetDeleted(): invalid status: %d", status)
	}
	if application == nil {
		t.Fatal("ApplicationsClient.GetDeleted(): application was nil")
	}
	return
}

func testApplicationsClient_Delete(t *testing.T, c *test.Test, id string) {
	status, err := c.ApplicationsClient.Delete(c.Context, id)
	if err != nil {
		t.Fatalf("ApplicationsClient.Delete(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ApplicationsClient.Delete(): invalid status: %d", status)
	}
}

func testApplicationsClient_DeletePermanently(t *testing.T, c *test.Test, id string) {
	status, err := c.ApplicationsClient.DeletePermanently(c.Context, id)
	if err != nil {
		t.Fatalf("ApplicationsClient.DeletePermanently(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ApplicationsClient.DeletePermanently(): invalid status: %d", status)
	}
}

func testApplicationsClient_RestoreDeleted(t *testing.T, c *test.Test, id string) {
	application, status, err := c.ApplicationsClient.RestoreDeleted(c.Context, id)
	if err != nil {
		t.Fatalf("ApplicationsClient.RestoreDeleted(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ApplicationsClient.RestoreDeleted(): invalid status: %d", status)
	}
	if application == nil {
		t.Fatal("ApplicationsClient.RestoreDeleted(): application was nil")
	}
	if application.ID == nil {
		t.Fatal("ApplicationsClient.RestoreDeleted(): application.ID was nil")
	}
	if *application.ID != id {
		t.Fatal("ApplicationsClient.RestoreDeleted(): application ids do not match")
	}
}

func testApplicationsClient_ListOwners(t *testing.T, c *test.Test, id string) (owners *[]string) {
	owners, status, err := c.ApplicationsClient.ListOwners(c.Context, id)
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

func testApplicationsClient_GetOwner(t *testing.T, c *test.Test, appId string, ownerId string) (owner *string) {
	owner, status, err := c.ApplicationsClient.GetOwner(c.Context, appId, ownerId)
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

func testApplicationsClient_AddOwners(t *testing.T, c *test.Test, a *msgraph.Application) {
	status, err := c.ApplicationsClient.AddOwners(c.Context, a)
	if err != nil {
		t.Fatalf("ApplicationsClient.AddOwners(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ApplicationsClient.AddOwners(): invalid status: %d", status)
	}
}

func testApplicationsClient_RemoveOwners(t *testing.T, c *test.Test, appId string, ownerIds *[]string) {
	status, err := c.ApplicationsClient.RemoveOwners(c.Context, appId, ownerIds)
	if err != nil {
		t.Fatalf("ApplicationsClient.RemoveOwners(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ApplicationsClient.RemoveOwners(): invalid status: %d", status)
	}
}

func testApplicationsClient_AddPassword(t *testing.T, c *test.Test, a *msgraph.Application) *msgraph.PasswordCredential {
	expiry := time.Now().Add(24 * 90 * time.Hour)
	pwd := msgraph.PasswordCredential{
		DisplayName: utils.StringPtr("test password"),
		EndDateTime: &expiry,
	}
	newPwd, status, err := c.ApplicationsClient.AddPassword(c.Context, *a.ID, pwd)
	if err != nil {
		t.Fatalf("ApplicationsClient.AddPassword(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ApplicationsClient.AddPassword(): invalid status: %d", status)
	}
	if newPwd.SecretText == nil || len(*newPwd.SecretText) == 0 {
		t.Fatalf("ApplicationsClient.AddPassword(): nil or empty secretText returned by API")
	}
	if *newPwd.DisplayName != *pwd.DisplayName {
		t.Fatalf("ApplicationsClient.AddPassword(): password names do not match")
	}
	return newPwd
}

func testApplicationsClient_RemovePassword(t *testing.T, c *test.Test, a *msgraph.Application, p *msgraph.PasswordCredential) {
	status, err := c.ApplicationsClient.RemovePassword(c.Context, *a.ID, *p.KeyId)
	if err != nil {
		t.Fatalf("ApplicationsClient.RemovePassword(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ApplicationsClient.RemovePassword(): invalid status: %d", status)
	}
}

func testApplicationsClient_ListDeleted(t *testing.T, c *test.Test, expectedId string) (deletedApps *[]msgraph.Application) {
	deletedApps, status, err := c.ApplicationsClient.ListDeleted(c.Context, odata.Query{
		Filter: fmt.Sprintf("id eq '%s'", expectedId),
		Top:    10,
	})
	if err != nil {
		t.Fatalf("ApplicationsClient.ListDeleted(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ApplicationsClient.ListDeleted(): invalid status: %d", status)
	}
	if deletedApps == nil {
		t.Fatal("ApplicationsClient.ListDeleted(): deletedApps was nil")
	}
	if len(*deletedApps) == 0 {
		t.Fatal("ApplicationsClient.ListDeleted(): expected at least 1 deleted application, was: 0")
	}
	found := false
	for _, app := range *deletedApps {
		if app.ID != nil && *app.ID == expectedId {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("ApplicationsClient.ListDeleted(): expected app ID %q in result", expectedId)
	}
	return
}

func testApplicationsClient_UploadLogo(t *testing.T, c *test.Test, a *msgraph.Application) {
	b, err := os.ReadFile(filepath.Join("..", "internal", "test", "testlogo.png"))
	if err != nil {
		t.Fatalf("reading testlogo.png: %v", err)
	}
	status, err := c.ApplicationsClient.UploadLogo(c.Context, *a.ID, "image/png", b)
	if err != nil {
		t.Fatalf("ApplicationsClient.UploadLogo(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ApplicationsClient.UploadLogo(): invalid status: %d", status)
	}
}

func testApplicationsClient_CreateFederatedIdentityCredential(t *testing.T, c *test.Test, applicationId string, credential msgraph.FederatedIdentityCredential) (newCredential *msgraph.FederatedIdentityCredential) {
	newCredential, status, err := c.ApplicationsClient.CreateFederatedIdentityCredential(c.Context, applicationId, credential)
	if err != nil {
		t.Fatalf("ApplicationsClient.CreateFederatedIdentityCredential(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ApplicationsClient.CreateFederatedIdentityCredential(): invalid status: %d", status)
	}
	if newCredential == nil {
		t.Fatal("ApplicationsClient.CreateFederatedIdentityCredential(): credential was nil")
	}
	if newCredential.ID == nil {
		t.Fatal("ApplicationsClient.CreateFederatedIdentityCredential(): credential.ID was nil")
	}
	return
}

func testApplicationsClient_UpdateFederatedIdentityCredential(t *testing.T, c *test.Test, applicationId string, credential msgraph.FederatedIdentityCredential) {
	status, err := c.ApplicationsClient.UpdateFederatedIdentityCredential(c.Context, applicationId, credential)
	if err != nil {
		t.Fatalf("ApplicationsClient.UpdateFederatedIdentityCredential(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ApplicationsClient.UpdateFederatedIdentityCredential(): invalid status: %d", status)
	}
}

func testApplicationsClient_ListFederatedIdentityCredentials(t *testing.T, c *test.Test, applicationId string) (credentials *[]msgraph.FederatedIdentityCredential) {
	credentials, status, err := c.ApplicationsClient.ListFederatedIdentityCredentials(c.Context, applicationId, odata.Query{})
	if err != nil {
		t.Fatalf("ApplicationsClient.ListFederatedIdentityCredentials(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ApplicationsClient.ListFederatedIdentityCredentials(): invalid status: %d", status)
	}
	if credentials == nil {
		t.Fatal("ApplicationsClient.ListFederatedIdentityCredentials(): credentials was nil")
	}
	return
}

func testApplicationsClient_GetFederatedIdentityCredential(t *testing.T, c *test.Test, applicationId, credentialId string) (credential *msgraph.FederatedIdentityCredential) {
	credential, status, err := c.ApplicationsClient.GetFederatedIdentityCredential(c.Context, applicationId, credentialId, odata.Query{})
	if err != nil {
		t.Fatalf("ApplicationsClient.GetFederatedIdentityCredential(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ApplicationsClient.GetFederatedIdentityCredential(): invalid status: %d", status)
	}
	if credential == nil {
		t.Fatal("ApplicationsClient.GetFederatedIdentityCredential(): credential was nil")
	}
	return
}

func testApplicationsClient_DeleteFederatedIdentityCredential(t *testing.T, c *test.Test, applicationId, credentialId string) {
	status, err := c.ApplicationsClient.DeleteFederatedIdentityCredential(c.Context, applicationId, credentialId)
	if err != nil {
		t.Fatalf("ApplicationsClient.DeleteFederatedIdentityCredential(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ApplicationsClient.DeleteFederatedIdentityCredential(): invalid status: %d", status)
	}
}
