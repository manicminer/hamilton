package msgraph

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/manicminer/hamilton/odata"
)

// SynchronizationJobClient performs operations on SynchronizationJobs.
type SynchronizationJobClient struct {
	BaseClient Client
}

// NewSynchronizationJobClient returns a new SynchronizationJobClient
func NewSynchronizationJobClient(tenantId string) *SynchronizationJobClient {
	return &SynchronizationJobClient{
		BaseClient: NewClient(VersionBeta, tenantId),
	}
}

// Copied from another function
// TODO REMOVE
// func (c *ServicePrincipalsClient) AddOwners(ctx context.Context, servicePrincipal *ServicePrincipal) (int, error) {
// 	var status int

// 	if servicePrincipal.ID == nil {
// 		return status, errors.New("cannot update service principal with nil ID")
// 	}
// 	if servicePrincipal.Owners == nil {
// 		return status, errors.New("cannot update service principal with nil Owners")
// 	}

// 	for _, owner := range *servicePrincipal.Owners {
// 		// don't fail if an owner already exists
// 		checkOwnerAlreadyExists := func(resp *http.Response, o *odata.OData) bool {
// 			if resp != nil && resp.StatusCode == http.StatusBadRequest && o != nil && o.Error != nil {
// 				return o.Error.Match(odata.ErrorAddedObjectReferencesAlreadyExist)
// 			}
// 			return false
// 		}

// 		body, err := json.Marshal(DirectoryObject{ODataId: owner.ODataId})
// 		if err != nil {
// 			return status, fmt.Errorf("json.Marshal(): %v", err)
// 		}

// 		_, status, _, err = c.BaseClient.Post(ctx, PostHttpRequestInput{
// 			Body:                   body,
// 			ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
// 			ValidStatusCodes:       []int{http.StatusNoContent},
// 			ValidStatusFunc:        checkOwnerAlreadyExists,
// 			Uri: Uri{
// 				Entity:      fmt.Sprintf("/servicePrincipals/%s/owners/$ref", *servicePrincipal.ID),
// 				HasTenantId: true,
// 			},
// 		})
// 		if err != nil {
// 			return status, fmt.Errorf("ServicePrincipalsClient.BaseClient.Post(): %v", err)
// 		}
// 	}

// 	return status, nil
// }

