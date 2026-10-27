package testing

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/vpngateway/v1/vpngateways"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

const GatewayID = "2ad3c5dc-8c7e-4ff0-8b54-4b4d4e2c7c6e"

const CreateRequest = `
{
    "vpn_gateway": {
        "project_id": "9145d91459d248b1b02fdaca97c6a75d",
        "name": "gateway",
        "description": "VPN gateway",
        "admin_state_up": true,
        "ha": false,
        "flavor_id": "b8b9cc5a-0b05-4b1e-9d0a-1d6c0f5a3e01",
        "availability_zone_hints": ["az1"],
        "routes": [
            {
                "destination": "10.20.0.0/16",
                "nexthop": "10.0.0.254"
            }
        ]
    }
}
`

const CreatedGatewayBody = `
{
    "vpn_gateway": {
        "id": "2ad3c5dc-8c7e-4ff0-8b54-4b4d4e2c7c6e",
        "project_id": "9145d91459d248b1b02fdaca97c6a75d",
        "name": "gateway",
        "description": "VPN gateway",
        "admin_state_up": true,
        "status": "PENDING_CREATE",
        "external_ipv4": null,
        "external_ipv6": null,
        "ha": false,
        "flavor_id": "b8b9cc5a-0b05-4b1e-9d0a-1d6c0f5a3e01",
        "routes": [
            {
                "destination": "10.20.0.0/16",
                "nexthop": "10.0.0.254"
            }
        ],
        "availability_zone_hints": ["az1"],
        "availability_zones": [],
        "revision_number": 0,
        "external_gateways": [],
        "internal_interfaces": [],
        "created_at": "2026-10-27T15:04:05.123456",
        "updated_at": null,
        "tags": []
    }
}
`

// GatewayBody is an ACTIVE, dual-stack gateway as an admin sees it, with its
// instances.
const GatewayBody = `
{
    "vpn_gateway": {
        "id": "2ad3c5dc-8c7e-4ff0-8b54-4b4d4e2c7c6e",
        "project_id": "9145d91459d248b1b02fdaca97c6a75d",
        "name": "gateway",
        "description": "VPN gateway",
        "admin_state_up": true,
        "status": "ACTIVE",
        "external_ipv4": "203.0.113.10",
        "external_ipv6": "2001:db8::10",
        "ha": false,
        "flavor_id": "b8b9cc5a-0b05-4b1e-9d0a-1d6c0f5a3e01",
        "routes": [
            {
                "destination": "10.20.0.0/16",
                "nexthop": "10.0.0.254"
            }
        ],
        "availability_zone_hints": ["az1"],
        "availability_zones": [],
        "revision_number": 3,
        "external_gateways": [
            {
                "network_id": "ae3b6c4b-3e0b-4c3e-9b8a-3f3f0d1c2e11",
                "port_id": "0c4b7a8e-5d8c-4b7a-9a3e-2f1e0d9c8b71",
                "fixed_ips": [
                    {"subnet_id": "a2f1f29d-571b-4533-907f-5803ab96ead1", "ip_address": "203.0.113.10"},
                    {"subnet_id": "c6d1a3b2-0e5f-4a8b-9c7d-6e5f4a3b2c11", "ip_address": "2001:db8::10"}
                ],
                "status": "ACTIVE"
            }
        ],
        "internal_interfaces": [
            {
                "network_id": "5b7e1f2a-3c4d-4e5f-8a9b-0c1d2e3f4a51",
                "port_id": "9e8d7c6b-5a4f-4e3d-8c2b-1a0f9e8d7c61",
                "fixed_ips": [
                    {"subnet_id": "d4c3b2a1-0f9e-4d8c-b7a6-5f4e3d2c1b01", "ip_address": "10.0.0.1"}
                ],
                "status": "ACTIVE"
            }
        ],
        "instances": [
            {
                "instance_id": "7f6e5d4c-3b2a-4190-8f7e-6d5c4b3a2911",
                "management_ip": "192.168.100.5",
                "image_id": "8ab9a0a5-7bce-4c2e-9a9e-5ff04e3a2f4d"
            }
        ],
        "created_at": "2026-10-27T15:04:05.123456",
        "updated_at": "2026-10-27T15:06:00",
        "tags": []
    }
}
`

