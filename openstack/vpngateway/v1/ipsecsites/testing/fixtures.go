package testing

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/vpngateway/v1/ipsecsites"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

const SiteID = "9d1e6c2a-4b3f-4a5e-8d7c-6b5a4f3e2d1c"

const CreateRequest = `
{
    "ipsec_site": {
        "project_id": "9145d91459d248b1b02fdaca97c6a75d",
        "name": "site",
        "description": "IPsec site",
        "vpn_gateway_id": "2ad3c5dc-8c7e-4ff0-8b54-4b4d4e2c7c6e",
        "routing": "bgp",
        "local_asn": 4200000001,
        "peer_asn": 64512,
        "ecmp_inbound": true,
        "local_ep_group_id": "6ecd9cf3-ca64-46c7-863f-f2eb1b9e838a",
        "peer_ep_group_id": "3e1f6a2b-7c8d-4e9f-a0b1-c2d3e4f5a6b7"
    }
}
`

const SiteBody = `
{
    "ipsec_site": {
        "id": "9d1e6c2a-4b3f-4a5e-8d7c-6b5a4f3e2d1c",
        "project_id": "9145d91459d248b1b02fdaca97c6a75d",
        "name": "site",
        "description": "IPsec site",
        "vpn_gateway_id": "2ad3c5dc-8c7e-4ff0-8b54-4b4d4e2c7c6e",
        "routing": "bgp",
        "local_asn": 4200000001,
        "peer_asn": 64512,
        "ecmp_inbound": true,
        "local_ep_group_id": "6ecd9cf3-ca64-46c7-863f-f2eb1b9e838a",
        "peer_ep_group_id": "3e1f6a2b-7c8d-4e9f-a0b1-c2d3e4f5a6b7",
        "created_at": "2026-10-27T15:04:05.123456",
        "updated_at": null
    }
}
`

const UpdateRequest = `
{
    "ipsec_site": {
        "name": "updatedname",
        "description": "updated site"
    }
}
`

// UpdatedSiteBody is a static site: its ASNs are null.
const UpdatedSiteBody = `
{
    "ipsec_site": {
        "id": "9d1e6c2a-4b3f-4a5e-8d7c-6b5a4f3e2d1c",
        "project_id": "9145d91459d248b1b02fdaca97c6a75d",
        "name": "updatedname",
        "description": "updated site",
        "vpn_gateway_id": "2ad3c5dc-8c7e-4ff0-8b54-4b4d4e2c7c6e",
        "routing": "static",
        "local_asn": null,
        "peer_asn": null,
        "ecmp_inbound": false,
        "local_ep_group_id": "6ecd9cf3-ca64-46c7-863f-f2eb1b9e838a",
        "peer_ep_group_id": "3e1f6a2b-7c8d-4e9f-a0b1-c2d3e4f5a6b7",
        "created_at": "2026-10-27T15:04:05",
        "updated_at": "2026-10-27T16:00:00"
    }
}
`

const ListFirstPageBody = `
{
    "ipsec_sites": [
        {
            "id": "9d1e6c2a-4b3f-4a5e-8d7c-6b5a4f3e2d1c",
            "name": "site"
        }
    ],
    "ipsec_sites_links": [
        {
            "rel": "next",
            "href": "%s/v1/ipsec_sites?limit=1&marker=9d1e6c2a-4b3f-4a5e-8d7c-6b5a4f3e2d1c"
        }
    ]
}
`

const ListLastPageBody = `
{
    "ipsec_sites": [
        {
            "id": "0a1b2c3d-4e5f-4a6b-8c7d-9e0f1a2b3c4d",
            "name": "site2"
        }
    ]
}
`

var (
	localASN = int64(4200000001)
	peerASN  = int64(64512)

	Site = ipsecsites.Site{
		ID:             SiteID,
		ProjectID:      "9145d91459d248b1b02fdaca97c6a75d",
		Name:           "site",
		Description:    "IPsec site",
		VPNGatewayID:   "2ad3c5dc-8c7e-4ff0-8b54-4b4d4e2c7c6e",
		Routing:        "bgp",
		LocalASN:       &localASN,
		PeerASN:        &peerASN,
		ECMPInbound:    true,
		LocalEPGroupID: "6ecd9cf3-ca64-46c7-863f-f2eb1b9e838a",
		PeerEPGroupID:  "3e1f6a2b-7c8d-4e9f-a0b1-c2d3e4f5a6b7",
		CreatedAt:      time.Date(2026, 10, 27, 15, 4, 5, 123456000, time.UTC),
	}

	UpdatedSite = ipsecsites.Site{
		ID:             SiteID,
		ProjectID:      "9145d91459d248b1b02fdaca97c6a75d",
		Name:           "updatedname",
		Description:    "updated site",
		VPNGatewayID:   "2ad3c5dc-8c7e-4ff0-8b54-4b4d4e2c7c6e",
		Routing:        "static",
		LocalEPGroupID: "6ecd9cf3-ca64-46c7-863f-f2eb1b9e838a",
		PeerEPGroupID:  "3e1f6a2b-7c8d-4e9f-a0b1-c2d3e4f5a6b7",
		CreatedAt:      time.Date(2026, 10, 27, 15, 4, 5, 0, time.UTC),
		UpdatedAt:      time.Date(2026, 10, 27, 16, 0, 0, 0, time.UTC),
	}

	ListedSites = []ipsecsites.Site{
		{ID: SiteID, Name: "site"},
		{ID: "0a1b2c3d-4e5f-4a6b-8c7d-9e0f1a2b3c4d", Name: "site2"},
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

func HandleSiteCreationSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	handle(t, fakeServer, "/v1/ipsec_sites", "POST", CreateRequest, http.StatusCreated, SiteBody)
}

func HandleSiteGetSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	handle(t, fakeServer, "/v1/ipsec_sites/"+SiteID, "GET", "", http.StatusOK, SiteBody)
}

func HandleSiteUpdateSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	handle(t, fakeServer, "/v1/ipsec_sites/"+SiteID, "PUT", UpdateRequest, http.StatusOK, UpdatedSiteBody)
}

func HandleSiteDeletionSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	handle(t, fakeServer, "/v1/ipsec_sites/"+SiteID, "DELETE", "", http.StatusNoContent, "")
}

// HandleSiteListSuccessfully serves two pages: the first answers the
// ListOpts query, the second its next link.
func HandleSiteListSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/v1/ipsec_sites", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		if err := r.ParseForm(); err != nil {
			t.Errorf("Failed to parse request form %v", err)
		}
		switch r.Form.Get("marker") {
		case "":
			th.CheckDeepEquals(t, url.Values{
				"project_id":     {"9145d91459d248b1b02fdaca97c6a75d"},
				"name":           {"site"},
				"vpn_gateway_id": {"2ad3c5dc-8c7e-4ff0-8b54-4b4d4e2c7c6e"},
				"limit":          {"1"},
				"page_reverse":   {"true"},
				"sort_key":       {"routing"},
				"sort_dir":       {"asc"},
				"fields":         {"id", "name"},
			}, r.Form)
			fmt.Fprintf(w, ListFirstPageBody, fakeServer.Server.URL)
		case SiteID:
			fmt.Fprint(w, ListLastPageBody)
		default:
			t.Fatalf("Unexpected marker: [%s]", r.Form.Get("marker"))
		}
	})
}
