// TODO
package msgraph_test

import (
	"fmt"
	"testing"
	// "time"

	// "github.com/hashicorp/go-uuid"

	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/msgraph"
	// "github.com/manicminer/hamilton/odata"
)

func TestSynchronizationClient(t *testing.T) {
	c := test.NewTest(t)
	defer c.CancelFunc()

	app := testApplicationsClient_Create(t, c, msgraph.Application{
		DisplayName: utils.StringPtr(fmt.Sprintf("test-serviceprincipal-%s", c.RandomString)),
	})

	sp := testServicePrincipalsClient_Create(t, c, msgraph.ServicePrincipal{
		AccountEnabled: utils.BoolPtr(true),
		AppId:          app.AppId,
		DisplayName:    app.DisplayName,
	})


	testServicePrincipalsClient_Delete(t, c, *sp.ID)

	testApplicationsClient_Delete(t, c, *app.ID)
}

