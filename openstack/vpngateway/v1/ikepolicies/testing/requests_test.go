package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/vpngateway/v1/ikepolicies"
	fake "github.com/gophercloud/gophercloud/v2/openstack/vpngateway/v1/testhelper"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestCreate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandlePolicyCreationSuccessfully(t, fakeServer)

	options := ikepolicies.CreateOpts{
		ProjectID:             "9145d91459d248b1b02fdaca97c6a75d",
		Name:                  "policy",
		Description:           "IKE policy",
		AuthAlgorithm:         ikepolicies.AuthAlgorithmSHA384,
		EncryptionAlgorithm:   ikepolicies.EncryptionAlgorithmAES256GCM,
		Phase1NegotiationMode: ikepolicies.Phase1NegotiationModeMain,
		IKEVersion:            ikepolicies.IKEVersionv2,
		PFS:                   ikepolicies.PFSGroup19,
		Lifetime: &ikepolicies.LifetimeOpts{
			Units: ikepolicies.UnitSeconds,
			Value: 7200,
		},
	}

	actual, err := ikepolicies.Create(context.TODO(), fake.ServiceClient(fakeServer), options).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, Policy, *actual)
}

func TestGet(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandlePolicyGetSuccessfully(t, fakeServer)

	actual, err := ikepolicies.Get(context.TODO(), fake.ServiceClient(fakeServer), PolicyID).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, Policy, *actual)
}

func TestDelete(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandlePolicyDeletionSuccessfully(t, fakeServer)

	res := ikepolicies.Delete(context.TODO(), fake.ServiceClient(fakeServer), PolicyID)
	th.AssertNoErr(t, res.Err)
}

func TestList(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandlePolicyListSuccessfully(t, fakeServer)

	options := ikepolicies.ListOpts{
		ProjectID:   "9145d91459d248b1b02fdaca97c6a75d",
		Name:        "policy",
		Limit:       1,
		PageReverse: true,
		SortKey:     "created_at",
		SortDir:     "asc",
		Fields:      []string{"id", "name"},
	}

	var actual []ikepolicies.Policy
	pages := 0
	err := ikepolicies.List(fake.ServiceClient(fakeServer), options).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		pages++
		policies, err := ikepolicies.ExtractPolicies(page)
		if err != nil {
			return false, err
		}
		actual = append(actual, policies...)
		return true, nil
	})
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 2, pages)
	th.CheckDeepEquals(t, []ikepolicies.Policy{Policy, SecondPolicy}, actual)
}

func TestListOptsMarker(t *testing.T) {
	query, err := ikepolicies.ListOpts{Marker: PolicyID}.ToPolicyListQuery()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "?marker="+PolicyID, query)
}

func TestUpdate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandlePolicyUpdateSuccessfully(t, fakeServer)

	name := "updatedname"
	description := "updated policy"
	options := ikepolicies.UpdateOpts{
		Name:                  &name,
		Description:           &description,
		AuthAlgorithm:         ikepolicies.AuthAlgorithmSHA512,
		EncryptionAlgorithm:   ikepolicies.EncryptionAlgorithmAES128,
		Phase1NegotiationMode: ikepolicies.Phase1NegotiationModeMain,
		IKEVersion:            ikepolicies.IKEVersionv2,
		PFS:                   ikepolicies.PFSGroup31,
		Lifetime: &ikepolicies.LifetimeOpts{
			Units: ikepolicies.UnitSeconds,
			Value: 3600,
		},
	}

	actual, err := ikepolicies.Update(context.TODO(), fake.ServiceClient(fakeServer), PolicyID, options).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, UpdatedPolicy, *actual)
}
