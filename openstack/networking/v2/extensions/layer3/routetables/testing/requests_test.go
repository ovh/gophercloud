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

const routeTableID = "b2f3c4d5-1a2b-3c4d-5e6f-7a8b9c0d1e2f"

func TestList(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/route_tables", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, RouteTablesListResult)
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
			RouteTable1,
			RouteTable2,
		}

		th.CheckDeepEquals(t, expected, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestListWithFilters(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/route_tables", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestFormValues(t, r, map[string]string{
			"id":         "e3f5f925-6c6c-4b1e-9c1e-7c6d4d7c3f2a",
			"name":       "public",
			"router_id":  "a5c3a4d4-5e0f-4c2a-bb02-e2723de5b8f1",
			"vpc_status": "ACTIVE",
			"tenant_id":  "4fd44f30292945e481c7b8a0c8908869",
			"project_id": "4fd44f30292945e481c7b8a0c8908869",
			"limit":      "10",
			"marker":     "b2f3c4d5-1a2b-3c4d-5e6f-7a8b9c0d1e2f",
			"sort_key":   "name",
			"sort_dir":   "asc",
		})

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, `{"route_tables": []}`)
	})

	listOpts := routetables.ListOpts{
		ID:        "e3f5f925-6c6c-4b1e-9c1e-7c6d4d7c3f2a",
		Name:      "public",
		RouterID:  "a5c3a4d4-5e0f-4c2a-bb02-e2723de5b8f1",
		VpcStatus: "ACTIVE",
		TenantID:  "4fd44f30292945e481c7b8a0c8908869",
		ProjectID: "4fd44f30292945e481c7b8a0c8908869",
		Limit:     10,
		Marker:    "b2f3c4d5-1a2b-3c4d-5e6f-7a8b9c0d1e2f",
		SortKey:   "name",
		SortDir:   "asc",
	}

	allPages, err := routetables.List(fake.ServiceClient(fakeServer), listOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	actual, err := routetables.ExtractRouteTables(allPages)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 0, len(actual))
}

func TestGet(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/route_tables/e3f5f925-6c6c-4b1e-9c1e-7c6d4d7c3f2a", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, RouteTableGetResult)
	})

	rt, err := routetables.Get(context.TODO(), fake.ServiceClient(fakeServer), "e3f5f925-6c6c-4b1e-9c1e-7c6d4d7c3f2a").Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, RouteTable1, *rt)
}

func TestCreate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/route_tables", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, RouteTableCreateRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprint(w, RouteTableCreateResult)
	})

	opts := routetables.CreateOpts{
		Name:        "private",
		Description: "Private route table",
		RouterID:    "a5c3a4d4-5e0f-4c2a-bb02-e2723de5b8f1",
		ProjectID:   "4fd44f30292945e481c7b8a0c8908869",
	}
	rt, err := routetables.Create(context.TODO(), fake.ServiceClient(fakeServer), opts).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, routeTableID, rt.ID)
	th.AssertEquals(t, "private", rt.Name)
	th.AssertEquals(t, "Private route table", rt.Description)
	th.AssertEquals(t, "a5c3a4d4-5e0f-4c2a-bb02-e2723de5b8f1", rt.RouterID)
	th.AssertEquals(t, 0, len(rt.Routes))
	th.AssertEquals(t, 0, len(rt.Subnets))
	th.AssertEquals(t, "PENDING_CREATE", rt.VpcStatus)
	th.AssertEquals(t, "4fd44f30292945e481c7b8a0c8908869", rt.TenantID)
	th.AssertEquals(t, "4fd44f30292945e481c7b8a0c8908869", rt.ProjectID)
}

func TestCreateRequiresRouterID(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	res := routetables.Create(context.TODO(), fake.ServiceClient(fakeServer), routetables.CreateOpts{Name: "private"})
	if res.Err == nil {
		t.Fatal("Expected an error when RouterID is missing")
	}
}

func TestUpdate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/route_tables/"+routeTableID, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, RouteTableUpdateRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, RouteTableUpdateResult)
	})

	name := "renamed"
	description := ""
	opts := routetables.UpdateOpts{
		Name:        &name,
		Description: &description,
	}
	rt, err := routetables.Update(context.TODO(), fake.ServiceClient(fakeServer), routeTableID, opts).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "renamed", rt.Name)
	th.AssertEquals(t, "", rt.Description)
	th.AssertEquals(t, "PENDING_UPDATE", rt.VpcStatus)
}

