package testing

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	fake "github.com/gophercloud/gophercloud/v2/openstack/networking/v2/common"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/routetables"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestList(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/route_tables", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, `
{
    "route_tables": [
        {
            "id": "e3f5f925-6c6c-4b1e-9c1e-7c6d4d7c3f2a",
            "name": "default",
            "description": "Default route table",
            "router_id": "a5c3a4d4-5e0f-4c2a-bb02-e2723de5b8f1",
            "is_default": true,
            "routes": [
                {"destination": "10.0.0.0/8", "nexthop": "", "nexthop_type": "local"}
            ],
            "subnets": ["08eae331-0402-425a-923c-34f7cfe39c1b"],
            "project_id": "4fd44f30292945e481c7b8a0c8908869"
        },
        {
            "id": "b2f3c4d5-1a2b-3c4d-5e6f-7a8b9c0d1e2f",
            "name": "custom",
            "description": "Custom route table",
            "router_id": "a5c3a4d4-5e0f-4c2a-bb02-e2723de5b8f1",
            "is_default": false,
            "routes": [
                {"destination": "192.168.1.0/24", "nexthop": "10.0.0.1", "nexthop_type": "ip"}
            ],
            "subnets": [],
            "project_id": "4fd44f30292945e481c7b8a0c8908869"
        }
    ]
}`)
	})

	count := 0

	err := routetables.List(fake.ServiceClient(fakeServer), routetables.ListOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++
		actual, err := routetables.ExtractRouteTables(page)
		if err != nil {
			t.Errorf("Failed to extract route tables: %v", err)
			return false, err
		}

		expected := []routetables.RouteTable{
			{
				ID:          "e3f5f925-6c6c-4b1e-9c1e-7c6d4d7c3f2a",
				Name:        "default",
				Description: "Default route table",
				RouterID:    "a5c3a4d4-5e0f-4c2a-bb02-e2723de5b8f1",
				IsDefault:   true,
				Routes: []routetables.Route{
					{Destination: "10.0.0.0/8", Nexthop: "", NexthopType: "local"},
				},
				Subnets:   []string{"08eae331-0402-425a-923c-34f7cfe39c1b"},
				ProjectID: "4fd44f30292945e481c7b8a0c8908869",
			},
			{
				ID:          "b2f3c4d5-1a2b-3c4d-5e6f-7a8b9c0d1e2f",
				Name:        "custom",
				Description: "Custom route table",
				RouterID:    "a5c3a4d4-5e0f-4c2a-bb02-e2723de5b8f1",
				IsDefault:   false,
				Routes: []routetables.Route{
					{Destination: "192.168.1.0/24", Nexthop: "10.0.0.1", NexthopType: "ip"},
				},
				Subnets:   []string{},
				ProjectID: "4fd44f30292945e481c7b8a0c8908869",
			},
		}

		th.CheckDeepEquals(t, expected, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestListWithFilter(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/route_tables", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		th.AssertEquals(t, "a5c3a4d4-5e0f-4c2a-bb02-e2723de5b8f1", r.URL.Query().Get("router_id"))
		th.AssertEquals(t, "true", r.URL.Query().Get("is_default"))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, `
{
    "route_tables": [
        {
            "id": "e3f5f925-6c6c-4b1e-9c1e-7c6d4d7c3f2a",
            "name": "default",
            "description": "",
            "router_id": "a5c3a4d4-5e0f-4c2a-bb02-e2723de5b8f1",
            "is_default": true,
            "routes": [],
            "subnets": [],
            "project_id": "4fd44f30292945e481c7b8a0c8908869"
        }
    ]
}`)
	})

	isDefault := true
	listOpts := routetables.ListOpts{
		RouterID:  "a5c3a4d4-5e0f-4c2a-bb02-e2723de5b8f1",
		IsDefault: &isDefault,
	}

	allPages, err := routetables.List(fake.ServiceClient(fakeServer), listOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allRouteTables, err := routetables.ExtractRouteTables(allPages)
	th.AssertNoErr(t, err)

	th.AssertEquals(t, 1, len(allRouteTables))
	th.AssertEquals(t, true, allRouteTables[0].IsDefault)
}

