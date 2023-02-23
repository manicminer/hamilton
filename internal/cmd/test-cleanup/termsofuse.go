package main

import (
	"fmt"
	"log"

	"github.com/manicminer/hamilton/msgraph"
)

func cleanupTermsOfUseAgreements() {
	client := msgraph.NewTermsOfUseAgreementClient(tenantId)
	client.BaseClient.Authorizer = authorizer

	result, _, err := client.List(ctx, fmt.Sprintf("startsWith(displayName, '%s')", displayNamePrefix))
	if err != nil {
		log.Println(err)
		return
	}
	if result == nil {
		log.Println("bad API response, nil terms of use agreements result received")
		return
	}
	for _, item := range *result {
		if item.ID == nil || item.DisplayName == nil {
			log.Println("Terms of use agreement returned with nil ID or DisplayName")
			continue
		}

		log.Printf("Deleting terms of use agreement %q (DisplayName: %q)\n", *item.ID, *item.DisplayName)
		_, err := client.Delete(ctx, *item.ID)
		if err != nil {
			log.Printf("Error when deleting terms of use agreement %q: %v\n", *item.ID, err)
		}
	}
}
