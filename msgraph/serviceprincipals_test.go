package msgraph_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/hashicorp/go-uuid"
	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/msgraph"
)

func TestServicePrincipalsClient(t *testing.T) {
	c := test.NewTest(t)
	defer c.CancelFunc()

	app := testApplicationsClient_Create(t, c, msgraph.Application{
		DisplayName: utils.StringPtr(fmt.Sprintf("test-serviceprincipal-%s", c.RandomString)),
	})

	sp := testServicePrincipalsClient_Create(t, c, msgraph.ServicePrincipal{
		AccountEnabled: utils.BoolPtr(true),
		AppId:          app.AppId,
		DisplayName:    app.DisplayName,
	})

	appChild := testApplicationsClient_Create(t, c, msgraph.Application{
		DisplayName: utils.StringPtr(fmt.Sprintf("test-serviceprincipal-child%s", c.RandomString)),
	})
	spChild := testServicePrincipalsClient_Create(t, c, msgraph.ServicePrincipal{
		AccountEnabled: utils.BoolPtr(true),
		AppId:          appChild.AppId,
		DisplayName:    appChild.DisplayName,
	})

	spChild.Owners = &msgraph.Owners{sp.DirectoryObject}
	testServicePrincipalsClient_AddOwners(t, c, spChild)
	testServicePrincipalsClient_ListOwners(t, c, *spChild.ID(), []string{*sp.ID()})
	testServicePrincipalsClient_GetOwner(t, c, *spChild.ID(), *sp.ID())
	testServicePrincipalsClient_Get(t, c, *sp.ID())
	sp.Tags = &([]string{"TestTag"})
	testServicePrincipalsClient_Update(t, c, *sp)
	pwd := testServicePrincipalsClient_AddPassword(t, c, sp)
	testServicePrincipalsClient_RemovePassword(t, c, sp, pwd)

	tsc := testServicePrincipalsClient_AddTokenSigningCertificate(t, c, sp)
	testServicePrincipalsClient_SetPreferredTokenSigningKeyThumbprint(t, c, sp, *tsc.Thumbprint)

	testServicePrincipalsClient_List(t, c, odata.Query{})

	claimsMappingPolicy := testClaimsMappingPolicyClient_Create(t, c, msgraph.ClaimsMappingPolicy{
		DisplayName: utils.StringPtr(fmt.Sprintf("test-claims-mapping-policy-%s", c.RandomString)),
		Definition: utils.ArrayStringPtr(
			[]string{
				"{\"ClaimsMappingPolicy\":{\"Version\":1,\"IncludeBasicClaimSet\":\"true\",\"ClaimsSchema\": [{\"Source\":\"user\",\"ID\":\"employeeid\",\"SamlClaimType\":\"http://schemas.xmlsoap.org/ws/2005/05/identity/claims/name\",\"JwtClaimType\":\"name\"},{\"Source\":\"company\",\"ID\":\"tenantcountry\",\"SamlClaimType\":\"http://schemas.xmlsoap.org/ws/2005/05/identity/claims/country\",\"JwtClaimType\":\"country\"}]}}",
			},
		),
	})

	sp.ClaimsMappingPolicies = &[]msgraph.ClaimsMappingPolicy{*claimsMappingPolicy}

	testServicePrincipalsClient_AssignClaimsMappingPolicy(t, c, sp)
	// ListClaimsMappingPolicy is called within RemoveClaimsMappingPolicy
	testServicePrincipalsClient_RemoveClaimsMappingPolicy(t, c, sp, []string{*claimsMappingPolicy.ID()})
	// A Second call tests that a remove call on an empty assignment list returns ok
	testServicePrincipalsClient_RemoveClaimsMappingPolicy(t, c, sp, []string{*claimsMappingPolicy.ID()})
	testClaimsMappingPolicyClient_Delete(t, c, *claimsMappingPolicy.ID())

	tokenIssuancePolicy := testTokenIssuancePolicyClient_Create(t, c, msgraph.TokenIssuancePolicy{
		DisplayName: utils.StringPtr(fmt.Sprintf("test-token-issuance-policy-%s", c.RandomString)),
		Definition: utils.ArrayStringPtr(
			[]string{
				"{\"TokenIssuancePolicy\":{\"Version\":1,\"SigningAlgorithm\":\"http://www.w3.org/2001/04/xmldsig-more#rsa-sha256\",\"TokenResponseSigningPolicy\":\"ResponseAndToken\",\"SamlTokenVersion\":\"2.0\",\"EmitSamlNameFormat\":false}}",
			},
		),
	})

	sp.TokenIssuancePolicies = &[]msgraph.TokenIssuancePolicy{*tokenIssuancePolicy}

	testServicePrincipalsClient_AssignTokenIssuancePolicy(t, c, sp)
	// ListTokenIssuancePolicy is called within RemoveTokenIssuancePolicy
	testServicePrincipalsClient_RemoveTokenIssuancePolicy(t, c, sp, []string{*tokenIssuancePolicy.Id})
	// A Second call tests that a remove call on an empty assignment list returns ok
	testServicePrincipalsClient_RemoveTokenIssuancePolicy(t, c, sp, []string{*tokenIssuancePolicy.Id})
	testTokenIssuancePolicyClient_Delete(t, c, *tokenIssuancePolicy.Id)

	newGroupParent := msgraph.Group{
		DisplayName:     utils.StringPtr("test-group-servicePrincipal-parent"),
		MailEnabled:     utils.BoolPtr(false),
		MailNickname:    utils.StringPtr(fmt.Sprintf("test-group-parent-%s", c.RandomString)),
		SecurityEnabled: utils.BoolPtr(true),
	}
	newGroupChild := msgraph.Group{
		DisplayName:     utils.StringPtr("test-group-servicePrincipal-child"),
		MailEnabled:     utils.BoolPtr(false),
		MailNickname:    utils.StringPtr(fmt.Sprintf("test-group-child-%s", c.RandomString)),
		SecurityEnabled: utils.BoolPtr(true),
	}

	groupParent := testGroupsClient_Create(t, c, newGroupParent)
	groupChild := testGroupsClient_Create(t, c, newGroupChild)
	groupParent.Members = &msgraph.Members{groupChild.DirectoryObject}
	testGroupsClient_AddMembers(t, c, groupParent)
	groupChild.Members = &msgraph.Members{sp.DirectoryObject}
	testGroupsClient_AddMembers(t, c, groupChild)

	testServicePrincipalsClient_ListGroupMemberships(t, c, *sp.ID())
	testServicePrincipalsClient_ListOwnedObjects(t, c, *sp.ID())

	testServicePrincipalsClient_RemoveOwners(t, c, *spChild.ID(), []string{*sp.ID()})
	testGroupsClient_Delete(t, c, *groupParent.ID())
	testGroupsClient_Delete(t, c, *groupChild.ID())

	testServicePrincipalsClient_Delete(t, c, *sp.ID())
	testServicePrincipalsClient_Delete(t, c, *spChild.ID())

	testApplicationsClient_Delete(t, c, *app.ID())
	testApplicationsClient_Delete(t, c, *appChild.ID())
}

