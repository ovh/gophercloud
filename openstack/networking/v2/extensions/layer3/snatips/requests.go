package snatips

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToSnatIPListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through
// the Neutron API. Filtering is achieved by passing in struct field values
// that map to the SNAT IP attributes you want to see returned.
// SortKey allows you to sort by a particular SNAT IP attribute.
// SortDir sets the direction, and is either `asc' or `desc'.
// Marker and Limit are used for the pagination.
type ListOpts struct {
	ID               string `q:"id"`
	NetworkID        string `q:"network_id"`
	AvailabilityZone string `q:"availability_zone"`
	IPAddress        string `q:"ip_address"`
	VpcID            string `q:"vpc_id"`
	Status           string `q:"status"`
	VpcStatus        string `q:"vpc_status"`
	TenantID         string `q:"tenant_id"`
	ProjectID        string `q:"project_id"`
	Limit            int    `q:"limit"`
	Marker           string `q:"marker"`
	SortKey          string `q:"sort_key"`
	SortDir          string `q:"sort_dir"`
}

// ToSnatIPListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToSnatIPListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List returns a Pager which allows you to iterate over a collection of
// SNAT IPs. It accepts a ListOpts struct, which allows you to filter and sort
// the returned collection for greater efficiency.
func List(c *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(c)
	if opts != nil {
		query, err := opts.ToSnatIPListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return SnatIPPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves a specific SNAT IP based on its ID.
func Get(ctx context.Context, c *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := c.Get(ctx, getURL(c, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToSnatIPCreateMap() (map[string]any, error)
}

// CreateOpts specifies parameters of a new SNAT IP reservation.
type CreateOpts struct {
	// NetworkID is the ID of the external network the address is reserved on.
	NetworkID string `json:"network_id" required:"true"`

	// AvailabilityZone is the availability zone the address is reserved in.
	AvailabilityZone string `json:"availability_zone" required:"true"`

	// IPAddress is the address to reserve. When omitted, the lowest free
	// address of the network's pools is taken, IPv4 before IPv6.
	IPAddress string `json:"ip_address,omitempty"`

	// Description is the human-readable description of the SNAT IP.
	Description string `json:"description,omitempty"`

	// ProjectID is the ID of the project that owns the SNAT IP. Only
	// administrative users can specify a project other than their own.
	ProjectID string `json:"project_id,omitempty"`
}

// ToSnatIPCreateMap constructs a request body from CreateOpts.
func (opts CreateOpts) ToSnatIPCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "snat_ip")
}

// Create accepts a CreateOpts struct and reserves a new SNAT IP using the
// values provided.
func Create(ctx context.Context, c *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToSnatIPCreateMap()
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
	ToSnatIPUpdateMap() (map[string]any, error)
}

// UpdateOpts represents options used to update a SNAT IP.
type UpdateOpts struct {
	// Description is the human-readable description of the SNAT IP.
	Description *string `json:"description,omitempty"`
}

// ToSnatIPUpdateMap builds a request body from UpdateOpts.
func (opts UpdateOpts) ToSnatIPUpdateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "snat_ip")
}

// Update accepts an UpdateOpts struct and updates an existing SNAT IP using
// the values provided.
func Update(ctx context.Context, c *gophercloud.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToSnatIPUpdateMap()
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

// Delete accepts a unique ID and releases the SNAT IP associated with it.
// The request fails with a 409 while a router gateway port carries the
// address, or while the reservation is transitioning.
func Delete(ctx context.Context, c *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := c.Delete(ctx, deleteURL(c, id), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
