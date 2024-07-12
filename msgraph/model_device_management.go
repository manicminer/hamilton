package msgraph

import (
	"encoding/json"
	"fmt"
	"time"
)

////////// device management from: https://github.com/microsoftgraph/msgraph-metadata/blob/master/openapi/v1.0/openapi.yaml

type Entity struct {
	// The unique identifier for an entity. Read-only.
	Id        *string `json:"id,omitempty"`
	OdataType string  `json:"@odata.type"`
}

// DeviceManagement struct for DeviceManagement
type DeviceManagement struct {
	Entity
	DeviceProtectionOverview *DeviceManagementDeviceProtectionOverview `json:"deviceProtectionOverview,omitempty"`
	// Intune Account Id for given tenant
	IntuneAccountId                  *string                                           `json:"intuneAccountId,omitempty"`
	IntuneBrand                      *DeviceManagementIntuneBrand                      `json:"intuneBrand,omitempty"`
	Settings                         *DeviceManagementSettings                         `json:"settings,omitempty"`
	SubscriptionState                *DeviceManagementSubscriptionState                `json:"subscriptionState,omitempty"`
	UserExperienceAnalyticsSettings  *DeviceManagementUserExperienceAnalyticsSettings  `json:"userExperienceAnalyticsSettings,omitempty"`
	WindowsMalwareOverview           *DeviceManagementWindowsMalwareOverview           `json:"windowsMalwareOverview,omitempty"`
	ApplePushNotificationCertificate *DeviceManagementApplePushNotificationCertificate `json:"applePushNotificationCertificate,omitempty"`
	// The Audit Events
	AuditEvents []AuditEvent `json:"auditEvents,omitempty"`
	// The list of Compliance Management Partners configured by the tenant.
	ComplianceManagementPartners []ComplianceManagementPartner        `json:"complianceManagementPartners,omitempty"`
	ConditionalAccessSettings    *OnPremisesConditionalAccessSettings `json:"conditionalAccessSettings,omitempty"`
	// The list of detected apps associated with a device.
	DetectedApps []DetectedApp `json:"detectedApps,omitempty"`
	// The list of device categories with the tenant.
	DeviceCategories []DeviceCategory `json:"deviceCategories,omitempty"`
	// The device compliance policies.
	DeviceCompliancePolicies                 []DeviceCompliancePolicy                                  `json:"deviceCompliancePolicies,omitempty"`
	DeviceCompliancePolicyDeviceStateSummary *DeviceManagementDeviceCompliancePolicyDeviceStateSummary `json:"deviceCompliancePolicyDeviceStateSummary,omitempty"`
	// The summary states of compliance policy settings for this account.
	DeviceCompliancePolicySettingStateSummaries []DeviceCompliancePolicySettingStateSummary              `json:"deviceCompliancePolicySettingStateSummaries,omitempty"`
	DeviceConfigurationDeviceStateSummaries     *DeviceManagementDeviceConfigurationDeviceStateSummaries `json:"deviceConfigurationDeviceStateSummaries,omitempty"`
	// The device configurations.
	DeviceConfigurations []DeviceConfiguration `json:"deviceConfigurations,omitempty"`
	// The list of device enrollment configurations
	DeviceEnrollmentConfigurations []DeviceEnrollmentConfiguration `json:"deviceEnrollmentConfigurations,omitempty"`
	// The list of Device Management Partners configured by the tenant.
	DeviceManagementPartners []DeviceManagementPartner `json:"deviceManagementPartners,omitempty"`
	// The list of Exchange Connectors configured by the tenant.
	ExchangeConnectors []DeviceManagementExchangeConnector `json:"exchangeConnectors,omitempty"`
	// Collection of imported Windows autopilot devices.
	ImportedWindowsAutopilotDeviceIdentities []ImportedWindowsAutopilotDeviceIdentity `json:"importedWindowsAutopilotDeviceIdentities,omitempty"`
	// The IOS software update installation statuses for this account.
	IosUpdateStatuses     []IosUpdateDeviceStatus                `json:"iosUpdateStatuses,omitempty"`
	ManagedDeviceOverview *DeviceManagementManagedDeviceOverview `json:"managedDeviceOverview,omitempty"`
	// The list of managed devices.
	ManagedDevices []ManagedDevice `json:"managedDevices,omitempty"`
	// The collection property of MobileAppTroubleshootingEvent.
	MobileAppTroubleshootingEvents []MobileAppTroubleshootingEvent `json:"mobileAppTroubleshootingEvents,omitempty"`
	// The list of Mobile threat Defense connectors configured by the tenant.
	MobileThreatDefenseConnectors []MobileThreatDefenseConnector `json:"mobileThreatDefenseConnectors,omitempty"`
	// The Notification Message Templates.
	NotificationMessageTemplates []NotificationMessageTemplate `json:"notificationMessageTemplates,omitempty"`
	// The remote assist partners.
	RemoteAssistancePartners []RemoteAssistancePartner `json:"remoteAssistancePartners,omitempty"`
	Reports                  *DeviceManagementReports  `json:"reports,omitempty"`
	// The Resource Operations.
	ResourceOperations []ResourceOperation `json:"resourceOperations,omitempty"`
	// The Role Assignments.
	RoleAssignments []DeviceAndAppManagementRoleAssignment `json:"roleAssignments,omitempty"`
	// The Role Definitions.
	RoleDefinitions             []RoleDefinition                             `json:"roleDefinitions,omitempty"`
	SoftwareUpdateStatusSummary *DeviceManagementSoftwareUpdateStatusSummary `json:"softwareUpdateStatusSummary,omitempty"`
	// The telecom expense management partners.
	TelecomExpenseManagementPartners []TelecomExpenseManagementPartner `json:"telecomExpenseManagementPartners,omitempty"`
	// The terms and conditions associated with device management of the company.
	TermsAndConditions []TermsAndConditions `json:"termsAndConditions,omitempty"`
	// The list of troubleshooting events for the tenant.
	TroubleshootingEvents []DeviceManagementTroubleshootingEvent `json:"troubleshootingEvents,omitempty"`
	// User experience analytics appHealth Application Performance
	UserExperienceAnalyticsAppHealthApplicationPerformance []UserExperienceAnalyticsAppHealthApplicationPerformance `json:"userExperienceAnalyticsAppHealthApplicationPerformance,omitempty"`
	// User experience analytics appHealth Application Performance by App Version details
	UserExperienceAnalyticsAppHealthApplicationPerformanceByAppVersionDetails []UserExperienceAnalyticsAppHealthAppPerformanceByAppVersionDetails `json:"userExperienceAnalyticsAppHealthApplicationPerformanceByAppVersionDetails,omitempty"`
	// User experience analytics appHealth Application Performance by App Version Device Id
	UserExperienceAnalyticsAppHealthApplicationPerformanceByAppVersionDeviceId []UserExperienceAnalyticsAppHealthAppPerformanceByAppVersionDeviceId `json:"userExperienceAnalyticsAppHealthApplicationPerformanceByAppVersionDeviceId,omitempty"`
	// User experience analytics appHealth Application Performance by OS Version
	UserExperienceAnalyticsAppHealthApplicationPerformanceByOSVersion []UserExperienceAnalyticsAppHealthAppPerformanceByOSVersion `json:"userExperienceAnalyticsAppHealthApplicationPerformanceByOSVersion,omitempty"`
	// User experience analytics appHealth Model Performance
	UserExperienceAnalyticsAppHealthDeviceModelPerformance []UserExperienceAnalyticsAppHealthDeviceModelPerformance `json:"userExperienceAnalyticsAppHealthDeviceModelPerformance,omitempty"`
	// User experience analytics appHealth Device Performance
	UserExperienceAnalyticsAppHealthDevicePerformance []UserExperienceAnalyticsAppHealthDevicePerformance `json:"userExperienceAnalyticsAppHealthDevicePerformance,omitempty"`
	// User experience analytics device performance details
	UserExperienceAnalyticsAppHealthDevicePerformanceDetails []UserExperienceAnalyticsAppHealthDevicePerformanceDetails `json:"userExperienceAnalyticsAppHealthDevicePerformanceDetails,omitempty"`
	// User experience analytics appHealth OS version Performance
	UserExperienceAnalyticsAppHealthOSVersionPerformance []UserExperienceAnalyticsAppHealthOSVersionPerformance    `json:"userExperienceAnalyticsAppHealthOSVersionPerformance,omitempty"`
	UserExperienceAnalyticsAppHealthOverview             *DeviceManagementUserExperienceAnalyticsAppHealthOverview `json:"userExperienceAnalyticsAppHealthOverview,omitempty"`
	// User experience analytics baselines
	UserExperienceAnalyticsBaselines []UserExperienceAnalyticsBaseline `json:"userExperienceAnalyticsBaselines,omitempty"`
	// User experience analytics categories
	UserExperienceAnalyticsCategories []UserExperienceAnalyticsCategory `json:"userExperienceAnalyticsCategories,omitempty"`
	// User experience analytics device performance
	UserExperienceAnalyticsDevicePerformance []UserExperienceAnalyticsDevicePerformance `json:"userExperienceAnalyticsDevicePerformance,omitempty"`
	// User experience analytics device scores
	UserExperienceAnalyticsDeviceScores []UserExperienceAnalyticsDeviceScores `json:"userExperienceAnalyticsDeviceScores,omitempty"`
	// User experience analytics device Startup History
	UserExperienceAnalyticsDeviceStartupHistory []UserExperienceAnalyticsDeviceStartupHistory `json:"userExperienceAnalyticsDeviceStartupHistory,omitempty"`
	// User experience analytics device Startup Processes
	UserExperienceAnalyticsDeviceStartupProcesses []UserExperienceAnalyticsDeviceStartupProcess `json:"userExperienceAnalyticsDeviceStartupProcesses,omitempty"`
	// User experience analytics device Startup Process Performance
	UserExperienceAnalyticsDeviceStartupProcessPerformance []UserExperienceAnalyticsDeviceStartupProcessPerformance `json:"userExperienceAnalyticsDeviceStartupProcessPerformance,omitempty"`
	// User experience analytics metric history
	UserExperienceAnalyticsMetricHistory []UserExperienceAnalyticsMetricHistory `json:"userExperienceAnalyticsMetricHistory,omitempty"`
	// User experience analytics model scores
	UserExperienceAnalyticsModelScores []UserExperienceAnalyticsModelScores             `json:"userExperienceAnalyticsModelScores,omitempty"`
	UserExperienceAnalyticsOverview    *DeviceManagementUserExperienceAnalyticsOverview `json:"userExperienceAnalyticsOverview,omitempty"`
	// User experience analytics device Startup Score History
	UserExperienceAnalyticsScoreHistory                            []UserExperienceAnalyticsScoreHistory                                           `json:"userExperienceAnalyticsScoreHistory,omitempty"`
	UserExperienceAnalyticsWorkFromAnywhereHardwareReadinessMetric *DeviceManagementUserExperienceAnalyticsWorkFromAnywhereHardwareReadinessMetric `json:"userExperienceAnalyticsWorkFromAnywhereHardwareReadinessMetric,omitempty"`
	// User experience analytics work from anywhere metrics.
	UserExperienceAnalyticsWorkFromAnywhereMetrics []UserExperienceAnalyticsWorkFromAnywhereMetric `json:"userExperienceAnalyticsWorkFromAnywhereMetrics,omitempty"`
	// The user experience analytics work from anywhere model performance
	UserExperienceAnalyticsWorkFromAnywhereModelPerformance []UserExperienceAnalyticsWorkFromAnywhereModelPerformance `json:"userExperienceAnalyticsWorkFromAnywhereModelPerformance,omitempty"`
	VirtualEndpoint                                         *DeviceManagementVirtualEndpoint                          `json:"virtualEndpoint,omitempty"`
	// The Windows autopilot device identities contained collection.
	WindowsAutopilotDeviceIdentities []WindowsAutopilotDeviceIdentity `json:"windowsAutopilotDeviceIdentities,omitempty"`
	// The windows information protection app learning summaries.
	WindowsInformationProtectionAppLearningSummaries []WindowsInformationProtectionAppLearningSummary `json:"windowsInformationProtectionAppLearningSummaries,omitempty"`
	// The windows information protection network learning summaries.
	WindowsInformationProtectionNetworkLearningSummaries []WindowsInformationProtectionNetworkLearningSummary `json:"windowsInformationProtectionNetworkLearningSummaries,omitempty"`
	// The list of affected malware in the tenant.
	WindowsMalwareInformation []WindowsMalwareInformation `json:"windowsMalwareInformation,omitempty"`
	OdataType                 string                      `json:"@odata.type"`
}

