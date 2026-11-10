package siteconnections

import "github.com/gophercloud/gophercloud/v2"

const resourcePath = "ipsec_site_connections"

func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}

func resourceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(resourcePath, id)
}

func pskURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(resourcePath, id, "psk")
}

func regeneratePSKURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(resourcePath, id, "regenerate_psk")
}