const UpdateRequest = `
{
    "vpn_gateway": {
        "name": "updatedname",
        "description": "updated gateway",
        "admin_state_up": false,
        "ha": true,
        "flavor_id": "1d2c3b4a-5e6f-4a7b-8c9d-0e1f2a3b4c51",
        "availability_zone_hints": ["az2"],
        "routes": [
            {
                "destination": "10.30.0.0/16",
                "nexthop": "10.0.0.253"
            }
        ]
    }
}
`

const UpdatedGatewayBody = `
{
    "vpn_gateway": {
        "id": "2ad3c5dc-8c7e-4ff0-8b54-4b4d4e2c7c6e",
        "project_id": "9145d91459d248b1b02fdaca97c6a75d",
        "name": "updatedname",
        "description": "updated gateway",
        "admin_state_up": false,
        "status": "ACTIVE",
        "external_ipv4": null,
        "external_ipv6": null,
        "ha": true,
        "flavor_id": "1d2c3b4a-5e6f-4a7b-8c9d-0e1f2a3b4c51",
        "routes": [
            {
                "destination": "10.30.0.0/16",
                "nexthop": "10.0.0.253"
            }
        ],
        "availability_zone_hints": ["az2"],
        "availability_zones": [],
        "revision_number": 1,
        "external_gateways": [],
        "internal_interfaces": [],
        "created_at": "2026-10-27T15:04:05",
        "updated_at": "2026-10-27T16:00:00",
        "tags": []
    }
}
`

const ClearRoutesRequest = `
{
    "vpn_gateway": {
        "routes": []
    }
}
`

const ListFirstPageBody = `
{
    "vpn_gateways": [
        {
            "id": "2ad3c5dc-8c7e-4ff0-8b54-4b4d4e2c7c6e",
            "name": "gateway",
            "status": "ACTIVE"
        }
    ],
    "vpn_gateways_links": [
        {
            "rel": "next",
            "href": "%s/v1/vpn_gateways?status=ACTIVE&limit=1&marker=2ad3c5dc-8c7e-4ff0-8b54-4b4d4e2c7c6e"
        }
    ]
}
`

const ListLastPageBody = `
{
    "vpn_gateways": [
        {
            "id": "4e5f6a7b-8c9d-4e0f-a1b2-c3d4e5f6a7b1",
            "name": "gateway2",
            "status": "ACTIVE"
        }
    ]
}
`

const InterfaceRequest = `
{
    "subnet_id": "d4c3b2a1-0f9e-4d8c-b7a6-5f4e3d2c1b01"
}
`

const AddInterfaceBody = `
{
    "id": "2ad3c5dc-8c7e-4ff0-8b54-4b4d4e2c7c6e",
    "project_id": "9145d91459d248b1b02fdaca97c6a75d",
    "network_id": null,
    "port_id": null,
    "subnet_id": "d4c3b2a1-0f9e-4d8c-b7a6-5f4e3d2c1b01",
    "subnet_ids": ["d4c3b2a1-0f9e-4d8c-b7a6-5f4e3d2c1b01"],
    "tags": []
}
`

const RemoveInterfaceBody = `
{
    "id": "2ad3c5dc-8c7e-4ff0-8b54-4b4d4e2c7c6e",
    "project_id": "9145d91459d248b1b02fdaca97c6a75d",
    "network_id": "5b7e1f2a-3c4d-4e5f-8a9b-0c1d2e3f4a51",
    "port_id": "9e8d7c6b-5a4f-4e3d-8c2b-1a0f9e8d7c61",
    "subnet_id": "d4c3b2a1-0f9e-4d8c-b7a6-5f4e3d2c1b01",
    "subnet_ids": ["d4c3b2a1-0f9e-4d8c-b7a6-5f4e3d2c1b01"],
    "tags": []
}
`

const AddExternalGatewaysRequest = `
{
    "router": {
        "external_gateways": [
            {"subnet_id": "a2f1f29d-571b-4533-907f-5803ab96ead1"},
            {"subnet_id": "c6d1a3b2-0e5f-4a8b-9c7d-6e5f4a3b2c11"}
        ]
    }
}
`

