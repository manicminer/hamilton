package models

// CountryNamedLocation describes an Country Named Location object.
type CountryNamedLocation struct {
	ODataType                         *string   `json:"@odata.type,omitempty"`
	ID                                *string   `json:"id,omitempty"`
	CountriesAndRegions               *[]string `json:"countriesAndRegions,omitempty"`
	DisplayName                       *string   `json:"displayName,omitempty"`
	IncludeUnknownCountriesAndRegions *bool     `json:"includeUnknownCountriesAndRegions,omitempty"`
	// CreatedDateTime                   *time.Time `json:"createdDateTime,omitempty"`
	// ModifiedDateTime                  *time.Time `json:"modifiedDateTime,omitempty"`
}
