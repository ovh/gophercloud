package vpngateways

import (
	"encoding/json"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// Route is a static route of a VPN gateway.
type Route struct {
	// DestinationCIDR is the destination network of the route.
	DestinationCIDR string `json:"destination"`

	// NextHop is the IP address the route forwards to.
	NextHop string `json:"nexthop"`
}

// FixedIP is an IP address allocated to a port of a VPN gateway.
type FixedIP struct {
	// SubnetID is the subnet the address belongs to.
	SubnetID string `json:"subnet_id"`

	// IPAddress is the allocated address. It is empty until the port exists.
	IPAddress string `json:"ip_address"`
}

// Attachment is an external gateway or an internal interface of a VPN
// gateway: its port on a network.
type Attachment struct {
	// NetworkID is the network the port is on.
	NetworkID string `json:"network_id"`

	// PortID is the port of the VPN gateway. It is empty until the port
	// exists.
	PortID string `json:"port_id"`

	// FixedIPs is the list of addresses of the port.
	FixedIPs []FixedIP `json:"fixed_ips"`

	// Status is the provisioning status of the attachment.
	Status string `json:"status"`
}

// Instance is a compute instance backing a VPN gateway.
type Instance struct {
	// InstanceID is the ID of the compute instance.
	InstanceID string `json:"instance_id"`

	// ManagementIP is the address of the instance on the management network.
	ManagementIP string `json:"management_ip"`

	// ImageID is the image the instance was booted or last rebuilt with.
	ImageID string `json:"image_id"`
}

// VPNGateway is a VPN gateway, shaped like a Neutron router.
type VPNGateway struct {
	// ID is the ID of the VPN gateway.
	ID string `json:"id"`

	// ProjectID is the ID of the project owning the VPN gateway.
	ProjectID string `json:"project_id"`

	// Name is the human readable name of the VPN gateway.
	Name string `json:"name"`

	// Description is the human readable description of the VPN gateway.
	Description string `json:"description"`

	// AdminStateUp is the administrative state of the VPN gateway.
	AdminStateUp bool `json:"admin_state_up"`

	// Status is the provisioning status of the VPN gateway.
	Status string `json:"status"`

	// ExternalIPv4 is the first IPv4 address of the active external gateway.
	ExternalIPv4 string `json:"external_ipv4"`

	// ExternalIPv6 is the first IPv6 address of the active external gateway.
	ExternalIPv6 string `json:"external_ipv6"`

	// HA tells whether the VPN gateway is highly available.
	HA bool `json:"ha"`

	// FlavorID is the ID of the flavor of the VPN gateway.
	FlavorID string `json:"flavor_id"`

	// Routes is the list of static routes of the VPN gateway.
	Routes []Route `json:"routes"`

	// AvailabilityZoneHints is the list of requested availability zones.
	AvailabilityZoneHints []string `json:"availability_zone_hints"`

	// AvailabilityZones is the list of availability zones of the VPN gateway.
	AvailabilityZones []string `json:"availability_zones"`

	// RevisionNumber is the revision number of the VPN gateway.
	RevisionNumber int `json:"revision_number"`

	// ExternalGateways is the list of external network attachments.
	ExternalGateways []Attachment `json:"external_gateways"`

	// InternalInterfaces is the list of internal subnet attachments.
	InternalInterfaces []Attachment `json:"internal_interfaces"`

	// Instances is the list of compute instances backing the VPN gateway.
	// It is only returned to callers allowed the
	// vpn_gateway:get:admin_attributes policy, admins by default.
	Instances []Instance `json:"instances"`

	// Tags is the list of tags of the VPN gateway.
	Tags []string `json:"tags"`

	// CreatedAt is the time the VPN gateway was created.
	CreatedAt time.Time `json:"-"`

	// UpdatedAt is the time the VPN gateway was last updated.
	UpdatedAt time.Time `json:"-"`
}

func (r *VPNGateway) UnmarshalJSON(b []byte) error {
	type tmp VPNGateway
	var s struct {
		tmp
		CreatedAt gophercloud.JSONRFC3339MilliNoZ `json:"created_at"`
		UpdatedAt gophercloud.JSONRFC3339MilliNoZ `json:"updated_at"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*r = VPNGateway(s.tmp)
	r.CreatedAt = time.Time(s.CreatedAt)
	r.UpdatedAt = time.Time(s.UpdatedAt)
	return nil
}

type commonResult struct {
	gophercloud.Result
}

// Extract is a function that accepts a result and extracts a VPN gateway.
func (r commonResult) Extract() (*VPNGateway, error) {
	var s struct {
		VPNGateway *VPNGateway `json:"vpn_gateway"`
	}
	err := r.ExtractInto(&s)
	return s.VPNGateway, err
}

// VPNGatewayPage is the page returned by a pager when traversing over a
// collection of VPN gateways.
type VPNGatewayPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of VPN gateways has
// reached the end of a page and the pager seeks to traverse over a new one.
// In order to do this, it needs to construct the next page's URL.
func (r VPNGatewayPage) NextPageURL() (string, error) {
	var s struct {
		Links []gophercloud.Link `json:"vpn_gateways_links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gophercloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a VPNGatewayPage struct is empty.
func (r VPNGatewayPage) IsEmpty() (bool, error) {
	if r.StatusCode == 204 {
		return true, nil
	}

	is, err := ExtractVPNGateways(r)
	return len(is) == 0, err
}

// ExtractVPNGateways accepts a Page struct, specifically a VPNGatewayPage
// struct, and extracts the elements into a slice of VPNGateway structs.
func ExtractVPNGateways(r pagination.Page) ([]VPNGateway, error) {
	var s struct {
		VPNGateways []VPNGateway `json:"vpn_gateways"`
	}
	err := (r.(VPNGatewayPage)).ExtractInto(&s)
	return s.VPNGateways, err
}

// CreateResult represents the result of a Create operation. Call its Extract method to
// interpret it as a VPNGateway.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a Get operation. Call its Extract method to
// interpret it as a VPNGateway.
type GetResult struct {
	commonResult
}

// UpdateResult represents the result of an update operation. Call its Extract method
// to interpret it as a VPNGateway.
type UpdateResult struct {
	commonResult
}

// DeleteResult represents the results of a Delete operation. Call its ExtractErr method
// to determine whether the operation succeeded or failed.
type DeleteResult struct {
	gophercloud.ErrResult
}

// AddExternalGatewaysResult represents the result of an AddExternalGateways
// operation. Call its Extract method to interpret it as a VPNGateway.
type AddExternalGatewaysResult struct {
	commonResult
}

// ActionResult represents the result of a Rebuild or Resync operation. Call
// its Extract method to interpret it as a VPNGateway.
type ActionResult struct {
	commonResult
}

// InterfaceInfo represents information about an internal interface of a VPN
// gateway, as returned by AddInterface and RemoveInterface.
type InterfaceInfo struct {
	// ID is the ID of the VPN gateway.
	ID string `json:"id"`

	// ProjectID is the ID of the project owning the VPN gateway.
	ProjectID string `json:"project_id"`

	// NetworkID is the network of the interface. It may be empty until the
	// port exists.
	NetworkID string `json:"network_id"`

	// PortID is the port of the interface. It is empty until the port exists.
	PortID string `json:"port_id"`

	// SubnetID is the subnet of the interface.
	SubnetID string `json:"subnet_id"`

	// SubnetIDs is the list of subnets of the interface.
	SubnetIDs []string `json:"subnet_ids"`

	// Tags is the list of tags of the interface.
	Tags []string `json:"tags"`
}

// InterfaceResult represents the result of interface operations, such as
// AddInterface and RemoveInterface. Call its Extract method to interpret
// the result as an InterfaceInfo.
type InterfaceResult struct {
	gophercloud.Result
}

// Extract is a function that accepts a result and extracts an information
// object.
func (r InterfaceResult) Extract() (*InterfaceInfo, error) {
	var s InterfaceInfo
	err := r.ExtractInto(&s)
	return &s, err
}
