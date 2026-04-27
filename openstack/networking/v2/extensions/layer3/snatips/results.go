package snatips

import (
	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

type commonResult struct {
	gophercloud.Result
}

// Extract is a function that accepts a result and extracts a SNAT IP resource.
func (r commonResult) Extract() (*SnatIP, error) {
	var s SnatIP
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonResult) ExtractInto(v any) error {
	return r.ExtractIntoStructPtr(v, "snat_ip")
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a SnatIP.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a SnatIP.
type GetResult struct {
	commonResult
}

// UpdateResult represents the result of an update operation. Call its Extract
// method to interpret it as a SnatIP.
type UpdateResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation. Call its
// ExtractErr method to determine if the request succeeded or failed.
type DeleteResult struct {
	gophercloud.ErrResult
}

// SnatIP represents a public address reserved for the gateway port of a VPC
// router.
type SnatIP struct {
	// ID is the unique identifier of the SNAT IP.
	ID string `json:"id"`

	// NetworkID is the ID of the external network the address is reserved on.
	NetworkID string `json:"network_id"`

	// AvailabilityZone is the availability zone the address is reserved in.
	AvailabilityZone string `json:"availability_zone"`

	// IPAddress is the reserved address.
	IPAddress string `json:"ip_address"`

	// Description is the human-readable description of the SNAT IP.
	Description string `json:"description"`

	// VpcID is the VPC of the router whose gateway port carries the address.
	// It is empty while no gateway carries it.
	VpcID string `json:"vpc_id"`

	// Status is ACTIVE once a router gateway port carries the address, and
	// DOWN while none does.
	Status string `json:"status"`

	// VpcStatus is the sync state of the reservation with the backend:
	// ACTIVE, PENDING_CREATE, PENDING_UPDATE, PENDING_DELETE or ERROR.
	VpcStatus string `json:"vpc_status"`

	// TenantID is the project owner of the SNAT IP.
	TenantID string `json:"tenant_id"`

	// ProjectID is the project owner of the SNAT IP.
	ProjectID string `json:"project_id"`
}

// SnatIPPage is the page returned by a pager when traversing over a
// collection of SNAT IPs.
type SnatIPPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of SNAT IPs has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r SnatIPPage) NextPageURL() (string, error) {
	var s struct {
		Links []gophercloud.Link `json:"snat_ips_links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gophercloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a SnatIPPage struct is empty.
func (r SnatIPPage) IsEmpty() (bool, error) {
	if r.StatusCode == 204 {
		return true, nil
	}

	is, err := ExtractSnatIPs(r)
	return len(is) == 0, err
}

// ExtractSnatIPs accepts a Page struct, specifically a SnatIPPage struct,
// and extracts the elements into a slice of SnatIP structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractSnatIPs(r pagination.Page) ([]SnatIP, error) {
	var s []SnatIP
	err := ExtractSnatIPsInto(r, &s)
	return s, err
}

// ExtractSnatIPsInto extracts the elements of a SnatIPPage into the provided
// slice.
func ExtractSnatIPsInto(r pagination.Page, v any) error {
	return r.(SnatIPPage).ExtractIntoSlicePtr(v, "snat_ips")
}
