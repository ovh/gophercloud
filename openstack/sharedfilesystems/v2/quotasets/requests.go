package quotasets

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
)

// Get returns data about a previously created QuotaSet.
func Get(ctx context.Context, client *gophercloud.ServiceClient, tenantID string) (r GetResult) {
	resp, err := client.Get(ctx, getURL(client, tenantID), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// GetDetail returns detailed Networking Quotas for a project.
func GetDetail(ctx context.Context, client *gophercloud.ServiceClient, projectID string) (r GetDetailResult) {
	resp, err := client.Get(ctx, getDetailURL(client, projectID), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Updates the quotas for the given tenantID.
func Update(ctx context.Context, client *gophercloud.ServiceClient, tenantID string, opts UpdateOptsBuilder) (r UpdateResult) {
	reqBody, err := opts.ToManillaQuotaUpdateMap()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Put(ctx, updateURL(client, tenantID), reqBody, &r.Body, &gophercloud.RequestOpts{OkCodes: []int{200}})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Resets the quotas for the given tenant to their default values.
func Delete(ctx context.Context, client *gophercloud.ServiceClient, tenantID string) (r DeleteResult) {
	resp, err := client.Delete(ctx, deleteURL(client, tenantID), &gophercloud.RequestOpts{OkCodes: []int{200, 202, 204}})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Get returns data about a previously created QuotaSet for a share type.
func GetByShareType(ctx context.Context, client *gophercloud.ServiceClient, tenantID string, share_type string) (r GetResult) {
	resp, err := client.Get(ctx, getURLbyShareType(client, tenantID, share_type), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Updates the quotas for a given sharetype.
func UpdateByShareType(ctx context.Context, client *gophercloud.ServiceClient, tenantID string, share_type string, opts UpdateOptsBuilder) (r UpdateResult) {
	reqBody, err := opts.ToManillaQuotaUpdateMap()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Put(ctx, updateURLByShareType(client, tenantID, share_type), reqBody, &r.Body, &gophercloud.RequestOpts{OkCodes: []int{200}})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Get returns data about a previously created QuotaSet for a user id.
func GetByUser(ctx context.Context, client *gophercloud.ServiceClient, tenantID string, user_id string) (r GetResult) {
	resp, err := client.Get(ctx, getURLbyUser(client, tenantID, user_id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Updates the quotas for a given user.
func UpdateByUser(ctx context.Context, client *gophercloud.ServiceClient, tenantID string, user_id string, opts UpdateOptsBuilder) (r UpdateResult) {
	reqBody, err := opts.ToManillaQuotaUpdateMap()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Put(ctx, updateURLByUser(client, tenantID, user_id), reqBody, &r.Body, &gophercloud.RequestOpts{OkCodes: []int{200}})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Options for Updating the quotas of a Tenant.
// All int-values are pointers so they can be nil if they are not needed.
// You can use gophercloud.IntToPointer() for convenience
type UpdateOpts struct {
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

	// Replica Gigabytes is the total size of share replicas for the project in gigabytes.
	ReplicaGigabytes *int `json:"replica_gigabytes,omitempty"`

	// PerShareGigabytes is the maximum size of a share for the project in gigabytes.
	PerShareGigabytes *int `json:"per_share_gigabytes,omitempty"`

	// Backups is the maximum number of backups allowed for each project.
	Backups *int `json:"backups,omitempty"`

	// BackupsGigabytes is the maximum number of gigabytes for the backups allowed for each project.
	BackupsGigabytes *int `json:"backup_gigabytes,omitempty"`

	// EncryptionKeys is the maximum number of encryption keys allowed for each project.
	EncryptionKeys *int `json:"encryption_keys,omitempty"`
}

// UpdateOptsBuilder enables extensions to add parameters to the update request.
type UpdateOptsBuilder interface {
	// Extra specific name to prevent collisions with interfaces for other quotas
	// (e.g. neutron)
	ToManillaQuotaUpdateMap() (map[string]any, error)
}

// ToComputeManillaUpdateMap builds the update options into a serializable
// format.
func (opts UpdateOpts) ToManillaQuotaUpdateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "quota_set")

}