// DetectedApp struct for DetectedApp
type DetectedApp struct {
	Entity
	// The number of devices that have installed this application
	DeviceCount *int32 `json:"deviceCount,omitempty"`
	// Name of the discovered application. Read-only
	DisplayName *string                  `json:"displayName,omitempty"`
	Platform    *DetectedAppPlatformType `json:"platform,omitempty"`
	// Indicates the publisher of the discovered application. For example: 'Microsoft'.  The default value is an empty string.
	Publisher *string `json:"publisher,omitempty"`
	// Discovered application size in bytes. Read-only
	SizeInByte *int64 `json:"sizeInByte,omitempty"`
	// Version of the discovered application. Read-only
	Version *string `json:"version,omitempty"`
	// The devices that have the discovered application installed
	ManagedDevices []ManagedDevice `json:"managedDevices,omitempty"`
	OdataType      string          `json:"@odata.type"`
}

// ManagedDevice struct for ManagedDevice
type ManagedDevice struct {
	Entity
	// The code that allows the Activation Lock on managed device to be bypassed. Default, is Null (Non-Default property) for this property when returned as part of managedDevice entity in LIST call. To retrieve actual values GET call needs to be made, with device id and included in select parameter. Supports: $select. $Search is not supported. Read-only. This property is read-only.
	ActivationLockBypassCode *string `json:"activationLockBypassCode,omitempty"`
	// Android security patch level. This property is read-only.
	AndroidSecurityPatchLevel *string `json:"androidSecurityPatchLevel,omitempty"`
	// The unique identifier for the Azure Active Directory device. Read only. This property is read-only.
	AzureADDeviceId *string `json:"azureADDeviceId,omitempty"`
	// Whether the device is Azure Active Directory registered. This property is read-only.
	AzureADRegistered *bool `json:"azureADRegistered,omitempty"`
	// The DateTime when device compliance grace period expires. This property is read-only.
	ComplianceGracePeriodExpirationDateTime   *time.Time                                 `json:"complianceGracePeriodExpirationDateTime,omitempty"`
	ComplianceState                           *ComplianceState                           `json:"complianceState,omitempty"`
	ConfigurationManagerClientEnabledFeatures *ConfigurationManagerClientEnabledFeatures `json:"configurationManagerClientEnabledFeatures,omitempty"`
	// List of ComplexType deviceActionResult objects. This property is read-only.
	DeviceActionResults []DeviceActionResult `json:"deviceActionResults,omitempty"`
	// Device category display name. Default is an empty string. Supports $filter operator 'eq' and 'or'. This property is read-only.
	DeviceCategoryDisplayName    *string                       `json:"deviceCategoryDisplayName,omitempty"`
	DeviceEnrollmentType         *DeviceEnrollmentType         `json:"deviceEnrollmentType,omitempty"`
	DeviceHealthAttestationState *DeviceHealthAttestationState `json:"deviceHealthAttestationState,omitempty"`
	// Name of the device. This property is read-only.
	DeviceName              *string                  `json:"deviceName,omitempty"`
	DeviceRegistrationState *DeviceRegistrationState `json:"deviceRegistrationState,omitempty"`
	// Whether the device is Exchange ActiveSync activated. This property is read-only.
	EasActivated *bool `json:"easActivated,omitempty"`
	// Exchange ActivationSync activation time of the device. This property is read-only.
	EasActivationDateTime *time.Time `json:"easActivationDateTime,omitempty"`
	// Exchange ActiveSync Id of the device. This property is read-only.
	EasDeviceId *string `json:"easDeviceId,omitempty"`
	// Email(s) for the user associated with the device. This property is read-only.
	EmailAddress *string `json:"emailAddress,omitempty"`
	// Enrollment time of the device. Supports $filter operator 'lt' and 'gt'. This property is read-only.
	EnrolledDateTime *time.Time `json:"enrolledDateTime,omitempty"`
	// Name of the enrollment profile assigned to the device. Default value is empty string, indicating no enrollment profile was assgined. This property is read-only.
	EnrollmentProfileName *string `json:"enrollmentProfileName,omitempty"`
	// Indicates Ethernet MAC Address of the device. Default, is Null (Non-Default property) for this property when returned as part of managedDevice entity. Individual get call with select query options is needed to retrieve actual values. Example: deviceManagement/managedDevices({managedDeviceId})?$select=ethernetMacAddress Supports: $select. $Search is not supported. Read-only. This property is read-only.
	EthernetMacAddress        *string                                    `json:"ethernetMacAddress,omitempty"`
	ExchangeAccessState       *DeviceManagementExchangeAccessState       `json:"exchangeAccessState,omitempty"`
	ExchangeAccessStateReason *DeviceManagementExchangeAccessStateReason `json:"exchangeAccessStateReason,omitempty"`
	// Last time the device contacted Exchange. This property is read-only.
	ExchangeLastSuccessfulSyncDateTime *time.Time `json:"exchangeLastSuccessfulSyncDateTime,omitempty"`
	// Free Storage in Bytes. Default value is 0. Read-only. This property is read-only.
	FreeStorageSpaceInBytes *int64 `json:"freeStorageSpaceInBytes,omitempty"`
	// Integrated Circuit Card Identifier, it is A SIM card's unique identification number. Default is an empty string. To retrieve actual values GET call needs to be made, with device id and included in select parameter. Supports: $select. $Search is not supported. Read-only. This property is read-only.
	Iccid *string `json:"iccid,omitempty"`
	// IMEI. This property is read-only.
	Imei *string `json:"imei,omitempty"`
	// Device encryption status. This property is read-only.
	IsEncrypted *bool `json:"isEncrypted,omitempty"`
	// Device supervised status. This property is read-only.
	IsSupervised *bool `json:"isSupervised,omitempty"`
	// Whether the device is jail broken or rooted. Default is an empty string. Supports $filter operator 'eq' and 'or'. This property is read-only.
	JailBroken *string `json:"jailBroken,omitempty"`
	// The date and time that the device last completed a successful sync with Intune. Supports $filter operator 'lt' and 'gt'. This property is read-only.
	LastSyncDateTime *time.Time `json:"lastSyncDateTime,omitempty"`
	// Automatically generated name to identify a device. Can be overwritten to a user friendly name.
	ManagedDeviceName      *string                 `json:"managedDeviceName,omitempty"`
	ManagedDeviceOwnerType *ManagedDeviceOwnerType `json:"managedDeviceOwnerType,omitempty"`
	ManagementAgent        *ManagementAgentType    `json:"managementAgent,omitempty"`
	// Reports device management certificate expiration date. This property is read-only.
	ManagementCertificateExpirationDate *time.Time `json:"managementCertificateExpirationDate,omitempty"`
	// Manufacturer of the device. This property is read-only.
	Manufacturer *string `json:"manufacturer,omitempty"`
	// MEID. This property is read-only.
	Meid *string `json:"meid,omitempty"`
	// Model of the device. This property is read-only.
	Model *string `json:"model,omitempty"`
	// Notes on the device created by IT Admin. Default is null. To retrieve actual values GET call needs to be made, with device id and included in select parameter. Supports: $select. $Search is not supported.
	Notes *string `json:"notes,omitempty"`
	// Operating system of the device. Windows, iOS, etc. This property is read-only.
	OperatingSystem *string `json:"operatingSystem,omitempty"`
	// Operating system version of the device. This property is read-only.
	OsVersion                  *string                                  `json:"osVersion,omitempty"`
	PartnerReportedThreatState *ManagedDevicePartnerReportedHealthState `json:"partnerReportedThreatState,omitempty"`
	// Phone number of the device. This property is read-only.
	PhoneNumber *string `json:"phoneNumber,omitempty"`
	// Total Memory in Bytes. Default is 0. To retrieve actual values GET call needs to be made, with device id and included in select parameter. Supports: $select. Read-only. This property is read-only.
	PhysicalMemoryInBytes *int64 `json:"physicalMemoryInBytes,omitempty"`
	// An error string that identifies issues when creating Remote Assistance session objects. This property is read-only.
	RemoteAssistanceSessionErrorDetails *string `json:"remoteAssistanceSessionErrorDetails,omitempty"`
	// Url that allows a Remote Assistance session to be established with the device. Default is an empty string. To retrieve actual values GET call needs to be made, with device id and included in select parameter. This property is read-only.
	RemoteAssistanceSessionUrl *string `json:"remoteAssistanceSessionUrl,omitempty"`
	// Reports if the managed iOS device is user approval enrollment. This property is read-only.
	RequireUserEnrollmentApproval *bool `json:"requireUserEnrollmentApproval,omitempty"`
	// SerialNumber. This property is read-only.
	SerialNumber *string `json:"serialNumber,omitempty"`
	// Subscriber Carrier. This property is read-only.
	SubscriberCarrier *string `json:"subscriberCarrier,omitempty"`
	// Total Storage in Bytes. This property is read-only.
	TotalStorageSpaceInBytes *int64 `json:"totalStorageSpaceInBytes,omitempty"`
	// Unique Device Identifier for iOS and macOS devices. Default is an empty string. To retrieve actual values GET call needs to be made, with device id and included in select parameter. Supports: $select. $Search is not supported. Read-only. This property is read-only.
	Udid *string `json:"udid,omitempty"`
	// User display name. This property is read-only.
	UserDisplayName *string `json:"userDisplayName,omitempty"`
	// Unique Identifier for the user associated with the device. This property is read-only.
	UserId *string `json:"userId,omitempty"`
	// Device user principal name. This property is read-only.
	UserPrincipalName *string `json:"userPrincipalName,omitempty"`
	// Wi-Fi MAC. This property is read-only.
	WiFiMacAddress *string         `json:"wiFiMacAddress,omitempty"`
	DeviceCategory *DeviceCategory `json:"deviceCategory,omitempty"`
	// Device compliance policy states for this device.
	DeviceCompliancePolicyStates []DeviceCompliancePolicyState `json:"deviceCompliancePolicyStates,omitempty"`
	// Device configuration states for this device.
	DeviceConfigurationStates []DeviceConfigurationState `json:"deviceConfigurationStates,omitempty"`
	// List of log collection requests
	LogCollectionRequests []DeviceLogCollectionResponse `json:"logCollectionRequests,omitempty"`
	// The primary users associated with the managed device.
	Users                  []User                  `json:"users,omitempty"`
	WindowsProtectionState *WindowsProtectionState `json:"windowsProtectionState,omitempty"`
	OdataType              string                  `json:"@odata.type"`
}

