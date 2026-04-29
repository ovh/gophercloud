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

// ListOpts allows the filtering of paginated collections through the API.
type ListOpts struct {
	RouterID  string `q:"router_id"`
	IsDefault *bool  `q:"is_default"`
	Limit     int    `q:"limit"`
	Marker    string `q:"marker"`
}

// ToRouteTableListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToRouteTableListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), nil
}

// List returns a Pager which allows you to iterate over a collection of
// route tables.
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

	// RouterID is the ID of the router this route table belongs to. Required.
	RouterID string `json:"router_id" required:"true"`

	// IsDefault indicates whether this is the default route table.
	IsDefault *bool `json:"is_default,omitempty"`
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

// Delete accepts a route table ID and deletes it.
func Delete(ctx context.Context, c *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := c.Delete(ctx, resourceURL(c, id), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// AddRoutesOptsBuilder allows extensions to add additional parameters to the
// AddRoutes request.
type AddRoutesOptsBuilder interface {
	ToAddRoutesMap() (map[string]any, error)
}

// RouteOpts represents a route to add to a route table.
type RouteOpts struct {
	// Destination is the destination CIDR of the route. Required.
	Destination string `json:"destination" required:"true"`

	// Nexthop is the IP address of the next hop.
	Nexthop string `json:"nexthop,omitempty"`

	// NexthopType is the type of the nexthop. Required.
	NexthopType string `json:"nexthop_type" required:"true"`
}

// AddRoutesOpts represents options for adding routes to a route table.
type AddRoutesOpts struct {
	Routes []RouteOpts `json:"routes" required:"true"`
}

// ToAddRoutesMap builds a request body from AddRoutesOpts.
func (opts AddRoutesOpts) ToAddRoutesMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// AddRoutes adds routes to a route table.
func AddRoutes(ctx context.Context, c *gophercloud.ServiceClient, id string, opts AddRoutesOptsBuilder) (r AddRoutesResult) {
	b, err := opts.ToAddRoutesMap()
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

// RemoveRoutesOptsBuilder allows extensions to add additional parameters to the
// RemoveRoutes request.
type RemoveRoutesOptsBuilder interface {
	ToRemoveRoutesMap() (map[string]any, error)
}

// RemoveRouteOpts represents a route to remove from a route table.
type RemoveRouteOpts struct {
	// Destination is the destination CIDR of the route to remove. Required.
	Destination string `json:"destination" required:"true"`
}

// RemoveRoutesOpts represents options for removing routes from a route table.
type RemoveRoutesOpts struct {
	Routes []RemoveRouteOpts `json:"routes" required:"true"`
}

// ToRemoveRoutesMap builds a request body from RemoveRoutesOpts.
func (opts RemoveRoutesOpts) ToRemoveRoutesMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// RemoveRoutes removes routes from a route table.
func RemoveRoutes(ctx context.Context, c *gophercloud.ServiceClient, id string, opts RemoveRoutesOptsBuilder) (r RemoveRoutesResult) {
	b, err := opts.ToRemoveRoutesMap()
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
	ToAddSubnetsMap() (map[string]any, error)
}

// AddSubnetsOpts represents options for adding subnets to a route table.
type AddSubnetsOpts struct {
	Subnets []string `json:"subnets" required:"true"`
}

// ToAddSubnetsMap builds a request body from AddSubnetsOpts.
func (opts AddSubnetsOpts) ToAddSubnetsMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// AddSubnets adds subnet associations to a route table.
func AddSubnets(ctx context.Context, c *gophercloud.ServiceClient, id string, opts AddSubnetsOptsBuilder) (r AddSubnetsResult) {
	b, err := opts.ToAddSubnetsMap()
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

// RemoveSubnetsOptsBuilder allows extensions to add additional parameters to the
// RemoveSubnets request.
type RemoveSubnetsOptsBuilder interface {
	ToRemoveSubnetsMap() (map[string]any, error)
}

// RemoveSubnetsOpts represents options for removing subnets from a route table.
type RemoveSubnetsOpts struct {
	Subnets []string `json:"subnets" required:"true"`
}

// ToRemoveSubnetsMap builds a request body from RemoveSubnetsOpts.
func (opts RemoveSubnetsOpts) ToRemoveSubnetsMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// RemoveSubnets removes subnet associations from a route table.
func RemoveSubnets(ctx context.Context, c *gophercloud.ServiceClient, id string, opts RemoveSubnetsOptsBuilder) (r RemoveSubnetsResult) {
	b, err := opts.ToRemoveSubnetsMap()
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