// TODO: CHeck if we can use OData
// List returns a list of SynchronizationJobs, optionally queried using OData.
func (c *SynchronizationJobClient) List(ctx context.Context, query odata.Query, servicePrincipal *ServicePrincipal) (*[]SynchronizationJob, int, error) {
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		DisablePaging:    query.Top > 0,
		OData:            query,
		ValidStatusCodes: []int{http.StatusOK},
		Uri: Uri{
			Entity:      fmt.Sprintf("/servicePrincipals/%s/synchronization/jobs/", *servicePrincipal.ID),
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("SynchronizationJobClient.BaseClient.Get(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var data struct {
		SynchronizationJobs []SynchronizationJob `json:"value"`
	}
	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &data.SynchronizationJobs, status, nil
}

// Get retrieves a SynchronizationJob
func (c *SynchronizationJobClient) Get(ctx context.Context, id string, query odata.Query, servicePrincipal *ServicePrincipal) (*SynchronizationJob, int, error) {
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		OData:            query,
		ValidStatusCodes: []int{http.StatusOK},
		Uri: Uri{
			Entity:      fmt.Sprintf("/servicePrincipals/%s/synchronization/jobs/%s", *servicePrincipal.ID, id),
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("SynchronizationJobClient.BaseClient.Get(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var SynchronizationJob SynchronizationJob
	if err := json.Unmarshal(respBody, &SynchronizationJob); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &SynchronizationJob, status, nil
}

// Creates a SynchronizationJob.
func (c *SynchronizationJobClient) Create(ctx context.Context, synchronizationJob SynchronizationJob, servicePrincipal *ServicePrincipal) (*SynchronizationJob, int, error) {
	var status int

	body, err := json.Marshal(synchronizationJob)
	if err != nil {
		return nil, status, fmt.Errorf("json.Marshal(): %v", err)
	}

	resp, status, _, err := c.BaseClient.Post(ctx, PostHttpRequestInput{
		Body:             body,
		ValidStatusCodes: []int{http.StatusCreated},
		Uri: Uri{
			Entity:      fmt.Sprintf("/servicePrincipals/%s/synchronization/jobs/", *servicePrincipal.ID),
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("SynchronizationJobClient.BaseClient.Post(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var newSynchronizationJob SynchronizationJob
	if err := json.Unmarshal(respBody, &newSynchronizationJob); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &newSynchronizationJob, status, nil
}

// Starts a SynchronizationJob.
func (c *SynchronizationJobClient) Start(ctx context.Context, id string, synchronizationJob SynchronizationJob, servicePrincipal *ServicePrincipal) (*SynchronizationJob, int, error) {
	var status int

	if synchronizationJob.ID == nil {
		return nil, status, errors.New("SynchronizationJobClient.Start(): cannot start SynchronizationJob with nil ID")
	}

	body, err := json.Marshal(synchronizationJob)
	if err != nil {
		return nil, status, fmt.Errorf("json.Marshal(): %v", err)
	}

	resp, status, _, err := c.BaseClient.Post(ctx, PostHttpRequestInput{
		Body:             body,
		ValidStatusCodes: []int{http.StatusCreated},
		Uri: Uri{
			Entity:      fmt.Sprintf("/servicePrincipals/%s/synchronization/jobs/%s/start", *servicePrincipal.ID, id),
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("SynchronizationJobClient.BaseClient.Post(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var newSynchronizationJob SynchronizationJob
	if err := json.Unmarshal(respBody, &newSynchronizationJob); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &newSynchronizationJob, status, nil
}

// Delete
func (c *SynchronizationJobClient) Delete(ctx context.Context, meta interface{}) {
	// TODO
	// client := meta.(*clients.Client).AdministrativeUnits.AdministrativeUnitsClient
	// administrativeUnitId := d.Id()

	// _, status, err := client.Get(ctx, administrativeUnitId, odata.Query{})
	// if err != nil {
	// 	if status == http.StatusNotFound {
	// 		return tf.ErrorDiagPathF(fmt.Errorf("Administrative unit was not found"), "id", "Retrieving administrative unit with object ID %q", administrativeUnitId)
	// 	}
	// 	return tf.ErrorDiagPathF(err, "id", "Retrieving administrative unit with object ID: %q", administrativeUnitId)
	// }

	// if _, err := client.Delete(ctx, administrativeUnitId); err != nil {
	// 	return tf.ErrorDiagF(err, "Deleting administrative unit with object ID: %q", administrativeUnitId)
	// }

	// // Wait for administrative unit object to be deleted
	// if err := helpers.WaitForDeletion(ctx, func(ctx context.Context) (*bool, error) {
	// 	client.BaseClient.DisableRetries = true
	// 	if _, status, err := client.Get(ctx, administrativeUnitId, odata.Query{}); err != nil {
	// 		if status == http.StatusNotFound {
	// 			return utils.Bool(false), nil
	// 		}
	// 		return nil, err
	// 	}
	// 	return utils.Bool(true), nil
	// }); err != nil {
	// 	return tf.ErrorDiagF(err, "Waiting for deletion of administrative unit with object ID %q", administrativeUnitId)
	// }

	// return nil
}

// Pause
func (c *SynchronizationJobClient) Pause(ctx context.Context, id string, synchronizationJob SynchronizationJob, servicePrincipal *ServicePrincipal) (*SynchronizationJob, int, error) {
	var status int

	if synchronizationJob.ID == nil {
		return nil, status, errors.New("SynchronizationJobClient.Pause(): cannot pause SynchronizationJob with nil ID")
	}

	body, err := json.Marshal(synchronizationJob)
	if err != nil {
		return nil, status, fmt.Errorf("json.Marshal(): %v", err)
	}

	resp, status, _, err := c.BaseClient.Post(ctx, PostHttpRequestInput{
		Body:             body,
		ValidStatusCodes: []int{http.StatusCreated},
		Uri: Uri{
			Entity:      fmt.Sprintf("/servicePrincipals/%s/synchronization/jobs/%s/pause", *servicePrincipal.ID, id),
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("SynchronizationJobClient.BaseClient.Post(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var newSynchronizationJob SynchronizationJob
	if err := json.Unmarshal(respBody, &newSynchronizationJob); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &newSynchronizationJob, status, nil
}

// Restart, TODO Criteria
func (c *SynchronizationJobClient) Restart(ctx context.Context, id string, synchronizationJob SynchronizationJob, servicePrincipal *ServicePrincipal) (*SynchronizationJob, int, error) {
	var status int

	if synchronizationJob.ID == nil {
		return nil, status, errors.New("SynchronizationJobClient.Restart(): cannot restart SynchronizationJob with nil ID")
	}

	body, err := json.Marshal(synchronizationJob)
	if err != nil {
		return nil, status, fmt.Errorf("json.Marshal(): %v", err)
	}

	resp, status, _, err := c.BaseClient.Post(ctx, PostHttpRequestInput{
		Body:             body,
		ValidStatusCodes: []int{http.StatusCreated},
		Uri: Uri{
			Entity:      fmt.Sprintf("/servicePrincipals/%s/synchronization/jobs/%s/restart", *servicePrincipal.ID, id),
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("SynchronizationJobClient.BaseClient.Post(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var newSynchronizationJob SynchronizationJob
	if err := json.Unmarshal(respBody, &newSynchronizationJob); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &newSynchronizationJob, status, nil
}

// Provision on demand
func (c *SynchronizationJobClient) ProvisionOnDemand(ctx context.Context, id string, synchronizationJob SynchronizationJob, servicePrincipal *ServicePrincipal) (*SynchronizationJob, int, error) {
	var status int

	if synchronizationJob.ID == nil {
		return nil, status, errors.New("SynchronizationJobClient.ProvisionOnDemand(): cannot set provision on demand for SynchronizationJob with nil ID")
	}

	body, err := json.Marshal(synchronizationJob)
	if err != nil {
		return nil, status, fmt.Errorf("json.Marshal(): %v", err)
	}

	resp, status, _, err := c.BaseClient.Post(ctx, PostHttpRequestInput{
		Body:             body,
		ValidStatusCodes: []int{http.StatusCreated},
		Uri: Uri{
			Entity:      fmt.Sprintf("/servicePrincipals/%s/synchronization/jobs/%s/provisionOnDemand", *servicePrincipal.ID, id),
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("SynchronizationJobClient.BaseClient.Post(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var newSynchronizationJob SynchronizationJob
	if err := json.Unmarshal(respBody, &newSynchronizationJob); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &newSynchronizationJob, status, nil
}

// Validate credentials, TODO params useSavedCredentials / credentials
func (c *SynchronizationJobClient) ValidateCredentials(ctx context.Context, id string, synchronizationJob SynchronizationJob, servicePrincipal *ServicePrincipal) (*SynchronizationJob, int, error) {
	var status int

	if synchronizationJob.ID == nil {
		return nil, status, errors.New("SynchronizationJobClient.ValidateCredentials(): cannot validate credentials SynchronizationJob with nil ID")
	}

	body, err := json.Marshal(synchronizationJob)
	if err != nil {
		return nil, status, fmt.Errorf("json.Marshal(): %v", err)
	}

	resp, status, _, err := c.BaseClient.Post(ctx, PostHttpRequestInput{
		Body:             body,
		ValidStatusCodes: []int{http.StatusCreated},
		Uri: Uri{
			Entity:      fmt.Sprintf("/servicePrincipals/%s/synchronization/jobs/%s/validateCredentials", *servicePrincipal.ID, id),
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("SynchronizationJobClient.BaseClient.Post(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var newSynchronizationJob SynchronizationJob
	if err := json.Unmarshal(respBody, &newSynchronizationJob); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &newSynchronizationJob, status, nil
}
