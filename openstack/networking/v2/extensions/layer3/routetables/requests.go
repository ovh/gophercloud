package routetables

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToRouteTableListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through
// the Neutron API. Filtering is achieved by passing in struct field values
// that map to the route table attributes you want to see returned.
// SortKey allows you to sort by a particular route table attribute.
// SortDir sets the direction, and is either `asc' or `desc'.
// Marker and Limit are used for the pagination.
type ListOpts struct {
	ID        string `q:"id"`
	Name      string `q:"name"`
	RouterID  string `q:"router_id"`
	VpcStatus string `q:"vpc_status"`
	TenantID  string `q:"tenant_id"`
	ProjectID string `q:"project_id"`
	Limit     int    `q:"limit"`
	Marker    string `q:"marker"`
	SortKey   string `q:"sort_key"`
	SortDir   string `q:"sort_dir"`
}

// ToRouteTableListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToRouteTableListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List returns a Pager which allows you to iterate over a collection of
// route tables. It accepts a ListOpts struct, which allows you to filter and
// sort the returned collection for greater efficiency.
func List(c *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := rootURL(c)
	if opts != nil {
		query, err := opts.ToRouteTableListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return RouteTablePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves a specific route table based on its unique ID.
func Get(ctx context.Context, c *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := c.Get(ctx, resourceURL(c, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToRouteTableCreateMap() (map[string]any, error)
}

// CreateOpts represents options used to create a route table.
type CreateOpts struct {
	// Name is the human-readable name of the route table.
	Name string `json:"name,omitempty"`

	// Description is a human-readable description of the route table.
	Description string `json:"description,omitempty"`

	// RouterID is the ID of the VPC router this route table belongs to.
	RouterID string `json:"router_id" required:"true"`

	// ProjectID is the ID of the project that owns the route table. Only
	// administrative users can specify a project other than their own.
	ProjectID string `json:"project_id,omitempty"`
}

// ToRouteTableCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToRouteTableCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "route_table")
}

// Create accepts a CreateOpts struct and creates a new route table.
func Create(ctx context.Context, c *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToRouteTableCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Post(ctx, rootURL(c), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToRouteTableUpdateMap() (map[string]any, error)
}

// UpdateOpts represents options used to update a route table.
type UpdateOpts struct {
	// Name is the human-readable name of the route table.
	Name *string `json:"name,omitempty"`

	// Description is a human-readable description of the route table.
	Description *string `json:"description,omitempty"`
}

// ToRouteTableUpdateMap builds a request body from UpdateOpts.
func (opts UpdateOpts) ToRouteTableUpdateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "route_table")
}

// Update accepts an UpdateOpts struct and updates a route table.
func Update(ctx context.Context, c *gophercloud.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToRouteTableUpdateMap()
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

// Delete accepts a route table ID and deletes it. The request fails with a
// 409 while subnets are still associated with the route table.
func Delete(ctx context.Context, c *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := c.Delete(ctx, resourceURL(c, id), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// RouteOpts represents a route to add to a route table.
type RouteOpts struct {
	// Destination is the destination CIDR of the route.
	Destination string `json:"destination" required:"true"`

	// Nexthop is the IP address of the next hop. Only used with the ip
	// NexthopType.
	Nexthop string `json:"nexthop,omitempty"`

	// NexthopType is the type of the next hop: ip, local, blackhole,
	// unreachable or prohibit.
	NexthopType string `json:"nexthop_type" required:"true"`
}

// RemoveRouteOpts represents a route to remove from a route table.
type RemoveRouteOpts struct {
	// Destination is the destination CIDR of the route to remove.
	Destination string `json:"destination" required:"true"`
}

// AddRoutesOptsBuilder allows extensions to add additional parameters to the
// AddRoutes request.
type AddRoutesOptsBuilder interface {
	ToRouteTableAddRoutesMap() (map[string]any, error)
}

// AddRoutesOpts represents options for adding routes to a route table.
type AddRoutesOpts struct {
	// Routes is the list of routes to add.
	Routes []RouteOpts `json:"routes" required:"true"`
}

// ToRouteTableAddRoutesMap builds a request body from AddRoutesOpts.
func (opts AddRoutesOpts) ToRouteTableAddRoutesMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// AddRoutes adds routes to a route table.
func AddRoutes(ctx context.Context, c *gophercloud.ServiceClient, id string, opts AddRoutesOptsBuilder) (r AddRoutesResult) {
	b, err := opts.ToRouteTableAddRoutesMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Put(ctx, addRoutesURL(c, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// RemoveRoutesOptsBuilder allows extensions to add additional parameters to
// the RemoveRoutes request.
type RemoveRoutesOptsBuilder interface {
	ToRouteTableRemoveRoutesMap() (map[string]any, error)
}

// RemoveRoutesOpts represents options for removing routes from a route table.
type RemoveRoutesOpts struct {
	// Routes is the list of routes to remove, identified by destination.
	Routes []RemoveRouteOpts `json:"routes" required:"true"`
}

// ToRouteTableRemoveRoutesMap builds a request body from RemoveRoutesOpts.
func (opts RemoveRoutesOpts) ToRouteTableRemoveRoutesMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// RemoveRoutes removes routes from a route table.
func RemoveRoutes(ctx context.Context, c *gophercloud.ServiceClient, id string, opts RemoveRoutesOptsBuilder) (r RemoveRoutesResult) {
	b, err := opts.ToRouteTableRemoveRoutesMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Put(ctx, removeRoutesURL(c, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// AddSubnetsOptsBuilder allows extensions to add additional parameters to the
// AddSubnets request.
type AddSubnetsOptsBuilder interface {
	ToRouteTableAddSubnetsMap() (map[string]any, error)
}

// AddSubnetsOpts represents options for associating subnets with a route
// table.
type AddSubnetsOpts struct {
	// Subnets is the list of subnet IDs to associate. The subnets must be in
	// the VPC of the route table's router.
	Subnets []string `json:"subnets" required:"true"`

	// Routes is an optional list of routes added in the same operation.
	Routes []RouteOpts `json:"routes,omitempty"`
}

// ToRouteTableAddSubnetsMap builds a request body from AddSubnetsOpts.
func (opts AddSubnetsOpts) ToRouteTableAddSubnetsMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// AddSubnets associates subnets with a route table.
func AddSubnets(ctx context.Context, c *gophercloud.ServiceClient, id string, opts AddSubnetsOptsBuilder) (r AddSubnetsResult) {
	b, err := opts.ToRouteTableAddSubnetsMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Put(ctx, addSubnetsURL(c, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// RemoveSubnetsOptsBuilder allows extensions to add additional parameters to
// the RemoveSubnets request.
type RemoveSubnetsOptsBuilder interface {
	ToRouteTableRemoveSubnetsMap() (map[string]any, error)
}

// RemoveSubnetsOpts represents options for dissociating subnets from a route
// table.
type RemoveSubnetsOpts struct {
	// Subnets is the list of subnet IDs to dissociate.
	Subnets []string `json:"subnets" required:"true"`

	// Routes is an optional list of routes removed in the same operation,
	// identified by destination.
	Routes []RemoveRouteOpts `json:"routes,omitempty"`
}

// ToRouteTableRemoveSubnetsMap builds a request body from RemoveSubnetsOpts.
func (opts RemoveSubnetsOpts) ToRouteTableRemoveSubnetsMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// RemoveSubnets dissociates subnets from a route table.
func RemoveSubnets(ctx context.Context, c *gophercloud.ServiceClient, id string, opts RemoveSubnetsOptsBuilder) (r RemoveSubnetsResult) {
	b, err := opts.ToRouteTableRemoveSubnetsMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Put(ctx, removeSubnetsURL(c, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
