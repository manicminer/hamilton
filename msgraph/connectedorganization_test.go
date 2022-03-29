package msgraph_test

import (
	"os"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/msgraph"
	"github.com/manicminer/hamilton/odata"
)

func TestConnectedOrganizationClient(t *testing.T) {
	c := test.NewTest(t)
	defer c.CancelFunc()

	// The connected tenant needs to be a valid tenant, so set in environment config.
	// We can either create the connected organization with a tenant id or with a domain name, test both.
	connectedTenantId := os.Getenv("CONNECTED_TENANT_ID")
	connectedDomain := os.Getenv("CONNECTED_DOMAIN")

	if _, err := uuid.ParseUUID(connectedTenantId); err != nil {
		t.Fatalf("CONNECTED_TENANT_ID is %q. It must be specified as a UUID value.", connectedTenantId)
	}

	if connectedDomain == "" {
		t.Fatalf("CONNECTED_DOMAIN must be specified.")
	}

	// CREATE
	newConnectedOrg := testConnectedOrganizationClient_Create(t, c, GetTestConnectedOrganization(&connectedTenantId))

	// Now delete it
	if newConnectedOrg != nil && newConnectedOrg.ID != nil {
		testConnectedOrganizationClient_Delete(t, c, *newConnectedOrg.ID)
	}

	// and test again with a domain name
	newConnectedOrg = testConnectedOrganizationClient_Create(t, c, GetTestConnectedOrganization(&connectedDomain))

	// LIST
	connectedOrganizations := testConnectedOrganizationClient_List(t, c)

	listedNewOrg := false
	for _, v := range *connectedOrganizations {
		if *v.ID == *newConnectedOrg.ID {
			listedNewOrg = true
		}
	}
	if !listedNewOrg {
		t.Fatalf("Could not find newly created connected organization")
	}

	readConnectedOrg := testConnectedOrganizationClient_Get(t, c, *newConnectedOrg.ID)

	// GET
	if *(*readConnectedOrg.IdentitySources)[0].TenantId != connectedTenantId {
		t.Fatalf("The connected organization should have the source tenant id set, even when created with a domain name.")
	}

	// UPDATE
	newConnectedOrg.Description = utils.StringPtr("Changed description")
	testConnectedOrganizationClient_Update(t, c, newConnectedOrg)

	// DELETE
	if newConnectedOrg != nil && newConnectedOrg.ID != nil {
		testConnectedOrganizationClient_Delete(t, c, *newConnectedOrg.ID)
	}
}

func testConnectedOrganizationClient_Create(t *testing.T, c *test.Test, a *msgraph.ConnectedOrganization) (connectedOrganization *msgraph.ConnectedOrganization) {
	connectedOrganization, status, err := c.ConnectedOrganizationClient.Create(c.Context, *a)

	if err != nil {
		t.Fatalf("ConnectedOrganizationClient.Create(): %v", err)
	}

	if status < 200 || status >= 300 {
		t.Fatalf("ConnectedOrganizationClient.Create(): invalid status: %d", status)
	}

	if connectedOrganization == nil {
		t.Fatal("ConnectedOrganizationClient.Create(): connectedOrganization was nil")
	}

	if connectedOrganization.ID == nil {
		t.Fatal("ConnectedOrganizationClient.Create(): connectedOrganization.ID was nil")
	}

	return
}

func testConnectedOrganizationClient_List(t *testing.T, c *test.Test) (connectedOrganisations *[]msgraph.ConnectedOrganization) {
	connectedOrganisations, _, err := c.ConnectedOrganizationClient.List(c.Context, odata.Query{Top: 10})
	if err != nil {
		t.Fatalf("ConnectedOrganizationClient.List(): %v", err)
	}
	if connectedOrganisations == nil {
		t.Fatal("ConnectedOrganizationClient.List(): connectedOrganisations was nil")
	}
	return
}

func testConnectedOrganizationClient_Get(t *testing.T, c *test.Test, id string) (connectedOrganisation *msgraph.ConnectedOrganization) {
	connectedOrganisation, status, err := c.ConnectedOrganizationClient.Get(c.Context, id, odata.Query{})
	if err != nil {
		t.Fatalf("ConnectedOrganizationClient.Get(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ConnectedOrganizationClient.Get(): invalid status: %d", status)
	}
	if connectedOrganisation == nil {
		t.Fatal("ConnectedOrganizationClient.Get(): connectedOrganisation was nil")
	}
	return
}

func testConnectedOrganizationClient_Delete(t *testing.T, c *test.Test, id string) {
	status, err := c.ConnectedOrganizationClient.Delete(c.Context, id)
	if err != nil {
		t.Fatalf("ConnectedOrganizationClient.Delete(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ConnectedOrganizationClient.Delete(): invalid status: %d", status)
	}
}

func testConnectedOrganizationClient_Update(t *testing.T, c *test.Test, connectedOrganization *msgraph.ConnectedOrganization) {
	status, err := c.ConnectedOrganizationClient.Update(c.Context, *connectedOrganization)
	if err != nil {
		t.Fatalf("ConnectedOrganizationClient.Update(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ConnectedOrganizationClient.Update(): invalid status: %d", status)
	}
}

func GetTestConnectedOrganization(idOrDomain *string) *msgraph.ConnectedOrganization {

	var idSrcs []msgraph.IdentitySource

	_, err := uuid.ParseUUID(*idOrDomain)
	if err == nil {
		idSrcs = []msgraph.IdentitySource{{
			ODataType:   utils.StringPtr(odata.TypeAzureActiveDirectoryTenant),
			TenantId:    idOrDomain,
			DisplayName: utils.StringPtr("Test connected organization"),
		}}
	} else {
		idSrcs = []msgraph.IdentitySource{{
			ODataType:   utils.StringPtr(odata.TypeDomainIdentitySource),
			DomainName:  idOrDomain,
			DisplayName: utils.StringPtr("Test connected organization"),
		}}
	}

	testConnectedOrg := msgraph.ConnectedOrganization{
		Description:     utils.StringPtr("Test Connected Organization"),
		DisplayName:     utils.StringPtr("Test Organization"),
		IdentitySources: &idSrcs,
		State:           utils.StringPtr(msgraph.ConnectedOrganizationStateConfigured),
	}

	return &testConnectedOrg
}