func TestServicePrincipalsClient_AppRoleAssignments(t *testing.T) {
	c := test.NewTest(t)
	defer c.CancelFunc()

	// pre-generate uuid for a test app role
	testResourceAppRoleId, _ := uuid.GenerateUUID()
	// create a new test application with a test app role
	app := testApplicationsClient_Create(t, c, msgraph.Application{
		DisplayName: utils.StringPtr(fmt.Sprintf("test-serviceprincipal-appRoleAssignments-%s", c.RandomString)),
		AppRoles: &[]msgraph.AppRole{
			{
				ID:          utils.StringPtr(testResourceAppRoleId),
				DisplayName: utils.StringPtr(fmt.Sprintf("test-resourceApp-role-%s", c.RandomString)),
				IsEnabled:   utils.BoolPtr(true),
				Description: utils.StringPtr(fmt.Sprintf("test-resourceApp-role-description-%s", c.RandomString)),
				Value:       utils.StringPtr(fmt.Sprintf("test-resourceApp-role-value-%s", c.RandomString)),
				AllowedMemberTypes: &[]msgraph.AppRoleAllowedMemberType{
					msgraph.AppRoleAllowedMemberTypeUser,
					msgraph.AppRoleAllowedMemberTypeApplication,
				},
			},
		},
	})

	sp := testServicePrincipalsClient_Create(t, c, msgraph.ServicePrincipal{
		AccountEnabled: utils.BoolPtr(true),
		AppId:          app.AppId,
		DisplayName:    app.DisplayName,
	})
	testServicePrincipalsClient_Get(t, c, *sp.ID())
	sp.Tags = &([]string{"TestTag"})
	testServicePrincipalsClient_Update(t, c, *sp)
	testServicePrincipalsClient_List(t, c, odata.Query{})

	newGroupParent := msgraph.Group{
		DisplayName:     utils.StringPtr("test-group-parent-servicePrincipals-appRoleAssignments"),
		MailEnabled:     utils.BoolPtr(false),
		MailNickname:    utils.StringPtr(fmt.Sprintf("test-group-parent-%s", c.RandomString)),
		SecurityEnabled: utils.BoolPtr(true),
	}
	newGroupChild := msgraph.Group{
		DisplayName:     utils.StringPtr("test-group-child-servicePrincipals-appRoleAssignments"),
		MailEnabled:     utils.BoolPtr(false),
		MailNickname:    utils.StringPtr(fmt.Sprintf("test-group-child-%s", c.RandomString)),
		SecurityEnabled: utils.BoolPtr(true),
	}

	groupParent := testGroupsClient_Create(t, c, newGroupParent)
	groupChild := testGroupsClient_Create(t, c, newGroupChild)
	groupParent.Members = &msgraph.Members{groupChild.DirectoryObject}
	testGroupsClient_AddMembers(t, c, groupParent)
	groupChild.Members = &msgraph.Members{sp.DirectoryObject}
	testGroupsClient_AddMembers(t, c, groupChild)

	testServicePrincipalsClient_ListGroupMemberships(t, c, *sp.ID())

	// App Role Assignments
	appRoleAssignment := testServicePrincipalsClient_AssignAppRole(t, c, *groupParent.ID(), *sp.ID(), testResourceAppRoleId)
	// list resourceApp role assignments for a test group
	appRoleAssignments := testServicePrincipalsClient_ListAppRoleAssignments(t, c, *sp.ID())
	if len(*appRoleAssignments) == 0 {
		t.Fatal("expected at least one app role assignment assigned to the test group")
	}
	// removes app role assignment previously set to the test group
	testServicePrincipalsClient_RemoveAppRoleAssignment(t, c, *sp.ID(), *appRoleAssignment.Id)

	// remove all test resources
	testGroupsClient_Delete(t, c, *groupParent.ID())
	testGroupsClient_Delete(t, c, *groupChild.ID())
	testServicePrincipalsClient_Delete(t, c, *sp.ID())
	testApplicationsClient_Delete(t, c, *app.ID())

}

