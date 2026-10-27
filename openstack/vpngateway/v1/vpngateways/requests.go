package vpngateways

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

const (
	StatusPendingCreate = "PENDING_CREATE"
	StatusPendingUpdate = "PENDING_UPDATE"
	StatusPendingDelete = "PENDING_DELETE"
	StatusActive        = "ACTIVE"
	StatusError         = "ERROR"
)

// RouteOpts is a static route of a VPN gateway.
type RouteOpts struct {
	// DestinationCIDR is the destination network of the route.
	DestinationCIDR string `json:"destination" required:"true"`

	// NextHop is the IP address the route forwards to.
	NextHop string `json:"nexthop" required:"true"`
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToVPNGatewayCreateMap() (map[string]any, error)
}

// CreateOpts contains all the values needed to create a new VPN gateway.
type CreateOpts struct {
	// ProjectID specifies a project to own the VPN gateway. The caller must be
	// allowed to act on that project in order to set this. Otherwise, the
	// caller's project will be the owner.
	ProjectID string `json:"project_id,omitempty"`

	// Name is the human readable name of the VPN gateway, unique in its
	// project.
	Name string `json:"name,omitempty"`

	// Description is the human readable description of the VPN gateway.
	Description string `json:"description,omitempty"`

	// AdminStateUp is the administrative state of the VPN gateway.
	AdminStateUp *bool `json:"admin_state_up,omitempty"`

	// HA asks for a highly available VPN gateway.
	HA *bool `json:"ha,omitempty"`

	// FlavorID is the ID of the flavor of the VPN gateway.
	FlavorID string `json:"flavor_id,omitempty"`

	// AvailabilityZoneHints is the list of availability zones the VPN gateway
	// should be scheduled to.
	AvailabilityZoneHints []string `json:"availability_zone_hints,omitempty"`

	// Routes is the list of static routes of the VPN gateway.
	Routes []RouteOpts `json:"routes,omitempty"`
}

// ToVPNGatewayCreateMap casts a CreateOpts struct to a map.
func (opts CreateOpts) ToVPNGatewayCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "vpn_gateway")
}

