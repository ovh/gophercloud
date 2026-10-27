/*
Package vpngateways allows management and retrieval of VPN gateways in the
OpenStack VPN Gateway Service (ovpi3n). A VPN gateway is shaped like a Neutron
router: it is attached to an external network and to internal subnets, and
terminates the IPsec tunnels of its sites.

Creating, attaching and deleting are asynchronous: poll the gateway until its
Status, and the Status of each of its attachments, leaves PENDING_*.

Example to Create a VPN Gateway

	createOpts := vpngateways.CreateOpts{
		Name:        "gateway1",
		Description: "Description of gateway1",
	}

	gateway, err := vpngateways.Create(context.TODO(), vpnClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Show the details of a specific VPN Gateway by ID

	gateway, err := vpngateways.Get(context.TODO(), vpnClient, "2ad3c5dc-8c7e-4ff0-8b54-4b4d4e2c7c6e").Extract()
	if err != nil {
		panic(err)
	}

Example to List VPN Gateways

	allPages, err := vpngateways.List(vpnClient, vpngateways.ListOpts{Status: vpngateways.StatusActive}).AllPages(context.TODO())
	if err != nil {
		panic(err)
	}

	allGateways, err := vpngateways.ExtractVPNGateways(allPages)
	if err != nil {
		panic(err)
	}

Example to Update a VPN Gateway

	name := "updatedname"
	routes := []vpngateways.RouteOpts{
		{DestinationCIDR: "10.20.0.0/16", NextHop: "10.0.0.254"},
	}
	updateOpts := vpngateways.UpdateOpts{
		Name:   &name,
		Routes: &routes,
	}

	gateway, err := vpngateways.Update(context.TODO(), vpnClient, "2ad3c5dc-8c7e-4ff0-8b54-4b4d4e2c7c6e", updateOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Attach a VPN Gateway to an External Network

	opts := vpngateways.AddExternalGatewaysOpts{
		ExternalGateways: []vpngateways.ExternalGatewayOpts{
			{SubnetID: "a2f1f29d-571b-4533-907f-5803ab96ead1"},
		},
	}

	gateway, err := vpngateways.AddExternalGateways(context.TODO(), vpnClient, "2ad3c5dc-8c7e-4ff0-8b54-4b4d4e2c7c6e", opts).Extract()
	if err != nil {
		panic(err)
	}

Example to Add an Interface to a VPN Gateway

	opts := vpngateways.AddInterfaceOpts{
		SubnetID: "a2f1f29d-571b-4533-907f-5803ab96ead1",
	}

	iface, err := vpngateways.AddInterface(context.TODO(), vpnClient, "2ad3c5dc-8c7e-4ff0-8b54-4b4d4e2c7c6e", opts).Extract()
	if err != nil {
		panic(err)
	}

Example to Remove an Interface from a VPN Gateway

	opts := vpngateways.RemoveInterfaceOpts{
		SubnetID: "a2f1f29d-571b-4533-907f-5803ab96ead1",
	}

	iface, err := vpngateways.RemoveInterface(context.TODO(), vpnClient, "2ad3c5dc-8c7e-4ff0-8b54-4b4d4e2c7c6e", opts).Extract()
	if err != nil {
		panic(err)
	}

Example to Rebuild a VPN Gateway on a given Image

	opts := vpngateways.RebuildOpts{
		ImageID: "8ab9a0a5-7bce-4c2e-9a9e-5ff04e3a2f4d",
	}

	gateway, err := vpngateways.Rebuild(context.TODO(), vpnClient, "2ad3c5dc-8c7e-4ff0-8b54-4b4d4e2c7c6e", opts).Extract()
	if err != nil {
		panic(err)
	}

Example to Resync a VPN Gateway

	gateway, err := vpngateways.Resync(context.TODO(), vpnClient, "2ad3c5dc-8c7e-4ff0-8b54-4b4d4e2c7c6e").Extract()
	if err != nil {
		panic(err)
	}

Example to Delete a VPN Gateway

	err := vpngateways.Delete(context.TODO(), vpnClient, "2ad3c5dc-8c7e-4ff0-8b54-4b4d4e2c7c6e").ExtractErr()
	if err != nil {
		panic(err)
	}
*/
package vpngateways