func testServicePrincipalsClient_Create(t *testing.T, c *test.Test, sp msgraph.ServicePrincipal) (servicePrincipal *msgraph.ServicePrincipal) {
	servicePrincipal, status, err := c.ServicePrincipalsClient.Create(c.Context, sp)
	if err != nil {
		t.Fatalf("ServicePrincipalsClient.Create(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ServicePrincipalsClient.Create(): invalid status: %d", status)
	}
	if servicePrincipal == nil {
		t.Fatal("ServicePrincipalsClient.Create(): servicePrincipal was nil")
	}
	if servicePrincipal.ID() == nil {
		t.Fatal("ServicePrincipalsClient.Create(): servicePrincipal.ID was nil")
	}
	return
}

func testServicePrincipalsClient_Update(t *testing.T, c *test.Test, sp msgraph.ServicePrincipal) (servicePrincipal *msgraph.ServicePrincipal) {
	status, err := c.ServicePrincipalsClient.Update(c.Context, sp)
	if err != nil {
		t.Fatalf("ServicePrincipalsClient.Update(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ServicePrincipalsClient.Update(): invalid status: %d", status)
	}
	return
}

func testServicePrincipalsClient_List(t *testing.T, c *test.Test, query odata.Query) (servicePrincipals *[]msgraph.ServicePrincipal) {
	query.Top = 10
	servicePrincipals, _, err := c.ServicePrincipalsClient.List(c.Context, query)
	if err != nil {
		t.Fatalf("ServicePrincipalsClient.List(): %v", err)
	}
	if servicePrincipals == nil {
		t.Fatal("ServicePrincipalsClient.List(): servicePrincipals was nil")
	}
	return
}

