package models

import (
	goerrors "errors"
	"fmt"
	"time"

	"github.com/manicminer/hamilton/errors"
)

type Application struct {
	ID                         *string                              `json:"id,omitempty,readonly"`
	AddIns                     *[]ApplicationAddIn                  `json:"addIns,omitempty"`
	Api                        *ApplicationApi                      `json:"api,omitempty"`
	AppId                      *string                              `json:"appId,omitempty"`
	AppRoles                   *[]ApplicationAppRole                `json:"appRoles,omitempty"`
	CreatedDateTime            *time.Time                           `json:"createdDateTime,omitempty,readonly"`
	DeletedDateTime            *time.Time                           `json:"deletedDateTime,omitempty,readonly"`
	DisplayName                *string                              `json:"displayName,omitempty"`
	GroupMembershipClaims      *string                              `json:"groupMembershipClaims,omitempty"`
	IdentifierUris             *[]string                            `json:"identifierUris,omitempty"`
	Info                       *ApplicationInformationalUrl         `json:"info,omitempty"`
	IsFallbackPublicClient     *bool                                `json:"isFallbackPublicCLient,omitempty"`
	KeyCredentials             *[]KeyCredential                     `json:"keyCredentials,omitempty"`
	Oauth2RequiredPostResponse *bool                                `json:"oauth2RequiredPostResponse,omitempty"`
	OnPremisesPublishing       *ApplicationOnPremisesPublishing     `json:"onPremisePublishing,omitempty"`
	OptionalClaims             *ApplicationOptionalClaims           `json:"optionalClaims,omitempty"`
	ParentalControlSettings    *ParentalControlSettings             `json:"parentalControlSettings,omitempty"`
	PasswordCredentials        *[]PasswordCredential                `json:"passwordCredentials,omitempty"`
	PublicClient               *ApplicationPublicClient             `json:"publicClient,omitempty"`
	PublisherDomain            *string                              `json:"publisherDomain,omitempty"`
	RequiredResourceAccess     *[]ApplicationRequiredResourceAccess `json:"requiredResourceAccess,omitempty"`
	SignInAudience             ApplicationSignInAudience            `json:"signInAudience,omitempty"`
	Tags                       *[]string                            `json:"tags,omitempty"`
	TokenEncryptionKeyId       *string                              `json:"tokenEncryptionKeyId,omitempty"`
	Web                        *ApplicationWeb                      `json:"web,omitempty"`

	Owners *[]string `json:"owners@odata.bind,omitempty"`
}

func (a *Application) AppendOwner(endpoint string, apiVersion string, id string) {
	val := fmt.Sprintf("%s/%s/directoryObjects/%s", endpoint, apiVersion, id)
	var owners []string
	if a.Owners != nil {
		owners = *a.Owners
	}
	owners = append(owners, val)
	a.Owners = &owners
}

func (a *Application) AppendAppRole(role ApplicationAppRole) error {
	if role.ID == nil {
		return goerrors.New("ID of new role is nil")
	}

	cap := 1
	if a.AppRoles != nil {
		cap += len(*a.AppRoles)
	}

	newRoles := make([]ApplicationAppRole, 1, cap)
	newRoles[0] = role

	for _, v := range *a.AppRoles {
		if v.ID != nil && *v.ID == *role.ID {
			return &errors.AlreadyExistsError{"App Role", *role.ID}
		}
		newRoles = append(newRoles, v)
	}

	a.AppRoles = &newRoles
	return nil
}

func (a *Application) RemoveAppRole(role ApplicationAppRole) error {
	if role.ID == nil {
		return goerrors.New("ID of role is nil")
	}

	if a.AppRoles == nil {
		return goerrors.New("no roles to remove")
	}

	appRoles := make([]ApplicationAppRole, 0, len(*a.AppRoles))
	for _, v := range *a.AppRoles {
		if v.ID == nil || *v.ID != *role.ID {
			appRoles = append(appRoles, v)
		}
	}

	if len(appRoles) == len(*a.AppRoles) {
		return goerrors.New("could not find role to remove")
	}

	a.AppRoles = &appRoles
	return nil
}

