package ipsecsites

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

type Routing string

const (
	RoutingStatic Routing = "static"
	RoutingBGP    Routing = "bgp"
)

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToSiteCreateMap() (map[string]any, error)
}

// CreateOpts contains all the values needed to create a new IPsec site. A
// VPN gateway holds a single IPsec site.
type CreateOpts struct {
	// ProjectID specifies a project to own the IPsec site. The caller must be
	// allowed to act on that project in order to set this. Otherwise, the
	// caller's project will be the owner.
	ProjectID string `json:"project_id,omitempty"`

	// Name is the human readable name of the IPsec site.
	Name string `json:"name,omitempty"`

	// Description is the human readable description of the IPsec site.
	Description string `json:"description,omitempty"`

	// VPNGatewayID is the VPN gateway terminating the tunnels of the site.
	VPNGatewayID string `json:"vpn_gateway_id" required:"true"`

	// Routing is the routing scheme of the site. The default is static.
	Routing Routing `json:"routing,omitempty"`

	// LocalASN is the ASN of the VPN gateway. It is required with bgp
	// routing, and rejected otherwise.
	LocalASN int64 `json:"local_asn,omitempty"`

	// PeerASN is the ASN of the peer, different from LocalASN. It is
	// required with bgp routing, and rejected otherwise.
	PeerASN int64 `json:"peer_asn,omitempty"`

	// ECMPInbound advertises the networks equally on both tunnels, for an
	// ECMP-capable peer, instead of prepending the standby's advertisements.
	// It is only accepted with bgp routing.
	ECMPInbound *bool `json:"ecmp_inbound,omitempty"`

	// LocalEPGroupID is the endpoint group of the local networks.
	LocalEPGroupID string `json:"local_ep_group_id" required:"true"`

	// PeerEPGroupID is the endpoint group of the peer networks.
	PeerEPGroupID string `json:"peer_ep_group_id" required:"true"`
}

// ToSiteCreateMap casts a CreateOpts struct to a map.
func (opts CreateOpts) ToSiteCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "ipsec_site")
}

// Create accepts a CreateOpts struct and uses the values to create a new
// IPsec site.
func Create(ctx context.Context, c *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToSiteCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Post(ctx, rootURL(c), b, &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Get retrieves a particular IPsec site based on its unique ID.
func Get(ctx context.Context, c *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := c.Get(ctx, resourceURL(c, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete will permanently delete a particular IPsec site based on its unique
// ID. A site that still has connections cannot be deleted.
func Delete(ctx context.Context, c *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := c.Delete(ctx, resourceURL(c, id), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToSiteListQuery() (string, error)
}

// ListOpts allows the filtering, sorting and paging of IPsec sites. The
// service rejects any other query parameter.
type ListOpts struct {
	// ProjectID lists the sites of this project only. Naming a project other
	// than the caller's requires the get_all-global policy.
	ProjectID string `q:"project_id"`

	// Name lists the sites with this name only.
	Name string `q:"name"`

	// VPNGatewayID lists the sites of this VPN gateway only.
	VPNGatewayID string `q:"vpn_gateway_id"`

	// Limit is the maximum number of sites per page.
	Limit int `q:"limit"`

	// Marker is the ID of the last site of the previous page.
	Marker string `q:"marker"`

	// PageReverse pages backwards from Marker.
	PageReverse bool `q:"page_reverse"`

	// SortKey is the attribute to sort on.
	SortKey string `q:"sort_key"`

	// SortDir is the sort direction, asc or desc. It requires SortKey.
	SortDir string `q:"sort_dir"`

	// Fields restricts each returned site to these attributes.
	Fields []string `q:"fields"`
}

// ToSiteListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToSiteListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List returns a Pager which allows you to iterate over a collection of
// IPsec sites. It accepts a ListOpts struct, which allows you to filter
// the returned collection for greater efficiency.
func List(c *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := rootURL(c)
	if opts != nil {
		query, err := opts.ToSiteListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return SitePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToSiteUpdateMap() (map[string]any, error)
}

// UpdateOpts contains the values used when updating an IPsec site. The
// gateway, routing and networks of a site are fixed at creation.
type UpdateOpts struct {
	// Name is the human readable name of the IPsec site.
	Name *string `json:"name,omitempty"`

	// Description is the human readable description of the IPsec site.
	Description *string `json:"description,omitempty"`
}

// ToSiteUpdateMap casts an UpdateOpts struct to a map.
func (opts UpdateOpts) ToSiteUpdateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "ipsec_site")
}

// Update allows IPsec sites to be updated.
func Update(ctx context.Context, c *gophercloud.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToSiteUpdateMap()
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
