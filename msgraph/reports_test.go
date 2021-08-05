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
	testReports_GetCredentialUserRegistrationCount(t, c)
	testReports_GetCredentialUserRegistrationDetails(t, c)
	testReports_GetUserCredentialUsageDetails(t, c)
	testReports_GetCredentialUsageSummary(t, c)
	testReports_GetAuthenticationMethodsUsersRegisteredByMethod(t, c)

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

func testReports_GetCredentialUserRegistrationCount(t *testing.T, c ReportsClientTest) (report *[]msgraph.CredentialUserRegistrationCount) {
	report, status, err := c.client.GetCredentialUserRegistrationCount(c.connection.Context, odata.Query{})
	if status < 200 || status >= 300 {
		t.Fatalf("ReportsClient.GetCredentialUserRegistrationCount(): invalid status: %d", status)
	}

	if err != nil {
		t.Fatalf("ReportsClient.GetCredentialUserRegistrationCount(): %v", err)
	}

	if report == nil {
		t.Fatal("ReportsClient.GetCredentialUserRegistrationCount():report was nil")
	}
	return
}

func testReports_GetCredentialUserRegistrationDetails(t *testing.T, c ReportsClientTest) (report *[]msgraph.CredentialUserRegistrationDetails) {
	report, status, err := c.client.GetCredentialUserRegistrationDetails(c.connection.Context, odata.Query{})
	if status < 200 || status >= 300 {
		t.Fatalf("ReportsClient.GetCredentialUserRegistrationDetails(): invalid status: %d", status)
	}

	if err != nil {
		t.Fatalf("ReportsClient.GetCredentialUserRegistrationDetails(): %v", err)
	}

	if report == nil {
		t.Fatal("ReportsClient.GetCredentialUserRegistrationDetails():report was nil")
	}
	return
}

func testReports_GetUserCredentialUsageDetails(t *testing.T, c ReportsClientTest) (report *[]msgraph.UserCredentialUsageDetails) {
	report, status, err := c.client.GetUserCredentialUsageDetails(c.connection.Context, odata.Query{})
	if status < 200 || status >= 300 {
		t.Fatalf("ReportsClient.GetUserCredentialUsageDetails(): invalid status: %d", status)
	}

	if err != nil {
		t.Fatalf("ReportsClient.GetUserCredentialUsageDetails(): %v", err)
	}

	if report == nil {
		t.Fatal("ReportsClient.GetUserCredentialUsageDetails():report was nil")
	}
	return
}

func testReports_GetCredentialUsageSummary(t *testing.T, c ReportsClientTest) (report *[]msgraph.CredentialUsageSummary) {
	report, status, err := c.client.GetCredentialUsageSummary(c.connection.Context, msgraph.CredentialUsageSummaryPeriod1, odata.Query{})
	if status < 200 || status >= 300 {
		t.Fatalf("ReportsClient.GetCredentialUsageSummary(): invalid status: %d", status)
	}

	if err != nil {
		t.Fatalf("ReportsClient.GetCredentialUsageSummary(): %v", err)
	}

	if report == nil {
		t.Fatal("ReportsClient.GetCredentialUsageSummary():report was nil")
	}
	return
}

func testReports_GetAuthenticationMethodsUsersRegisteredByMethod(t *testing.T, c ReportsClientTest) (report *msgraph.UserRegistrationMethodSummary) {
	report, status, err := c.client.GetAuthenticationMethodsUsersRegisteredByMethod(c.connection.Context, odata.Query{})
	if status < 200 || status >= 300 {
		t.Fatalf("ReportsClient.GetAuthenticationMethodsUsersRegisteredByMethod(): invalid status: %d", status)
	}

	if err != nil {
		t.Fatalf("ReportsClient.GetAuthenticationMethodsUsersRegisteredByMethod(): %v", err)
	}

	if report == nil {
		t.Fatal("ReportsClient.GetAuthenticationMethodsUsersRegisteredByMethod():report was nil")
	}
	return
}
