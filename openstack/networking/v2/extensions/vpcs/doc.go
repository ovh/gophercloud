/*
Package vpcs provides the ability to retrieve and manage Virtual Private Clouds
(VPCs) through the vpc extension of the Neutron API.

A VPC owns an RFC1918 CIDR block that its subnets are allocated from. Networks
and routers are attached to a VPC at creation through their VpcID attribute.
Provisioning on the backend is asynchronous: a VPC is created in the
PENDING_CREATE VpcStatus and becomes ACTIVE once the backend is in sync.

Example of Listing VPCs

	listOpts := vpcs.ListOpts{
		VpcStatus: "ACTIVE",
	}

	allPages, err := vpcs.List(networkClient, listOpts).AllPages(context.TODO())
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

Example to Get a VPC

	vpcID := "6f1a5c3e-2b4d-4c8e-9f0a-1b2c3d4e5f60"
	vpc, err := vpcs.Get(context.TODO(), networkClient, vpcID).Extract()
	if err != nil {
		panic(err)
	}

Example to Create a VPC

	createOpts := vpcs.CreateOpts{
		Name:      "my_vpc",
		CIDRBlock: "10.0.0.0/16",
	}

	vpc, err := vpcs.Create(context.TODO(), networkClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Update a VPC

	vpcID := "6f1a5c3e-2b4d-4c8e-9f0a-1b2c3d4e5f60"
	newName := "new_name"
	updateOpts := vpcs.UpdateOpts{
		Name: &newName,
	}

	vpc, err := vpcs.Update(context.TODO(), networkClient, vpcID, updateOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Delete a VPC

	vpcID := "6f1a5c3e-2b4d-4c8e-9f0a-1b2c3d4e5f60"
	err := vpcs.Delete(context.TODO(), networkClient, vpcID).ExtractErr()
	if err != nil {
		panic(err)
	}
*/
package vpcs
