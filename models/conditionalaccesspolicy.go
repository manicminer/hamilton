package models

import (
	"time"
)

// ConditionalAccessPolicy describes an Conditional Access Policy object.
type ConditionalAccessPolicy struct {
	Conditions       *ConditionalAccessConditionSet    `json:"conditions,omitempty"`
	CreatedDateTime  *time.Time                        `json:"createdDateTime,omitempty"`
	DisplayName      *string                           `json:"displayName,omitempty"`
	GrantControls    *ConditionalAccessGrantControls   `json:"grantControls,omitempty"`
	ID               *string                           `json:"id,omitempty"`
	ModifiedDateTime *time.Time                        `json:"modifiedDateTime,omitempty"`
	SessionControls  *ConditionalAccessSessionControls `json:"sessionControls,omitempty"`
	State            *string                           `json:"state,omitempty"`
}

type ConditionalAccessConditionSet struct {
	Applications     *[]ConditionalAccessApplications `json:"applications,omitempty"`
	Users            *[]ConditionalAccessUsers        `json:"users,omitempty"`
	ClientAppTypes   *[]string                        `json:"clientAppTypes,omitempty"`
	Locations        *[]ConditionalAccessLocations    `json:"locations,omitempty"`
	Platforms        *[]ConditionalAccessPlatforms    `json:"platforms,omitempty"`
	SignInRiskLevels *[]string                        `json:"signInRiskLevels,omitempty"`
	UserRiskLevels   *[]string                        `json:"userRiskLevels,omitempty"`
}

type ConditionalAccessApplications struct {
	IncludeApplications *[]string `json:"includeApplications,omitempty"`
	ExcludeApplications *[]string `json:"excludeApplications,omitempty"`
	IncludeUserActions  *[]string `json:"includeUserActions,omitempty"`
}

type ConditionalAccessUsers struct {
	IncludeUsers  *[]string `json:"includeUsers,omitempty"`
	ExcludeUsers  *[]string `json:"excludeUsers,omitempty"`
	IncludeGroups *[]string `json:"includeGroups,omitempty"`
	ExcludeGroups *[]string `json:"excludeGroups,omitempty"`
	IncludeRoles  *[]string `json:"includeRoles,omitempty"`
	ExcludeRoles  *[]string `json:"excludeRoles,omitempty"`
}

type ConditionalAccessLocations struct {
	IncludeLocations *[]string `json:"includeLocations,omitempty"`
	ExcludeLocations *[]string `json:"excludeLocations,omitempty"`
}

type ConditionalAccessPlatforms struct {
	IncludePlatforms *[]string `json:"includePlatforms,omitempty"`
	ExcludePlatforms *[]string `json:"excludePlatforms,omitempty"`
}

type ConditionalAccessGrantControls struct {
	Operator                    *string   `json:"operator,omitempty"`
	BuiltInControls             *[]string `json:"builtInControls,omitempty"`
	CustomAuthenticationFactors *[]string `json:"customAuthenticationFactors,omitempty"`
	TermsOfUse                  *[]string `json:"termsOfUse,omitempty"`
}

type ConditionalAccessSessionControls struct {
	ApplicationEnforcedRestrictions *ApplicationEnforcedRestrictionsSessionControl `json:"applicationEnforcedRestrictions,omitempty"`
	CloudAppSecurity                *CloudAppSecurityControl                       `json:"cloudAppSecurity,omitempty"`
	PersistentBrowser               *PersistentBrowserSessionControl               `json:"persistentBrowser,omitempty"`
	SignInFrequency                 *SignInFrequencySessionControl                 `json:"signInFrequency,omitempty"`
}

type ApplicationEnforcedRestrictionsSessionControl struct {
	IsEnabled *bool `json:"isEnabled,omitempty"`
}

type CloudAppSecurityControl struct {
	IsEnabled            *bool   `json:"isEnabled,omitempty"`
	CloudAppSecurityType *string `json:"cloudAppSecurityType,omitempty"`
}

type PersistentBrowserSessionControl struct {
	IsEnabled *bool   `json:"isEnabled,omitempty"`
	Mode      *string `json:"mode,omitempty"`
}

type SignInFrequencySessionControl struct {
	IsEnabled *bool   `json:"isEnabled,omitempty"`
	Type      *string `json:"type,omitempty"`
	Value     *int32  `json:"value,omitempty"`
}