func TestDelete(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/route_tables/"+routeTableID, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.WriteHeader(http.StatusNoContent)
	})

	res := routetables.Delete(context.TODO(), fake.ServiceClient(fakeServer), routeTableID)
	th.AssertNoErr(t, res.Err)
}

func TestAddRoutes(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/route_tables/"+routeTableID+"/add_routes", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, AddRoutesRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, AddRoutesResult)
	})

	opts := routetables.AddRoutesOpts{
		Routes: []routetables.RouteOpts{
			{Destination: "192.168.1.0/24", Nexthop: "10.0.0.1", NexthopType: "ip"},
			{Destination: "192.168.2.0/24", NexthopType: "blackhole"},
		},
	}
	rt, err := routetables.AddRoutes(context.TODO(), fake.ServiceClient(fakeServer), routeTableID, opts).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, routeTableID, rt.ID)
	th.CheckDeepEquals(t, []routetables.Route{
		{Destination: "192.168.1.0/24", Nexthop: "10.0.0.1", NexthopType: "ip"},
		{Destination: "192.168.2.0/24", Nexthop: "", NexthopType: "blackhole"},
	}, rt.Routes)
	th.AssertEquals(t, "PENDING_UPDATE", rt.VpcStatus)
}

func TestAddRoutesRequiredFields(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	res := routetables.AddRoutes(context.TODO(), fake.ServiceClient(fakeServer), routeTableID, routetables.AddRoutesOpts{
		Routes: []routetables.RouteOpts{{Destination: "192.168.1.0/24"}},
	})
	if res.Err == nil {
		t.Fatal("Expected an error when NexthopType is missing")
	}
}

func TestRemoveRoutes(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/route_tables/"+routeTableID+"/remove_routes", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, RemoveRoutesRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, RemoveRoutesResult)
	})

	opts := routetables.RemoveRoutesOpts{
		Routes: []routetables.RemoveRouteOpts{
			{Destination: "192.168.1.0/24"},
		},
	}
	rt, err := routetables.RemoveRoutes(context.TODO(), fake.ServiceClient(fakeServer), routeTableID, opts).Extract()
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, []routetables.Route{
		{Destination: "192.168.2.0/24", Nexthop: "", NexthopType: "blackhole"},
	}, rt.Routes)
}

func TestAddSubnets(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/route_tables/"+routeTableID+"/add_subnets", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, AddSubnetsRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, AddSubnetsResult)
	})

	opts := routetables.AddSubnetsOpts{
		Subnets: []string{"08eae331-0402-425a-923c-34f7cfe39c1b"},
		Routes: []routetables.RouteOpts{
			{Destination: "10.0.1.0/24", NexthopType: "local"},
		},
	}
	rt, err := routetables.AddSubnets(context.TODO(), fake.ServiceClient(fakeServer), routeTableID, opts).Extract()
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, []string{"08eae331-0402-425a-923c-34f7cfe39c1b"}, rt.Subnets)
	th.CheckDeepEquals(t, []routetables.Route{
		{Destination: "10.0.1.0/24", Nexthop: "", NexthopType: "local"},
	}, rt.Routes)
}

func TestAddSubnetsWithoutRoutes(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/route_tables/"+routeTableID+"/add_subnets", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestJSONRequest(t, r, AddSubnetsWithoutRoutesRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, AddSubnetsResult)
	})

	opts := routetables.AddSubnetsOpts{
		Subnets: []string{"08eae331-0402-425a-923c-34f7cfe39c1b"},
	}
	_, err := routetables.AddSubnets(context.TODO(), fake.ServiceClient(fakeServer), routeTableID, opts).Extract()
	th.AssertNoErr(t, err)
}

func TestRemoveSubnets(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/route_tables/"+routeTableID+"/remove_subnets", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, RemoveSubnetsRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, RemoveSubnetsResult)
	})

	opts := routetables.RemoveSubnetsOpts{
		Subnets: []string{"08eae331-0402-425a-923c-34f7cfe39c1b"},
		Routes: []routetables.RemoveRouteOpts{
			{Destination: "10.0.1.0/24"},
		},
	}
	rt, err := routetables.RemoveSubnets(context.TODO(), fake.ServiceClient(fakeServer), routeTableID, opts).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, 0, len(rt.Subnets))
	th.AssertEquals(t, 0, len(rt.Routes))
}
