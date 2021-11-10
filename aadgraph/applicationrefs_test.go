package aadgraph_test

import (
	"testing"

	"github.com/manicminer/hamilton/aadgraph"
	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/internal/test"
)

type ApplicationRefsClientTest struct {
	connection   *test.Connection
	client       *aadgraph.ApplicationRefsClient
	randomString string
}

func TestApplicationRefsClient(t *testing.T) {
	c := ApplicationRefsClientTest{
		connection:   test.NewConnection(auth.TokenVersion1),
		randomString: test.RandomString(),
	}
	c.connection.Authorize(c.connection.AuthConfig.Environment.MsGraph)
	c.client = aadgraph.NewApplicationRefsClient(c.connection.AuthConfig.TenantID)
	c.client.BaseClient.Authorizer = c.connection.Authorizer

	t.Skipf("ApplicationRefs is a private endpoint and cannot be automatically tested")

	//appRef := testApplicationRefsClient_Get(t, c, environments.PublishedApis["AzureActiveDirectoryGraph"])
	//fmt.Printf("%+v", appRef)
	//appRef = testApplicationRefsClient_Get(t, c, environments.PublishedApis["MicrosoftGraph"])
	//fmt.Printf("%+v", appRef)
}

//func testApplicationRefsClient_Get(t *testing.T, c ApplicationRefsClientTest, id environments.ApiAppId) (appRef *aadgraph.ApplicationRef) {
//	appRef, status, err := c.client.Get(c.connection.Context, id)
//	if err != nil {
//		t.Fatalf("ApplicationRefsClient.Get(): %v", err)
//	}
//	if status < 200 || status >= 300 {
//		t.Fatalf("ApplicationRefsClient.Get(): invalid status: %d", status)
//	}
//	if appRef == nil {
//		t.Fatal("ApplicationRefsClient.Get(): appRef was nil")
//	}
//	return
//}
