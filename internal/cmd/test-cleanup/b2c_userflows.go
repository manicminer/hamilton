package main

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/manicminer/hamilton/msgraph"
)

func cleanupB2CUserFlows() {
	client := msgraph.NewB2CUserFlowClient(b2cTenantId)
	client.BaseClient.Authorizer = b2cAuthorizer

	result, _, err := client.List(ctx, odata.Query{Filter: fmt.Sprintf("startsWith(id, 'B2C_1_%s')", displayNamePrefix)})
	if err != nil {
		log.Println(err)
		return
	}
	if result == nil {
		log.Println("bad API response, nil B2C userflows result received")
		return
	}
	for _, item := range *result {
		if item.ID == nil {
			log.Println("B2C userflow returned with nil ID")
			continue
		}

		log.Printf("Deleting B2C userflow %q\n", *item.ID)
		_, err := client.Delete(ctx, *item.ID)
		if err != nil {
			log.Printf("Error when deleting B2C userflow %q: %v\n", *item.ID, err)
		}
	}
}
