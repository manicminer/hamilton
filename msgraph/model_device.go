package msgraph

// DeviceManagementDeviceProtectionOverview Hardware information of a given device.
type DeviceManagementDeviceProtectionOverview struct {
	// Indicates number of devices reporting as clean
	CleanDeviceCount *int32 `json:"cleanDeviceCount,omitempty"`
	// Indicates number of devices with critical failures
	CriticalFailuresDeviceCount *int32 `json:"criticalFailuresDeviceCount,omitempty"`
	// Indicates number of devices with inactive threat agent
	InactiveThreatAgentDeviceCount *int32 `json:"inactiveThreatAgentDeviceCount,omitempty"`
	// Indicates number of devices pending full scan
	PendingFullScanDeviceCount *int32 `json:"pendingFullScanDeviceCount,omitempty"`
	// Indicates number of devices with pending manual steps
	PendingManualStepsDeviceCount *int32 `json:"pendingManualStepsDeviceCount,omitempty"`
	// Indicates number of pending offline scan devices
	PendingOfflineScanDeviceCount *int32 `json:"pendingOfflineScanDeviceCount,omitempty"`
	// Indicates the number of devices that have a pending full scan. Valid values -2147483648 to 2147483647
	PendingQuickScanDeviceCount *int32 `json:"pendingQuickScanDeviceCount,omitempty"`
	// Indicates number of devices pending restart
	PendingRestartDeviceCount *int32 `json:"pendingRestartDeviceCount,omitempty"`
	// Indicates number of devices with an old signature
	PendingSignatureUpdateDeviceCount *int32 `json:"pendingSignatureUpdateDeviceCount,omitempty"`
	// Total device count.
	TotalReportedDeviceCount *int32 `json:"totalReportedDeviceCount,omitempty"`
	// Indicates number of devices with threat agent state as unknown
	UnknownStateThreatAgentDeviceCount *int32 `json:"unknownStateThreatAgentDeviceCount,omitempty"`
	OdataType                          string `json:"@odata.type"`
}

// DeviceManagementSettings struct for DeviceManagementSettings
type DeviceManagementSettings struct {
	// The number of days a device is allowed to go without checking in to remain compliant.
	DeviceComplianceCheckinThresholdDays *int32 `json:"deviceComplianceCheckinThresholdDays,omitempty"`
	// Is feature enabled or not for scheduled action for rule.
	IsScheduledActionEnabled *bool `json:"isScheduledActionEnabled,omitempty"`
	// Device should be noncompliant when there is no compliance policy targeted when this is true
	SecureByDefault *bool  `json:"secureByDefault,omitempty"`
	OdataType       string `json:"@odata.type"`
}
