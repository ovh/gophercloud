package testing

import (
	"context"
	"testing"

	fake "github.com/gophercloud/gophercloud/v2/openstack/vpngateway/v1/testhelper"
	"github.com/gophercloud/gophercloud/v2/openstack/vpngateway/v1/vpngateways"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestCreate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGatewayCreationSuccessfully(t, fakeServer)

	adminStateUp := true
	ha := false
	options := vpngateways.CreateOpts{
		ProjectID:             "9145d91459d248b1b02fdaca97c6a75d",
		Name:                  "gateway",
		Description:           "VPN gateway",
		AdminStateUp:          &adminStateUp,
		HA:                    &ha,
		FlavorID:              "b8b9cc5a-0b05-4b1e-9d0a-1d6c0f5a3e01",
		AvailabilityZoneHints: []string{"az1"},
		Routes: []vpngateways.RouteOpts{
			{DestinationCIDR: "10.20.0.0/16", NextHop: "10.0.0.254"},
		},
	}

	actual, err := vpngateways.Create(context.TODO(), fake.ServiceClient(fakeServer), options).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, CreatedGateway, *actual)
}

func TestCreateRequiresRouteParts(t *testing.T) {
	options := vpngateways.CreateOpts{
		Routes: []vpngateways.RouteOpts{{DestinationCIDR: "10.20.0.0/16"}},
	}
	_, err := options.ToVPNGatewayCreateMap()
	if err == nil {
		t.Fatal("expected an error for a route without a next hop")
	}
}

func TestGet(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGatewayGetSuccessfully(t, fakeServer)

	actual, err := vpngateways.Get(context.TODO(), fake.ServiceClient(fakeServer), GatewayID).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, Gateway, *actual)
}

func TestDelete(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGatewayDeletionSuccessfully(t, fakeServer)

	res := vpngateways.Delete(context.TODO(), fake.ServiceClient(fakeServer), GatewayID)
	th.AssertNoErr(t, res.Err)
}

func TestList(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGatewayListSuccessfully(t, fakeServer)

	options := vpngateways.ListOpts{
		ProjectID:   "9145d91459d248b1b02fdaca97c6a75d",
		Name:        "gateway",
		Status:      vpngateways.StatusActive,
		FlavorID:    "b8b9cc5a-0b05-4b1e-9d0a-1d6c0f5a3e01",
		Limit:       1,
		PageReverse: true,
		SortKey:     "name",
		SortDir:     "desc",
		Fields:      []string{"id", "name", "status"},
	}

	var actual []vpngateways.VPNGateway
	pages := 0
	err := vpngateways.List(fake.ServiceClient(fakeServer), options).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		pages++
		gateways, err := vpngateways.ExtractVPNGateways(page)
		if err != nil {
			return false, err
		}
		actual = append(actual, gateways...)
		return true, nil
	})
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 2, pages)
	th.CheckDeepEquals(t, ListedGateways, actual)
}

func TestListOptsMarker(t *testing.T) {
	query, err := vpngateways.ListOpts{Marker: GatewayID}.ToVPNGatewayListQuery()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "?marker="+GatewayID, query)
}

func TestUpdate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGatewayUpdateSuccessfully(t, fakeServer)

	name := "updatedname"
	description := "updated gateway"
	adminStateUp := false
	ha := true
	flavorID := "1d2c3b4a-5e6f-4a7b-8c9d-0e1f2a3b4c51"
	azHints := []string{"az2"}
	routes := []vpngateways.RouteOpts{
		{DestinationCIDR: "10.30.0.0/16", NextHop: "10.0.0.253"},
	}
	options := vpngateways.UpdateOpts{
		Name:                  &name,
		Description:           &description,
		AdminStateUp:          &adminStateUp,
		HA:                    &ha,
		FlavorID:              &flavorID,
		AvailabilityZoneHints: &azHints,
		Routes:                &routes,
	}

	actual, err := vpngateways.Update(context.TODO(), fake.ServiceClient(fakeServer), GatewayID, options).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, UpdatedGateway, *actual)
}

func TestUpdateClearRoutes(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGatewayClearRoutesSuccessfully(t, fakeServer)

	routes := []vpngateways.RouteOpts{}
	options := vpngateways.UpdateOpts{
		Routes: &routes,
	}

	_, err := vpngateways.Update(context.TODO(), fake.ServiceClient(fakeServer), GatewayID, options).Extract()
	th.AssertNoErr(t, err)
}

func TestAddInterface(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleAddInterfaceSuccessfully(t, fakeServer)

	options := vpngateways.AddInterfaceOpts{
		SubnetID: "d4c3b2a1-0f9e-4d8c-b7a6-5f4e3d2c1b01",
	}

	actual, err := vpngateways.AddInterface(context.TODO(), fake.ServiceClient(fakeServer), GatewayID, options).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, AddedInterface, *actual)
}

func TestRemoveInterface(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleRemoveInterfaceSuccessfully(t, fakeServer)

	options := vpngateways.RemoveInterfaceOpts{
		SubnetID: "d4c3b2a1-0f9e-4d8c-b7a6-5f4e3d2c1b01",
	}

	actual, err := vpngateways.RemoveInterface(context.TODO(), fake.ServiceClient(fakeServer), GatewayID, options).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, RemovedInterface, *actual)
}

func TestAddExternalGateways(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleAddExternalGatewaysSuccessfully(t, fakeServer)

	options := vpngateways.AddExternalGatewaysOpts{
		ExternalGateways: []vpngateways.ExternalGatewayOpts{
			{SubnetID: "a2f1f29d-571b-4533-907f-5803ab96ead1"},
			{SubnetID: "c6d1a3b2-0e5f-4a8b-9c7d-6e5f4a3b2c11"},
		},
	}

	actual, err := vpngateways.AddExternalGateways(context.TODO(), fake.ServiceClient(fakeServer), GatewayID, options).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, AttachingGateway, *actual)
}

func TestRebuild(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleRebuildSuccessfully(t, fakeServer, RebuildRequest)

	options := vpngateways.RebuildOpts{
		ImageID: "8ab9a0a5-7bce-4c2e-9a9e-5ff04e3a2f4d",
	}

	actual, err := vpngateways.Rebuild(context.TODO(), fake.ServiceClient(fakeServer), GatewayID, options).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, PendingUpdateGateway, *actual)
}

func TestRebuildNewestImage(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleRebuildSuccessfully(t, fakeServer, `{}`)

	actual, err := vpngateways.Rebuild(context.TODO(), fake.ServiceClient(fakeServer), GatewayID, vpngateways.RebuildOpts{}).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, PendingUpdateGateway, *actual)
}

func TestResync(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleResyncSuccessfully(t, fakeServer)

	actual, err := vpngateways.Resync(context.TODO(), fake.ServiceClient(fakeServer), GatewayID).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, PendingUpdateGateway, *actual)
}
