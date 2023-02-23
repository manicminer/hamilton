package msgraph_test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/msgraph"
)

func TestAdministrativeUnitsClient(t *testing.T) {
	c := test.NewTest(t)
	defer c.CancelFunc()

	newAdministrativeUnit := msgraph.AdministrativeUnit{
		DisplayName: utils.StringPtr("test-administrativeUnit"),
		Description: msgraph.NullableString("test-administrativeUnit-description"),
	}
	administrativeUnit := testAdministrativeUnitsClient_Create(t, c, newAdministrativeUnit)
	testAdministrativeUnitsClient_Get(t, c, *administrativeUnit.ID)

	administrativeUnit.DisplayName = utils.StringPtr(fmt.Sprintf("test-updated-administrativeUnit-%s", c.RandomString))
	administrativeUnit.Description = msgraph.NullableString("")
	testAdministrativeUnitsClient_Update(t, c, *administrativeUnit)

	user := testUsersClient_Create(t, c, msgraph.User{
		AccountEnabled:    utils.BoolPtr(true),
		DisplayName:       utils.StringPtr("test-user"),
		MailNickname:      utils.StringPtr(fmt.Sprintf("test-user-%s", c.RandomString)),
		UserPrincipalName: utils.StringPtr(fmt.Sprintf("test-user-%s@%s", c.RandomString, c.Connections["default"].DomainName)),
		PasswordProfile: &msgraph.UserPasswordProfile{
			Password: utils.StringPtr(fmt.Sprintf("IrPa55w0rd%s", c.RandomString)),
		},
	})

	testAdministrativeUnitsClient_AddMembers(t, c, *administrativeUnit.ID, &msgraph.Members{user.DirectoryObject})
	testAdministrativeUnitsClient_ListMembers(t, c, *administrativeUnit.ID)
	testAdministrativeUnitsClient_GetMember(t, c, *administrativeUnit.ID, *user.ID())
	testAdministrativeUnitsClient_RemoveMembers(t, c, *administrativeUnit.ID, &([]string{*user.ID()}))

	self := testDirectoryObjectsClient_Get(t, c, c.Claims.ObjectId)
	group := msgraph.Group{
		DisplayName:     utils.StringPtr("test-group"),
		MailEnabled:     utils.BoolPtr(false),
		MailNickname:    utils.StringPtr(fmt.Sprintf("test-group-%s", c.RandomString)),
		SecurityEnabled: utils.BoolPtr(true),
		Owners:          &msgraph.Owners{*self},
		Members:         &msgraph.Members{*self},
	}
	createdGroup := testAdministrativeUnitsClient_CreateGroup(t, c, *administrativeUnit.ID, &group)
	testAdministrativeUnitsClient_RemoveMembers(t, c, *administrativeUnit.ID, &([]string{*createdGroup.ID()}))
	testGroup_Delete(t, c, createdGroup)

	directoryRoleTemplates := testDirectoryRoleTemplatesClient_List(t, c)
	var helpdeskAdministratorRoleId string
	for _, template := range *directoryRoleTemplates {
		if strings.EqualFold(*template.DisplayName, "Helpdesk administrator") {
			helpdeskAdministratorRoleId = *template.ID
		}
	}
	testDirectoryRolesClient_Activate(t, c, helpdeskAdministratorRoleId)
	directoryRole := testDirectoryRolesClient_GetByTemplateId(t, c, helpdeskAdministratorRoleId)

	membership := testAdministrativeUnitsClient_AddScopedRoleMember(t, c, *administrativeUnit.ID, msgraph.ScopedRoleMembership{
		RoleId:         directoryRole.ID(),
		RoleMemberInfo: &msgraph.Identity{Id: user.ID()},
	})
	testAdministrativeUnitsClient_ListScopedRoleMembers(t, c, *administrativeUnit.ID)
	testAdministrativeUnitsClient_GetRoleScopedMember(t, c, *administrativeUnit.ID, *membership.Id)
	testAdministrativeUnitsClient_RemoveScopedRoleMember(t, c, *administrativeUnit.ID, *membership.Id)

	testAdministrativeUnitsClient_List(t, c)
	testAdministrativeUnitsClient_Delete(t, c, *administrativeUnit.ID)
	testUsersClient_Delete(t, c, *user.ID())
}