func (a *Application) UpdateAppRole(role ApplicationAppRole) error {
	if role.ID == nil {
		return goerrors.New("ID of role is nil")
	}

	if a.AppRoles == nil {
		return goerrors.New("no roles to update")
	}

	appRoles := *a.AppRoles
	for i, v := range appRoles {
		if v.ID != nil && *v.ID == *role.ID {
			appRoles[i] = role
			break
		}
	}

	a.AppRoles = &appRoles
	return nil
}

type ApplicationAddIn struct {
	ID         *string          `json:"id,omitempty"`
	Properties *[]AddInKeyValue `json:"properties,omitempty"`
	Type       *string          `json:"type,omitempty"`
}

type AddInKeyValue struct {
	Key   *string `json:"key,omitempty"`
	Value *string `json:"value,omitempty"`
}

type ApplicationApi struct {
	AcceptMappedClaims          *bool                                     `json:"acceptMappedClaims,omitempty"`
	KnownClientApplications     *[]string                                 `json:"knownClientApplications,omitempty"`
	OAuth2PermissionScopes      *[]ApplicationApiPermissionScope          `json:"oauth2PermissionScopes,omitempty"`
	PreAuthorizedApplications   *[]ApplicationApiPreAuthorizedApplication `json:"preAuthorizedApplications,omitempty"`
	RequestedAccessTokenVersion *int32                                    `json:"requestedAccessTokenVersion,omitempty"`
}

type ApplicationApiPermissionScope struct {
	ID                      *string `json:"id,omitempty"`
	AdminConsentDescription *string `json:"adminConsentDescription,omitempty"`
	AdminConsentDisplayName *string `json:"adminConsentDisplayName,omitempty"`
	IsEnabled               *bool   `json:"isEnabled,omitempty"`
	Type                    *string `json:"type,omitempty"`
	UserConsentDescription  *string `json:"userConsentDescription,omitempty"`
	UserConsentDisplayName  *string `json:"userConsentDisplayName,omitempty"`
	Value                   *string `json:"value,omitempty"`
}

type ApplicationApiPreAuthorizedApplication struct {
	AppId         *string   `json:"appId,omitempty"`
	PermissionIds *[]string `json:"permissionIds,omitempty"`
}

type ApplicationAppRole struct {
	ID                 *string   `json:"id,omitempty"`
	AllowedMemberTypes *[]string `json:"allowedMemberTypes,omitempty"`
	Description        *string   `json:"description,omitempty"`
	DisplayName        *string   `json:"displayName,omitempty"`
	IsEnabled          *bool     `json:"isEnabled,omitempty"`
	Origin             *string   `json:"origin,omitempty"`
	Value              *string   `json:"value,omitempty"`
}

type ApplicationImplicitGrantSettings struct {
	EnableAccessTokenIssuance *bool `json:"enableAccessTokenIssuance,omitempty"`
	EnableIdTokenIssuance     *bool `json:"enableIdTokenIssuance,omitempty"`
}

type ApplicationInformationalUrl struct {
	LogoUrl             *string `json:"logoUrl,omitempty`
	MarketingUrl        *string `json:"marketingUrl"`
	PrivacyStatementUrl *string `json:"privacyStatementUrl"`
	SupportUrl          *string `json:"supportUrl"`
	TermsOfServiceUrl   *string `json:"termsOfServiceUrl"`
}

type ApplicationKerberosSignOnSettings struct {
	ServicePrincipalName       *string `json:"kerberosServicePrincipalName,omitempty"`
	SignOnMappingAttributeType *string `jsonL:"kerberosSignOnMappingAttributeType,omitempty"`
}

