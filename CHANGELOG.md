## v0.70.0 (June 6, 2024)

- Support for [updating a user photo](https://learn.microsoft.com/en-us/graph/api/profilephoto-update?view=graph-rest-1.0) ([#227](https://github.com/manicminer/hamilton/pull/227))
- Enable `RetryOn404ConsistencyFailureFunc` in the `Instantiate()` method for the `ApplicationTemplatesClient` ([#283](https://github.com/manicminer/hamilton/pull/283))
- Bugfix: the `AccessReviewSettings` field in the `AccessPackageAssignmentPolicy` model is nullable ([#284](https://github.com/manicminer/hamilton/pull/284))

## v0.69.0 (May 17, 2024)

- Support for [Attribute Sets](https://learn.microsoft.com/en-us/graph/api/resources/attributeset?view=graph-rest-1.0) and [Custom Security Attribute Definitions](https://learn.microsoft.com/en-us/graph/api/resources/customsecurityattributedefinition?view=graph-rest-1.0) ([#281](https://github.com/manicminer/hamilton/pull/281))

## v0.68.0 (May 16, 2024)

- Removed `omitempty` from struct tags for `ApplicationEnforcedRestrictions`, `CloudAppSecurity`, `PersistentBrowser`, and `SignInFrequency` fields in the `ConditionalAccessSessionControls` model ([#282](https://github.com/manicminer/hamilton/pull/282))

## v0.67.0 (March 28, 2024)

- Base Client: improve error visibility by returning the error for failed requests, when the response body is missing ([#280](https://github.com/manicminer/hamilton/pull/280))
- Support for [listing Access Package Resource Roles](https://learn.microsoft.com/en-us/graph/api/accesspackagecatalog-list-accesspackageresourceroles?view=graph-rest-beta) ([#278](https://github.com/manicminer/hamilton/pull/278))
- Support for PIM Role Management [Policies](https://learn.microsoft.com/en-us/graph/api/resources/unifiedrolemanagementpolicy?view=graph-rest-1.0), [Rules](https://learn.microsoft.com/en-us/graph/api/resources/unifiedrolemanagementpolicyrule?view=graph-rest-1.0), and [Assignments](https://learn.microsoft.com/en-us/graph/api/resources/unifiedrolemanagementpolicyassignment?view=graph-rest-1.0) ([#277](https://github.com/manicminer/hamilton/pull/277))
- Support for PIM Group Eligibility [Schedules](https://learn.microsoft.com/en-us/graph/api/resources/unifiedroleeligibilityschedule?view=graph-rest-1.0), [Instances](https://learn.microsoft.com/en-us/graph/api/resources/unifiedroleeligibilityscheduleinstance?view=graph-rest-1.0), and [Requests](https://learn.microsoft.com/en-us/graph/api/resources/unifiedroleeligibilityschedulerequest?view=graph-rest-1.0) ([#277](https://github.com/manicminer/hamilton/pull/277))
- Support for PIM Group Assignment [Schedules]([https://learn.microsoft.com/en-us/graph/api/resources/unifiedroleeligibilityschedule?view=graph-rest-1.0](https://learn.microsoft.com/en-us/graph/api/resources/privilegedaccessgroupassignmentschedule?view=graph-rest-1.0)), [Instances]([https://learn.microsoft.com/en-us/graph/api/resources/unifiedroleeligibilityscheduleinstance?view=graph-rest-1.0](https://learn.microsoft.com/en-us/graph/api/resources/privilegedaccessgroupassignmentscheduleinstance?view=graph-rest-1.0)), and [Requests]([https://learn.microsoft.com/en-us/graph/api/resources/unifiedroleeligibilityschedulerequest?view=graph-rest-1.0](https://learn.microsoft.com/en-us/graph/api/resources/privilegedaccessgroupassignmentschedulerequest?view=graph-rest-1.0)) ([#277](https://github.com/manicminer/hamilton/pull/277))
- Support for the `ApplicationFilter` field in the `ConditionalAccessApplications` model ([#268](https://github.com/manicminer/hamilton/pull/268))
- Added `SkipExchangeInstantOn` to supported `msgraph.GroupResourceBehaviorOption` values ([#275](https://github.com/manicminer/hamilton/pull/275))

⚠️ BREAKING CHANGES:

- `ExpirationPattern.Duration` has changed from a `*time.Duration` to a `*string` ([#276](https://github.com/manicminer/hamilton/pull/277))

## v0.66.0 (January 25, 2024)

- This is a maintenance release to update to the latest published module for `github.com/hashicorp/go-azure-sdk/sdk` ([#272](https://github.com/manicminer/hamilton/pull/272))

## 0.65.0 (October 26, 2023)

- Support for configuring [GuestsOrExternalUsers](https://learn.microsoft.com/en-us/graph/api/resources/conditionalaccessusers?view=graph-rest-1.0) for Conditional Access Policies ([#262](https://github.com/manicminer/hamilton/pull/262))
- Support for the `AuthenticationType` and `FrequencyInterval` fields in the `SignInFrequencySessionControl` model ([#263](https://github.com/manicminer/hamilton/pull/263))

## 0.64.0 (October 18, 2023)

- Add the `UpdateAllowedCombinations()` method to `AuthenticationStrengthPoliciesClient` ([#257](https://github.com/manicminer/hamilton/pull/257))
- Support for the `AppMetadata` field in the `ServicePrincipal` model ([#259](https://github.com/manicminer/hamilton/pull/259))
- Add the `SetFallbackPublicClient()` method to `ApplicationsClient` ([#260](https://github.com/manicminer/hamilton/pull/260))

⚠️ BREAKING CHANGES:

- `InformationalUrl.LogoUrl` has changed from a `*string` to a `*StringNullWhenEmpty` ([#260](https://github.com/manicminer/hamilton/pull/260))
- `InformationalUrl.MarketingUrl` has changed from a `*string` to a `*StringNullWhenEmpty` ([#260](https://github.com/manicminer/hamilton/pull/260))
- `InformationalUrl.PrivacyStatementUrl` has changed from a `*string` to a `*StringNullWhenEmpty` ([#260](https://github.com/manicminer/hamilton/pull/260))
- `InformationalUrl.SupportUrl` has changed from a `*string` to a `*StringNullWhenEmpty` ([#260](https://github.com/manicminer/hamilton/pull/260))
- `InformationalUrl.TermsOfServiceUrl` has changed from a `*string` to a `*StringNullWhenEmpty` ([#260](https://github.com/manicminer/hamilton/pull/260))

## 0.63.0 (July 27, 2023)

- Support for [Authentication Strength Policies](https://learn.microsoft.com/en-us/graph/api/resources/authenticationstrengthpolicy?view=graph-rest-1.0) ([#249](https://github.com/manicminer/hamilton/pull/249))
- Support `Manager` value for `AccessReviewReviewerType` ([#251](https://github.com/manicminer/hamilton/pull/251))
- Bugfix: Support for specifying `null` values for `grantControls` or `sessionControls` within [Conditional Access Policies](https://learn.microsoft.com/en-us/graph/api/resources/conditionalaccesspolicy?view=graph-rest-1.0) ([#250](https://github.com/manicminer/hamilton/pull/250))
- Bugfix: correct typo in names of `AccessReviewRecurrenceType` values, and in value of `AccessPackageRequestStateDelivered` ([#252](https://github.com/manicminer/hamilton/pull/252))

## 0.62.0 (July 14, 2023)

- Support for deleting an `accessPackageResourceRoleScope` ([#245](https://github.com/manicminer/hamilton/pull/245))
- Support for additional group behaviors `CalendarMemberReadOnly` and `ConnectorsDisabled` ([#247](https://github.com/manicminer/hamilton/pull/247))
- Support for the `ServicePrincipalRiskLevels` field in the `ConditionalAccessConditionSet` model ([#246](https://github.com/manicminer/hamilton/pull/246))
- Support for [Role Eligibility Schedule Requests](https://learn.microsoft.com/en-us/graph/api/resources/unifiedroleeligibilityschedulerequest?view=graph-rest-1.0) ([#204](https://github.com/manicminer/hamilton/pull/204))
- Bugfix: fake a 404 response when an `accessPackageResourceRoleScope` could not be found ([#245](https://github.com/manicminer/hamilton/pull/245))
- Bugfix: correctly unmarshal the `onPremisesPublishing` field in the `Application` model ([#244](https://github.com/manicminer/hamilton/pull/244))

## 0.61.0 (April 13, 2023)

- Support additional types for UserflowAttributeDataType ([#241](https://github.com/manicminer/hamilton/pull/241))
- Retry on 404 when creating directory role assignments ([#242](https://github.com/manicminer/hamilton/pull/242))

## 0.60.0 (April 3, 2023)

- dependencies: updating to `v0.20230331.1143618` of `github.com/hashicorp/go-azure-sdk` ([#228](https://github.com/manicminer/hamilton/pull/228))
- Add the `GetMembers()` method to the `GroupsClient` ([#236](https://github.com/manicminer/hamilton/pull/236))
- Support for the `ClientApplications` field in the `ConditionalAccessConditionSet` model ([#235](https://github.com/manicminer/hamilton/pull/235))
- Support for the `ServiceManagementReference` field in the `Application` model ([#233](https://github.com/manicminer/hamilton/pull/233))
- Support for the `LastNonInteractiveSignInDateTime` and `LastNonInteractiveSignInRequestId` fields in the `SignInActivity` model ([#237](https://github.com/manicminer/hamilton/pull/237))
- Support for managing [Token Issuance Policies](https://learn.microsoft.com/en-us/graph/api/resources/tokenissuancepolicy?view=graph-rest-1.0) for service principals ([#215](https://github.com/manicminer/hamilton/pull/215))
- Support for managing [Windows Autopilot Deployment Profiles](https://learn.microsoft.com/en-us/graph/api/resources/intune-enrollment-azureadwindowsautopilotdeploymentprofile?view=graph-rest-beta) ([#228](https://github.com/manicminer/hamilton/pull/228))

## 0.59.0 (March 1, 2023)

- Bugfix: Allow the `SynchronizationJobClient{}.ProvisionOnDemand()` method to recognise HTTP 200 responses ([#226](https://github.com/manicminer/hamilton/pull/226))

⚠️ BREAKING CHANGES:

- This release removes support for specifying the tenant ID as part of the request URL, as this causes some issues with newer APIs and is no longer supported by Microsoft Graph ([#230](https://github.com/manicminer/hamilton/pull/230))

## 0.58.0 (February 23, 2023)

⚠️ BREAKING CHANGES:

- This release removes the `auth`, `environments` and `odata` packages, replacing them with equivalent packages from the `github.com/hashicorp/go-azure-sdk` module.
- In order to use the `msgraph` clients, you will now need to make use of the newer authorizers from the `github.com/hashicorp/go-azure-sdk/sdk/auth` package. The [example](https://github.com/manicminer/hamilton/blob/main/example/example.go) in this repository have been updated accordingly.

## 0.57.1 (February 21, 2023)

- Bugfix: `Notes` in the `Application` model has changed from a `*string` to a `*StringNullWhenEmpty` ([#225](https://github.com/manicminer/hamilton/pull/225))

## 0.57.0 (February 21, 2023)

- Add the `ListAdministrativeUnitMemberships()` method to `GroupsClient` ([#220](https://github.com/manicminer/hamilton/pull/220))
- Support for `Notes` field in the `Application` model ([#218](https://github.com/manicminer/hamilton/pull/218))
- Bugfix: accommodate field mis-naming for `Oauth2RequiredPostResponse` in `Application` model (see [upstream bug report](https://github.com/microsoftgraph/msgraph-metadata/issues/273)) ([#221](https://github.com/manicminer/hamilton/pull/221))

## 0.56.0 (February 13, 2023)

- Auth: support for reading PKCS#12 bundles containing additional CA certificates ([#212](https://github.com/manicminer/hamilton/pull/212))
- Environments: add `MicrosoftOffice`, `MicrosoftTeams`, `MicrosoftTeamsWebClient`, `Office365SuiteUx`, `OfficeHome`, `OfficeUwpPwa` and `OssRdbmsPostgreSqlFlexibleServerAadAuthentication` to `PublishedApis` ([#216](https://github.com/manicminer/hamilton/pull/216))
- Support for [Access Package Assignment Requests](https://learn.microsoft.com/en-us/graph/api/resources/accesspackageassignmentrequest?view=graph-rest-1.0) ([#210](https://github.com/manicminer/hamilton/pull/210))
- Support for [Role Definitions](https://learn.microsoft.com/en-us/graph/api/resources/unifiedroledefinition?view=graph-rest-1.0) and [Role Assignments](https://learn.microsoft.com/en-us/graph/api/rbacapplication-post-roleassignments?view=graph-rest-1.0) for [Entitlement Management](https://learn.microsoft.com/en-us/graph/api/resources/entitlementmanagement-overview?view=graph-rest-1.0) ([#208](https://github.com/manicminer/hamilton/pull/208))
- Support for the [`DisableResilienceDefaults` session control](https://learn.microsoft.com/en-us/graph/api/resources/conditionalaccesssessioncontrols?view=graph-rest-1.0) for Conditional Access Policies ([#207](https://github.com/manicminer/hamilton/pull/207))
- Support for [Terms of Use Agreements](https://learn.microsoft.com/en-us/graph/api/resources/agreement?view=graph-rest-1.0) ([#209](https://github.com/manicminer/hamilton/pull/209))
- Bugfix: `PreferredTokenSigningKeyThumbprint` in the `ServicePrincipal` model has changed from a `*string` to a `*StringNullWhenEmpty` ([#211](https://github.com/manicminer/hamilton/pull/211))

## 0.55.0 (January 27, 2023)

- Support for [creating groups in an administrative unit](https://learn.microsoft.com/en-us/azure/active-directory/roles/admin-units-members-add#create-a-new-group-in-an-administrative-unit-2) ([#206](https://github.com/manicminer/hamilton/pull/206))

## 0.54.0 (January 13, 2023)

- Support for `Description` field in the `Application` model ([#205](https://github.com/manicminer/hamilton/pull/205))

## 0.53.0 (January 11, 2023)

- Add the App ID for `Microsoft.StorageSync` to `PublishedApis` ([#200](https://github.com/manicminer/hamilton/pull/200))
- Bugfix: Fix casing of values for `OnPremisesGroupType` type ([#199](https://github.com/manicminer/hamilton/pull/199))

## 0.52.0 (November 30, 2022)

- Bugfix: Use `eq` over `startsWith` in `msgraph.AccessPackageResourceClient{}.Get()` to improve accuracy ([#194](https://github.com/manicminer/hamilton/pull/194))
- Support for reading the [`objectId` field](https://learn.microsoft.com/en-us/graph/api/applicationtemplate-instantiate?view=graph-rest-beta&tabs=http#response-1) in API responses for `directoryObjects`
- Support for `WritebackConfiguration` in the `Group` model ([#197](https://github.com/manicminer/hamilton/pull/197))

⚠️ BREAKING CHANGES:

- The `ID` field of the `DirectoryObject` model has been renamed to `Id` and a method `ID()` has been introduced ([#198](https://github.com/manicminer/hamilton/pull/198))

## 0.51.0 (October 27, 2022)

- Support for [Connected Organizations](https://learn.microsoft.com/en-gb/graph/api/resources/connectedorganization?view=graph-rest-1.0) ([#156](https://github.com/manicminer/hamilton/pull/156))
- Support for listing [transitive group members](https://learn.microsoft.com/en-us/graph/api/group-list-transitivemembers?view=graph-rest-1.0) ([#191](https://github.com/manicminer/hamilton/pull/191))
- Bugfix: Add nil slice check in `AccessPackageResourceClient.Get()` ([#187](https://github.com/manicminer/hamilton/pull/187))
- Bugfix: `AccessPackageResource.Description` has changed from a `*bool` to a `*string` ([#187](https://github.com/manicminer/hamilton/pull/187))

## 0.50.0 (October 14, 2022)

- Environments: add Synapse API in USGovernment ([#186](https://github.com/manicminer/hamilton/pull/186))

## 0.49.0 (September 29, 2022)

- Service Principals: support for the `Oauth2PermissionScopes` field (json `oauth2PermissionScopes`), which is used by the v1.0 API ([#183](https://github.com/manicminer/hamilton/pull/183))

## 0.48.0 (September 29, 2022)

- Bug fix: `SynchronizationTaskExecution.CountEntitled` has changed from a `string` to an `int64` ([#172](https://github.com/manicminer/hamilton/pull/172))
- Support for [B2C User Flows](https://learn.microsoft.com/en-us/graph/api/resources/b2cidentityuserflow?view=graph-rest-beta) ([#179](https://github.com/manicminer/hamilton/pull/179))
- Support for [User Flow Attributes](https://learn.microsoft.com/en-us/graph/api/resources/identityuserflowattribute?view=graph-rest-1.0) ([#182](https://github.com/manicminer/hamilton/pull/182))
- Add an `AdditionalData` field of type `map[string]interface{}` to the `DirectoryObject` model, for returning additional untyped fields ([#171](https://github.com/manicminer/hamilton/pull/171))
- `AppRoleAssignmentsClient.List()` - support odata query parameters ([#181](https://github.com/manicminer/hamilton/pull/181))
- Environments: add new well-known App IDs `MicrosoftAzureFrontDoor and `MicrosoftAzureFrontDoorCdn` ([#175](https://github.com/manicminer/hamilton/pull/175))
- OData: Support for a `ConsistencyLevel` header with the value `session` ([#174](https://github.com/manicminer/hamilton/pull/174))

⚠️ BREAKING CHANGES:

- `Group.GroupTypes` has changed from a `[]GroupType` to a `*[]GroupType` ([#160](https://github.com/manicminer/hamilton/pull/160))
- `Group.ResourceBehaviorOptions` has changed from a `[]GroupResourceBehaviorOption` to a `*[]GroupResourceBehaviorOption` ([#160](https://github.com/manicminer/hamilton/pull/160))
- `Group.ResourceProvisioningOptions` has changed from a `[]GroupResourceProvisioningOption` to a `*[]GroupResourceProvisioningOption` ([#160](https://github.com/manicminer/hamilton/pull/160))

## 0.47.1 (August 30, 2022)

- Bugfix: Add missing configuration checks for OIDC methods in the `auth.Config.NewAuthorizer()` method ([#173](https://github.com/manicminer/hamilton/pull/173))

## 0.47.0 (August 25, 2022)

- Support for OIDC federated authentication by supplying an ID token directly ([#166](https://github.com/manicminer/hamilton/pull/166))
- Support for [Azure AD Synchronization](https://docs.microsoft.com/en-us/graph/api/resources/synchronization-overview?view=graph-rest-beta) ([#167](https://github.com/manicminer/hamilton/pull/167))

## 0.46.0 (April 27, 2022)

- Added Azure Security Insights to `environments.PublishedApis` ([#162](https://github.com/manicminer/hamilton/pull/162))
- Added `linux` to supported `msgraph.ConditionalAccessDevicePlatform` values ([#163](https://github.com/manicminer/hamilton/pull/163))
- Added `SubscribeMembersToCalendarEventsDisabled` to supported `msgraph.GroupResourceBehaviorOption` values ([#163](https://github.com/manicminer/hamilton/pull/163))

## 0.45.0 (April 21, 2022)

⚠️ BREAKING CHANGES:

- Removed `omitempty` from the JSON struct tag for the `Locations` and `Platforms` fields of the `msgraph.ConditionalAccessConditionSet` model ([#161](https://github.com/manicminer/hamilton/pull/161))

## 0.44.0 (April 19, 2022)

- Bugfix: Set the correct URL for `environments.KeyVaultUSGovEndpoint` ([#157](https://github.com/manicminer/hamilton/pull/157))
- Support for [Token Signing Certificates](https://docs.microsoft.com/en-us/graph/api/serviceprincipal-addtokensigningcertificate?view=graph-rest-1.0&tabs=http) for service principals ([#151](https://github.com/manicminer/hamilton/pull/151), [#158](https://github.com/manicminer/hamilton/pull/158))

## 0.43.0 (March 15, 2022)

⚠️ BREAKING CHANGES:

- This release removes the following, which have been replaced by compatible equivalents in the new [hamilton-autorest](https://github.com/manicminer/hamilton-autorest) module ([#154](https://github.com/manicminer/hamilton/pull/154))
  - `auth.AutorestAuthorizerWrapper{}` struct
  - `auth.CachedAuthorizer{}.BearerAuthorizerCallback()` method
  - `auth.CachedAuthorizer{}.WithAuthorization()` method
  - `auth.NewAutorestAuthorizerWrapper()` function
  - `auth.ServicePrincipalToken` interface
  - `environments.EnvironmentFromMetadata()` function

## 0.42.0 (March 9, 2022)

- Broaden the regular expression used for fixing up bad oData IDs when marshaling an `odata.Id` ([#152](https://github.com/manicminer/hamilton/pull/152))
- Support for [Claims Mapping Policies](https://docs.microsoft.com/en-us/graph/api/resources/claimsmappingpolicy?view=graph-rest-1.0) ([#147](https://github.com/manicminer/hamilton/pull/147))

## 0.41.1 (February 3, 2022)

⚠️ BREAKING CHANGES:

- Bug fix: `UnifiedRoleDefinition.Description` has changed from a `*string` to a `*StringNullWhenEmpty` ([#148](https://github.com/manicminer/hamilton/pull/148))
- Bug fix: `UnifiedRolePermission.Condition` has changed from a `*string` to a `*StringNullWhenEmpty` ([#148](https://github.com/manicminer/hamilton/pull/148))

## 0.41.0 (February 3, 2022)

- Support for selecting GitHub OIDC authentication when using the `auth.NewAuthorizer()` helper function ([#145](https://github.com/manicminer/hamilton/pull/145))
- Bump supported Go version to 1.17.6 ([#145](https://github.com/manicminer/hamilton/pull/145))

## 0.40.1 (January 28, 2022)

- Bug fix: Correct the type for `AllowExternalSenders` field in the `Group` model ([#143](https://github.com/manicminer/hamilton/pull/143))
- `GroupsClient{}.Update()` - Don't include the ID in the body when updating a group, as this prevents some Unified group fields from being updated ([#143](https://github.com/manicminer/hamilton/pull/143))

## 0.40.0 (January 26, 2022)

- Add a new authorizer `GitHubOIDCAuthorizer` which supports [OIDC token exchange](https://docs.github.com/en/actions/deployment/security-hardening-your-deployments/about-security-hardening-with-openid-connect) for authenticating to Azure Active Directory ([#142](https://github.com/manicminer/hamilton/pull/142))
- Support v1.0 API for [Entitlement Management](https://docs.microsoft.com/en-us/graph/api/resources/entitlementmanagement-overview?view=graph-rest-1.0) ([#133](https://github.com/manicminer/hamilton/pull/133))
  - `AccessPackageQuestion` model - add the `Choices` and `IsSingleLineQuestion` fields
  - `AccessPackageCatalog` model - add the `State` field
  - `AssignmentReviewSettings` model - add the `IsAccessRecommendationEnabled`, `IsApprovalJustificationRequired` and `AccessReviewTimeoutBehavior` fields
  - `UserSet` model - add the `ManagerLevel` field
  - New model: `AccessPackageMultipleChoiceQuestions`
- Support for [Role Definitions](https://docs.microsoft.com/en-us/graph/api/resources/unifiedroledefinition?view=graph-rest-1.0) via the unified role management endpoint ([#137](https://github.com/manicminer/hamilton/pull/137))
- Support for [Role Assignments](https://docs.microsoft.com/en-us/graph/api/resources/unifiedroleassignment?view=graph-rest-1.0) via the unified role management endpoint ([#137](https://github.com/manicminer/hamilton/pull/137))

⚠️ BREAKING CHANGES:

- `AccessPackage` model - the `CatalogId` field is replaced by the `Catalog` field
- `AssignmentReviewSettings` model - the `RecurrenceType` field now has a custom type
- `AssignmentReviewSettings` model - the `ReviewerType` field now has a custom type

## 0.39.0 (January 7, 2022)

- Support for [Federated Identity Credentials](https://docs.microsoft.com/en-us/graph/api/resources/federatedidentitycredential?view=graph-rest-beta) (beta-only) ([#134](https://github.com/manicminer/hamilton/pull/134))
- Bug fix: corrected the `DisplayName` struct tag for the `GroupAssignedLabel` model ([#135](https://github.com/manicminer/hamilton/pull/135))
- Bug fix: fixed a typo in the constant `AccessPackageResourceRequestTypeAdminRemove` (was `AccessPackageResourceRequestTypeAdmminRemove`) ([#135](https://github.com/manicminer/hamilton/pull/135))

## 0.38.0 (December 8, 2021)

- Add a helper function `environments.EnvironmentFromMetadata()` which is intended to substitute the [`azure.EnvironmentFromURL()` function from go-autorest](https://github.com/Azure/go-autorest/blob/v14.2.0/autorest/azure/metadata_environment.go#L96-L141) ([#131](https://github.com/manicminer/hamilton/pull/131))
- Fix an incorrect API ID for KeyVault ([#131](https://github.com/manicminer/hamilton/pull/131))
- Improve support for dynamic group memberships ([#132](https://github.com/manicminer/hamilton/pull/132))

⚠️ BREAKING CHANGES:

- Bug fix: `Group.MembershipRule` has changed from a `*string` to a `*StringNullWhenEmpty` ([#132](https://github.com/manicminer/hamilton/pull/132))

## 0.37.0 (November 29, 2021)

- Add some missing API endpoints for national cloud environments ([#129](https://github.com/manicminer/hamilton/pull/129))
- Add an `Api{}.IsAvailable()` method to determine whether a service is supported for an environment  ([#129](https://github.com/manicminer/hamilton/pull/129))
- Fix an incorrect hostname for `environments.KeyVaultUSGovEndpoint` ([#128](https://github.com/manicminer/hamilton/pull/128))
- Support for `autorest.BearerAuthorizerCallback` in `auth.CachedAuthorizer` ([#130](https://github.com/manicminer/hamilton/pull/130))

## 0.36.1 (November 25, 2021)

- Fix an incorrect enum value for `ConditionalAccessDevicePlatformAll` ([#127](https://github.com/manicminer/hamilton/pull/127))

## 0.36.0 (November 25, 2021)

- Support for [administrative units](https://docs.microsoft.com/en-us/graph/api/resources/administrativeunit?view=graph-rest-beta) ([#124](https://github.com/manicminer/hamilton/pull/124))
- Support for [delegated permission grants](https://docs.microsoft.com/en-us/graph/api/resources/oauth2permissiongrant?view=graph-rest-1.0) ([#126](https://github.com/manicminer/hamilton/pull/126))
- Conditional Access Policies: support for `devices` and `deviceStates` in policy `conditions` ([#125](https://github.com/manicminer/hamilton/pull/125))
- Conditional Access Policies: add type aliases and constants for enum values ([#125](https://github.com/manicminer/hamilton/pull/125))

## 0.35.0 (November 16, 2021)

- Auth package refactoring ([#123](https://github.com/manicminer/hamilton/pull/123))
  - Remove the `auth.Api` type and instead use `environments.Api` directly
  - Use the resource URI instead of the friendly name for Azure CLI auth tokens

- Add the `AuxiliaryTokens()` method to the `auth.Authorizer` interface to support obtaining tokens for additional tenants ([#123](https://github.com/manicminer/hamilton/pull/123))
- Expand support in `auth.AutorestAuthorizerWrapper` to support any `autorest.Authorizer` ([#123](https://github.com/manicminer/hamilton/pull/123))
  - `autorest.BearerAuthorizer` and `autorest.MultiTenantBearerAuthorizer` are fully supported with access tokens, refresh tokens and expiry
  - Other authorizers can supply access tokens only
- Support auxiliary tenants with client secret and client certificate authorizers ([#123](https://github.com/manicminer/hamilton/pull/123))

- Implement the `autorest.Authorizer` interface with `auth.CachedAuthorizer` (which wraps all supported Authorizers) ([#123](https://github.com/manicminer/hamilton/pull/123))
  - This allows authorizers to be used with https://github.com/Azure/go-autorest, with multi-tenant support, with the exception of `auth.MsiAuthorizer`

- Export environment configs for more management plane APIs ([#123](https://github.com/manicminer/hamilton/pull/123))
  - Resource Manager
  - Batch Management
  - Data Lake
  - Gallery
  - KeyVault
  - Operational Insights
  - OSS RDBMS
  - Service Bus
  - Service Management (Azure Classic)
  - SQL Database
  - Storage
  - Synapse

- Refactor and tidy up tests for the `msgraph` package ([#123](https://github.com/manicminer/hamilton/pull/123))

- Say goodbye to Azure Germany 🇩🇪 ([#123](https://github.com/manicminer/hamilton/pull/123))

⚠️ BREAKING CHANGES:

- The signatures for `auth.NewClientCertificateAuthorizer`, `auth.NewClientSecretAuthorizer` and `auth.NewAzureCliAuthorizer` have changed to accommodate passing additional tenant IDs for multi-tenant authorization ([#123](https://github.com/manicminer/hamilton/pull/123))

## 0.34.0 (November 12, 2021)

- Remove a surplus configuration check when using Managed Identity authentication, which improves compatibility with Azure Cloud Shell ([#119](https://github.com/manicminer/hamilton/pull/119))
- Add a new authorizer `AutorestAuthorizerWrapper` which supports obtaining tokens from go-autorest via `autorest.BearerAuthorizer` ([#120](https://github.com/manicminer/hamilton/pull/120))

## 0.33.0 (October 14, 2021)

- Support for specifying the client ID when using managed identity authentication ([#115](https://github.com/manicminer/hamilton/pull/115))
- Mitigation for breaking API changes around the `@odata.id` field ([#114](https://github.com/manicminer/hamilton/pull/114))
  - If `@odata.id` is returned in the form `objectType('GUID')` (i.e. not a valid URI), then attempt to reconstruct a URI
  - This currently hardcodes the `graph.microsoft.com` host in the generated URI but this does not appear to be a problem for other clouds
  - This field is exported in all structs that reference it, so it's possible to override this if necessary
- Support for running `msgraph` tests in national clouds ([#114](https://github.com/manicminer/hamilton/pull/114))

⚠️ BREAKING CHANGES:

- The signatures for the `auth.NewMsiAuthorizer()` and `auth.NewMsiConfig()` functions have changed to accommodate the client ID ([#115](https://github.com/manicminer/hamilton/pull/115))

## 0.32.0 (October 6, 2021)

- Support for setting OData-related HTTP headers
  - Implement a new way to pass the entire `odata.Query` object as part of request inputs
  - Update all existing clients to pass `odata.Query` in full
    - The existing method of passing a `url.Values`map still works, maintains compatibility and can be used for passing non-odata related query parameters
  - Support setting [OData-related HTTP headers](https://docs.oasis-open.org/odata/odata/v4.01/odata-v4.01-part1-protocol.html#_Toc31358856) including `OData-Version` / `OData-MaxVersion`, [odata-json parameters](https://docs.oasis-open.org/odata/odata-json-format/v4.01/odata-json-format-v4.01.html#_Toc38457724) on the `Accept` header, and [the `ConsistencyLevel` header](https://developer.microsoft.com/en-us/identity/blogs/microsoft-graph-advanced-queries-for-directory-objects-are-now-generally-available/) which isn't strictly in the odata 'standard' but heavily related
  - Set the `odata.metadata` parameter to `full` when retrieving directory objects to [ensure the `@odata.id` field is returned](https://docs.oasis-open.org/odata/odata-json-format/v4.0/errata03/os/odata-json-format-v4.0-errata03-os-complete.html#_Toc453766619)
    - This dramatically increases the payload volume so we don't default it everywhere
- Support for [assigning](https://docs.microsoft.com/en-us/graph/api/user-post-manager?view=graph-rest-beta&tabs=http) and [retrieving](https://docs.microsoft.com/en-us/graph/api/user-list-manager?view=graph-rest-beta&tabs=http) a user's manager ([#111](https://github.com/manicminer/hamilton/pull/111))
- Add application ID for "Azure VPN" to environments package ([#113](https://github.com/manicminer/hamilton/pull/113))

## 0.31.1 (September 30, 2021)

- Bug fix: `User{}.EmployeeType` is a nullable string ([#110](https://github.com/manicminer/hamilton/pull/110))

## 0.31.0 (September 30, 2021)

- Add support for [Entitlement Management](https://docs.microsoft.com/en-us/graph/api/resources/entitlementmanagement-root?view=graph-rest-beta) (beta-only API) ([#93](https://github.com/manicminer/hamilton/pull/93))
- Bug fix: handle inconsistent 400 error when listing sign-in reports with an OData filter ([#108](https://github.com/manicminer/hamilton/pull/108))
- Bug fix: work around an API consistency issue when creating service principals for new applications that have not fully replicated ([#109](https://github.com/manicminer/hamilton/pull/109))

## 0.30.0 (September 22, 2021)

- Support for the [appRolesAssignedTo](https://docs.microsoft.com/en-us/graph/api/serviceprincipal-post-approleassignedto?view=graph-rest-1.0) endpoint ([#107](https://github.com/manicminer/hamilton/pull/107))
- Bug fix: `odata.Odata{}` - the `Count` field is now a `*int` ([#105](https://github.com/manicminer/hamilton/pull/105))

## 0.29.0 (September 15, 2021)

- Applications: add consistency check for roles/scopes that may be in the process of being disabled, when updating an application ([#102](https://github.com/manicminer/hamilton/pull/102))
- Applications: support for uploading application logos via the `ApplicationsClient{}.UploadLogo()` method([#103](https://github.com/manicminer/hamilton/pull/103))
- Directory Roles: add the `DirectoryROlesClient{}.GetByTemplateId()` method for retrieving roles by their template ID ([#101](https://github.com/manicminer/hamilton/pull/101))
- `User` model: support [EmployeeOrgData](https://docs.microsoft.com/en-us/graph/api/resources/employeeorgdata?view=graph-rest-beta) ([#99](https://github.com/manicminer/hamilton/pull/99))

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

