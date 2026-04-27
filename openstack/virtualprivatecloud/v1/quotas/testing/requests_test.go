package testing

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	fake "github.com/gophercloud/gophercloud/v2/openstack/virtualprivatecloud/v1/common"
	"github.com/gophercloud/gophercloud/v2/openstack/virtualprivatecloud/v1/quotas"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestGetDefault(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v1/quotas/a1b2c3d4-e5f6-7890-abcd-ef1234567890", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, GetResponseDefault)
	})

	q, err := quotas.Get(context.TODO(), fake.ServiceClient(fakeServer), "a1b2c3d4-e5f6-7890-abcd-ef1234567890").Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "a1b2c3d4-e5f6-7890-abcd-ef1234567890", q.ProjectID)
	th.AssertEquals(t, 5, q.MaxVpcs)
	th.AssertEquals(t, "config_default", q.Source)
}

func TestGetOverride(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v1/quotas/a1b2c3d4-e5f6-7890-abcd-ef1234567890", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, GetResponseOverride)
	})

	q, err := quotas.Get(context.TODO(), fake.ServiceClient(fakeServer), "a1b2c3d4-e5f6-7890-abcd-ef1234567890").Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "a1b2c3d4-e5f6-7890-abcd-ef1234567890", q.ProjectID)
	th.AssertEquals(t, 10, q.MaxVpcs)
	th.AssertEquals(t, "project_override", q.Source)
}

func TestSet(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v1/quotas/a1b2c3d4-e5f6-7890-abcd-ef1234567890", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, SetRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, SetResponse)
	})

	opts := quotas.SetOpts{MaxVpcs: 10}
	q, err := quotas.Set(context.TODO(), fake.ServiceClient(fakeServer), "a1b2c3d4-e5f6-7890-abcd-ef1234567890", opts).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "a1b2c3d4-e5f6-7890-abcd-ef1234567890", q.ProjectID)
	th.AssertEquals(t, 10, q.MaxVpcs)
}

func TestDelete(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v1/quotas/a1b2c3d4-e5f6-7890-abcd-ef1234567890", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.WriteHeader(http.StatusNoContent)
	})

	err := quotas.Delete(context.TODO(), fake.ServiceClient(fakeServer), "a1b2c3d4-e5f6-7890-abcd-ef1234567890").ExtractErr()
	th.AssertNoErr(t, err)
}
