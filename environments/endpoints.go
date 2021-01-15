package environments

import (
	"fmt"
	"golang.org/x/oauth2"
)

type AzureADEndpoint string

const (
	AzureADGlobal  AzureADEndpoint = "https://login.microsoftonline.com"
	AzureADUSGov   AzureADEndpoint = "https://login.microsoftonline.us"
	AzureADGermany AzureADEndpoint = "https://login.microsoftonline.de"
	AzureADChina   AzureADEndpoint = "https://login.chinacloudapi.cn"
)

func AzureAD(endpoint AzureADEndpoint, tenant string) oauth2.Endpoint {
	if tenant == "" {
		tenant = "common"
	}
	return oauth2.Endpoint{
		AuthURL:  fmt.Sprintf("%s/%s/oauth2/v2.0/authorize", endpoint, tenant),
		TokenURL: fmt.Sprintf("%s/%s/oauth2/v2.0/token", endpoint, tenant),
	}
}

type MsGraphEndpoint string

const (
	MsGraphGlobal  MsGraphEndpoint = "https://graph.microsoft.com"
	MsGraphUSGovL4 MsGraphEndpoint = "https://graph.microsoft.us"
	MsGraphUSGovL5 MsGraphEndpoint = "https://dod-graph.microsoft.us"
	MsGraphGermany MsGraphEndpoint = "https://graph.microsoft.de"
	MsGraphChina   MsGraphEndpoint = "https://microsoftgraph.chinacloudapi.cn"
)