func TestGet(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/route_tables/e3f5f925-6c6c-4b1e-9c1e-7c6d4d7c3f2a", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, `
{
    "route_table": {
        "id": "e3f5f925-6c6c-4b1e-9c1e-7c6d4d7c3f2a",
        "name": "default",
        "description": "Default route table",
        "router_id": "a5c3a4d4-5e0f-4c2a-bb02-e2723de5b8f1",
        "is_default": true,
        "routes": [
            {"destination": "10.0.0.0/8", "nexthop": "", "nexthop_type": "local"},
            {"destination": "0.0.0.0/0", "nexthop": "10.0.0.1", "nexthop_type": "ip"}
        ],
        "subnets": ["08eae331-0402-425a-923c-34f7cfe39c1b", "54d6f61d-db07-451c-9ab3-b9609b6b6f0b"],
        "project_id": "4fd44f30292945e481c7b8a0c8908869"
    }
}`)
	})

	rt, err := routetables.Get(context.TODO(), fake.ServiceClient(fakeServer), "e3f5f925-6c6c-4b1e-9c1e-7c6d4d7c3f2a").Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "e3f5f925-6c6c-4b1e-9c1e-7c6d4d7c3f2a", rt.ID)
	th.AssertEquals(t, "default", rt.Name)
	th.AssertEquals(t, "Default route table", rt.Description)
	th.AssertEquals(t, "a5c3a4d4-5e0f-4c2a-bb02-e2723de5b8f1", rt.RouterID)
	th.AssertEquals(t, true, rt.IsDefault)
	th.AssertEquals(t, 2, len(rt.Routes))
	th.AssertEquals(t, "10.0.0.0/8", rt.Routes[0].Destination)
	th.AssertEquals(t, "local", rt.Routes[0].NexthopType)
	th.AssertEquals(t, "0.0.0.0/0", rt.Routes[1].Destination)
	th.AssertEquals(t, "10.0.0.1", rt.Routes[1].Nexthop)
	th.AssertEquals(t, "ip", rt.Routes[1].NexthopType)
	th.AssertEquals(t, 2, len(rt.Subnets))
	th.AssertEquals(t, "4fd44f30292945e481c7b8a0c8908869", rt.ProjectID)
}

func TestCreate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/route_tables", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
    "route_table": {
        "name": "custom-rt",
        "description": "Custom route table for web tier",
        "router_id": "a5c3a4d4-5e0f-4c2a-bb02-e2723de5b8f1"
    }
}`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprint(w, `
{
    "route_table": {
        "id": "b2f3c4d5-1a2b-3c4d-5e6f-7a8b9c0d1e2f",
        "name": "custom-rt",
        "description": "Custom route table for web tier",
        "router_id": "a5c3a4d4-5e0f-4c2a-bb02-e2723de5b8f1",
        "is_default": false,
        "routes": [],
        "subnets": [],
        "project_id": "4fd44f30292945e481c7b8a0c8908869"
    }
}`)
	})

	options := routetables.CreateOpts{
		Name:        "custom-rt",
		Description: "Custom route table for web tier",
		RouterID:    "a5c3a4d4-5e0f-4c2a-bb02-e2723de5b8f1",
	}
	rt, err := routetables.Create(context.TODO(), fake.ServiceClient(fakeServer), options).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "b2f3c4d5-1a2b-3c4d-5e6f-7a8b9c0d1e2f", rt.ID)
	th.AssertEquals(t, "custom-rt", rt.Name)
	th.AssertEquals(t, "Custom route table for web tier", rt.Description)
	th.AssertEquals(t, "a5c3a4d4-5e0f-4c2a-bb02-e2723de5b8f1", rt.RouterID)
	th.AssertEquals(t, false, rt.IsDefault)
	th.AssertEquals(t, 0, len(rt.Routes))
	th.AssertEquals(t, 0, len(rt.Subnets))
}

