package test

import (
	"context"
	"log"
	"os"
	"strconv"
	"testing"

	"golang.org/x/oauth2"

	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/environments"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/msgraph"
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
	environment           = os.Getenv("AZURE_ENVIRONMENT")
	idTokenRequestUrl     = os.Getenv("ACTIONS_ID_TOKEN_REQUEST_URL")
	idTokenRequestToken   = os.Getenv("ACTIONS_ID_TOKEN_REQUEST_TOKEN")
	retryMax              = envDefault("RETRY_MAX", "14")
)

type Connection struct {
	AuthConfig *auth.Config
	Authorizer auth.Authorizer
	DomainName string
}

// NewConnection configures and returns a Connection for use in tests.
func NewConnection(tokenVersion auth.TokenVersion, tenantId, tenantDomain string) *Connection {
	env, err := environments.EnvironmentFromString(environment)
	if err != nil {
		log.Fatal(err)
	}

	t := Connection{
		AuthConfig: &auth.Config{
			Environment:            env,
			Version:                tokenVersion,
			TenantID:               tenantId,
			ClientID:               clientId,
			ClientCertData:         utils.Base64DecodeCertificate(clientCertificate),
			ClientCertPath:         clientCertificatePath,
			ClientCertPassword:     clientCertPassword,
			ClientSecret:           clientSecret,
			IDTokenRequestURL:      idTokenRequestUrl,
			IDTokenRequestToken:    idTokenRequestToken,
			EnableClientCertAuth:   true,
			EnableClientSecretAuth: true,
			EnableAzureCliToken:    true,
			EnableGitHubOIDCAuth:   true,
		},
		DomainName: tenantDomain,
	}

	return &t
}

// Authorize configures an Authorizer for the Connection
func (c *Connection) Authorize(ctx context.Context, api environments.Api) {
	var err error
	c.Authorizer, err = c.AuthConfig.NewAuthorizer(ctx, api)
	if err != nil {
		log.Fatal(err)
	}
}

type Test struct {
	Context      context.Context
	CancelFunc   context.CancelFunc
	Connections  map[string]*Connection
	RandomString string

	Claims auth.Claims
	Token  *oauth2.Token

	AccessPackageAssignmentPolicyClient       *msgraph.AccessPackageAssignmentPolicyClient
	AccessPackageCatalogClient                *msgraph.AccessPackageCatalogClient
	AccessPackageClient                       *msgraph.AccessPackageClient
	AccessPackageResourceClient               *msgraph.AccessPackageResourceClient
	AccessPackageResourceRequestClient        *msgraph.AccessPackageResourceRequestClient
	AccessPackageResourceRoleScopeClient      *msgraph.AccessPackageResourceRoleScopeClient
	AdministrativeUnitsClient                 *msgraph.AdministrativeUnitsClient
	ApplicationTemplatesClient                *msgraph.ApplicationTemplatesClient
	ApplicationsClient                        *msgraph.ApplicationsClient
	AppRoleAssignedToClient                   *msgraph.AppRoleAssignedToClient
	AuthenticationMethodsClient               *msgraph.AuthenticationMethodsClient
	B2CUserFlowClient                         *msgraph.B2CUserFlowClient
	ClaimsMappingPolicyClient                 *msgraph.ClaimsMappingPolicyClient
	ConditionalAccessPoliciesClient           *msgraph.ConditionalAccessPoliciesClient
	ConnectedOrganizationClient               *msgraph.ConnectedOrganizationClient
	DelegatedPermissionGrantsClient           *msgraph.DelegatedPermissionGrantsClient
	DirectoryAuditReportsClient               *msgraph.DirectoryAuditReportsClient
	DirectoryObjectsClient                    *msgraph.DirectoryObjectsClient
	DirectoryRoleTemplatesClient              *msgraph.DirectoryRoleTemplatesClient
	DirectoryRolesClient                      *msgraph.DirectoryRolesClient
	DomainsClient                             *msgraph.DomainsClient
	GroupsAppRoleAssignmentsClient            *msgraph.AppRoleAssignmentsClient
	GroupsClient                              *msgraph.GroupsClient
	IdentityProvidersClient                   *msgraph.IdentityProvidersClient
	InvitationsClient                         *msgraph.InvitationsClient
	MeClient                                  *msgraph.MeClient
	NamedLocationsClient                      *msgraph.NamedLocationsClient
	ReportsClient                             *msgraph.ReportsClient
	RoleAssignmentsClient                     *msgraph.RoleAssignmentsClient
	RoleDefinitionsClient                     *msgraph.RoleDefinitionsClient
	SchemaExtensionsClient                    *msgraph.SchemaExtensionsClient
	ServicePrincipalsAppRoleAssignmentsClient *msgraph.AppRoleAssignmentsClient
	ServicePrincipalsClient                   *msgraph.ServicePrincipalsClient
	SignInReportsClient                       *msgraph.SignInReportsClient
	SynchronizationJobClient                  *msgraph.SynchronizationJobClient
	UserFlowAttributesClient                  *msgraph.UserFlowAttributesClient
	UsersAppRoleAssignmentsClient             *msgraph.AppRoleAssignmentsClient
	UsersClient                               *msgraph.UsersClient
}

