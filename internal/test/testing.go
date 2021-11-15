package test

import (
	"context"
	"log"
	"os"
	"testing"

	"golang.org/x/oauth2"

	"github.com/manicminer/hamilton/msgraph"

	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/environments"
	"github.com/manicminer/hamilton/internal/utils"
)

var (
	tenantId              = os.Getenv("TENANT_ID")
	tenantDomain          = os.Getenv("TENANT_DOMAIN")
	clientId              = os.Getenv("CLIENT_ID")
	clientCertificate     = os.Getenv("CLIENT_CERTIFICATE")
	clientCertificatePath = os.Getenv("CLIENT_CERTIFICATE_PATH")
	clientCertPassword    = os.Getenv("CLIENT_CERTIFICATE_PASSWORD")
	clientSecret          = os.Getenv("CLIENT_SECRET")
	environment           = os.Getenv("AZURE_ENVIRONMENT")
)

type Connection struct {
	AuthConfig *auth.Config
	Authorizer auth.Authorizer
	DomainName string
}

// NewConnection configures and returns a Connection for use in tests.
func NewConnection(tokenVersion auth.TokenVersion) *Connection {
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
			EnableClientCertAuth:   true,
			EnableClientSecretAuth: true,
			EnableAzureCliToken:    true,
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
	Connection   *Connection
	RandomString string

	Claims auth.Claims
	Token  *oauth2.Token

	AccessPackageAssignmentPolicyClient       *msgraph.AccessPackageAssignmentPolicyClient
	AccessPackageCatalogClient                *msgraph.AccessPackageCatalogClient
	AccessPackageClient                       *msgraph.AccessPackageClient
	AccessPackageResourceClient               *msgraph.AccessPackageResourceClient
	AccessPackageResourceRequestClient        *msgraph.AccessPackageResourceRequestClient
	AccessPackageResourceRoleScopeClient      *msgraph.AccessPackageResourceRoleScopeClient
	ApplicationTemplatesClient                *msgraph.ApplicationTemplatesClient
	ApplicationsClient                        *msgraph.ApplicationsClient
	AppRoleAssignedToClient                   *msgraph.AppRoleAssignedToClient
	AuthenticationMethodsClient               *msgraph.AuthenticationMethodsClient
	ConditionalAccessPoliciesClient           *msgraph.ConditionalAccessPoliciesClient
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
	SchemaExtensionsClient                    *msgraph.SchemaExtensionsClient
	ServicePrincipalsAppRoleAssignmentsClient *msgraph.AppRoleAssignmentsClient
	ServicePrincipalsClient                   *msgraph.ServicePrincipalsClient
	SignInReportsClient                       *msgraph.SignInReportsClient
	UsersAppRoleAssignmentsClient             *msgraph.AppRoleAssignmentsClient
	UsersClient                               *msgraph.UsersClient
}

