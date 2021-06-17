package msgraph_test

import (
	"testing"

	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/msgraph"
)

type SignInReportsClientTest struct {
	connection   *test.Connection
	client       *msgraph.SignInReportsClient
	randomString string
}

func TestSignInReportsTest(t *testing.T) {
	c := SignInReportsClientTest{
		connection:   test.NewConnection(auth.MsGraph, auth.TokenVersion2),
		randomString: test.RandomString(),
	}
	c.client = msgraph.NewSignInLogsClient(c.connection.AuthConfig.TenantID)
	c.client.BaseClient.Authorizer = c.connection.Authorizer

	auditLogs := testSignInReports_List(t, c)
	testSignInReports_Get(t, c, *(*auditLogs)[0].ID)
}

func testSignInReports_List(t *testing.T, c TestSignInReportsTest) (signInLogs *[]msgraph.SignInReport) {
	signInLogs, status, err := c.client.List(c.connection.Context)

	if status < 200 || status >= 300 {
		t.Fatalf("SignInReportsClient.List(): invalid status: %d", status)
	}

	if err != nil {
		t.Fatalf("SignInReportsClient.List(): %v", err)
	}

	if signInLogs == nil {
		t.Fatal("SignInReportsClient.List():logs was nil")
	}
	return
}

func testSignInReports_Get(t *testing.T, c TestSignInReportsTest, id string) (signInLog *msgraph.SignInReport) {
	signInLog, status, err := c.client.Get(c.connection.Context, id)
	if err != nil {
		t.Fatalf("SignInReportsClient.Get(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("SignInReportsClient.Get(): invalid status: %d", status)
	}
	if signInLog == nil {
		t.Fatal("SignInReportsClient.Get(): domain was nil")
	}
	return
}