// DeviceLogCollectionResponse struct for DeviceLogCollectionResponse
type DeviceLogCollectionResponse struct {
	Entity
	// The User Principal Name (UPN) of the user that enrolled the device.
	EnrolledByUser *string `json:"enrolledByUser,omitempty"`
	// The DateTime of the expiration of the logs.
	ExpirationDateTimeUTC *time.Time `json:"expirationDateTimeUTC,omitempty"`
	// The UPN for who initiated the request.
	InitiatedByUserPrincipalName *string `json:"initiatedByUserPrincipalName,omitempty"`
	// Indicates Intune device unique identifier.
	ManagedDeviceId *string `json:"managedDeviceId,omitempty"`
	// The DateTime the request was received.
	ReceivedDateTimeUTC *time.Time `json:"receivedDateTimeUTC,omitempty"`
	// The DateTime of the request.
	RequestedDateTimeUTC *time.Time                           `json:"requestedDateTimeUTC,omitempty"`
	SizeInKB             *DeviceLogCollectionResponseSizeInKB `json:"sizeInKB,omitempty"`
	Status               *AppLogUploadState                   `json:"status,omitempty"`
	OdataType            string                               `json:"@odata.type"`
}

// DeviceCategory struct for DeviceCategory
type DeviceCategory struct {
	Entity
	// Optional description for the device category.
	Description *string `json:"description,omitempty"`
	// Display name for the device category.
	DisplayName *string `json:"displayName,omitempty"`
	OdataType   string  `json:"@odata.type"`
}

// DeviceConfigurationSettingState Device Configuration Setting State for a given device.
type DeviceConfigurationSettingState struct {
	// Current value of setting on device
	CurrentValue *string `json:"currentValue,omitempty"`
	// Error code for the setting
	ErrorCode *int64 `json:"errorCode,omitempty"`
	// Error description
	ErrorDescription *string `json:"errorDescription,omitempty"`
	// Name of setting instance that is being reported.
	InstanceDisplayName *string `json:"instanceDisplayName,omitempty"`
	// The setting that is being reported
	Setting *string `json:"setting,omitempty"`
	// Localized/user friendly setting name that is being reported
	SettingName *string `json:"settingName,omitempty"`
	// Contributing policies
	Sources []SettingSource   `json:"sources,omitempty"`
	State   *ComplianceStatus `json:"state,omitempty"`
	// UserEmail
	UserEmail *string `json:"userEmail,omitempty"`
	// UserId
	UserId *string `json:"userId,omitempty"`
	// UserName
	UserName *string `json:"userName,omitempty"`
	// UserPrincipalName.
	UserPrincipalName *string `json:"userPrincipalName,omitempty"`
	OdataType         string  `json:"@odata.type"`
}

// DeviceCompliancePolicyState struct for DeviceCompliancePolicyState
type DeviceCompliancePolicyState struct {
	Entity
	// The name of the policy for this policyBase
	DisplayName  *string             `json:"displayName,omitempty"`
	PlatformType *PolicyPlatformType `json:"platformType,omitempty"`
	// Count of how many setting a policy holds
	SettingCount  *int32                               `json:"settingCount,omitempty"`
	SettingStates []DeviceCompliancePolicySettingState `json:"settingStates,omitempty"`
	State         *ComplianceStatus                    `json:"state,omitempty"`
	// The version of the policy
	Version   *int32 `json:"version,omitempty"`
	OdataType string `json:"@odata.type"`
}

// DeviceCompliancePolicySettingState Device Compilance Policy Setting State for a given device.
type DeviceCompliancePolicySettingState struct {
	// Current value of setting on device
	CurrentValue *string `json:"currentValue,omitempty"`
	// Error code for the setting
	ErrorCode *int64 `json:"errorCode,omitempty"`
	// Error description
	ErrorDescription *string `json:"errorDescription,omitempty"`
	// Name of setting instance that is being reported.
	InstanceDisplayName *string `json:"instanceDisplayName,omitempty"`
	// The setting that is being reported
	Setting *string `json:"setting,omitempty"`
	// Localized/user friendly setting name that is being reported
	SettingName *string `json:"settingName,omitempty"`
	// Contributing policies
	Sources []SettingSource   `json:"sources,omitempty"`
	State   *ComplianceStatus `json:"state,omitempty"`
	// UserEmail
	UserEmail *string `json:"userEmail,omitempty"`
	// UserId
	UserId *string `json:"userId,omitempty"`
	// UserName
	UserName *string `json:"userName,omitempty"`
	// UserPrincipalName.
	UserPrincipalName *string `json:"userPrincipalName,omitempty"`
	OdataType         string  `json:"@odata.type"`
}

// DeviceConfigurationState struct for DeviceConfigurationState
type DeviceConfigurationState struct {
	Entity
	// The name of the policy for this policyBase
	DisplayName  *string             `json:"displayName,omitempty"`
	PlatformType *PolicyPlatformType `json:"platformType,omitempty"`
	// Count of how many setting a policy holds
	SettingCount  *int32                            `json:"settingCount,omitempty"`
	SettingStates []DeviceConfigurationSettingState `json:"settingStates,omitempty"`
	State         *ComplianceStatus                 `json:"state,omitempty"`
	// The version of the policy
	Version   *int32 `json:"version,omitempty"`
	OdataType string `json:"@odata.type"`
}

// ComplianceStatus the model 'ComplianceStatus'
type ComplianceStatus string

// List of microsoft.graph.complianceStatus
const (
	COMPLIANCESTATUS_UNKNOWN        ComplianceStatus = "unknown"
	COMPLIANCESTATUS_NOT_APPLICABLE ComplianceStatus = "notApplicable"
	COMPLIANCESTATUS_COMPLIANT      ComplianceStatus = "compliant"
	COMPLIANCESTATUS_REMEDIATED     ComplianceStatus = "remediated"
	COMPLIANCESTATUS_NON_COMPLIANT  ComplianceStatus = "nonCompliant"
	COMPLIANCESTATUS_ERROR          ComplianceStatus = "error"
	COMPLIANCESTATUS_CONFLICT       ComplianceStatus = "conflict"
	COMPLIANCESTATUS_NOT_ASSIGNED   ComplianceStatus = "notAssigned"
)

// All allowed values of ComplianceStatus enum
var AllowedComplianceStatusEnumValues = []ComplianceStatus{
	"unknown",
	"notApplicable",
	"compliant",
	"remediated",
	"nonCompliant",
	"error",
	"conflict",
	"notAssigned",
}

func (v *ComplianceStatus) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := ComplianceStatus(value)
	for _, existing := range AllowedComplianceStatusEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid ComplianceStatus", value)
}

// SettingSource struct for SettingSource
type SettingSource struct {
	// Not yet documented
	DisplayName *string `json:"displayName,omitempty"`
	// Not yet documented
	Id         *string            `json:"id,omitempty"`
	SourceType *SettingSourceType `json:"sourceType,omitempty"`
	OdataType  string             `json:"@odata.type"`
}

// SettingSourceType the model 'SettingSourceType'
type SettingSourceType string

// List of microsoft.graph.settingSourceType
const (
	SETTINGSOURCETYPE_DEVICE_CONFIGURATION SettingSourceType = "deviceConfiguration"
	SETTINGSOURCETYPE_DEVICE_INTENT        SettingSourceType = "deviceIntent"
)

// All allowed values of SettingSourceType enum
var AllowedSettingSourceTypeEnumValues = []SettingSourceType{
	"deviceConfiguration",
	"deviceIntent",
}

func (v *SettingSourceType) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := SettingSourceType(value)
	for _, existing := range AllowedSettingSourceTypeEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid SettingSourceType", value)
}

// PolicyPlatformType Supported platform types for policies.
type PolicyPlatformType string

// List of microsoft.graph.policyPlatformType
const (
	POLICYPLATFORMTYPE_ANDROID             PolicyPlatformType = "android"
	POLICYPLATFORMTYPE_ANDROID_FOR_WORK    PolicyPlatformType = "androidForWork"
	POLICYPLATFORMTYPE_I_OS                PolicyPlatformType = "iOS"
	POLICYPLATFORMTYPE_MAC_OS              PolicyPlatformType = "macOS"
	POLICYPLATFORMTYPE_WINDOWS_PHONE81     PolicyPlatformType = "windowsPhone81"
	POLICYPLATFORMTYPE_WINDOWS81_AND_LATER PolicyPlatformType = "windows81AndLater"
	POLICYPLATFORMTYPE_WINDOWS10_AND_LATER PolicyPlatformType = "windows10AndLater"
	POLICYPLATFORMTYPE_ALL                 PolicyPlatformType = "all"
)

// All allowed values of PolicyPlatformType enum
var AllowedPolicyPlatformTypeEnumValues = []PolicyPlatformType{
	"android",
	"androidForWork",
	"iOS",
	"macOS",
	"windowsPhone81",
	"windows81AndLater",
	"windows10AndLater",
	"all",
}

func (v *PolicyPlatformType) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := PolicyPlatformType(value)
	for _, existing := range AllowedPolicyPlatformTypeEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid PolicyPlatformType", value)
}

// ManagedDevicePartnerReportedHealthState Available health states for the Device Health API
type ManagedDevicePartnerReportedHealthState string