// AddExternalGatewaysBody shows the external gateway before the worker has
// created its port: no port, and fixed IPs without addresses.
const AddExternalGatewaysBody = `
{
    "vpn_gateway": {
        "id": "2ad3c5dc-8c7e-4ff0-8b54-4b4d4e2c7c6e",
        "project_id": "9145d91459d248b1b02fdaca97c6a75d",
        "name": "gateway",
        "status": "ACTIVE",
        "external_ipv4": null,
        "external_ipv6": null,
        "external_gateways": [
            {
                "network_id": "ae3b6c4b-3e0b-4c3e-9b8a-3f3f0d1c2e11",
                "port_id": null,
                "fixed_ips": [
                    {"subnet_id": "a2f1f29d-571b-4533-907f-5803ab96ead1"},
                    {"subnet_id": "c6d1a3b2-0e5f-4a8b-9c7d-6e5f4a3b2c11"}
                ],
                "status": "PENDING_CREATE"
            }
        ],
        "internal_interfaces": [],
        "created_at": "2026-10-27T15:04:05",
        "updated_at": null
    }
}
`

const RebuildRequest = `
{
    "image_id": "8ab9a0a5-7bce-4c2e-9a9e-5ff04e3a2f4d"
}
`

const PendingUpdateGatewayBody = `
{
    "vpn_gateway": {
        "id": "2ad3c5dc-8c7e-4ff0-8b54-4b4d4e2c7c6e",
        "project_id": "9145d91459d248b1b02fdaca97c6a75d",
        "name": "gateway",
        "status": "PENDING_UPDATE",
        "created_at": "2026-10-27T15:04:05",
        "updated_at": "2026-10-27T17:00:00"
    }
}
`

