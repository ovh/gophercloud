//go:build acceptance || identity || endpointgroups
// +build acceptance identity endpointgroups

package v3

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/endpointgroups"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestEndpointGroupsList(t *testing.T) {
	clients.SkipRelease(t, "stable/mitaka")
	clients.SkipRelease(t, "stable/newton")
	clients.SkipRelease(t, "stable/ocata")
	clients.SkipRelease(t, "stable/pike")
	clients.SkipRelease(t, "stable/queens")
	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	allPages, err := endpointgroups.List(client, nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allEndpointGroups, err := endpointgroups.ExtractEndpointGroups(allPages)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, len(allEndpointGroups), 0)
}