// List of microsoft.graph.managedDevicePartnerReportedHealthState
const (
	MANAGEDDEVICEPARTNERREPORTEDHEALTHSTATE_UNKNOWN         ManagedDevicePartnerReportedHealthState = "unknown"
	MANAGEDDEVICEPARTNERREPORTEDHEALTHSTATE_ACTIVATED       ManagedDevicePartnerReportedHealthState = "activated"
	MANAGEDDEVICEPARTNERREPORTEDHEALTHSTATE_DEACTIVATED     ManagedDevicePartnerReportedHealthState = "deactivated"
	MANAGEDDEVICEPARTNERREPORTEDHEALTHSTATE_SECURED         ManagedDevicePartnerReportedHealthState = "secured"
	MANAGEDDEVICEPARTNERREPORTEDHEALTHSTATE_LOW_SEVERITY    ManagedDevicePartnerReportedHealthState = "lowSeverity"
	MANAGEDDEVICEPARTNERREPORTEDHEALTHSTATE_MEDIUM_SEVERITY ManagedDevicePartnerReportedHealthState = "mediumSeverity"
	MANAGEDDEVICEPARTNERREPORTEDHEALTHSTATE_HIGH_SEVERITY   ManagedDevicePartnerReportedHealthState = "highSeverity"
	MANAGEDDEVICEPARTNERREPORTEDHEALTHSTATE_UNRESPONSIVE    ManagedDevicePartnerReportedHealthState = "unresponsive"
	MANAGEDDEVICEPARTNERREPORTEDHEALTHSTATE_COMPROMISED     ManagedDevicePartnerReportedHealthState = "compromised"
	MANAGEDDEVICEPARTNERREPORTEDHEALTHSTATE_MISCONFIGURED   ManagedDevicePartnerReportedHealthState = "misconfigured"
)

// All allowed values of ManagedDevicePartnerReportedHealthState enum
var AllowedManagedDevicePartnerReportedHealthStateEnumValues = []ManagedDevicePartnerReportedHealthState{
	"unknown",
	"activated",
	"deactivated",
	"secured",
	"lowSeverity",
	"mediumSeverity",
	"highSeverity",
	"unresponsive",
	"compromised",
	"misconfigured",
}

func (v *ManagedDevicePartnerReportedHealthState) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := ManagedDevicePartnerReportedHealthState(value)
	for _, existing := range AllowedManagedDevicePartnerReportedHealthStateEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid ManagedDevicePartnerReportedHealthState", value)
}

// ManagementAgentType the model 'ManagementAgentType'
type ManagementAgentType string

// List of microsoft.graph.managementAgentType
const (
	MANAGEMENTAGENTTYPE_EAS                                   ManagementAgentType = "eas"
	MANAGEMENTAGENTTYPE_MDM                                   ManagementAgentType = "mdm"
	MANAGEMENTAGENTTYPE_EAS_MDM                               ManagementAgentType = "easMdm"
	MANAGEMENTAGENTTYPE_INTUNE_CLIENT                         ManagementAgentType = "intuneClient"
	MANAGEMENTAGENTTYPE_EAS_INTUNE_CLIENT                     ManagementAgentType = "easIntuneClient"
	MANAGEMENTAGENTTYPE_CONFIGURATION_MANAGER_CLIENT          ManagementAgentType = "configurationManagerClient"
	MANAGEMENTAGENTTYPE_CONFIGURATION_MANAGER_CLIENT_MDM      ManagementAgentType = "configurationManagerClientMdm"
	MANAGEMENTAGENTTYPE_CONFIGURATION_MANAGER_CLIENT_MDM_EAS  ManagementAgentType = "configurationManagerClientMdmEas"
	MANAGEMENTAGENTTYPE_UNKNOWN                               ManagementAgentType = "unknown"
	MANAGEMENTAGENTTYPE_JAMF                                  ManagementAgentType = "jamf"
	MANAGEMENTAGENTTYPE_GOOGLE_CLOUD_DEVICE_POLICY_CONTROLLER ManagementAgentType = "googleCloudDevicePolicyController"
	MANAGEMENTAGENTTYPE_MICROSOFT365_MANAGED_MDM              ManagementAgentType = "microsoft365ManagedMdm"
	MANAGEMENTAGENTTYPE_MS_SENSE                              ManagementAgentType = "msSense"
)

// All allowed values of ManagementAgentType enum
var AllowedManagementAgentTypeEnumValues = []ManagementAgentType{
	"eas",
	"mdm",
	"easMdm",
	"intuneClient",
	"easIntuneClient",
	"configurationManagerClient",
	"configurationManagerClientMdm",
	"configurationManagerClientMdmEas",
	"unknown",
	"jamf",
	"googleCloudDevicePolicyController",
	"microsoft365ManagedMdm",
	"msSense",
}

func (v *ManagementAgentType) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := ManagementAgentType(value)
	for _, existing := range AllowedManagementAgentTypeEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid ManagementAgentType", value)
}

// ManagedDeviceOwnerType Owner type of device.
type ManagedDeviceOwnerType string

// List of microsoft.graph.managedDeviceOwnerType
const (
	MANAGEDDEVICEOWNERTYPE_UNKNOWN  ManagedDeviceOwnerType = "unknown"
	MANAGEDDEVICEOWNERTYPE_COMPANY  ManagedDeviceOwnerType = "company"
	MANAGEDDEVICEOWNERTYPE_PERSONAL ManagedDeviceOwnerType = "personal"
)

// All allowed values of ManagedDeviceOwnerType enum
var AllowedManagedDeviceOwnerTypeEnumValues = []ManagedDeviceOwnerType{
	"unknown",
	"company",
	"personal",
}

func (v *ManagedDeviceOwnerType) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := ManagedDeviceOwnerType(value)
	for _, existing := range AllowedManagedDeviceOwnerTypeEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid ManagedDeviceOwnerType", value)
}

// DeviceManagementExchangeAccessStateReason Device Exchange Access State Reason.
type DeviceManagementExchangeAccessStateReason string

// List of microsoft.graph.deviceManagementExchangeAccessStateReason
const (
	DEVICEMANAGEMENTEXCHANGEACCESSSTATEREASON_NONE                                DeviceManagementExchangeAccessStateReason = "none"
	DEVICEMANAGEMENTEXCHANGEACCESSSTATEREASON_UNKNOWN                             DeviceManagementExchangeAccessStateReason = "unknown"
	DEVICEMANAGEMENTEXCHANGEACCESSSTATEREASON_EXCHANGE_GLOBAL_RULE                DeviceManagementExchangeAccessStateReason = "exchangeGlobalRule"
	DEVICEMANAGEMENTEXCHANGEACCESSSTATEREASON_EXCHANGE_INDIVIDUAL_RULE            DeviceManagementExchangeAccessStateReason = "exchangeIndividualRule"
	DEVICEMANAGEMENTEXCHANGEACCESSSTATEREASON_EXCHANGE_DEVICE_RULE                DeviceManagementExchangeAccessStateReason = "exchangeDeviceRule"
	DEVICEMANAGEMENTEXCHANGEACCESSSTATEREASON_EXCHANGE_UPGRADE                    DeviceManagementExchangeAccessStateReason = "exchangeUpgrade"
	DEVICEMANAGEMENTEXCHANGEACCESSSTATEREASON_EXCHANGE_MAILBOX_POLICY             DeviceManagementExchangeAccessStateReason = "exchangeMailboxPolicy"
	DEVICEMANAGEMENTEXCHANGEACCESSSTATEREASON_OTHER                               DeviceManagementExchangeAccessStateReason = "other"
	DEVICEMANAGEMENTEXCHANGEACCESSSTATEREASON_COMPLIANT                           DeviceManagementExchangeAccessStateReason = "compliant"
	DEVICEMANAGEMENTEXCHANGEACCESSSTATEREASON_NOT_COMPLIANT                       DeviceManagementExchangeAccessStateReason = "notCompliant"
	DEVICEMANAGEMENTEXCHANGEACCESSSTATEREASON_NOT_ENROLLED                        DeviceManagementExchangeAccessStateReason = "notEnrolled"
	DEVICEMANAGEMENTEXCHANGEACCESSSTATEREASON_UNKNOWN_LOCATION                    DeviceManagementExchangeAccessStateReason = "unknownLocation"
	DEVICEMANAGEMENTEXCHANGEACCESSSTATEREASON_MFA_REQUIRED                        DeviceManagementExchangeAccessStateReason = "mfaRequired"
	DEVICEMANAGEMENTEXCHANGEACCESSSTATEREASON_AZURE_AD_BLOCK_DUE_TO_ACCESS_POLICY DeviceManagementExchangeAccessStateReason = "azureADBlockDueToAccessPolicy"
	DEVICEMANAGEMENTEXCHANGEACCESSSTATEREASON_COMPROMISED_PASSWORD                DeviceManagementExchangeAccessStateReason = "compromisedPassword"
	DEVICEMANAGEMENTEXCHANGEACCESSSTATEREASON_DEVICE_NOT_KNOWN_WITH_MANAGED_APP   DeviceManagementExchangeAccessStateReason = "deviceNotKnownWithManagedApp"
)

// All allowed values of DeviceManagementExchangeAccessStateReason enum
var AllowedDeviceManagementExchangeAccessStateReasonEnumValues = []DeviceManagementExchangeAccessStateReason{
	"none",
	"unknown",
	"exchangeGlobalRule",
	"exchangeIndividualRule",
	"exchangeDeviceRule",
	"exchangeUpgrade",
	"exchangeMailboxPolicy",
	"other",
	"compliant",
	"notCompliant",
	"notEnrolled",
	"unknownLocation",
	"mfaRequired",
	"azureADBlockDueToAccessPolicy",
	"compromisedPassword",
	"deviceNotKnownWithManagedApp",
}

func (v *DeviceManagementExchangeAccessStateReason) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := DeviceManagementExchangeAccessStateReason(value)
	for _, existing := range AllowedDeviceManagementExchangeAccessStateReasonEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid DeviceManagementExchangeAccessStateReason", value)
}

// DeviceManagementExchangeAccessState Device Exchange Access State.
type DeviceManagementExchangeAccessState string

// List of microsoft.graph.deviceManagementExchangeAccessState
const (
	DEVICEMANAGEMENTEXCHANGEACCESSSTATE_NONE        DeviceManagementExchangeAccessState = "none"
	DEVICEMANAGEMENTEXCHANGEACCESSSTATE_UNKNOWN     DeviceManagementExchangeAccessState = "unknown"
	DEVICEMANAGEMENTEXCHANGEACCESSSTATE_ALLOWED     DeviceManagementExchangeAccessState = "allowed"
	DEVICEMANAGEMENTEXCHANGEACCESSSTATE_BLOCKED     DeviceManagementExchangeAccessState = "blocked"
	DEVICEMANAGEMENTEXCHANGEACCESSSTATE_QUARANTINED DeviceManagementExchangeAccessState = "quarantined"
)

// All allowed values of DeviceManagementExchangeAccessState enum
var AllowedDeviceManagementExchangeAccessStateEnumValues = []DeviceManagementExchangeAccessState{
	"none",
	"unknown",
	"allowed",
	"blocked",
	"quarantined",
}

func (v *DeviceManagementExchangeAccessState) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := DeviceManagementExchangeAccessState(value)
	for _, existing := range AllowedDeviceManagementExchangeAccessStateEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid DeviceManagementExchangeAccessState", value)
}

