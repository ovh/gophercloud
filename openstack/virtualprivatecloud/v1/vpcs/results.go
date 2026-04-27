package vpcs

import (
	"encoding/json"
	"time"

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

// DeleteResult represents the result of a delete operation. Call its Extract
// method to interpret it as a Vpc (Orion returns 202 with the VPC in DELETING
// state).
type DeleteResult struct {
	commonResult
}

// Vpc represents an Orion Virtual Private Cloud resource.
type Vpc struct {
	// UUID for the VPC.
	ID string `json:"id"`

	// OpenStack project ID (from auth token).
	ProjectID string `json:"project_id"`

	// Human-readable VPC name.
	Name string `json:"name"`

	// Optional description.
	Description string `json:"description"`

	// IPv4 CIDR block.
	CIDRBlock string `json:"cidr_block"`

	// VPC status: CREATING, READY, UPDATING, ERROR, DELETING, DELETED.
	Status string `json:"status"`

	// Creation timestamp.
	CreatedAt time.Time `json:"-"`

	// Last update timestamp.
	UpdatedAt time.Time `json:"-"`
}

func (r *Vpc) UnmarshalJSON(b []byte) error {
	type tmp Vpc

	// Support for time format without timezone suffix
	var s1 struct {
		tmp
		CreatedAt gophercloud.JSONRFC3339NoZ `json:"created_at"`
		UpdatedAt gophercloud.JSONRFC3339NoZ `json:"updated_at"`
	}

	err := json.Unmarshal(b, &s1)
	if err == nil {
		*r = Vpc(s1.tmp)
		r.CreatedAt = time.Time(s1.CreatedAt)
		r.UpdatedAt = time.Time(s1.UpdatedAt)

		return nil
	}

	// Support for standard RFC3339 time format
	var s2 struct {
		tmp
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	err = json.Unmarshal(b, &s2)
	if err != nil {
		return err
	}

	*r = Vpc(s2.tmp)
	r.CreatedAt = time.Time(s2.CreatedAt)
	r.UpdatedAt = time.Time(s2.UpdatedAt)

	return nil
}

// VpcPage is the page returned by a pager when traversing over a collection
// of VPCs. Orion does not paginate, so this is always a single page.
type VpcPage struct {
	pagination.SinglePageBase
}

// IsEmpty checks whether a VpcPage struct is empty.
func (r VpcPage) IsEmpty() (bool, error) {
	if r.StatusCode == 204 {
		return true, nil
	}

	is, err := ExtractVpcs(r)
	return len(is) == 0, err
}

// ExtractVpcs accepts a Page struct, specifically a VpcPage struct,
// and extracts the elements into a slice of Vpc structs.
func ExtractVpcs(r pagination.Page) ([]Vpc, error) {
	var s []Vpc
	err := ExtractVpcsInto(r, &s)
	return s, err
}

// ExtractVpcsInto extracts VPCs from a page into the provided slice.
func ExtractVpcsInto(r pagination.Page, v any) error {
	return r.(VpcPage).ExtractIntoSlicePtr(v, "vpcs")
}
