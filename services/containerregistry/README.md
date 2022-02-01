# Azure Container Registry

This is a small Azure Container Registry client for Go based on the [documentation from Microsoft](https://docs.microsoft.com/en-us/rest/api/containerregistry/). The idea is to use the standard library as far as it's feasible and to be an alternative to using [autorest](https://github.com/Azure/go-autorest) and [Azure SDK for Go](https://github.com/Azure/azure-sdk-for-go).

An alternative to this is using the [azure-sdk-fo-go](https://github.com/Azure/azure-sdk-for-go/tree/main/services/preview/containerregistry/runtime/2019-08-15-preview/containerregistry) implementation of the API.

## Implemented APIs

The following APIs are implemented.

### Access Token

| Name                                                                                 | Method                  |
| ------------------------------------------------------------------------------------ | ----------------------- |
| [Get](https://docs.microsoft.com/en-us/rest/api/containerregistry/access-tokens/get) | `ExchangeAccessToken()` |

### Refresh Token

| Name                                                                                                              | Method                   |
| ----------------------------------------------------------------------------------------------------------------- | ------------------------ |
| [Get From Exchange](https://docs.microsoft.com/en-us/rest/api/containerregistry/refresh-tokens/get-from-exchange) | `ExchangeRefreshToken()` |

### Repository (Catalog)

| Name                                                                                                          | Method                      |
| ------------------------------------------------------------------------------------------------------------- | --------------------------- |
| [Get List](https://docs.microsoft.com/en-us/rest/api/containerregistry/repository/get-list)                   | `CatalogList()`             |
| [Get Attributes](https://docs.microsoft.com/en-us/rest/api/containerregistry/repository/get-attributes)       | `CatalogGetAttributes()`    |
| [Update Attributes](https://docs.microsoft.com/en-us/rest/api/containerregistry/repository/update-attributes) | `CatalogUpdateAttributes()` |
| [Delete](https://docs.microsoft.com/en-us/rest/api/containerregistry/repository/delete)                       | `CatalogDelete()`           |

### Tag

| Name                                                                                                   | Method                  |
| ------------------------------------------------------------------------------------------------------ | ----------------------- |
| [Get List](https://docs.microsoft.com/en-us/rest/api/containerregistry/tag/get-list)                   | `TagList()`             |
| [Get Attributes](https://docs.microsoft.com/en-us/rest/api/containerregistry/tag/get-attributes)       | `TagGetAttributes()`    |
| [Update Attributes](https://docs.microsoft.com/en-us/rest/api/containerregistry/tag/update-attributes) | `TagUpdateAttributes()` |
| [Delete](https://docs.microsoft.com/en-us/rest/api/containerregistry/tag/delete)                       | `TagDelete()`           |

### Manifest

| Name                                                                                                         | Method                       |
| ------------------------------------------------------------------------------------------------------------ | ---------------------------- |
| [Get List](https://docs.microsoft.com/en-us/rest/api/containerregistry/manifests/get-list)                   | `ManifestList()`             |
| [Get](https://docs.microsoft.com/en-us/rest/api/containerregistry/manifests/get)                             | `ManifestGet()`              |
| [Get Attributes](https://docs.microsoft.com/en-us/rest/api/containerregistry/manifests/get-attributes)       | `ManifestGetAttributes()`    |
| [Update Attributes](https://docs.microsoft.com/en-us/rest/api/containerregistry/manifests/update-attributes) | `ManifestUpdateAttributes()` |
| [Delete](https://docs.microsoft.com/en-us/rest/api/containerregistry/manifests/delete)                       | `ManifestDelete()`           |

## Not implemented APIs (yet)

The following APIs haven't been implemented yet.

### Manifest

| Name                                                                                   |
| -------------------------------------------------------------------------------------- |
| [Create](https://docs.microsoft.com/en-us/rest/api/containerregistry/manifests/create) |

### Blob

| Name                                                                                            |
| ----------------------------------------------------------------------------------------------- |
| [Cancel Upload](https://docs.microsoft.com/en-us/rest/api/containerregistry/blob/cancel-upload) |
| [Check](https://docs.microsoft.com/en-us/rest/api/containerregistry/blob/check)                 |
| [Check Chunk](https://docs.microsoft.com/en-us/rest/api/containerregistry/blob/check-chunk)     |
| [Delete](https://docs.microsoft.com/en-us/rest/api/containerregistry/blob/delete)               |
| [End Upload](https://docs.microsoft.com/en-us/rest/api/containerregistry/blob/end-upload)       |
| [Get](https://docs.microsoft.com/en-us/rest/api/containerregistry/blob/get)                     |
| [Get Chunk](https://docs.microsoft.com/en-us/rest/api/containerregistry/blob/get-chunk)         |
| [Get Status](https://docs.microsoft.com/en-us/rest/api/containerregistry/blob/get-status)       |
| [Mount](https://docs.microsoft.com/en-us/rest/api/containerregistry/blob/mount)                 |
| [Start Upload](https://docs.microsoft.com/en-us/rest/api/containerregistry/blob/start-upload)   |
| [Upload](https://docs.microsoft.com/en-us/rest/api/containerregistry/blob/upload)               |