var (
	CreatedGateway = vpngateways.VPNGateway{
		ID:           GatewayID,
		ProjectID:    "9145d91459d248b1b02fdaca97c6a75d",
		Name:         "gateway",
		Description:  "VPN gateway",
		AdminStateUp: true,
		Status:       vpngateways.StatusPendingCreate,
		FlavorID:     "b8b9cc5a-0b05-4b1e-9d0a-1d6c0f5a3e01",
		Routes: []vpngateways.Route{
			{DestinationCIDR: "10.20.0.0/16", NextHop: "10.0.0.254"},
		},
		AvailabilityZoneHints: []string{"az1"},
		AvailabilityZones:     []string{},
		ExternalGateways:      []vpngateways.Attachment{},
		InternalInterfaces:    []vpngateways.Attachment{},
		Tags:                  []string{},
		CreatedAt:             time.Date(2026, 10, 27, 15, 4, 5, 123456000, time.UTC),
	}

	Gateway = vpngateways.VPNGateway{
		ID:           GatewayID,
		ProjectID:    "9145d91459d248b1b02fdaca97c6a75d",
		Name:         "gateway",
		Description:  "VPN gateway",
		AdminStateUp: true,
		Status:       vpngateways.StatusActive,
		ExternalIPv4: "203.0.113.10",
		ExternalIPv6: "2001:db8::10",
		FlavorID:     "b8b9cc5a-0b05-4b1e-9d0a-1d6c0f5a3e01",
		Routes: []vpngateways.Route{
			{DestinationCIDR: "10.20.0.0/16", NextHop: "10.0.0.254"},
		},
		AvailabilityZoneHints: []string{"az1"},
		AvailabilityZones:     []string{},
		RevisionNumber:        3,
		ExternalGateways: []vpngateways.Attachment{
			{
				NetworkID: "ae3b6c4b-3e0b-4c3e-9b8a-3f3f0d1c2e11",
				PortID:    "0c4b7a8e-5d8c-4b7a-9a3e-2f1e0d9c8b71",
				FixedIPs: []vpngateways.FixedIP{
					{SubnetID: "a2f1f29d-571b-4533-907f-5803ab96ead1", IPAddress: "203.0.113.10"},
					{SubnetID: "c6d1a3b2-0e5f-4a8b-9c7d-6e5f4a3b2c11", IPAddress: "2001:db8::10"},
				},
				Status: vpngateways.StatusActive,
			},
		},
		InternalInterfaces: []vpngateways.Attachment{
			{
				NetworkID: "5b7e1f2a-3c4d-4e5f-8a9b-0c1d2e3f4a51",
				PortID:    "9e8d7c6b-5a4f-4e3d-8c2b-1a0f9e8d7c61",
				FixedIPs: []vpngateways.FixedIP{
					{SubnetID: "d4c3b2a1-0f9e-4d8c-b7a6-5f4e3d2c1b01", IPAddress: "10.0.0.1"},
				},
				Status: vpngateways.StatusActive,
			},
		},
		Instances: []vpngateways.Instance{
			{
				InstanceID:   "7f6e5d4c-3b2a-4190-8f7e-6d5c4b3a2911",
				ManagementIP: "192.168.100.5",
				ImageID:      "8ab9a0a5-7bce-4c2e-9a9e-5ff04e3a2f4d",
			},
		},
		Tags:      []string{},
		CreatedAt: time.Date(2026, 10, 27, 15, 4, 5, 123456000, time.UTC),
		UpdatedAt: time.Date(2026, 10, 27, 15, 6, 0, 0, time.UTC),
	}

	UpdatedGateway = vpngateways.VPNGateway{
		ID:           GatewayID,
		ProjectID:    "9145d91459d248b1b02fdaca97c6a75d",
		Name:         "updatedname",
		Description:  "updated gateway",
		AdminStateUp: false,
		Status:       vpngateways.StatusActive,
		HA:           true,
		FlavorID:     "1d2c3b4a-5e6f-4a7b-8c9d-0e1f2a3b4c51",
		Routes: []vpngateways.Route{
			{DestinationCIDR: "10.30.0.0/16", NextHop: "10.0.0.253"},
		},
		AvailabilityZoneHints: []string{"az2"},
		AvailabilityZones:     []string{},
		RevisionNumber:        1,
		ExternalGateways:      []vpngateways.Attachment{},
		InternalInterfaces:    []vpngateways.Attachment{},
		Tags:                  []string{},
		CreatedAt:             time.Date(2026, 10, 27, 15, 4, 5, 0, time.UTC),
		UpdatedAt:             time.Date(2026, 10, 27, 16, 0, 0, 0, time.UTC),
	}

	ListedGateways = []vpngateways.VPNGateway{
		{ID: GatewayID, Name: "gateway", Status: vpngateways.StatusActive},
		{ID: "4e5f6a7b-8c9d-4e0f-a1b2-c3d4e5f6a7b1", Name: "gateway2", Status: vpngateways.StatusActive},
	}

	AddedInterface = vpngateways.InterfaceInfo{
		ID:        GatewayID,
		ProjectID: "9145d91459d248b1b02fdaca97c6a75d",
		SubnetID:  "d4c3b2a1-0f9e-4d8c-b7a6-5f4e3d2c1b01",
		SubnetIDs: []string{"d4c3b2a1-0f9e-4d8c-b7a6-5f4e3d2c1b01"},
		Tags:      []string{},
	}

	RemovedInterface = vpngateways.InterfaceInfo{
		ID:        GatewayID,
		ProjectID: "9145d91459d248b1b02fdaca97c6a75d",
		NetworkID: "5b7e1f2a-3c4d-4e5f-8a9b-0c1d2e3f4a51",
		PortID:    "9e8d7c6b-5a4f-4e3d-8c2b-1a0f9e8d7c61",
		SubnetID:  "d4c3b2a1-0f9e-4d8c-b7a6-5f4e3d2c1b01",
		SubnetIDs: []string{"d4c3b2a1-0f9e-4d8c-b7a6-5f4e3d2c1b01"},
		Tags:      []string{},
	}

	AttachingGateway = vpngateways.VPNGateway{
		ID:        GatewayID,
		ProjectID: "9145d91459d248b1b02fdaca97c6a75d",
		Name:      "gateway",
		Status:    vpngateways.StatusActive,
		ExternalGateways: []vpngateways.Attachment{
			{
				NetworkID: "ae3b6c4b-3e0b-4c3e-9b8a-3f3f0d1c2e11",
				FixedIPs: []vpngateways.FixedIP{
					{SubnetID: "a2f1f29d-571b-4533-907f-5803ab96ead1"},
					{SubnetID: "c6d1a3b2-0e5f-4a8b-9c7d-6e5f4a3b2c11"},
				},
				Status: vpngateways.StatusPendingCreate,
			},
		},
		InternalInterfaces: []vpngateways.Attachment{},
		CreatedAt:          time.Date(2026, 10, 27, 15, 4, 5, 0, time.UTC),
	}

	PendingUpdateGateway = vpngateways.VPNGateway{
		ID:        GatewayID,
		ProjectID: "9145d91459d248b1b02fdaca97c6a75d",
		Name:      "gateway",
		Status:    vpngateways.StatusPendingUpdate,
		CreatedAt: time.Date(2026, 10, 27, 15, 4, 5, 0, time.UTC),
		UpdatedAt: time.Date(2026, 10, 27, 17, 0, 0, 0, time.UTC),
	}
)