func testAdministrativeUnitsClient_Create(t *testing.T, c *test.Test, g msgraph.AdministrativeUnit) (administrativeUnit *msgraph.AdministrativeUnit) {
	administrativeUnit, status, err := c.AdministrativeUnitsClient.Create(c.Context, g)
	if err != nil {
		t.Fatalf("AdministrativeUnitsClient.Create(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AdministrativeUnitsClient.Create(): invalid status: %d", status)
	}
	if administrativeUnit == nil {
		t.Fatal("AdministrativeUnitsClient.Create(): administrativeUnit was nil")
	}
	if administrativeUnit.ID == nil {
		t.Fatal("AdministrativeUnitsClient.Create(): administrativeUnit.ID was nil")
	}
	return
}

func testAdministrativeUnitsClient_Update(t *testing.T, c *test.Test, g msgraph.AdministrativeUnit) {
	status, err := c.AdministrativeUnitsClient.Update(c.Context, g)
	if err != nil {
		t.Fatalf("AdministrativeUnitsClient.Update(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AdministrativeUnitsClient.Update(): invalid status: %d", status)
	}
}

func testAdministrativeUnitsClient_List(t *testing.T, c *test.Test) (administrativeUnits *[]msgraph.AdministrativeUnit) {
	administrativeUnits, _, err := c.AdministrativeUnitsClient.List(c.Context, odata.Query{Top: 10})
	if err != nil {
		t.Fatalf("AdministrativeUnitsClient.List(): %v", err)
	}
	if administrativeUnits == nil {
		t.Fatal("AdministrativeUnitsClient.List(): administrativeUnits was nil")
	}
	return
}

func testAdministrativeUnitsClient_Get(t *testing.T, c *test.Test, id string) (administrativeUnit *msgraph.AdministrativeUnit) {
	administrativeUnit, status, err := c.AdministrativeUnitsClient.Get(c.Context, id, odata.Query{})
	if err != nil {
		t.Fatalf("AdministrativeUnitsClient.Get(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AdministrativeUnitsClient.Get(): invalid status: %d", status)
	}
	if administrativeUnit == nil {
		t.Fatal("AdministrativeUnitsClient.Get(): administrativeUnit was nil")
	}
	return
}

func testAdministrativeUnitsClient_Delete(t *testing.T, c *test.Test, id string) {
	status, err := c.AdministrativeUnitsClient.Delete(c.Context, id)
	if err != nil {
		t.Fatalf("AdministrativeUnitsClient.Delete(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AdministrativeUnitsClient.Delete(): invalid status: %d", status)
	}
}

func testAdministrativeUnitsClient_ListMembers(t *testing.T, c *test.Test, administrativeUnitId string) (members *[]string) {
	members, status, err := c.AdministrativeUnitsClient.ListMembers(c.Context, administrativeUnitId)
	if err != nil {
		t.Fatalf("AdministrativeUnitsClient.ListMembers(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AdministrativeUnitsClient.ListMembers(): invalid status: %d", status)
	}
	if members == nil {
		t.Fatal("AdministrativeUnitsClient.ListMembers(): members was nil")
	}
	if len(*members) == 0 {
		t.Fatal("AdministrativeUnitsClient.ListMembers(): members was empty")
	}
	return
}

func testAdministrativeUnitsClient_GetMember(t *testing.T, c *test.Test, administrativeUnitId string, memberId string) (member *string) {
	member, status, err := c.AdministrativeUnitsClient.GetMember(c.Context, administrativeUnitId, memberId)
	if err != nil {
		t.Fatalf("AdministrativeUnitsClient.GetMember(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AdministrativeUnitsClient.GetMember(): invalid status: %d", status)
	}
	if member == nil {
		t.Fatal("AdministrativeUnitsClient.GetMember(): member was nil")
	}
	return
}

