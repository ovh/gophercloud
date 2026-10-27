/*
Package endpointgroups allows management and retrieval of VPN endpoint groups
in the OpenStack VPN Gateway Service (ovpi3n).

Example to Create an Endpoint Group

	createOpts := endpointgroups.CreateOpts{
		Name: "peer-networks",
		Type: endpointgroups.TypeCIDR,
		CIDRs: []string{
			"10.2.0.0/24",
			"10.3.0.0/24",
		},
	}

	group, err := endpointgroups.Create(context.TODO(), vpnClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Show the details of a specific Endpoint Group by ID

	group, err := endpointgroups.Get(context.TODO(), vpnClient, "6ecd9cf3-ca64-46c7-863f-f2eb1b9e838a").Extract()
	if err != nil {
		panic(err)
	}

Example to Delete an Endpoint Group

	err := endpointgroups.Delete(context.TODO(), vpnClient, "5291b189-fd84-46e5-84bd-78f40c05d69c").ExtractErr()
	if err != nil {
		panic(err)
	}

Example to List Endpoint Groups

	allPages, err := endpointgroups.List(vpnClient, nil).AllPages(context.TODO())
	if err != nil {
		panic(err)
	}

	allGroups, err := endpointgroups.ExtractEndpointGroups(allPages)
	if err != nil {
		panic(err)
	}

Example to Update an Endpoint Group

	name := "updatedname"
	description := "updated endpoint group"
	updateOpts := endpointgroups.UpdateOpts{
		Name:        &name,
		Description: &description,
	}

	updatedGroup, err := endpointgroups.Update(context.TODO(), vpnClient, "5c561d9d-eaea-45f6-ae3e-08d1a7080828", updateOpts).Extract()
	if err != nil {
		panic(err)
	}
*/
package endpointgroups
