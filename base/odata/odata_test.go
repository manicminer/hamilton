package odata_test

import (
	"encoding/json"
	"testing"

	"github.com/manicminer/hamilton/base/odata"
)

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
			t.Errorf("test case %d: expected %q, got %q", n, s, c.expected)
		}
	}
}