// DeviceRegistrationState Device registration status.
type DeviceRegistrationState string

// List of microsoft.graph.deviceRegistrationState
const (
	DEVICEREGISTRATIONSTATE_NOT_REGISTERED                    DeviceRegistrationState = "notRegistered"
	DEVICEREGISTRATIONSTATE_REGISTERED                        DeviceRegistrationState = "registered"
	DEVICEREGISTRATIONSTATE_REVOKED                           DeviceRegistrationState = "revoked"
	DEVICEREGISTRATIONSTATE_KEY_CONFLICT                      DeviceRegistrationState = "keyConflict"
	DEVICEREGISTRATIONSTATE_APPROVAL_PENDING                  DeviceRegistrationState = "approvalPending"
	DEVICEREGISTRATIONSTATE_CERTIFICATE_RESET                 DeviceRegistrationState = "certificateReset"
	DEVICEREGISTRATIONSTATE_NOT_REGISTERED_PENDING_ENROLLMENT DeviceRegistrationState = "notRegisteredPendingEnrollment"
	DEVICEREGISTRATIONSTATE_UNKNOWN                           DeviceRegistrationState = "unknown"
)

// All allowed values of DeviceRegistrationState enum
var AllowedDeviceRegistrationStateEnumValues = []DeviceRegistrationState{
	"notRegistered",
	"registered",
	"revoked",
	"keyConflict",
	"approvalPending",
	"certificateReset",
	"notRegisteredPendingEnrollment",
	"unknown",
}

func (v *DeviceRegistrationState) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := DeviceRegistrationState(value)
	for _, existing := range AllowedDeviceRegistrationStateEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid DeviceRegistrationState", value)
}

// DeviceHealthAttestationState struct for DeviceHealthAttestationState
type DeviceHealthAttestationState struct {
	// TWhen an Attestation Identity Key (AIK) is present on a device, it indicates that the device has an endorsement key (EK) certificate.
	AttestationIdentityKey *string `json:"attestationIdentityKey,omitempty"`
	// On or Off of BitLocker Drive Encryption
	BitLockerStatus *string `json:"bitLockerStatus,omitempty"`
	// The security version number of the Boot Application
	BootAppSecurityVersion *string `json:"bootAppSecurityVersion,omitempty"`
	// When bootDebugging is enabled, the device is used in development and testing
	BootDebugging *string `json:"bootDebugging,omitempty"`
	// The security version number of the Boot Application
	BootManagerSecurityVersion *string `json:"bootManagerSecurityVersion,omitempty"`
	// The version of the Boot Manager
	BootManagerVersion *string `json:"bootManagerVersion,omitempty"`
	// The Boot Revision List that was loaded during initial boot on the attested device
	BootRevisionListInfo *string `json:"bootRevisionListInfo,omitempty"`
	// When code integrity is enabled, code execution is restricted to integrity verified code
	CodeIntegrity *string `json:"codeIntegrity,omitempty"`
	// The version of the Boot Manager
	CodeIntegrityCheckVersion *string `json:"codeIntegrityCheckVersion,omitempty"`
	// The Code Integrity policy that is controlling the security of the boot environment
	CodeIntegrityPolicy *string `json:"codeIntegrityPolicy,omitempty"`
	// The DHA report version. (Namespace version)
	ContentNamespaceUrl *string `json:"contentNamespaceUrl,omitempty"`
	// The HealthAttestation state schema version
	ContentVersion *string `json:"contentVersion,omitempty"`
	// DEP Policy defines a set of hardware and software technologies that perform additional checks on memory
	DataExcutionPolicy *string `json:"dataExcutionPolicy,omitempty"`
	// The DHA report version. (Namespace version)
	DeviceHealthAttestationStatus *string `json:"deviceHealthAttestationStatus,omitempty"`
	// ELAM provides protection for the computers in your network when they start up
	EarlyLaunchAntiMalwareDriverProtection *string `json:"earlyLaunchAntiMalwareDriverProtection,omitempty"`
	// This attribute indicates if DHA is supported for the device
	HealthAttestationSupportedStatus *string `json:"healthAttestationSupportedStatus,omitempty"`
	// This attribute appears if DHA-Service detects an integrity issue
	HealthStatusMismatchInfo *string `json:"healthStatusMismatchInfo,omitempty"`
	// The DateTime when device was evaluated or issued to MDM
	IssuedDateTime *time.Time `json:"issuedDateTime,omitempty"`
	// The Timestamp of the last update.
	LastUpdateDateTime *string `json:"lastUpdateDateTime,omitempty"`
	// When operatingSystemKernelDebugging is enabled, the device is used in development and testing
	OperatingSystemKernelDebugging *string `json:"operatingSystemKernelDebugging,omitempty"`
	// The Operating System Revision List that was loaded during initial boot on the attested device
	OperatingSystemRevListInfo *string `json:"operatingSystemRevListInfo,omitempty"`
	// The measurement that is captured in PCR[0]
	Pcr0 *string `json:"pcr0,omitempty"`
	// Informational attribute that identifies the HASH algorithm that was used by TPM
	PcrHashAlgorithm *string `json:"pcrHashAlgorithm,omitempty"`
	// The number of times a PC device has hibernated or resumed
	ResetCount *int64 `json:"resetCount,omitempty"`
	// The number of times a PC device has rebooted
	RestartCount *int64 `json:"restartCount,omitempty"`
	// Safe mode is a troubleshooting option for Windows that starts your computer in a limited state
	SafeMode *string `json:"safeMode,omitempty"`
	// When Secure Boot is enabled, the core components must have the correct cryptographic signatures
	SecureBoot *string `json:"secureBoot,omitempty"`
	// Fingerprint of the Custom Secure Boot Configuration Policy
	SecureBootConfigurationPolicyFingerPrint *string `json:"secureBootConfigurationPolicyFingerPrint,omitempty"`
	// When test signing is allowed, the device does not enforce signature validation during boot
	TestSigning *string `json:"testSigning,omitempty"`
	// The security version number of the Boot Application
	TpmVersion *string `json:"tpmVersion,omitempty"`
	// VSM is a container that protects high value assets from a compromised kernel
	VirtualSecureMode *string `json:"virtualSecureMode,omitempty"`
	// Operating system running with limited services that is used to prepare a computer for Windows
	WindowsPE *string `json:"windowsPE,omitempty"`
	OdataType string  `json:"@odata.type"`
}

// DeviceEnrollmentType Possible ways of adding a mobile device to management.
type DeviceEnrollmentType string

// List of microsoft.graph.deviceEnrollmentType
const (
	DEVICEENROLLMENTTYPE_UNKNOWN                                    DeviceEnrollmentType = "unknown"
	DEVICEENROLLMENTTYPE_USER_ENROLLMENT                            DeviceEnrollmentType = "userEnrollment"
	DEVICEENROLLMENTTYPE_DEVICE_ENROLLMENT_MANAGER                  DeviceEnrollmentType = "deviceEnrollmentManager"
	DEVICEENROLLMENTTYPE_APPLE_BULK_WITH_USER                       DeviceEnrollmentType = "appleBulkWithUser"
	DEVICEENROLLMENTTYPE_APPLE_BULK_WITHOUT_USER                    DeviceEnrollmentType = "appleBulkWithoutUser"
	DEVICEENROLLMENTTYPE_WINDOWS_AZURE_AD_JOIN                      DeviceEnrollmentType = "windowsAzureADJoin"
	DEVICEENROLLMENTTYPE_WINDOWS_BULK_USERLESS                      DeviceEnrollmentType = "windowsBulkUserless"
	DEVICEENROLLMENTTYPE_WINDOWS_AUTO_ENROLLMENT                    DeviceEnrollmentType = "windowsAutoEnrollment"
	DEVICEENROLLMENTTYPE_WINDOWS_BULK_AZURE_DOMAIN_JOIN             DeviceEnrollmentType = "windowsBulkAzureDomainJoin"
	DEVICEENROLLMENTTYPE_WINDOWS_CO_MANAGEMENT                      DeviceEnrollmentType = "windowsCoManagement"
	DEVICEENROLLMENTTYPE_WINDOWS_AZURE_AD_JOIN_USING_DEVICE_AUTH    DeviceEnrollmentType = "windowsAzureADJoinUsingDeviceAuth"
	DEVICEENROLLMENTTYPE_APPLE_USER_ENROLLMENT                      DeviceEnrollmentType = "appleUserEnrollment"
	DEVICEENROLLMENTTYPE_APPLE_USER_ENROLLMENT_WITH_SERVICE_ACCOUNT DeviceEnrollmentType = "appleUserEnrollmentWithServiceAccount"
)

// All allowed values of DeviceEnrollmentType enum
var AllowedDeviceEnrollmentTypeEnumValues = []DeviceEnrollmentType{
	"unknown",
	"userEnrollment",
	"deviceEnrollmentManager",
	"appleBulkWithUser",
	"appleBulkWithoutUser",
	"windowsAzureADJoin",
	"windowsBulkUserless",
	"windowsAutoEnrollment",
	"windowsBulkAzureDomainJoin",
	"windowsCoManagement",
	"windowsAzureADJoinUsingDeviceAuth",
	"appleUserEnrollment",
	"appleUserEnrollmentWithServiceAccount",
}

func (v *DeviceEnrollmentType) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := DeviceEnrollmentType(value)
	for _, existing := range AllowedDeviceEnrollmentTypeEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid DeviceEnrollmentType", value)
}

// ActionResult Device action result
type DeviceActionResult struct {
	// Action name
	ActionName  *string      `json:"actionName,omitempty"`
	ActionState *ActionState `json:"actionState,omitempty"`
	// Time the action state was last updated
	LastUpdatedDateTime *time.Time `json:"lastUpdatedDateTime,omitempty"`
	// Time the action was initiated
	StartDateTime *time.Time `json:"startDateTime,omitempty"`
	OdataType     string     `json:"@odata.type"`
}

// ActionState State of the action on the device
type ActionState string

// List of microsoft.graph.actionState
const (
	ACTIONSTATE_NONE          ActionState = "none"
	ACTIONSTATE_PENDING       ActionState = "pending"
	ACTIONSTATE_CANCELED      ActionState = "canceled"
	ACTIONSTATE_ACTIVE        ActionState = "active"
	ACTIONSTATE_DONE          ActionState = "done"
	ACTIONSTATE_FAILED        ActionState = "failed"
	ACTIONSTATE_NOT_SUPPORTED ActionState = "notSupported"
)

// All allowed values of ActionState enum
var AllowedActionStateEnumValues = []ActionState{
	"none",
	"pending",
	"canceled",
	"active",
	"done",
	"failed",
	"notSupported",
}

func (v *ActionState) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := ActionState(value)
	for _, existing := range AllowedActionStateEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid ActionState", value)
}