func testAdministrativeUnitsClient_AddMembers(t *testing.T, c *test.Test, administrativeUnitId string, members *msgraph.Members) {
	status, err := c.AdministrativeUnitsClient.AddMembers(c.Context, administrativeUnitId, members)
	if err != nil {
		t.Fatalf("AdministrativeUnitsClient.AddMembers(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AdministrativeUnitsClient.AddMembers(): invalid status: %d", status)
	}
}

func testAdministrativeUnitsClient_RemoveMembers(t *testing.T, c *test.Test, administrativeUnitId string, memberIds *[]string) {
	status, err := c.AdministrativeUnitsClient.RemoveMembers(c.Context, administrativeUnitId, memberIds)
	if err != nil {
		t.Fatalf("AdministrativeUnitsClient.RemoveMembers(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AdministrativeUnitsClient.RemoveMembers(): invalid status: %d", status)
	}
}

func testAdministrativeUnitsClient_CreateGroup(t *testing.T, c *test.Test, administrativeUnitId string, g *msgraph.Group) (group *msgraph.Group) {
	group, status, err := c.AdministrativeUnitsClient.CreateGroup(c.Context, administrativeUnitId, g)
	if err != nil {
		t.Fatalf("AdministrativeUnitsClient.CreateGroup(): %v", err)
	}
	if status != http.StatusCreated {
		t.Fatalf("AdministrativeUnitsClient.CreateGroup(): invalid status: %d", status)
	}
	return group
}

func testAdministrativeUnitsClient_ListScopedRoleMembers(t *testing.T, c *test.Test, administrativeUnitId string) (memberships *[]msgraph.ScopedRoleMembership) {
	memberships, status, err := c.AdministrativeUnitsClient.ListScopedRoleMembers(c.Context, administrativeUnitId, odata.Query{})
	if err != nil {
		t.Fatalf("AdministrativeUnitsClient.ListScopedRoleMembers(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AdministrativeUnitsClient.ListScopedRoleMembers(): invalid status: %d", status)
	}
	if memberships == nil {
		t.Fatal("AdministrativeUnitsClient.ListScopedRoleMembers(): members was nil")
	}
	if len(*memberships) == 0 {
		t.Fatal("AdministrativeUnitsClient.ListScopedRoleMembers(): members was empty")
	}
	return
}

func testAdministrativeUnitsClient_GetRoleScopedMember(t *testing.T, c *test.Test, administrativeUnitId string, memberId string) (membership *msgraph.ScopedRoleMembership) {
	member, status, err := c.AdministrativeUnitsClient.GetScopedRoleMember(c.Context, administrativeUnitId, memberId, odata.Query{})
	if err != nil {
		t.Fatalf("AdministrativeUnitsClient.GetScopedRoleMember(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AdministrativeUnitsClient.GetScopedRoleMember(): invalid status: %d", status)
	}
	if member == nil {
		t.Fatal("AdministrativeUnitsClient.GetScopedRoleMember(): member was nil")
	}
	return
}

func testAdministrativeUnitsClient_AddScopedRoleMember(t *testing.T, c *test.Test, administrativeUnitId string, member msgraph.ScopedRoleMembership) (membership *msgraph.ScopedRoleMembership) {
	membership, status, err := c.AdministrativeUnitsClient.AddScopedRoleMember(c.Context, administrativeUnitId, member)
	if err != nil {
		t.Fatalf("AdministrativeUnitsClient.AddScopedRoleMember(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AdministrativeUnitsClient.AddScopedRoleMember(): invalid status: %d", status)
	}
	return
}

func testAdministrativeUnitsClient_RemoveScopedRoleMember(t *testing.T, c *test.Test, administrativeUnitId string, membershipId string) {
	status, err := c.AdministrativeUnitsClient.RemoveScopedRoleMembers(c.Context, administrativeUnitId, membershipId)
	if err != nil {
		t.Fatalf("AdministrativeUnitsClient.RemoveScopedRoleMembers(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AdministrativeUnitsClient.RemoveScopedRoleMembers(): invalid status: %d", status)
	}
}
