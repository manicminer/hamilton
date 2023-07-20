package msgraph_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/manicminer/hamilton/internal/test"
	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/msgraph"
)

func TestAuthenticationStrengthPolicyClient(t *testing.T) {
	c := test.NewTest(t)
	defer c.CancelFunc()

	policy := testAuthenticationStrengthPoliciesClient_Create(t, c, msgraph.AuthenticationStrengthPolicy{
		DisplayName:         utils.StringPtr(fmt.Sprintf("test-policy-%s", c.RandomString)),
		Description:         utils.StringPtr("FIDO2"),
		AllowedCombinations: &[]string{"password", "hardwareOath"},
	},
	)

	updatePolicy := msgraph.AuthenticationStrengthPolicy{
		ID:          policy.ID,
		DisplayName: utils.StringPtr(fmt.Sprintf("test-policy-updated-%s", c.RandomString)),
	}
	testAuthenticationStrengthPolicysClient_Update(t, c, updatePolicy)

	testAuthenticationStrengthPolicysClient_List(t, c)
	testAuthenticationStrengthPolicysClient_Get(t, c, *policy.ID)
	testAuthenticationStrengthPolicysClient_Delete(t, c, *policy.ID)

}

func testAuthenticationStrengthPoliciesClient_Create(t *testing.T, c *test.Test, a msgraph.AuthenticationStrengthPolicy) (authenticationStrengthPolicy *msgraph.AuthenticationStrengthPolicy) {
	authenticationStrengthPolicy, status, err := c.AuthenticationStrengthPoliciesClient.Create(c.Context, a)
	if err != nil {
		t.Fatalf("AuthenticationStrengthPolicyClient.Create(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AuthenticationStrengthPolicyClient.Create(): invalid status: %d", status)
	}
	if authenticationStrengthPolicy == nil {
		t.Fatal("AuthenticationStrengthPolicyClient.Create(): authenticationStrengthPolicy was nil")
	}
	if authenticationStrengthPolicy.ID == nil {
		t.Fatal("AuthenticationStrengthPolicyClient.Create(): authenticationStrengthPolicy.ID was nil")
	}
	return
}

func testAuthenticationStrengthPolicysClient_Get(t *testing.T, c *test.Test, id string) (policy *msgraph.AuthenticationStrengthPolicy) {
	policy, status, err := c.AuthenticationStrengthPoliciesClient.Get(c.Context, id, odata.Query{})
	if err != nil {
		t.Fatalf("AuthenticationStrengthPolicyClient.Get(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AuthenticationStrengthPolicyClient.Get(): invalid status: %d", status)
	}
	if policy == nil {
		t.Fatal("AuthenticationStrengthPolicyClient.Get(): policy was nil")
	}
	return
}

func testAuthenticationStrengthPolicysClient_Update(t *testing.T, c *test.Test, policy msgraph.AuthenticationStrengthPolicy) {
	status, err := c.AuthenticationStrengthPoliciesClient.Update(c.Context, policy)
	if err != nil {
		t.Fatalf("AuthenticationStrengthPolicyClient.Update(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AuthenticationStrengthPolicyClient.Update(): invalid status: %d", status)
	}
}

func testAuthenticationStrengthPolicysClient_List(t *testing.T, c *test.Test) (policies *[]msgraph.AuthenticationStrengthPolicy) {
	policies, _, err := c.AuthenticationStrengthPoliciesClient.List(c.Context, odata.Query{Top: 10})
	if err != nil {
		t.Fatalf("AuthenticationStrengthPolicyClient.List(): %v", err)
	}
	if policies == nil {
		t.Fatal("AuthenticationStrengthPolicyClient.List(): policies was nil")
	}
	return
}

func testAuthenticationStrengthPolicysClient_Delete(t *testing.T, c *test.Test, id string) {
	status, err := c.AuthenticationStrengthPoliciesClient.Delete(c.Context, id)
	if err != nil {
		t.Fatalf("AuthenticationStrengthPolicyClient.Delete(): %v", err)
	}
	if status < 200 || status >= 300 {
		t.Fatalf("AuthenticationStrengthPolicyClient.Delete(): invalid status: %d", status)
	}
}
