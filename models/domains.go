package models

import (
	"time"
)

type Domain struct {
	ID                               *string   `json:"id,omitempty,readonly"`
	AuthenticationType               *string   `json:"authenticationType,omitempty,readonly"`
	IsAdminManaged                   *bool     `json:"isAdminManaged,omitempty,readonly"`
	IsDefault                        *bool     `json:"isDefault,omitempty,readonly"`
	IsInitial                        *bool     `json:"isInitial,omitempty,readonly"`
	IsRoot                           *bool     `json:"isRoot,omitempty,readonly"`
	IsVerified                       *bool     `json:"isVerified,omitempty,readonly"`
	PasswordNotificationWindowInDays *int      `json:"passwordNotificationWindowInDays,omitempty"`
	PasswordValidityPeriodInDays     *int      `json:"passwordValidityPeriodInDays,omitempty"`
	SupportedServices                *[]string `json:"supportedServices,omitempty,readonly"`

	State *DomainState `json:"state,omitempty,readonly"`
}

type DomainState struct {
	LastActionDateTime *time.Time `json:"lastActionDateTime,omitempty,readonly"`
	Operation          *string    `json:"operation,omitempty,readonly"`
	Status             *string    `json:"status,omitempty,readonly"`
}
