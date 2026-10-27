package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/vpngateway/v1/endpointgroups"
	fake "github.com/gophercloud/gophercloud/v2/openstack/vpngateway/v1/testhelper"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestCreate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGroupCreationSuccessfully(t, fakeServer)

	options := endpointgroups.CreateOpts{
		ProjectID:   "9145d91459d248b1b02fdaca97c6a75d",
		Name:        "peers",
		Description: "peer networks",
		Type:        endpointgroups.TypeCIDR,
		CIDRs:       []string{"10.2.0.0/24", "fd00:2::/64"},
	}

	actual, err := endpointgroups.Create(context.TODO(), fake.ServiceClient(fakeServer), options).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, Group, *actual)
}

func TestCreateRequiresTypeAndCIDRs(t *testing.T) {
	_, err := endpointgroups.CreateOpts{CIDRs: []string{"10.2.0.0/24"}}.ToEndpointGroupCreateMap()
	if err == nil {
		t.Fatal("expected an error without a type")
	}
	_, err = endpointgroups.CreateOpts{Type: endpointgroups.TypeCIDR}.ToEndpointGroupCreateMap()
	if err == nil {
		t.Fatal("expected an error without cidrs")
	}
}

func TestGet(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGroupGetSuccessfully(t, fakeServer)

	actual, err := endpointgroups.Get(context.TODO(), fake.ServiceClient(fakeServer), GroupID).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, Group, *actual)
}

func TestDelete(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGroupDeletionSuccessfully(t, fakeServer)

	res := endpointgroups.Delete(context.TODO(), fake.ServiceClient(fakeServer), GroupID)
	th.AssertNoErr(t, res.Err)
}

func TestList(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGroupListSuccessfully(t, fakeServer)

	options := endpointgroups.ListOpts{
		ProjectID:   "9145d91459d248b1b02fdaca97c6a75d",
		Name:        "peers",
		Limit:       1,
		PageReverse: true,
		SortKey:     "created_at",
		SortDir:     "asc",
		Fields:      []string{"id", "cidrs"},
	}

	var actual []endpointgroups.EndpointGroup
	pages := 0
	err := endpointgroups.List(fake.ServiceClient(fakeServer), options).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		pages++
		groups, err := endpointgroups.ExtractEndpointGroups(page)
		if err != nil {
			return false, err
		}
		actual = append(actual, groups...)
		return true, nil
	})
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 2, pages)
	th.CheckDeepEquals(t, []endpointgroups.EndpointGroup{Group, SecondGroup}, actual)
}

func TestListOptsMarker(t *testing.T) {
	query, err := endpointgroups.ListOpts{Marker: GroupID}.ToEndpointGroupListQuery()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "?marker="+GroupID, query)
}

func TestUpdate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGroupUpdateSuccessfully(t, fakeServer)

	name := "updatedname"
	description := "updated endpoint group"
	options := endpointgroups.UpdateOpts{
		Name:        &name,
		Description: &description,
		CIDRs:       []string{"10.4.0.0/16"},
	}

	actual, err := endpointgroups.Update(context.TODO(), fake.ServiceClient(fakeServer), GroupID, options).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, UpdatedGroup, *actual)
}
