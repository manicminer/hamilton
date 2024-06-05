package main

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/manicminer/hamilton/msgraph"
)

func cleanupAccessPackageAssignmentPolicies() {
	client := msgraph.NewAccessPackageAssignmentPolicyClient()
	client.BaseClient.Authorizer = authorizer

	result, _, err := client.List(ctx, odata.Query{Filter: fmt.Sprintf("startsWith(displayName, '%s')", displayNamePrefix)})
	if err != nil {
		log.Println(err)
		return
	}
	if result == nil {
		log.Println("bad API response, nil access package assignment policies result received")
		return
	}
	for _, item := range *result {
		if item.ID == nil || item.DisplayName == nil {
			log.Println("Access package assignment policy returned with nil ID or DisplayName")
			continue
		}

		log.Printf("Deleting access package assignment policy %q (DisplayName: %q)\n", *item.ID, *item.DisplayName)
		_, err := client.Delete(ctx, *item.ID)
		if err != nil {
			log.Printf("Error when deleting access package assignment policy %q: %v\n", *item.ID, err)
		}
	}
}
