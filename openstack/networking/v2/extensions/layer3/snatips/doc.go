/*
Package snatips provides the ability to retrieve and manage SNAT IPs through
the vpc extension of the Neutron API.

A SNAT IP is a public address reserved on an external network, in a given
availability zone, so that the gateway port of a VPC router can be given this
very address. Its Status is ACTIVE once a router gateway port carries the
address, and DOWN while none does.

Example of Listing SNAT IPs

	listOpts := snatips.ListOpts{
		NetworkID: "8f9a1c2b-3d4e-4f5a-9b6c-7d8e9f0a1b2c",
	}

	allPages, err := snatips.List(networkClient, listOpts).AllPages(context.TODO())
	if err != nil {
		panic(err)
	}

	allSnatIPs, err := snatips.ExtractSnatIPs(allPages)
	if err != nil {
		panic(err)
	}

	for _, snatIP := range allSnatIPs {
		fmt.Printf("%+v\n", snatIP)
	}

Example to Get a SNAT IP

	snatIPID := "2c4e6a8b-1d3f-4a5c-8e7b-9d0f1a2b3c4d"
	snatIP, err := snatips.Get(context.TODO(), networkClient, snatIPID).Extract()
	if err != nil {
		panic(err)
	}

Example to Reserve a SNAT IP

	createOpts := snatips.CreateOpts{
		NetworkID:        "8f9a1c2b-3d4e-4f5a-9b6c-7d8e9f0a1b2c",
		AvailabilityZone: "eu-west-par-a",
	}

	snatIP, err := snatips.Create(context.TODO(), networkClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Update a SNAT IP

	snatIPID := "2c4e6a8b-1d3f-4a5c-8e7b-9d0f1a2b3c4d"
	description := "NAT for production"
	updateOpts := snatips.UpdateOpts{
		Description: &description,
	}

	snatIP, err := snatips.Update(context.TODO(), networkClient, snatIPID, updateOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Release a SNAT IP

	snatIPID := "2c4e6a8b-1d3f-4a5c-8e7b-9d0f1a2b3c4d"
	err := snatips.Delete(context.TODO(), networkClient, snatIPID).ExtractErr()
	if err != nil {
		panic(err)
	}
*/
package snatips
