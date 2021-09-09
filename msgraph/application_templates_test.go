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

type ApplicationTemplatesClientTest struct {
	connection   *test.Connection
	client       *msgraph.ApplicationTemplatesClient
	randomString string
}

const testApplicationTemplateId = "4601ed45-8ff3-4599-8377-b6649007e876" // Marketo

func TestApplicationTemplatesClient(t *testing.T) {
	rs := test.RandomString()
	c := ApplicationTemplatesClientTest{
		connection:   test.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: rs,
	}
	c.client = msgraph.NewApplicationTemplatesClient(c.connection.AuthConfig.TenantID)
	c.client.BaseClient.Authorizer = c.connection.Authorizer

	testApplicationTemplatesClient_List(t, c, odata.Query{})
	testApplicationTemplatesClient_List(t, c, odata.Query{Filter: fmt.Sprintf("categories/any(c:contains(c, '%s'))", msgraph.ApplicationTemplateCategoryEducation)})
	template := testApplicationTemplatesClient_Get(t, c, testApplicationTemplateId)
	app := testApplicationTemplatesClient_Instantiate(t, c, msgraph.ApplicationTemplate{
		ID:          template.ID,
		DisplayName: utils.StringPtr(fmt.Sprintf("test-applicationTemplate-%s", c.randomString)),
	})

	s := ServicePrincipalsClientTest{
		connection:   test.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: rs,
	}
	s.client = msgraph.NewServicePrincipalsClient(c.connection.AuthConfig.TenantID)
	s.client.BaseClient.Authorizer = c.connection.Authorizer

	testServicePrincipalsClient_Delete(t, s, *app.ServicePrincipal.ID)

	a := ApplicationsClientTest{
		connection:   test.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: rs,
	}
	a.client = msgraph.NewApplicationsClient(c.connection.AuthConfig.TenantID)
	a.client.BaseClient.Authorizer = c.connection.Authorizer

	testApplicationsClient_Delete(t, a, *app.Application.ID)
	testApplicationsClient_DeletePermanently(t, a, *app.Application.ID)
}

func testApplicationTemplatesClient_List(t *testing.T, c ApplicationTemplatesClientTest, o odata.Query) (applicationTemplates []msgraph.ApplicationTemplate) {
	result, _, err := c.client.List(c.connection.Context, o)
	if err != nil {
		t.Fatalf("ApplicationTemplatesClient.List(): %v", err)
	}
	if result == nil {
		t.Fatal("ApplicationsTemplateClient.List(): result was nil")
	}
	applicationTemplates = *result
	if len(applicationTemplates) == 0 {
		t.Fatal("ApplicationsTemplateClient.List(): applicationTemplates was empty")
	}
	if applicationTemplates[0].ID == nil {
		t.Fatal("ApplicationsTemplateClient.List(): first result of applicationTemplates has nil ID")
	}
	return
}

func testApplicationTemplatesClient_Get(t *testing.T, c ApplicationTemplatesClientTest, id string) (applicationTemplate *msgraph.ApplicationTemplate) {
	applicationTemplate, status, err := c.client.Get(c.connection.Context, id, odata.Query{})
	if err != nil {
		t.Fatalf("ApplicationTemplatesClient.Get(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ApplicationTemplatesClient.Get(): invalid status: %d", status)
	}
	if applicationTemplate == nil {
		t.Fatal("ApplicationTemplatesClient.Get(): applicationTemplate was nil")
	}
	if applicationTemplate.ID == nil {
		t.Fatal("ApplicationTemplatesClient.Get(): applicationTemplate.ID was nil")
	}
	return
}

func testApplicationTemplatesClient_Instantiate(t *testing.T, c ApplicationTemplatesClientTest, a msgraph.ApplicationTemplate) (applicationTemplate *msgraph.ApplicationTemplate) {
	applicationTemplate, status, err := c.client.Instantiate(c.connection.Context, a)
	if err != nil {
		t.Fatalf("ApplicationTemplatesClient.Instantiate(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ApplicationTemplatesClient.Instantiate(): invalid status: %d", status)
	}
	if applicationTemplate == nil {
		t.Fatal("ApplicationsTemplateClient.Instantiate(): applicationTemplate was nil")
	}
	if applicationTemplate.Application == nil {
		t.Fatal("ApplicationTemplatesClient.Instantiate(): applicationTemplate.Application was nil")
	}
	if applicationTemplate.ServicePrincipal == nil {
		t.Fatal("ApplicationTemplatesClient.Instantiate(): applicationTemplate.ServicePrincipal was nil")
	}
	return
}
