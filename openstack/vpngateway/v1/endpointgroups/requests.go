package endpointgroups

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

type EndpointType string

const (
	TypeCIDR EndpointType = "cidr"
)

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToEndpointGroupCreateMap() (map[string]any, error)
}

// CreateOpts contains all the values needed to create a new endpoint group.
type CreateOpts struct {
	// ProjectID specifies a project to own the endpoint group. The caller must
	// be allowed to act on that project in order to set this. Otherwise, the
	// caller's project will be the owner.
	ProjectID string `json:"project_id,omitempty"`

	// Name is the human readable name of the endpoint group.
	Name string `json:"name,omitempty"`

	// Description is the human readable description of the endpoint group.
	Description string `json:"description,omitempty"`

	// Type is the type of the endpoints in the group. Only cidr is accepted.
	Type EndpointType `json:"type" required:"true"`

	// CIDRs is the list of networks of the group, at least one.
	CIDRs []string `json:"cidrs" required:"true"`
}

// ToEndpointGroupCreateMap casts a CreateOpts struct to a map.
func (opts CreateOpts) ToEndpointGroupCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "endpoint_group")
}

// Create accepts a CreateOpts struct and uses the values to create a new
// endpoint group.
func Create(ctx context.Context, c *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToEndpointGroupCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Post(ctx, rootURL(c), b, &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Get retrieves a particular endpoint group based on its unique ID.
func Get(ctx context.Context, c *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := c.Get(ctx, resourceURL(c, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToEndpointGroupListQuery() (string, error)
}

// ListOpts allows the filtering, sorting and paging of endpoint groups. The
// service rejects any other query parameter.
type ListOpts struct {
	// ProjectID lists the endpoint groups of this project only. Naming a
	// project other than the caller's requires the get_all-global policy.
	ProjectID string `q:"project_id"`

	// Name lists the endpoint groups with this name only.
	Name string `q:"name"`

	// Limit is the maximum number of endpoint groups per page.
	Limit int `q:"limit"`

	// Marker is the ID of the last endpoint group of the previous page.
	Marker string `q:"marker"`

	// PageReverse pages backwards from Marker.
	PageReverse bool `q:"page_reverse"`

	// SortKey is the attribute to sort on.
	SortKey string `q:"sort_key"`

	// SortDir is the sort direction, asc or desc. It requires SortKey.
	SortDir string `q:"sort_dir"`

	// Fields restricts each returned endpoint group to these attributes.
	Fields []string `q:"fields"`
}

// ToEndpointGroupListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToEndpointGroupListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List returns a Pager which allows you to iterate over a collection of
// endpoint groups. It accepts a ListOpts struct, which allows you to filter
// the returned collection for greater efficiency.
func List(c *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := rootURL(c)
	if opts != nil {
		query, err := opts.ToEndpointGroupListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return EndpointGroupPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Delete will permanently delete a particular endpoint group based on its
// unique ID. A group still used by an IPsec site cannot be deleted.
func Delete(ctx context.Context, c *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := c.Delete(ctx, resourceURL(c, id), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToEndpointGroupUpdateMap() (map[string]any, error)
}

// UpdateOpts contains the values used when updating an endpoint group.
type UpdateOpts struct {
	// Name is the human readable name of the endpoint group.
	Name *string `json:"name,omitempty"`

	// Description is the human readable description of the endpoint group.
	Description *string `json:"description,omitempty"`

	// CIDRs replaces the networks of the group. They cannot change while an
	// IPsec site uses the group.
	CIDRs []string `json:"cidrs,omitempty"`
}

// ToEndpointGroupUpdateMap casts an UpdateOpts struct to a map.
func (opts UpdateOpts) ToEndpointGroupUpdateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "endpoint_group")
}

// Update allows endpoint groups to be updated.
func Update(ctx context.Context, c *gophercloud.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToEndpointGroupUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Put(ctx, resourceURL(c, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
