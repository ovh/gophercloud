package testing

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	fake "github.com/gophercloud/gophercloud/v2/openstack/networking/v2/common"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/snatips"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestList(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/snat_ips", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, SnatIPsListResult)
	})

	count := 0

	err := snatips.List(fake.ServiceClient(fakeServer), snatips.ListOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++
		actual, err := snatips.ExtractSnatIPs(page)
		if err != nil {
			t.Errorf("Failed to extract snat ips: %v", err)
			return false, nil
		}

		expected := []snatips.SnatIP{
			SnatIP1,
			SnatIP2,
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

	fakeServer.Mux.HandleFunc("/v2.0/snat_ips", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestFormValues(t, r, map[string]string{
			"id":                "2c4e6a8b-1d3f-4a5c-8e7b-9d0f1a2b3c4d",
			"network_id":        "8f9a1c2b-3d4e-4f5a-9b6c-7d8e9f0a1b2c",
			"availability_zone": "eu-west-par-a",
			"ip_address":        "203.0.113.10",
			"vpc_id":            "6f1a5c3e-2b4d-4c8e-9f0a-1b2c3d4e5f60",
			"status":            "ACTIVE",
			"vpc_status":        "ACTIVE",
			"tenant_id":         "4fd44f30292945e481c7b8a0c8908869",
			"project_id":        "4fd44f30292945e481c7b8a0c8908869",
			"limit":             "10",
			"marker":            "7e9f1a3b-5c6d-4e8f-a0b1-c2d3e4f5a6b7",
			"sort_key":          "ip_address",
			"sort_dir":          "desc",
		})

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, `{"snat_ips": []}`)
	})

	listOpts := snatips.ListOpts{
		ID:               "2c4e6a8b-1d3f-4a5c-8e7b-9d0f1a2b3c4d",
		NetworkID:        "8f9a1c2b-3d4e-4f5a-9b6c-7d8e9f0a1b2c",
		AvailabilityZone: "eu-west-par-a",
		IPAddress:        "203.0.113.10",
		VpcID:            "6f1a5c3e-2b4d-4c8e-9f0a-1b2c3d4e5f60",
		Status:           "ACTIVE",
		VpcStatus:        "ACTIVE",
		TenantID:         "4fd44f30292945e481c7b8a0c8908869",
		ProjectID:        "4fd44f30292945e481c7b8a0c8908869",
		Limit:            10,
		Marker:           "7e9f1a3b-5c6d-4e8f-a0b1-c2d3e4f5a6b7",
		SortKey:          "ip_address",
		SortDir:          "desc",
	}

	allPages, err := snatips.List(fake.ServiceClient(fakeServer), listOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	actual, err := snatips.ExtractSnatIPs(allPages)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 0, len(actual))
}

func TestGet(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/snat_ips/2c4e6a8b-1d3f-4a5c-8e7b-9d0f1a2b3c4d", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, SnatIPGetResult)
	})

	s, err := snatips.Get(context.TODO(), fake.ServiceClient(fakeServer), "2c4e6a8b-1d3f-4a5c-8e7b-9d0f1a2b3c4d").Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, SnatIP1, *s)
}

func TestCreate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/snat_ips", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, SnatIPCreateRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprint(w, SnatIPCreateResult)
	})

	opts := snatips.CreateOpts{
		NetworkID:        "8f9a1c2b-3d4e-4f5a-9b6c-7d8e9f0a1b2c",
		AvailabilityZone: "eu-west-par-a",
		IPAddress:        "203.0.113.10",
		Description:      "NAT for production",
		ProjectID:        "4fd44f30292945e481c7b8a0c8908869",
	}
	s, err := snatips.Create(context.TODO(), fake.ServiceClient(fakeServer), opts).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "2c4e6a8b-1d3f-4a5c-8e7b-9d0f1a2b3c4d", s.ID)
	th.AssertEquals(t, "8f9a1c2b-3d4e-4f5a-9b6c-7d8e9f0a1b2c", s.NetworkID)
	th.AssertEquals(t, "eu-west-par-a", s.AvailabilityZone)
	th.AssertEquals(t, "203.0.113.10", s.IPAddress)
	th.AssertEquals(t, "NAT for production", s.Description)
	th.AssertEquals(t, "", s.VpcID)
	th.AssertEquals(t, "DOWN", s.Status)
	th.AssertEquals(t, "PENDING_CREATE", s.VpcStatus)
	th.AssertEquals(t, "4fd44f30292945e481c7b8a0c8908869", s.TenantID)
	th.AssertEquals(t, "4fd44f30292945e481c7b8a0c8908869", s.ProjectID)
}

func TestCreateRequiredFields(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	res := snatips.Create(context.TODO(), fake.ServiceClient(fakeServer), snatips.CreateOpts{
		AvailabilityZone: "eu-west-par-a",
	})
	if res.Err == nil {
		t.Fatal("Expected an error when NetworkID is missing")
	}

	res = snatips.Create(context.TODO(), fake.ServiceClient(fakeServer), snatips.CreateOpts{
		NetworkID: "8f9a1c2b-3d4e-4f5a-9b6c-7d8e9f0a1b2c",
	})
	if res.Err == nil {
		t.Fatal("Expected an error when AvailabilityZone is missing")
	}
}

func TestUpdate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/snat_ips/2c4e6a8b-1d3f-4a5c-8e7b-9d0f1a2b3c4d", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, SnatIPUpdateRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, SnatIPUpdateResult)
	})

	description := "NAT for staging"
	opts := snatips.UpdateOpts{
		Description: &description,
	}
	s, err := snatips.Update(context.TODO(), fake.ServiceClient(fakeServer), "2c4e6a8b-1d3f-4a5c-8e7b-9d0f1a2b3c4d", opts).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "NAT for staging", s.Description)
}

func TestDelete(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/snat_ips/2c4e6a8b-1d3f-4a5c-8e7b-9d0f1a2b3c4d", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.WriteHeader(http.StatusNoContent)
	})

	res := snatips.Delete(context.TODO(), fake.ServiceClient(fakeServer), "2c4e6a8b-1d3f-4a5c-8e7b-9d0f1a2b3c4d")
	th.AssertNoErr(t, res.Err)
}
