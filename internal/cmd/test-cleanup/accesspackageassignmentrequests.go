package main

import (
	"log"

	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/manicminer/hamilton/msgraph"
)

func cleanupAccessPackageAssignmentRequests() {
	client := msgraph.NewAccessPackageAssignmentRequestClient(tenantId)
	client.BaseClient.Authorizer = authorizer

	result, _, err := client.List(ctx, odata.Query{})
	if err != nil {
		log.Println(err)
		return
	}
	if result == nil {
		log.Println("bad API response, nil access package assignment requests result received")
		return
	}
	for _, item := range *result {
		if item.ID == nil {
			log.Println("Access package assignment requests returned with nil ID")
			continue
		}

		log.Printf("Deleting access package assignment request %q\n", *item.ID)
		_, err := client.Delete(ctx, *item.ID)
		if err != nil {
			log.Printf("Error when deleting access package assignment request %q: %v\n", *item.ID, err)
		}
	}
}
