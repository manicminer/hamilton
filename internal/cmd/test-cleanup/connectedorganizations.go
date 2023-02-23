package main

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/manicminer/hamilton/msgraph"
)

func cleanupConnectedOrganizations() {
	client := msgraph.NewConnectedOrganizationClient(tenantId)
	client.BaseClient.Authorizer = authorizer

	result, _, err := client.List(ctx, odata.Query{Filter: fmt.Sprintf("startsWith(displayName, '%s')", displayNamePrefix)})
	if err != nil {
		log.Println(err)
		return
	}
	if result == nil {
		log.Println("bad API response, nil connected organizations result received")
		return
	}
	for _, item := range *result {
		if item.ID == nil || item.DisplayName == nil {
			log.Println("Connected organization returned with nil ID or DisplayName")
			continue
		}

		log.Printf("Deleting connected organization %q (DisplayName: %q)\n", *item.ID, *item.DisplayName)
		_, err := client.Delete(ctx, *item.ID)
		if err != nil {
			log.Printf("Error when deleting connected organization %q: %v\n", *item.ID, err)
		}
	}
}
