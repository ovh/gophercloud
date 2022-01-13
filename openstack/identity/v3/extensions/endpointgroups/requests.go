package endpointgroups

import (
	"net/http"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// Get retrieves details on a single endpoint group, by ID.
func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(resourceURL(client, id), &r.Body, nil)
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
func CreateProjectAssociation(client *gophercloud.ServiceClient, id string, projectId string) (r CreateProjectAssociationResult) {
	resp, err := client.Put(projectAssociationURL(client, id, projectId), nil, nil, &gophercloud.RequestOpts{
		OkCodes: []int{http.StatusNoContent},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// CheckProjectAssociation checks if an endpoint group is associated to a project.
func CheckProjectAssociation(client *gophercloud.ServiceClient, id string, projectId string) (r CheckProjectAssociationResult) {
	resp, err := client.Head(projectAssociationURL(client, id, projectId), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// DeleteProjectAssociation deletes an endpoint group to a project association.
func DeleteProjectAssociation(client *gophercloud.ServiceClient, id string, projectId string) (r DeleteProjectAssociationResult) {
	resp, err := client.Delete(projectAssociationURL(client, id, projectId), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
