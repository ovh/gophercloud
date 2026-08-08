package quotas

import "github.com/gophercloud/gophercloud/v2"

func resourceURL(c *gophercloud.ServiceClient, projectID string) string {
	return c.ServiceURL("v1", "quotas", projectID)
}
