package testing

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/keymanager/v1/quotas"
	fake "github.com/gophercloud/gophercloud/v2/openstack/networking/v2/common"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestGet(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/quotas", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, GetResponseRaw)
	})

	q, err := quotas.Get(t.Context(), client.ServiceClient(fakeServer)).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, GetResponse, q)
}

func TestListQuotas(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/project-quotas", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, ListResponseRaw)
	})

	count := 0
	err := quotas.List(client.ServiceClient(fakeServer), nil).EachPage(t.Context(), func(ctx context.Context, page pagination.Page) (bool, error) {
		count++

		actual, err := quotas.ExtractQuotas(page)
		th.AssertNoErr(t, err)

		th.AssertDeepEquals(t, ExpectedQuotasSlice, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.AssertEquals(t, count, 1)
}

func TestListOrdersAllPages(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/project-quotas", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, ListResponseRaw)
	})

	allPages, err := quotas.List(client.ServiceClient(fakeServer), nil).AllPages(t.Context())
	th.AssertNoErr(t, err)
	actual, err := quotas.ExtractQuotas(allPages)
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedQuotasSlice, actual)
}

func TestGetProjectQuota(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/project-quotas/0a73845280574ad389c292f6a74afa76", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, GetProjectResponseRaw)
	})

	q, err := quotas.GetProjectQuota(t.Context(), client.ServiceClient(fakeServer), "0a73845280574ad389c292f6a74afa76").Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, GetResponse, q)
}

func TestUpdate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/project-quotas/0a73845280574ad389c292f6a74afa76", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestJSONRequest(t, r, `
		{
			"project_quotas": {
				"secrets": 10,
				"containers": 14,
				"consumers": 15
			}
		}`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	})

	err := quotas.Update(t.Context(), client.ServiceClient(fakeServer), "0a73845280574ad389c292f6a74afa76", quotas.UpdateOpts{
		Secrets:    gophercloud.IntToPointer(10),
		Orders:     nil,
		Containers: gophercloud.IntToPointer(14),
		Consumers:  gophercloud.IntToPointer(15),
		CAS:        nil,
	}).Err

	th.AssertNoErr(t, err)
}

func TestDelete(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/project-quotas/0a73845280574ad389c292f6a74afa76", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	})

	err := quotas.Delete(t.Context(), client.ServiceClient(fakeServer), "0a73845280574ad389c292f6a74afa76").Err

	th.AssertNoErr(t, err)
}
