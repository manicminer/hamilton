package models

import (
	"fmt"
	"time"
)

// Group describes a Group object.
type Group struct {
	ID                            *string                             `json:"id,omitempty,readonly"`
	AllowExternalSenders          *string                             `json:"allowExternalSenders,omitempty"`
	AssignedLabels                *[]GroupAssignedLabel               `json:"assignedLabels,omitempty"`
	AssignedLicenses              *[]GroupAssignedLicense             `json:"assignLicenses,omitempty"`
	AutoSubscribeNewMembers       *bool                               `json:"autoSubscribeNewMembers,omitempty"`
	Classification                *string                             `json:"classification,omitempty"`
	CreatedDateTime               *time.Time                          `json:"createdDateTime,omitempty,readonly"`
	DeletedDateTime               *time.Time                          `json:"deletedDateTime,omitempty,readonly"`
	Description                   *string                             `json:"description,omitempty"`
	DisplayName                   *string                             `json:"displayName,omitempty"`
	ExpirationDateTime            *time.Time                          `json:"expirationDateTime,omitempty,readonly"`
	GroupTypes                    *[]string                           `json:"groupTypes,omitempty"`
	HasMembersWithLicenseErrors   *bool                               `json:"hasMembersWithLicenseErrors,omitempty"`
	HideFromAddressLists          *bool                               `json:"hideFromAddressLists,omitempty"`
	HideFromOutlookClients        *bool                               `json:"hideFromOutlookClients,omitempty"`
	IsSubscribedByMail            *bool                               `json:"isSubscribedByMail,omitempty"`
	LicenseProcessingState        *string                             `json:"licenseProcessingState,omitempty,readonly"`
	Mail                          *string                             `json:"mail,omitempty"`
	MailEnabled                   *bool                               `json:"mailEnabled,omitempty"`
	MailNickname                  *string                             `json:"mailNickname,omitempty"`
	MembershipRule                *string                             `json:"membershipRule,omitempty"`
	MembershipRuleProcessingState *string                             `json:"membershipRuleProcessingState,omitempty"`
	OnPremisesDomainName          *string                             `json:"onPremisesDomainName,omitempty,readonly"`
	OnPremisesLastSyncDateTime    *time.Time                          `json:"onPremisesLastSyncDateTime,omitempty,readonly"`
	OnPremisesNetBiosName         *string                             `json:"onPremisesNetBiosName,omitempty,readonly"`
	OnPremisesProvisioningErrors  *[]GroupOnPremisesProvisioningError `json:"onPremisesProvisioningErrors,omitempty"`
	OnPremisesSamAccountName      *string                             `json:"onPremisesSamAccountName,omitempty,readonly"`
	OnPremisesSecurityIdentifier  *string                             `json:"onPremisesSecurityIdentifier,omitempty,readonly"`
	OnPremisesSyncEnabled         *bool                               `json:"onPremisesSyncEnabled,omitempty,readonly"`
	PreferredDataLocation         *string                             `json:"preferredDataLocation,omitempty"`
	PreferredLanguage             *string                             `json:"preferredLanguage,omitempty"`
	ProxyAddresses                *[]string                           `json:"proxyAddresses,omitempty,readonly"`
	RenewedDateTime               *time.Time                          `json:"renewedDateTime,omitempty,readonly"`
	SecurityEnabled               *bool                               `json:"securityEnabled,omitempty"`
	SecurityIdentifier            *string                             `json:"securityIdentifier,omitempty"`
	Theme                         *string                             `json:"theme,omitempty"`
	UnseenCount                   *int                                `json:"unseenCount,omitempty,readonly"`
	Visibility                    *string                             `json:"visibility,omitempty"`

	Members *[]string `json:"members@odata.bind,omitempty"`
	Owners  *[]string `json:"owners@odata.bind,omitempty"`
}

// AppendMember appends a new member object URI to the Members slice.
func (g *Group) AppendMember(endpoint string, apiVersion string, id string) {
	val := fmt.Sprintf("%s/%s/directoryObjects/%s", endpoint, apiVersion, id)
	var members []string
	if g.Members != nil {
		members = *g.Members
	}
	members = append(members, val)
	g.Members = &members
}

// AppendOwner appends a new owner object URI to the Owners slice.
func (g *Group) AppendOwner(endpoint string, apiVersion string, id string) {
	val := fmt.Sprintf("%s/%s/directoryObjects/%s", endpoint, apiVersion, id)
	var owners []string
	if g.Owners != nil {
		owners = *g.Owners
	}
	owners = append(owners, val)
	g.Owners = &owners
}

type GroupAssignedLabel struct {
	LabelId     *string `json:"labelId, omitempty"`
	DisplayName *string `json:"displayNanme,omitempty"`
}

type GroupAssignedLicense struct {
	DisabledPlans *[]string `json:"disabledPlans, omitempty"`
	SkuId         *string   `json:"skuId,omitempty"`
}

type GroupOnPremisesProvisioningError struct {
	Category             *string   `json:"category,omitempty"`
	OccurredDateTime     time.Time `json:"occurredDateTime,omitempty"`
	PropertyCausingError *string   `json:"propertyCausingError,omitempty"`
	Value                *string   `json:"value,omitempty"`
}
