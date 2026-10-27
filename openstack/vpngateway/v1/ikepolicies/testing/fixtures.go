package testing

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/vpngateway/v1/ikepolicies"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

const PolicyID = "f2b08c1e-aa81-4668-8ae1-1401bcb0576c"

const CreateRequest = `
{
    "ike_policy": {
        "project_id": "9145d91459d248b1b02fdaca97c6a75d",
        "name": "policy",
        "description": "IKE policy",
        "auth_algorithm": "sha384",
        "encryption_algorithm": "aes-256-gcm-16",
        "phase1_negotiation_mode": "main",
        "ike_version": "v2",
        "pfs": "group19",
        "lifetime": {
            "units": "seconds",
            "value": 7200
        }
    }
}
`

// PolicyBody carries a timestamp with microseconds and a null updated_at, as
// the service renders a freshly created row.
const PolicyBody = `
{
    "ike_policy": {
        "id": "f2b08c1e-aa81-4668-8ae1-1401bcb0576c",
        "project_id": "9145d91459d248b1b02fdaca97c6a75d",
        "name": "policy",
        "description": "IKE policy",
        "auth_algorithm": "sha384",
        "encryption_algorithm": "aes-256-gcm-16",
        "phase1_negotiation_mode": "main",
        "ike_version": "v2",
        "pfs": "group19",
        "lifetime": {
            "units": "seconds",
            "value": 7200
        },
        "created_at": "2026-10-27T15:04:05.123456",
        "updated_at": null
    }
}
`

const UpdateRequest = `
{
    "ike_policy": {
        "name": "updatedname",
        "description": "updated policy",
        "auth_algorithm": "sha512",
        "encryption_algorithm": "aes-128",
        "phase1_negotiation_mode": "main",
        "ike_version": "v2",
        "pfs": "group31",
        "lifetime": {
            "units": "seconds",
            "value": 3600
        }
    }
}
`

// UpdatedPolicyBody carries timestamps without a fraction, as the service
// renders a datetime whose microseconds are zero.
const UpdatedPolicyBody = `
{
    "ike_policy": {
        "id": "f2b08c1e-aa81-4668-8ae1-1401bcb0576c",
        "project_id": "9145d91459d248b1b02fdaca97c6a75d",
        "name": "updatedname",
        "description": "updated policy",
        "auth_algorithm": "sha512",
        "encryption_algorithm": "aes-128",
        "phase1_negotiation_mode": "main",
        "ike_version": "v2",
        "pfs": "group31",
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
    "ike_policies": [
        {
            "id": "f2b08c1e-aa81-4668-8ae1-1401bcb0576c",
            "project_id": "9145d91459d248b1b02fdaca97c6a75d",
            "name": "policy",
            "description": "IKE policy",
            "auth_algorithm": "sha384",
            "encryption_algorithm": "aes-256-gcm-16",
            "phase1_negotiation_mode": "main",
            "ike_version": "v2",
            "pfs": "group19",
            "lifetime": {
                "units": "seconds",
                "value": 7200
            },
            "created_at": "2026-10-27T15:04:05.123456",
            "updated_at": null
        }
    ],
    "ike_policies_links": [
        {
            "rel": "next",
            "href": "%s/v1/ike_policies?name=policy&limit=1&marker=f2b08c1e-aa81-4668-8ae1-1401bcb0576c"
        }
    ]
}
`

const ListLastPageBody = `
{
    "ike_policies": [
        {
            "id": "5c561d9d-eaea-45f6-ae3e-08d1a7080828",
            "project_id": "9145d91459d248b1b02fdaca97c6a75d",
            "name": "policy",
            "description": null,
            "auth_algorithm": "sha256",
            "encryption_algorithm": "aes-256",
            "phase1_negotiation_mode": "main",
            "ike_version": "v2",
            "pfs": "group14",
            "lifetime": {
                "units": "seconds",
                "value": 10800
            },
            "created_at": "2026-10-27T15:10:00",
            "updated_at": null
        }
    ],
    "ike_policies_links": [
        {
            "rel": "previous",
            "href": "%s/v1/ike_policies?name=policy&limit=1&marker=5c561d9d-eaea-45f6-ae3e-08d1a7080828&page_reverse=True"
        }
    ]
}
`

var (
	Policy = ikepolicies.Policy{
		ID:                    PolicyID,
		ProjectID:             "9145d91459d248b1b02fdaca97c6a75d",
		Name:                  "policy",
		Description:           "IKE policy",
		AuthAlgorithm:         "sha384",
		EncryptionAlgorithm:   "aes-256-gcm-16",
		Phase1NegotiationMode: "main",
		IKEVersion:            "v2",
		PFS:                   "group19",
		Lifetime:              ikepolicies.Lifetime{Units: "seconds", Value: 7200},
		CreatedAt:             time.Date(2026, 10, 27, 15, 4, 5, 123456000, time.UTC),
	}

	UpdatedPolicy = ikepolicies.Policy{
		ID:                    PolicyID,
		ProjectID:             "9145d91459d248b1b02fdaca97c6a75d",
		Name:                  "updatedname",
		Description:           "updated policy",
		AuthAlgorithm:         "sha512",
		EncryptionAlgorithm:   "aes-128",
		Phase1NegotiationMode: "main",
		IKEVersion:            "v2",
		PFS:                   "group31",
		Lifetime:              ikepolicies.Lifetime{Units: "seconds", Value: 3600},
		CreatedAt:             time.Date(2026, 10, 27, 15, 4, 5, 0, time.UTC),
		UpdatedAt:             time.Date(2026, 10, 27, 16, 0, 0, 0, time.UTC),
	}

	SecondPolicy = ikepolicies.Policy{
		ID:                    "5c561d9d-eaea-45f6-ae3e-08d1a7080828",
		ProjectID:             "9145d91459d248b1b02fdaca97c6a75d",
		Name:                  "policy",
		AuthAlgorithm:         "sha256",
		EncryptionAlgorithm:   "aes-256",
		Phase1NegotiationMode: "main",
		IKEVersion:            "v2",
		PFS:                   "group14",
		Lifetime:              ikepolicies.Lifetime{Units: "seconds", Value: 10800},
		CreatedAt:             time.Date(2026, 10, 27, 15, 10, 0, 0, time.UTC),
	}
)

func HandlePolicyCreationSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/v1/ike_policies", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, CreateRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, PolicyBody)
	})
}

func HandlePolicyGetSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/v1/ike_policies/"+PolicyID, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, PolicyBody)
	})
}

func HandlePolicyUpdateSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/v1/ike_policies/"+PolicyID, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, UpdateRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, UpdatedPolicyBody)
	})
}

func HandlePolicyDeletionSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/v1/ike_policies/"+PolicyID, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.WriteHeader(http.StatusNoContent)
	})
}

// HandlePolicyListSuccessfully serves two pages: the first answers the
// ListOpts query, the second its next link.
func HandlePolicyListSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/v1/ike_policies", func(w http.ResponseWriter, r *http.Request) {
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
