/*
Package ipsecpolicies allows management and retrieval of IPsec policies in the
OpenStack VPN Gateway Service (ovpi3n).

Example to Create an IPsec policy

	createOpts := ipsecpolicies.CreateOpts{
		Name:                "ipsecpolicy1",
		Description:         "Description of ipsecpolicy1",
		EncryptionAlgorithm: ipsecpolicies.EncryptionAlgorithmAES256GCM,
		PFS:                 ipsecpolicies.PFSGroup19,
	}

	policy, err := ipsecpolicies.Create(context.TODO(), vpnClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Show the details of a specific IPsec policy by ID

	policy, err := ipsecpolicies.Get(context.TODO(), vpnClient, "f2b08c1e-aa81-4668-8ae1-1401bcb0576c").Extract()
	if err != nil {
		panic(err)
	}

Example to Delete a Policy

	err := ipsecpolicies.Delete(context.TODO(), vpnClient, "5291b189-fd84-46e5-84bd-78f40c05d69c").ExtractErr()
	if err != nil {
		panic(err)
	}

Example to Update an IPsec policy

	name := "updatedname"
	description := "updated policy"
	updateOpts := ipsecpolicies.UpdateOpts{
		Name:        &name,
		Description: &description,
		Lifetime: &ipsecpolicies.LifetimeOpts{
			Units: ipsecpolicies.UnitSeconds,
			Value: 1800,
		},
	}
	updatedPolicy, err := ipsecpolicies.Update(context.TODO(), vpnClient, "5c561d9d-eaea-45f6-ae3e-08d1a7080828", updateOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to List IPsec policies

	allPages, err := ipsecpolicies.List(vpnClient, nil).AllPages(context.TODO())
	if err != nil {
		panic(err)
	}

	allPolicies, err := ipsecpolicies.ExtractPolicies(allPages)
	if err != nil {
		panic(err)
	}
*/
package ipsecpolicies
