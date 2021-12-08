package environments

import (
	"reflect"
	"testing"
)

func TestEnvironmentFromMetadata(t *testing.T) {
	env, err := EnvironmentFromMetadata("https://management.azure.com/")
	if err != nil {
		t.Fatal(err)
	}
	if env == nil {
		t.Fatal("env was nil")
	}

	expected := Environment{
		AzureADEndpoint: AzureADEndpoint("https://login.microsoftonline.com"),
		AadGraph: Api{
			AppId:    PublishedApis["AzureActiveDirectoryGraph"],
			Endpoint: ApiEndpoint("https://graph.windows.net"),
		},
		KeyVault: Api{
			AppId:    PublishedApis["AzureKeyVault"],
			Endpoint: ApiEndpoint("https://vault.azure.com"),
		},
		ResourceManager: Api{
			AppId:    PublishedApis["AzureServiceManagement"],
			Endpoint: "https://management.azure.com",
		},
		Storage: Api{
			AppId:    PublishedApis["AzureStorage"],
			Endpoint: ApiEndpoint("https://storage.azure.com"),
		},
	}

	if !reflect.DeepEqual(*env, expected) {
		t.Fatalf("expected: %#v, got: %#v", expected, *env)
	}
}
