package quotas

import (
	"github.com/gophercloud/gophercloud/v2"
)

// Quota represents VPC quota information for a project.
type Quota struct {
	// OpenStack project ID.
	ProjectID string `json:"project_id"`

	// Maximum number of VPCs for this project.
	MaxVpcs int `json:"max_vpcs"`

	// Whether this is the global default or a per-project override.
	// Possible values: "config_default", "project_override".
	Source string `json:"source"`
}

// QuotaResult represents the response from a Set operation (no source field).
type QuotaResult struct {
	// OpenStack project ID.
	ProjectID string `json:"project_id"`

	// Maximum number of VPCs for this project.
	MaxVpcs int `json:"max_vpcs"`
}

// GetResult represents the result of a Get operation.
type GetResult struct {
	gophercloud.Result
}

// Extract interprets a GetResult as a Quota.
func (r GetResult) Extract() (*Quota, error) {
	var s Quota
	err := r.ExtractIntoStructPtr(&s, "quota")
	return &s, err
}

// SetResult represents the result of a Set operation.
type SetResult struct {
	gophercloud.Result
}

// Extract interprets a SetResult as a QuotaResult.
func (r SetResult) Extract() (*QuotaResult, error) {
	var s QuotaResult
	err := r.ExtractIntoStructPtr(&s, "quota")
	return &s, err
}

// DeleteResult represents the result of a Delete operation.
type DeleteResult struct {
	gophercloud.ErrResult
}
