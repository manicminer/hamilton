package odata_test

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/manicminer/hamilton/odata"
)

func TestQueryHeaders(t *testing.T) {
	type testCase struct {
		query    odata.Query
		expected http.Header
	}
	testCases := []testCase{
		{
			query: odata.Query{},
			expected: http.Header{
				"Accept":           []string{"application/json; charset=utf-8; IEEE754Compatible=false"},
				"Odata-Maxversion": []string{odata.ODataVersion},
				"Odata-Version":    []string{odata.ODataVersion},
			},
		},
		{
			query: odata.Query{
				ConsistencyLevel: odata.ConsistencyLevelEventual,
				Metadata:         odata.MetadataMinimal,
			},
			expected: http.Header{
				"Accept":           []string{"application/json; charset=utf-8; IEEE754Compatible=false; odata.metadata=minimal"},
				"Consistencylevel": []string{"eventual"},
				"Odata-Maxversion": []string{odata.ODataVersion},
				"Odata-Version":    []string{odata.ODataVersion},
			},
		},
	}
	for n, c := range testCases {
		v := c.query.Headers()
		if !reflect.DeepEqual(v, c.expected) {
			t.Errorf("test case %d for %s: expected %#v, got %#v", n, "Headers()", c.expected, v)
		}
		c.expected.Set("Testing", "foo")
		v = c.query.AppendHeaders(http.Header{"Testing": []string{"foo"}})
		if !reflect.DeepEqual(v, c.expected) {
			t.Errorf("test case %d for %s: expected %#v, got %#v", n, "AppendHeaders()", c.expected, v)
		}
	}
}

func TestQueryValues(t *testing.T) {
	type testCase struct {
		query    odata.Query
		expected url.Values
	}
	testCases := []testCase{
		{
			query:    odata.Query{},
			expected: url.Values{},
		},
		{
			query: odata.Query{
				Count:  true,
				Format: odata.FormatAtom,
				Skip:   20,
				Top:    10,
			},
			expected: url.Values{
				"$count":  []string{"true"},
				"$format": []string{"atom"},
				"$skip":   []string{"20"},
				"$top":    []string{"10"},
			},
		},
		{
			query: odata.Query{
				OrderBy: odata.OrderBy{
					Field:     "displayName",
					Direction: "desc",
				},
			},
			expected: url.Values{
				"$orderby": []string{"displayName desc"},
			},
		},
		{
			query: odata.Query{
				Expand: odata.Expand{
					Relationship: "children",
					Select:       []string{"id", "childName"},
				},
			},
			expected: url.Values{
				"$expand": []string{"children($select=id,childName)"},
			},
		},
		{
			query: odata.Query{
				Filter: "startsWith(displayName,'Widgets')",
			},
			expected: url.Values{
				"$filter": []string{"startsWith(displayName,'Widgets')"},
			},
		},
		{
			query: odata.Query{
				Search: "displayName:Astley",
			},
			expected: url.Values{
				"$search": []string{`"displayName:Astley"`},
			},
		},
		{
			query: odata.Query{
				Select: []string{"id", "userPrincipalName"},
			},
			expected: url.Values{
				"$select": []string{"id,userPrincipalName"},
			},
		},
	}
	for n, c := range testCases {
		v := c.query.Values()
		if !reflect.DeepEqual(v, c.expected) {
			t.Errorf("test case %d: expected %#v, got %#v", n, c.expected, v)
		}
	}
}
