package msgraph

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/manicminer/hamilton/odata"
)

type DeviceManagementClient struct {
	BaseClient Client
}

// NewDeviceManagementClient returns a new DeviceManagementClient.
func NewDeviceManagementClient(tenantId string) *DeviceManagementClient {
	return &DeviceManagementClient{
		BaseClient: NewClient(VersionBeta, tenantId),
	}
}

// Create creates a new device management configuration policy.
func (c *DeviceManagementClient) Create(ctx context.Context, deviceManagementConfigurationPolicy DeviceManagementConfigurationPolicy) (*DeviceManagementConfigurationPolicy, int, error) {
	var status int

	body, err := json.Marshal(deviceManagementConfigurationPolicy)
	if err != nil {
		return nil, status, fmt.Errorf("json.Marshal(): %v", err)
	}

	resp, status, _, err := c.BaseClient.Post(ctx, PostHttpRequestInput{
		Body:             body,
		ValidStatusCodes: []int{http.StatusCreated},
		Uri: Uri{
			Entity:      "/deviceManagement/configurationPolicies",
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("DeviceManagementClient.BaseClient.Post(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var newDeviceManagementConfigurationPolicy DeviceManagementConfigurationPolicy
	if err := json.Unmarshal(respBody, &newDeviceManagementConfigurationPolicy); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	// Get Settings (they're not returned as a part of the post)
	newDeviceManagementConfigurationPolicy.Settings, status, err = c.GetSettings(ctx, *newDeviceManagementConfigurationPolicy.ID)

	return &newDeviceManagementConfigurationPolicy, status, err
}

func (c *DeviceManagementClient) GetSettings(ctx context.Context, id string) (*[]DeviceManagementConfigurationSetting, int, error) {
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		ValidStatusCodes:       []int{http.StatusOK},
		Uri: Uri{
			Entity:      fmt.Sprintf("/deviceManagement/configurationPolicies/%s/settings", id),
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("DeviceManagementClient.BaseClient.Get(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	type settingInstance struct {
		Instance *json.RawMessage `json:"settingInstance"`
	}

	var data struct {
		Settings *[]settingInstance `json:"value"`
	}

	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	// The Graph API returns a mixture of types, this loop matches up the result to the appropriate model
	var ret []DeviceManagementConfigurationSetting

	if data.Settings == nil {
		// Treat this as no result
		return &ret, status, nil
	}

	for _, setting := range *data.Settings {

		instance, err := c.ExpandSettingInstance(setting.Instance)
		if err != nil {
			return nil, status, fmt.Errorf("DeviceManagementClient.ExpandSettingInstance(): %v", err)
		}

		ret = append(ret, DeviceManagementConfigurationSetting{SettingInstance: instance})
	}

	return &ret, status, err
}

func (c *DeviceManagementClient) ExpandSettingInstance(setting *json.RawMessage) (*DeviceManagementConfigurationSettingInstance, error) {
	var o odata.OData

	if err := json.Unmarshal(*setting, &o); err != nil {
		return nil, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	if o.Type == nil {
		return nil, nil
	}

	var ret DeviceManagementConfigurationSettingInstance

	type childInstance struct {
		json.RawMessage
	}

	switch *o.Type {
	case odata.TypeDeviceManagementConfigurationSettingInstance:
		var set BaseDeviceManagementConfigurationSettingInstance
		if err := json.Unmarshal(*setting, &set); err != nil {
			return nil, fmt.Errorf("json.Unmarshal(): %v", err)
		}
		ret = set
	case odata.TypeDeviceManagementConfigurationGroupSettingCollectionInstance:
		var set DeviceManagementConfigurationGroupSettingCollectionInstance
		if err := json.Unmarshal(*setting, &set); err != nil {
			return nil, fmt.Errorf("json.Unmarshal(): %v", err)
		}

		var gscv []DeviceManagementConfigurationGroupSettingValue

		for _, groupSetting := range *set.GroupSettingCollectionValue {
			if groupSetting.Children != nil {

				byte, err := json.Marshal(groupSetting.Children)
				if err != nil {
					return nil, fmt.Errorf("json.Marshal(): %v", err)
				}

				var raw []childInstance
				if err := json.Unmarshal(byte, &raw); err != nil {
					return nil, fmt.Errorf("json.Unmarshal(): %v", err)
				}

				var children []DeviceManagementConfigurationSettingInstance
				for _, child := range raw {

					expandedChild, err := c.ExpandSettingInstance(&child.RawMessage)
					if err != nil {
						return nil, fmt.Errorf("DeviceManagementClient.ExpandSettingInstance(): %v", err)
					}

					children = append(children, expandedChild)
				}

				groupSetting.Children = children

			}
			gscv = append(gscv, groupSetting)

		}
		set.GroupSettingCollectionValue = &gscv

		ret = set

	case odata.TypeDeviceManagementConfigurationChoiceSettingInstance:
		var set DeviceManagementConfigurationChoiceSettingInstance
		if err := json.Unmarshal(*setting, &set); err != nil {
			return nil, fmt.Errorf("json.Unmarshal(): %v", err)
		}

		if set.ChoiceSettingValue != nil {
			if set.ChoiceSettingValue.Children != nil {

				byte, err := json.Marshal(set.ChoiceSettingValue.Children)
				if err != nil {
					return nil, fmt.Errorf("json.Marshal(): %v", err)
				}

				var raw []childInstance
				if err := json.Unmarshal(byte, &raw); err != nil {
					return nil, fmt.Errorf("json.Unmarshal(): %v", err)
				}

				var children []DeviceManagementConfigurationSettingInstance
				for _, child := range raw {

					expandedChild, err := c.ExpandSettingInstance(&child.RawMessage)
					if err != nil {
						return nil, fmt.Errorf("DeviceManagementClient.BaseClient.ExpandSettingInstance(): %v", err)
					}

					children = append(children, expandedChild)
				}

				set.ChoiceSettingValue.Children = children

			}

		}
		ret = set
	case odata.TypeDeviceManagementConfigurationSimpleSettingInstance:
		var set DeviceManagementConfigurationSimpleSettingInstance
		if err := json.Unmarshal(*setting, &set); err != nil {
			return nil, fmt.Errorf("json.Unmarshal(): %v", err)
		}
		ret = set
	}

	return &ret, nil
}