func handle(t *testing.T, fakeServer th.FakeServer, path, method, request string, status int, response string) {
	fakeServer.Mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, method)
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		if request != "" {
			th.TestJSONRequest(t, r, request)
		}

		if response != "" {
			w.Header().Add("Content-Type", "application/json")
		}
		w.WriteHeader(status)
		fmt.Fprint(w, response)
	})
}

func HandleGatewayCreationSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	handle(t, fakeServer, "/v1/vpn_gateways", "POST", CreateRequest, http.StatusCreated, CreatedGatewayBody)
}

func HandleGatewayGetSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	handle(t, fakeServer, "/v1/vpn_gateways/"+GatewayID, "GET", "", http.StatusOK, GatewayBody)
}

func HandleGatewayUpdateSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	handle(t, fakeServer, "/v1/vpn_gateways/"+GatewayID, "PUT", UpdateRequest, http.StatusOK, UpdatedGatewayBody)
}

func HandleGatewayClearRoutesSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	handle(t, fakeServer, "/v1/vpn_gateways/"+GatewayID, "PUT", ClearRoutesRequest, http.StatusOK, UpdatedGatewayBody)
}

func HandleGatewayDeletionSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	handle(t, fakeServer, "/v1/vpn_gateways/"+GatewayID, "DELETE", "", http.StatusNoContent, "")
}

func HandleAddInterfaceSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	handle(t, fakeServer, "/v1/vpn_gateways/"+GatewayID+"/add_router_interface", "PUT", InterfaceRequest, http.StatusOK, AddInterfaceBody)
}

func HandleRemoveInterfaceSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	handle(t, fakeServer, "/v1/vpn_gateways/"+GatewayID+"/remove_router_interface", "PUT", InterfaceRequest, http.StatusOK, RemoveInterfaceBody)
}

func HandleAddExternalGatewaysSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	handle(t, fakeServer, "/v1/vpn_gateways/"+GatewayID+"/add_external_gateways", "PUT", AddExternalGatewaysRequest, http.StatusOK, AddExternalGatewaysBody)
}

func HandleRebuildSuccessfully(t *testing.T, fakeServer th.FakeServer, request string) {
	handle(t, fakeServer, "/v1/vpn_gateways/"+GatewayID+"/rebuild", "PUT", request, http.StatusAccepted, PendingUpdateGatewayBody)
}

// HandleResyncSuccessfully checks that a resync carries no body.
func HandleResyncSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/v1/vpn_gateways/"+GatewayID+"/resync", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		body, err := io.ReadAll(r.Body)
		th.AssertNoErr(t, err)
		th.AssertEquals(t, "", string(body))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprint(w, PendingUpdateGatewayBody)
	})
}

// HandleGatewayListSuccessfully serves two pages: the first answers the
// ListOpts query, the second its next link.
func HandleGatewayListSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/v1/vpn_gateways", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		if err := r.ParseForm(); err != nil {
			t.Errorf("Failed to parse request form %v", err)
		}
		switch r.Form.Get("marker") {
		case "":
			th.CheckDeepEquals(t, url.Values{
				"project_id":   {"9145d91459d248b1b02fdaca97c6a75d"},
				"name":         {"gateway"},
				"status":       {"ACTIVE"},
				"flavor_id":    {"b8b9cc5a-0b05-4b1e-9d0a-1d6c0f5a3e01"},
				"limit":        {"1"},
				"page_reverse": {"true"},
				"sort_key":     {"name"},
				"sort_dir":     {"desc"},
				"fields":       {"id", "name", "status"},
			}, r.Form)
			fmt.Fprintf(w, ListFirstPageBody, fakeServer.Server.URL)
		case GatewayID:
			fmt.Fprint(w, ListLastPageBody)
		default:
			t.Fatalf("Unexpected marker: [%s]", r.Form.Get("marker"))
		}
	})
}
