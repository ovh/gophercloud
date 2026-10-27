package ipsecsites

import "github.com/gophercloud/gophercloud/v2"

const resourcePath = "ipsec_sites"

func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}

func resourceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(resourcePath, id)
}
