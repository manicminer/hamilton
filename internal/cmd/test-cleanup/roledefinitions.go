package main

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/manicminer/hamilton/msgraph"
)

func cleanupRoleDefinitions() {
	client := msgraph.NewRoleDefinitionsClient(tenantId)
	client.BaseClient.Authorizer = authorizer

	result, _, err := client.List(ctx, odata.Query{Filter: fmt.Sprintf("startsWith(displayName, '%s')", displayNamePrefix)})
	if err != nil {
		log.Println(err)
		return
	}
	if result == nil {
		log.Println("bad API response, nil role definitions result received")
		return
	}
	for _, item := range *result {
		if item.ID() == nil || item.DisplayName == nil {
			log.Println("Role definition returned with nil ID or DisplayName")
			continue
		}

		log.Printf("Deleting role definition %q (DisplayName: %q)\n", *item.ID(), *item.DisplayName)
		_, err := client.Delete(ctx, *item.ID())
		if err != nil {
			log.Printf("Error when deleting role definition %q: %v\n", *item.ID(), err)
		}
	}
}
