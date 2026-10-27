package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/vpngateway/v1/siteconnections"
	fake "github.com/gophercloud/gophercloud/v2/openstack/vpngateway/v1/testhelper"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestCreate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleConnectionCreationSuccessfully(t, fakeServer)

	adminStateUp := true
	options := siteconnections.CreateOpts{
		ProjectID:     "9145d91459d248b1b02fdaca97c6a75d",
		Name:          "tunnel",
		Description:   "IPsec site connection",
		IPsecSiteID:   "9d1e6c2a-4b3f-4a5e-8d7c-6b5a4f3e2d1c",
		PeerAddress:   "198.51.100.20",
		PeerID:        "vpn.peer.example.com",
		LocalID:       "@gateway.example.com",
		AuthMode:      siteconnections.AuthModePSK,
		PSK:           "a-long-pre-shared-key",
		Priority:      1,
		AdminStateUp:  &adminStateUp,
		IKEPolicyID:   "f2b08c1e-aa81-4668-8ae1-1401bcb0576c",
		IPsecPolicyID: "e3b0c442-98fc-4c14-9afb-f4c8996fb924",
		DPD: &siteconnections.DPDOpts{
			Action:   siteconnections.DPDActionRestart,
			Interval: 10,
			Timeout:  60,
		},
		MTU:            1380,
		Initiator:      siteconnections.InitiatorResponseOnly,
		TunnelAddress:  "169.254.10.1/30",
		BGPPeerAddress: "169.254.10.2",
	}

	actual, err := siteconnections.Create(context.TODO(), fake.ServiceClient(fakeServer), options).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, Connection, *actual)
}

func TestCreateRequiresFields(t *testing.T) {
	complete := siteconnections.CreateOpts{
		IPsecSiteID:   "9d1e6c2a-4b3f-4a5e-8d7c-6b5a4f3e2d1c",
		PeerAddress:   "198.51.100.20",
		PSK:           "a-long-pre-shared-key",
		IKEPolicyID:   "f2b08c1e-aa81-4668-8ae1-1401bcb0576c",
		IPsecPolicyID: "e3b0c442-98fc-4c14-9afb-f4c8996fb924",
	}
	_, err := complete.ToConnectionCreateMap()
	th.AssertNoErr(t, err)

	for _, missing := range []func(*siteconnections.CreateOpts){
		func(o *siteconnections.CreateOpts) { o.IPsecSiteID = "" },
		func(o *siteconnections.CreateOpts) { o.PeerAddress = "" },
		func(o *siteconnections.CreateOpts) { o.PSK = "" },
		func(o *siteconnections.CreateOpts) { o.IKEPolicyID = "" },
		func(o *siteconnections.CreateOpts) { o.IPsecPolicyID = "" },
	} {
		opts := complete
		missing(&opts)
		if _, err := opts.ToConnectionCreateMap(); err == nil {
			t.Errorf("expected an error for %+v", opts)
		}
	}
}

func TestGet(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleConnectionGetSuccessfully(t, fakeServer)

	actual, err := siteconnections.Get(context.TODO(), fake.ServiceClient(fakeServer), ConnectionID).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, Connection, *actual)
}

func TestDelete(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleConnectionDeletionSuccessfully(t, fakeServer)

	res := siteconnections.Delete(context.TODO(), fake.ServiceClient(fakeServer), ConnectionID)
	th.AssertNoErr(t, res.Err)
}

func TestList(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleConnectionListSuccessfully(t, fakeServer)

	options := siteconnections.ListOpts{
		ProjectID:   "9145d91459d248b1b02fdaca97c6a75d",
		Name:        "tunnel",
		Status:      siteconnections.StatusActive,
		IPsecSiteID: "9d1e6c2a-4b3f-4a5e-8d7c-6b5a4f3e2d1c",
		Limit:       1,
		PageReverse: true,
		SortKey:     "priority",
		SortDir:     "asc",
		Fields:      []string{"id", "name", "status"},
	}

	var actual []siteconnections.Connection
	pages := 0
	err := siteconnections.List(fake.ServiceClient(fakeServer), options).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		pages++
		connections, err := siteconnections.ExtractConnections(page)
		if err != nil {
			return false, err
		}
		actual = append(actual, connections...)
		return true, nil
	})
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 2, pages)
	th.CheckDeepEquals(t, ListedConnections, actual)
}

func TestListOptsMarker(t *testing.T) {
	query, err := siteconnections.ListOpts{Marker: ConnectionID}.ToConnectionListQuery()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "?marker="+ConnectionID, query)
}

func TestUpdate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleConnectionUpdateSuccessfully(t, fakeServer)

	name := "updatedname"
	description := "updated connection"
	adminStateUp := false
	options := siteconnections.UpdateOpts{
		Name:          &name,
		Description:   &description,
		PeerAddress:   "198.51.100.30",
		PeerID:        "198.51.100.30",
		LocalID:       "203.0.113.10",
		PSK:           "another-pre-shared-key",
		Priority:      2,
		AdminStateUp:  &adminStateUp,
		IKEPolicyID:   "0f1e2d3c-4b5a-4968-8778-695a4b3c2d1e",
		IPsecPolicyID: "1e2d3c4b-5a69-4877-9685-a4b3c2d1e0f1",
		DPD: &siteconnections.DPDOpts{
			Action: siteconnections.DPDActionHold,
		},
		MTU:            1400,
		Initiator:      siteconnections.InitiatorBiDirectional,
		TunnelAddress:  "169.254.20.1/30",
		BGPPeerAddress: "169.254.20.2",
	}

	actual, err := siteconnections.Update(context.TODO(), fake.ServiceClient(fakeServer), ConnectionID, options).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, UpdatedConnection, *actual)
}
