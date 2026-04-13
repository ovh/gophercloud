package endpointgroups

import (
	"context"
	"net/http"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// Get retrieves details on a single endpoint group, by ID.
func Get(ctx context.Context, client *gophercloud.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(ctx, resourceURL(client, id), &r.Body, nil)
	return
}

// ListOptsBuilder allows extensions to add additional parameters to
// the List request
type ListOptsBuilder interface {
	ToEndpointGroupListQuery() (string, error)
}

// ListOpts provides options to filter the List results.
type ListOpts struct {
	// Name filters the response by endpoint group name.
	Name string `q:"name"`
}

// ToEndpointGroupListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToEndpointGroupListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List enumerates the endpoint groups.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := rootURL(client)
	if opts != nil {
		query, err := opts.ToEndpointGroupListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return EndpointGroupPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// ListForProjects enumerates the endpoint groups associated to a project.
func ListForProjects(client *gophercloud.ServiceClient, projectId string) pagination.Pager {
	return pagination.NewPager(client, listEndpointGroupsAssociationURL(client, projectId), func(r pagination.PageResult) pagination.Page {
		return EndpointGroupPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// CreateProjectAssociation creates an endpoint group to a project association.
func CreateProjectAssociation(ctx context.Context, client *gophercloud.ServiceClient, id string, projectId string) (r CreateProjectAssociationResult) {
	resp, err := client.Put(ctx, projectAssociationURL(client, id, projectId), nil, nil, &gophercloud.RequestOpts{
		OkCodes: []int{http.StatusNoContent},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// CheckProjectAssociation checks if an endpoint group is associated to a project.
func CheckProjectAssociation(ctx context.Context, client *gophercloud.ServiceClient, id string, projectId string) (r CheckProjectAssociationResult) {
	resp, err := client.Head(ctx, projectAssociationURL(client, id, projectId), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// DeleteProjectAssociation deletes an endpoint group to a project association.
func DeleteProjectAssociation(ctx context.Context, client *gophercloud.ServiceClient, id string, projectId string) (r DeleteProjectAssociationResult) {
	resp, err := client.Delete(ctx, projectAssociationURL(client, id, projectId), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// CreateOptsBuilder allows extensions to add additional parameters to
// the Create request.
type CreateOptsBuilder interface {
	ToGroupCreateMap() (map[string]interface{}, error)
}

// CreateOpts provides options used to create a group.
type CreateOpts struct {
	// Name is the name of the new endpoint group.
	Name string `json:"name" required:"true"`

	// Description is a description of the endpoint group.
	Description *string `json:"description,omitempty"`

	// Filters are the filters of the endpoint group.
	Filters map[string]interface{} `json:"filters,omitempty"`
}

// ToGroupCreateMap formats a CreateOpts into a create request.
func (opts CreateOpts) ToGroupCreateMap() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "endpoint_group")
	if err != nil {
		return nil, err
	}

	return b, nil
}

// Create creates a new endpoint group.
func Create(ctx context.Context, client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToGroupCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, createURL(client), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to
// the Update request.
type UpdateOptsBuilder interface {
	ToGroupUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts provides options for updating an endpoint group.
type UpdateOpts struct {
	// Name is the name of the endpoint group.
	Name string `json:"name,omitempty"`

	// Description is a description of the endpoint group.
	Description *string `json:"description,omitempty"`

	// Filters are the filters of the endpoint group.
	Filters map[string]interface{} `json:"filters,omitempty"`
}

// ToGroupUpdateMap formats a UpdateOpts into an update request.
func (opts UpdateOpts) ToGroupUpdateMap() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "endpoint_group")
	if err != nil {
		return nil, err
	}

	return b, nil
}

// Update updates an existing Endpoint Group.
func Update(ctx context.Context, client *gophercloud.ServiceClient, groupID string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToGroupUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Patch(ctx, updateURL(client, groupID), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete deletes an endpoint group.
func Delete(ctx context.Context, client *gophercloud.ServiceClient, groupID string) (r DeleteResult) {
	resp, err := client.Delete(ctx, deleteURL(client, groupID), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
