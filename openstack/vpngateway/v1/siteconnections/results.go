package siteconnections

import (
	"encoding/json"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// DPD is the dead peer detection settings of a connection.
type DPD struct {
	// Action is the action taken when the peer is dead.
	Action string `json:"action"`

	// Interval is the idle time in seconds before a liveness check.
	Interval int `json:"interval"`

	// Timeout is the time in seconds before the peer is declared dead.
	Timeout int `json:"timeout"`
}

// Connection is an IPsec site connection: one tunnel to an IPsec site.
type Connection struct {
	// ID is the ID of the connection.
	ID string `json:"id"`

	// ProjectID is the ID of the project owning the connection.
	ProjectID string `json:"project_id"`

	// Name is the human readable name of the connection.
	Name string `json:"name"`

	// Description is the human readable description of the connection.
	Description string `json:"description"`

	// IPsecSiteID is the IPsec site of the connection.
	IPsecSiteID string `json:"ipsec_site_id"`

	// PeerAddress is the IP address of the peer.
	PeerAddress string `json:"peer_address"`

	// PeerID is the IKE identity of the peer.
	PeerID string `json:"peer_id"`

	// LocalID is the local IKE identity. It is empty when the gateway uses
	// its external address.
	LocalID string `json:"local_id"`

	// AuthMode is the authentication mode.
	AuthMode string `json:"auth_mode"`

	// PSKSet tells whether a pre-shared key is set. The key itself is never
	// returned.
	PSKSet bool `json:"psk_set"`

	// Priority is the priority of the tunnel within its site, lower is
	// preferred.
	Priority int `json:"priority"`

	// IKEPolicyID is the IKE policy of the connection.
	IKEPolicyID string `json:"ikepolicy_id"`

	// IPsecPolicyID is the IPsec policy of the connection.
	IPsecPolicyID string `json:"ipsecpolicy_id"`

	// DPD is the dead peer detection settings of the connection.
	DPD DPD `json:"dpd"`

	// MTU is the MTU of the tunnel.
	MTU int `json:"mtu"`

	// Initiator tells whether the gateway initiates the connection.
	Initiator string `json:"initiator"`

	// TunnelAddress is the address of the gateway inside the tunnel, on a
	// bgp site.
	TunnelAddress string `json:"tunnel_address"`

	// BGPPeerAddress is the address of the peer inside the tunnel network,
	// on a bgp site.
	BGPPeerAddress string `json:"bgp_peer_address"`

	// AdminStateUp is the administrative state of the connection.
	AdminStateUp bool `json:"admin_state_up"`

	// Status is the provisioning status of the connection. ACTIVE means the
	// tunnel is configured on the gateway, not that it is up.
	Status string `json:"status"`

	// CreatedAt is the time the connection was created.
	CreatedAt time.Time `json:"-"`

	// UpdatedAt is the time the connection was last updated.
	UpdatedAt time.Time `json:"-"`
}

func (r *Connection) UnmarshalJSON(b []byte) error {
	type tmp Connection
	var s struct {
		tmp
		CreatedAt gophercloud.JSONRFC3339MilliNoZ `json:"created_at"`
		UpdatedAt gophercloud.JSONRFC3339MilliNoZ `json:"updated_at"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*r = Connection(s.tmp)
	r.CreatedAt = time.Time(s.CreatedAt)
	r.UpdatedAt = time.Time(s.UpdatedAt)
	return nil
}

type commonResult struct {
	gophercloud.Result
}

// Extract is a function that accepts a result and extracts an IPsec site
// connection.
func (r commonResult) Extract() (*Connection, error) {
	var s struct {
		Connection *Connection `json:"ipsec_site_connection"`
	}
	err := r.ExtractInto(&s)
	return s.Connection, err
}

// ConnectionPage is the page returned by a pager when traversing over a
// collection of IPsec site connections.
type ConnectionPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of IPsec site
// connections has reached the end of a page and the pager seeks to traverse
// over a new one. In order to do this, it needs to construct the next page's
// URL.
func (r ConnectionPage) NextPageURL() (string, error) {
	var s struct {
		Links []gophercloud.Link `json:"ipsec_site_connections_links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gophercloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a ConnectionPage struct is empty.
func (r ConnectionPage) IsEmpty() (bool, error) {
	if r.StatusCode == 204 {
		return true, nil
	}

	is, err := ExtractConnections(r)
	return len(is) == 0, err
}

// ExtractConnections accepts a Page struct, specifically a ConnectionPage
// struct, and extracts the elements into a slice of Connection structs.
func ExtractConnections(r pagination.Page) ([]Connection, error) {
	var s struct {
		Connections []Connection `json:"ipsec_site_connections"`
	}
	err := (r.(ConnectionPage)).ExtractInto(&s)
	return s.Connections, err
}

// CreateResult represents the result of a Create operation. Call its Extract method to
// interpret it as a Connection.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a Get operation. Call its Extract method to
// interpret it as a Connection.
type GetResult struct {
	commonResult
}

// DeleteResult represents the results of a Delete operation. Call its ExtractErr method
// to determine whether the operation succeeded or failed.
type DeleteResult struct {
	gophercloud.ErrResult
}

// UpdateResult represents the result of an update operation. Call its Extract method
// to interpret it as a Connection.
type UpdateResult struct {
	commonResult
}
