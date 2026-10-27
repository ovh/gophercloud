package ipsecpolicies

import (
	"encoding/json"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// Policy is an IPsec policy.
type Policy struct {
	// ID is the ID of the policy.
	ID string `json:"id"`

	// ProjectID is the ID of the project owning the policy.
	ProjectID string `json:"project_id"`

	// Name is the human readable name of the policy.
	Name string `json:"name"`

	// Description is the human readable description of the policy.
	Description string `json:"description"`

	// TransformProtocol is the transform protocol.
	TransformProtocol string `json:"transform_protocol"`

	// AuthAlgorithm is the authentication hash algorithm. It is nil with an
	// AES-GCM cipher, which carries its own integrity.
	AuthAlgorithm *string `json:"auth_algorithm"`

	// EncryptionAlgorithm is the encryption algorithm.
	EncryptionAlgorithm string `json:"encryption_algorithm"`

	// EncapsulationMode is the encapsulation mode.
	EncapsulationMode string `json:"encapsulation_mode"`

	// PFS is the Perfect forward secrecy group.
	PFS string `json:"pfs"`

	// Lifetime is the lifetime of the security association.
	Lifetime Lifetime `json:"lifetime"`

	// CreatedAt is the time the policy was created.
	CreatedAt time.Time `json:"-"`

	// UpdatedAt is the time the policy was last updated.
	UpdatedAt time.Time `json:"-"`
}

// Lifetime is the lifetime of a security association.
type Lifetime struct {
	// Units is the unit of the lifetime.
	Units string `json:"units"`

	// Value is the lifetime value.
	Value int `json:"value"`
}

func (r *Policy) UnmarshalJSON(b []byte) error {
	type tmp Policy
	var s struct {
		tmp
		CreatedAt gophercloud.JSONRFC3339MilliNoZ `json:"created_at"`
		UpdatedAt gophercloud.JSONRFC3339MilliNoZ `json:"updated_at"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*r = Policy(s.tmp)
	r.CreatedAt = time.Time(s.CreatedAt)
	r.UpdatedAt = time.Time(s.UpdatedAt)
	return nil
}

type commonResult struct {
	gophercloud.Result
}

// Extract is a function that accepts a result and extracts an IPsec Policy.
func (r commonResult) Extract() (*Policy, error) {
	var s struct {
		Policy *Policy `json:"ipsec_policy"`
	}
	err := r.ExtractInto(&s)
	return s.Policy, err
}

// PolicyPage is the page returned by a pager when traversing over a
// collection of Policies.
type PolicyPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of IPsec policies has
// reached the end of a page and the pager seeks to traverse over a new one.
// In order to do this, it needs to construct the next page's URL.
func (r PolicyPage) NextPageURL() (string, error) {
	var s struct {
		Links []gophercloud.Link `json:"ipsec_policies_links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gophercloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a PolicyPage struct is empty.
func (r PolicyPage) IsEmpty() (bool, error) {
	if r.StatusCode == 204 {
		return true, nil
	}

	is, err := ExtractPolicies(r)
	return len(is) == 0, err
}

// ExtractPolicies accepts a Page struct, specifically a PolicyPage struct,
// and extracts the elements into a slice of Policy structs.
func ExtractPolicies(r pagination.Page) ([]Policy, error) {
	var s struct {
		Policies []Policy `json:"ipsec_policies"`
	}
	err := (r.(PolicyPage)).ExtractInto(&s)
	return s.Policies, err
}

// CreateResult represents the result of a Create operation. Call its Extract method to
// interpret it as a Policy.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a Get operation. Call its Extract method to
// interpret it as a Policy.
type GetResult struct {
	commonResult
}

// DeleteResult represents the results of a Delete operation. Call its ExtractErr method
// to determine whether the operation succeeded or failed.
type DeleteResult struct {
	gophercloud.ErrResult
}

// UpdateResult represents the result of an update operation. Call its Extract method
// to interpret it as a Policy.
type UpdateResult struct {
	commonResult
}
