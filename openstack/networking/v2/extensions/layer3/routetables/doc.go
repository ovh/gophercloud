/*
Package routetables provides the ability to retrieve and manage route tables
through the vpc extension of the Neutron API.

A route table belongs to a VPC router and holds routes and subnet
associations. Routes and subnets are managed through the AddRoutes,
RemoveRoutes, AddSubnets and RemoveSubnets actions. A route table cannot be
deleted while subnets are still associated with it.

Example of Listing Route Tables

	listOpts := routetables.ListOpts{
		RouterID: "a5c3a4d4-5e0f-4c2a-bb02-e2723de5b8f1",
	}

	allPages, err := routetables.List(networkClient, listOpts).AllPages(context.TODO())
	if err != nil {
		panic(err)
	}

	allRouteTables, err := routetables.ExtractRouteTables(allPages)
	if err != nil {
		panic(err)
	}

	for _, routeTable := range allRouteTables {
		fmt.Printf("%+v\n", routeTable)
	}

Example to Create a Route Table

	createOpts := routetables.CreateOpts{
		Name:     "private",
		RouterID: "a5c3a4d4-5e0f-4c2a-bb02-e2723de5b8f1",
	}

	routeTable, err := routetables.Create(context.TODO(), networkClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Add Routes to a Route Table

	routeTableID := "e3f5f925-6c6c-4b1e-9c1e-7c6d4d7c3f2a"
	addOpts := routetables.AddRoutesOpts{
		Routes: []routetables.RouteOpts{
			{
				Destination: "192.168.1.0/24",
				Nexthop:     "10.0.0.1",
				NexthopType: "ip",
			},
		},
	}

	routeTable, err := routetables.AddRoutes(context.TODO(), networkClient, routeTableID, addOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Associate Subnets with a Route Table

	routeTableID := "e3f5f925-6c6c-4b1e-9c1e-7c6d4d7c3f2a"
	addOpts := routetables.AddSubnetsOpts{
		Subnets: []string{"08eae331-0402-425a-923c-34f7cfe39c1b"},
	}

	routeTable, err := routetables.AddSubnets(context.TODO(), networkClient, routeTableID, addOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Delete a Route Table

	routeTableID := "e3f5f925-6c6c-4b1e-9c1e-7c6d4d7c3f2a"
	err := routetables.Delete(context.TODO(), networkClient, routeTableID).ExtractErr()
	if err != nil {
		panic(err)
	}
*/
package routetables
