package v1

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/subnets"
	"github.com/gophercloud/gophercloud/v2/openstack/vpngateway/v1/vpngateways"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

// ExternalSubnetID returns the ID of an IPv4 subnet of the external network
// named by OS_EXTGW_ID, skipping the test when it is unset.
func ExternalSubnetID(t *testing.T, networkClient *gophercloud.ServiceClient) string {
	networkID := os.Getenv("OS_EXTGW_ID")
	if networkID == "" {
		t.Skip("OS_EXTGW_ID must be set to attach a VPN gateway to an external network")
	}

	allPages, err := subnets.List(networkClient, subnets.ListOpts{
		NetworkID: networkID,
		IPVersion: 4,
	}).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	allSubnets, err := subnets.ExtractSubnets(allPages)
	th.AssertNoErr(t, err)
	if len(allSubnets) == 0 {
		t.Fatalf("External network %s has no IPv4 subnet", networkID)
	}
	return allSubnets[0].ID
}

// CreateVPNGateway will create a VPN gateway with a random name and wait for
// it to be ACTIVE. An error will be returned if the gateway could not be
// created.
func CreateVPNGateway(t *testing.T, client *gophercloud.ServiceClient) (*vpngateways.VPNGateway, error) {
	gatewayName := tools.RandomString("TESTACC-", 8)

	t.Logf("Attempting to create VPN gateway %s", gatewayName)

	createOpts := vpngateways.CreateOpts{
		Name:        gatewayName,
		Description: "VPN gateway",
		Routes: []vpngateways.RouteOpts{
			{DestinationCIDR: "10.20.0.0/16", NextHop: "10.0.0.254"},
		},
	}

	gateway, err := vpngateways.Create(context.TODO(), client, createOpts).Extract()
	if err != nil {
		return gateway, err
	}

	if err := WaitForVPNGatewayStatus(client, gateway.ID, vpngateways.StatusActive); err != nil {
		return gateway, err
	}

	t.Logf("Successfully created VPN gateway %s", gatewayName)

	th.AssertEquals(t, gatewayName, gateway.Name)
	th.AssertEquals(t, "VPN gateway", gateway.Description)

	return vpngateways.Get(context.TODO(), client, gateway.ID).Extract()
}

// DeleteVPNGateway will delete a VPN gateway with a specified ID once it has
// settled, and wait for it to be gone. A fatal error will occur if the
// delete was not successful. This works best when used as a deferred
// function.
func DeleteVPNGateway(t *testing.T, client *gophercloud.ServiceClient, gatewayID string) {
	t.Logf("Attempting to delete VPN gateway: %s", gatewayID)

	err := tools.WaitFor(func(ctx context.Context) (bool, error) {
		gateway, err := vpngateways.Get(ctx, client, gatewayID).Extract()
		if err != nil {
			return false, err
		}
		return !strings.HasPrefix(gateway.Status, "PENDING_"), nil
	})
	if err != nil {
		t.Fatalf("VPN gateway %s did not settle: %v", gatewayID, err)
	}

	err = vpngateways.Delete(context.TODO(), client, gatewayID).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete VPN gateway %s: %v", gatewayID, err)
	}

	err = tools.WaitFor(func(ctx context.Context) (bool, error) {
		_, err := vpngateways.Get(ctx, client, gatewayID).Extract()
		if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
			return true, nil
		}
		return false, err
	})
	if err != nil {
		t.Fatalf("VPN gateway %s was not deleted: %v", gatewayID, err)
	}

	t.Logf("Deleted VPN gateway: %s", gatewayID)
}

// WaitForVPNGatewayStatus will wait for a VPN gateway to reach a status. It
// stops early when the gateway goes to ERROR instead.
func WaitForVPNGatewayStatus(client *gophercloud.ServiceClient, gatewayID, status string) error {
	return tools.WaitFor(func(ctx context.Context) (bool, error) {
		gateway, err := vpngateways.Get(ctx, client, gatewayID).Extract()
		if err != nil {
			return false, err
		}
		if gateway.Status == status {
			return true, nil
		}
		if gateway.Status == vpngateways.StatusError {
			return false, fmt.Errorf("VPN gateway %s went to ERROR", gatewayID)
		}
		return false, nil
	})
}

// WaitForAttachments will wait for every external gateway and internal
// interface of a VPN gateway to be ACTIVE, and returns the gateway.
func WaitForAttachments(client *gophercloud.ServiceClient, gatewayID string) (*vpngateways.VPNGateway, error) {
	var gateway *vpngateways.VPNGateway
	err := tools.WaitFor(func(ctx context.Context) (bool, error) {
		var err error
		gateway, err = vpngateways.Get(ctx, client, gatewayID).Extract()
		if err != nil {
			return false, err
		}
		attachments := append(append([]vpngateways.Attachment{}, gateway.ExternalGateways...), gateway.InternalInterfaces...)
		for _, attachment := range attachments {
			switch attachment.Status {
			case vpngateways.StatusActive:
			case vpngateways.StatusError:
				return false, fmt.Errorf("VPN gateway %s has an attachment in ERROR on network %s", gatewayID, attachment.NetworkID)
			default:
				return false, nil
			}
		}
		return true, nil
	})
	return gateway, err
}
