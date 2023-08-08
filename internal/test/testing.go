package test

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-sdk/sdk/auth"
	"github.com/hashicorp/go-azure-sdk/sdk/claims"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/msgraph"
	"golang.org/x/oauth2"
)

func envDefault(envVarName, defaultValue string) string {
	if v := os.Getenv(envVarName); v != "" {
		return v
	}
	return defaultValue
}

var (
	defaultTenantId       = envDefault("DEFAULT_TENANT_ID", "6df54acb-f3cd-4734-85e3-7511ade57a02")
	defaultTenantDomain   = envDefault("DEFAULT_TENANT_DOMAIN", "hamiltontesting2.onmicrosoft.com")
	b2cTenantId           = envDefault("B2C_TENANT_ID", "1f54b558-9265-4964-a119-3eab591f05a9")
	b2cTenantDomain       = envDefault("B2C_TENANT_DOMAIN", "hamiltontestingb2c.onmicrosoft.com")
	connectedTenantId     = envDefault("CONNECTED_TENANT_ID", "df78c1c2-114a-425c-b711-091b30dea0ac")
	connectedTenantDomain = envDefault("CONNECTED_TENANT_DOMAIN", "hamiltontesting3.onmicrosoft.com")

	clientId              = envDefault("CLIENT_ID", "c072182f-ead2-4e94-aa7f-a5975558c945")
	clientCertificate     = os.Getenv("CLIENT_CERTIFICATE")
	clientCertificatePath = os.Getenv("CLIENT_CERTIFICATE_PATH")
	clientCertPassword    = os.Getenv("CLIENT_CERTIFICATE_PASSWORD")
	clientSecret          = os.Getenv("CLIENT_SECRET")
	environment           = envDefault("AZURE_ENVIRONMENT", "global")
	retryMax              = envDefault("RETRY_MAX", "14")
)

type Connection struct {
	AuthConfig *auth.Credentials
	Authorizer auth.Authorizer
	DomainName string
}

// NewConnection configures and returns a Connection for use in tests.
func NewConnection(tenantId, tenantDomain string) *Connection {
	env, err := environments.FromName(environment)
	if err != nil {
		log.Fatal(err)
	}

	t := Connection{
		AuthConfig: &auth.Credentials{
			Environment:               *env,
			TenantID:                  tenantId,
			ClientID:                  clientId,
			ClientCertificateData:     utils.Base64DecodeCertificate(clientCertificate),
			ClientCertificatePath:     clientCertificatePath,
			ClientCertificatePassword: clientCertPassword,
			ClientSecret:              clientSecret,
			EnableAuthenticatingUsingClientCertificate: true,
			EnableAuthenticatingUsingClientSecret:      true,
			EnableAuthenticatingUsingAzureCLI:          true,
		},
		DomainName: tenantDomain,
	}

	return &t
}

// Authorize configures an Authorizer for the Connection
func (c *Connection) Authorize(ctx context.Context, api environments.Api) {
	var err error
	c.Authorizer, err = auth.NewAuthorizerFromCredentials(ctx, *c.AuthConfig, api)
	if err != nil {
		log.Fatal(err)
	}
}

