package testing

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/vpngateway/v1/siteconnections"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

const ConnectionID = "5c561d9d-eaea-45f6-ae3e-08d1a7080828"

const CreateRequest = `
{
    "ipsec_site_connection": {
        "project_id": "9145d91459d248b1b02fdaca97c6a75d",
        "name": "tunnel",
        "description": "IPsec site connection",
        "ipsec_site_id": "9d1e6c2a-4b3f-4a5e-8d7c-6b5a4f3e2d1c",
        "peer_address": "198.51.100.20",
        "peer_id": "vpn.peer.example.com",
        "local_id": "@gateway.example.com",
        "auth_mode": "psk",
        "psk": "a-long-pre-shared-key",
        "priority": 1,
        "admin_state_up": true,
        "ikepolicy_id": "f2b08c1e-aa81-4668-8ae1-1401bcb0576c",
        "ipsecpolicy_id": "e3b0c442-98fc-4c14-9afb-f4c8996fb924",
        "dpd": {
            "action": "restart",
            "interval": 10,
            "timeout": 60
        },
        "mtu": 1380,
        "initiator": "response-only",
        "tunnel_address": "169.254.10.1/30",
        "bgp_peer_address": "169.254.10.2"
    }
}
`

// ConnectionBody never carries the PSK: psk_set says one is set.
const ConnectionBody = `
{
    "ipsec_site_connection": {
        "id": "5c561d9d-eaea-45f6-ae3e-08d1a7080828",
        "project_id": "9145d91459d248b1b02fdaca97c6a75d",
        "name": "tunnel",
        "description": "IPsec site connection",
        "ipsec_site_id": "9d1e6c2a-4b3f-4a5e-8d7c-6b5a4f3e2d1c",
        "peer_address": "198.51.100.20",
        "peer_id": "vpn.peer.example.com",
        "local_id": "@gateway.example.com",
        "auth_mode": "psk",
        "psk_set": true,
        "priority": 1,
        "ikepolicy_id": "f2b08c1e-aa81-4668-8ae1-1401bcb0576c",
        "ipsecpolicy_id": "e3b0c442-98fc-4c14-9afb-f4c8996fb924",
        "dpd": {
            "action": "restart",
            "interval": 10,
            "timeout": 60
        },
        "mtu": 1380,
        "initiator": "response-only",
        "tunnel_address": "169.254.10.1/30",
        "bgp_peer_address": "169.254.10.2",
        "admin_state_up": true,
        "status": "PENDING_CREATE",
        "created_at": "2026-10-27T15:04:05.123456",
        "updated_at": null
    }
}
`

const UpdateRequest = `
{
    "ipsec_site_connection": {
        "name": "updatedname",
        "description": "updated connection",
        "peer_address": "198.51.100.30",
        "peer_id": "198.51.100.30",
        "local_id": "203.0.113.10",
        "psk": "another-pre-shared-key",
        "priority": 2,
        "admin_state_up": false,
        "ikepolicy_id": "0f1e2d3c-4b5a-4968-8778-695a4b3c2d1e",
        "ipsecpolicy_id": "1e2d3c4b-5a69-4877-9685-a4b3c2d1e0f1",
        "dpd": {
            "action": "hold"
        },
        "mtu": 1400,
        "initiator": "bi-directional",
        "tunnel_address": "169.254.20.1/30",
        "bgp_peer_address": "169.254.20.2"
    }
}
`

// UpdatedConnectionBody shows the DPD parameters the update left unset
// keeping their value.
const UpdatedConnectionBody = `
{
    "ipsec_site_connection": {
        "id": "5c561d9d-eaea-45f6-ae3e-08d1a7080828",
        "project_id": "9145d91459d248b1b02fdaca97c6a75d",
        "name": "updatedname",
        "description": "updated connection",
        "ipsec_site_id": "9d1e6c2a-4b3f-4a5e-8d7c-6b5a4f3e2d1c",
        "peer_address": "198.51.100.30",
        "peer_id": "198.51.100.30",
        "local_id": "203.0.113.10",
        "auth_mode": "psk",
        "psk_set": true,
        "priority": 2,
        "ikepolicy_id": "0f1e2d3c-4b5a-4968-8778-695a4b3c2d1e",
        "ipsecpolicy_id": "1e2d3c4b-5a69-4877-9685-a4b3c2d1e0f1",
        "dpd": {
            "action": "hold",
            "interval": 10,
            "timeout": 60
        },
        "mtu": 1400,
        "initiator": "bi-directional",
        "tunnel_address": "169.254.20.1/30",
        "bgp_peer_address": "169.254.20.2",
        "admin_state_up": false,
        "status": "PENDING_UPDATE",
        "created_at": "2026-10-27T15:04:05",
        "updated_at": "2026-10-27T16:00:00"
    }
}
`

const ListFirstPageBody = `
{
    "ipsec_site_connections": [
        {
            "id": "5c561d9d-eaea-45f6-ae3e-08d1a7080828",
            "name": "tunnel",
            "status": "ACTIVE"
        }
    ],
    "ipsec_site_connections_links": [
        {
            "rel": "next",
            "href": "%s/v1/ipsec_site_connections?limit=1&marker=5c561d9d-eaea-45f6-ae3e-08d1a7080828"
        }
    ]
}
`

