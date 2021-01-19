package base

import "encoding/json"

type OData struct {
	Context      *string            `json:"@odata.context"`
	MetadataEtag *string            `json:"@odata.metadataEtag"`
	Type         *string            `json:"@odata.type"`
	Count        *string            `json:"@odata.count"`
	NextLink     *string            `json:"@odata.nextLink"`
	Delta        *string            `json:"@odata.delta"`
	DeltaLink    *string            `json:"@odata.deltaLink"`
	Id           *string            `json:"@odata.id"`
	Etag         *string            `json:"@odata.etag"`
	Value        *[]json.RawMessage `json:"value"`
}
