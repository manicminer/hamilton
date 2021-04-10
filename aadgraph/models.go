package aadgraph

import (
	"encoding/json"

	"github.com/manicminer/hamilton/environments"
	"github.com/manicminer/hamilton/msgraph"
)

type ApplicationRef struct {
	AppId                   *environments.ApiAppId            `json:"appId,omitempty"`
	AppCategory             *json.RawMessage                  `json:"appCategory"`
	AppContextId            *string                           `json:"appContextId"`
	AppData                 *json.RawMessage                  `json:"appData"`
	AppRoles                *[]msgraph.AppRole                `json:"appRoles,omitempty"`
	AvailableToOtherTenants *bool                             `json:"availableToOtherTenants"`
	DisplayName             *string                           `json:"displayName,omitempty"`
	ErrorUrl                *string                           `json:"errorUrl"`
	Homepage                *string                           `json:"homepage"`
	IdentifierUris          *[]string                         `json:"identifierUris,omitempty"`
	KnownClientApplications *[]string                         `json:"knownClientApplications"`
	LogoutUrl               *string                           `json:"logoutUrl,omitempty"`
	LogoUrl                 *string                           `json:"logoUrl,omitempty"`
	OAuth2Permissions       *[]msgraph.PermissionScope        `json:"oauth2Permissions,omitempty"`
	PublisherDomain         *string                           `json:"publisherDomain,omitempty"`
	PublisherName           *string                           `json:"publisherName,omitempty"`
	PublicClient            *bool                             `json:"publicClient"`
	ReplyUrls               *[]string                         `json:"replyUrls,omitempty"`
	RequiredResourceAccess  *[]msgraph.RequiredResourceAccess `json:"requiredResourceAccess,omitempty"`
	SamlMetadataUrl         *string                           `json:"samlMetadataUrl"`
	SupportsConvergence     *bool                             `json:"supportsConvergence"`
	VerifiedPublisher       *msgraph.VerifiedPublisher        `json:"verifiedPublisher"`
}
