package msgraph_test

import (
	"testing"

	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/msgraph"
	"github.com/manicminer/hamilton/odata"
)

func TestDeviceManagementClient(t *testing.T) {
	c := test.NewTest(t)
	defer c.CancelFunc()

	newDeviceManagementConfigurationPolicy := msgraph.DeviceManagementConfigurationPolicy{
		Name:            utils.StringPtr("BasicTest3"),
		Description:     utils.StringPtr("BasicTest3"),
		Platforms:       utils.StringPtr("iOS"),
		Technologies:    utils.StringPtr("mdm,appleRemoteManagement"),
		RoleScopeTagIds: &[]string{"0"},
		Settings: &[]msgraph.DeviceManagementConfigurationSetting{
			{
				ODataType: utils.StringPtr(odata.TypeDeviceManagementConfigurationSetting),
				SettingInstance: &msgraph.DeviceManagementConfigurationGroupSettingCollectionInstance{
					BaseDeviceManagementConfigurationSettingInstance: &msgraph.BaseDeviceManagementConfigurationSettingInstance{
						ODataType:           utils.StringPtr(odata.TypeDeviceManagementConfigurationGroupSettingCollectionInstance),
						SettingDefinitionId: utils.StringPtr("com.apple.applicationaccess_com.apple.applicationaccess"),
					},
					GroupSettingCollectionValue: &[]msgraph.DeviceManagementConfigurationGroupSettingValue{
						{
							Children: &[]msgraph.DeviceManagementConfigurationChoiceSettingInstance{
								{
									BaseDeviceManagementConfigurationSettingInstance: &msgraph.BaseDeviceManagementConfigurationSettingInstance{
										ODataType:           utils.StringPtr(odata.TypeDeviceManagementConfigurationChoiceSettingInstance),
										SettingDefinitionId: utils.StringPtr("com.apple.applicationaccess_allowairprint"),
									},
									ChoiceSettingValue: &msgraph.DeviceManagementConfigurationChoiceSettingValue{
										DeviceManagementConfigurationSettingValue: &msgraph.DeviceManagementConfigurationSettingValue{
											ODataType: utils.StringPtr(odata.TypeDeviceManagementConfigurationChoiceSettingValue),
										},
										Value: utils.StringPtr("com.apple.applicationaccess_allowairprint_true"),
									},
								},
							},
						},
					},
				},
			},
		},
	}

	c.DeviceManagementClient.Create(c.Context, newDeviceManagementConfigurationPolicy)

}
