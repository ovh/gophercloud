//go:build acceptance || vpngateway || ipsecpolicies

package v1

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/vpngateway/v1/ipsecpolicies"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestIPsecPoliciesList(t *testing.T) {
	client := NewClient(t)

	allPages, err := ipsecpolicies.List(client, nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allPolicies, err := ipsecpolicies.ExtractPolicies(allPages)
	th.AssertNoErr(t, err)

	for _, policy := range allPolicies {
		tools.PrintResource(t, policy)
	}
}

func TestIPsecPoliciesCRUD(t *testing.T) {
	client := NewClient(t)

	policy, err := CreateIPsecPolicy(t, client)
	th.AssertNoErr(t, err)
	defer DeleteIPsecPolicy(t, client, policy.ID)

	tools.PrintResource(t, policy)

	newPolicy, err := ipsecpolicies.Get(context.TODO(), client, policy.ID).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, *policy, *newPolicy)

	listOpts := ipsecpolicies.ListOpts{
		Name:    policy.Name,
		SortKey: "created_at",
		SortDir: "desc",
	}
	allPages, err := ipsecpolicies.List(client, listOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	allPolicies, err := ipsecpolicies.ExtractPolicies(allPages)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 1, len(allPolicies))
	th.AssertEquals(t, policy.ID, allPolicies[0].ID)

	// Leaving AES-GCM without naming an auth algorithm restores the default.
	description := tools.RandomString("TESTACC-desc-", 8)
	updateOpts := ipsecpolicies.UpdateOpts{
		Description:         &description,
		EncryptionAlgorithm: ipsecpolicies.EncryptionAlgorithmAES256,
		PFS:                 ipsecpolicies.PFSGroup14,
		Lifetime: &ipsecpolicies.LifetimeOpts{
			Units: ipsecpolicies.UnitSeconds,
			Value: 3600,
		},
	}

	updatedPolicy, err := ipsecpolicies.Update(context.TODO(), client, policy.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, updatedPolicy)

	th.AssertEquals(t, policy.Name, updatedPolicy.Name)
	th.AssertEquals(t, description, updatedPolicy.Description)
	th.AssertEquals(t, "aes-256", updatedPolicy.EncryptionAlgorithm)
	th.AssertEquals(t, "group14", updatedPolicy.PFS)
	th.AssertEquals(t, 3600, updatedPolicy.Lifetime.Value)
	if updatedPolicy.AuthAlgorithm == nil {
		t.Fatal("Expected the default auth algorithm after leaving AES-GCM")
	}
	th.AssertEquals(t, "sha256", *updatedPolicy.AuthAlgorithm)

	updateOpts = ipsecpolicies.UpdateOpts{
		AuthAlgorithm: ipsecpolicies.AuthAlgorithmSHA384,
	}

	updatedPolicy, err = ipsecpolicies.Update(context.TODO(), client, policy.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "sha384", *updatedPolicy.AuthAlgorithm)
}
