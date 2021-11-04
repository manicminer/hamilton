package environments

type ApiEndpoint string

const (
	MsGraphGlobalEndpoint  ApiEndpoint = "https://graph.microsoft.com"
	MsGraphGermanyEndpoint ApiEndpoint = "https://graph.microsoft.de"
	MsGraphChinaEndpoint   ApiEndpoint = "https://microsoftgraph.chinacloudapi.cn"
	MsGraphUSGovL4Endpoint ApiEndpoint = "https://graph.microsoft.us"
	MsGraphUSGovL5Endpoint ApiEndpoint = "https://dod-graph.microsoft.us"
	MsGraphCanaryEndpoint  ApiEndpoint = "https://canary.graph.microsoft.com"

	AadGraphGlobalEndpoint  ApiEndpoint = "https://graph.windows.net"
	AadGraphGermanyEndpoint ApiEndpoint = "https://graph.cloudapi.de"
	AadGraphChinaEndpoint   ApiEndpoint = "https://graph.chinacloudapi.cn"
	AadGraphUSGovEndpoint   ApiEndpoint = "https://graph.microsoftazure.us"

	ResourceManagerPublicEndpoint  ApiEndpoint = "https://management.azure.com"
	ResourceManagerGermanyEndpoint ApiEndpoint = "https://management.microsoftazure.de"
	ResourceManagerChinaEndpoint   ApiEndpoint = "https://management.chinacloudapi.cn"
	ResourceManagerUSGovEndpoint   ApiEndpoint = "https://management.usgovcloudapi.net"

	BatchManagementPublicEndpoint  ApiEndpoint = "https://batch.core.windows.net"
	BatchManagementGermanyEndpoint ApiEndpoint = "https://batch.cloudapi.de"
	BatchManagementChinaEndpoint   ApiEndpoint = "https://batch.chinacloudapi.cn"
	BatchManagementUSGovEndpoint   ApiEndpoint = "https://batch.core.usgovcloudapi.net"

	DataLakePublicEndpoint ApiEndpoint = "https://datalake.azure.net"

	GalleryPublicEndpoint  ApiEndpoint = "https://gallery.azure.com"
	GalleryGermanyEndpoint ApiEndpoint = "https://gallery.cloudapi.de"
	GalleryChinaEndpoint   ApiEndpoint = "https://gallery.chinacloudapi.cn"
	GalleryUSGovEndpoint   ApiEndpoint = "https://gallery.usgovcloudapi.net"

	KeyVaultPublicEndpoint  ApiEndpoint = "https://vault.azure.net"
	KeyVaultGermanyEndpoint ApiEndpoint = "https://vault.azure.net"
	KeyVaultChinaEndpoint   ApiEndpoint = "https://vault.azure.cn"
	KeyVaultUSGovEndpoint   ApiEndpoint = "https://vault.microsoftazure.de"

	OperationalInsightsPublicEndpoint ApiEndpoint = "https://api.loganalytics.io"
	OperationalInsightsUSGovEndpoint  ApiEndpoint = "https://api.loganalytics.us"

	OSSRDBMSPublicEndpoint  ApiEndpoint = "https://ossrdbms-aad.database.windows.net"
	OSSRDBMSGermanyEndpoint ApiEndpoint = "https://ossrdbms-aad.database.cloudapi.de"
	OSSRDBMSChinaEndpoint   ApiEndpoint = "https://ossrdbms-aad.database.chinacloudapi.cn"
	OSSRDBMSUSGovEndpoint   ApiEndpoint = "https://ossrdbms-aad.database.usgovcloudapi.net"

	ServiceBusPublicEndpoint  ApiEndpoint = "https://servicebus.windows.net"
	ServiceBusGermanyEndpoint ApiEndpoint = "https://servicebus.cloudapi.de"
	ServiceBusChinaEndpoint   ApiEndpoint = "https://servicebus.chinacloudapi.cn"
	ServiceBusUSGovEndpoint   ApiEndpoint = "https://servicebus.usgovcloudapi.net"

	ServiceManagementPublicEndpoint  ApiEndpoint = "https://management.core.windows.net"
	ServiceManagementGermanyEndpoint ApiEndpoint = "https://management.core.cloudapi.de"
	ServiceManagementChinaEndpoint   ApiEndpoint = "https://management.core.chinacloudapi.cn"
	ServiceManagementUSGovEndpoint   ApiEndpoint = "https://management.core.usgovcloudapi.net"

	SQLDatabasePublicEndpoint  ApiEndpoint = "https://database.windows.net"
	SQLDatabaseGermanyEndpoint ApiEndpoint = "https://database.cloudapi.de"
	SQLDatabaseChinaEndpoint   ApiEndpoint = "https://database.chinacloudapi.cn"
	SQLDatabaseUSGovEndpoint   ApiEndpoint = "https://database.usgovcloudapi.net"

	StoragePublicEndpoint ApiEndpoint = "https://storage.azure.com"

	SynapsePublicEndpoint ApiEndpoint = "https://dev.azuresynapse.net"
)