func NewTest(t *testing.T) (c *Test) {
	ctx := context.Background()
	var cancel context.CancelFunc

	if deadline, ok := t.Deadline(); ok {
		ctx, cancel = context.WithDeadline(context.Background(), deadline)
	}

	c = &Test{
		Context:      ctx,
		CancelFunc:   cancel,
		RandomString: RandomString(),
	}

	c.Connection = NewConnection(auth.TokenVersion2)
	c.Connection.Authorize(ctx, c.Connection.AuthConfig.Environment.MsGraph)

	var err error
	c.Token, err = c.Connection.Authorizer.Token()
	if err != nil {
		t.Fatalf("could not acquire access token: %v", err)
	}

	c.Claims, err = auth.ParseClaims(c.Token)
	if err != nil {
		t.Fatalf("could not parse claims: %v", err)
	}

	c.AccessPackageAssignmentPolicyClient = msgraph.NewAccessPackageAssignmentPolicyClient(c.Connection.AuthConfig.TenantID)
	c.AccessPackageAssignmentPolicyClient.BaseClient.Authorizer = c.Connection.Authorizer
	c.AccessPackageAssignmentPolicyClient.BaseClient.Endpoint = c.Connection.AuthConfig.Environment.MsGraph.Endpoint

	c.AccessPackageCatalogClient = msgraph.NewAccessPackageCatalogClient(c.Connection.AuthConfig.TenantID)
	c.AccessPackageCatalogClient.BaseClient.Authorizer = c.Connection.Authorizer
	c.AccessPackageCatalogClient.BaseClient.Endpoint = c.Connection.AuthConfig.Environment.MsGraph.Endpoint

	c.AccessPackageClient = msgraph.NewAccessPackageClient(c.Connection.AuthConfig.TenantID)
	c.AccessPackageClient.BaseClient.Authorizer = c.Connection.Authorizer
	c.AccessPackageClient.BaseClient.Endpoint = c.Connection.AuthConfig.Environment.MsGraph.Endpoint

	c.AccessPackageResourceClient = msgraph.NewAccessPackageResourceClient(c.Connection.AuthConfig.TenantID)
	c.AccessPackageResourceClient.BaseClient.Authorizer = c.Connection.Authorizer
	c.AccessPackageResourceClient.BaseClient.Endpoint = c.Connection.AuthConfig.Environment.MsGraph.Endpoint

	c.AccessPackageResourceRequestClient = msgraph.NewAccessPackageResourceRequestClient(c.Connection.AuthConfig.TenantID)
	c.AccessPackageResourceRequestClient.BaseClient.Authorizer = c.Connection.Authorizer
	c.AccessPackageResourceRequestClient.BaseClient.Endpoint = c.Connection.AuthConfig.Environment.MsGraph.Endpoint

	c.AccessPackageResourceRoleScopeClient = msgraph.NewAccessPackageResourceRoleScopeClient(c.Connection.AuthConfig.TenantID)
	c.AccessPackageResourceRoleScopeClient.BaseClient.Authorizer = c.Connection.Authorizer
	c.AccessPackageResourceRoleScopeClient.BaseClient.Endpoint = c.Connection.AuthConfig.Environment.MsGraph.Endpoint

	c.ApplicationTemplatesClient = msgraph.NewApplicationTemplatesClient(c.Connection.AuthConfig.TenantID)
	c.ApplicationTemplatesClient.BaseClient.Authorizer = c.Connection.Authorizer
	c.ApplicationTemplatesClient.BaseClient.Endpoint = c.Connection.AuthConfig.Environment.MsGraph.Endpoint

	c.ApplicationsClient = msgraph.NewApplicationsClient(c.Connection.AuthConfig.TenantID)
	c.ApplicationsClient.BaseClient.Authorizer = c.Connection.Authorizer
	c.ApplicationsClient.BaseClient.Endpoint = c.Connection.AuthConfig.Environment.MsGraph.Endpoint

	c.AppRoleAssignedToClient = msgraph.NewAppRoleAssignedToClient(c.Connection.AuthConfig.TenantID)
	c.AppRoleAssignedToClient.BaseClient.Authorizer = c.Connection.Authorizer
	c.AppRoleAssignedToClient.BaseClient.Endpoint = c.Connection.AuthConfig.Environment.MsGraph.Endpoint

	c.AuthenticationMethodsClient = msgraph.NewAuthenticationMethodsClient(c.Connection.AuthConfig.TenantID)
	c.AuthenticationMethodsClient.BaseClient.Authorizer = c.Connection.Authorizer
	c.AuthenticationMethodsClient.BaseClient.Endpoint = c.Connection.AuthConfig.Environment.MsGraph.Endpoint

	c.ConditionalAccessPoliciesClient = msgraph.NewConditionalAccessPoliciesClient(c.Connection.AuthConfig.TenantID)
	c.ConditionalAccessPoliciesClient.BaseClient.Authorizer = c.Connection.Authorizer
	c.ConditionalAccessPoliciesClient.BaseClient.Endpoint = c.Connection.AuthConfig.Environment.MsGraph.Endpoint

	c.DirectoryAuditReportsClient = msgraph.NewDirectoryAuditReportsClient(c.Connection.AuthConfig.TenantID)
	c.DirectoryAuditReportsClient.BaseClient.Authorizer = c.Connection.Authorizer
	c.DirectoryAuditReportsClient.BaseClient.Endpoint = c.Connection.AuthConfig.Environment.MsGraph.Endpoint

	c.DirectoryObjectsClient = msgraph.NewDirectoryObjectsClient(c.Connection.AuthConfig.TenantID)
	c.DirectoryObjectsClient.BaseClient.Authorizer = c.Connection.Authorizer
	c.DirectoryObjectsClient.BaseClient.Endpoint = c.Connection.AuthConfig.Environment.MsGraph.Endpoint

	c.DirectoryRoleTemplatesClient = msgraph.NewDirectoryRoleTemplatesClient(c.Connection.AuthConfig.TenantID)
	c.DirectoryRoleTemplatesClient.BaseClient.Authorizer = c.Connection.Authorizer
	c.DirectoryRoleTemplatesClient.BaseClient.Endpoint = c.Connection.AuthConfig.Environment.MsGraph.Endpoint

	c.DirectoryRolesClient = msgraph.NewDirectoryRolesClient(c.Connection.AuthConfig.TenantID)
	c.DirectoryRolesClient.BaseClient.Authorizer = c.Connection.Authorizer
	c.DirectoryRolesClient.BaseClient.Endpoint = c.Connection.AuthConfig.Environment.MsGraph.Endpoint

	c.DomainsClient = msgraph.NewDomainsClient(c.Connection.AuthConfig.TenantID)
	c.DomainsClient.BaseClient.Authorizer = c.Connection.Authorizer
	c.DomainsClient.BaseClient.Endpoint = c.Connection.AuthConfig.Environment.MsGraph.Endpoint

	c.GroupsAppRoleAssignmentsClient = msgraph.NewGroupsAppRoleAssignmentsClient(c.Connection.AuthConfig.TenantID)
	c.GroupsAppRoleAssignmentsClient.BaseClient.Authorizer = c.Connection.Authorizer
	c.GroupsAppRoleAssignmentsClient.BaseClient.Endpoint = c.Connection.AuthConfig.Environment.MsGraph.Endpoint

	c.GroupsClient = msgraph.NewGroupsClient(c.Connection.AuthConfig.TenantID)
	c.GroupsClient.BaseClient.Authorizer = c.Connection.Authorizer
	c.GroupsClient.BaseClient.Endpoint = c.Connection.AuthConfig.Environment.MsGraph.Endpoint

	c.IdentityProvidersClient = msgraph.NewIdentityProvidersClient(c.Connection.AuthConfig.TenantID)
	c.IdentityProvidersClient.BaseClient.Authorizer = c.Connection.Authorizer
	c.IdentityProvidersClient.BaseClient.Endpoint = c.Connection.AuthConfig.Environment.MsGraph.Endpoint

	c.InvitationsClient = msgraph.NewInvitationsClient(c.Connection.AuthConfig.TenantID)
	c.InvitationsClient.BaseClient.Authorizer = c.Connection.Authorizer
	c.InvitationsClient.BaseClient.Endpoint = c.Connection.AuthConfig.Environment.MsGraph.Endpoint

	c.MeClient = msgraph.NewMeClient(c.Connection.AuthConfig.TenantID)
	c.MeClient.BaseClient.Authorizer = c.Connection.Authorizer
	c.MeClient.BaseClient.Endpoint = c.Connection.AuthConfig.Environment.MsGraph.Endpoint

	c.NamedLocationsClient = msgraph.NewNamedLocationsClient(c.Connection.AuthConfig.TenantID)
	c.NamedLocationsClient.BaseClient.Authorizer = c.Connection.Authorizer
	c.NamedLocationsClient.BaseClient.Endpoint = c.Connection.AuthConfig.Environment.MsGraph.Endpoint

	c.ReportsClient = msgraph.NewReportsClient(c.Connection.AuthConfig.TenantID)
	c.ReportsClient.BaseClient.Authorizer = c.Connection.Authorizer
	c.ReportsClient.BaseClient.Endpoint = c.Connection.AuthConfig.Environment.MsGraph.Endpoint

	c.SchemaExtensionsClient = msgraph.NewSchemaExtensionsClient(c.Connection.AuthConfig.TenantID)
	c.SchemaExtensionsClient.BaseClient.Authorizer = c.Connection.Authorizer
	c.SchemaExtensionsClient.BaseClient.Endpoint = c.Connection.AuthConfig.Environment.MsGraph.Endpoint

	c.ServicePrincipalsAppRoleAssignmentsClient = msgraph.NewServicePrincipalsAppRoleAssignmentsClient(c.Connection.AuthConfig.TenantID)
	c.ServicePrincipalsAppRoleAssignmentsClient.BaseClient.Authorizer = c.Connection.Authorizer
	c.ServicePrincipalsAppRoleAssignmentsClient.BaseClient.Endpoint = c.Connection.AuthConfig.Environment.MsGraph.Endpoint

	c.ServicePrincipalsClient = msgraph.NewServicePrincipalsClient(c.Connection.AuthConfig.TenantID)
	c.ServicePrincipalsClient.BaseClient.Authorizer = c.Connection.Authorizer
	c.ServicePrincipalsClient.BaseClient.Endpoint = c.Connection.AuthConfig.Environment.MsGraph.Endpoint

	c.SignInReportsClient = msgraph.NewSignInReportsClient(c.Connection.AuthConfig.TenantID)
	c.SignInReportsClient.BaseClient.Authorizer = c.Connection.Authorizer
	c.SignInReportsClient.BaseClient.Endpoint = c.Connection.AuthConfig.Environment.MsGraph.Endpoint

	c.UsersAppRoleAssignmentsClient = msgraph.NewUsersAppRoleAssignmentsClient(c.Connection.AuthConfig.TenantID)
	c.UsersAppRoleAssignmentsClient.BaseClient.Authorizer = c.Connection.Authorizer
	c.UsersAppRoleAssignmentsClient.BaseClient.Endpoint = c.Connection.AuthConfig.Environment.MsGraph.Endpoint

	c.UsersClient = msgraph.NewUsersClient(c.Connection.AuthConfig.TenantID)
	c.UsersClient.BaseClient.Authorizer = c.Connection.Authorizer
	c.UsersClient.BaseClient.Endpoint = c.Connection.AuthConfig.Environment.MsGraph.Endpoint

	return
}
