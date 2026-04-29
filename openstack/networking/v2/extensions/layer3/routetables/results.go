package routetables

import (
	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// Route represents a route entry in a route table.
type Route struct {
	// Destination is the destination CIDR of the route.
	Destination string `json:"destination"`

	// Nexthop is the IP address of the next hop. It is empty for non-ip
	// NexthopType values.
	Nexthop string `json:"nexthop"`

	// NexthopType is the type of the next hop: local, ip, blackhole,
	// unreachable or prohibit.
	NexthopType string `json:"nexthop_type"`
}

// RouteTable represents a route table resource.
type RouteTable struct {
	// ID is the unique identifier of the route table.
	ID string `json:"id"`

	// Name is the human-readable name of the route table.
	Name string `json:"name"`

	// Description is a human-readable description of the route table.
	Description string `json:"description"`

	// RouterID is the ID of the VPC router this route table belongs to.
	RouterID string `json:"router_id"`

	// Routes is the list of routes in this route table.
	Routes []Route `json:"routes"`

	// Subnets is the list of subnet IDs associated with this route table.
	Subnets []string `json:"subnets"`

	// VpcStatus is the sync state of the route table with the backend:
	// ACTIVE, PENDING_CREATE, PENDING_UPDATE, PENDING_DELETE or ERROR.
	VpcStatus string `json:"vpc_status"`

	// TenantID is the project owner of the route table.
	TenantID string `json:"tenant_id"`

	// ProjectID is the project owner of the route table.
	ProjectID string `json:"project_id"`
}

type commonResult struct {
	gophercloud.Result
}

// Extract interprets the result as a RouteTable.
func (r commonResult) Extract() (*RouteTable, error) {
	var s RouteTable
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonResult) ExtractInto(v any) error {
	return r.ExtractIntoStructPtr(v, "route_table")
}

// actionResult is the result of the add/remove routes and subnets actions,
// which return the route table without the "route_table" wrapper.
type actionResult struct {
	gophercloud.Result
}

// Extract interprets the result as a RouteTable.
func (r actionResult) Extract() (*RouteTable, error) {
	var s RouteTable
	err := r.ExtractInto(&s)
	return &s, err
}

func (r actionResult) ExtractInto(v any) error {
	return r.ExtractIntoStructPtr(v, "")
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a RouteTable.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a RouteTable.
type GetResult struct {
	commonResult
}

// UpdateResult represents the result of an update operation. Call its Extract
// method to interpret it as a RouteTable.
type UpdateResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation. Call its
// ExtractErr method to determine if the request succeeded or failed.
type DeleteResult struct {
	gophercloud.ErrResult
}

// AddRoutesResult represents the result of an add routes operation. Call its
// Extract method to interpret it as a RouteTable.
type AddRoutesResult struct {
	actionResult
}

// RemoveRoutesResult represents the result of a remove routes operation.
// Call its Extract method to interpret it as a RouteTable.
type RemoveRoutesResult struct {
	actionResult
}

// AddSubnetsResult represents the result of an add subnets operation. Call
// its Extract method to interpret it as a RouteTable.
type AddSubnetsResult struct {
	actionResult
}

// RemoveSubnetsResult represents the result of a remove subnets operation.
// Call its Extract method to interpret it as a RouteTable.
type RemoveSubnetsResult struct {
	actionResult
}

// RouteTablePage is the page returned by a pager when traversing over a
// collection of route tables.
type RouteTablePage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of route tables has
// reached the end of a page and the pager seeks to traverse over a new one.
func (r RouteTablePage) NextPageURL() (string, error) {
	var s struct {
		Links []gophercloud.Link `json:"route_tables_links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gophercloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a RouteTablePage struct is empty.
func (r RouteTablePage) IsEmpty() (bool, error) {
	if r.StatusCode == 204 {
		return true, nil
	}

	is, err := ExtractRouteTables(r)
	return len(is) == 0, err
}

// ExtractRouteTables accepts a Page struct, specifically a RouteTablePage,
// and extracts the elements into a slice of RouteTable structs.
func ExtractRouteTables(r pagination.Page) ([]RouteTable, error) {
	var s []RouteTable
	err := ExtractRouteTablesInto(r, &s)
	return s, err
}

// ExtractRouteTablesInto extracts the elements of a RouteTablePage into the
// provided slice.
func ExtractRouteTablesInto(r pagination.Page, v any) error {
	return r.(RouteTablePage).ExtractIntoSlicePtr(v, "route_tables")
}
