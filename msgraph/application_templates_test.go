package msgraph_test

import (
	"fmt"
	"testing"

	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/msgraph"
	"github.com/manicminer/hamilton/odata"
)

const testApplicationTemplateId = "4601ed45-8ff3-4599-8377-b6649007e876" // Marketo

func TestApplicationTemplatesClient(t *testing.T) {
	c := test.NewTest(t)
	defer c.CancelFunc()

	testApplicationTemplatesClient_List(t, c, odata.Query{})
	testApplicationTemplatesClient_List(t, c, odata.Query{Filter: fmt.Sprintf("categories/any(c:contains(c, '%s'))", msgraph.ApplicationTemplateCategoryEducation)})
	template := testApplicationTemplatesClient_Get(t, c, testApplicationTemplateId)
	app := testApplicationTemplatesClient_Instantiate(t, c, msgraph.ApplicationTemplate{
		ID:          template.ID,
		DisplayName: utils.StringPtr(fmt.Sprintf("test-applicationTemplate-%s", c.RandomString)),
	})

	testServicePrincipalsClient_Delete(t, c, *app.ServicePrincipal.ID)

	testApplicationsClient_Delete(t, c, *app.Application.ID)
	testApplicationsClient_DeletePermanently(t, c, *app.Application.ID)
}

func testApplicationTemplatesClient_List(t *testing.T, c *test.Test, o odata.Query) (applicationTemplates []msgraph.ApplicationTemplate) {
	result, _, err := c.ApplicationTemplatesClient.List(c.Context, o)
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

func testApplicationTemplatesClient_Get(t *testing.T, c *test.Test, id string) (applicationTemplate *msgraph.ApplicationTemplate) {
	applicationTemplate, status, err := c.ApplicationTemplatesClient.Get(c.Context, id, odata.Query{})
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

func testApplicationTemplatesClient_Instantiate(t *testing.T, c *test.Test, a msgraph.ApplicationTemplate) (applicationTemplate *msgraph.ApplicationTemplate) {
	applicationTemplate, status, err := c.ApplicationTemplatesClient.Instantiate(c.Context, a)
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