const ListLastPageBody = `
{
    "ipsec_site_connections": [
        {
            "id": "6d7e8f90-a1b2-4c3d-8e4f-5a6b7c8d9e0f",
            "name": "tunnel2",
            "status": "ACTIVE"
        }
    ]
}
`

var (
	Connection = siteconnections.Connection{
		ID:            ConnectionID,
		ProjectID:     "9145d91459d248b1b02fdaca97c6a75d",
		Name:          "tunnel",
		Description:   "IPsec site connection",
		IPsecSiteID:   "9d1e6c2a-4b3f-4a5e-8d7c-6b5a4f3e2d1c",
		PeerAddress:   "198.51.100.20",
		PeerID:        "vpn.peer.example.com",
		LocalID:       "@gateway.example.com",
		AuthMode:      "psk",
		PSKSet:        true,
		Priority:      1,
		IKEPolicyID:   "f2b08c1e-aa81-4668-8ae1-1401bcb0576c",
		IPsecPolicyID: "e3b0c442-98fc-4c14-9afb-f4c8996fb924",
		DPD: siteconnections.DPD{
			Action:   "restart",
			Interval: 10,
			Timeout:  60,
		},
		MTU:            1380,
		Initiator:      "response-only",
		TunnelAddress:  "169.254.10.1/30",
		BGPPeerAddress: "169.254.10.2",
		AdminStateUp:   true,
		Status:         siteconnections.StatusPendingCreate,
		CreatedAt:      time.Date(2026, 10, 27, 15, 4, 5, 123456000, time.UTC),
	}

	UpdatedConnection = siteconnections.Connection{
		ID:            ConnectionID,
		ProjectID:     "9145d91459d248b1b02fdaca97c6a75d",
		Name:          "updatedname",
		Description:   "updated connection",
		IPsecSiteID:   "9d1e6c2a-4b3f-4a5e-8d7c-6b5a4f3e2d1c",
		PeerAddress:   "198.51.100.30",
		PeerID:        "198.51.100.30",
		LocalID:       "203.0.113.10",
		AuthMode:      "psk",
		PSKSet:        true,
		Priority:      2,
		IKEPolicyID:   "0f1e2d3c-4b5a-4968-8778-695a4b3c2d1e",
		IPsecPolicyID: "1e2d3c4b-5a69-4877-9685-a4b3c2d1e0f1",
		DPD: siteconnections.DPD{
			Action:   "hold",
			Interval: 10,
			Timeout:  60,
		},
		MTU:            1400,
		Initiator:      "bi-directional",
		TunnelAddress:  "169.254.20.1/30",
		BGPPeerAddress: "169.254.20.2",
		AdminStateUp:   false,
		Status:         siteconnections.StatusPendingUpdate,
		CreatedAt:      time.Date(2026, 10, 27, 15, 4, 5, 0, time.UTC),
		UpdatedAt:      time.Date(2026, 10, 27, 16, 0, 0, 0, time.UTC),
	}

	ListedConnections = []siteconnections.Connection{
		{ID: ConnectionID, Name: "tunnel", Status: siteconnections.StatusActive},
		{ID: "6d7e8f90-a1b2-4c3d-8e4f-5a6b7c8d9e0f", Name: "tunnel2", Status: siteconnections.StatusActive},
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

func HandleConnectionCreationSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	handle(t, fakeServer, "/v1/ipsec_site_connections", "POST", CreateRequest, http.StatusCreated, ConnectionBody)
}

func HandleConnectionGetSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	handle(t, fakeServer, "/v1/ipsec_site_connections/"+ConnectionID, "GET", "", http.StatusOK, ConnectionBody)
}

func HandleConnectionUpdateSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	handle(t, fakeServer, "/v1/ipsec_site_connections/"+ConnectionID, "PUT", UpdateRequest, http.StatusOK, UpdatedConnectionBody)
}

func HandleConnectionDeletionSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	handle(t, fakeServer, "/v1/ipsec_site_connections/"+ConnectionID, "DELETE", "", http.StatusNoContent, "")
}

// HandleConnectionListSuccessfully serves two pages: the first answers the
// ListOpts query, the second its next link.
func HandleConnectionListSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/v1/ipsec_site_connections", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		if err := r.ParseForm(); err != nil {
			t.Errorf("Failed to parse request form %v", err)
		}
		switch r.Form.Get("marker") {
		case "":
			th.CheckDeepEquals(t, url.Values{
				"project_id":    {"9145d91459d248b1b02fdaca97c6a75d"},
				"name":          {"tunnel"},
				"status":        {"ACTIVE"},
				"ipsec_site_id": {"9d1e6c2a-4b3f-4a5e-8d7c-6b5a4f3e2d1c"},
				"limit":         {"1"},
				"page_reverse":  {"true"},
				"sort_key":      {"priority"},
				"sort_dir":      {"asc"},
				"fields":        {"id", "name", "status"},
			}, r.Form)
			fmt.Fprintf(w, ListFirstPageBody, fakeServer.Server.URL)
		case ConnectionID:
			fmt.Fprint(w, ListLastPageBody)
		default:
			t.Fatalf("Unexpected marker: [%s]", r.Form.Get("marker"))
		}
	})
}