type Test struct {
	Context      context.Context
	CancelFunc   context.CancelFunc
	Connections  map[string]*Connection
	RandomString string

	Claims *claims.Claims
	Token  *oauth2.Token

	AccessPackageAssignmentPolicyClient                  *msgraph.AccessPackageAssignmentPolicyClient
	AccessPackageAssignmentRequestClient                 *msgraph.AccessPackageAssignmentRequestClient
	AccessPackageCatalogClient                           *msgraph.AccessPackageCatalogClient
	AccessPackageClient                                  *msgraph.AccessPackageClient
	AccessPackageResourceClient                          *msgraph.AccessPackageResourceClient
	AccessPackageResourceRequestClient                   *msgraph.AccessPackageResourceRequestClient
	AccessPackageResourceRoleScopeClient                 *msgraph.AccessPackageResourceRoleScopeClient
	AdministrativeUnitsClient                            *msgraph.AdministrativeUnitsClient
	ApplicationTemplatesClient                           *msgraph.ApplicationTemplatesClient
	ApplicationsClient                                   *msgraph.ApplicationsClient
	AppRoleAssignedToClient                              *msgraph.AppRoleAssignedToClient
	AuthenticationMethodsClient                          *msgraph.AuthenticationMethodsClient
	AuthenticationStrengthPoliciesClient                 *msgraph.AuthenticationStrengthPoliciesClient
	B2CUserFlowClient                                    *msgraph.B2CUserFlowClient
	ClaimsMappingPolicyClient                            *msgraph.ClaimsMappingPolicyClient
	ConditionalAccessPoliciesClient                      *msgraph.ConditionalAccessPoliciesClient
	ConnectedOrganizationClient                          *msgraph.ConnectedOrganizationClient
	DelegatedPermissionGrantsClient                      *msgraph.DelegatedPermissionGrantsClient
	DirectoryAuditReportsClient                          *msgraph.DirectoryAuditReportsClient
	DirectoryObjectsClient                               *msgraph.DirectoryObjectsClient
	DirectoryRoleTemplatesClient                         *msgraph.DirectoryRoleTemplatesClient
	DirectoryRolesClient                                 *msgraph.DirectoryRolesClient
	DomainsClient                                        *msgraph.DomainsClient
	EntitlementRoleAssignmentsClient                     *msgraph.EntitlementRoleAssignmentsClient
	EntitlementRoleDefinitionsClient                     *msgraph.EntitlementRoleDefinitionsClient
	GroupsAppRoleAssignmentsClient                       *msgraph.AppRoleAssignmentsClient
	GroupsClient                                         *msgraph.GroupsClient
	IdentityProvidersClient                              *msgraph.IdentityProvidersClient
	InvitationsClient                                    *msgraph.InvitationsClient
	MeClient                                             *msgraph.MeClient
	NamedLocationsClient                                 *msgraph.NamedLocationsClient
	PrivilegedAccessGroupClient                          *msgraph.PrivilegedAccessGroupClient
	PrivilegedAccessGroupAssignmentScheduleClient        *msgraph.PrivilegedAccessGroupAssignmentScheduleClient
	PrivilegedAccessGroupAssignmentScheduleRequestClient *msgraph.PrivilegedAccessGroupAssignmentScheduleRequestClient
	ReportsClient                                        *msgraph.ReportsClient
	RoleAssignmentsClient                                *msgraph.RoleAssignmentsClient
	RoleDefinitionsClient                                *msgraph.RoleDefinitionsClient
	RoleEligibilityScheduleRequestClient                 *msgraph.RoleEligibilityScheduleRequestClient
	SchemaExtensionsClient                               *msgraph.SchemaExtensionsClient
	ServicePrincipalsAppRoleAssignmentsClient            *msgraph.AppRoleAssignmentsClient
	ServicePrincipalsClient                              *msgraph.ServicePrincipalsClient
	SignInReportsClient                                  *msgraph.SignInReportsClient
	SynchronizationJobClient                             *msgraph.SynchronizationJobClient
	TermsOfUseAgreementClient                            *msgraph.TermsOfUseAgreementClient
	TokenIssuancePolicyClient                            *msgraph.TokenIssuancePolicyClient
	UserFlowAttributesClient                             *msgraph.UserFlowAttributesClient
	UsersAppRoleAssignmentsClient                        *msgraph.AppRoleAssignmentsClient
	UsersClient                                          *msgraph.UsersClient
	WindowsAutopilotDeploymentProfilesClient             *msgraph.WindowsAutopilotDeploymentProfilesClient
}

