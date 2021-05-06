package msgraph_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/msgraph"
)

func TestApplicationAppRoleModel(t *testing.T) {
	a := msgraph.Application{
		DisplayName: utils.StringPtr(fmt.Sprintf("test-application-%s", test.RandomString())),
		AppRoles:    &[]msgraph.AppRole{},
	}
	testAppRoleId, _ := uuid.GenerateUUID()
	appRole := &msgraph.AppRole{
		ID:        utils.StringPtr(testAppRoleId),
		IsEnabled: utils.BoolPtr(true),
	}
	a = testApplicationAppRole_append(t, a, appRole)
	a = testApplicationAppRole_update(t, a, appRole)
	testApplicationAppRole_remove(t, a, appRole)

}

func testApplicationAppRole_append(t *testing.T, application msgraph.Application, appRole *msgraph.AppRole) msgraph.Application {
	if err := application.AppendAppRole(*appRole); err != nil {
		t.Fatalf("Application.AppendAppRole(): %t", err)
	}
	if application.AppRoles == nil {
		t.Fatal("Application.AppendAppRole(): application.AppRoles was nil")
	}
	if len(*application.AppRoles) != 1 {
		t.Fatal("Application.AppendAppRole(): application.AppRoles did not contain 1 app role")
	}
	if (*application.AppRoles)[0].ID != appRole.ID {
		t.Fatal("Application.AppendAppRole(): application.AppRoles does not contain new app role")
	}
	return application
}

func testApplicationAppRole_update(t *testing.T, application msgraph.Application, appRole *msgraph.AppRole) msgraph.Application {
	appRole.IsEnabled = utils.BoolPtr(false)
	if err := application.UpdateAppRole(*appRole); err != nil {
		t.Fatalf("Application.UpdateAppRole(): %t", err)
	}
	if application.AppRoles == nil {
		t.Fatal("Application.UpdateAppRole(): application.AppRoles was nil")
	}
	if len(*application.AppRoles) != 1 {
		t.Fatal("Application.UpdateAppRole(): application.AppRoles did not contain 1 app role")
	}
	if (*application.AppRoles)[0].IsEnabled != appRole.IsEnabled {
		t.Fatal("Application.UpdateAppRole(): application.AppRoles does not contain updated role")
	}
	return application
}

func testApplicationAppRole_remove(t *testing.T, application msgraph.Application, appRole *msgraph.AppRole) {
	if err := application.RemoveAppRole(*appRole); err != nil {
		t.Fatalf("Application.RemoveAppRole(): %t", err)
	}
	if application.AppRoles == nil {
		t.Fatal("Application.RemoveAppRole(): application.AppRoles was nil")
	}
	if len(*application.AppRoles) != 0 {
		t.Fatal("Application.RemoveAppRole(): application.AppRoles is not empty")
	}
}

func TestServicePrincipalAppRoleModel(t *testing.T) {
	s := msgraph.ServicePrincipal{
		DisplayName: utils.StringPtr(fmt.Sprintf("test-servicePrincipal-%s", test.RandomString())),
		AppRoles:    &[]msgraph.AppRole{},
	}
	testAppRoleId, _ := uuid.GenerateUUID()
	appRole := &msgraph.AppRole{
		ID:        utils.StringPtr(testAppRoleId),
		IsEnabled: utils.BoolPtr(true),
	}
	s = testServicePrincipalAppRole_append(t, s, appRole)
	s = testServicePrincipalAppRole_update(t, s, appRole)
	testServicePrincipalAppRole_remove(t, s, appRole)

}

func testServicePrincipalAppRole_append(t *testing.T, servicePrincipal msgraph.ServicePrincipal, appRole *msgraph.AppRole) msgraph.ServicePrincipal {
	if err := servicePrincipal.AppendAppRole(*appRole); err != nil {
		t.Fatalf("ServicePrincipal.AppendAppRole(): %t", err)
	}
	if servicePrincipal.AppRoles == nil {
		t.Fatal("ServicePrincipal.AppendAppRole(): servicePrincipal.AppRoles was nil")
	}
	if len(*servicePrincipal.AppRoles) != 1 {
		t.Fatal("ServicePrincipal.AppendAppRole(): servicePrincipal.AppRoles did not contain 1 app role")
	}
	if (*servicePrincipal.AppRoles)[0].ID != appRole.ID {
		t.Fatal("ServicePrincipal.AppendAppRole(): servicePrincipal.AppRoles does not contain new app role")
	}
	return servicePrincipal
}

func testServicePrincipalAppRole_update(t *testing.T, servicePrincipal msgraph.ServicePrincipal, appRole *msgraph.AppRole) msgraph.ServicePrincipal {
	appRole.IsEnabled = utils.BoolPtr(false)
	if err := servicePrincipal.UpdateAppRole(*appRole); err != nil {
		t.Fatalf("ServicePrincipal.UpdateAppRole(): %t", err)
	}
	if servicePrincipal.AppRoles == nil {
		t.Fatal("ServicePrincipal.UpdateAppRole(): servicePrincipal.AppRoles was nil")
	}
	if len(*servicePrincipal.AppRoles) != 1 {
		t.Fatal("ServicePrincipal.UpdateAppRole(): servicePrincipal.AppRoles did not contain 1 app role")
	}
	if (*servicePrincipal.AppRoles)[0].IsEnabled != appRole.IsEnabled {
		t.Fatal("ServicePrincipal.UpdateAppRole(): servicePrincipal.AppRoles does not contain updated role")
	}
	return servicePrincipal
}

func testServicePrincipalAppRole_remove(t *testing.T, servicePrincipal msgraph.ServicePrincipal, appRole *msgraph.AppRole) {
	if err := servicePrincipal.RemoveAppRole(*appRole); err != nil {
		t.Fatalf("ServicePrincipal.RemoveAppRole(): %t", err)
	}
	if servicePrincipal.AppRoles == nil {
		t.Fatal("ServicePrincipal.RemoveAppRole(): servicePrincipal.AppRoles was nil")
	}
	if len(*servicePrincipal.AppRoles) != 0 {
		t.Fatal("ServicePrincipal.RemoveAppRole(): servicePrincipal.AppRoles is not empty")
	}
}
