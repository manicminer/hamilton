package environments

type Environment struct {
	AzureADEndpoint AzureADEndpoint
	MsGraphEndpoint MsGraphEndpoint
}

var (
	Global = Environment{
		AzureADEndpoint: AzureADGlobal,
		MsGraphEndpoint: MsGraphGlobal,
	}

	USGovernmentL4 = Environment{
		AzureADEndpoint: AzureADUSGov,
		MsGraphEndpoint: MsGraphUSGovL4,
	}

	USGovernmentL5 = Environment{
		AzureADEndpoint: AzureADUSGov,
		MsGraphEndpoint: MsGraphUSGovL5,
	}

	Germany = Environment {
		AzureADEndpoint: AzureADGermany,
		MsGraphEndpoint: MsGraphGermany,
	}

	China = Environment {
		AzureADEndpoint: AzureADChina,
		MsGraphEndpoint: MsGraphChina,
	}
)