package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/manicminer/hamilton/msgraph"
)

func cleanupAccessPackageCatalogs() {
	client := msgraph.NewAccessPackageCatalogClient(tenantId)
	client.BaseClient.Authorizer = authorizer

	result, _, err := client.List(ctx, odata.Query{Filter: fmt.Sprintf("startsWith(displayName, '%s')", displayNamePrefix)})
	if err != nil {
		log.Println(err)
		return
	}
	if result == nil {
		log.Println("bad API response, nil access package catalogs result received")
		return
	}
	for _, item := range *result {
		if item.ID == nil || item.DisplayName == nil {
			log.Println("Access package catalog returned with nil ID or DisplayName")
			continue
		}

		if strings.HasPrefix(strings.ToLower(*item.DisplayName), strings.ToLower(displayNamePrefix)) {
			log.Printf("Deleting access package catalog %q (DisplayName: %q)\n", *item.ID, *item.DisplayName)
			_, err := client.Delete(ctx, *item.ID)
			if err != nil {
				log.Printf("Error when deleting access package catalog %q: %v\n", *item.ID, err)
			}
		}
	}
}
