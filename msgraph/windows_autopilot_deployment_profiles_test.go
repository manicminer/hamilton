package msgraph_test

import (
	"fmt"
	"testing"

	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/msgraph"
	"github.com/manicminer/hamilton/odata"
)

func TestWindowsAutopilotDeploymentProfilesClient(t *testing.T) {
	c := test.NewTest(t)
	defer c.CancelFunc()

	profile := testWindowsAutopilotDeploymentProfilesClient_Create(t, c, msgraph.WindowsAutopilotDeploymentProfile{
		ODataType: utils.StringPtr(odata.TypeAzureADWindowsAutopilotDeploymentProfile),
		// displayName must not contain hyphen characters
		DisplayName: utils.StringPtr(fmt.Sprintf("test windows autopilot deployment profile %s", c.RandomString)),
		Description: utils.StringPtr("Test Windows Autopilot Deployment Profile"),
		OutOfBoxExperienceSettings: &msgraph.OutOfBoxExperienceSettings{
			HidePrivacySettings:       utils.BoolPtr(true),
			HideEULA:                  utils.BoolPtr(true),
			UserType:                  utils.StringPtr(msgraph.WindowsUserTypeStandard),
			DeviceUsageType:           utils.StringPtr(msgraph.DeviceUsageTypeSingleUser),
			SkipKeyboardSelectionPage: utils.BoolPtr(true),
			HideEscapeLink:            utils.BoolPtr(true),
		},
		EnrollmentStatusScreenSettings: &msgraph.WindowsEnrollmentStatusScreenSettings{
			HideInstallationProgress:                         utils.BoolPtr(true),
			AllowDeviceUseBeforeProfileAndAppInstallComplete: utils.BoolPtr(true),
			BlockDeviceSetupRetryByUser:                      utils.BoolPtr(true),
			AllowLogCollectionOnInstallFailure:               utils.BoolPtr(true),
			CustomErrorMessage:                               utils.StringPtr("Test Custom Error Message"),
			InstallProgressTimeoutInMinutes:                  utils.Int32Ptr(15),
			AllowDeviceUseOnInstallFailure:                   utils.BoolPtr(true),
		},
		ExtractHardwareHash:                    utils.BoolPtr(true),
		DeviceType:                             utils.StringPtr(msgraph.WindowsAutopilotDeviceTypeWindowsPc),
		EnableWhiteGlove:                       utils.BoolPtr(false),
		Language:                               utils.StringPtr("os-default"),
		HybridAzureADJoinSkipConnectivityCheck: utils.BoolPtr(false),
	})

	updateProfile := msgraph.WindowsAutopilotDeploymentProfile{
		ODataType:   utils.StringPtr(odata.TypeAzureADWindowsAutopilotDeploymentProfile),
		ID:          profile.ID,
		Description: utils.StringPtr("Test Windows Autopilot Deployment Profile Update"),
	}

	testWindowsAutopilotDeploymentProfilesClient_Update(t, c, updateProfile)

	testWindowsAutopilotDeploymentProfilesClient_List(t, c)
	testWindowsAutopilotDeploymentProfilesClient_Get(t, c, *profile.ID)
	testWindowsAutopilotDeploymentProfilesClient_Delete(t, c, *profile.ID)
}

func testWindowsAutopilotDeploymentProfilesClient_Create(t *testing.T, c *test.Test, profile msgraph.WindowsAutopilotDeploymentProfile) (windowsAutopilotDeploymentProfile *msgraph.WindowsAutopilotDeploymentProfile) {
	windowsAutopilotDeploymentProfile, status, err := c.WindowsAutopilotDeploymentProfilesClient.Create(c.Context, profile)
	if err != nil {
		t.Fatalf("WindowsAutopilotDeploymentProfilesClient.Create(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("WindowsAutopilotDeploymentProfilesClient.Create(): invalid status: %d", status)
	}
	if windowsAutopilotDeploymentProfile == nil {
		t.Fatal("WindowsAutopilotDeploymentProfilesClient.Create(): windowsAutopilotDeploymentProfile was nil")
	}
	if windowsAutopilotDeploymentProfile.ID == nil {
		t.Fatal("WindowsAutopilotDeploymentProfilesClient.Create(): windowsAutopilotDeploymentProfile.ID was nil")
	}
	return
}

func testWindowsAutopilotDeploymentProfilesClient_Get(t *testing.T, c *test.Test, id string) (profile *msgraph.WindowsAutopilotDeploymentProfile) {
	policy, status, err := c.WindowsAutopilotDeploymentProfilesClient.Get(c.Context, id, odata.Query{})
	if err != nil {
		t.Fatalf("WindowsAutopilotDeploymentProfilesClient.Get(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("WindowsAutopilotDeploymentProfilesClient.Get(): invalid status: %d", status)
	}
	if policy == nil {
		t.Fatal("WindowsAutopilotDeploymentProfilesClient.Get(): policy was nil")
	}
	return
}

func testWindowsAutopilotDeploymentProfilesClient_Update(t *testing.T, c *test.Test, profile msgraph.WindowsAutopilotDeploymentProfile) {
	status, err := c.WindowsAutopilotDeploymentProfilesClient.Update(c.Context, profile)
	if err != nil {
		t.Fatalf("WindowsAutopilotDeploymentProfilesClient.Update(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("WindowsAutopilotDeploymentProfilesClient.Update(): invalid status: %d", status)
	}
}

func testWindowsAutopilotDeploymentProfilesClient_List(t *testing.T, c *test.Test) (profiles *[]msgraph.WindowsAutopilotDeploymentProfile) {
	policies, _, err := c.WindowsAutopilotDeploymentProfilesClient.List(c.Context, odata.Query{Top: 10})
	if err != nil {
		t.Fatalf("WindowsAutopilotDeploymentProfilesClient.List(): %v", err)
	}
	if policies == nil {
		t.Fatal("WindowsAutopilotDeploymentProfilesClient.List(): profiles was nil")
	}
	return
}

func testWindowsAutopilotDeploymentProfilesClient_Delete(t *testing.T, c *test.Test, id string) {
	status, err := c.WindowsAutopilotDeploymentProfilesClient.Delete(c.Context, id)
	if err != nil {
		t.Fatalf("WindowsAutopilotDeploymentProfilesClient.Delete(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("WindowsAutopilotDeploymentProfilesClient.Delete(): invalid status: %d", status)
	}
}
