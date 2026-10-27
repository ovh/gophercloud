/*
Package siteconnections allows management and retrieval of IPsec site
connections in the OpenStack VPN Gateway Service (ovpi3n). A connection is one
tunnel to an IPsec site; a site fails over between at most two connections,
the lower Priority being preferred.

Connections are provisioned asynchronously: poll a connection until its
Status leaves PENDING_*. ACTIVE means the tunnel is configured on the VPN
gateway, not that it is up.

Example to Create an IPsec Site Connection

	createOpts := siteconnections.CreateOpts{
		Name:          "tunnel1",
		IPsecSiteID:   "9d1e6c2a-4b3f-4a5e-8d7c-6b5a4f3e2d1c",
		PeerAddress:   "198.51.100.20",
		PSK:           "a-long-pre-shared-key",
		IKEPolicyID:   "f2b08c1e-aa81-4668-8ae1-1401bcb0576c",
		IPsecPolicyID: "e3b0c442-98fc-4c14-9afb-f4c8996fb924",
		Priority:      1,
		DPD: &siteconnections.DPDOpts{
			Action:   siteconnections.DPDActionRestart,
			Interval: 10,
		},
	}

	connection, err := siteconnections.Create(context.TODO(), vpnClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Show the details of a specific IPsec Site Connection by ID

	connection, err := siteconnections.Get(context.TODO(), vpnClient, "5c561d9d-eaea-45f6-ae3e-08d1a7080828").Extract()
	if err != nil {
		panic(err)
	}

Example to List the Connections of an IPsec Site

	listOpts := siteconnections.ListOpts{
		IPsecSiteID: "9d1e6c2a-4b3f-4a5e-8d7c-6b5a4f3e2d1c",
	}

	allPages, err := siteconnections.List(vpnClient, listOpts).AllPages(context.TODO())
	if err != nil {
		panic(err)
	}

	allConnections, err := siteconnections.ExtractConnections(allPages)
	if err != nil {
		panic(err)
	}

Example to Update an IPsec Site Connection

	updateOpts := siteconnections.UpdateOpts{
		PSK: "a-new-pre-shared-key",
		MTU: 1380,
	}

	connection, err := siteconnections.Update(context.TODO(), vpnClient, "5c561d9d-eaea-45f6-ae3e-08d1a7080828", updateOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Delete an IPsec Site Connection

	err := siteconnections.Delete(context.TODO(), vpnClient, "5c561d9d-eaea-45f6-ae3e-08d1a7080828").ExtractErr()
	if err != nil {
		panic(err)
	}
*/
package siteconnections
