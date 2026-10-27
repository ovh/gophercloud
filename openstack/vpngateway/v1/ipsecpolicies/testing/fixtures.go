package testing

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/vpngateway/v1/ipsecpolicies"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

const PolicyID = "e3b0c442-98fc-4c14-9afb-f4c8996fb924"

const CreateRequest = `
{
    "ipsec_policy": {
        "project_id": "9145d91459d248b1b02fdaca97c6a75d",
        "name": "policy",
        "description": "IPsec policy",
        "transform_protocol": "esp",
        "encryption_algorithm": "aes-256-gcm-16",
        "encapsulation_mode": "tunnel",
        "pfs": "group20",
        "lifetime": {
            "units": "seconds",
            "value": 1800
        }
    }
}
`

// PolicyBody is an AES-GCM policy: its auth_algorithm is null.
const PolicyBody = `
{
    "ipsec_policy": {
        "id": "e3b0c442-98fc-4c14-9afb-f4c8996fb924",
        "project_id": "9145d91459d248b1b02fdaca97c6a75d",
        "name": "policy",
        "description": "IPsec policy",
        "transform_protocol": "esp",
        "auth_algorithm": null,
        "encryption_algorithm": "aes-256-gcm-16",
        "encapsulation_mode": "tunnel",
        "pfs": "group20",
        "lifetime": {
            "units": "seconds",
            "value": 1800
        },
        "created_at": "2026-10-27T15:04:05.123456",
        "updated_at": null
    }
}
`

const UpdateRequest = `
{
    "ipsec_policy": {
        "name": "updatedname",
        "description": "updated policy",
        "transform_protocol": "esp",
        "auth_algorithm": "sha256",
        "encryption_algorithm": "aes-256",
        "encapsulation_mode": "tunnel",
        "pfs": "group14",
        "lifetime": {
            "units": "seconds",
            "value": 3600
        }
    }
}
`

const UpdatedPolicyBody = `
{
    "ipsec_policy": {
        "id": "e3b0c442-98fc-4c14-9afb-f4c8996fb924",
        "project_id": "9145d91459d248b1b02fdaca97c6a75d",
        "name": "updatedname",
        "description": "updated policy",
        "transform_protocol": "esp",
        "auth_algorithm": "sha256",
        "encryption_algorithm": "aes-256",
        "encapsulation_mode": "tunnel",
        "pfs": "group14",
        "lifetime": {
            "units": "seconds",
            "value": 3600
        },
        "created_at": "2026-10-27T15:04:05",
        "updated_at": "2026-10-27T16:00:00"
    }
}
`

const ListFirstPageBody = `
{
    "ipsec_policies": [
        {
            "id": "e3b0c442-98fc-4c14-9afb-f4c8996fb924",
            "project_id": "9145d91459d248b1b02fdaca97c6a75d",
            "name": "policy",
            "description": "IPsec policy",
            "transform_protocol": "esp",
            "auth_algorithm": null,
            "encryption_algorithm": "aes-256-gcm-16",
            "encapsulation_mode": "tunnel",
            "pfs": "group20",
            "lifetime": {
                "units": "seconds",
                "value": 1800
            },
            "created_at": "2026-10-27T15:04:05.123456",
            "updated_at": null
        }
    ],
    "ipsec_policies_links": [
        {
            "rel": "next",
            "href": "%s/v1/ipsec_policies?name=policy&limit=1&marker=e3b0c442-98fc-4c14-9afb-f4c8996fb924"
        }
    ]
}
`

