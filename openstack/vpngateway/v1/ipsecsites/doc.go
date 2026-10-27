/*
Package ipsecsites allows management and retrieval of IPsec sites in the
OpenStack VPN Gateway Service (ovpi3n). An IPsec site holds what every tunnel
to one remote site shares: the VPN gateway, the routing scheme and the local
and peer networks. Its tunnels are IPsec site connections.

Example to Create a static IPsec Site

	createOpts := ipsecsites.CreateOpts{
		Name:           "site1",
		VPNGatewayID:   "2ad3c5dc-8c7e-4ff0-8b54-4b4d4e2c7c6e",
		Routing:        ipsecsites.RoutingStatic,
		LocalEPGroupID: "6ecd9cf3-ca64-46c7-863f-f2eb1b9e838a",
		PeerEPGroupID:  "3e1f6a2b-7c8d-4e9f-a0b1-c2d3e4f5a6b7",
	}

	site, err := ipsecsites.Create(context.TODO(), vpnClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Create a BGP IPsec Site

	createOpts := ipsecsites.CreateOpts{
		Name:           "site1",
		VPNGatewayID:   "2ad3c5dc-8c7e-4ff0-8b54-4b4d4e2c7c6e",
		Routing:        ipsecsites.RoutingBGP,
		LocalASN:       65001,
		PeerASN:        64512,
		LocalEPGroupID: "6ecd9cf3-ca64-46c7-863f-f2eb1b9e838a",
		PeerEPGroupID:  "3e1f6a2b-7c8d-4e9f-a0b1-c2d3e4f5a6b7",
	}

	site, err := ipsecsites.Create(context.TODO(), vpnClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Show the details of a specific IPsec Site by ID

	site, err := ipsecsites.Get(context.TODO(), vpnClient, "9d1e6c2a-4b3f-4a5e-8d7c-6b5a4f3e2d1c").Extract()
	if err != nil {
		panic(err)
	}

Example to List the IPsec Sites of a VPN Gateway

	listOpts := ipsecsites.ListOpts{
		VPNGatewayID: "2ad3c5dc-8c7e-4ff0-8b54-4b4d4e2c7c6e",
	}

	allPages, err := ipsecsites.List(vpnClient, listOpts).AllPages(context.TODO())
	if err != nil {
		panic(err)
	}

	allSites, err := ipsecsites.ExtractSites(allPages)
	if err != nil {
		panic(err)
	}

Example to Update an IPsec Site

	name := "updatedname"
	updateOpts := ipsecsites.UpdateOpts{
		Name: &name,
	}

	site, err := ipsecsites.Update(context.TODO(), vpnClient, "9d1e6c2a-4b3f-4a5e-8d7c-6b5a4f3e2d1c", updateOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Delete an IPsec Site

	err := ipsecsites.Delete(context.TODO(), vpnClient, "9d1e6c2a-4b3f-4a5e-8d7c-6b5a4f3e2d1c").ExtractErr()
	if err != nil {
		panic(err)
	}
*/
package ipsecsites
