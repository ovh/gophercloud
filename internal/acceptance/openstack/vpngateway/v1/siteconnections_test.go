//go:build acceptance || vpngateway || siteconnections

package v1

import (
	"context"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/vpngateway/v1/ipsecsites"
	"github.com/gophercloud/gophercloud/v2/openstack/vpngateway/v1/siteconnections"
	"github.com/gophercloud/gophercloud/v2/openstack/vpngateway/v1/vpngateways"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestSiteConnectionsList(t *testing.T) {
	client := NewClient(t)

	allPages, err := siteconnections.List(client, nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allConnections, err := siteconnections.ExtractConnections(allPages)
	th.AssertNoErr(t, err)

	for _, connection := range allConnections {
		tools.PrintResource(t, connection)
	}
}

func TestSiteConnectionsCRUD(t *testing.T) {
	client := NewClient(t)

	networkClient, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)
	externalSubnetID := ExternalSubnetID(t, networkClient)

	ikePolicy, err := CreateIKEPolicy(t, client)
	th.AssertNoErr(t, err)
	defer DeleteIKEPolicy(t, client, ikePolicy.ID)

	ipsecPolicy, err := CreateIPsecPolicy(t, client)
	th.AssertNoErr(t, err)
	defer DeleteIPsecPolicy(t, client, ipsecPolicy.ID)

	localGroup, err := CreateEndpointGroup(t, client, []string{"10.1.0.0/24"})
	th.AssertNoErr(t, err)
	defer DeleteEndpointGroup(t, client, localGroup.ID)

	peerGroup, err := CreateEndpointGroup(t, client, []string{"10.2.0.0/24"})
	th.AssertNoErr(t, err)
	defer DeleteEndpointGroup(t, client, peerGroup.ID)

	gateway, err := CreateVPNGateway(t, client)
	th.AssertNoErr(t, err)
	defer DeleteVPNGateway(t, client, gateway.ID)

	_, err = vpngateways.AddExternalGateways(context.TODO(), client, gateway.ID, vpngateways.AddExternalGatewaysOpts{
		ExternalGateways: []vpngateways.ExternalGatewayOpts{
			{SubnetID: externalSubnetID},
		},
	}).Extract()
	th.AssertNoErr(t, err)
	_, err = WaitForAttachments(client, gateway.ID)
	th.AssertNoErr(t, err)

	site, err := CreateIPsecSite(t, client, ipsecsites.CreateOpts{
		VPNGatewayID:   gateway.ID,
		LocalEPGroupID: localGroup.ID,
		PeerEPGroupID:  peerGroup.ID,
	})
	th.AssertNoErr(t, err)
	defer DeleteIPsecSite(t, client, site.ID)

	// The peers are documentation addresses: ACTIVE means the tunnel is
	// configured on the gateway, not that it is up.
	adminStateUp := true
	primary, err := CreateSiteConnection(t, client, siteconnections.CreateOpts{
		Description:   "primary tunnel",
		IPsecSiteID:   site.ID,
		PeerAddress:   "198.51.100.20",
		PeerID:        "vpn.peer.example.com",
		AuthMode:      siteconnections.AuthModePSK,
		PSK:           tools.RandomString("psk-", 16),
		Priority:      1,
		AdminStateUp:  &adminStateUp,
		IKEPolicyID:   ikePolicy.ID,
		IPsecPolicyID: ipsecPolicy.ID,
		DPD: &siteconnections.DPDOpts{
			Action:   siteconnections.DPDActionHold,
			Interval: 10,
		},
		MTU:       1380,
		Initiator: siteconnections.InitiatorBiDirectional,
	})
	th.AssertNoErr(t, err)
	defer DeleteSiteConnection(t, client, primary.ID)

	tools.PrintResource(t, primary)

	th.AssertEquals(t, site.ProjectID, primary.ProjectID)
	th.AssertEquals(t, "vpn.peer.example.com", primary.PeerID)
	th.AssertEquals(t, 1, primary.Priority)
	th.AssertEquals(t, 1380, primary.MTU)
	th.AssertDeepEquals(t, siteconnections.DPD{Action: "hold", Interval: 10, Timeout: 150}, primary.DPD)

	// Left unset, the priority follows the primary's and the rest defaults.
	standby, err := CreateSiteConnection(t, client, siteconnections.CreateOpts{
		IPsecSiteID:   site.ID,
		PeerAddress:   "198.51.100.21",
		PSK:           tools.RandomString("psk-", 16),
		IKEPolicyID:   ikePolicy.ID,
		IPsecPolicyID: ipsecPolicy.ID,
	})
	th.AssertNoErr(t, err)
	defer DeleteSiteConnection(t, client, standby.ID)

	tools.PrintResource(t, standby)

	th.AssertEquals(t, 2, standby.Priority)
	th.AssertEquals(t, "198.51.100.21", standby.PeerID)
	th.AssertEquals(t, 1400, standby.MTU)
	th.AssertEquals(t, "bi-directional", standby.Initiator)
	th.AssertDeepEquals(t, siteconnections.DPD{Action: "restart", Interval: 30, Timeout: 150}, standby.DPD)

	// A site fails over between two tunnels at most.
	_, err = siteconnections.Create(context.TODO(), client, siteconnections.CreateOpts{
		IPsecSiteID:   site.ID,
		PeerAddress:   "198.51.100.22",
		PSK:           tools.RandomString("psk-", 16),
		IKEPolicyID:   ikePolicy.ID,
		IPsecPolicyID: ipsecPolicy.ID,
	}).Extract()
	if !gophercloud.ResponseCodeIs(err, http.StatusConflict) {
		t.Fatalf("Expected a 409 for a third connection, got %v", err)
	}

	listOpts := siteconnections.ListOpts{
		IPsecSiteID: site.ID,
		SortKey:     "priority",
		SortDir:     "asc",
	}
	allPages, err := siteconnections.List(client, listOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	allConnections, err := siteconnections.ExtractConnections(allPages)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 2, len(allConnections))
	th.AssertEquals(t, primary.ID, allConnections[0].ID)
	th.AssertEquals(t, standby.ID, allConnections[1].ID)

	name := tools.RandomString("TESTACC-", 8)
	updatedConnection, err := siteconnections.Update(context.TODO(), client, primary.ID, siteconnections.UpdateOpts{
		Name: &name,
		PSK:  tools.RandomString("psk-", 16),
		MTU:  1400,
		DPD: &siteconnections.DPDOpts{
			Action: siteconnections.DPDActionRestart,
		},
	}).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, siteconnections.StatusPendingUpdate, updatedConnection.Status)

	err = WaitForSiteConnectionStatus(client, primary.ID, siteconnections.StatusActive)
	th.AssertNoErr(t, err)

	updatedConnection, err = siteconnections.Get(context.TODO(), client, primary.ID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, updatedConnection)

	th.AssertEquals(t, name, updatedConnection.Name)
	th.AssertEquals(t, 1400, updatedConnection.MTU)
	th.AssertEquals(t, true, updatedConnection.PSKSet)
	th.AssertDeepEquals(t, siteconnections.DPD{Action: "restart", Interval: 10, Timeout: 150}, updatedConnection.DPD)
}
