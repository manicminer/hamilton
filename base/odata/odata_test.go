package odata_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/manicminer/hamilton/base/odata"
	"github.com/manicminer/hamilton/internal/utils"
)

func TestOData(t *testing.T) {
	type testCase struct {
		response string
		expected odata.OData
	}
	testCases := []testCase{
		{
			response: `{
  "@odata.context": "https://graph.microsoft.com/beta/$metadata#servicePrincipals",
  "@odata.nextLink": "https://graph.microsoft.com/beta/26e25406-6564-4a26-98ee-c71ba03235ad/servicePrincipals?$skiptoken=X%274453707402000100000035536572766963655072696E636970616C5F31326430653134382D663634382D343233382D383566312D34336331643937353963313035536572766963655072696E636970616C5F31326430653134382D663634382D343233382D383566312D3433633164393735396331300000000000000000000000%27",
  "value": [
    {
      "id": "00000000-0000-0000-0000-000000000000",
      "deletedDateTime": null,
      "accountEnabled": true,
      "createdDateTime": "2020-07-08T01:22:29Z"
    }
  ]
}`,
			expected: odata.OData{
				Context:  utils.StringPtr("https://graph.microsoft.com/beta/$metadata#servicePrincipals"),
				NextLink: utils.StringPtr("https://graph.microsoft.com/beta/26e25406-6564-4a26-98ee-c71ba03235ad/servicePrincipals?$skiptoken=X%274453707402000100000035536572766963655072696E636970616C5F31326430653134382D663634382D343233382D383566312D34336331643937353963313035536572766963655072696E636970616C5F31326430653134382D663634382D343233382D383566312D3433633164393735396331300000000000000000000000%27"),
				Value: &[]json.RawMessage{[]byte(`{
      "id": "00000000-0000-0000-0000-000000000000",
      "deletedDateTime": null,
      "accountEnabled": true,
      "createdDateTime": "2020-07-08T01:22:29Z"
    }`)},
			},
		},
	}
	for n, c := range testCases {
		var o odata.OData
		err := json.Unmarshal([]byte(c.response), &o)
		if err != nil {
			t.Errorf("test case %d: JSON unmarshalling failed: %v", n, err)
			continue
		}
		if !reflect.DeepEqual(o, c.expected) {
			t.Errorf("test case %d: expected %#v, got %#v", n, c.expected, o)
		}
	}
}

func TestError(t *testing.T) {
	type testCase struct {
		response string
		expected string
	}
	testCases := []testCase{
		{
			response: `{
  "error": {
    "code": "Service_ServiceUnavailable",
    "message": "Service is temporarily unavailable. Please wait and retry again.",
    "innerError": {
      "date": "2021-01-24T14:52:27",
      "request-id": "7c974c85-e572-43ff-9633-f2dddf28725a",
      "client-request-id": "7c974c85-e572-43ff-9633-f2dddf28725a"
    }
  }
}`,
			expected: "Service_ServiceUnavailable: Service is temporarily unavailable. Please wait and retry again.",
		},
		{
			response: `{
  "odata.error": {
    "code": "Request_InvalidDataContractVersion",
    "message": {
      "lang": "en",
      "value": "The specified api-version is invalid. The value must exactly match a supported version."
    },
    "requestId": "e3a05e86-92ae-4e7e-9635-a3f62342da5b",
    "date": "2021-01-24T15:37:05"
  }
}`,
			expected: "Request_InvalidDataContractVersion: The specified api-version is invalid. The value must exactly match a supported version.",
		},
	}

	for n, c := range testCases {
		var o odata.OData
		err := json.Unmarshal([]byte(c.response), &o)
		if err != nil {
			t.Errorf("test case %d: JSON unmarshalling failed: %v", n, err)
			continue
		}
		if o.Error == nil {
			t.Errorf("test case %d: Error field was nil", n)
			continue
		}
		if s := o.Error.String(); s != c.expected {
			t.Errorf("test case %d: expected %q, got %q", n, c.expected, s)
		}
	}
}
