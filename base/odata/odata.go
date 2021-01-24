package odata

import (
	"encoding/json"
	"fmt"
)

type OData struct {
	Context      *string `json:"@odata.context"`
	MetadataEtag *string `json:"@odata.metadataEtag"`
	Type         *string `json:"@odata.type"`
	Count        *string `json:"@odata.count"`
	NextLink     *string `json:"@odata.nextLink"`
	Delta        *string `json:"@odata.delta"`
	DeltaLink    *string `json:"@odata.deltaLink"`
	Id           *string `json:"@odata.id"`
	Etag         *string `json:"@odata.etag"`

	Error      *Error `json:"-"`

	Value *[]json.RawMessage `json:"value"`
}

func (o *OData) UnmarshalJSON(data []byte) error {
	// Perform unmarshalling using a local type
	type odata OData
	var o2 odata
	if err := json.Unmarshal(data, &o2); err != nil {
		return err
	}
	*o = OData(o2)

	// Look for errors in the "error" and "odata.error" fields
	var e map[string]json.RawMessage
	if err := json.Unmarshal(data, &e); err != nil {
		return err
	}
	for _, k := range []string{"error", "odata.error"} {
		if v, ok := e[k]; ok {
			var e2 Error
			if err := json.Unmarshal(v, &e2); err != nil {
				return err
			}
			o.Error = &e2
		}
	}
	return nil
}

// Error is used to unmarshal an error response from Microsoft Graph.
type Error struct {
	Code            *string      `json:"code"`
	Date            *string      `json:"date"`
	Message         *interface{} `json:"message"` // sometimes a string, sometimes an object :/
	ClientRequestId *string      `json:"client-request-id"`
	RequestId       *string      `json:"request-id"`
	InnerError      *Error       `json:"innerError"` // nested errors
}

func (e Error) String() (s string) {
	if e.Code != nil {
		s = *e.Code
	}
	if e.Message != nil {
		msg := *e.Message
		if v, ok := msg.(string); ok {
			if v != "" {
				s = fmt.Sprintf("%s: %s", s, v)
			}
		} else if m, ok := msg.(map[string]interface{}); ok {
			if v, ok := m["value"]; ok {
				if vs, ok := v.(string); ok && vs != "" {
					s = fmt.Sprintf("%s: %s", s, vs)
				}
			}
		}
	}
	return
}