// Create accepts a CreateOpts struct and uses the values to create a new
// VPN gateway. The gateway starts in PENDING_CREATE until its instance is
// booted.
func Create(ctx context.Context, c *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToVPNGatewayCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Post(ctx, rootURL(c), b, &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Get retrieves a particular VPN gateway based on its unique ID.
func Get(ctx context.Context, c *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := c.Get(ctx, resourceURL(c, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete will delete a particular VPN gateway based on its unique ID. The
// gateway moves to PENDING_DELETE and disappears once it is torn down.
func Delete(ctx context.Context, c *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := c.Delete(ctx, resourceURL(c, id), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToVPNGatewayListQuery() (string, error)
}

// ListOpts allows the filtering, sorting and paging of VPN gateways. The
// service rejects any other query parameter.
type ListOpts struct {
	// ProjectID lists the VPN gateways of this project only. Naming a project
	// other than the caller's requires the get_all-global policy.
	ProjectID string `q:"project_id"`

	// Name lists the VPN gateways with this name only.
	Name string `q:"name"`

	// Status lists the VPN gateways with this status only.
	Status string `q:"status"`

	// FlavorID lists the VPN gateways of this flavor only.
	FlavorID string `q:"flavor_id"`

	// Limit is the maximum number of VPN gateways per page.
	Limit int `q:"limit"`

	// Marker is the ID of the last VPN gateway of the previous page.
	Marker string `q:"marker"`

	// PageReverse pages backwards from Marker.
	PageReverse bool `q:"page_reverse"`

	// SortKey is the attribute to sort on.
	SortKey string `q:"sort_key"`

	// SortDir is the sort direction, asc or desc. It requires SortKey.
	SortDir string `q:"sort_dir"`

	// Fields restricts each returned VPN gateway to these attributes.
	Fields []string `q:"fields"`
}

// ToVPNGatewayListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToVPNGatewayListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List returns a Pager which allows you to iterate over a collection of
// VPN gateways. It accepts a ListOpts struct, which allows you to filter
// the returned collection for greater efficiency.
func List(c *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := rootURL(c)
	if opts != nil {
		query, err := opts.ToVPNGatewayListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return VPNGatewayPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToVPNGatewayUpdateMap() (map[string]any, error)
}

// UpdateOpts contains the values used when updating a VPN gateway. The
// gateway must not be PENDING_*.
type UpdateOpts struct {
	// Name is the human readable name of the VPN gateway.
	Name *string `json:"name,omitempty"`

	// Description is the human readable description of the VPN gateway.
	Description *string `json:"description,omitempty"`

	// AdminStateUp is the administrative state of the VPN gateway.
	AdminStateUp *bool `json:"admin_state_up,omitempty"`

	// HA asks for a highly available VPN gateway.
	HA *bool `json:"ha,omitempty"`

	// FlavorID is the ID of the flavor of the VPN gateway.
	FlavorID *string `json:"flavor_id,omitempty"`

	// AvailabilityZoneHints replaces the availability zone hints.
	AvailabilityZoneHints *[]string `json:"availability_zone_hints,omitempty"`

	// Routes replaces the static routes. An empty slice removes them all.
	Routes *[]RouteOpts `json:"routes,omitempty"`
}

// ToVPNGatewayUpdateMap casts an UpdateOpts struct to a map.
func (opts UpdateOpts) ToVPNGatewayUpdateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "vpn_gateway")
}

// Update allows VPN gateways to be updated.
func Update(ctx context.Context, c *gophercloud.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToVPNGatewayUpdateMap()
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

// AddInterfaceOptsBuilder allows extensions to add additional parameters to
// the AddInterface request.
type AddInterfaceOptsBuilder interface {
	ToVPNGatewayAddInterfaceMap() (map[string]any, error)
}

// AddInterfaceOpts represents the options for adding an internal interface
// to a VPN gateway.
type AddInterfaceOpts struct {
	// SubnetID is the subnet the interface is attached to.
	SubnetID string `json:"subnet_id" required:"true"`
}

// ToVPNGatewayAddInterfaceMap builds a request body from AddInterfaceOpts.
func (opts AddInterfaceOpts) ToVPNGatewayAddInterfaceMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// AddInterface attaches an internal subnet to a VPN gateway. The port is
// created asynchronously: the interface shows in the gateway's
// InternalInterfaces, PENDING_CREATE until it is ACTIVE.
func AddInterface(ctx context.Context, c *gophercloud.ServiceClient, id string, opts AddInterfaceOptsBuilder) (r InterfaceResult) {
	b, err := opts.ToVPNGatewayAddInterfaceMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Put(ctx, addInterfaceURL(c, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// RemoveInterfaceOptsBuilder allows extensions to add additional parameters
// to the RemoveInterface request.
type RemoveInterfaceOptsBuilder interface {
	ToVPNGatewayRemoveInterfaceMap() (map[string]any, error)
}

// RemoveInterfaceOpts represents the options for removing an internal
// interface from a VPN gateway.
type RemoveInterfaceOpts struct {
	// SubnetID is the subnet the interface is attached to.
	SubnetID string `json:"subnet_id" required:"true"`
}

// ToVPNGatewayRemoveInterfaceMap builds a request body from
// RemoveInterfaceOpts.
func (opts RemoveInterfaceOpts) ToVPNGatewayRemoveInterfaceMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// RemoveInterface detaches an internal subnet from a VPN gateway. The
// interface is PENDING_DELETE until its port is removed.
func RemoveInterface(ctx context.Context, c *gophercloud.ServiceClient, id string, opts RemoveInterfaceOptsBuilder) (r InterfaceResult) {
	b, err := opts.ToVPNGatewayRemoveInterfaceMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Put(ctx, removeInterfaceURL(c, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ExternalGatewayOpts names a subnet of the external network to attach.
type ExternalGatewayOpts struct {
	// SubnetID is a subnet of the external network.
	SubnetID string `json:"subnet_id" required:"true"`
}

// AddExternalGatewaysOptsBuilder allows extensions to add additional
// parameters to the AddExternalGateways request.
type AddExternalGatewaysOptsBuilder interface {
	ToVPNGatewayAddExternalGatewaysMap() (map[string]any, error)
}

// AddExternalGatewaysOpts represents the options for attaching a VPN gateway
// to an external network.
type AddExternalGatewaysOpts struct {
	// ExternalGateways names one or two subnets, at most one IPv4 and one
	// IPv6, all on the same external network.
	ExternalGateways []ExternalGatewayOpts `json:"external_gateways" required:"true"`
}

// ToVPNGatewayAddExternalGatewaysMap builds a request body from
// AddExternalGatewaysOpts.
func (opts AddExternalGatewaysOpts) ToVPNGatewayAddExternalGatewaysMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "router")
}

// AddExternalGateways attaches a VPN gateway to an external network. A
// gateway holds a single external network. The port is created
// asynchronously: the external gateway is PENDING_CREATE until it is ACTIVE.
func AddExternalGateways(ctx context.Context, c *gophercloud.ServiceClient, id string, opts AddExternalGatewaysOptsBuilder) (r AddExternalGatewaysResult) {
	b, err := opts.ToVPNGatewayAddExternalGatewaysMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Put(ctx, addExternalGatewaysURL(c, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// RebuildOptsBuilder allows extensions to add additional parameters to the
// Rebuild request.
type RebuildOptsBuilder interface {
	ToVPNGatewayRebuildMap() (map[string]any, error)
}

// RebuildOpts represents the options for rebuilding a VPN gateway.
type RebuildOpts struct {
	// ImageID is the Glance image to rebuild with. Without it, the service
	// picks its newest gateway image.
	ImageID string `json:"image_id,omitempty"`
}

// ToVPNGatewayRebuildMap builds a request body from RebuildOpts.
func (opts RebuildOpts) ToVPNGatewayRebuildMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// Rebuild rebuilds the instances of a VPN gateway on an image, then pushes
// its configuration again. It is an admin action by default. The gateway is
// PENDING_UPDATE until the rebuild ends.
func Rebuild(ctx context.Context, c *gophercloud.ServiceClient, id string, opts RebuildOptsBuilder) (r ActionResult) {
	b, err := opts.ToVPNGatewayRebuildMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Put(ctx, rebuildURL(c, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Resync pushes the configuration of a VPN gateway to its instances again.
// It is an admin action by default. The gateway is PENDING_UPDATE until the
// resync ends.
func Resync(ctx context.Context, c *gophercloud.ServiceClient, id string) (r ActionResult) {
	resp, err := c.Put(ctx, resyncURL(c, id), nil, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
