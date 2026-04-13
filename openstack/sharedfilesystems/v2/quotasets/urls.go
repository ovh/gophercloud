package quotasets

import "github.com/gophercloud/gophercloud/v2"

const resourcePath = "quota-sets"
const resourcePathDetail = "detail"

func getURL(c *gophercloud.ServiceClient, tenantID string) string {
	return c.ServiceURL(resourcePath, tenantID)
}

func getDetailURL(c *gophercloud.ServiceClient, tenantID string) string {
	return c.ServiceURL(resourcePath, tenantID, resourcePathDetail)
}

func updateURL(c *gophercloud.ServiceClient, tenantID string) string {
	return c.ServiceURL(resourcePath, tenantID)
}

func getURLbyShareType(c *gophercloud.ServiceClient, tenantID string, share_type string) string {
	return c.ServiceURL(resourcePath, tenantID) + "?share_type=" + share_type
}

func updateURLByShareType(c *gophercloud.ServiceClient, tenantID string, share_type string) string {
	return c.ServiceURL(resourcePath, tenantID) + "?share_type=" + share_type
}

func getURLbyUser(c *gophercloud.ServiceClient, tenantID string, user_id string) string {
	return c.ServiceURL(resourcePath, tenantID) + "?user_id=" + user_id
}

func updateURLByUser(c *gophercloud.ServiceClient, tenantID string, user_id string) string {
	return c.ServiceURL(resourcePath, tenantID) + "?user_id=" + user_id
}
