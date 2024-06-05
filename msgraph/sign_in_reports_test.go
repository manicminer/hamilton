package msgraph_test

import (
	"testing"

	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/msgraph"
)

func TestSignInReportsTest(t *testing.T) {
	c := test.NewTest(t)
	defer c.CancelFunc()

	signInLogs := testSignInReports_List(t, c)
	if *signInLogs != nil && len(*signInLogs) > 0 {
		testSignInReports_Get(t, c, *(*signInLogs)[0].Id)
	}
}

func testSignInReports_List(t *testing.T, c *test.Test) (signInLogs *[]msgraph.SignInReport) {
	signInLogs, status, err := c.SignInReportsClient.List(c.Context, odata.Query{Top: 10})

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

func testSignInReports_Get(t *testing.T, c *test.Test, id string) (signInLog *msgraph.SignInReport) {
	signInLog, status, err := c.SignInReportsClient.Get(c.Context, id, odata.Query{})
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