func NewTest(t *testing.T) (c *Test) {
	ctx := context.Background()
	var cancel context.CancelFunc = func() {}

	if deadline, ok := t.Deadline(); ok {
		ctx, cancel = context.WithDeadline(context.Background(), deadline)
	}

	c = &Test{
		Context:      ctx,
		CancelFunc:   cancel,
		Connections:  make(map[string]*Connection),
		RandomString: RandomString(),
	}

	conn := NewConnection(auth.TokenVersion2, defaultTenantId, defaultTenantDomain)
	conn.Authorize(ctx, conn.AuthConfig.Environment.MsGraph)
	c.Connections["default"] = conn

	conn2 := NewConnection(auth.TokenVersion2, b2cTenantId, b2cTenantDomain)
	conn2.Authorize(ctx, conn.AuthConfig.Environment.MsGraph)
	c.Connections["b2c"] = conn2

	conn3 := NewConnection(auth.TokenVersion2, connectedTenantId, connectedTenantDomain)
	conn3.Authorize(ctx, conn.AuthConfig.Environment.MsGraph)
	c.Connections["connected"] = conn3

	var err error
	c.Token, err = conn.Authorizer.Token()
	if err != nil {
		t.Fatalf("could not acquire access token: %v", err)
	}

	c.Claims, err = auth.ParseClaims(c.Token)
	if err != nil {
		t.Fatalf("could not parse claims: %v", err)
	}

	retry, err := strconv.Atoi(retryMax)
	if err != nil {
		t.Fatalf("invalid retry count %q: %v", retryMax, err)
	}

	c.AccessPackageAssignmentPolicyClient = msgraph.NewAccessPackageAssignmentPolicyClient(c.Connections["default"].AuthConfig.TenantID)
	c.AccessPackageAssignmentPolicyClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.AccessPackageAssignmentPolicyClient.BaseClient.Endpoint = c.Connections["default"].AuthConfig.Environment.MsGraph.Endpoint
	c.AccessPackageAssignmentPolicyClient.BaseClient.RetryableClient.RetryMax = retry

	c.AccessPackageCatalogClient = msgraph.NewAccessPackageCatalogClient(c.Connections["default"].AuthConfig.TenantID)
	c.AccessPackageCatalogClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.AccessPackageCatalogClient.BaseClient.Endpoint = c.Connections["default"].AuthConfig.Environment.MsGraph.Endpoint
	c.AccessPackageCatalogClient.BaseClient.RetryableClient.RetryMax = retry

	c.AccessPackageClient = msgraph.NewAccessPackageClient(c.Connections["default"].AuthConfig.TenantID)
	c.AccessPackageClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.AccessPackageClient.BaseClient.Endpoint = c.Connections["default"].AuthConfig.Environment.MsGraph.Endpoint
	c.AccessPackageClient.BaseClient.RetryableClient.RetryMax = retry

	c.AccessPackageResourceClient = msgraph.NewAccessPackageResourceClient(c.Connections["default"].AuthConfig.TenantID)
	c.AccessPackageResourceClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.AccessPackageResourceClient.BaseClient.Endpoint = c.Connections["default"].AuthConfig.Environment.MsGraph.Endpoint
	c.AccessPackageResourceClient.BaseClient.RetryableClient.RetryMax = retry

	c.AccessPackageResourceRequestClient = msgraph.NewAccessPackageResourceRequestClient(c.Connections["default"].AuthConfig.TenantID)
	c.AccessPackageResourceRequestClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.AccessPackageResourceRequestClient.BaseClient.Endpoint = c.Connections["default"].AuthConfig.Environment.MsGraph.Endpoint
	c.AccessPackageAssignmentPolicyClient.BaseClient.RetryableClient.RetryMax = retry

	c.AccessPackageResourceRoleScopeClient = msgraph.NewAccessPackageResourceRoleScopeClient(c.Connections["default"].AuthConfig.TenantID)
	c.AccessPackageResourceRoleScopeClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.AccessPackageResourceRoleScopeClient.BaseClient.Endpoint = c.Connections["default"].AuthConfig.Environment.MsGraph.Endpoint
	c.AccessPackageResourceRoleScopeClient.BaseClient.RetryableClient.RetryMax = retry

	c.AdministrativeUnitsClient = msgraph.NewAdministrativeUnitsClient(c.Connections["default"].AuthConfig.TenantID)
	c.AdministrativeUnitsClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.AdministrativeUnitsClient.BaseClient.Endpoint = c.Connections["default"].AuthConfig.Environment.MsGraph.Endpoint
	c.AdministrativeUnitsClient.BaseClient.RetryableClient.RetryMax = retry

	c.ApplicationTemplatesClient = msgraph.NewApplicationTemplatesClient(c.Connections["default"].AuthConfig.TenantID)
	c.ApplicationTemplatesClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.ApplicationTemplatesClient.BaseClient.Endpoint = c.Connections["default"].AuthConfig.Environment.MsGraph.Endpoint
	c.ApplicationTemplatesClient.BaseClient.RetryableClient.RetryMax = retry

	c.ApplicationsClient = msgraph.NewApplicationsClient(c.Connections["default"].AuthConfig.TenantID)
	c.ApplicationsClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.ApplicationsClient.BaseClient.Endpoint = c.Connections["default"].AuthConfig.Environment.MsGraph.Endpoint
	c.ApplicationsClient.BaseClient.RetryableClient.RetryMax = retry
	c.ApplicationsClient.BaseClient.ApiVersion = msgraph.Version10

	c.AppRoleAssignedToClient = msgraph.NewAppRoleAssignedToClient(c.Connections["default"].AuthConfig.TenantID)
	c.AppRoleAssignedToClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.AppRoleAssignedToClient.BaseClient.Endpoint = c.Connections["default"].AuthConfig.Environment.MsGraph.Endpoint
	c.AppRoleAssignedToClient.BaseClient.RetryableClient.RetryMax = retry

	c.AuthenticationMethodsClient = msgraph.NewAuthenticationMethodsClient(c.Connections["default"].AuthConfig.TenantID)
	c.AuthenticationMethodsClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.AuthenticationMethodsClient.BaseClient.Endpoint = c.Connections["default"].AuthConfig.Environment.MsGraph.Endpoint
	c.AuthenticationMethodsClient.BaseClient.RetryableClient.RetryMax = retry

	c.B2CUserFlowClient = msgraph.NewB2CUserFlowClient(c.Connections["b2c"].AuthConfig.TenantID)
	c.B2CUserFlowClient.BaseClient.Authorizer = c.Connections["b2c"].Authorizer
	c.B2CUserFlowClient.BaseClient.Endpoint = c.Connections["b2c"].AuthConfig.Environment.MsGraph.Endpoint
	c.B2CUserFlowClient.BaseClient.RetryableClient.RetryMax = retry

	c.ClaimsMappingPolicyClient = msgraph.NewClaimsMappingPolicyClient(c.Connections["default"].AuthConfig.TenantID)
	c.ClaimsMappingPolicyClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.ClaimsMappingPolicyClient.BaseClient.Endpoint = c.Connections["default"].AuthConfig.Environment.MsGraph.Endpoint
	c.ClaimsMappingPolicyClient.BaseClient.RetryableClient.RetryMax = retry

	c.ConditionalAccessPoliciesClient = msgraph.NewConditionalAccessPoliciesClient(c.Connections["default"].AuthConfig.TenantID)
	c.ConditionalAccessPoliciesClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.ConditionalAccessPoliciesClient.BaseClient.Endpoint = c.Connections["default"].AuthConfig.Environment.MsGraph.Endpoint
	c.ConditionalAccessPoliciesClient.BaseClient.RetryableClient.RetryMax = retry

	c.ConnectedOrganizationClient = msgraph.NewConnectedOrganizationClient(c.Connections["default"].AuthConfig.TenantID)
	c.ConnectedOrganizationClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.ConnectedOrganizationClient.BaseClient.Endpoint = c.Connections["default"].AuthConfig.Environment.MsGraph.Endpoint
	c.ConnectedOrganizationClient.BaseClient.RetryableClient.RetryMax = retry

	c.DelegatedPermissionGrantsClient = msgraph.NewDelegatedPermissionGrantsClient(c.Connections["default"].AuthConfig.TenantID)
	c.DelegatedPermissionGrantsClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.DelegatedPermissionGrantsClient.BaseClient.Endpoint = c.Connections["default"].AuthConfig.Environment.MsGraph.Endpoint
	c.DelegatedPermissionGrantsClient.BaseClient.RetryableClient.RetryMax = retry

	c.DirectoryAuditReportsClient = msgraph.NewDirectoryAuditReportsClient(c.Connections["default"].AuthConfig.TenantID)
	c.DirectoryAuditReportsClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.DirectoryAuditReportsClient.BaseClient.Endpoint = c.Connections["default"].AuthConfig.Environment.MsGraph.Endpoint
	c.DirectoryAuditReportsClient.BaseClient.RetryableClient.RetryMax = retry

	c.DirectoryObjectsClient = msgraph.NewDirectoryObjectsClient(c.Connections["default"].AuthConfig.TenantID)
	c.DirectoryObjectsClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.DirectoryObjectsClient.BaseClient.Endpoint = c.Connections["default"].AuthConfig.Environment.MsGraph.Endpoint
	c.DirectoryObjectsClient.BaseClient.RetryableClient.RetryMax = retry

	c.DirectoryRoleTemplatesClient = msgraph.NewDirectoryRoleTemplatesClient(c.Connections["default"].AuthConfig.TenantID)
	c.DirectoryRoleTemplatesClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.DirectoryRoleTemplatesClient.BaseClient.Endpoint = c.Connections["default"].AuthConfig.Environment.MsGraph.Endpoint
	c.DirectoryRoleTemplatesClient.BaseClient.RetryableClient.RetryMax = retry

	c.DirectoryRolesClient = msgraph.NewDirectoryRolesClient(c.Connections["default"].AuthConfig.TenantID)
	c.DirectoryRolesClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.DirectoryRolesClient.BaseClient.Endpoint = c.Connections["default"].AuthConfig.Environment.MsGraph.Endpoint
	c.DirectoryRolesClient.BaseClient.RetryableClient.RetryMax = retry

	c.DomainsClient = msgraph.NewDomainsClient(c.Connections["default"].AuthConfig.TenantID)
	c.DomainsClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.DomainsClient.BaseClient.Endpoint = c.Connections["default"].AuthConfig.Environment.MsGraph.Endpoint
	c.DomainsClient.BaseClient.RetryableClient.RetryMax = retry

	c.GroupsAppRoleAssignmentsClient = msgraph.NewGroupsAppRoleAssignmentsClient(c.Connections["default"].AuthConfig.TenantID)
	c.GroupsAppRoleAssignmentsClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.GroupsAppRoleAssignmentsClient.BaseClient.Endpoint = c.Connections["default"].AuthConfig.Environment.MsGraph.Endpoint
	c.GroupsAppRoleAssignmentsClient.BaseClient.RetryableClient.RetryMax = retry

	c.GroupsClient = msgraph.NewGroupsClient(c.Connections["default"].AuthConfig.TenantID)
	c.GroupsClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.GroupsClient.BaseClient.Endpoint = c.Connections["default"].AuthConfig.Environment.MsGraph.Endpoint
	c.GroupsClient.BaseClient.RetryableClient.RetryMax = retry

	c.IdentityProvidersClient = msgraph.NewIdentityProvidersClient(c.Connections["default"].AuthConfig.TenantID)
	c.IdentityProvidersClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.IdentityProvidersClient.BaseClient.Endpoint = c.Connections["default"].AuthConfig.Environment.MsGraph.Endpoint
	c.IdentityProvidersClient.BaseClient.RetryableClient.RetryMax = retry

	c.InvitationsClient = msgraph.NewInvitationsClient(c.Connections["default"].AuthConfig.TenantID)
	c.InvitationsClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.InvitationsClient.BaseClient.Endpoint = c.Connections["default"].AuthConfig.Environment.MsGraph.Endpoint
	c.InvitationsClient.BaseClient.RetryableClient.RetryMax = retry

	c.MeClient = msgraph.NewMeClient(c.Connections["default"].AuthConfig.TenantID)
	c.MeClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.MeClient.BaseClient.Endpoint = c.Connections["default"].AuthConfig.Environment.MsGraph.Endpoint
	c.MeClient.BaseClient.RetryableClient.RetryMax = retry

	c.NamedLocationsClient = msgraph.NewNamedLocationsClient(c.Connections["default"].AuthConfig.TenantID)
	c.NamedLocationsClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.NamedLocationsClient.BaseClient.Endpoint = c.Connections["default"].AuthConfig.Environment.MsGraph.Endpoint
	c.NamedLocationsClient.BaseClient.RetryableClient.RetryMax = retry

	c.ReportsClient = msgraph.NewReportsClient(c.Connections["default"].AuthConfig.TenantID)
	c.ReportsClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.ReportsClient.BaseClient.Endpoint = c.Connections["default"].AuthConfig.Environment.MsGraph.Endpoint
	c.ReportsClient.BaseClient.RetryableClient.RetryMax = retry

	c.RoleAssignmentsClient = msgraph.NewRoleAssignmentsClient(c.Connections["default"].AuthConfig.TenantID)
	c.RoleAssignmentsClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.RoleAssignmentsClient.BaseClient.Endpoint = c.Connections["default"].AuthConfig.Environment.MsGraph.Endpoint
	c.RoleAssignmentsClient.BaseClient.RetryableClient.RetryMax = retry

	c.RoleDefinitionsClient = msgraph.NewRoleDefinitionsClient(c.Connections["default"].AuthConfig.TenantID)
	c.RoleDefinitionsClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.RoleDefinitionsClient.BaseClient.Endpoint = c.Connections["default"].AuthConfig.Environment.MsGraph.Endpoint
	c.RoleDefinitionsClient.BaseClient.RetryableClient.RetryMax = retry

	c.SchemaExtensionsClient = msgraph.NewSchemaExtensionsClient(c.Connections["default"].AuthConfig.TenantID)
	c.SchemaExtensionsClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.SchemaExtensionsClient.BaseClient.Endpoint = c.Connections["default"].AuthConfig.Environment.MsGraph.Endpoint
	c.SchemaExtensionsClient.BaseClient.RetryableClient.RetryMax = retry

	c.ServicePrincipalsAppRoleAssignmentsClient = msgraph.NewServicePrincipalsAppRoleAssignmentsClient(c.Connections["default"].AuthConfig.TenantID)
	c.ServicePrincipalsAppRoleAssignmentsClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.ServicePrincipalsAppRoleAssignmentsClient.BaseClient.Endpoint = c.Connections["default"].AuthConfig.Environment.MsGraph.Endpoint
	c.ServicePrincipalsAppRoleAssignmentsClient.BaseClient.RetryableClient.RetryMax = retry

	c.ServicePrincipalsClient = msgraph.NewServicePrincipalsClient(c.Connections["default"].AuthConfig.TenantID)
	c.ServicePrincipalsClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.ServicePrincipalsClient.BaseClient.Endpoint = c.Connections["default"].AuthConfig.Environment.MsGraph.Endpoint
	c.ServicePrincipalsClient.BaseClient.RetryableClient.RetryMax = retry
	c.ServicePrincipalsClient.BaseClient.ApiVersion = msgraph.Version10

	c.SignInReportsClient = msgraph.NewSignInReportsClient(c.Connections["default"].AuthConfig.TenantID)
	c.SignInReportsClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.SignInReportsClient.BaseClient.Endpoint = c.Connections["default"].AuthConfig.Environment.MsGraph.Endpoint
	c.SignInReportsClient.BaseClient.RetryableClient.RetryMax = retry

	c.SynchronizationJobClient = msgraph.NewSynchronizationJobClient(c.Connections["default"].AuthConfig.TenantID)
	c.SynchronizationJobClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.SynchronizationJobClient.BaseClient.Endpoint = c.Connections["default"].AuthConfig.Environment.MsGraph.Endpoint
	c.SynchronizationJobClient.BaseClient.RetryableClient.RetryMax = retry

	c.UserFlowAttributesClient = msgraph.NewUserFlowAttributesClient(c.Connections["b2c"].AuthConfig.TenantID)
	c.UserFlowAttributesClient.BaseClient.Authorizer = c.Connections["b2c"].Authorizer
	c.UserFlowAttributesClient.BaseClient.Endpoint = c.Connections["b2c"].AuthConfig.Environment.MsGraph.Endpoint
	c.UserFlowAttributesClient.BaseClient.RetryableClient.RetryMax = retry

	c.UsersAppRoleAssignmentsClient = msgraph.NewUsersAppRoleAssignmentsClient(c.Connections["default"].AuthConfig.TenantID)
	c.UsersAppRoleAssignmentsClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.UsersAppRoleAssignmentsClient.BaseClient.Endpoint = c.Connections["default"].AuthConfig.Environment.MsGraph.Endpoint
	c.UsersAppRoleAssignmentsClient.BaseClient.RetryableClient.RetryMax = retry

	c.UsersClient = msgraph.NewUsersClient(c.Connections["default"].AuthConfig.TenantID)
	c.UsersClient.BaseClient.Authorizer = c.Connections["default"].Authorizer
	c.UsersClient.BaseClient.Endpoint = c.Connections["default"].AuthConfig.Environment.MsGraph.Endpoint
	c.UsersClient.BaseClient.RetryableClient.RetryMax = retry

	return
}
