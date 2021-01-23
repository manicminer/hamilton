package environments

type Environment struct {
	AADGraphEndpoint AADGraphEndpoint
	AzureADEndpoint  AzureADEndpoint
	MsGraphEndpoint  MsGraphEndpoint
}

var (
	Global = Environment{
		AADGraphEndpoint: AADGraphGlobal,
		AzureADEndpoint:  AzureADGlobal,
		MsGraphEndpoint:  MsGraphGlobal,
	}

	USGovernmentL4 = Environment{
		AADGraphEndpoint: AADGraphUSGov,
		AzureADEndpoint:  AzureADUSGov,
		MsGraphEndpoint:  MsGraphUSGovL4,
	}

	USGovernmentL5 = Environment{
		AADGraphEndpoint: AADGraphUSGov,
		AzureADEndpoint:  AzureADUSGov,
		MsGraphEndpoint:  MsGraphUSGovL5,
	}

	Germany = Environment{
		AADGraphEndpoint: AADGraphGermany,
		AzureADEndpoint:  AzureADGermany,
		MsGraphEndpoint:  MsGraphGermany,
	}

	China = Environment{
		AADGraphEndpoint: AADGraphChina,
		AzureADEndpoint:  AzureADChina,
		MsGraphEndpoint:  MsGraphChina,
	}

	Canary = Environment{
		AzureADEndpoint: AzureADGlobal,
		MsGraphEndpoint: MsGraphCanary,
	}
)
