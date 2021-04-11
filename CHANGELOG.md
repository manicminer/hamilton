## 0.10.0 (April 10, 2021)

⚠️ BREAKING CHANGES:

- This release refactors various packages to make for a better import experience.
- `base`, `clients` and `models` packages have been combined into a single `msgraph` package.
- `base/aadgraph` package has been moved to `aadgraph`.
- `base/odata` package has been moved to `odata`.

## 0.9.0 (March 1, 2021)

- Add support for guest user invitations ([#21](https://github.com/manicminer/hamilton/pull/21))

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

