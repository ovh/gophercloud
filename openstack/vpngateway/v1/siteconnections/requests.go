package siteconnections

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

type AuthMode string
type Initiator string
type DPDAction string

const (
	StatusPendingCreate = "PENDING_CREATE"
	StatusPendingUpdate = "PENDING_UPDATE"
	StatusPendingDelete = "PENDING_DELETE"
	StatusActive        = "ACTIVE"
	StatusError         = "ERROR"

	AuthModePSK            AuthMode  = "psk"
	InitiatorBiDirectional Initiator = "bi-directional"
	InitiatorResponseOnly  Initiator = "response-only"
	DPDActionHold          DPDAction = "hold"
	DPDActionClear         DPDAction = "clear"
	DPDActionDisabled      DPDAction = "disabled"
	DPDActionRestart       DPDAction = "restart"
	DPDActionRestartByPeer DPDAction = "restart-by-peer"
)

// DPDOpts contains the dead peer detection settings of a connection. A
// parameter left unset keeps its current value, or its default on create.
type DPDOpts struct {
	// Action is the action to take when the peer is dead.
	// The default is restart.
	Action DPDAction `json:"action,omitempty"`

	// Interval is the idle time in seconds before a liveness check.
	// The default is 30.
	Interval int `json:"interval,omitempty"`

	// Timeout is the time in seconds before the peer is declared dead. It
	// applies to IKEv1 only. The default is 150.
	Timeout int `json:"timeout,omitempty"`
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToConnectionCreateMap() (map[string]any, error)
}

// CreateOpts contains all the values needed to create a new IPsec site
// connection: one tunnel to an IPsec site, at most two per site.
type CreateOpts struct {
	// ProjectID must be the project of the IPsec site when set. The
	// connection is always owned by the project of its site.
	ProjectID string `json:"project_id,omitempty"`

	// Name is the human readable name of the connection.
	Name string `json:"name,omitempty"`

	// Description is the human readable description of the connection.
	Description string `json:"description,omitempty"`

	// IPsecSiteID is the IPsec site of the connection.
	IPsecSiteID string `json:"ipsec_site_id" required:"true"`

	// PeerAddress is the IP address of the peer.
	PeerAddress string `json:"peer_address" required:"true"`

	// PeerID is the IKE identity of the peer: an IP address, an FQDN,
	// user@fqdn or a distinguished name. The default is PeerAddress.
	PeerID string `json:"peer_id,omitempty"`

	// LocalID is the local IKE identity. The default is the external
	// address of the VPN gateway.
	LocalID string `json:"local_id,omitempty"`

	// AuthMode is the authentication mode. Only psk is accepted.
	AuthMode AuthMode `json:"auth_mode,omitempty"`

	// PSK is the pre-shared key, 8 to 255 printable ASCII characters. It is
	// never returned.
	PSK string `json:"psk" required:"true"`

	// Priority picks the preferred tunnel of the site: lower is preferred,
	// unique per site. The default is one more than the highest in use.
	Priority int `json:"priority,omitempty"`

	// AdminStateUp is the administrative state of the connection.
	AdminStateUp *bool `json:"admin_state_up,omitempty"`

	// IKEPolicyID is the IKE policy of the connection.
	IKEPolicyID string `json:"ikepolicy_id" required:"true"`

	// IPsecPolicyID is the IPsec policy of the connection.
	IPsecPolicyID string `json:"ipsecpolicy_id" required:"true"`

	// DPD is the dead peer detection settings of the connection.
	DPD *DPDOpts `json:"dpd,omitempty"`

	// MTU is the MTU of the tunnel, at least 68. The default is 1400.
	MTU int `json:"mtu,omitempty"`

	// Initiator tells whether the gateway initiates the connection.
	// The default is bi-directional.
	Initiator Initiator `json:"initiator,omitempty"`

	// TunnelAddress is the address of the gateway inside the tunnel, with
	// its prefix (e.g. 169.254.10.1/30). It is required on a bgp site, and
	// rejected on a static one.
	TunnelAddress string `json:"tunnel_address,omitempty"`

	// BGPPeerAddress is the address of the peer inside the tunnel network.
	// It is required on a bgp site, and rejected on a static one.
	BGPPeerAddress string `json:"bgp_peer_address,omitempty"`
}

// ToConnectionCreateMap casts a CreateOpts struct to a map.
func (opts CreateOpts) ToConnectionCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "ipsec_site_connection")
}

