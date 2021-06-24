package main

import (
	"fmt"
	"log"

	"github.com/manicminer/hamilton/msgraph"
	"github.com/manicminer/hamilton/odata"
)

func cleanupApplications() {
	appsClient := msgraph.NewApplicationsClient(tenantId)
	appsClient.BaseClient.Authorizer = authorizer

	apps, _, err := appsClient.List(ctx, odata.Query{Filter: fmt.Sprintf("startsWith(displayName, '%s')", displayNamePrefix)})
	if err != nil {
		log.Println(err)
		return
	}
	if apps == nil {
		log.Println("bad API response, nil apps result received")
		return
	}
	for _, app := range *apps {
		if app.ID == nil || app.AppId == nil || app.DisplayName == nil {
			log.Println("App returned with nil ID, AppId or DisplayName")
			continue
		}

		log.Printf("Deleting application %q (AppID: %q, DisplayName: %q)\n", *app.ID, *app.AppId, *app.DisplayName)
		_, err := appsClient.Delete(ctx, *app.ID)
		if err != nil {
			log.Printf("Error when deleting application %q: %v\n", *app.ID, err)
		}

		log.Printf("Permanently deleting application %q (AppID: %q, DisplayName: %q)\n", *app.ID, *app.AppId, *app.DisplayName)
		_, err = appsClient.DeletePermanently(ctx, *app.ID)
		if err != nil {
			log.Printf("Error when permanently deleting application %q: %v\n", *app.ID, err)
		}
	}
}
