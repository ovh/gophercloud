package testing

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/vpngateway/v1/endpointgroups"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

const GroupID = "6ecd9cf3-ca64-46c7-863f-f2eb1b9e838a"

const CreateRequest = `
{
    "endpoint_group": {
        "project_id": "9145d91459d248b1b02fdaca97c6a75d",
        "name": "peers",
        "description": "peer networks",
        "type": "cidr",
        "cidrs": [
            "10.2.0.0/24",
            "fd00:2::/64"
        ]
    }
}
`

const GroupBody = `
{
    "endpoint_group": {
        "id": "6ecd9cf3-ca64-46c7-863f-f2eb1b9e838a",
        "project_id": "9145d91459d248b1b02fdaca97c6a75d",
        "name": "peers",
        "description": "peer networks",
        "type": "cidr",
        "cidrs": [
            "10.2.0.0/24",
            "fd00:2::/64"
        ],
        "created_at": "2026-10-27T15:04:05.123456",
        "updated_at": null
    }
}
`

const UpdateRequest = `
{
    "endpoint_group": {
        "name": "updatedname",
        "description": "updated endpoint group",
        "cidrs": [
            "10.4.0.0/16"
        ]
    }
}
`

const UpdatedGroupBody = `
{
    "endpoint_group": {
        "id": "6ecd9cf3-ca64-46c7-863f-f2eb1b9e838a",
        "project_id": "9145d91459d248b1b02fdaca97c6a75d",
        "name": "updatedname",
        "description": "updated endpoint group",
        "type": "cidr",
        "cidrs": [
            "10.4.0.0/16"
        ],
        "created_at": "2026-10-27T15:04:05",
        "updated_at": "2026-10-27T16:00:00"
    }
}
`

const ListFirstPageBody = `
{
    "endpoint_groups": [
        {
            "id": "6ecd9cf3-ca64-46c7-863f-f2eb1b9e838a",
            "project_id": "9145d91459d248b1b02fdaca97c6a75d",
            "name": "peers",
            "description": "peer networks",
            "type": "cidr",
            "cidrs": [
                "10.2.0.0/24",
                "fd00:2::/64"
            ],
            "created_at": "2026-10-27T15:04:05.123456",
            "updated_at": null
        }
    ],
    "endpoint_groups_links": [
        {
            "rel": "next",
            "href": "%s/v1/endpoint_groups?name=peers&limit=1&marker=6ecd9cf3-ca64-46c7-863f-f2eb1b9e838a"
        }
    ]
}
`

const ListLastPageBody = `
{
    "endpoint_groups": [
        {
            "id": "5c561d9d-eaea-45f6-ae3e-08d1a7080828",
            "project_id": "9145d91459d248b1b02fdaca97c6a75d",
            "name": "peers",
            "description": null,
            "type": "cidr",
            "cidrs": [
                "192.168.0.0/24"
            ],
            "created_at": "2026-10-27T15:10:00",
            "updated_at": null
        }
    ],
    "endpoint_groups_links": [
        {
            "rel": "previous",
            "href": "%s/v1/endpoint_groups?name=peers&limit=1&marker=5c561d9d-eaea-45f6-ae3e-08d1a7080828&page_reverse=True"
        }
    ]
}
`

var (
	Group = endpointgroups.EndpointGroup{
		ID:          GroupID,
		ProjectID:   "9145d91459d248b1b02fdaca97c6a75d",
		Name:        "peers",
		Description: "peer networks",
		Type:        "cidr",
		CIDRs:       []string{"10.2.0.0/24", "fd00:2::/64"},
		CreatedAt:   time.Date(2026, 10, 27, 15, 4, 5, 123456000, time.UTC),
	}

	UpdatedGroup = endpointgroups.EndpointGroup{
		ID:          GroupID,
		ProjectID:   "9145d91459d248b1b02fdaca97c6a75d",
		Name:        "updatedname",
		Description: "updated endpoint group",
		Type:        "cidr",
		CIDRs:       []string{"10.4.0.0/16"},
		CreatedAt:   time.Date(2026, 10, 27, 15, 4, 5, 0, time.UTC),
		UpdatedAt:   time.Date(2026, 10, 27, 16, 0, 0, 0, time.UTC),
	}

	SecondGroup = endpointgroups.EndpointGroup{
		ID:        "5c561d9d-eaea-45f6-ae3e-08d1a7080828",
		ProjectID: "9145d91459d248b1b02fdaca97c6a75d",
		Name:      "peers",
		Type:      "cidr",
		CIDRs:     []string{"192.168.0.0/24"},
		CreatedAt: time.Date(2026, 10, 27, 15, 10, 0, 0, time.UTC),
	}
)

func HandleGroupCreationSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/v1/endpoint_groups", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, CreateRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, GroupBody)
	})
}

func HandleGroupGetSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/v1/endpoint_groups/"+GroupID, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, GroupBody)
	})
}

func HandleGroupUpdateSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/v1/endpoint_groups/"+GroupID, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, UpdateRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, UpdatedGroupBody)
	})
}

func HandleGroupDeletionSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/v1/endpoint_groups/"+GroupID, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.WriteHeader(http.StatusNoContent)
	})
}

// HandleGroupListSuccessfully serves two pages: the first answers the
// ListOpts query, the second its next link.
func HandleGroupListSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/v1/endpoint_groups", func(w http.ResponseWriter, r *http.Request) {
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
				"name":         {"peers"},
				"limit":        {"1"},
				"page_reverse": {"true"},
				"sort_key":     {"created_at"},
				"sort_dir":     {"asc"},
				"fields":       {"id", "cidrs"},
			}, r.Form)
			fmt.Fprintf(w, ListFirstPageBody, fakeServer.Server.URL)
		case GroupID:
			fmt.Fprintf(w, ListLastPageBody, fakeServer.Server.URL)
		default:
			t.Fatalf("Unexpected marker: [%s]", r.Form.Get("marker"))
		}
	})
}
