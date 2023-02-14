package main

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/manicminer/hamilton/msgraph"
)

func cleanupNamedLocations() {
	namedLocationsClient := msgraph.NewNamedLocationsClient(tenantId)
	namedLocationsClient.BaseClient.Authorizer = authorizer

	namedLocations, _, err := namedLocationsClient.List(ctx, odata.Query{Filter: fmt.Sprintf("startsWith(displayName, '%s')", displayNamePrefix)})
	if err != nil {
		log.Println(err)
		return
	}
	if namedLocations == nil {
		log.Println("bad API response, nil namedLocations result received")
		return
	}
	for _, namedLocation := range *namedLocations {
		var id string
		if countryNamedLocation, ok := namedLocation.(msgraph.CountryNamedLocation); ok {
			if countryNamedLocation.ID == nil {
				log.Println("Country Named Location returned with nil ID")
				continue
			}
			id = *countryNamedLocation.ID
		} else if ipNamedLocation, ok := namedLocation.(msgraph.IPNamedLocation); ok {
			if ipNamedLocation.ID == nil {
				log.Println("IP Named Location returned with nil ID")
				continue
			}
			id = *ipNamedLocation.ID
		}

		log.Printf("Deleting named location %q\n", id)
		_, err := namedLocationsClient.Delete(ctx, id)
		if err != nil {
			log.Printf("Error when deleting named location %q: %v\n", id, err)
		}
	}
}
