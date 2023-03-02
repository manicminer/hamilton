package msgraph_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/msgraph"
)

func TestTokenPolicyClient(t *testing.T) {
	c := test.NewTest(t)
	defer c.CancelFunc()
	policy := testTokenIssuancePolicyClient_Create(t, c, msgraph.TokenIssuancePolicy{
		DisplayName: utils.StringPtr(fmt.Sprintf("test-token-issuance-policy-%s", c.RandomString)),
		Definition: utils.ArrayStringPtr(
			[]string{
				"{\"TokenIssuancePolicy\":{\"Version\":1,\"SigningAlgorithm\":\"http://www.w3.org/2001/04/xmldsig-more#rsa-sha256\",\"TokenResponseSigningPolicy\":\"ResponseAndToken\",\"SamlTokenVersion\":\"2.0\",\"EmitSamlNameFormat\":false}}",
			},
		),
	})
	testTokenIssuancePolicyClient_List(t, c)
	testTokenIssuancePolicyClient_Get(t, c, *policy.ID())
	policy.DisplayName = utils.StringPtr(fmt.Sprintf("test-token-issuance-policy-%s", c.RandomString))
	testTokenIssuancePolicyClient_Update(t, c, *policy)
	testTokenIssuancePolicyClient_Delete(t, c, *policy.ID())
}

func testTokenIssuancePolicyClient_Create(t *testing.T, c *test.Test, p msgraph.TokenIssuancePolicy) (policy *msgraph.TokenIssuancePolicy) {
	policy, status, err := c.TokenIssuancePolicyClient.Create(c.Context, p)
	if err != nil {
		t.Fatalf("TokenIssuancePolicyClient.Create(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("TokenIssuancePolicyClient.Create(): invalid status: %d", status)
	}
	if policy == nil {
		t.Fatal("TokenIssuancePolicyClient.Create(): policy was nil")
	}
	if policy.ID() == nil {
		t.Fatal("TokenIssuancePolicyClient.Create(): policy.ID was nil")
	}
	return
}

func testTokenIssuancePolicyClient_List(t *testing.T, c *test.Test) (policy *[]msgraph.ClaimsMappingPolicy) {
	policies, _, err := c.TokenIssuancePolicyClient.List(c.Context, odata.Query{Top: 10})
	if err != nil {
		t.Fatalf("TokenIssuancePolicyClient.List(): %v", err)
	}
	if policies == nil {
		t.Fatal("TokenIssuancePolicyClient.List(): TokenIssuancePolicies was nil")
	}
	return
}

func testTokenIssuancePolicyClient_Get(t *testing.T, c *test.Test, id string) (policy *msgraph.TokenIssuancePolicy) {
	policies, status, err := c.TokenIssuancePolicyClient.Get(c.Context, id, odata.Query{})
	if err != nil {
		t.Fatalf("TokenIssuancePolicyClient.Get(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("TokenIssuancePolicyClient.Get(): invalid status: %d", status)
	}
	if policies == nil {
		t.Fatal("TokenIssuancePolicyClient.Get(): policies was nil")
	}
	return
}

func testTokenIssuancePolicyClient_Update(t *testing.T, c *test.Test, p msgraph.TokenIssuancePolicy) {
	status, err := c.TokenIssuancePolicyClient.Update(c.Context, p)
	if err != nil {
		t.Fatalf("TokenIssuancePolicyClient.Update(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("TokenIssuancePolicyClient.Update(): invalid status: %d", status)
	}
}

func testTokenIssuancePolicyClient_Delete(t *testing.T, c *test.Test, id string) {
	status, err := c.TokenIssuancePolicyClient.Delete(c.Context, id)
	if err != nil {
		t.Fatalf("TokenIssuancePolicyClient.Delete(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("TokenIssuancePolicyClient.Delete(): invalid status: %d", status)
	}
}