// Create accepts a CreateOpts struct and uses the values to create a new
// IPsec site connection. The connection is PENDING_CREATE until the tunnel
// is provisioned on the VPN gateway.
func Create(ctx context.Context, c *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToConnectionCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Post(ctx, rootURL(c), b, &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Get retrieves a particular IPsec site connection based on its unique ID.
func Get(ctx context.Context, c *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := c.Get(ctx, resourceURL(c, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete will delete a particular IPsec site connection based on its unique
// ID. The connection is PENDING_DELETE until its tunnel is torn down.
func Delete(ctx context.Context, c *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := c.Delete(ctx, resourceURL(c, id), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToConnectionListQuery() (string, error)
}

// ListOpts allows the filtering, sorting and paging of IPsec site
// connections. The service rejects any other query parameter.
type ListOpts struct {
	// ProjectID lists the connections of this project only. Naming a project
	// other than the caller's requires the get_all-global policy.
	ProjectID string `q:"project_id"`

	// Name lists the connections with this name only.
	Name string `q:"name"`

	// Status lists the connections with this status only.
	Status string `q:"status"`

	// IPsecSiteID lists the connections of this IPsec site only.
	IPsecSiteID string `q:"ipsec_site_id"`

	// Limit is the maximum number of connections per page.
	Limit int `q:"limit"`

	// Marker is the ID of the last connection of the previous page.
	Marker string `q:"marker"`

	// PageReverse pages backwards from Marker.
	PageReverse bool `q:"page_reverse"`

	// SortKey is the attribute to sort on.
	SortKey string `q:"sort_key"`

	// SortDir is the sort direction, asc or desc. It requires SortKey.
	SortDir string `q:"sort_dir"`

	// Fields restricts each returned connection to these attributes.
	Fields []string `q:"fields"`
}

// ToConnectionListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToConnectionListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List returns a Pager which allows you to iterate over a collection of
// IPsec site connections. It accepts a ListOpts struct, which allows you to
// filter the returned collection for greater efficiency.
func List(c *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := rootURL(c)
	if opts != nil {
		query, err := opts.ToConnectionListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return ConnectionPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToConnectionUpdateMap() (map[string]any, error)
}

// UpdateOpts contains the values used when updating an IPsec site
// connection. The site of a connection is fixed at creation. An update
// re-provisions the tunnel: the connection is PENDING_UPDATE until it is done.
type UpdateOpts struct {
	// Name is the human readable name of the connection.
	Name *string `json:"name,omitempty"`

	// Description is the human readable description of the connection.
	Description *string `json:"description,omitempty"`

	// PeerAddress is the IP address of the peer.
	PeerAddress string `json:"peer_address,omitempty"`

	// PeerID is the IKE identity of the peer.
	PeerID string `json:"peer_id,omitempty"`

	// LocalID is the local IKE identity.
	LocalID string `json:"local_id,omitempty"`

	// PSK replaces the pre-shared key.
	PSK string `json:"psk,omitempty"`

	// Priority is the priority of the tunnel within its site.
	Priority int `json:"priority,omitempty"`

	// AdminStateUp is the administrative state of the connection.
	AdminStateUp *bool `json:"admin_state_up,omitempty"`

	// IKEPolicyID is the IKE policy of the connection.
	IKEPolicyID string `json:"ikepolicy_id,omitempty"`

	// IPsecPolicyID is the IPsec policy of the connection.
	IPsecPolicyID string `json:"ipsecpolicy_id,omitempty"`

	// DPD changes the dead peer detection parameters it sets.
	DPD *DPDOpts `json:"dpd,omitempty"`

	// MTU is the MTU of the tunnel.
	MTU int `json:"mtu,omitempty"`

	// Initiator tells whether the gateway initiates the connection.
	Initiator Initiator `json:"initiator,omitempty"`

	// TunnelAddress is the address of the gateway inside the tunnel. It is
	// only accepted on a bgp site.
	TunnelAddress string `json:"tunnel_address,omitempty"`

	// BGPPeerAddress is the address of the peer inside the tunnel network.
	// It is only accepted on a bgp site.
	BGPPeerAddress string `json:"bgp_peer_address,omitempty"`
}

// ToConnectionUpdateMap casts an UpdateOpts struct to a map.
func (opts UpdateOpts) ToConnectionUpdateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "ipsec_site_connection")
}

// Update allows IPsec site connections to be updated.
func Update(ctx context.Context, c *gophercloud.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToConnectionUpdateMap()
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
