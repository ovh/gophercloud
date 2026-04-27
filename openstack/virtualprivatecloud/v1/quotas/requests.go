package quotas

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
)

// Get returns the VPC quota for a given project ID.
func Get(ctx context.Context, client *gophercloud.ServiceClient, projectID string) (r GetResult) {
	resp, err := client.Get(ctx, resourceURL(client, projectID), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// SetOptsBuilder allows extensions to add additional parameters to the
// Set request.
type SetOptsBuilder interface {
	ToQuotaSetMap() (map[string]any, error)
}

// SetOpts represents options used to set VPC quotas for a project.
type SetOpts struct {
	// Maximum number of VPCs for this project.
	MaxVpcs int `json:"max_vpcs"`
}

// ToQuotaSetMap builds a request body from SetOpts.
func (opts SetOpts) ToQuotaSetMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "quota")
}

// Set creates or updates a per-project VPC quota override. Requires admin role.
func Set(ctx context.Context, c *gophercloud.ServiceClient, projectID string, opts SetOptsBuilder) (r SetResult) {
	b, err := opts.ToQuotaSetMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Put(ctx, resourceURL(c, projectID), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete removes the per-project quota override. The project reverts to the
// config default. Requires admin role.
func Delete(ctx context.Context, c *gophercloud.ServiceClient, projectID string) (r DeleteResult) {
	resp, err := c.Delete(ctx, resourceURL(c, projectID), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
