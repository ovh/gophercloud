package vpcs

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToVpcListQuery() (string, error)
}

// ListOpts allows the filtering of paginated collections through the API.
// Filtering is achieved by passing in struct field values that map to the VPC
// attributes you want to see returned.
type ListOpts struct {
	// Filter by VPC name (exact match).
	Name string `q:"name"`

	// Filter by VPC status.
	Status string `q:"status"`
}

// ToVpcListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToVpcListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List returns a Pager which allows you to iterate over a collection of
// VPCs. It accepts a ListOpts struct, which allows you to filter the
// returned collection for greater efficiency.
func List(c *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(c)
	if opts != nil {
		query, err := opts.ToVpcListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return VpcPage{pagination.SinglePageBase(r)}
	})
}

// Get retrieves a specific VPC based on its unique ID.
func Get(ctx context.Context, c *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := c.Get(ctx, getURL(c, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToVpcCreateMap() (map[string]any, error)
}

// CreateOpts represents options used to create a VPC.
type CreateOpts struct {
	// Human-readable VPC name. Required.
	Name string `json:"name"`

	// IPv4 CIDR block (RFC 1918, /16 to /28). Required.
	CIDRBlock string `json:"cidr_block"`

	// Optional description.
	Description string `json:"description,omitempty"`
}

// ToVpcCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToVpcCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "vpc")
}

// Create accepts a CreateOpts struct and creates a new VPC using the values
// provided. This is an asynchronous operation that returns 202.
func Create(ctx context.Context, c *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToVpcCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Post(ctx, createURL(c), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToVpcUpdateMap() (map[string]any, error)
}

// UpdateOpts represents options used to update a VPC.
type UpdateOpts struct {
	// New VPC name.
	Name *string `json:"name,omitempty"`

	// New description.
	Description *string `json:"description,omitempty"`
}

// ToVpcUpdateMap builds a request body from UpdateOpts.
func (opts UpdateOpts) ToVpcUpdateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "vpc")
}

// Update accepts an UpdateOpts struct and updates an existing VPC using the
// values provided. Only allowed when the VPC is in READY state.
func Update(ctx context.Context, c *gophercloud.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToVpcUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Put(ctx, updateURL(c, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete accepts a unique ID and deletes the VPC associated with it.
// This is an asynchronous operation that returns 202 with the VPC in
// DELETING state.
func Delete(ctx context.Context, c *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := c.Delete(ctx, deleteURL(c, id), &gophercloud.RequestOpts{
		OkCodes:      []int{202},
		JSONResponse: &r.Body,
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
