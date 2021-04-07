package models

type NamedLocation struct {
	ODataType   *string `json:"@odata.type,omitempty"`
	ID          *string `json:"id,omitempty"`
	DisplayName *string `json:"displayName,omitempty"`
	// CreatedDateTime  *time.Time `json:"createdDateTime,omitempty"`
	// ModifiedDateTime *time.Time `json:"modifiedDateTime,omitempty"`
}

// CountryNamedLocation describes an Country Named Location object.
type CountryNamedLocation struct {
	*NamedLocation
	CountriesAndRegions               *[]string `json:"countriesAndRegions,omitempty"`
	IncludeUnknownCountriesAndRegions *bool     `json:"includeUnknownCountriesAndRegions,omitempty"`
}

// IPNamedLocation describes an IP Named Location object.
type IPNamedLocation struct {
	*NamedLocation
	IPRanges  *[]IPNamedLocationIPRange `json:"ipRanges,omitempty"`
	IsTrusted *bool                     `json:"isTrusted,omitempty"`
}

type IPNamedLocationIPRange struct {
	CIDRAddress *string `json:"cidrAddress,omitempty"`
}