// ConfigurationManagerClientEnabledFeatures configuration Manager client enabled features
type ConfigurationManagerClientEnabledFeatures struct {
	// Whether compliance policy is managed by Intune
	CompliancePolicy *bool `json:"compliancePolicy,omitempty"`
	// Whether device configuration is managed by Intune
	DeviceConfiguration *bool `json:"deviceConfiguration,omitempty"`
	// Whether inventory is managed by Intune
	Inventory *bool `json:"inventory,omitempty"`
	// Whether modern application is managed by Intune
	ModernApps *bool `json:"modernApps,omitempty"`
	// Whether resource access is managed by Intune
	ResourceAccess *bool `json:"resourceAccess,omitempty"`
	// Whether Windows Update for Business is managed by Intune
	WindowsUpdateForBusiness *bool  `json:"windowsUpdateForBusiness,omitempty"`
	OdataType                string `json:"@odata.type"`
}

// ComplianceState Compliance state.
type ComplianceState string

// List of microsoft.graph.complianceState
const (
	COMPLIANCESTATE_UNKNOWN         ComplianceState = "unknown"
	COMPLIANCESTATE_COMPLIANT       ComplianceState = "compliant"
	COMPLIANCESTATE_NONCOMPLIANT    ComplianceState = "noncompliant"
	COMPLIANCESTATE_CONFLICT        ComplianceState = "conflict"
	COMPLIANCESTATE_ERROR           ComplianceState = "error"
	COMPLIANCESTATE_IN_GRACE_PERIOD ComplianceState = "inGracePeriod"
	COMPLIANCESTATE_CONFIG_MANAGER  ComplianceState = "configManager"
)

// All allowed values of ComplianceState enum
var AllowedComplianceStateEnumValues = []ComplianceState{
	"unknown",
	"compliant",
	"noncompliant",
	"conflict",
	"error",
	"inGracePeriod",
	"configManager",
}

func (v *ComplianceState) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := ComplianceState(value)
	for _, existing := range AllowedComplianceStateEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid ComplianceState", value)
}

// DetectedAppPlatformType Indicates the operating system / platform of the discovered application.  Some possible values are Windows, iOS, macOS. The default value is unknown (0).
type DetectedAppPlatformType string

// List of microsoft.graph.detectedAppPlatformType
const (
	DETECTEDAPPPLATFORMTYPE_UNKNOWN                             DetectedAppPlatformType = "unknown"
	DETECTEDAPPPLATFORMTYPE_WINDOWS                             DetectedAppPlatformType = "windows"
	DETECTEDAPPPLATFORMTYPE_WINDOWS_MOBILE                      DetectedAppPlatformType = "windowsMobile"
	DETECTEDAPPPLATFORMTYPE_WINDOWS_HOLOGRAPHIC                 DetectedAppPlatformType = "windowsHolographic"
	DETECTEDAPPPLATFORMTYPE_IOS                                 DetectedAppPlatformType = "ios"
	DETECTEDAPPPLATFORMTYPE_MAC_OS                              DetectedAppPlatformType = "macOS"
	DETECTEDAPPPLATFORMTYPE_CHROME_OS                           DetectedAppPlatformType = "chromeOS"
	DETECTEDAPPPLATFORMTYPE_ANDROID_OSP                         DetectedAppPlatformType = "androidOSP"
	DETECTEDAPPPLATFORMTYPE_ANDROID_DEVICE_ADMINISTRATOR        DetectedAppPlatformType = "androidDeviceAdministrator"
	DETECTEDAPPPLATFORMTYPE_ANDROID_WORK_PROFILE                DetectedAppPlatformType = "androidWorkProfile"
	DETECTEDAPPPLATFORMTYPE_ANDROID_DEDICATED_AND_FULLY_MANAGED DetectedAppPlatformType = "androidDedicatedAndFullyManaged"
	DETECTEDAPPPLATFORMTYPE_UNKNOWN_FUTURE_VALUE                DetectedAppPlatformType = "unknownFutureValue"
)

// All allowed values of DetectedAppPlatformType enum
var AllowedDetectedAppPlatformTypeEnumValues = []DetectedAppPlatformType{
	"unknown",
	"windows",
	"windowsMobile",
	"windowsHolographic",
	"ios",
	"macOS",
	"chromeOS",
	"androidOSP",
	"androidDeviceAdministrator",
	"androidWorkProfile",
	"androidDedicatedAndFullyManaged",
	"unknownFutureValue",
}

func (v *DetectedAppPlatformType) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := DetectedAppPlatformType(value)
	for _, existing := range AllowedDetectedAppPlatformTypeEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid DetectedAppPlatformType", value)
}

// OnPremisesConditionalAccessSettings struct for OnPremisesConditionalAccessSettings
type OnPremisesConditionalAccessSettings struct {
	Entity
	// Indicates if on premises conditional access is enabled for this organization
	Enabled *bool `json:"enabled,omitempty"`
	// User groups that will be exempt by on premises conditional access. All users in these groups will be exempt from the conditional access policy.
	ExcludedGroups []string `json:"excludedGroups,omitempty"`
	// User groups that will be targeted by on premises conditional access. All users in these groups will be required to have mobile device managed and compliant for mail access.
	IncludedGroups []string `json:"includedGroups,omitempty"`
	// Override the default access rule when allowing a device to ensure access is granted.
	OverrideDefaultRule *bool  `json:"overrideDefaultRule,omitempty"`
	OdataType           string `json:"@odata.type"`
}

// ComplianceManagementPartner struct for ComplianceManagementPartner
type ComplianceManagementPartner struct {
	Entity
	// User groups which enroll Android devices through partner.
	AndroidEnrollmentAssignments []ComplianceManagementPartnerAssignment `json:"androidEnrollmentAssignments,omitempty"`
	// Partner onboarded for Android devices.
	AndroidOnboarded *bool `json:"androidOnboarded,omitempty"`
	// Partner display name
	DisplayName *string `json:"displayName,omitempty"`
	// User groups which enroll ios devices through partner.
	IosEnrollmentAssignments []ComplianceManagementPartnerAssignment `json:"iosEnrollmentAssignments,omitempty"`
	// Partner onboarded for ios devices.
	IosOnboarded *bool `json:"iosOnboarded,omitempty"`
	// Timestamp of last heartbeat after admin onboarded to the compliance management partner
	LastHeartbeatDateTime *time.Time `json:"lastHeartbeatDateTime,omitempty"`
	// User groups which enroll Mac devices through partner.
	MacOsEnrollmentAssignments []ComplianceManagementPartnerAssignment `json:"macOsEnrollmentAssignments,omitempty"`
	// Partner onboarded for Mac devices.
	MacOsOnboarded *bool                               `json:"macOsOnboarded,omitempty"`
	PartnerState   *DeviceManagementPartnerTenantState `json:"partnerState,omitempty"`
	OdataType      string                              `json:"@odata.type"`
}

// DeviceManagementPartnerTenantState Partner state of this tenant.
type DeviceManagementPartnerTenantState string

// List of microsoft.graph.deviceManagementPartnerTenantState
const (
	DEVICEMANAGEMENTPARTNERTENANTSTATE_UNKNOWN      DeviceManagementPartnerTenantState = "unknown"
	DEVICEMANAGEMENTPARTNERTENANTSTATE_UNAVAILABLE  DeviceManagementPartnerTenantState = "unavailable"
	DEVICEMANAGEMENTPARTNERTENANTSTATE_ENABLED      DeviceManagementPartnerTenantState = "enabled"
	DEVICEMANAGEMENTPARTNERTENANTSTATE_TERMINATED   DeviceManagementPartnerTenantState = "terminated"
	DEVICEMANAGEMENTPARTNERTENANTSTATE_REJECTED     DeviceManagementPartnerTenantState = "rejected"
	DEVICEMANAGEMENTPARTNERTENANTSTATE_UNRESPONSIVE DeviceManagementPartnerTenantState = "unresponsive"
)

// All allowed values of DeviceManagementPartnerTenantState enum
var AllowedDeviceManagementPartnerTenantStateEnumValues = []DeviceManagementPartnerTenantState{
	"unknown",
	"unavailable",
	"enabled",
	"terminated",
	"rejected",
	"unresponsive",
}

func (v *DeviceManagementPartnerTenantState) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := DeviceManagementPartnerTenantState(value)
	for _, existing := range AllowedDeviceManagementPartnerTenantStateEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid DeviceManagementPartnerTenantState", value)
}

// ComplianceManagementPartnerAssignment User group targeting for Compliance Management Partner
type ComplianceManagementPartnerAssignment struct {
	Target    *DeviceAndAppManagementAssignmentTarget `json:"target,omitempty"`
	OdataType string                                  `json:"@odata.type"`
}

// AuditEvent struct for AuditEvent
type AuditEvent struct {
	Entity
	// Friendly name of the activity.
	Activity *string `json:"activity,omitempty"`
	// The date time in UTC when the activity was performed.
	ActivityDateTime *time.Time `json:"activityDateTime,omitempty"`
	// The HTTP operation type of the activity.
	ActivityOperationType *string `json:"activityOperationType,omitempty"`
	// The result of the activity.
	ActivityResult *string `json:"activityResult,omitempty"`
	// The type of activity that was being performed.
	ActivityType *string     `json:"activityType,omitempty"`
	Actor        *AuditActor `json:"actor,omitempty"`
	// Audit category.
	Category *string `json:"category,omitempty"`
	// Component name.
	ComponentName *string `json:"componentName,omitempty"`
	// The client request Id that is used to correlate activity within the system.
	CorrelationId *string `json:"correlationId,omitempty"`
	// Event display name.
	DisplayName *string `json:"displayName,omitempty"`
	// Resources being modified.
	Resources []AuditResource `json:"resources,omitempty"`
	OdataType string          `json:"@odata.type"`
}

// AuditResource A class containing the properties for Audit Resource.
type AuditResource struct {
	// Audit resource's type.
	AuditResourceType *string `json:"auditResourceType,omitempty"`
	// Display name.
	DisplayName *string `json:"displayName,omitempty"`
	// List of modified properties.
	ModifiedProperties []AuditProperty `json:"modifiedProperties,omitempty"`
	// Audit resource's Id.
	ResourceId *string `json:"resourceId,omitempty"`
	OdataType  string  `json:"@odata.type"`
}

// AuditProperty A class containing the properties for Audit Property.
type AuditProperty struct {
	// Display name.
	DisplayName *string `json:"displayName,omitempty"`
	// New value.
	NewValue *string `json:"newValue,omitempty"`
	// Old value.
	OldValue  *string `json:"oldValue,omitempty"`
	OdataType string  `json:"@odata.type"`
}

