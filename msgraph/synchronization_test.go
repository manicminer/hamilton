package msgraph_test

import (
	"fmt"
	"testing"

	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/msgraph"
)

const testApplicationTemplateIdDatabricks = "9c9818d2-2900-49e8-8ba4-22688be7c675" // Databricks SCIM connector

func TestSynchronizationClient(t *testing.T) {
	c := test.NewTest(t)
	defer c.CancelFunc()

	template := testApplicationTemplatesClient_Get(t, c, testApplicationTemplateIdDatabricks)
	app := testApplicationTemplatesClient_Instantiate(t, c, msgraph.ApplicationTemplate{
		ID:          template.ID,
		DisplayName: utils.StringPtr(fmt.Sprintf("test-applicationTemplate-%s", c.RandomString)),
	})
	testSynchronizationJobClient_SetSecrets(t, c, msgraph.SynchronizationSecret{
		Credentials: &[]msgraph.SynchronizationSecretKeyStringValuePair{
			{
				Key:   utils.StringPtr("BaseAddress"),
				Value: utils.StringPtr("https://test-address.azuredatabricks.net"),
			},
			{
				Key:   utils.StringPtr("SecretToken"),
				Value: utils.StringPtr("dummy-token"),
			},
		},
	}, *app.ServicePrincipal.ID())
	testSynchronizationJobClient_GetSecrets(t, c, *app.ServicePrincipal.ID())
	job := testSynchronizationJobClient_Create(t, c, msgraph.SynchronizationJob{
		Schedule: &msgraph.SynchronizationSchedule{
			State: utils.StringPtr("Disabled"),
		},
		TemplateId: utils.StringPtr("dataBricks"),
	}, *app.ServicePrincipal.ID())

	testSynchronizationJobClient_Get(t, c, *job.ID, *app.ServicePrincipal.ID())
	testSynchronizationJobClient_Start(t, c, *job.ID, *app.ServicePrincipal.ID())
	testSynchronizationJobClient_List(t, c, *app.ServicePrincipal.ID())
	testSynchronizationJobClient_Pause(t, c, *job.ID, *app.ServicePrincipal.ID())
	testSynchronizationJobClient_Restart(t, c, *job.ID, msgraph.SynchronizationJobRestartCriteria{}, *app.ServicePrincipal.ID())
	testSynchronizationJobClient_Delete(t, c, *job.ID, *app.ServicePrincipal.ID())

	// We don't test validateCredentials as this requires provisioning a valid enterprise application

	testServicePrincipalsClient_Delete(t, c, *app.ServicePrincipal.ID())

	testApplicationsClient_Delete(t, c, *app.Application.ID())
	testApplicationsClient_DeletePermanently(t, c, *app.Application.ID())
}

func testSynchronizationJobClient_GetSecrets(t *testing.T, c *test.Test, servicePrincipalId string) (synchronizationSecret *msgraph.SynchronizationSecret) {
	synchronizationSecret, status, err := c.SynchronizationJobClient.GetSecrets(c.Context, servicePrincipalId)
	if err != nil {
		t.Fatalf("SynchronizationJobClient.GetSecrets(): %v", err)
	}

	if status < 200 || status >= 300 {
		t.Fatalf("SynchronizationJobClient.GetSecrets(): invalid status: %d", status)
	}
	return
}

func testSynchronizationJobClient_SetSecrets(t *testing.T, c *test.Test, s msgraph.SynchronizationSecret, servicePrincipalId string) {
	status, err := c.SynchronizationJobClient.SetSecrets(c.Context, s, servicePrincipalId)
	if err != nil {
		t.Fatalf("SynchronizationJobClient.SetSecrets(): %v", err)
	}

	if status < 200 || status >= 300 {
		t.Fatalf("SynchronizationJobClient.SetSecrets(): invalid status: %d", status)
	}
}

func testSynchronizationJobClient_Create(t *testing.T, c *test.Test, a msgraph.SynchronizationJob, servicePrincipalId string) (synchronizationJob *msgraph.SynchronizationJob) {
	synchronizationJob, status, err := c.SynchronizationJobClient.Create(c.Context, a, servicePrincipalId)
	if err != nil {
		t.Fatalf("SynchronizationJobClient.Create(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("SynchronizationJobClient.Create(): invalid status: %d", status)
	}
	if synchronizationJob == nil {
		t.Fatal("SynchronizationJobClient.Create(): synchronizationJob was nil")
	}
	return synchronizationJob
}

func testSynchronizationJobClient_List(t *testing.T, c *test.Test, servicePrincipalId string) (synchronizationJobs *[]msgraph.SynchronizationJob) {
	synchronizationJobs, status, err := c.SynchronizationJobClient.List(c.Context, servicePrincipalId)
	if err != nil {
		t.Fatalf("SynchronizationJobClient.List(): %v", err)
	}

	if status < 200 || status >= 300 {
		t.Fatalf("SynchronizationJobClient.List(): invalid status: %d", status)
	}

	return synchronizationJobs
}

func testSynchronizationJobClient_Get(t *testing.T, c *test.Test, jobId string, servicePrincipalId string) (synchronizationJob *msgraph.SynchronizationJob) {
	synchronizationJob, status, err := c.SynchronizationJobClient.Get(c.Context, jobId, servicePrincipalId)
	if err != nil {
		t.Fatalf("SynchronizationJobClient.Get(): %v", err)
	}

	if status < 200 || status >= 300 {
		t.Fatalf("SynchronizationJobClient.Get(): invalid status: %d", status)
	}

	if synchronizationJob == nil {
		t.Fatalf("SynchronizationJobClient.Get(): synchronizationJob was nil")
	}
	return synchronizationJob
}

func testSynchronizationJobClient_Start(t *testing.T, c *test.Test, jobId string, servicePrincipalId string) {
	status, err := c.SynchronizationJobClient.Start(c.Context, jobId, servicePrincipalId)
	if err != nil {
		t.Fatalf("SynchronizationJobClient.Start(): %v", err)
	}

	if status < 200 || status >= 300 {
		t.Fatalf("SynchronizationJobClient.Start(): invalid status: %d", status)
	}
}

func testSynchronizationJobClient_Pause(t *testing.T, c *test.Test, jobId string, servicePrincipalId string) {
	status, err := c.SynchronizationJobClient.Pause(c.Context, jobId, servicePrincipalId)
	if err != nil {
		t.Fatalf("SynchronizationJobClient.Pause(): %v", err)
	}

	if status < 200 || status >= 300 {
		t.Fatalf("SynchronizationJobClient.Pause(): invalid status: %d", status)
	}
}

func testSynchronizationJobClient_Restart(t *testing.T, c *test.Test, jobId string, synchronizationJobRestartCriteria msgraph.SynchronizationJobRestartCriteria, servicePrincipalId string) {
	status, err := c.SynchronizationJobClient.Restart(c.Context, jobId, synchronizationJobRestartCriteria, servicePrincipalId)
	if err != nil {
		t.Fatalf("SynchronizationJobClient.Restart(): %v", err)
	}

	if status < 200 || status >= 300 {
		t.Fatalf("SynchronizationJobClient.Restart(): invalid status: %d", status)
	}
}

func testSynchronizationJobClient_Delete(t *testing.T, c *test.Test, jobId string, servicePrincipalId string) {
	status, err := c.SynchronizationJobClient.Delete(c.Context, jobId, servicePrincipalId)
	if err != nil {
		t.Fatalf("SynchronizationJobClient.Delete(): %v", err)
	}

	if status < 200 || status >= 300 {
		t.Fatalf("SynchronizationJobClient.Delete(): invalid status: %d", status)
	}
}
