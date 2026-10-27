//go:build acceptance || vpngateway || ikepolicies

package v1

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/vpngateway/v1/ikepolicies"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestIKEPoliciesList(t *testing.T) {
	client := NewClient(t)

	allPages, err := ikepolicies.List(client, nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allPolicies, err := ikepolicies.ExtractPolicies(allPages)
	th.AssertNoErr(t, err)

	for _, policy := range allPolicies {
		tools.PrintResource(t, policy)
	}
}

func TestIKEPoliciesCRUD(t *testing.T) {
	client := NewClient(t)

	policy, err := CreateIKEPolicy(t, client)
	th.AssertNoErr(t, err)
	defer DeleteIKEPolicy(t, client, policy.ID)

	tools.PrintResource(t, policy)

	newPolicy, err := ikepolicies.Get(context.TODO(), client, policy.ID).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, *policy, *newPolicy)

	listOpts := ikepolicies.ListOpts{
		Name:   policy.Name,
		Fields: []string{"id", "name"},
	}
	allPages, err := ikepolicies.List(client, listOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	allPolicies, err := ikepolicies.ExtractPolicies(allPages)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 1, len(allPolicies))
	th.AssertEquals(t, policy.ID, allPolicies[0].ID)
	th.AssertEquals(t, "", allPolicies[0].PFS)

	name := ""
	description := tools.RandomString("TESTACC-desc-", 8)
	updateOpts := ikepolicies.UpdateOpts{
		Name:                &name,
		Description:         &description,
		AuthAlgorithm:       ikepolicies.AuthAlgorithmSHA512,
		EncryptionAlgorithm: ikepolicies.EncryptionAlgorithmAES128,
		PFS:                 ikepolicies.PFSGroup31,
		Lifetime: &ikepolicies.LifetimeOpts{
			Units: ikepolicies.UnitSeconds,
			Value: 3600,
		},
	}

	updatedPolicy, err := ikepolicies.Update(context.TODO(), client, policy.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, updatedPolicy)

	th.AssertEquals(t, name, updatedPolicy.Name)
	th.AssertEquals(t, description, updatedPolicy.Description)
	th.AssertEquals(t, "sha512", updatedPolicy.AuthAlgorithm)
	th.AssertEquals(t, "aes-128", updatedPolicy.EncryptionAlgorithm)
	th.AssertEquals(t, "group31", updatedPolicy.PFS)
	th.AssertEquals(t, "seconds", updatedPolicy.Lifetime.Units)
	th.AssertEquals(t, 3600, updatedPolicy.Lifetime.Value)
}
