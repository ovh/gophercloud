package testing

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	fake "github.com/gophercloud/gophercloud/v2/openstack/networking/v2/common"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/vpcs"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestList(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/vpcs", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, VpcsListResult)
	})

	count := 0

	err := vpcs.List(fake.ServiceClient(fakeServer), vpcs.ListOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++
		actual, err := vpcs.ExtractVpcs(page)
		if err != nil {
			t.Errorf("Failed to extract vpcs: %v", err)
			return false, nil
		}

		expected := []vpcs.Vpc{
			Vpc1,
			Vpc2,
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

	fakeServer.Mux.HandleFunc("/v2.0/vpcs", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestFormValues(t, r, map[string]string{
			"id":         "6f1a5c3e-2b4d-4c8e-9f0a-1b2c3d4e5f60",
			"name":       "vpc-prod",
			"cidr_block": "10.0.0.0/16",
			"vpc_status": "ACTIVE",
			"tenant_id":  "4fd44f30292945e481c7b8a0c8908869",
			"project_id": "4fd44f30292945e481c7b8a0c8908869",
			"limit":      "10",
			"marker":     "a3b2c1d0-9e8f-4a7b-8c6d-5e4f3a2b1c0d",
			"sort_key":   "name",
			"sort_dir":   "asc",
		})

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, `{"vpcs": []}`)
	})

	listOpts := vpcs.ListOpts{
		ID:        "6f1a5c3e-2b4d-4c8e-9f0a-1b2c3d4e5f60",
		Name:      "vpc-prod",
		CIDRBlock: "10.0.0.0/16",
		VpcStatus: "ACTIVE",
		TenantID:  "4fd44f30292945e481c7b8a0c8908869",
		ProjectID: "4fd44f30292945e481c7b8a0c8908869",
		Limit:     10,
		Marker:    "a3b2c1d0-9e8f-4a7b-8c6d-5e4f3a2b1c0d",
		SortKey:   "name",
		SortDir:   "asc",
	}

	allPages, err := vpcs.List(fake.ServiceClient(fakeServer), listOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	actual, err := vpcs.ExtractVpcs(allPages)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 0, len(actual))
}

func TestGet(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/vpcs/6f1a5c3e-2b4d-4c8e-9f0a-1b2c3d4e5f60", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, VpcGetResult)
	})

	v, err := vpcs.Get(context.TODO(), fake.ServiceClient(fakeServer), "6f1a5c3e-2b4d-4c8e-9f0a-1b2c3d4e5f60").Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, Vpc1, *v)
}

func TestCreate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/vpcs", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, VpcCreateRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprint(w, VpcCreateResult)
	})

	opts := vpcs.CreateOpts{
		Name:        "vpc-prod",
		Description: "Production VPC",
		CIDRBlock:   "10.0.0.0/16",
		ProjectID:   "4fd44f30292945e481c7b8a0c8908869",
	}
	v, err := vpcs.Create(context.TODO(), fake.ServiceClient(fakeServer), opts).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "6f1a5c3e-2b4d-4c8e-9f0a-1b2c3d4e5f60", v.ID)
	th.AssertEquals(t, "vpc-prod", v.Name)
	th.AssertEquals(t, "Production VPC", v.Description)
	th.AssertEquals(t, "10.0.0.0/16", v.CIDRBlock)
	th.AssertEquals(t, "PENDING_CREATE", v.VpcStatus)
	th.AssertEquals(t, "4fd44f30292945e481c7b8a0c8908869", v.TenantID)
	th.AssertEquals(t, "4fd44f30292945e481c7b8a0c8908869", v.ProjectID)
}

func TestCreateRequiresCIDRBlock(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	res := vpcs.Create(context.TODO(), fake.ServiceClient(fakeServer), vpcs.CreateOpts{Name: "vpc-prod"})
	if res.Err == nil {
		t.Fatal("Expected an error when CIDRBlock is missing")
	}
}

func TestUpdate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/vpcs/6f1a5c3e-2b4d-4c8e-9f0a-1b2c3d4e5f60", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, VpcUpdateRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, VpcUpdateResult)
	})

	newName := "vpc-renamed"
	newDescription := ""
	opts := vpcs.UpdateOpts{
		Name:        &newName,
		Description: &newDescription,
	}
	v, err := vpcs.Update(context.TODO(), fake.ServiceClient(fakeServer), "6f1a5c3e-2b4d-4c8e-9f0a-1b2c3d4e5f60", opts).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "vpc-renamed", v.Name)
	th.AssertEquals(t, "", v.Description)
	th.AssertEquals(t, "PENDING_UPDATE", v.VpcStatus)
}

func TestDelete(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/vpcs/6f1a5c3e-2b4d-4c8e-9f0a-1b2c3d4e5f60", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.WriteHeader(http.StatusNoContent)
	})

	res := vpcs.Delete(context.TODO(), fake.ServiceClient(fakeServer), "6f1a5c3e-2b4d-4c8e-9f0a-1b2c3d4e5f60")
	th.AssertNoErr(t, res.Err)
}
