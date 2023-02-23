package main

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/manicminer/hamilton/msgraph"
)

func cleanupSchemaExtensions() {
	schemaExtensionsClient := msgraph.NewSchemaExtensionsClient(tenantId)
	schemaExtensionsClient.BaseClient.Authorizer = authorizer

	schemaExtensions, _, err := schemaExtensionsClient.List(ctx, odata.Query{Filter: fmt.Sprintf("status eq '%s'", msgraph.SchemaExtensionStatusInDevelopment)})
	if err != nil {
		log.Println(err)
		return
	}
	if schemaExtensions == nil {
		log.Println("bad API response, nil SchemaExtensions result received")
		return
	}
	for _, schemaExtension := range *schemaExtensions {
		if schemaExtension.ID == nil {
			log.Println("Schema Extensions returned with nil ID")
			continue
		}

		log.Printf("Deleting schema extension %q\n", *schemaExtension.ID)
		_, err := schemaExtensionsClient.Delete(ctx, *schemaExtension.ID)
		if err != nil {
			log.Printf("Error when deleting schema extension %q: %v\n", *schemaExtension.ID, err)
		}
	}
}
