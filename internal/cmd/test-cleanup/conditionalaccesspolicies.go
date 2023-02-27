package main

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/manicminer/hamilton/msgraph"
)

func cleanupConditionalAccessPolicies() {
	conditionalAccessPoliciesClient := msgraph.NewConditionalAccessPoliciesClient()
	conditionalAccessPoliciesClient.BaseClient.Authorizer = authorizer

	conditionalAccessPolicies, _, err := conditionalAccessPoliciesClient.List(ctx, odata.Query{Filter: fmt.Sprintf("startsWith(displayName, '%s')", displayNamePrefix)})
	if err != nil {
		log.Println(err)
		return
	}
	if conditionalAccessPolicies == nil {
		log.Println("bad API response, nil ConditionalAccessPolicies result received")
		return
	}
	for _, conditionalAccessPolicies := range *conditionalAccessPolicies {
		if conditionalAccessPolicies.ID == nil || conditionalAccessPolicies.DisplayName == nil {
			log.Println("Conditional Access Policy returned with nil ID or DisplayName")
			continue
		}

		log.Printf("Deleting conditional access policy %q (DisplayName: %q)\n", *conditionalAccessPolicies.ID, *conditionalAccessPolicies.DisplayName)
		_, err := conditionalAccessPoliciesClient.Delete(ctx, *conditionalAccessPolicies.ID)
		if err != nil {
			log.Printf("Error when deleting conditional access policy %q: %v\n", *conditionalAccessPolicies.ID, err)
		}
	}
}
