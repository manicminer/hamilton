package msgraph_test

import (
	"fmt"
	"testing"

	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/msgraph"
	"github.com/manicminer/hamilton/odata"
)

func TestClaimsMappingPolicyClient(t *testing.T) {
	c := test.NewTest(t)
	defer c.CancelFunc()
	policy := testClaimsMappingPolicyClient_Create(t, c, msgraph.ClaimsMappingPolicy{
		DisplayName: utils.StringPtr(fmt.Sprintf("test-claims-mapping-policy-%s", c.RandomString)),
		Definition: utils.ArrayStringPtr(
			[]string{
				"{\"ClaimsMappingPolicy\":{\"Version\":1,\"IncludeBasicClaimSet\":\"true\",\"ClaimsSchema\": [{\"Source\":\"user\",\"ID\":\"employeeid\",\"SamlClaimType\":\"http://schemas.xmlsoap.org/ws/2005/05/identity/claims/name\",\"JwtClaimType\":\"name\"},{\"Source\":\"company\",\"ID\":\"tenantcountry\",\"SamlClaimType\":\"http://schemas.xmlsoap.org/ws/2005/05/identity/claims/country\",\"JwtClaimType\":\"country\"}]}}",
			},
		),
	})
	testClaimsMappingPolicyClient_List(t, c)
	testClaimsMappingPolicyClient_Get(t, c, *policy.ID())
	policy.DisplayName = utils.StringPtr(fmt.Sprintf("test-claims-mapping-policy-%s", c.RandomString))
	testClaimsMappingPolicyClient_Update(t, c, *policy)
	testClaimsMappingPolicyClient_Delete(t, c, *policy.ID())
}

func testClaimsMappingPolicyClient_Create(t *testing.T, c *test.Test, p msgraph.ClaimsMappingPolicy) (policy *msgraph.ClaimsMappingPolicy) {
	policy, status, err := c.ClaimsMappingPolicyClient.Create(c.Context, p)
	if err != nil {
		t.Fatalf("ClaimsMappingPolicyClient.Create(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ClaimsMappingPolicyClient.Create(): invalid status: %d", status)
	}
	if policy == nil {
		t.Fatal("ClaimsMappingPolicyClient.Create(): policy was nil")
	}
	if policy.ID() == nil {
		t.Fatal("ClaimsMappingPolicyClient.Create(): policy.ID was nil")
	}
	return
}

func testClaimsMappingPolicyClient_List(t *testing.T, c *test.Test) (policy *[]msgraph.ClaimsMappingPolicy) {
	policies, _, err := c.ClaimsMappingPolicyClient.List(c.Context, odata.Query{Top: 10})
	if err != nil {
		t.Fatalf("ClaimsMappingPolicy.List(): %v", err)
	}
	if policies == nil {
		t.Fatal("ClaimsMappingPolicy.List(): ClaimsMappingPolicies was nil")
	}
	return
}

func testClaimsMappingPolicyClient_Get(t *testing.T, c *test.Test, id string) (ClaimsMappingPolicy *msgraph.ClaimsMappingPolicy) {
	policies, status, err := c.ClaimsMappingPolicyClient.Get(c.Context, id, odata.Query{})
	if err != nil {
		t.Fatalf("ClaimsMappingPolicyClient.Get(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ClaimsMappingPolicyClient.Get(): invalid status: %d", status)
	}
	if policies == nil {
		t.Fatal("ClaimsMappingPolicyClient.Get(): policies was nil")
	}
	return
}

func testClaimsMappingPolicyClient_Update(t *testing.T, c *test.Test, g msgraph.ClaimsMappingPolicy) {
	status, err := c.ClaimsMappingPolicyClient.Update(c.Context, g)
	if err != nil {
		t.Fatalf("ClaimsMappingPolicyClient.Update(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ClaimsMappingPolicyClient.Update(): invalid status: %d", status)
	}
}

func testClaimsMappingPolicyClient_Delete(t *testing.T, c *test.Test, id string) {
	status, err := c.ClaimsMappingPolicyClient.Delete(c.Context, id)
	if err != nil {
		t.Fatalf("ClaimsMappingPolicyClient.Delete(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("ClaimsMappingPolicyClient.Delete(): invalid status: %d", status)
	}
}
