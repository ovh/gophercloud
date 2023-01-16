package endpointgroups

import "github.com/gophercloud/gophercloud"

const (
	endpointGroupPath             = "OS-EP-FILTER/endpoint_groups"
	endpointGroupsAssociationPath = "OS-EP-FILTER/projects"
)

func rootURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL(endpointGroupPath)
}

func resourceURL(client *gophercloud.ServiceClient, endpointGroupID string) string {
	return client.ServiceURL(endpointGroupPath, endpointGroupID)
}

func projectAssociationURL(client *gophercloud.ServiceClient, endpointGroupID string, projectID string) string {
	return client.ServiceURL(endpointGroupPath, endpointGroupID, "projects", projectID)
}

func listEndpointGroupsAssociationURL(client *gophercloud.ServiceClient, projectID string) string {
	return client.ServiceURL(endpointGroupsAssociationPath, projectID, "endpoint_groups")
}

func createURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL(endpointGroupPath)
}

func updateURL(client *gophercloud.ServiceClient, endpointGroupID string) string {
	return client.ServiceURL(endpointGroupPath, endpointGroupID)
}

func deleteURL(client *gophercloud.ServiceClient, endpointGroupID string) string {
	return client.ServiceURL(endpointGroupPath, endpointGroupID)
}
