package ikepolicies

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

type AuthAlgorithm string
type EncryptionAlgorithm string
type PFS string
type Unit string
type IKEVersion string
type Phase1NegotiationMode string

const (
	AuthAlgorithmSHA256          AuthAlgorithm         = "sha256"
	AuthAlgorithmSHA384          AuthAlgorithm         = "sha384"
	AuthAlgorithmSHA512          AuthAlgorithm         = "sha512"
	EncryptionAlgorithmAES128    EncryptionAlgorithm   = "aes-128"
	EncryptionAlgorithmAES192    EncryptionAlgorithm   = "aes-192"
	EncryptionAlgorithmAES256    EncryptionAlgorithm   = "aes-256"
	EncryptionAlgorithmAES128GCM EncryptionAlgorithm   = "aes-128-gcm-16"
	EncryptionAlgorithmAES256GCM EncryptionAlgorithm   = "aes-256-gcm-16"
	PFSGroup14                   PFS                   = "group14"
	PFSGroup15                   PFS                   = "group15"
	PFSGroup16                   PFS                   = "group16"
	PFSGroup19                   PFS                   = "group19"
	PFSGroup20                   PFS                   = "group20"
	PFSGroup21                   PFS                   = "group21"
	PFSGroup31                   PFS                   = "group31"
	UnitSeconds                  Unit                  = "seconds"
	IKEVersionv2                 IKEVersion            = "v2"
	Phase1NegotiationModeMain    Phase1NegotiationMode = "main"
)

// LifetimeOpts is the lifetime of the IKE security association.
// The service stores it as a whole, so an update must carry both parts to
// keep the units.
type LifetimeOpts struct {
	// Units is the unit of the lifetime. Only seconds is accepted.
	Units Unit `json:"units,omitempty"`

	// Value is the lifetime value, at least 60.
	Value int `json:"value"`
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToPolicyCreateMap() (map[string]any, error)
}

// CreateOpts contains all the values needed to create a new IKE policy.
type CreateOpts struct {
	// ProjectID specifies a project to own the IKE policy. The caller must be
	// allowed to act on that project in order to set this. Otherwise, the
	// caller's project will be the owner.
	ProjectID string `json:"project_id,omitempty"`

	// Name is the human readable name of the policy.
	Name string `json:"name,omitempty"`

	// Description is the human readable description of the policy.
	Description string `json:"description,omitempty"`

	// AuthAlgorithm is the authentication hash algorithm. With an AES-GCM
	// cipher it is the PRF only. The default is sha256.
	AuthAlgorithm AuthAlgorithm `json:"auth_algorithm,omitempty"`

	// EncryptionAlgorithm is the encryption algorithm. The default is aes-256.
	EncryptionAlgorithm EncryptionAlgorithm `json:"encryption_algorithm,omitempty"`

	// Phase1NegotiationMode is the IKE mode. Only main is accepted.
	Phase1NegotiationMode Phase1NegotiationMode `json:"phase1_negotiation_mode,omitempty"`

	// IKEVersion is the IKE version. Only v2 is accepted.
	IKEVersion IKEVersion `json:"ike_version,omitempty"`

	// PFS is the Perfect forward secrecy group. The default is group14.
	PFS PFS `json:"pfs,omitempty"`

	// Lifetime is the lifetime of the security association.
	// The default is 10800 seconds.
	Lifetime *LifetimeOpts `json:"lifetime,omitempty"`
}

// ToPolicyCreateMap casts a CreateOpts struct to a map.
func (opts CreateOpts) ToPolicyCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "ike_policy")
}

// Create accepts a CreateOpts struct and uses the values to create a new
// IKE policy.
func Create(ctx context.Context, c *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToPolicyCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Post(ctx, rootURL(c), b, &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Get retrieves a particular IKE policy based on its unique ID.
func Get(ctx context.Context, c *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := c.Get(ctx, resourceURL(c, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete will permanently delete a particular IKE policy based on its
// unique ID. A policy still used by an IPsec site connection cannot be
// deleted.
func Delete(ctx context.Context, c *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := c.Delete(ctx, resourceURL(c, id), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToPolicyListQuery() (string, error)
}

// ListOpts allows the filtering, sorting and paging of IKE policies. The
// service rejects any other query parameter.
type ListOpts struct {
	// ProjectID lists the policies of this project only. Naming a project
	// other than the caller's requires the get_all-global policy.
	ProjectID string `q:"project_id"`

	// Name lists the policies with this name only.
	Name string `q:"name"`

	// Limit is the maximum number of policies per page.
	Limit int `q:"limit"`

	// Marker is the ID of the last policy of the previous page.
	Marker string `q:"marker"`

	// PageReverse pages backwards from Marker.
	PageReverse bool `q:"page_reverse"`

	// SortKey is the attribute to sort on.
	SortKey string `q:"sort_key"`

	// SortDir is the sort direction, asc or desc. It requires SortKey.
	SortDir string `q:"sort_dir"`

	// Fields restricts each returned policy to these attributes.
	Fields []string `q:"fields"`
}

// ToPolicyListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToPolicyListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List returns a Pager which allows you to iterate over a collection of
// IKE policies. It accepts a ListOpts struct, which allows you to filter
// the returned collection for greater efficiency.
func List(c *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := rootURL(c)
	if opts != nil {
		query, err := opts.ToPolicyListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return PolicyPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToPolicyUpdateMap() (map[string]any, error)
}

// UpdateOpts contains the values used when updating an IKE policy.
type UpdateOpts struct {
	// Name is the human readable name of the policy.
	Name *string `json:"name,omitempty"`

	// Description is the human readable description of the policy.
	Description *string `json:"description,omitempty"`

	// AuthAlgorithm is the authentication hash algorithm.
	AuthAlgorithm AuthAlgorithm `json:"auth_algorithm,omitempty"`

	// EncryptionAlgorithm is the encryption algorithm.
	EncryptionAlgorithm EncryptionAlgorithm `json:"encryption_algorithm,omitempty"`

	// Phase1NegotiationMode is the IKE mode. Only main is accepted.
	Phase1NegotiationMode Phase1NegotiationMode `json:"phase1_negotiation_mode,omitempty"`

	// IKEVersion is the IKE version. Only v2 is accepted.
	IKEVersion IKEVersion `json:"ike_version,omitempty"`

	// PFS is the Perfect forward secrecy group.
	PFS PFS `json:"pfs,omitempty"`

	// Lifetime replaces the lifetime of the security association.
	Lifetime *LifetimeOpts `json:"lifetime,omitempty"`
}

// ToPolicyUpdateMap casts an UpdateOpts struct to a map.
func (opts UpdateOpts) ToPolicyUpdateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "ike_policy")
}

// Update allows IKE policies to be updated.
func Update(ctx context.Context, c *gophercloud.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToPolicyUpdateMap()
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