// AuditActor A class containing the properties for Audit Actor.
type AuditActor struct {
	// Name of the Application.
	ApplicationDisplayName *string `json:"applicationDisplayName,omitempty"`
	// AAD Application Id.
	ApplicationId *string `json:"applicationId,omitempty"`
	// Actor Type.
	AuditActorType *string `json:"auditActorType,omitempty"`
	// IPAddress.
	IpAddress *string `json:"ipAddress,omitempty"`
	// Service Principal Name (SPN).
	ServicePrincipalName *string `json:"servicePrincipalName,omitempty"`
	// User Id.
	UserId *string `json:"userId,omitempty"`
	// List of user permissions when the audit was performed.
	UserPermissions []*string `json:"userPermissions,omitempty"`
	// User Principal Name (UPN).
	UserPrincipalName *string `json:"userPrincipalName,omitempty"`
	OdataType         string  `json:"@odata.type"`
}

// DeviceManagementApplePushNotificationCertificate struct for ApplePushNotificationCertificate
type DeviceManagementApplePushNotificationCertificate struct {
	Entity
	// Apple Id of the account used to create the MDM push certificate.
	AppleIdentifier *string `json:"appleIdentifier,omitempty"`
	// Not yet documented
	Certificate *string `json:"certificate,omitempty"`
	// Certificate serial number. This property is read-only.
	CertificateSerialNumber *string `json:"certificateSerialNumber,omitempty"`
	// The reason the certificate upload failed.
	CertificateUploadFailureReason *string `json:"certificateUploadFailureReason,omitempty"`
	// The certificate upload status.
	CertificateUploadStatus *string `json:"certificateUploadStatus,omitempty"`
	// The expiration date and time for Apple push notification certificate.
	ExpirationDateTime *time.Time `json:"expirationDateTime,omitempty"`
	// Last modified date and time for Apple push notification certificate.
	LastModifiedDateTime *time.Time `json:"lastModifiedDateTime,omitempty"`
	// Topic Id.
	TopicIdentifier *string `json:"topicIdentifier,omitempty"`
	OdataType       string  `json:"@odata.type"`
}

// DeviceManagementSubscriptionState Tenant mobile device management subscription state.
type DeviceManagementSubscriptionState string

// List of microsoft.graph.deviceManagementSubscriptionState
const (
	DEVICEMANAGEMENTSUBSCRIPTIONSTATE_PENDING    DeviceManagementSubscriptionState = "pending"
	DEVICEMANAGEMENTSUBSCRIPTIONSTATE_ACTIVE     DeviceManagementSubscriptionState = "active"
	DEVICEMANAGEMENTSUBSCRIPTIONSTATE_WARNING    DeviceManagementSubscriptionState = "warning"
	DEVICEMANAGEMENTSUBSCRIPTIONSTATE_DISABLED   DeviceManagementSubscriptionState = "disabled"
	DEVICEMANAGEMENTSUBSCRIPTIONSTATE_DELETED    DeviceManagementSubscriptionState = "deleted"
	DEVICEMANAGEMENTSUBSCRIPTIONSTATE_BLOCKED    DeviceManagementSubscriptionState = "blocked"
	DEVICEMANAGEMENTSUBSCRIPTIONSTATE_LOCKED_OUT DeviceManagementSubscriptionState = "lockedOut"
)

// All allowed values of DeviceManagementSubscriptionState enum
var AllowedDeviceManagementSubscriptionStateEnumValues = []DeviceManagementSubscriptionState{
	"pending",
	"active",
	"warning",
	"disabled",
	"deleted",
	"blocked",
	"lockedOut",
}

func (v *DeviceManagementSubscriptionState) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := DeviceManagementSubscriptionState(value)
	for _, existing := range AllowedDeviceManagementSubscriptionStateEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid DeviceManagementSubscriptionState", value)
}

// UserExperienceAnalyticsSettings The user experience analytics insight is the recomendation to improve the user experience analytics score.
type DeviceManagementUserExperienceAnalyticsSettings struct {
	// When TRUE, indicates Tenant attach is configured properly and System Center Configuration Manager (SCCM) tenant attached devices will show up in endpoint analytics reporting. When FALSE, indicates Tenant attach is not configured. FALSE by default.
	ConfigurationManagerDataConnectorConfigured *bool  `json:"configurationManagerDataConnectorConfigured,omitempty"`
	OdataType                                   string `json:"@odata.type"`
}

// DeviceManagementWindowsMalwareOverview Windows device malware overview.
type DeviceManagementWindowsMalwareOverview struct {
	// List of device counts per malware category
	MalwareCategorySummary []WindowsMalwareCategoryCount `json:"malwareCategorySummary,omitempty"`
	// Count of devices with malware detected in the last 30 days
	MalwareDetectedDeviceCount *int32 `json:"malwareDetectedDeviceCount,omitempty"`
	// List of device counts per malware execution state
	MalwareExecutionStateSummary []WindowsMalwareExecutionStateCount `json:"malwareExecutionStateSummary,omitempty"`
	// List of device counts per malware
	MalwareNameSummary []WindowsMalwareNameCount `json:"malwareNameSummary,omitempty"`
	// List of active malware counts per malware severity
	MalwareSeveritySummary []WindowsMalwareSeverityCount `json:"malwareSeveritySummary,omitempty"`
	// List of device counts per malware state
	MalwareStateSummary []WindowsMalwareStateCount `json:"malwareStateSummary,omitempty"`
	// List of device counts with malware per windows OS version
	OsVersionsSummary []OsVersionCount `json:"osVersionsSummary,omitempty"`
	// Count of all distinct malwares detected across all devices. Valid values -2147483648 to 2147483647
	TotalDistinctMalwareCount *int32 `json:"totalDistinctMalwareCount,omitempty"`
	// Count of all malware detections across all devices. Valid values -2147483648 to 2147483647
	TotalMalwareCount *int32 `json:"totalMalwareCount,omitempty"`
	OdataType         string `json:"@odata.type"`
}

// OsVersionCount Count of devices with malware for each OS version
type OsVersionCount struct {
	// Count of devices with malware for the OS version
	DeviceCount *int32 `json:"deviceCount,omitempty"`
	// The Timestamp of the last update for the device count in UTC
	LastUpdateDateTime *time.Time `json:"lastUpdateDateTime,omitempty"`
	// OS version
	OsVersion *string `json:"osVersion,omitempty"`
	OdataType string  `json:"@odata.type"`
}

// WindowsMalwareStateCount Windows Malware State Summary.
type WindowsMalwareStateCount struct {
	// Count of devices with malware detections for this malware State
	DeviceCount *int32 `json:"deviceCount,omitempty"`
	// Count of distinct malwares for this malware State. Valid values -2147483648 to 2147483647
	DistinctMalwareCount *int32 `json:"distinctMalwareCount,omitempty"`
	// The Timestamp of the last update for the device count in UTC
	LastUpdateDateTime *time.Time `json:"lastUpdateDateTime,omitempty"`
	// Count of total malware detections for this malware State. Valid values -2147483648 to 2147483647
	MalwareDetectionCount *int32                     `json:"malwareDetectionCount,omitempty"`
	State                 *WindowsMalwareThreatState `json:"state,omitempty"`
	OdataType             string                     `json:"@odata.type"`
}

// WindowsMalwareThreatState Malware threat status
type WindowsMalwareThreatState string

// List of microsoft.graph.windowsMalwareThreatState
const (
	WINDOWSMALWARETHREATSTATE_ACTIVE                                WindowsMalwareThreatState = "active"
	WINDOWSMALWARETHREATSTATE_ACTION_FAILED                         WindowsMalwareThreatState = "actionFailed"
	WINDOWSMALWARETHREATSTATE_MANUAL_STEPS_REQUIRED                 WindowsMalwareThreatState = "manualStepsRequired"
	WINDOWSMALWARETHREATSTATE_FULL_SCAN_REQUIRED                    WindowsMalwareThreatState = "fullScanRequired"
	WINDOWSMALWARETHREATSTATE_REBOOT_REQUIRED                       WindowsMalwareThreatState = "rebootRequired"
	WINDOWSMALWARETHREATSTATE_REMEDIATED_WITH_NON_CRITICAL_FAILURES WindowsMalwareThreatState = "remediatedWithNonCriticalFailures"
	WINDOWSMALWARETHREATSTATE_QUARANTINED                           WindowsMalwareThreatState = "quarantined"
	WINDOWSMALWARETHREATSTATE_REMOVED                               WindowsMalwareThreatState = "removed"
	WINDOWSMALWARETHREATSTATE_CLEANED                               WindowsMalwareThreatState = "cleaned"
	WINDOWSMALWARETHREATSTATE_ALLOWED                               WindowsMalwareThreatState = "allowed"
	WINDOWSMALWARETHREATSTATE_NO_STATUS_CLEARED                     WindowsMalwareThreatState = "noStatusCleared"
)

// All allowed values of WindowsMalwareThreatState enum
var AllowedWindowsMalwareThreatStateEnumValues = []WindowsMalwareThreatState{
	"active",
	"actionFailed",
	"manualStepsRequired",
	"fullScanRequired",
	"rebootRequired",
	"remediatedWithNonCriticalFailures",
	"quarantined",
	"removed",
	"cleaned",
	"allowed",
	"noStatusCleared",
}

func (v *WindowsMalwareThreatState) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := WindowsMalwareThreatState(value)
	for _, existing := range AllowedWindowsMalwareThreatStateEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid WindowsMalwareThreatState", value)
}

// WindowsMalwareSeverityCount Windows Malware Severity Count Summary
type WindowsMalwareSeverityCount struct {
	// Count of distinct malwares for this malware State. Valid values -2147483648 to 2147483647
	DistinctMalwareCount *int32 `json:"distinctMalwareCount,omitempty"`
	// The Timestamp of the last update for the WindowsMalwareSeverityCount in UTC
	LastUpdateDateTime *time.Time `json:"lastUpdateDateTime,omitempty"`
	// Count of threats detections for this malware severity. Valid values -2147483648 to 2147483647
	MalwareDetectionCount *int32                  `json:"malwareDetectionCount,omitempty"`
	Severity              *WindowsMalwareSeverity `json:"severity,omitempty"`
	OdataType             string                  `json:"@odata.type"`
}

// WindowsMalwareSeverity Malware severity
type WindowsMalwareSeverity string

// List of microsoft.graph.windowsMalwareSeverity
const (
	WINDOWSMALWARESEVERITY_UNKNOWN  WindowsMalwareSeverity = "unknown"
	WINDOWSMALWARESEVERITY_LOW      WindowsMalwareSeverity = "low"
	WINDOWSMALWARESEVERITY_MODERATE WindowsMalwareSeverity = "moderate"
	WINDOWSMALWARESEVERITY_HIGH     WindowsMalwareSeverity = "high"
	WINDOWSMALWARESEVERITY_SEVERE   WindowsMalwareSeverity = "severe"
)

// All allowed values of WindowsMalwareSeverity enum
var AllowedWindowsMalwareSeverityEnumValues = []WindowsMalwareSeverity{
	"unknown",
	"low",
	"moderate",
	"high",
	"severe",
}

