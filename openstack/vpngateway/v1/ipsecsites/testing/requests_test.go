package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/vpngateway/v1/ipsecsites"
	fake "github.com/gophercloud/gophercloud/v2/openstack/vpngateway/v1/testhelper"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestCreate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleSiteCreationSuccessfully(t, fakeServer)

	ecmpInbound := true
	options := ipsecsites.CreateOpts{
		ProjectID:      "9145d91459d248b1b02fdaca97c6a75d",
		Name:           "site",
		Description:    "IPsec site",
		VPNGatewayID:   "2ad3c5dc-8c7e-4ff0-8b54-4b4d4e2c7c6e",
		Routing:        ipsecsites.RoutingBGP,
		LocalASN:       4200000001,
		PeerASN:        64512,
		ECMPInbound:    &ecmpInbound,
		LocalEPGroupID: "6ecd9cf3-ca64-46c7-863f-f2eb1b9e838a",
		PeerEPGroupID:  "3e1f6a2b-7c8d-4e9f-a0b1-c2d3e4f5a6b7",
	}

	actual, err := ipsecsites.Create(context.TODO(), fake.ServiceClient(fakeServer), options).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, Site, *actual)
}

func TestCreateRequiresGatewayAndGroups(t *testing.T) {
	complete := ipsecsites.CreateOpts{
		VPNGatewayID:   "2ad3c5dc-8c7e-4ff0-8b54-4b4d4e2c7c6e",
		LocalEPGroupID: "6ecd9cf3-ca64-46c7-863f-f2eb1b9e838a",
		PeerEPGroupID:  "3e1f6a2b-7c8d-4e9f-a0b1-c2d3e4f5a6b7",
	}
	_, err := complete.ToSiteCreateMap()
	th.AssertNoErr(t, err)

	for _, missing := range []func(*ipsecsites.CreateOpts){
		func(o *ipsecsites.CreateOpts) { o.VPNGatewayID = "" },
		func(o *ipsecsites.CreateOpts) { o.LocalEPGroupID = "" },
		func(o *ipsecsites.CreateOpts) { o.PeerEPGroupID = "" },
	} {
		opts := complete
		missing(&opts)
		if _, err := opts.ToSiteCreateMap(); err == nil {
			t.Errorf("expected an error for %+v", opts)
		}
	}
}

func TestGet(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleSiteGetSuccessfully(t, fakeServer)

	actual, err := ipsecsites.Get(context.TODO(), fake.ServiceClient(fakeServer), SiteID).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, Site, *actual)
}

func TestDelete(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleSiteDeletionSuccessfully(t, fakeServer)

	res := ipsecsites.Delete(context.TODO(), fake.ServiceClient(fakeServer), SiteID)
	th.AssertNoErr(t, res.Err)
}

func TestList(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleSiteListSuccessfully(t, fakeServer)

	options := ipsecsites.ListOpts{
		ProjectID:    "9145d91459d248b1b02fdaca97c6a75d",
		Name:         "site",
		VPNGatewayID: "2ad3c5dc-8c7e-4ff0-8b54-4b4d4e2c7c6e",
		Limit:        1,
		PageReverse:  true,
		SortKey:      "routing",
		SortDir:      "asc",
		Fields:       []string{"id", "name"},
	}

	var actual []ipsecsites.Site
	pages := 0
	err := ipsecsites.List(fake.ServiceClient(fakeServer), options).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		pages++
		sites, err := ipsecsites.ExtractSites(page)
		if err != nil {
			return false, err
		}
		actual = append(actual, sites...)
		return true, nil
	})
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 2, pages)
	th.CheckDeepEquals(t, ListedSites, actual)
}

func TestListOptsMarker(t *testing.T) {
	query, err := ipsecsites.ListOpts{Marker: SiteID}.ToSiteListQuery()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "?marker="+SiteID, query)
}

func TestUpdate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleSiteUpdateSuccessfully(t, fakeServer)

	name := "updatedname"
	description := "updated site"
	options := ipsecsites.UpdateOpts{
		Name:        &name,
		Description: &description,
	}

	actual, err := ipsecsites.Update(context.TODO(), fake.ServiceClient(fakeServer), SiteID, options).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, UpdatedSite, *actual)
}
