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

// ListOpts allows the filtering and sorting of paginated collections through
// the Neutron API. Filtering is achieved by passing in struct field values
// that map to the VPC attributes you want to see returned.
// SortKey allows you to sort by a particular VPC attribute.
// SortDir sets the direction, and is either `asc' or `desc'.
// Marker and Limit are used for the pagination.
type ListOpts struct {
	ID        string `q:"id"`
	Name      string `q:"name"`
	CIDRBlock string `q:"cidr_block"`
	VpcStatus string `q:"vpc_status"`
	TenantID  string `q:"tenant_id"`
	ProjectID string `q:"project_id"`
	Limit     int    `q:"limit"`
	Marker    string `q:"marker"`
	SortKey   string `q:"sort_key"`
	SortDir   string `q:"sort_dir"`
}

// ToVpcListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToVpcListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List returns a Pager which allows you to iterate over a collection of
// VPCs. It accepts a ListOpts struct, which allows you to filter and sort
// the returned collection for greater efficiency.
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
		return VpcPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves a specific VPC based on its ID.
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

// CreateOpts specifies parameters of a new VPC.
type CreateOpts struct {
	// Name is the human-readable name of the VPC.
	Name string `json:"name,omitempty"`

	// Description is the human-readable description of the VPC.
	Description string `json:"description,omitempty"`

	// CIDRBlock is the RFC1918 range the VPC allocates its subnets from.
	// Its prefix must be between /16 and /28. It cannot be updated.
	CIDRBlock string `json:"cidr_block" required:"true"`

	// ProjectID is the ID of the project that owns the VPC. Only
	// administrative users can specify a project other than their own.
	ProjectID string `json:"project_id,omitempty"`
}

// ToVpcCreateMap constructs a request body from CreateOpts.
func (opts CreateOpts) ToVpcCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "vpc")
}

// Create accepts a CreateOpts struct and creates a new VPC using the values
// provided. The VPC is PENDING_CREATE until the backend provisions it.
func Create(ctx context.Context, c *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToVpcCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Post(ctx, createURL(c), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
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
	// Name is the human-readable name of the VPC.
	Name *string `json:"name,omitempty"`

	// Description is the human-readable description of the VPC.
	Description *string `json:"description,omitempty"`
}

// ToVpcUpdateMap builds a request body from UpdateOpts.
func (opts UpdateOpts) ToVpcUpdateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "vpc")
}

// Update accepts an UpdateOpts struct and updates an existing VPC using the
// values provided. The request fails with a 409 while the VPC is
// transitioning.
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

// Delete accepts a unique ID and deletes the VPC associated with it. The
// request fails with a 409 while networks, subnets, ports or a router are
// still attached to the VPC, or while it is transitioning.
func Delete(ctx context.Context, c *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := c.Delete(ctx, deleteURL(c, id), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
