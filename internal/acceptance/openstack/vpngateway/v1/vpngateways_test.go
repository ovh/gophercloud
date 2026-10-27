//go:build acceptance || vpngateway || vpngateways

package v1

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	networking "github.com/gophercloud/gophercloud/v2/internal/acceptance/openstack/networking/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/vpngateway/v1/vpngateways"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestVPNGatewaysList(t *testing.T) {
	client := NewClient(t)

	allPages, err := vpngateways.List(client, nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allGateways, err := vpngateways.ExtractVPNGateways(allPages)
	th.AssertNoErr(t, err)

	for _, gateway := range allGateways {
		tools.PrintResource(t, gateway)
	}
}

func TestVPNGatewaysCRUD(t *testing.T) {
	client := NewClient(t)

	gateway, err := CreateVPNGateway(t, client)
	th.AssertNoErr(t, err)
	defer DeleteVPNGateway(t, client, gateway.ID)

	tools.PrintResource(t, gateway)

	th.AssertEquals(t, vpngateways.StatusActive, gateway.Status)
	th.AssertDeepEquals(t, []vpngateways.Route{
		{DestinationCIDR: "10.20.0.0/16", NextHop: "10.0.0.254"},
	}, gateway.Routes)

	listOpts := vpngateways.ListOpts{
		Name:   gateway.Name,
		Status: vpngateways.StatusActive,
		Fields: []string{"id", "name", "status"},
	}
	allPages, err := vpngateways.List(client, listOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	allGateways, err := vpngateways.ExtractVPNGateways(allPages)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 1, len(allGateways))
	th.AssertEquals(t, gateway.ID, allGateways[0].ID)
	th.AssertEquals(t, "", allGateways[0].ProjectID)

	name := tools.RandomString("TESTACC-", 8)
	description := tools.RandomString("TESTACC-desc-", 8)
	routes := []vpngateways.RouteOpts{
		{DestinationCIDR: "10.30.0.0/16", NextHop: "10.0.0.253"},
	}
	updateOpts := vpngateways.UpdateOpts{
		Name:        &name,
		Description: &description,
		Routes:      &routes,
	}

	updatedGateway, err := vpngateways.Update(context.TODO(), client, gateway.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, updatedGateway)

	th.AssertEquals(t, name, updatedGateway.Name)
	th.AssertEquals(t, description, updatedGateway.Description)
	th.AssertDeepEquals(t, []vpngateways.Route{
		{DestinationCIDR: "10.30.0.0/16", NextHop: "10.0.0.253"},
	}, updatedGateway.Routes)

	noRoutes := []vpngateways.RouteOpts{}
	updatedGateway, err = vpngateways.Update(context.TODO(), client, gateway.ID, vpngateways.UpdateOpts{Routes: &noRoutes}).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 0, len(updatedGateway.Routes))
	th.AssertEquals(t, name, updatedGateway.Name)
}

func TestVPNGatewaysAttachments(t *testing.T) {
	client := NewClient(t)

	networkClient, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)
	externalSubnetID := ExternalSubnetID(t, networkClient)

	network, err := networking.CreateNetwork(t, networkClient)
	th.AssertNoErr(t, err)
	defer networking.DeleteNetwork(t, networkClient, network.ID)

	subnet, err := networking.CreateSubnet(t, networkClient, network.ID)
	th.AssertNoErr(t, err)
	defer networking.DeleteSubnet(t, networkClient, subnet.ID)

	gateway, err := CreateVPNGateway(t, client)
	th.AssertNoErr(t, err)
	defer DeleteVPNGateway(t, client, gateway.ID)

	externalOpts := vpngateways.AddExternalGatewaysOpts{
		ExternalGateways: []vpngateways.ExternalGatewayOpts{
			{SubnetID: externalSubnetID},
		},
	}
	gateway, err = vpngateways.AddExternalGateways(context.TODO(), client, gateway.ID, externalOpts).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 1, len(gateway.ExternalGateways))

	gateway, err = WaitForAttachments(client, gateway.ID)
	th.AssertNoErr(t, err)
	tools.PrintResource(t, gateway)
	if gateway.ExternalIPv4 == "" {
		t.Fatal("Expected an external IPv4 address once the external gateway is ACTIVE")
	}

	iface, err := vpngateways.AddInterface(context.TODO(), client, gateway.ID, vpngateways.AddInterfaceOpts{
		SubnetID: subnet.ID,
	}).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, iface)
	th.AssertEquals(t, gateway.ID, iface.ID)
	th.AssertEquals(t, subnet.ID, iface.SubnetID)

	gateway, err = WaitForAttachments(client, gateway.ID)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 1, len(gateway.InternalInterfaces))
	th.AssertEquals(t, network.ID, gateway.InternalInterfaces[0].NetworkID)
	th.AssertEquals(t, subnet.ID, gateway.InternalInterfaces[0].FixedIPs[0].SubnetID)

	iface, err = vpngateways.RemoveInterface(context.TODO(), client, gateway.ID, vpngateways.RemoveInterfaceOpts{
		SubnetID: subnet.ID,
	}).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, subnet.ID, iface.SubnetID)

	err = tools.WaitFor(func(ctx context.Context) (bool, error) {
		gateway, err = vpngateways.Get(ctx, client, gateway.ID).Extract()
		if err != nil {
			return false, err
		}
		return len(gateway.InternalInterfaces) == 0, nil
	})
	th.AssertNoErr(t, err)
}

func TestVPNGatewaysResync(t *testing.T) {
	clients.RequireAdmin(t)

	client := NewClient(t)

	gateway, err := CreateVPNGateway(t, client)
	th.AssertNoErr(t, err)
	defer DeleteVPNGateway(t, client, gateway.ID)

	if len(gateway.Instances) == 0 {
		t.Fatal("Expected an admin to see the gateway instances")
	}

	gateway, err = vpngateways.Resync(context.TODO(), client, gateway.ID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, vpngateways.StatusPendingUpdate, gateway.Status)

	err = WaitForVPNGatewayStatus(client, gateway.ID, vpngateways.StatusActive)
	th.AssertNoErr(t, err)
}

func TestVPNGatewaysRebuild(t *testing.T) {
	clients.RequireAdmin(t)
	clients.RequireLong(t)

	client := NewClient(t)

	gateway, err := CreateVPNGateway(t, client)
	th.AssertNoErr(t, err)
	defer DeleteVPNGateway(t, client, gateway.ID)

	imageID := gateway.Instances[0].ImageID

	gateway, err = vpngateways.Rebuild(context.TODO(), client, gateway.ID, vpngateways.RebuildOpts{
		ImageID: imageID,
	}).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, vpngateways.StatusPendingUpdate, gateway.Status)

	err = WaitForVPNGatewayStatus(client, gateway.ID, vpngateways.StatusActive)
	th.AssertNoErr(t, err)

	gateway, err = vpngateways.Get(context.TODO(), client, gateway.ID).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, gateway)
	th.AssertEquals(t, imageID, gateway.Instances[0].ImageID)
}
