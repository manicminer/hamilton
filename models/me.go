package models

// Me describes the authenticated user.
type Me struct {
	ID                *string `json:"id,readonly"`
	DisplayName       *string `json:"displayName,readonly"`
	UserPrincipalName *string `json:"userPrincipalName,readonly"`
}
