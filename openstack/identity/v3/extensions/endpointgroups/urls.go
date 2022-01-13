package endpointgroups

import "github.com/gophercloud/gophercloud"

const (
	endpointGroupPath             = "OS-EP-FILTER/endpoint_groups"
	endpointGroupsAssociationPath = "OS-EP-FILTER/projects"
)

func rootURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL(endpointGroupPath)
}

func resourceURL(client *gophercloud.ServiceClient, endointGroupID string) string {
	return client.ServiceURL(endpointGroupPath, endointGroupID)
}

func projectAssociationURL(client *gophercloud.ServiceClient, endpointGroupId string, projectId string) string {
	return client.ServiceURL(endpointGroupPath, endpointGroupId, "projects", projectId)
}

func listEndpointGroupsAssociationURL(client *gophercloud.ServiceClient, projectId string) string {
	return client.ServiceURL(endpointGroupsAssociationPath, projectId, "endpoint_groups")
}
