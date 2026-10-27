package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/vpngateway/v1/ipsecpolicies"
	fake "github.com/gophercloud/gophercloud/v2/openstack/vpngateway/v1/testhelper"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestCreate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandlePolicyCreationSuccessfully(t, fakeServer)

	options := ipsecpolicies.CreateOpts{
		ProjectID:           "9145d91459d248b1b02fdaca97c6a75d",
		Name:                "policy",
		Description:         "IPsec policy",
		TransformProtocol:   ipsecpolicies.TransformProtocolESP,
		EncryptionAlgorithm: ipsecpolicies.EncryptionAlgorithmAES256GCM,
		EncapsulationMode:   ipsecpolicies.EncapsulationModeTunnel,
		PFS:                 ipsecpolicies.PFSGroup20,
		Lifetime: &ipsecpolicies.LifetimeOpts{
			Units: ipsecpolicies.UnitSeconds,
			Value: 1800,
		},
	}

	actual, err := ipsecpolicies.Create(context.TODO(), fake.ServiceClient(fakeServer), options).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, Policy, *actual)
}

func TestGet(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandlePolicyGetSuccessfully(t, fakeServer)

	actual, err := ipsecpolicies.Get(context.TODO(), fake.ServiceClient(fakeServer), PolicyID).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, Policy, *actual)
}

func TestDelete(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandlePolicyDeletionSuccessfully(t, fakeServer)

	res := ipsecpolicies.Delete(context.TODO(), fake.ServiceClient(fakeServer), PolicyID)
	th.AssertNoErr(t, res.Err)
}

func TestList(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandlePolicyListSuccessfully(t, fakeServer)

	options := ipsecpolicies.ListOpts{
		ProjectID:   "9145d91459d248b1b02fdaca97c6a75d",
		Name:        "policy",
		Limit:       1,
		PageReverse: true,
		SortKey:     "created_at",
		SortDir:     "asc",
		Fields:      []string{"id", "name"},
	}

	var actual []ipsecpolicies.Policy
	pages := 0
	err := ipsecpolicies.List(fake.ServiceClient(fakeServer), options).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		pages++
		policies, err := ipsecpolicies.ExtractPolicies(page)
		if err != nil {
			return false, err
		}
		actual = append(actual, policies...)
		return true, nil
	})
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 2, pages)
	th.CheckDeepEquals(t, []ipsecpolicies.Policy{Policy, SecondPolicy}, actual)
}

func TestListOptsMarker(t *testing.T) {
	query, err := ipsecpolicies.ListOpts{Marker: PolicyID}.ToPolicyListQuery()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "?marker="+PolicyID, query)
}

func TestUpdate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandlePolicyUpdateSuccessfully(t, fakeServer)

	name := "updatedname"
	description := "updated policy"
	options := ipsecpolicies.UpdateOpts{
		Name:                &name,
		Description:         &description,
		TransformProtocol:   ipsecpolicies.TransformProtocolESP,
		AuthAlgorithm:       ipsecpolicies.AuthAlgorithmSHA256,
		EncryptionAlgorithm: ipsecpolicies.EncryptionAlgorithmAES256,
		EncapsulationMode:   ipsecpolicies.EncapsulationModeTunnel,
		PFS:                 ipsecpolicies.PFSGroup14,
		Lifetime: &ipsecpolicies.LifetimeOpts{
			Units: ipsecpolicies.UnitSeconds,
			Value: 3600,
		},
	}

	actual, err := ipsecpolicies.Update(context.TODO(), fake.ServiceClient(fakeServer), PolicyID, options).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, UpdatedPolicy, *actual)
}
