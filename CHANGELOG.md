## 0.29.0 (Unreleased)

## 0.28.2 (September 10, 2021)

- Bug fix: Correctly handle HTTP responses after retries have been exhausted for a request, so that the correct status and error can be returned ([#100](https://github.com/manicminer/hamilton/pull/100))

## 0.28.1 (September 9, 2021)

- Bug fix: Try to detect when running in Azure Cloud Shell and avoid specifying the tenant ID for Azure CLI authentication ([#98](https://github.com/manicminer/hamilton/pull/98))
- Bug fix: Use the correct base64 decoder when parsing token claims ([#97](https://github.com/manicminer/hamilton/pull/97))

⚠️ BREAKING CHANGES:

- Bug fix: `User.PasswordPolicies` has changed from a `*string` to a `*StringNullWhenEmpty` ([#96](https://github.com/manicminer/hamilton/pull/96))

## 0.28.0 (September 7, 2021)

- Support for [application templates](https://docs.microsoft.com/en-us/graph/api/resources/applicationtemplate?view=graph-rest-1.0) ([#95](https://github.com/manicminer/hamilton/pull/95))

## 0.27.0 (September 2, 2022)

- Add some value types for `ConditionalAccessPolicyState` and `InvitedUserType` ([#94](https://github.com/manicminer/hamilton/pull/94))

## 0.26.0 (September 1, 2021)

- `auth.CachedAuthorizer` - export this type and its `Source` field so that consumers can inspect it ([#90](https://github.com/manicminer/hamilton/pull/90))
- Bugfix: set the struct tag for `ServicePrincipal.Owners` field so it is marshaled correctly  ([#91](https://github.com/manicminer/hamilton/pull/91))

⚠️ BREAKING CHANGES:

- The `auth.CachedAuthorizer()` function has been renamed to `auth.NewCachedAuthorizer()` ([#90](https://github.com/manicminer/hamilton/pull/90))

## 0.25.0 (August 24, 2021)

- Support for [authentication methods](https://docs.microsoft.com/en-us/graph/api/resources/authenticationmethods-overview?view=graph-rest-beta) ([#89](https://github.com/manicminer/hamilton/pull/89))

## 0.24.0 (August 17, 2021)

- When authenticating using Azure CLI, access tokens are now cached to avoid repeatedly invoking `az` to get the latest token ([#88](https://github.com/manicminer/hamilton/pull/88))
- Support for [authentication methods usage reports](https://docs.microsoft.com/en-us/graph/api/resources/authenticationmethods-usage-insights-overview?view=graph-rest-beta) ([#85](https://github.com/manicminer/hamilton/pull/85))
- Support for [generic directory objects](https://docs.microsoft.com/en-us/graph/api/resources/directoryobject?view=graph-rest-beta) ([#86](https://github.com/manicminer/hamilton/pull/86))
- Add the `MemberOf` field to the `User` struct ([#84](https://github.com/manicminer/hamilton/pull/84))

⚠️ BREAKING CHANGES:

- The `ID` field of the `Application`, `DirectoryRole`, `Group`, `ServicePrincipal` and `User` models has been removed and is now a field of the embedded `DirectoryObject` struct ([#86](https://github.com/manicminer/hamilton/pull/86))
- The `Members` and/or `Owners` fields of the `Application`, `DirectoryRole`, `Group` and `ServicePrincipal` models have changed from a `*[]string` to a `*Members` and `*Owners` respectively ([#86](https://github.com/manicminer/hamilton/pull/86))
  - The `Members` and `Owners` types are based on `[]DirectoryObject` and have methods to marshal/unmarshal the `ODataId` fields of the contained `DirectoryObject`s
- The `AppendMember()` and/or `AppendOwner()` methods of the `Application`, `Group` and `ServicePrincipal` models are no longer required and have been removed ([#86](https://github.com/manicminer/hamilton/pull/86))

## 0.23.1 (July 21, 2021)

- Disable the default logger for `retryablehttp.Client{}` ([#83](https://github.com/manicminer/hamilton/pull/83))

## 0.23.0 (July 21, 2021)

- Support for schema extension data for Groups and Users  ([#81](https://github.com/manicminer/hamilton/pull/81))
  - Marshaling of schema extension data is handled automatically by the Group and User structs, enabling use of the existing `Update()` methods on the respective clients.
  - Unmarshaling is handled by either the provided `msgraph.SchemaExtensionMap` type, or a custom type supplied by the caller. Such a custom type must have an explicit `UnmarshalJSON()` method to satisfy the `SchemaExtensionProperties` interface. Both approaches have examples in the `TestSchemaExtensionsClient()` test.
- Support for injecting and sequencing middleware functions for manipulating and/or copying requests and responses ([#78](https://github.com/manicminer/hamilton/pull/78))
  - See [example.go](https://github.com/manicminer/hamilton/blob/main/example/example.go) for an example that logs requests and responses
- Request retry handling for rate limiting, server errors and replication delays is now handled by [go-retryablehttp](https://github.com/hashicorp/go-retryablehttp) ([#78](https://github.com/manicminer/hamilton/pull/78))
- `msgraph.Client{}.HttpClient` is now exported so callers can supply their own `http.Client` ([#78](https://github.com/manicminer/hamilton/pull/78))

⚠️ BREAKING CHANGES:

- Support `odata.Query{}` in more client methods ([#80](https://github.com/manicminer/hamilton/pull/80))
  - `ApplicationsClient{}.Get()`
  - `ApplicationsClient{}.GetDeleted()`
  - `ApplicationsClient{}.ListExtensions()`
  - `ConditionalAccessPolicyClient{}.Get()`
  - `DirectoryAuditReportsClient{}.Get()`
  - `DomainsClient{}.List()`
  - `DomainsClient{}.Get()`
  - `GroupsClient{}.Get()`
  - `GroupsClient{}.GetDeleted()`
  - `MeClient{}.Get()`
  - `MeClient{}.GetProfile()`
  - `NamedLocationsClient{}.Get()`
  - `NamedLocationsClient{}.GetCountry()`
  - `NamedLocationsClient{}.GetIP()`
  - `SchemaExtensionsClient{}.Get()`
  - `ServicePrincipalsClient{}.Get()`
  - `ServicePrincipalsClient{}.ListAppRoleAssignments()`
  - `SignInReportsClient{}.Get()`
  - `UsersClient{}.Get()`
  - `UsersClient{}.GetDeleted()`

## 0.22.0 (July 13, 2021)

- `msgraph.ServicePrincipal{}` now supports the `Description` field ([#77](https://github.com/manicminer/hamilton/pull/77))
- `msgraph.ServicePrincipal{}` now supports the `Notes` field ([#77](https://github.com/manicminer/hamilton/pull/77))
- `msgraph.ServicePrincipal{}` now supports the `SamlMetadataUrl` field ([#77](https://github.com/manicminer/hamilton/pull/77))

⚠️ BREAKING CHANGES:

- `environments.ApiAppId` is now a type alias
- `msgraph.ServicePrincipal{}.LoginUrl` is now a `StringNullWhenEmpty` type ([#77](https://github.com/manicminer/hamilton/pull/77))
- `msgraph.ServicePrincipal{}.PreferredSingleSignOnMode` is now a type alias pointer (formerly a string pointer) ([#77](https://github.com/manicminer/hamilton/pull/77))

## 0.21.0 (July 6, 2021)

- `msgraph.User{}` now supports the `AgeGroup` field ([#76](https://github.com/manicminer/hamilton/pull/76))
- `msgraph.User{}` now supports the `ConsentProvidedForMinor` field ([#76](https://github.com/manicminer/hamilton/pull/76))

⚠️ BREAKING CHANGES:

- `msgraph.Application{}.SignInAudience` is now a pointer reference ([#76](https://github.com/manicminer/hamilton/pull/76))
- `msgraph.ServicePrincipal{}.SignInAudience` is now a pointer reference ([#76](https://github.com/manicminer/hamilton/pull/76))
- `msgraph.Group{}.ResourceBehaviorOptions` is now a custom type ([#75](https://github.com/manicminer/hamilton/pull/75))
- `msgraph.Group{}.ResourceProvisioningOptions` is now a custom type ([#75](https://github.com/manicminer/hamilton/pull/75))
- `msgraph.Group{}.Theme` is now a custom type ([#75](https://github.com/manicminer/hamilton/pull/75))
- `msgraph.Group{}.Visibility` is now a custom type ([#75](https://github.com/manicminer/hamilton/pull/75))
- `msgraph.User{}.EmployeeId` is now a `StringNullWhenEmpty` type ([#76](https://github.com/manicminer/hamilton/pull/76))
- `msgraph.User{}.FaxNumber` is now a `StringNullWhenEmpty` type ([#76](https://github.com/manicminer/hamilton/pull/76))
- `msgraph.User{}.Mail` is now a `StringNullWhenEmpty` type ([#76](https://github.com/manicminer/hamilton/pull/76))
- `msgraph.User{}.PreferredLanguage` is now a `StringNullWhenEmpty` type ([#76](https://github.com/manicminer/hamilton/pull/76))
- `msgraph.ApplicationExtensionTargetObject` is now a type alias ([#76](https://github.com/manicminer/hamilton/pull/76))
- `msgraph.AppRoleAllowedMemberType` is now a type alias ([#76](https://github.com/manicminer/hamilton/pull/76))
- `msgraph.BodyType` is now a type alias ([#76](https://github.com/manicminer/hamilton/pull/76))
- `msgraph.ExtensionSchemaPropertyDataType` is now a type alias ([#76](https://github.com/manicminer/hamilton/pull/76))
- `msgraph.GroupType` is now a type alias ([#76](https://github.com/manicminer/hamilton/pull/76))
- `msgraph.GroupMembershipClaim` is now a type alias ([#76](https://github.com/manicminer/hamilton/pull/76))
- `msgraph.KeyCredentialType` is now a type alias ([#76](https://github.com/manicminer/hamilton/pull/76))
- `msgraph.KeyCredentialUsage` is now a type alias ([#76](https://github.com/manicminer/hamilton/pull/76))
- `msgraph.PermissionScopeType` is now a type alias ([#76](https://github.com/manicminer/hamilton/pull/76))
- `msgraph.ResourceAccessType` is now a type alias ([#76](https://github.com/manicminer/hamilton/pull/76))
- `msgraph.SignInAudience` is now a type alias ([#76](https://github.com/manicminer/hamilton/pull/76))

## 0.20.0 (July 1, 2021)

- Support the `spa` field for applications ([#74](https://github.com/manicminer/hamilton/pull/74))

## 0.19.0 (June 29, 2021)

- Support for [schema extensions](https://docs.microsoft.com/en-us/graph/api/resources/schemaextension?view=graph-rest-beta) ([#68](https://github.com/manicminer/hamilton/pull/68))
- Support for retrieving `SignInActivity` for users ([#72](https://github.com/manicminer/hamilton/pull/72))

⚠️ BREAKING CHANGES:

- Support for passing the raw bytes of a PKCS#12 bundle when using client certificate authentication. This alters the method signature of `auth.NewClientCertificateAuthorizer()` but does not affect the use of a PFX file read from the filesystem. See [#65](https://github.com/manicminer/hamilton/pull/65) for details and example usage.

## 0.18.0 (June 22, 2021)

- Support for [application extensions](https://docs.microsoft.com/en-us/graph/api/resources/extensionproperty?view=graph-rest-beta) ([#61](https://github.com/manicminer/hamilton/pull/61))
- Support for [directory audit and sign-in reports](https://docs.microsoft.com/en-us/graph/api/resources/azure-ad-auditlog-overview?view=graph-rest-beta) ([#61](https://github.com/manicminer/hamilton/pull/61))

⚠️ BREAKING CHANGES:

- This release introduces support for [OData query parameters](https://docs.microsoft.com/en-us/graph/query-parameters) via a new type `odata.Query{}`. Instead of accepting just a filter string, all clients now accept an instance of `odata.Query{}` on relevant List methods which encapsulates any combination of odata queries such as `$filter`, `$search`, `$top` etc. All documented parameters are supported and wrapped lightly where appropriate.  ([#63](https://github.com/manicminer/hamilton/pull/63))
- Updating to this release will require changes to affected method calls, for example:
  ```go
  apps, status, err := appsClient.List(ctx, odata.Query{
  	Filter: fmt.Sprintf("startsWith(displayName,'%s')", searchTerm),
  	OrderBy: odata.OrderBy{
  		Field:     "displayName",
  		Direction: "asc",
  	},
  	Top: 10,
  })
  ```
- Where an empty filter string was previously specified, it should be replaced with an empty `odata.Query{}` struct:
  ```go
  apps, status, err := appsClient.List(ctx, odata.Query{})
  ```


## 0.17.0 (June 15, 2021)

- Support for [restoring deleted applications/users/groups](https://docs.microsoft.com/en-us/graph/api/directory-deleteditems-restore?view=graph-rest-1.0) ([#58](https://github.com/manicminer/hamilton/pull/58))
- Support `PersonalMicrosoftAccount` for the `SignInAudience` field for Applications ([#59](https://github.com/manicminer/hamilton/pull/59))

⚠️ BREAKING CHANGES:

- This release adds a new type alias `StringNullWhenEmpty` which has replaced several existing field string types
- It enables zeroing field values that don't accept empty strings. See ([#59](https://github.com/manicminer/hamilton/pull/59)) for details and example usage

## 0.16.0 (June 08, 2021)

BEHAVIORAL CHANGES:

- This release implements a retry mechanism for some types of failed requests where the likely cause is indicated to be replication delays in Azure Active Directory ([#57](https://github.com/manicminer/hamilton/pull/57))
- Client methods which retrieve, update or delete _single_, _mutable_ objects will all exert this retry mechanism, and may take up to 2 minutes to return (successfully or not)
- To opt out of this behavior, simply set the `BaseClient.DisableRetries` field to `true` on your client(s), for example:
  ```go
  client := msgraph.NewApplicationsClient(tenantId)
  client.BaseClient.DisableRetries = true
  ```

## 0.15.0 (June 01, 2021)

- Bug fix: Set correct OData types when updating named locations ([#55](https://github.com/manicminer/hamilton/pull/55))
- Support for [permanently deleting](https://docs.microsoft.com/en-us/graph/api/directory-deleteditems-delete?view=graph-rest-1.0&tabs=http) applications, groups and service principals ([#54](https://github.com/manicminer/hamilton/pull/54))
- Add a `NamedLocationsClient{}.Get()` method ([#56](https://github.com/manicminer/hamilton/pull/56))

## 0.14.1 (May 28, 2021)

- Bug fix: Restore a missing field `OnPremisesImmutableId` in the User model ([#53](https://github.com/manicminer/hamilton/pull/53))

## 0.14.0 (May 27, 2021)

- Bug fix: Correctly marshal the request body for `ApplicationsClient{}.AddPassword()` and `ServicePrincipalsClient{}.AddPassword()` ([#49](https://github.com/manicminer/hamilton/pull/49))
- Bug fix: Resolve a potential race condition where a cached access token might be refreshed multiple times unnecessarily ([#46](https://github.com/manicminer/hamilton/pull/46))
- Support for [app role assignments](https://docs.microsoft.com/en-us/graph/api/resources/approleassignment?view=graph-rest-1.0) using the [appRolesAssignedTo](https://docs.microsoft.com/en-us/graph/api/serviceprincipal-list-approleassignedto?view=graph-rest-1.0&tabs=http) and [appRoleAssignments](https://docs.microsoft.com/en-us/graph/api/user-list-approleassignments?view=graph-rest-1.0&tabs=http) endpoints ([#39](https://github.com/manicminer/hamilton/pull/39))
- Support for listing [deleted applications, groups and users](https://docs.microsoft.com/en-us/graph/api/directory-deleteditems-list?view=graph-rest-beta) ([#48](https://github.com/manicminer/hamilton/pull/48))
- Support for retrieving [deleted applications, groups and users](https://docs.microsoft.com/en-us/graph/api/directory-deleteditems-get?view=graph-rest-beta) ([#51](https://github.com/manicminer/hamilton/pull/51))

## 0.13.0 (May 18, 2021)

- Bug fix: Don't clear `GroupMembershipClaims` when nil for an Application ([#40](https://github.com/manicminer/hamilton/pull/40))
- Bug fix: Handle empty OData error collections ([#43](https://github.com/manicminer/hamilton/pull/43))
- Support for sending emails from the authenticated user principal or a specified user ([#37](https://github.com/manicminer/hamilton/pull/37))
- Support for the [ownedObjects endpoint](https://docs.microsoft.com/en-us/graph/api/serviceprincipal-list-ownedobjects?view=graph-rest-beta&tabs=http) for service principals ([#38](https://github.com/manicminer/hamilton/pull/38))
- Support for managing [identity providers](https://docs.microsoft.com/en-us/graph/api/resources/identityproviderbase?view=graph-rest-beta) ([#41](https://github.com/manicminer/hamilton/pull/41))
- Support [adding](https://docs.microsoft.com/en-us/graph/api/application-addpassword?view=graph-rest-beta&tabs=http) and [removing](https://docs.microsoft.com/en-us/graph/api/application-removepassword?view=graph-rest-beta&tabs=http) application passwords ([#44](https://github.com/manicminer/hamilton/pull/44))
- Support [adding](https://docs.microsoft.com/en-us/graph/api/serviceprincipal-addpassword?view=graph-rest-beta&tabs=http) and [removing](https://docs.microsoft.com/en-us/graph/api/serviceprincipal-removepassword?view=graph-rest-beta&tabs=http) service principal passwords ([#45](https://github.com/manicminer/hamilton/pull/45))

## 0.12.0 (April 23, 2021)

- Support for [managing Directory Roles](https://docs.microsoft.com/en-us/graph/api/resources/directoryrole?view=graph-rest-beta) ([#30](https://github.com/manicminer/hamilton/pull/30))
- Support for [activating Directory Roles](https://docs.microsoft.com/en-us/graph/api/directoryrole-post-directoryroles?view=graph-rest-beta&tabs=http) ([#31](https://github.com/manicminer/hamilton/pull/31))
- Support for [App Role Assignments](https://docs.microsoft.com/en-us/graph/api/group-post-approleassignments?view=graph-rest-1.0&tabs=http) ([#32](https://github.com/manicminer/hamilton/pull/32))
- Restore the retry mechanism previously introduced in v0.8.0
- Use the `odata` package for parsing common error messages
- Handle some additional errors, mainly for `ioutil.Read*()`
- Add more `ValidStatusFunc`s for gracefully handling existing owner and member refs
- Remove an unused struct field `auth.ClientCredentialsConfig{}.Expires`

⚠️ BREAKING CHANGES:

- `msgraph.Application{}.GroupMembershipClaims` is now a custom type
- `msgraph.Application{}.SignInAudience` is now a custom type
- `msgraph.AppRole{}.AllowedMemberTypes` is now a custom type
- `msgraph.KeyCredential{}.Usage` is now a custom type
- `msgraph.PermissionScope{}.Type` is now a custom type
- `msgraph.ResourceAccess{}.Type` is now a custom type
- `msgraph.ServicePrincipal{}.SignInAudience` is now a custom type

## 0.11.0 (April 13, 2021)

- Support for [Conditional Access Policies](https://docs.microsoft.com/en-us/graph/api/resources/conditionalaccesspolicy?view=graph-rest-beta) ([#23](https://github.com/manicminer/hamilton/pull/23))
- Support for [Named Locations](https://docs.microsoft.com/en-us/graph/api/resources/namedlocation?view=graph-rest-beta) (IP-based and Country-based) ([#24](https://github.com/manicminer/hamilton/pull/24))
- Support for [Directory Role Templates](https://docs.microsoft.com/en-us/graph/api/resources/directoryroletemplate?view=graph-rest-beta) ([#27](https://github.com/manicminer/hamilton/pull/27))
- Set a default User Agent string if not provided by the caller
- Improved error handling

## 0.10.0 (April 10, 2021)

⚠️ BREAKING CHANGES:

- This release refactors various packages to make for a better import experience.
- `base`, `clients` and `models` packages have been combined into a single `msgraph` package.
- `base/aadgraph` package has been moved to `aadgraph`.
- `base/odata` package has been moved to `odata`.

## 0.9.0 (March 1, 2021)

- Add support for [guest user invitations](https://docs.microsoft.com/en-us/graph/api/invitation-post?view=graph-rest-beta&tabs=http) ([#21](https://github.com/manicminer/hamilton/pull/21))

## 0.8.0 (February 2, 2021)

- Exponential backoff for handling rate limited and failed requests to MS Graph and AAD Graph

## 0.7.0 (January 27, 2021)

- Check for supported `az` command version when using Azure CLI authentication
- Remove dependency on deprecated package golang.org/x/oauth2/jws
- Merge the `auth/internal/microsoft` package into `auth` now that it's stable
- Validate the MSI auth configuration before returning an MsiAuthorizer - ensure the metadata endpoint is reachable

## 0.6.0 (January 26, 2021)

- Support authentication using VM managed identity.
- Add App ID for Teams Services API.

## 0.5.0 (January 24, 2021)

- All responses from Microsoft Graph and Azure Active Directory Graph are now parsed for OData metadata. Calls to `base.Client.Delete()`, `base.Client.Get()`, `base.Client.Patch()`, `base.Client.Post()` and `base.client.Put()` each now return OData metadata in addition to the complete response.
- Support for v1 and v2 access tokens from Microsoft Identity Platform. Defaults to v2 tokens.
- Support for acquiring access tokens for Microsoft Graph or Azure Active Directory graph. Since the MSID platform only supports scopes from a single API per token, these must be requested separately if using both APIs.
- Token claims parsed now includes scopes (`scp` claim)
- Export app IDs for several published APIs from Microsoft. These can be reliably consumed as `environments.PublishedApis`.
- Support for querying Azure Active Directory Graph API
    - This is intended as a stopgap solution for when it's not possible to perform an action using Microsoft Graph.
    - A number of endpoints do not yet have equivalents in MS Graph, notably those used by the Azure Portal.
    - There is only a base client at present.

⚠️ BREAKING CHANGES:

- Method signature for `auth.Config.NewAuthorizer()` has changed to include the API to request tokens for.
- Corresponding function signatures for `auth.NewAzureCliAuthorizer()`, `auth.NewClientCertificateAuthorizer()` and `auth.NewClientSecretAuthorizer()` also now include an `api` argument.
- The `auth.NewAzureCliConfig()` function also now includes an `api` argument.
- Functions implementing `base.ValidStatusFunc` must now accept a second argument as the pointer to a `base.odata.OData` struct.
- The `environments.MsGraphEndpoint` type has been removed in favor of `environments.ApiEndpoint`.
- The `endpoint` argument for `models.Application.AppendOwner()`, `models.Group.AppendMember()` and `models.Group.AppendOwner()` methods should now be an `environments.ApiEndpoint`.
- The environments package now exports `Api` structs for each national cloud and API combination, e.g. `environments.MsGraphGermany`.
- The `Environment` structs exports in the environments package have been changed to reference `Api`s and no longer include `MsGraphEndpoint`.

## 0.4.0 (January 19, 2021)

- Adds the `ServicePrincipalsClient.ListGroupMemberships()` method.
- Adds the `UsersClient.ListGroupMemberships()` method.
- Pagination handling: multiple pages of results with OData metadata are now automatically retrieved and merged together in the BaseClient for GET requests.

## 0.3.0 (January 18, 2021)

- Methods on `models.ApplcationApi` to manage `Oauth2PermissionScopes`.
- Tests for `auth` and `clients` packages.

## 0.2.0 (January 15, 2021)

Add support for all national clouds:

- Global: graph.microsoft.com
- Germany: graph.microsoft.de
- China: microsoftgraph.chinacloudapi.cn
- US Government L4: graph.microsoft.us
- US Government L5 (DOD): dod-graph.microsoft.us

Note that this is a breaking change from v0.1.0 as the signatures for all the clients have changed.
If you are using the global cloud, you do not need to specify this when creating a new client as it is the default.
However, you do need to specify a cloud environment when acquiring an access token using auth.NewAuthorizer.


## 0.1.0 (January 13, 2021)

Initial release. Working support for:

- Applications
- Domains
- Groups
- Service Principals
- Users

