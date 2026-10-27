//go:build acceptance || vpngateway || ipsecsites

package v1

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/vpngateway/v1/ipsecsites"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestIPsecSitesList(t *testing.T) {
	client := NewClient(t)

	allPages, err := ipsecsites.List(client, nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allSites, err := ipsecsites.ExtractSites(allPages)
	th.AssertNoErr(t, err)

	for _, site := range allSites {
		tools.PrintResource(t, site)
	}
}

func TestIPsecSitesCRUD(t *testing.T) {
	client := NewClient(t)

	localGroup, err := CreateEndpointGroup(t, client, []string{"10.1.0.0/24"})
	th.AssertNoErr(t, err)
	defer DeleteEndpointGroup(t, client, localGroup.ID)

	peerGroup, err := CreateEndpointGroup(t, client, []string{"10.2.0.0/24"})
	th.AssertNoErr(t, err)
	defer DeleteEndpointGroup(t, client, peerGroup.ID)

	gateway, err := CreateVPNGateway(t, client)
	th.AssertNoErr(t, err)
	defer DeleteVPNGateway(t, client, gateway.ID)

	site, err := CreateIPsecSite(t, client, ipsecsites.CreateOpts{
		Description:    "static site",
		VPNGatewayID:   gateway.ID,
		LocalEPGroupID: localGroup.ID,
		PeerEPGroupID:  peerGroup.ID,
	})
	th.AssertNoErr(t, err)

	tools.PrintResource(t, site)

	th.AssertEquals(t, "static", site.Routing)
	if site.LocalASN != nil || site.PeerASN != nil {
		t.Fatal("Expected no ASNs on a static site")
	}

	newSite, err := ipsecsites.Get(context.TODO(), client, site.ID).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, *site, *newSite)

	allPages, err := ipsecsites.List(client, ipsecsites.ListOpts{VPNGatewayID: gateway.ID}).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	allSites, err := ipsecsites.ExtractSites(allPages)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 1, len(allSites))
	th.AssertEquals(t, site.ID, allSites[0].ID)

	name := tools.RandomString("TESTACC-", 8)
	description := tools.RandomString("TESTACC-desc-", 8)
	updatedSite, err := ipsecsites.Update(context.TODO(), client, site.ID, ipsecsites.UpdateOpts{
		Name:        &name,
		Description: &description,
	}).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, updatedSite)

	th.AssertEquals(t, name, updatedSite.Name)
	th.AssertEquals(t, description, updatedSite.Description)

	// A gateway holds a single site: replace the static site with a bgp one.
	DeleteIPsecSite(t, client, site.ID)

	ecmpInbound := true
	site, err = CreateIPsecSite(t, client, ipsecsites.CreateOpts{
		VPNGatewayID:   gateway.ID,
		Routing:        ipsecsites.RoutingBGP,
		LocalASN:       4200000001,
		PeerASN:        64512,
		ECMPInbound:    &ecmpInbound,
		LocalEPGroupID: localGroup.ID,
		PeerEPGroupID:  peerGroup.ID,
	})
	th.AssertNoErr(t, err)
	defer DeleteIPsecSite(t, client, site.ID)

	tools.PrintResource(t, site)

	th.AssertEquals(t, "bgp", site.Routing)
	th.AssertEquals(t, int64(4200000001), *site.LocalASN)
	th.AssertEquals(t, int64(64512), *site.PeerASN)
	th.AssertEquals(t, true, site.ECMPInbound)
}
