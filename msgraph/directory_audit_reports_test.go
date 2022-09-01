package msgraph_test

import (
	"testing"

	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/msgraph"
	"github.com/manicminer/hamilton/odata"
)

func TestDirectoryAuditReportsTest(t *testing.T) {
	c := test.NewTest(t)
	defer c.CancelFunc()

	auditLogs := testDirectoryAuditReports_List(t, c)
	testDirectoryAuditReports_Get(t, c, *(*auditLogs)[0].Id)
}

func testDirectoryAuditReports_List(t *testing.T, c *test.Test) (dirLogs *[]msgraph.DirectoryAudit) {
	dirLogs, status, err := c.DirectoryAuditReportsClient.List(c.Context, odata.Query{Top: 10})

	if status < 200 || status >= 300 {
		t.Fatalf("DirectoryAuditReportsClient.List(): invalid status: %d", status)
	}

	if err != nil {
		t.Fatalf("DirectoryAuditReportsClient.List(): %v", err)
	}

	if dirLogs == nil {
		t.Fatal("DirectoryAuditReportsClient.List():logs was nil")
	}
	return dirLogs
}

func testDirectoryAuditReports_Get(t *testing.T, c *test.Test, id string) (dirLog *msgraph.DirectoryAudit) {
	dirLog, status, err := c.DirectoryAuditReportsClient.Get(c.Context, id, odata.Query{})
	if err != nil {
		t.Fatalf("DirectoryAuditReportsClient.Get(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("DirectoryAuditReportsClient.Get(): invalid status: %d", status)
	}
	if dirLog == nil {
		t.Fatal("DirectoryAuditReportsClient.Get(): domain was nil")
	}
	return dirLog
}
