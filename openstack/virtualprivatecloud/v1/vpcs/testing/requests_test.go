package testing

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	fake "github.com/gophercloud/gophercloud/v2/openstack/virtualprivatecloud/v1/common"
	"github.com/gophercloud/gophercloud/v2/openstack/virtualprivatecloud/v1/vpcs"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestList(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v1/vpcs", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, ListResponse)
	})

	client := fake.ServiceClient(fakeServer)
	count := 0

	err := vpcs.List(client, vpcs.ListOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++
		actual, err := vpcs.ExtractVpcs(page)
		if err != nil {
			t.Errorf("Failed to extract vpcs: %v", err)
			return false, err
		}

		th.CheckDeepEquals(t, ExpectedVpcSlice, actual)

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

	fakeServer.Mux.HandleFunc("/v1/vpcs", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		th.AssertEquals(t, "ACTIVE", r.URL.Query().Get("status"))
		th.AssertEquals(t, "my-vpc", r.URL.Query().Get("name"))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, ListResponse)
	})

	client := fake.ServiceClient(fakeServer)

	allPages, err := vpcs.List(client, vpcs.ListOpts{
		Name:   "my-vpc",
		Status: "ACTIVE",
	}).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allVpcs, err := vpcs.ExtractVpcs(allPages)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 2, len(allVpcs))
}

func TestGet(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v1/vpcs/f47ac10b-58cc-4372-a567-0e02b2c3d479", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, GetResponse)
	})

	n, err := vpcs.Get(context.TODO(), fake.ServiceClient(fakeServer), "f47ac10b-58cc-4372-a567-0e02b2c3d479").Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &Vpc1, n)
	th.AssertEquals(t, n.CreatedAt.Format(time.RFC3339), "2026-03-04T12:00:00Z")
	th.AssertEquals(t, n.UpdatedAt.Format(time.RFC3339), "2026-03-04T12:05:00Z")
}

func TestCreate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v1/vpcs", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, CreateRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)

		fmt.Fprint(w, CreateResponse)
	})

	options := vpcs.CreateOpts{
		Name:        "my-vpc",
		CIDRBlock:   "10.0.0.0/16",
		Description: "Production VPC",
	}
	n, err := vpcs.Create(context.TODO(), fake.ServiceClient(fakeServer), options).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "f47ac10b-58cc-4372-a567-0e02b2c3d479", n.ID)
	th.AssertEquals(t, "CREATING", n.Status)
	th.AssertEquals(t, "my-vpc", n.Name)
	th.AssertEquals(t, "10.0.0.0/16", n.CIDRBlock)
	th.AssertEquals(t, n.CreatedAt.Format(time.RFC3339), "2026-03-04T12:00:00Z")
}

func TestUpdate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v1/vpcs/f47ac10b-58cc-4372-a567-0e02b2c3d479", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, UpdateRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, UpdateResponse)
	})

	name := "my-vpc-renamed"
	options := vpcs.UpdateOpts{Name: &name}
	n, err := vpcs.Update(context.TODO(), fake.ServiceClient(fakeServer), "f47ac10b-58cc-4372-a567-0e02b2c3d479", options).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "my-vpc-renamed", n.Name)
	th.AssertEquals(t, "ACTIVE", n.Status)
	th.AssertEquals(t, n.UpdatedAt.Format(time.RFC3339), "2026-03-04T12:10:00Z")
}

func TestDelete(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v1/vpcs/f47ac10b-58cc-4372-a567-0e02b2c3d479", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)

		fmt.Fprint(w, DeleteResponse)
	})

	n, err := vpcs.Delete(context.TODO(), fake.ServiceClient(fakeServer), "f47ac10b-58cc-4372-a567-0e02b2c3d479").Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "DELETING", n.Status)
	th.AssertEquals(t, "f47ac10b-58cc-4372-a567-0e02b2c3d479", n.ID)
	th.AssertEquals(t, n.UpdatedAt.Format(time.RFC3339), "2026-03-04T12:15:00Z")
}
