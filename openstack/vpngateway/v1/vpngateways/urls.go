package vpngateways

import "github.com/gophercloud/gophercloud/v2"

const resourcePath = "vpn_gateways"

func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}

func resourceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(resourcePath, id)
}

func addInterfaceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(resourcePath, id, "add_router_interface")
}

func removeInterfaceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(resourcePath, id, "remove_router_interface")
}

func addExternalGatewaysURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(resourcePath, id, "add_external_gateways")
}

func rebuildURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(resourcePath, id, "rebuild")
}

func resyncURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(resourcePath, id, "resync")
}
