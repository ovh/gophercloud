/*
Package quotas provides the ability to retrieve and manage VPC quotas through
the Orion API.

Example to Get a Quota

	projectID := "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
	quota, err := quotas.Get(context.TODO(), orionClient, projectID).Extract()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Max VPCs: %d (source: %s)\n", quota.MaxVpcs, quota.Source)

Example to Set a Quota (admin only)

	projectID := "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
	setOpts := quotas.SetOpts{
		MaxVpcs: 10,
	}

	result, err := quotas.Set(context.TODO(), orionClient, projectID, setOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Delete a Quota Override (admin only)

	projectID := "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
	err := quotas.Delete(context.TODO(), orionClient, projectID).ExtractErr()
	if err != nil {
		panic(err)
	}
*/
package quotas
