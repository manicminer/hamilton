package main

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/manicminer/hamilton/msgraph"
)

func cleanupGroups() {
	groupsClient := msgraph.NewGroupsClient(tenantId)
	groupsClient.BaseClient.Authorizer = authorizer

	groups, _, err := groupsClient.List(ctx, odata.Query{Filter: fmt.Sprintf("startsWith(displayName, '%s')", displayNamePrefix)})
	if err != nil {
		log.Println(err)
		return
	}
	if groups == nil {
		log.Println("bad API response, nil groups result received")
		return
	}
	for _, group := range *groups {
		if group.ID() == nil || group.DisplayName == nil {
			log.Println("Group returned with nil ID or DisplayName")
			continue
		}

		log.Printf("Deleting group %q (DisplayName: %q)\n", *group.ID(), *group.DisplayName)
		_, err := groupsClient.Delete(ctx, *group.ID())
		if err != nil {
			log.Printf("Error when deleting group %q: %v\n", *group.ID(), err)
		}

		if group.HasTypes([]msgraph.GroupType{msgraph.GroupTypeUnified}) {
			log.Printf("Permanently deleting group %q (DisplayName: %q)\n", *group.ID(), *group.DisplayName)
			_, err = groupsClient.DeletePermanently(ctx, *group.ID())
			if err != nil {
				log.Printf("Error when permanently deleting group %q: %v\n", *group.ID(), err)
			}
		}
	}
}
