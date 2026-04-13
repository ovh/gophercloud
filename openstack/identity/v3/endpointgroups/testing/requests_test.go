package testing

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/endpointgroups"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestGetEndpointGroup(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/OS-EP-FILTER/endpoint_groups/24", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		fmt.Fprintf(w, `
		{
			"endpoint_group": {
			    "id": "24",
			    "filters": {
				    "interface": "public",
					"service_id": "1234",
					"region_id": "5678"
			    },
			    "name": "endpointgroup1",
			    "description": "public endpoint group 1",
			    "links": {
					"self": "https://localhost:5000/v3/OS-EP-FILTER/endpoint_groups/24"
			    }
			}
		}
		`)
	})

	actual, err := endpointgroups.Get(t.Context(), client.ServiceClient(fakeServer), "24").Extract()
	if err != nil {
		t.Fatalf("Failed to extract EndpointGroup: %v", err)
	}

	expected := &endpointgroups.EndpointGroup{
		ID: "24",
		Filters: map[string]interface{}{
			"interface":  "public",
			"service_id": "1234",
			"region_id":  "5678",
		},
		Name:        "endpointgroup1",
		Description: "public endpoint group 1",
	}
	th.AssertDeepEquals(t, expected, actual)
}

func TestListEndpointGroups(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/OS-EP-FILTER/endpoint_groups", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `
		{
			"endpoint_groups": [
				{
					"id": "24",
					"filters": {
						"interface": "public",
						"service_id": "1234",
						"region_id": "5678"
					},
					"name": "endpointgroup1",
					"description": "public endpoint group 1",
					"links": {
						"self": "https://localhost:5000/v3/OS-EP-FILTER/endpoint_groups/24"
					}
				},
				{
					"id": "25",
					"filters": {
						"interface": "internal"
					},
					"name": "endpointgroup2",
					"description": "internal endpoint group 1",
					"links": {
						"self": "https://localhost:5000/v3/OS-EP-FILTER/endpoint_groups/25"
					}
				}
			]
		}
		`)
	})

	err := endpointgroups.List(client.ServiceClient(fakeServer), endpointgroups.ListOpts{}).EachPage(t.Context(), func(ctx context.Context, page pagination.Page) (bool, error) {
		actual, err := endpointgroups.ExtractEndpointGroups(page)
		if err != nil {
			t.Errorf("Failed to extract EndpointGroups: %v", err)
			return false, err
		}

		expected := []endpointgroups.EndpointGroup{
			{
				ID: "24",
				Filters: map[string]interface{}{
					"interface":  "public",
					"service_id": "1234",
					"region_id":  "5678",
				},
				Name:        "endpointgroup1",
				Description: "public endpoint group 1",
			},
			{
				ID: "25",
				Filters: map[string]interface{}{
					"interface": "internal",
				},
				Name:        "endpointgroup2",
				Description: "internal endpoint group 1",
			},
		}
		th.AssertDeepEquals(t, expected, actual)
		return true, nil
	})
	th.AssertNoErr(t, err)
}

func TestListEndpointGroupsForProject(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/OS-EP-FILTER/projects/42/endpoint_groups", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `
		{
			"endpoint_groups": [
				{
				    "id": "24",
				    "filters": {
						"interface": "public",
						"service_id": "1234",
						"region_id": "5678"
				    },
				    "name": "endpointgroup1",
				    "description": "public endpoint group 1",
				    "links": {
						"self": "https://localhost:5000/v3/OS-EP-FILTER/endpoint_groups/24"
				    }
				},
				{
				    "id": "25",
				    "filters": {
						"interface": "internal"
				    },
				    "name": "endpointgroup2",
				    "description": "internal endpoint group 1",
				    "links": {
						"self": "https://localhost:5000/v3/OS-EP-FILTER/endpoint_groups/25"
				    }
				}
			]
		}
		`)
	})

	err := endpointgroups.ListForProjects(client.ServiceClient(fakeServer), "42").EachPage(t.Context(), func(ctx context.Context, page pagination.Page) (bool, error) {
		actual, err := endpointgroups.ExtractEndpointGroups(page)
		if err != nil {
			t.Errorf("Failed to extract EndpointGroups: %v", err)
			return false, err
		}

		expected := []endpointgroups.EndpointGroup{
			{
				ID: "24",
				Filters: map[string]interface{}{
					"interface":  "public",
					"service_id": "1234",
					"region_id":  "5678",
				},
				Name:        "endpointgroup1",
				Description: "public endpoint group 1",
			},
			{
				ID: "25",
				Filters: map[string]interface{}{
					"interface": "internal",
				},
				Name:        "endpointgroup2",
				Description: "internal endpoint group 1",
			},
		}
		th.AssertDeepEquals(t, expected, actual)
		return true, nil
	})
	th.AssertNoErr(t, err)
}

func TestCreateProjectAssociation(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/OS-EP-FILTER/endpoint_groups/24/projects/42", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.WriteHeader(http.StatusNoContent)
	})

	err := endpointgroups.CreateProjectAssociation(t.Context(), client.ServiceClient(fakeServer), "24", "42").ExtractErr()
	th.AssertNoErr(t, err)
}

func TestDeleteProjectAssociation(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/OS-EP-FILTER/endpoint_groups/24/projects/42", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.WriteHeader(http.StatusNoContent)
	})

	err := endpointgroups.DeleteProjectAssociation(t.Context(), client.ServiceClient(fakeServer), "24", "42").ExtractErr()
	th.AssertNoErr(t, err)
}

func TestCheckProjectAssociation(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/OS-EP-FILTER/endpoint_groups/24/projects/42", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "HEAD")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.WriteHeader(http.StatusOK)
	})

	err := endpointgroups.CheckProjectAssociation(t.Context(), client.ServiceClient(fakeServer), "24", "42").ExtractErr()
	th.AssertNoErr(t, err)
}
