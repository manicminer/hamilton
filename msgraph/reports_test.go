package msgraph_test

import (
	"testing"

	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/msgraph"
	"github.com/manicminer/hamilton/odata"
)

type ReportsClientTest struct {
	connection *test.Connection
	client     *msgraph.ReportsClient
}

func TestReports(t *testing.T) {
	c := ReportsClientTest{
		connection: test.NewConnection(auth.MsGraph, auth.TokenVersion2),
	}
	c.client = msgraph.NewReportsClient(c.connection.AuthConfig.TenantID)
	c.client.BaseClient.Authorizer = c.connection.Authorizer
	testReports_GetAuthenticationMethodsUsersRegisteredByFeature(t, c)

}

func testReports_GetAuthenticationMethodsUsersRegisteredByFeature(t *testing.T, c ReportsClientTest) (report *msgraph.UserRegistrationFeatureSummary) {
	report, status, err := c.client.GetAuthenticationMethodsUsersRegisteredByFeature(c.connection.Context, odata.Query{})
	if status < 200 || status >= 300 {
		t.Fatalf("ReportsClient.GetAuthenticationMethodsUsersRegisteredByFeature(): invalid status: %d", status)
	}

	if err != nil {
		t.Fatalf("ReportsClient.GetAuthenticationMethodsUsersRegisteredByFeature(): %v", err)
	}

	if report == nil {
		t.Fatal("ReportsClient.GetAuthenticationMethodsUsersRegisteredByFeature():report was nil")
	}
	return
}
