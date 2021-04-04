package models

// IPNamedLocation describes an IP Named Location object.
type IPNamedLocation struct {
	ODataType   *string                   `json:"@odata.type,omitempty"`
	ID          *string                   `json:"id,omitempty"`
	DisplayName *string                   `json:"displayName,omitempty"`
	IPRanges    *[]IPNamedLocationIPRange `json:"ipRanges,omitempty"`
	IsTrusted   *bool                     `json:"isTrusted,omitempty"`
	// CreatedDateTime  *time.Time                `json:"createdDateTime,omitempty"`
	// ModifiedDateTime *time.Time                `json:"modifiedDateTime,omitempty"`
}

type IPNamedLocationIPRange struct {
	CIDRAddress *string `json:"cidrAddress,omitempty"`
}