// WindowsMalwareNameCount Malware name device count
type WindowsMalwareNameCount struct {
	// Count of devices with malware dectected for this malware
	DeviceCount *int32 `json:"deviceCount,omitempty"`
	// The Timestamp of the last update for the device count in UTC
	LastUpdateDateTime *time.Time `json:"lastUpdateDateTime,omitempty"`
	// The unique identifier. This is malware identifier
	MalwareIdentifier *string `json:"malwareIdentifier,omitempty"`
	// Malware name
	Name      *string `json:"name,omitempty"`
	OdataType string  `json:"@odata.type"`
}

// WindowsMalwareCategoryCount Malware category device count
type WindowsMalwareCategoryCount struct {
	// Count of active malware detections for this malware category. Valid values -2147483648 to 2147483647
	ActiveMalwareDetectionCount *int32                  `json:"activeMalwareDetectionCount,omitempty"`
	Category                    *WindowsMalwareCategory `json:"category,omitempty"`
	// Count of devices with malware detections for this malware category
	DeviceCount *int32 `json:"deviceCount,omitempty"`
	// Count of distinct active malwares for this malware category. Valid values -2147483648 to 2147483647
	DistinctActiveMalwareCount *int32 `json:"distinctActiveMalwareCount,omitempty"`
	// The Timestamp of the last update for the device count in UTC
	LastUpdateDateTime *time.Time `json:"lastUpdateDateTime,omitempty"`
	OdataType          string     `json:"@odata.type"`
}

// WindowsMalwareCategory Malware category id
type WindowsMalwareCategory string

// List of microsoft.graph.windowsMalwareCategory
const (
	WINDOWSMALWARECATEGORY_INVALID                      WindowsMalwareCategory = "invalid"
	WINDOWSMALWARECATEGORY_ADWARE                       WindowsMalwareCategory = "adware"
	WINDOWSMALWARECATEGORY_SPYWARE                      WindowsMalwareCategory = "spyware"
	WINDOWSMALWARECATEGORY_PASSWORD_STEALER             WindowsMalwareCategory = "passwordStealer"
	WINDOWSMALWARECATEGORY_TROJAN_DOWNLOADER            WindowsMalwareCategory = "trojanDownloader"
	WINDOWSMALWARECATEGORY_WORM                         WindowsMalwareCategory = "worm"
	WINDOWSMALWARECATEGORY_BACKDOOR                     WindowsMalwareCategory = "backdoor"
	WINDOWSMALWARECATEGORY_REMOTE_ACCESS_TROJAN         WindowsMalwareCategory = "remoteAccessTrojan"
	WINDOWSMALWARECATEGORY_TROJAN                       WindowsMalwareCategory = "trojan"
	WINDOWSMALWARECATEGORY_EMAIL_FLOODER                WindowsMalwareCategory = "emailFlooder"
	WINDOWSMALWARECATEGORY_KEYLOGGER                    WindowsMalwareCategory = "keylogger"
	WINDOWSMALWARECATEGORY_DIALER                       WindowsMalwareCategory = "dialer"
	WINDOWSMALWARECATEGORY_MONITORING_SOFTWARE          WindowsMalwareCategory = "monitoringSoftware"
	WINDOWSMALWARECATEGORY_BROWSER_MODIFIER             WindowsMalwareCategory = "browserModifier"
	WINDOWSMALWARECATEGORY_COOKIE                       WindowsMalwareCategory = "cookie"
	WINDOWSMALWARECATEGORY_BROWSER_PLUGIN               WindowsMalwareCategory = "browserPlugin"
	WINDOWSMALWARECATEGORY_AOL_EXPLOIT                  WindowsMalwareCategory = "aolExploit"
	WINDOWSMALWARECATEGORY_NUKER                        WindowsMalwareCategory = "nuker"
	WINDOWSMALWARECATEGORY_SECURITY_DISABLER            WindowsMalwareCategory = "securityDisabler"
	WINDOWSMALWARECATEGORY_JOKE_PROGRAM                 WindowsMalwareCategory = "jokeProgram"
	WINDOWSMALWARECATEGORY_HOSTILE_ACTIVE_X_CONTROL     WindowsMalwareCategory = "hostileActiveXControl"
	WINDOWSMALWARECATEGORY_SOFTWARE_BUNDLER             WindowsMalwareCategory = "softwareBundler"
	WINDOWSMALWARECATEGORY_STEALTH_NOTIFIER             WindowsMalwareCategory = "stealthNotifier"
	WINDOWSMALWARECATEGORY_SETTINGS_MODIFIER            WindowsMalwareCategory = "settingsModifier"
	WINDOWSMALWARECATEGORY_TOOL_BAR                     WindowsMalwareCategory = "toolBar"
	WINDOWSMALWARECATEGORY_REMOTE_CONTROL_SOFTWARE      WindowsMalwareCategory = "remoteControlSoftware"
	WINDOWSMALWARECATEGORY_TROJAN_FTP                   WindowsMalwareCategory = "trojanFtp"
	WINDOWSMALWARECATEGORY_POTENTIAL_UNWANTED_SOFTWARE  WindowsMalwareCategory = "potentialUnwantedSoftware"
	WINDOWSMALWARECATEGORY_ICQ_EXPLOIT                  WindowsMalwareCategory = "icqExploit"
	WINDOWSMALWARECATEGORY_TROJAN_TELNET                WindowsMalwareCategory = "trojanTelnet"
	WINDOWSMALWARECATEGORY_EXPLOIT                      WindowsMalwareCategory = "exploit"
	WINDOWSMALWARECATEGORY_FILESHARING_PROGRAM          WindowsMalwareCategory = "filesharingProgram"
	WINDOWSMALWARECATEGORY_MALWARE_CREATION_TOOL        WindowsMalwareCategory = "malwareCreationTool"
	WINDOWSMALWARECATEGORY_REMOTE_CONTROL_SOFTWARE2     WindowsMalwareCategory = "remote_Control_Software"
	WINDOWSMALWARECATEGORY_TOOL                         WindowsMalwareCategory = "tool"
	WINDOWSMALWARECATEGORY_TROJAN_DENIAL_OF_SERVICE     WindowsMalwareCategory = "trojanDenialOfService"
	WINDOWSMALWARECATEGORY_TROJAN_DROPPER               WindowsMalwareCategory = "trojanDropper"
	WINDOWSMALWARECATEGORY_TROJAN_MASS_MAILER           WindowsMalwareCategory = "trojanMassMailer"
	WINDOWSMALWARECATEGORY_TROJAN_MONITORING_SOFTWARE   WindowsMalwareCategory = "trojanMonitoringSoftware"
	WINDOWSMALWARECATEGORY_TROJAN_PROXY_SERVER          WindowsMalwareCategory = "trojanProxyServer"
	WINDOWSMALWARECATEGORY_VIRUS                        WindowsMalwareCategory = "virus"
	WINDOWSMALWARECATEGORY_KNOWN                        WindowsMalwareCategory = "known"
	WINDOWSMALWARECATEGORY_UNKNOWN                      WindowsMalwareCategory = "unknown"
	WINDOWSMALWARECATEGORY_SPP                          WindowsMalwareCategory = "spp"
	WINDOWSMALWARECATEGORY_BEHAVIOR                     WindowsMalwareCategory = "behavior"
	WINDOWSMALWARECATEGORY_VULNERABILITY                WindowsMalwareCategory = "vulnerability"
	WINDOWSMALWARECATEGORY_POLICY                       WindowsMalwareCategory = "policy"
	WINDOWSMALWARECATEGORY_ENTERPRISE_UNWANTED_SOFTWARE WindowsMalwareCategory = "enterpriseUnwantedSoftware"
	WINDOWSMALWARECATEGORY_RANSOM                       WindowsMalwareCategory = "ransom"
	WINDOWSMALWARECATEGORY_HIPS_RULE                    WindowsMalwareCategory = "hipsRule"
)

// All allowed values of WindowsMalwareCategory enum
var AllowedWindowsMalwareCategoryEnumValues = []WindowsMalwareCategory{
	"invalid",
	"adware",
	"spyware",
	"passwordStealer",
	"trojanDownloader",
	"worm",
	"backdoor",
	"remoteAccessTrojan",
	"trojan",
	"emailFlooder",
	"keylogger",
	"dialer",
	"monitoringSoftware",
	"browserModifier",
	"cookie",
	"browserPlugin",
	"aolExploit",
	"nuker",
	"securityDisabler",
	"jokeProgram",
	"hostileActiveXControl",
	"softwareBundler",
	"stealthNotifier",
	"settingsModifier",
	"toolBar",
	"remoteControlSoftware",
	"trojanFtp",
	"potentialUnwantedSoftware",
	"icqExploit",
	"trojanTelnet",
	"exploit",
	"filesharingProgram",
	"malwareCreationTool",
	"remote_Control_Software",
	"tool",
	"trojanDenialOfService",
	"trojanDropper",
	"trojanMassMailer",
	"trojanMonitoringSoftware",
	"trojanProxyServer",
	"virus",
	"known",
	"unknown",
	"spp",
	"behavior",
	"vulnerability",
	"policy",
	"enterpriseUnwantedSoftware",
	"ransom",
	"hipsRule",
}

func (v *WindowsMalwareCategory) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := WindowsMalwareCategory(value)
	for _, existing := range AllowedWindowsMalwareCategoryEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid WindowsMalwareCategory", value)
}

// WindowsMalwareExecutionStateCount Windows malware execution state summary.
type WindowsMalwareExecutionStateCount struct {
	// Count of devices with malware detections for this malware execution state
	DeviceCount    *int32                        `json:"deviceCount,omitempty"`
	ExecutionState *WindowsMalwareExecutionState `json:"executionState,omitempty"`
	// The Timestamp of the last update for the device count in UTC
	LastUpdateDateTime *time.Time `json:"lastUpdateDateTime,omitempty"`
	OdataType          string     `json:"@odata.type"`
}

// WindowsMalwareExecutionState Malware execution status
type WindowsMalwareExecutionState string

// List of microsoft.graph.windowsMalwareExecutionState
const (
	WINDOWSMALWAREEXECUTIONSTATE_UNKNOWN     WindowsMalwareExecutionState = "unknown"
	WINDOWSMALWAREEXECUTIONSTATE_BLOCKED     WindowsMalwareExecutionState = "blocked"
	WINDOWSMALWAREEXECUTIONSTATE_ALLOWED     WindowsMalwareExecutionState = "allowed"
	WINDOWSMALWAREEXECUTIONSTATE_RUNNING     WindowsMalwareExecutionState = "running"
	WINDOWSMALWAREEXECUTIONSTATE_NOT_RUNNING WindowsMalwareExecutionState = "notRunning"
)

// All allowed values of WindowsMalwareExecutionState enum
var AllowedWindowsMalwareExecutionStateEnumValues = []WindowsMalwareExecutionState{
	"unknown",
	"blocked",
	"allowed",
	"running",
	"notRunning",
}

func (v *WindowsMalwareExecutionState) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := WindowsMalwareExecutionState(value)
	for _, existing := range AllowedWindowsMalwareExecutionStateEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid WindowsMalwareExecutionState", value)
}
