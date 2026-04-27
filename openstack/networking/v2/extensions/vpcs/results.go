package vpcs

import (
	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

type commonResult struct {
	gophercloud.Result
}

// Extract is a function that accepts a result and extracts a VPC resource.
func (r commonResult) Extract() (*Vpc, error) {
	var s Vpc
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonResult) ExtractInto(v any) error {
	return r.ExtractIntoStructPtr(v, "vpc")
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a Vpc.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Vpc.
type GetResult struct {
	commonResult
}

// UpdateResult represents the result of an update operation. Call its Extract
// method to interpret it as a Vpc.
type UpdateResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation. Call its
// ExtractErr method to determine if the request succeeded or failed.
type DeleteResult struct {
	gophercloud.ErrResult
}

// Vpc represents a Neutron Virtual Private Cloud.
type Vpc struct {
	// ID is the unique identifier of the VPC.
	ID string `json:"id"`

	// Name is the human-readable name of the VPC.
	Name string `json:"name"`

	// Description is the human-readable description of the VPC.
	Description string `json:"description"`

	// CIDRBlock is the RFC1918 range the VPC allocates its subnets from.
	CIDRBlock string `json:"cidr_block"`

	// VpcStatus is the sync state of the VPC with the backend: ACTIVE,
	// PENDING_CREATE, PENDING_UPDATE, PENDING_DELETE or ERROR.
	VpcStatus string `json:"vpc_status"`

	// TenantID is the project owner of the VPC.
	TenantID string `json:"tenant_id"`

	// ProjectID is the project owner of the VPC.
	ProjectID string `json:"project_id"`
}

// VpcPage is the page returned by a pager when traversing over a collection
// of VPCs.
type VpcPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of VPCs has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r VpcPage) NextPageURL() (string, error) {
	var s struct {
		Links []gophercloud.Link `json:"vpcs_links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gophercloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a VpcPage struct is empty.
func (r VpcPage) IsEmpty() (bool, error) {
	if r.StatusCode == 204 {
		return true, nil
	}

	is, err := ExtractVpcs(r)
	return len(is) == 0, err
}

// ExtractVpcs accepts a Page struct, specifically a VpcPage struct, and
// extracts the elements into a slice of Vpc structs. In other words, a
// generic collection is mapped into a relevant slice.
func ExtractVpcs(r pagination.Page) ([]Vpc, error) {
	var s []Vpc
	err := ExtractVpcsInto(r, &s)
	return s, err
}

// ExtractVpcsInto extracts the elements of a VpcPage into the provided slice.
func ExtractVpcsInto(r pagination.Page, v any) error {
	return r.(VpcPage).ExtractIntoSlicePtr(v, "vpcs")
}
