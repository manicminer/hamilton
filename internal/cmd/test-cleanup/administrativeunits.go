package main

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/manicminer/hamilton/msgraph"
)

func cleanupAdministrativeUnits() {
	administrativeUnitsClient := msgraph.NewAdministrativeUnitsClient(tenantId)
	administrativeUnitsClient.BaseClient.Authorizer = authorizer

	administrativeUnits, _, err := administrativeUnitsClient.List(ctx, odata.Query{Filter: fmt.Sprintf("startsWith(displayName, '%s')", displayNamePrefix)})
	if err != nil {
		log.Println(err)
		return
	}
	if administrativeUnits == nil {
		log.Println("bad API response, nil administrativeUnits result received")
		return
	}
	for _, au := range *administrativeUnits {
		if au.ID == nil || au.DisplayName == nil {
			log.Println("Group returned with nil ID or DisplayName")
			continue
		}

		log.Printf("Deleting administrative unit %q (DisplayName: %q)\n", *au.ID, *au.DisplayName)
		_, err := administrativeUnitsClient.Delete(ctx, *au.ID)
		if err != nil {
			log.Printf("Error when deleting administrative group %q: %v\n", *au.ID, err)
		}
	}
}
