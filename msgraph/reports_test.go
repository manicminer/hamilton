package msgraph_test

import (
	"testing"

	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/msgraph"
)

func TestReports(t *testing.T) {
	c := test.NewTest(t)
	defer c.CancelFunc()

	testReports_GetAuthenticationMethodsUsersRegisteredByFeature(t, c)
	testReports_GetCredentialUserRegistrationCount(t, c)
	testReports_GetCredentialUserRegistrationDetails(t, c)
	testReports_GetUserCredentialUsageDetails(t, c)
	testReports_GetCredentialUsageSummary(t, c)
	testReports_GetAuthenticationMethodsUsersRegisteredByMethod(t, c)

}

func testReports_GetAuthenticationMethodsUsersRegisteredByFeature(t *testing.T, c *test.Test) (report *msgraph.UserRegistrationFeatureSummary) {
	report, status, err := c.ReportsClient.GetAuthenticationMethodsUsersRegisteredByFeature(c.Context, odata.Query{})
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

func testReports_GetCredentialUserRegistrationCount(t *testing.T, c *test.Test) (report *[]msgraph.CredentialUserRegistrationCount) {
	report, status, err := c.ReportsClient.GetCredentialUserRegistrationCount(c.Context, odata.Query{})
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

func testReports_GetCredentialUserRegistrationDetails(t *testing.T, c *test.Test) (report *[]msgraph.CredentialUserRegistrationDetails) {
	report, status, err := c.ReportsClient.GetCredentialUserRegistrationDetails(c.Context, odata.Query{})
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

func testReports_GetUserCredentialUsageDetails(t *testing.T, c *test.Test) (report *[]msgraph.UserCredentialUsageDetails) {
	report, status, err := c.ReportsClient.GetUserCredentialUsageDetails(c.Context, odata.Query{})
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

func testReports_GetCredentialUsageSummary(t *testing.T, c *test.Test) (report *[]msgraph.CredentialUsageSummary) {
	report, status, err := c.ReportsClient.GetCredentialUsageSummary(c.Context, msgraph.CredentialUsageSummaryPeriod1, odata.Query{})
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

func testReports_GetAuthenticationMethodsUsersRegisteredByMethod(t *testing.T, c *test.Test) (report *msgraph.UserRegistrationMethodSummary) {
	report, status, err := c.ReportsClient.GetAuthenticationMethodsUsersRegisteredByMethod(c.Context, odata.Query{})
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