func testServicePrincipalsClient_Get(t *testing.T, c *test.Test, id string) (servicePrincipal *msgraph.ServicePrincipal) {
	servicePrincipal, status, err := c.ServicePrincipalsClient.Get(c.Context, id, odata.Query{})
	if err != nil {
		t.Fatalf("ServicePrincipalsClient.Get(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ServicePrincipalsClient.Get(): invalid status: %d", status)
	}
	if servicePrincipal == nil {
		t.Fatal("ServicePrincipalsClient.Get(): servicePrincipal was nil")
	}
	return
}

func testServicePrincipalsClient_Delete(t *testing.T, c *test.Test, id string) {
	status, err := c.ServicePrincipalsClient.Delete(c.Context, id)
	if err != nil {
		t.Fatalf("ServicePrincipalsClient.Delete(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ServicePrincipalsClient.Delete(): invalid status: %d", status)
	}
}

func testServicePrincipalsClient_ListGroupMemberships(t *testing.T, c *test.Test, id string) (groups *[]msgraph.Group) {
	groups, _, err := c.ServicePrincipalsClient.ListGroupMemberships(c.Context, id, odata.Query{})
	if err != nil {
		t.Fatalf("ServicePrincipalsClient.ListGroupMemberships(): %v", err)
	}

	if groups == nil {
		t.Fatal("ServicePrincipalsClient.ListGroupMemberships(): groups was nil")
	}

	if len(*groups) != 2 {
		t.Fatalf("ServicePrincipalsClient.ListGroupMemberships(): expected groups length 2. was: %d", len(*groups))
	}

	return
}

func testServicePrincipalsClient_AddPassword(t *testing.T, c *test.Test, a *msgraph.ServicePrincipal) *msgraph.PasswordCredential {
	expiry := time.Now().Add(24 * 90 * time.Hour)
	pwd := msgraph.PasswordCredential{
		DisplayName: utils.StringPtr("test password"),
		EndDateTime: &expiry,
	}
	newPwd, status, err := c.ServicePrincipalsClient.AddPassword(c.Context, *a.ID(), pwd)
	if err != nil {
		t.Fatalf("ServicePrincipalsClient.AddPassword(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ServicePrincipalsClient.AddPassword(): invalid status: %d", status)
	}
	if newPwd.SecretText == nil || len(*newPwd.SecretText) == 0 {
		t.Fatalf("ServicePrincipalsClient.AddPassword(): nil or empty secretText returned by API")
	}
	return newPwd
}

func testServicePrincipalsClient_AddTokenSigningCertificate(t *testing.T, c *test.Test, a *msgraph.ServicePrincipal) *msgraph.KeyCredential {
	expiry := time.Now().Add(24 * 90 * time.Hour)
	tsc := msgraph.KeyCredential{
		DisplayName: utils.StringPtr("cn=test cert"),
		EndDateTime: &expiry,
	}
	newKey, status, err := c.ServicePrincipalsClient.AddTokenSigningCertificate(c.Context, *a.ID(), tsc)
	if err != nil {
		t.Fatalf("ServicePrincipalsClient.AddTokenSigningCertificate(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ServicePrincipalsClient.AddTokenSigningCertificate(): invalid status: %d", status)
	}

	if newKey.Thumbprint == nil || len(*newKey.Thumbprint) == 0 {
		t.Fatalf("ServicePrincipalsClient.AddTokenSigningCertificate(): nil or empty thumbprint returned by API")
	}

	return newKey
}

