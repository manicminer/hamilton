package models

type Me struct {
	ID                *string `json:"id,readonly"`
	DisplayName       *string `json:"displayName,readonly"`
	UserPrincipalName *string `json:"userPrincipalName,readonly"`
}