type ApiCliName string

const (
	MsGraphCliName         ApiCliName = "ms-graph"
	AadGraphCliName        ApiCliName = "aad-graph"
	ResourceManagerCliName ApiCliName = "arm"
	BatchCliName           ApiCliName = "batch"
	DataLakeCliName        ApiCliName = "data-lake"
	OSSRDBMSCliName        ApiCliName = "oss-rdbms"
)

// API represent an API configuration for Microsoft Graph or Azure Active Directory Graph.
type Api struct {
	// The Application ID for the API.
	AppId ApiAppId

	// The Azure CLI codename for the API. Used with `az account get-access-token`.
	CliName ApiCliName

	// The endpoint for the API, including scheme.
	Endpoint ApiEndpoint
}

var (
	MsGraphGlobal = Api{
		AppId:    PublishedApis["MicrosoftGraph"],
		CliName:  MsGraphCliName,
		Endpoint: MsGraphGlobalEndpoint,
	}

	MsGraphGermany = Api{
		AppId:    PublishedApis["MicrosoftGraph"],
		CliName:  MsGraphCliName,
		Endpoint: MsGraphGermanyEndpoint,
	}

	MsGraphChina = Api{
		AppId:    PublishedApis["MicrosoftGraph"],
		CliName:  MsGraphCliName,
		Endpoint: MsGraphChinaEndpoint,
	}

	MsGraphUSGovL4 = Api{
		AppId:    PublishedApis["MicrosoftGraph"],
		CliName:  MsGraphCliName,
		Endpoint: MsGraphUSGovL4Endpoint,
	}

	MsGraphUSGovL5 = Api{
		AppId:    PublishedApis["MicrosoftGraph"],
		CliName:  MsGraphCliName,
		Endpoint: MsGraphUSGovL5Endpoint,
	}

	MsGraphCanary = Api{
		AppId:    PublishedApis["MicrosoftGraph"],
		CliName:  MsGraphCliName,
		Endpoint: MsGraphCanaryEndpoint,
	}

	AadGraphGlobal = Api{
		AppId:    PublishedApis["AzureActiveDirectoryGraph"],
		CliName:  AadGraphCliName,
		Endpoint: AadGraphGlobalEndpoint,
	}

	AadGraphGermany = Api{
		AppId:    PublishedApis["AzureActiveDirectoryGraph"],
		CliName:  AadGraphCliName,
		Endpoint: AadGraphGermanyEndpoint,
	}

	AadGraphChina = Api{
		AppId:    PublishedApis["AzureActiveDirectoryGraph"],
		CliName:  AadGraphCliName,
		Endpoint: AadGraphChinaEndpoint,
	}

	AadGraphUSGov = Api{
		AppId:    PublishedApis["AzureActiveDirectoryGraph"],
		CliName:  AadGraphCliName,
		Endpoint: AadGraphUSGovEndpoint,
	}

	ResourceManagerPublic = Api{
		AppId:    PublishedApis["AzureServiceManagement"],
		CliName:  ResourceManagerCliName,
		Endpoint: ResourceManagerPublicEndpoint,
	}

	ResourceManagerGermany = Api{
		AppId:    PublishedApis["AzureServiceManagement"],
		CliName:  ResourceManagerCliName,
		Endpoint: ResourceManagerGermanyEndpoint,
	}

	ResourceManagerChina = Api{
		AppId:    PublishedApis["AzureServiceManagement"],
		CliName:  ResourceManagerCliName,
		Endpoint: ResourceManagerChinaEndpoint,
	}

	ResourceManagerUSGov = Api{
		AppId:    PublishedApis["AzureServiceManagement"],
		CliName:  ResourceManagerCliName,
		Endpoint: ResourceManagerUSGovEndpoint,
	}
)
