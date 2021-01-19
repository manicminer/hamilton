## 0.4.0 (January 19, 2021)

- Adds the `UsersCLient.ListGroupMemberships()` method.
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

