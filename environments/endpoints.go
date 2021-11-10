package environments

type AzureADEndpoint string

const (
	AzureADGlobal  AzureADEndpoint = "https://login.microsoftonline.com"
	AzureADUSGov   AzureADEndpoint = "https://login.microsoftonline.us"
	AzureADGermany AzureADEndpoint = "https://login.microsoftonline.de"
	AzureADChina   AzureADEndpoint = "https://login.chinacloudapi.cn"
)

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
