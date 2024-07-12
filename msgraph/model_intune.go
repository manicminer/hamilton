package msgraph

// DeviceManagementIntuneBrand intuneBrand contains data which is used in customizing the appearance of the Company Portal applications as well as the end user web portal.
type DeviceManagementIntuneBrand struct {
	// Email address of the person/organization responsible for IT support.
	ContactITEmailAddress *string `json:"contactITEmailAddress,omitempty"`
	// Name of the person/organization responsible for IT support.
	ContactITName *string `json:"contactITName,omitempty"`
	// Text comments regarding the person/organization responsible for IT support.
	ContactITNotes *string `json:"contactITNotes,omitempty"`
	// Phone number of the person/organization responsible for IT support.
	ContactITPhoneNumber *string      `json:"contactITPhoneNumber,omitempty"`
	DarkBackgroundLogo   *MimeContent `json:"darkBackgroundLogo,omitempty"`
	// Company/organization name that is displayed to end users.
	DisplayName         *string      `json:"displayName,omitempty"`
	LightBackgroundLogo *MimeContent `json:"lightBackgroundLogo,omitempty"`
	// Display name of the company/organization’s IT helpdesk site.
	OnlineSupportSiteName *string `json:"onlineSupportSiteName,omitempty"`
	// URL to the company/organization’s IT helpdesk site.
	OnlineSupportSiteUrl *string `json:"onlineSupportSiteUrl,omitempty"`
	// URL to the company/organization’s privacy policy.
	PrivacyUrl *string `json:"privacyUrl,omitempty"`
	// Boolean that represents whether the administrator-supplied display name will be shown next to the logo image.
	ShowDisplayNameNextToLogo *bool `json:"showDisplayNameNextToLogo,omitempty"`
	// Boolean that represents whether the administrator-supplied logo images are shown or not shown.
	ShowLogo *bool `json:"showLogo,omitempty"`
	// Boolean that represents whether the administrator-supplied display name will be shown next to the logo image.
	ShowNameNextToLogo *bool     `json:"showNameNextToLogo,omitempty"`
	ThemeColor         *RgbColor `json:"themeColor,omitempty"`
	OdataType          string    `json:"@odata.type"`
}

// MimeContent Contains properties for a generic mime content.
type MimeContent struct {
	// Indicates the content mime type.
	Type *string `json:"type,omitempty"`
	// The byte array that contains the actual content.
	Value     *string `json:"value,omitempty"`
	OdataType string  `json:"@odata.type"`
}

// RgbColor Color in RGB.
type RgbColor struct {
	// Blue value
	B *int32 `json:"b,omitempty"`
	// Green value
	G *int32 `json:"g,omitempty"`
	// Red value
	R         *int32 `json:"r,omitempty"`
	OdataType string `json:"@odata.type"`
}