type ApplicationOnPremisesPublishing struct {
	AlternateUrl                  *string `json:"alternateUrl,omitempty"`
	ApplicationServerTimeout      *string `json:"applicationServerTimeout,omitempty"`
	ApplicationType               *string `json:"applicationType,omitempty"`
	ExternalAuthenticationType    *string `json:"externalAuthenticationType,omitempty"`
	ExternalUrl                   *string `json:"externalUrl,omitempty"`
	InternalUrl                   *string `json:"internalUrl,omitempty"`
	IsHttpOnlyCookieEnabled       *bool   `json:"isHttpOnlyCookieEnabled,omitempty"`
	IsOnPremPublishingEnabled     *bool   `json:"isOnPremPublishingEnabled,omitempty"`
	IsPersistentCookieEnabled     *bool   `json:"isPersistentCookieEnabled,omitempty"`
	IsSecureCookieEnabled         *bool   `json:"isSecureCookieEnabled,omitempty"`
	IsTranslateHostHeaderEnabled  *bool   `json:"isTranslateHostHeaderEnabled,omitempty"`
	IsTranslateLinksInBodyEnabled *bool   `json:"isTranslateLinksInBodyEnabled,omitempty"`

	SingleSignOnSettings                     *ApplicationOnPremisesPublishingSingleSignOn                             `json:"singleSignOnSettings,omitempty"`
	VerifiedCustomDomainCertificatesMetadata *ApplicationOnPremisesPublishingVerifiedCustomDomainCertificatesMetadata `json:"verifiedCustomDomainCertificatesMetadata,omitempty"`
	VerifiedCustomDomainKeyCredential        *KeyCredential                                                           `json:"verifiedCustomDomainKeyCredential,omitempty"`
	VerifiedCustomDomainPasswordCredential   *PasswordCredential                                                      `json:"verifiedCustomDomainPasswordCredential,omitempty"`
}

type ApplicationOnPremisesPublishingSingleSignOn struct {
	KerberosSignOnSettings *ApplicationKerberosSignOnSettings `json:"kerberosSignOnSettings,omitempty"`
	SingleSignOnMode       *string                            `json:"singleSignOnMode,omitempty"`
}

type ApplicationOnPremisesPublishingVerifiedCustomDomainCertificatesMetadata struct {
	ExpiryDate  *time.Time `json:"expiryDate,omitempty,readonly"`
	IssueDate   *time.Time `json:"issueDate,omitempty,readonly"`
	IssuerName  *string    `json:"issuerName,omitempty"`
	SubjectName *string    `json:"subjectName,omitempty"`
	Thumbprint  *string    `json:"thumbprint,omitempty"`
}

type ApplicationOptionalClaim struct {
	AdditionalProperties *[]string `json:"additionalProperties,omitempty"`
	Essential            *bool     `json:"essential,omitempty"`
	Name                 *string   `json:"name,omitempty"`
	Source               *string   `json:"source,omitempty"`
}

type ApplicationOptionalClaims struct {
	AccessToken *[]ApplicationOptionalClaim `json:"accessToken,omitempty"`
	IdToken     *[]ApplicationOptionalClaim `json:"idToken,omitempty"`
	Saml2Token  *[]ApplicationOptionalClaim `json:"saml2Token,omitempty"`
}

type ApplicationPublicClient struct {
	RedirectUris *[]string `json:"redirectUris,omitempty"`
}

type ApplicationRequiredResourceAccess struct {
	ResourceAccess *[]ApplicationResourceAccess `json:"resourceAccess,omitempty"`
	ResourceAppId  *string                      `json:"resourceAppId,omitempty"`
}

type ApplicationResourceAccess struct {
	ID   *string `json:"id,omitempty"`
	Type *string `json:"type,omitempty"`
}

type ApplicationSignInAudience string

const (
	SignInAudienceAzureADMyOrg                       ApplicationSignInAudience = "AzureADMyOrg"
	SignInAudienceAzureADMultipleOrgs                ApplicationSignInAudience = "AzureADMultipleOrgs"
	SignInAudienceAzureADandPersonalMicrosoftAccount ApplicationSignInAudience = "AzureADandPersonalMicrosoftAccount"
)

type ApplicationWeb struct {
	HomePageUrl           *string                           `json:"homePageUrl"`
	ImplicitGrantSettings *ApplicationImplicitGrantSettings `json:"implicitGrantSettings,omitempty"`
	LogoutUrl             *string                           `json:"logoutUrl"`
	RedirectUris          *[]string                         `json:"redirectUris,omitempty"`
}

type ParentalControlSettings struct {
	CountriesBlockedForMinors *[]string `json:"countriesBlockedForMinors,omitempty"`
	LegalAgeGroupRule         *string   `json:"legalAgeGroupRule,omitempty"`
}
