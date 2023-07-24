package main

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/manicminer/hamilton/msgraph"
)

func cleanupAuthenticationStrengthPolicies() {
	client := msgraph.NewAuthenticationStrengthPoliciesClient()
	client.BaseClient.Authorizer = authorizer

	authStrengthPolicies, _, err := client.List(ctx, odata.Query{Filter: fmt.Sprintf("startsWith(displayName, '%s')", displayNamePrefix)})
	if err != nil {
		log.Println(err)
		return
	}
	if authStrengthPolicies == nil {
		log.Println("bad API response, nil authStrengthPolicies result received")
		return
	}
	for _, policy := range *authStrengthPolicies {
		if policy.ID == nil || policy.DisplayName == nil {
			log.Println("Authentication Strength Policy returned with nil ID or DisplayName")
			continue
		}

		log.Printf("Deleting Authentication Strength Policy %q (DisplayName: %q)\n", *policy.ID, *policy.DisplayName)
		_, err := client.Delete(ctx, *policy.ID)
		if err != nil {
			log.Printf("Error when deleting Authentication Strength Policy %q: %v\n", *policy.ID, err)
		}
	}
}
