package quotasets

import (
	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// QuotaSet is a set of operational limits that allow for control of manila
// usage.
type QuotaSet struct {
	// Gigabytes is the total size of share storage for the project in gigabytes.
	Gigabytes *int `json:"gigabytes,omitempty"`

	// Snapshots is the total number of share snapshots for the project.
	Snapshots *int `json:"snapshots,omitempty"`

	// Shares is the total number of shares for the project.
	Shares *int `json:"shares,omitempty"`

	// SnapshotGigabytes is the total size of share snapshots for the project in gigabytes.
	SnapshotGigabytes *int `json:"snapshot_gigabytes,omitempty"`

	// Share network is the total number of share networks for the project.
	ShareNetworks *int `json:"share_networks,omitempty"`

	// Share groups is the total number of share groups for the project.
	ShareGroups *int `json:"share_groups,omitempty"`

	// Share group snapshots is the total number of share group snapshots for the project.
	ShareGroupSnapshots *int `json:"share_group_snapshots,omitempty"`

	// Share Replicas is the total number of share replicas for the project.
	ShareReplicas *int `json:"share_replicas,omitempty"`

	// Share Replica Gigabytes is the total size of share replicas for the project in gigabytes.
	ShareReplicaGigabytes *int `json:"share_replica_gigabytes,omitempty"`

	// PerShareGigabytes is the maximum size of a share for the project in gigabytes.
	PerShareGigabytes *int `json:"per_share_gigabytes,omitempty"`

	// Backups is the maximum number of backups allowed for each project.
	Backups *int `json:"backups,omitempty"`

	// BackupsGigabytes is the maximum number of gigabytes for the backups allowed for each project.
	BackupsGigabytes *int `json:"backup_gigabytes,omitempty"`
}

// QuotaDetailSet represents details of both operational limits of shares file system resources for a project
// and the current usage of those resources.
type QuotaDetailSet struct {
	// Gigabytes is the total size of share storage for the project in gigabytes.
	Gigabytes QuotaDetail `json:"gigabytes,omitempty"`

	// Snapshots is the total number of share snapshots for the project.
	Snapshots QuotaDetail `json:"snapshots,omitempty"`

	// Shares is the total number of shares for the project.
	Shares QuotaDetail `json:"shares,omitempty"`

	// SnapshotGigabytes is the total size of share snapshots for the project in gigabytes.
	SnapshotGigabytes QuotaDetail `json:"snapshot_gigabytes,omitempty"`

	// Share network is the total number of share networks for the project.
	ShareNetworks QuotaDetail `json:"share_networks,omitempty"`

	// Share groups is the total number of share groups for the project.
	ShareGroups QuotaDetail `json:"share_groups,omitempty"`

	// Share group snapshots is the total number of share group snapshots for the project.
	ShareGroupSnapshots QuotaDetail `json:"share_group_snapshots,omitempty"`

	// Share Replicas is the total number of share replicas for the project.
	ShareReplicas QuotaDetail `json:"share_replicas,omitempty"`

	// Share Replica Gigabytes is the total size of share replicas for the project in gigabytes.
	ShareReplicaGigabytes QuotaDetail `json:"share_replica_gigabytes,omitempty"`

	// PerShareGigabytes is the maximum size of a share for the project in gigabytes.
	PerShareGigabytes QuotaDetail `json:"per_share_gigabytes,omitempty"`

	// Backups is the maximum number of backups allowed for each project.
	Backups QuotaDetail `json:"backups,omitempty"`

	// BackupsGigabytes is the maximum number of gigabytes for the backups allowed for each project.
	BackupsGigabytes QuotaDetail `json:"backup_gigabytes,omitempty"`
}

// QuotaDetail is a set of details about a single operational limit that allows
// for control of shared file system usage.
type QuotaDetail struct {
	// InUse is the current number of provisioned/allocated resources of the
	// given type.
	InUse int `json:"in_use"`

	// Reserved is a transitional state when a claim against quota has been made
	// but the resource is not yet fully online.
	Reserved int `json:"reserved"`

	// Limit is the maximum number of a given resource that can be
	// allocated/provisioned.  This is what "quota" usually refers to.
	Limit int `json:"limit"`
}

// QuotaSetPage stores a single page of all QuotaSet results from a List call.
type QuotaSetPage struct {
	pagination.SinglePageBase
}

// IsEmpty determines whether or not a QuotaSetPage is empty.
func (page QuotaSetPage) IsEmpty() (bool, error) {
	ks, err := ExtractQuotaSets(page)
	return len(ks) == 0, err
}

// ExtractQuotaSets interprets a page of results as a slice of QuotaSets.
func ExtractQuotaSets(r pagination.Page) ([]QuotaSet, error) {
	var s struct {
		QuotaSets []QuotaSet `json:"quotas"`
	}
	err := (r.(QuotaSetPage)).ExtractInto(&s)
	return s.QuotaSets, err
}

type quotaResult struct {
	gophercloud.Result
}

type quotaDetailResult struct {
	gophercloud.Result
}

// Extract is a function that accepts a result and extracts a QuotaDetailSet resource.
func (r quotaDetailResult) Extract() (*QuotaDetailSet, error) {
	var s struct {
		Quota *QuotaDetailSet `json:"quota_set"`
	}
	err := r.ExtractInto(&s)
	return s.Quota, err
}

// Extract is a method that attempts to interpret any QuotaSet resource response
// as a QuotaSet struct.
func (r quotaResult) Extract() (*QuotaSet, error) {
	var s struct {
		QuotaSet *QuotaSet `json:"quota_set"`
	}
	err := r.ExtractInto(&s)
	return s.QuotaSet, err
}

// GetResult is the response from a Get operation. Call its Extract method to
// interpret it as a QuotaSet.
type GetResult struct {
	quotaResult
}

// UpdateResult is the response from a Update operation. Call its Extract method
// to interpret it as a QuotaSet.
type UpdateResult struct {
	quotaResult
}

// GetDetailResult represents the detailed result of a get operation. Call its Extract
// method to interpret it as a Quota.
type GetDetailResult struct {
	quotaDetailResult
}
