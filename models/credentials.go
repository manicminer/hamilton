package models

import "time"

type KeyCredential struct {
	CustomKeyIdentifier *string    `json:"customKeyIdentifier,omitempty"`
	DisplayName         *string    `json:"displayName,omitempty"`
	EndDateTime         *time.Time `json:"endDateTime,omitempty,readonly"`
	KeyId               *string    `json:"keyId,omitempty"`
	StartDateTime       *time.Time `json:"startDateTime,omitempty,readonly"`
	Type                *string    `json:"type,omitempty"`
	Usage               *string    `json:"usage,omitempty"`
	Key                 *string    `json:"key,omitempty"`
}

type PasswordCredential struct {
	CustomKeyIdentifier *string    `json:"customKeyIdentifier,omitempty"`
	DisplayName         *string    `json:"displayName,omitempty"`
	EndDateTime         *time.Time `json:"endDateTime,omitempty,readonly"`
	Hint                *string    `json:"hint,omitempty"`
	KeyId               *string    `json:"keyId,omitempty"`
	SecretText          *string    `json:"secretText,omitempty"`
	StartDateTime       *time.Time `json:"startDateTime,omitempty,readonly"`
}