const ListLastPageBody = `
{
    "ipsec_policies": [
        {
            "id": "5c561d9d-eaea-45f6-ae3e-08d1a7080828",
            "project_id": "9145d91459d248b1b02fdaca97c6a75d",
            "name": "policy",
            "description": null,
            "transform_protocol": "esp",
            "auth_algorithm": "sha256",
            "encryption_algorithm": "aes-256",
            "encapsulation_mode": "tunnel",
            "pfs": "group14",
            "lifetime": {
                "units": "seconds",
                "value": 3600
            },
            "created_at": "2026-10-27T15:10:00",
            "updated_at": null
        }
    ],
    "ipsec_policies_links": [
        {
            "rel": "previous",
            "href": "%s/v1/ipsec_policies?name=policy&limit=1&marker=5c561d9d-eaea-45f6-ae3e-08d1a7080828&page_reverse=True"
        }
    ]
}
`

var sha256 = "sha256"

var (
	Policy = ipsecpolicies.Policy{
		ID:                  PolicyID,
		ProjectID:           "9145d91459d248b1b02fdaca97c6a75d",
		Name:                "policy",
		Description:         "IPsec policy",
		TransformProtocol:   "esp",
		EncryptionAlgorithm: "aes-256-gcm-16",
		EncapsulationMode:   "tunnel",
		PFS:                 "group20",
		Lifetime:            ipsecpolicies.Lifetime{Units: "seconds", Value: 1800},
		CreatedAt:           time.Date(2026, 10, 27, 15, 4, 5, 123456000, time.UTC),
	}

	UpdatedPolicy = ipsecpolicies.Policy{
		ID:                  PolicyID,
		ProjectID:           "9145d91459d248b1b02fdaca97c6a75d",
		Name:                "updatedname",
		Description:         "updated policy",
		TransformProtocol:   "esp",
		AuthAlgorithm:       &sha256,
		EncryptionAlgorithm: "aes-256",
		EncapsulationMode:   "tunnel",
		PFS:                 "group14",
		Lifetime:            ipsecpolicies.Lifetime{Units: "seconds", Value: 3600},
		CreatedAt:           time.Date(2026, 10, 27, 15, 4, 5, 0, time.UTC),
		UpdatedAt:           time.Date(2026, 10, 27, 16, 0, 0, 0, time.UTC),
	}

	SecondPolicy = ipsecpolicies.Policy{
		ID:                  "5c561d9d-eaea-45f6-ae3e-08d1a7080828",
		ProjectID:           "9145d91459d248b1b02fdaca97c6a75d",
		Name:                "policy",
		TransformProtocol:   "esp",
		AuthAlgorithm:       &sha256,
		EncryptionAlgorithm: "aes-256",
		EncapsulationMode:   "tunnel",
		PFS:                 "group14",
		Lifetime:            ipsecpolicies.Lifetime{Units: "seconds", Value: 3600},
		CreatedAt:           time.Date(2026, 10, 27, 15, 10, 0, 0, time.UTC),
	}
)

func HandlePolicyCreationSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/v1/ipsec_policies", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, CreateRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, PolicyBody)
	})
}

func HandlePolicyGetSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/v1/ipsec_policies/"+PolicyID, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, PolicyBody)
	})
}

func HandlePolicyUpdateSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/v1/ipsec_policies/"+PolicyID, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, UpdateRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, UpdatedPolicyBody)
	})
}

func HandlePolicyDeletionSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/v1/ipsec_policies/"+PolicyID, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.WriteHeader(http.StatusNoContent)
	})
}

// HandlePolicyListSuccessfully serves two pages: the first answers the
// ListOpts query, the second its next link.
func HandlePolicyListSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/v1/ipsec_policies", func(w http.ResponseWriter, r *http.Request) {
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
				"name":         {"policy"},
				"limit":        {"1"},
				"page_reverse": {"true"},
				"sort_key":     {"created_at"},
				"sort_dir":     {"asc"},
				"fields":       {"id", "name"},
			}, r.Form)
			fmt.Fprintf(w, ListFirstPageBody, fakeServer.Server.URL)
		case PolicyID:
			fmt.Fprintf(w, ListLastPageBody, fakeServer.Server.URL)
		default:
			t.Fatalf("Unexpected marker: [%s]", r.Form.Get("marker"))
		}
	})
}
