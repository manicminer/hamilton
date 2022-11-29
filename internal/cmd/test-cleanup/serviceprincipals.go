package main

import (
	"fmt"
	"log"

	"github.com/manicminer/hamilton/msgraph"
	"github.com/manicminer/hamilton/odata"
)

func cleanupServicePrincipals() {
	servicePrincipalsClient := msgraph.NewServicePrincipalsClient(tenantId)
	servicePrincipalsClient.BaseClient.Authorizer = authorizer

	servicePrincipals, _, err := servicePrincipalsClient.List(ctx, odata.Query{Filter: fmt.Sprintf("startsWith(displayName, '%s')", displayNamePrefix)})
	if err != nil {
		log.Println(err)
		return
	}
	if servicePrincipals == nil {
		log.Println("bad API response, nil ServicePrincipals result received")
		return
	}
	for _, servicePrincipal := range *servicePrincipals {
		if servicePrincipal.ID() == nil || servicePrincipal.DisplayName == nil {
			log.Println("Service Principal returned with nil ID or DisplayName")
			continue
		}

		log.Printf("Deleting service principal %q (DisplayName: %q)\n", *servicePrincipal.ID(), *servicePrincipal.DisplayName)
		_, err := servicePrincipalsClient.Delete(ctx, *servicePrincipal.ID())
		if err != nil {
			log.Printf("Error when deleting service principal %q: %v\n", *servicePrincipal.ID(), err)
		}
	}
}
