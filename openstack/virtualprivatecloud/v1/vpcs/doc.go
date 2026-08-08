/*
Package vpcs provides functionality for managing Orion VPC resources.

A VPC (Virtual Private Cloud) is an isolated virtual network with an IPv4 CIDR
block. VPCs are managed by the Orion service. CREATE and DELETE are asynchronous
(return 202), while UPDATE is synchronous (return 200).

Example to List VPCs

	listOpts := vpcs.ListOpts{
		Status: "READY",
	}

	allPages, err := vpcs.List(orionClient, listOpts).AllPages(context.TODO())
	if err != nil {
		panic(err)
	}

	allVpcs, err := vpcs.ExtractVpcs(allPages)
	if err != nil {
		panic(err)
	}

	for _, vpc := range allVpcs {
		fmt.Printf("%+v\n", vpc)
	}

Example to Create a VPC

	createOpts := vpcs.CreateOpts{
		Name:      "my-vpc",
		CIDRBlock: "10.0.0.0/16",
	}

	vpc, err := vpcs.Create(context.TODO(), orionClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Get a VPC

	vpc, err := vpcs.Get(context.TODO(), orionClient, "f47ac10b-58cc-4372-a567-0e02b2c3d479").Extract()
	if err != nil {
		panic(err)
	}

Example to Update a VPC

	name := "my-vpc-renamed"
	updateOpts := vpcs.UpdateOpts{
		Name: &name,
	}

	vpc, err := vpcs.Update(context.TODO(), orionClient, "f47ac10b-58cc-4372-a567-0e02b2c3d479", updateOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Delete a VPC

	vpc, err := vpcs.Delete(context.TODO(), orionClient, "f47ac10b-58cc-4372-a567-0e02b2c3d479").Extract()
	if err != nil {
		panic(err)
	}
*/
package vpcs