func testServicePrincipalsClient_SetPreferredTokenSigningKeyThumbprint(t *testing.T, c *test.Test, a *msgraph.ServicePrincipal, thumbprint string) {

	status, err := c.ServicePrincipalsClient.SetPreferredTokenSigningKeyThumbprint(c.Context, *a.ID(), thumbprint)
	if err != nil {
		t.Fatalf("ServicePrincipalsClient.AddTokenSigningCertificate(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ServicePrincipalsClient.AddTokenSigningCertificate(): invalid status: %d", status)
	}
}

func testServicePrincipalsClient_RemovePassword(t *testing.T, c *test.Test, a *msgraph.ServicePrincipal, p *msgraph.PasswordCredential) {
	status, err := c.ServicePrincipalsClient.RemovePassword(c.Context, *a.ID(), *p.KeyId)
	if err != nil {
		t.Fatalf("ServicePrincipalsClient.RemovePassword(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ServicePrincipalsClient.RemovePassword(): invalid status: %d", status)
	}
}

func testServicePrincipalsClient_AddOwners(t *testing.T, c *test.Test, sp *msgraph.ServicePrincipal) {
	status, err := c.ServicePrincipalsClient.AddOwners(c.Context, sp)
	if err != nil {
		t.Fatalf("ServicePrincipalsClient.AddOwners(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ServicePrincipalsClient.AddOwners(): invalid status: %d", status)
	}
}

func testServicePrincipalsClient_ListOwnedObjects(t *testing.T, c *test.Test, id string) (ownedObjects *[]string) {
	ownedObjects, _, err := c.ServicePrincipalsClient.ListOwnedObjects(c.Context, id)
	if err != nil {
		t.Fatalf("ServicePrincipalsClient.ListOwnedObjects(): %v", err)
	}

	if ownedObjects == nil {
		t.Fatal("ServicePrincipalsClient.ListOwnedObjects(): ownedObjects was nil")
	}

	if len(*ownedObjects) != 1 {
		t.Fatalf("ServicePrincipalsClient.ListOwnedObjects(): expected ownedObjects length 1. was: %d", len(*ownedObjects))
	}
	return
}

func testServicePrincipalsClient_ListOwners(t *testing.T, c *test.Test, id string, expected []string) (owners *[]string) {
	owners, status, err := c.ServicePrincipalsClient.ListOwners(c.Context, id)
	if err != nil {
		t.Fatalf("ServicePrincipalsClient.ListOwners(): %v", err)
	}

	if status < 200 || status >= 300 {
		t.Fatalf("ServicePrincipalsClient.ListOwners(): invalid status: %d", status)
	}

	ownersExpected := len(expected)

	if len(*owners) < ownersExpected {
		t.Fatalf("ServicePrincipalsClient.ListOwners(): expected at least %d owner. has: %d", ownersExpected, len(*owners))
	}

	var ownersFound int

	for _, e := range expected {
		for _, o := range *owners {
			if e == o {
				ownersFound++
				continue
			}
		}
	}

	if ownersFound < ownersExpected {
		t.Fatalf("ServicePrincipalsClient.ListOwners(): expected %d matching owners. found: %d", ownersExpected, ownersFound)
	}
	return
}

func testServicePrincipalsClient_GetOwner(t *testing.T, c *test.Test, spId, ownerId string) (owner *string) {
	owner, status, err := c.ServicePrincipalsClient.GetOwner(c.Context, spId, ownerId)
	if err != nil {
		t.Fatalf("ServicePrincipalsClient.GetOwner(): %v", err)
	}

	if status < 200 || status >= 300 {
		t.Fatalf("ServicePrincipalsClient.GetOwner(): invalid status: %d", status)
	}

	if owner == nil {
		t.Fatalf("ServicePrincipalsClient.GetOwner(): owner was nil")
	}
	return
}

func testServicePrincipalsClient_RemoveOwners(t *testing.T, c *test.Test, spId string, ownerIds []string) {
	_, err := c.ServicePrincipalsClient.RemoveOwners(c.Context, spId, &ownerIds)
	if err != nil {
		t.Fatalf("ServicePrincipalsClient.RemoveOwners(): %v", err)
	}
}