func NewTest(t *testing.T) (c *Test) {
	ctx := context.Background()
	var cancel context.CancelFunc

	if deadline, ok := t.Deadline(); ok {
		ctx, cancel = context.WithDeadline(ctx, deadline)
	} else {
		ctx, cancel = context.WithTimeout(ctx, 5*time.Minute)
	}

	c = &Test{
		Context:      ctx,
		CancelFunc:   cancel,
		Connections:  make(map[string]*Connection),
		RandomString: RandomString(),
	}

	conn := NewConnection(defaultTenantId, defaultTenantDomain)
	conn.Authorize(ctx, conn.AuthConfig.Environment.MicrosoftGraph)
	c.Connections["default"] = conn

	conn2 := NewConnection(b2cTenantId, b2cTenantDomain)
	conn2.Authorize(ctx, conn.AuthConfig.Environment.MicrosoftGraph)
	c.Connections["b2c"] = conn2

	conn3 := NewConnection(connectedTenantId, connectedTenantDomain)
	conn3.Authorize(ctx, conn.AuthConfig.Environment.MicrosoftGraph)
	c.Connections["connected"] = conn3

	var err error
	c.Token, err = conn.Authorizer.Token(ctx, &http.Request{})
	if err != nil {
		t.Fatalf("could not acquire access token: %v", err)
	}

	c.Claims, err = claims.ParseClaims(c.Token)
	if err != nil {
		t.Fatalf("could not parse claims: %v", err)
	}

	retry, err := strconv.Atoi(retryMax)
	if err != nil {
		t.Fatalf("invalid retry count %q: %v", retryMax, err)
	}

	endpoint, ok := c.Connections["default"].AuthConfig.Environment.MicrosoftGraph.Endpoint()
	if !ok {
		t.Fatalf("could not configure MS Graph endpoint for environment %q", c.Connections["default"].AuthConfig.Environment.Name)
	}

	c.AccessPackageAssignmentPolicyClient = msgraph.NewAccessPackageAssignmentPolicyClient()
	c.AccessPackageAssignmentPolicyClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.AccessPackageAssignmentPolicyClient.BaseClient.Endpoint = *endpoint
	c.AccessPackageAssignmentPolicyClient.BaseClient.RetryableClient.RetryMax = retry

	c.AccessPackageAssignmentRequestClient = msgraph.NewAccessPackageAssignmentRequestClient()
	c.AccessPackageAssignmentRequestClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.AccessPackageAssignmentRequestClient.BaseClient.Endpoint = *endpoint
	c.AccessPackageAssignmentRequestClient.BaseClient.RetryableClient.RetryMax = retry

	c.AccessPackageCatalogClient = msgraph.NewAccessPackageCatalogClient()
	c.AccessPackageCatalogClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.AccessPackageCatalogClient.BaseClient.Endpoint = *endpoint
	c.AccessPackageCatalogClient.BaseClient.RetryableClient.RetryMax = retry

	c.AccessPackageClient = msgraph.NewAccessPackageClient()
	c.AccessPackageClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.AccessPackageClient.BaseClient.Endpoint = *endpoint
	c.AccessPackageClient.BaseClient.RetryableClient.RetryMax = retry

	c.AccessPackageResourceClient = msgraph.NewAccessPackageResourceClient()
	c.AccessPackageResourceClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.AccessPackageResourceClient.BaseClient.Endpoint = *endpoint
	c.AccessPackageResourceClient.BaseClient.RetryableClient.RetryMax = retry

	c.AccessPackageResourceRequestClient = msgraph.NewAccessPackageResourceRequestClient()
	c.AccessPackageResourceRequestClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.AccessPackageResourceRequestClient.BaseClient.Endpoint = *endpoint
	c.AccessPackageAssignmentPolicyClient.BaseClient.RetryableClient.RetryMax = retry

	c.AccessPackageResourceRoleScopeClient = msgraph.NewAccessPackageResourceRoleScopeClient()
	c.AccessPackageResourceRoleScopeClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.AccessPackageResourceRoleScopeClient.BaseClient.Endpoint = *endpoint
	c.AccessPackageResourceRoleScopeClient.BaseClient.RetryableClient.RetryMax = retry

	c.AdministrativeUnitsClient = msgraph.NewAdministrativeUnitsClient()
	c.AdministrativeUnitsClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.AdministrativeUnitsClient.BaseClient.Endpoint = *endpoint
	c.AdministrativeUnitsClient.BaseClient.RetryableClient.RetryMax = retry

	c.ApplicationTemplatesClient = msgraph.NewApplicationTemplatesClient()
	c.ApplicationTemplatesClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.ApplicationTemplatesClient.BaseClient.Endpoint = *endpoint
	c.ApplicationTemplatesClient.BaseClient.RetryableClient.RetryMax = retry

	c.ApplicationsClient = msgraph.NewApplicationsClient()
	c.ApplicationsClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.ApplicationsClient.BaseClient.Endpoint = *endpoint
	c.ApplicationsClient.BaseClient.RetryableClient.RetryMax = retry
	c.ApplicationsClient.BaseClient.ApiVersion = msgraph.Version10

	c.AppRoleAssignedToClient = msgraph.NewAppRoleAssignedToClient()
	c.AppRoleAssignedToClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.AppRoleAssignedToClient.BaseClient.Endpoint = *endpoint
	c.AppRoleAssignedToClient.BaseClient.RetryableClient.RetryMax = retry

	c.AuthenticationMethodsClient = msgraph.NewAuthenticationMethodsClient()
	c.AuthenticationMethodsClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.AuthenticationMethodsClient.BaseClient.Endpoint = *endpoint
	c.AuthenticationMethodsClient.BaseClient.RetryableClient.RetryMax = retry

	c.AuthenticationStrengthPoliciesClient = msgraph.NewAuthenticationStrengthPoliciesClient()
	c.AuthenticationStrengthPoliciesClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.AuthenticationStrengthPoliciesClient.BaseClient.Endpoint = *endpoint
	c.AuthenticationStrengthPoliciesClient.BaseClient.RetryableClient.RetryMax = retry

	c.B2CUserFlowClient = msgraph.NewB2CUserFlowClient()
	c.B2CUserFlowClient.BaseClient.Authorizer = c.Connections["b2c"].Authorizer
	c.B2CUserFlowClient.BaseClient.Endpoint = *endpoint
	c.B2CUserFlowClient.BaseClient.RetryableClient.RetryMax = retry

	c.ClaimsMappingPolicyClient = msgraph.NewClaimsMappingPolicyClient()
	c.ClaimsMappingPolicyClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.ClaimsMappingPolicyClient.BaseClient.Endpoint = *endpoint
	c.ClaimsMappingPolicyClient.BaseClient.RetryableClient.RetryMax = retry

	c.ConditionalAccessPoliciesClient = msgraph.NewConditionalAccessPoliciesClient()
	c.ConditionalAccessPoliciesClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.ConditionalAccessPoliciesClient.BaseClient.Endpoint = *endpoint
	c.ConditionalAccessPoliciesClient.BaseClient.RetryableClient.RetryMax = retry

	c.ConnectedOrganizationClient = msgraph.NewConnectedOrganizationClient()
	c.ConnectedOrganizationClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.ConnectedOrganizationClient.BaseClient.Endpoint = *endpoint
	c.ConnectedOrganizationClient.BaseClient.RetryableClient.RetryMax = retry

	c.DelegatedPermissionGrantsClient = msgraph.NewDelegatedPermissionGrantsClient()
	c.DelegatedPermissionGrantsClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.DelegatedPermissionGrantsClient.BaseClient.Endpoint = *endpoint
	c.DelegatedPermissionGrantsClient.BaseClient.RetryableClient.RetryMax = retry

	c.DirectoryAuditReportsClient = msgraph.NewDirectoryAuditReportsClient()
	c.DirectoryAuditReportsClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.DirectoryAuditReportsClient.BaseClient.Endpoint = *endpoint
	c.DirectoryAuditReportsClient.BaseClient.RetryableClient.RetryMax = retry

	c.DirectoryObjectsClient = msgraph.NewDirectoryObjectsClient()
	c.DirectoryObjectsClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.DirectoryObjectsClient.BaseClient.Endpoint = *endpoint
	c.DirectoryObjectsClient.BaseClient.RetryableClient.RetryMax = retry

	c.DirectoryRoleTemplatesClient = msgraph.NewDirectoryRoleTemplatesClient()
	c.DirectoryRoleTemplatesClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.DirectoryRoleTemplatesClient.BaseClient.Endpoint = *endpoint
	c.DirectoryRoleTemplatesClient.BaseClient.RetryableClient.RetryMax = retry

	c.DirectoryRolesClient = msgraph.NewDirectoryRolesClient()
	c.DirectoryRolesClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.DirectoryRolesClient.BaseClient.Endpoint = *endpoint
	c.DirectoryRolesClient.BaseClient.RetryableClient.RetryMax = retry

	c.DomainsClient = msgraph.NewDomainsClient()
	c.DomainsClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.DomainsClient.BaseClient.Endpoint = *endpoint
	c.DomainsClient.BaseClient.RetryableClient.RetryMax = retry

	c.EntitlementRoleAssignmentsClient = msgraph.NewEntitlementRoleAssignmentsClient()
	c.EntitlementRoleAssignmentsClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.EntitlementRoleAssignmentsClient.BaseClient.Endpoint = *endpoint
	c.EntitlementRoleAssignmentsClient.BaseClient.RetryableClient.RetryMax = retry

	c.EntitlementRoleDefinitionsClient = msgraph.NewEntitlementRoleDefinitionsClient()
	c.EntitlementRoleDefinitionsClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.EntitlementRoleDefinitionsClient.BaseClient.Endpoint = *endpoint
	c.EntitlementRoleDefinitionsClient.BaseClient.RetryableClient.RetryMax = retry

	c.GroupsAppRoleAssignmentsClient = msgraph.NewGroupsAppRoleAssignmentsClient()
	c.GroupsAppRoleAssignmentsClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.GroupsAppRoleAssignmentsClient.BaseClient.Endpoint = *endpoint
	c.GroupsAppRoleAssignmentsClient.BaseClient.RetryableClient.RetryMax = retry

	c.GroupsClient = msgraph.NewGroupsClient()
	c.GroupsClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.GroupsClient.BaseClient.Endpoint = *endpoint
	c.GroupsClient.BaseClient.RetryableClient.RetryMax = retry

	c.IdentityProvidersClient = msgraph.NewIdentityProvidersClient()
	c.IdentityProvidersClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.IdentityProvidersClient.BaseClient.Endpoint = *endpoint
	c.IdentityProvidersClient.BaseClient.RetryableClient.RetryMax = retry

	c.InvitationsClient = msgraph.NewInvitationsClient()
	c.InvitationsClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.InvitationsClient.BaseClient.Endpoint = *endpoint
	c.InvitationsClient.BaseClient.RetryableClient.RetryMax = retry

	c.MeClient = msgraph.NewMeClient()
	c.MeClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.MeClient.BaseClient.Endpoint = *endpoint
	c.MeClient.BaseClient.RetryableClient.RetryMax = retry

	c.NamedLocationsClient = msgraph.NewNamedLocationsClient()
	c.NamedLocationsClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.NamedLocationsClient.BaseClient.Endpoint = *endpoint
	c.NamedLocationsClient.BaseClient.RetryableClient.RetryMax = retry

	c.PrivilegedAccessGroupClient = msgraph.NewPrivilegedAccessGroupClient()
	c.PrivilegedAccessGroupClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.PrivilegedAccessGroupClient.BaseClient.Endpoint = *endpoint
	c.PrivilegedAccessGroupClient.BaseClient.RetryableClient.RetryMax = retry

	c.ReportsClient = msgraph.NewReportsClient()
	c.ReportsClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.ReportsClient.BaseClient.Endpoint = *endpoint
	c.ReportsClient.BaseClient.RetryableClient.RetryMax = retry

	c.RoleAssignmentsClient = msgraph.NewRoleAssignmentsClient()
	c.RoleAssignmentsClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.RoleAssignmentsClient.BaseClient.Endpoint = *endpoint
	c.RoleAssignmentsClient.BaseClient.RetryableClient.RetryMax = retry

	c.RoleDefinitionsClient = msgraph.NewRoleDefinitionsClient()
	c.RoleDefinitionsClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.RoleDefinitionsClient.BaseClient.Endpoint = *endpoint
	c.RoleDefinitionsClient.BaseClient.RetryableClient.RetryMax = retry

	c.RoleEligibilityScheduleRequestClient = msgraph.NewRoleEligibilityScheduleRequestClient()
	c.RoleEligibilityScheduleRequestClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.RoleEligibilityScheduleRequestClient.BaseClient.Endpoint = *endpoint
	c.RoleEligibilityScheduleRequestClient.BaseClient.RetryableClient.RetryMax = retry

	c.SchemaExtensionsClient = msgraph.NewSchemaExtensionsClient()
	c.SchemaExtensionsClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.SchemaExtensionsClient.BaseClient.Endpoint = *endpoint
	c.SchemaExtensionsClient.BaseClient.RetryableClient.RetryMax = retry

	c.ServicePrincipalsAppRoleAssignmentsClient = msgraph.NewServicePrincipalsAppRoleAssignmentsClient()
	c.ServicePrincipalsAppRoleAssignmentsClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.ServicePrincipalsAppRoleAssignmentsClient.BaseClient.Endpoint = *endpoint
	c.ServicePrincipalsAppRoleAssignmentsClient.BaseClient.RetryableClient.RetryMax = retry

	c.ServicePrincipalsClient = msgraph.NewServicePrincipalsClient()
	c.ServicePrincipalsClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.ServicePrincipalsClient.BaseClient.Endpoint = *endpoint
	c.ServicePrincipalsClient.BaseClient.RetryableClient.RetryMax = retry
	c.ServicePrincipalsClient.BaseClient.ApiVersion = msgraph.Version10

	c.SignInReportsClient = msgraph.NewSignInReportsClient()
	c.SignInReportsClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.SignInReportsClient.BaseClient.Endpoint = *endpoint
	c.SignInReportsClient.BaseClient.RetryableClient.RetryMax = retry

	c.SynchronizationJobClient = msgraph.NewSynchronizationJobClient()
	c.SynchronizationJobClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.SynchronizationJobClient.BaseClient.Endpoint = *endpoint
	c.SynchronizationJobClient.BaseClient.RetryableClient.RetryMax = retry

	c.TermsOfUseAgreementClient = msgraph.NewTermsOfUseAgreementClient()
	c.TermsOfUseAgreementClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.TermsOfUseAgreementClient.BaseClient.Endpoint = *endpoint
	c.TermsOfUseAgreementClient.BaseClient.RetryableClient.RetryMax = retry

	c.TokenIssuancePolicyClient = msgraph.NewTokenIssuancePolicyClient()
	c.TokenIssuancePolicyClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.TokenIssuancePolicyClient.BaseClient.Endpoint = *endpoint
	c.TokenIssuancePolicyClient.BaseClient.RetryableClient.RetryMax = retry

	c.UserFlowAttributesClient = msgraph.NewUserFlowAttributesClient()
	c.UserFlowAttributesClient.BaseClient.Authorizer = c.Connections["b2c"].Authorizer
	c.UserFlowAttributesClient.BaseClient.Endpoint = *endpoint
	c.UserFlowAttributesClient.BaseClient.RetryableClient.RetryMax = retry

	c.UsersAppRoleAssignmentsClient = msgraph.NewUsersAppRoleAssignmentsClient()
	c.UsersAppRoleAssignmentsClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.UsersAppRoleAssignmentsClient.BaseClient.Endpoint = *endpoint
	c.UsersAppRoleAssignmentsClient.BaseClient.RetryableClient.RetryMax = retry

	c.UsersClient = msgraph.NewUsersClient()
	c.UsersClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.UsersClient.BaseClient.Endpoint = *endpoint
	c.UsersClient.BaseClient.RetryableClient.RetryMax = retry

	c.WindowsAutopilotDeploymentProfilesClient = msgraph.NewWindowsAutopilotDeploymentProfilesClient()
	c.WindowsAutopilotDeploymentProfilesClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.WindowsAutopilotDeploymentProfilesClient.BaseClient.Endpoint = *endpoint
	c.WindowsAutopilotDeploymentProfilesClient.BaseClient.RetryableClient.RetryMax = retry

	return
}
