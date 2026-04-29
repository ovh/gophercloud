package routetables

import (
	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// Route represents a route entry in a route table.
type Route struct {
	// Destination is the destination CIDR of the route.
	Destination string `json:"destination"`

	// Nexthop is the IP address of the next hop. Can be empty for non-IP nexthop types.
	Nexthop string `json:"nexthop"`

	// NexthopType is the type of the nexthop: local, ip, blackhole, unreachable, prohibit.
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

	// RouterID is the ID of the router this route table belongs to.
	RouterID string `json:"router_id"`

	// IsDefault indicates whether this is the default route table for the router.
	IsDefault bool `json:"is_default"`

	// Routes is the list of routes in this route table.
	Routes []Route `json:"routes"`

	// Subnets is the list of subnet IDs associated with this route table.
	Subnets []string `json:"subnets"`

	// ProjectID is the project owner of the route table.
	ProjectID string `json:"project_id"`
}

type commonResult struct {
	gophercloud.Result
}

// Extract interprets the result as a RouteTable.
func (r commonResult) Extract() (*RouteTable, error) {
	var s struct {
		RouteTable *RouteTable `json:"route_table"`
	}
	err := r.ExtractInto(&s)
	return s.RouteTable, err
}

// CreateResult represents the result of a create operation.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation.
type GetResult struct {
	commonResult
}

// UpdateResult represents the result of an update operation.
type UpdateResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation.
type DeleteResult struct {
	gophercloud.ErrResult
}

// AddRoutesResult represents the result of an add routes operation.
type AddRoutesResult struct {
	commonResult
}

// RemoveRoutesResult represents the result of a remove routes operation.
type RemoveRoutesResult struct {
	commonResult
}

// AddSubnetsResult represents the result of an add subnets operation.
type AddSubnetsResult struct {
	commonResult
}

// RemoveSubnetsResult represents the result of a remove subnets operation.
type RemoveSubnetsResult struct {
	commonResult
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
	err := r.(RouteTablePage).ExtractIntoSlicePtr(&s, "route_tables")
	return s, err
}
