package ipsecsites

import (
	"encoding/json"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// Site is an IPsec site: what every tunnel to one remote site shares.
type Site struct {
	// ID is the ID of the IPsec site.
	ID string `json:"id"`

	// ProjectID is the ID of the project owning the IPsec site.
	ProjectID string `json:"project_id"`

	// Name is the human readable name of the IPsec site.
	Name string `json:"name"`

	// Description is the human readable description of the IPsec site.
	Description string `json:"description"`

	// VPNGatewayID is the VPN gateway terminating the tunnels of the site.
	VPNGatewayID string `json:"vpn_gateway_id"`

	// Routing is the routing scheme of the site, static or bgp.
	Routing string `json:"routing"`

	// LocalASN is the ASN of the VPN gateway. It is nil with static routing.
	LocalASN *int64 `json:"local_asn"`

	// PeerASN is the ASN of the peer. It is nil with static routing.
	PeerASN *int64 `json:"peer_asn"`

	// ECMPInbound tells whether the networks are advertised equally on both
	// tunnels.
	ECMPInbound bool `json:"ecmp_inbound"`

	// LocalEPGroupID is the endpoint group of the local networks.
	LocalEPGroupID string `json:"local_ep_group_id"`

	// PeerEPGroupID is the endpoint group of the peer networks.
	PeerEPGroupID string `json:"peer_ep_group_id"`

	// CreatedAt is the time the IPsec site was created.
	CreatedAt time.Time `json:"-"`

	// UpdatedAt is the time the IPsec site was last updated.
	UpdatedAt time.Time `json:"-"`
}

func (r *Site) UnmarshalJSON(b []byte) error {
	type tmp Site
	var s struct {
		tmp
		CreatedAt gophercloud.JSONRFC3339MilliNoZ `json:"created_at"`
		UpdatedAt gophercloud.JSONRFC3339MilliNoZ `json:"updated_at"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*r = Site(s.tmp)
	r.CreatedAt = time.Time(s.CreatedAt)
	r.UpdatedAt = time.Time(s.UpdatedAt)
	return nil
}

type commonResult struct {
	gophercloud.Result
}

// Extract is a function that accepts a result and extracts an IPsec site.
func (r commonResult) Extract() (*Site, error) {
	var s struct {
		Site *Site `json:"ipsec_site"`
	}
	err := r.ExtractInto(&s)
	return s.Site, err
}

// SitePage is the page returned by a pager when traversing over a
// collection of IPsec sites.
type SitePage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of IPsec sites has
// reached the end of a page and the pager seeks to traverse over a new one.
// In order to do this, it needs to construct the next page's URL.
func (r SitePage) NextPageURL() (string, error) {
	var s struct {
		Links []gophercloud.Link `json:"ipsec_sites_links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gophercloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a SitePage struct is empty.
func (r SitePage) IsEmpty() (bool, error) {
	if r.StatusCode == 204 {
		return true, nil
	}

	is, err := ExtractSites(r)
	return len(is) == 0, err
}

// ExtractSites accepts a Page struct, specifically a SitePage struct, and
// extracts the elements into a slice of Site structs.
func ExtractSites(r pagination.Page) ([]Site, error) {
	var s struct {
		Sites []Site `json:"ipsec_sites"`
	}
	err := (r.(SitePage)).ExtractInto(&s)
	return s.Sites, err
}

// CreateResult represents the result of a Create operation. Call its Extract method to
// interpret it as a Site.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a Get operation. Call its Extract method to
// interpret it as a Site.
type GetResult struct {
	commonResult
}

// DeleteResult represents the results of a Delete operation. Call its ExtractErr method
// to determine whether the operation succeeded or failed.
type DeleteResult struct {
	gophercloud.ErrResult
}

// UpdateResult represents the result of an update operation. Call its Extract method
// to interpret it as a Site.
type UpdateResult struct {
	commonResult
}