func testServicePrincipalsClient_AssignClaimsMappingPolicy(t *testing.T, c *test.Test, sp *msgraph.ServicePrincipal) {
	status, err := c.ServicePrincipalsClient.AssignClaimsMappingPolicy(c.Context, sp)
	if err != nil {
		t.Fatalf("ServicePrincipalsClient.AssignClaimsMappingPolicy(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ServicePrincipalsClient.AssignClaimsMappingPolicy(): invalid status: %d", status)
	}
}

func testServicePrincipalsClient_RemoveClaimsMappingPolicy(t *testing.T, c *test.Test, sp *msgraph.ServicePrincipal, policyIds []string) {
	status, err := c.ServicePrincipalsClient.RemoveClaimsMappingPolicy(c.Context, sp, &policyIds)
	if err != nil {
		t.Fatalf("ServicePrincipalsClient.RemoveClaimsMappingPolicy(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ServicePrincipalsClient.RemoveClaimsMappingPolicy(): invalid status: %d", status)
	}
}

func testServicePrincipalsClient_AssignTokenIssuancePolicy(t *testing.T, c *test.Test, sp *msgraph.ServicePrincipal) {
	status, err := c.ServicePrincipalsClient.AssignTokenIssuancePolicy(c.Context, sp)
	if err != nil {
		t.Fatalf("ServicePrincipalsClient.AssignTokenIssuancePolicy(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ServicePrincipalsClient.AssignTokenIssuancePolicy(): invalid status: %d", status)
	}
}

func testServicePrincipalsClient_RemoveTokenIssuancePolicy(t *testing.T, c *test.Test, sp *msgraph.ServicePrincipal, policyIds []string) {
	status, err := c.ServicePrincipalsClient.RemoveTokenIssuancePolicy(c.Context, sp, &policyIds)
	if err != nil {
		t.Fatalf("ServicePrincipalsClient.RemoveTokenIssuancePolicy(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ServicePrincipalsClient.RemoveTokenIssuancePolicy(): invalid status: %d", status)
	}
}

func testServicePrincipalsClient_AssignAppRole(t *testing.T, c *test.Test, principalId, resourceId, appRoleId string) (appRoleAssignment *msgraph.AppRoleAssignment) {
	appRoleAssignment, status, err := c.ServicePrincipalsClient.AssignAppRoleForResource(c.Context, principalId, resourceId, appRoleId)
	if err != nil {
		t.Fatalf("ServicePrincipalsClient.AssignAppRoleForResource(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ServicePrincipalsClient.AssignAppRoleForResource(): invalid status: %d", status)
	}
	if appRoleAssignment == nil {
		t.Fatal("ServicePrincipalsClient.AssignAppRoleForResource(): appRoleAssignment was nil")
	}
	if appRoleAssignment.Id == nil {
		t.Fatal("ServicePrincipalsClient.AssignAppRoleForResource(): appRoleAssignment.Id was nil")
	}
	return
}

func testServicePrincipalsClient_ListAppRoleAssignments(t *testing.T, c *test.Test, resourceId string) (appRoleAssignments *[]msgraph.AppRoleAssignment) {
	appRoleAssignments, _, err := c.ServicePrincipalsClient.ListAppRoleAssignments(c.Context, resourceId, odata.Query{})
	if err != nil {
		t.Fatalf("ServicePrincipalsClient.ListAppRoleAssignments(): %v", err)
	}
	if appRoleAssignments == nil {
		t.Fatal("ServicePrincipalsClient.ListAppRoleAssignments(): appRoleAssignments was nil")
	}
	return
}

func testServicePrincipalsClient_RemoveAppRoleAssignment(t *testing.T, c *test.Test, resourceId, appRoleAssignmentId string) {
	status, err := c.ServicePrincipalsClient.RemoveAppRoleAssignment(c.Context, resourceId, appRoleAssignmentId)
	if err != nil {
		t.Fatalf("ServicePrincipalsClient.RemoveAppRoleAssignment(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ServicePrincipalsClient.RemoveAppRoleAssignment(): invalid status: %d", status)
	}
}