func TestUpdate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/route_tables/b2f3c4d5-1a2b-3c4d-5e6f-7a8b9c0d1e2f", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
    "route_table": {
        "name": "renamed-rt",
        "description": "Updated description"
    }
}`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, `
{
    "route_table": {
        "id": "b2f3c4d5-1a2b-3c4d-5e6f-7a8b9c0d1e2f",
        "name": "renamed-rt",
        "description": "Updated description",
        "router_id": "a5c3a4d4-5e0f-4c2a-bb02-e2723de5b8f1",
        "is_default": false,
        "routes": [],
        "subnets": [],
        "project_id": "4fd44f30292945e481c7b8a0c8908869"
    }
}`)
	})

	name := "renamed-rt"
	description := "Updated description"
	options := routetables.UpdateOpts{
		Name:        &name,
		Description: &description,
	}
	rt, err := routetables.Update(context.TODO(), fake.ServiceClient(fakeServer), "b2f3c4d5-1a2b-3c4d-5e6f-7a8b9c0d1e2f", options).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "renamed-rt", rt.Name)
	th.AssertEquals(t, "Updated description", rt.Description)
}

func TestDelete(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/route_tables/b2f3c4d5-1a2b-3c4d-5e6f-7a8b9c0d1e2f", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})

	res := routetables.Delete(context.TODO(), fake.ServiceClient(fakeServer), "b2f3c4d5-1a2b-3c4d-5e6f-7a8b9c0d1e2f")
	th.AssertNoErr(t, res.Err)
}

func TestAddRoutes(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/route_tables/e3f5f925-6c6c-4b1e-9c1e-7c6d4d7c3f2a/add_routes", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestJSONRequest(t, r, `
{
    "routes": [
        {"destination": "192.168.1.0/24", "nexthop": "10.0.0.1", "nexthop_type": "ip"},
        {"destination": "172.16.0.0/12", "nexthop_type": "blackhole"}
    ]
}`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, `
{
    "route_table": {
        "id": "e3f5f925-6c6c-4b1e-9c1e-7c6d4d7c3f2a",
        "name": "default",
        "description": "",
        "router_id": "a5c3a4d4-5e0f-4c2a-bb02-e2723de5b8f1",
        "is_default": true,
        "routes": [
            {"destination": "10.0.0.0/8", "nexthop": "", "nexthop_type": "local"},
            {"destination": "192.168.1.0/24", "nexthop": "10.0.0.1", "nexthop_type": "ip"},
            {"destination": "172.16.0.0/12", "nexthop": "", "nexthop_type": "blackhole"}
        ],
        "subnets": [],
        "project_id": "4fd44f30292945e481c7b8a0c8908869"
    }
}`)
	})

	opts := routetables.AddRoutesOpts{
		Routes: []routetables.RouteOpts{
			{Destination: "192.168.1.0/24", Nexthop: "10.0.0.1", NexthopType: "ip"},
			{Destination: "172.16.0.0/12", NexthopType: "blackhole"},
		},
	}
	rt, err := routetables.AddRoutes(context.TODO(), fake.ServiceClient(fakeServer), "e3f5f925-6c6c-4b1e-9c1e-7c6d4d7c3f2a", opts).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, 3, len(rt.Routes))
	th.AssertEquals(t, "192.168.1.0/24", rt.Routes[1].Destination)
	th.AssertEquals(t, "10.0.0.1", rt.Routes[1].Nexthop)
	th.AssertEquals(t, "ip", rt.Routes[1].NexthopType)
	th.AssertEquals(t, "172.16.0.0/12", rt.Routes[2].Destination)
	th.AssertEquals(t, "blackhole", rt.Routes[2].NexthopType)
}

func TestRemoveRoutes(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/route_tables/e3f5f925-6c6c-4b1e-9c1e-7c6d4d7c3f2a/remove_routes", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestJSONRequest(t, r, `
{
    "routes": [
        {"destination": "192.168.1.0/24"}
    ]
}`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, `
{
    "route_table": {
        "id": "e3f5f925-6c6c-4b1e-9c1e-7c6d4d7c3f2a",
        "name": "default",
        "description": "",
        "router_id": "a5c3a4d4-5e0f-4c2a-bb02-e2723de5b8f1",
        "is_default": true,
        "routes": [
            {"destination": "10.0.0.0/8", "nexthop": "", "nexthop_type": "local"}
        ],
        "subnets": [],
        "project_id": "4fd44f30292945e481c7b8a0c8908869"
    }
}`)
	})

	opts := routetables.RemoveRoutesOpts{
		Routes: []routetables.RemoveRouteOpts{
			{Destination: "192.168.1.0/24"},
		},
	}
	rt, err := routetables.RemoveRoutes(context.TODO(), fake.ServiceClient(fakeServer), "e3f5f925-6c6c-4b1e-9c1e-7c6d4d7c3f2a", opts).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, 1, len(rt.Routes))
	th.AssertEquals(t, "10.0.0.0/8", rt.Routes[0].Destination)
}

func TestAddSubnets(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/route_tables/e3f5f925-6c6c-4b1e-9c1e-7c6d4d7c3f2a/add_subnets", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestJSONRequest(t, r, `
{
    "subnets": [
        "08eae331-0402-425a-923c-34f7cfe39c1b",
        "54d6f61d-db07-451c-9ab3-b9609b6b6f0b"
    ]
}`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, `
{
    "route_table": {
        "id": "e3f5f925-6c6c-4b1e-9c1e-7c6d4d7c3f2a",
        "name": "default",
        "description": "",
        "router_id": "a5c3a4d4-5e0f-4c2a-bb02-e2723de5b8f1",
        "is_default": true,
        "routes": [],
        "subnets": ["08eae331-0402-425a-923c-34f7cfe39c1b", "54d6f61d-db07-451c-9ab3-b9609b6b6f0b"],
        "project_id": "4fd44f30292945e481c7b8a0c8908869"
    }
}`)
	})

	opts := routetables.AddSubnetsOpts{
		Subnets: []string{
			"08eae331-0402-425a-923c-34f7cfe39c1b",
			"54d6f61d-db07-451c-9ab3-b9609b6b6f0b",
		},
	}
	rt, err := routetables.AddSubnets(context.TODO(), fake.ServiceClient(fakeServer), "e3f5f925-6c6c-4b1e-9c1e-7c6d4d7c3f2a", opts).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, 2, len(rt.Subnets))
	th.AssertEquals(t, "08eae331-0402-425a-923c-34f7cfe39c1b", rt.Subnets[0])
	th.AssertEquals(t, "54d6f61d-db07-451c-9ab3-b9609b6b6f0b", rt.Subnets[1])
}

func TestRemoveSubnets(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/route_tables/e3f5f925-6c6c-4b1e-9c1e-7c6d4d7c3f2a/remove_subnets", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestJSONRequest(t, r, `
{
    "subnets": [
        "54d6f61d-db07-451c-9ab3-b9609b6b6f0b"
    ]
}`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, `
{
    "route_table": {
        "id": "e3f5f925-6c6c-4b1e-9c1e-7c6d4d7c3f2a",
        "name": "default",
        "description": "",
        "router_id": "a5c3a4d4-5e0f-4c2a-bb02-e2723de5b8f1",
        "is_default": true,
        "routes": [],
        "subnets": ["08eae331-0402-425a-923c-34f7cfe39c1b"],
        "project_id": "4fd44f30292945e481c7b8a0c8908869"
    }
}`)
	})

	opts := routetables.RemoveSubnetsOpts{
		Subnets: []string{"54d6f61d-db07-451c-9ab3-b9609b6b6f0b"},
	}
	rt, err := routetables.RemoveSubnets(context.TODO(), fake.ServiceClient(fakeServer), "e3f5f925-6c6c-4b1e-9c1e-7c6d4d7c3f2a", opts).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, 1, len(rt.Subnets))
	th.AssertEquals(t, "08eae331-0402-425a-923c-34f7cfe39c1b", rt.Subnets[0])
}
