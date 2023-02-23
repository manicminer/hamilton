package main

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/manicminer/hamilton/msgraph"
)

func cleanupUsers() {
	usersClient := msgraph.NewUsersClient(tenantId)
	usersClient.BaseClient.Authorizer = authorizer

	users, _, err := usersClient.List(ctx, odata.Query{Filter: fmt.Sprintf("startsWith(displayName, '%s')", displayNamePrefix)})
	if err != nil {
		log.Println(err)
		return
	}
	if users == nil {
		log.Println("bad API response, nil users result received")
		return
	}
	for _, user := range *users {
		if user.ID() == nil || user.DisplayName == nil {
			log.Println("User returned with nil ID or DisplayName")
			continue
		}

		log.Printf("Deleting user %q (DisplayName: %q)\n", *user.ID(), *user.DisplayName)
		_, err := usersClient.Delete(ctx, *user.ID())
		if err != nil {
			log.Printf("Error when deleting user %q: %v\n", *user.ID(), err)
		}

		log.Printf("Permanently deleting user %q (DisplayName: %q)\n", *user.ID(), *user.DisplayName)
		_, err = usersClient.DeletePermanently(ctx, *user.ID())
		if err != nil {
			log.Printf("Error when permanently deleting user %q: %v\n", *user.ID(), err)
		}
	}
}
