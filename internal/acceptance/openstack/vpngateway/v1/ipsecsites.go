package v1

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/vpngateway/v1/ipsecsites"
	"github.com/gophercloud/gophercloud/v2/openstack/vpngateway/v1/siteconnections"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

// CreateIPsecSite will create an IPsec site from createOpts with a random
// name. An error will be returned if the site could not be created.
func CreateIPsecSite(t *testing.T, client *gophercloud.ServiceClient, createOpts ipsecsites.CreateOpts) (*ipsecsites.Site, error) {
	createOpts.Name = tools.RandomString("TESTACC-", 8)

	t.Logf("Attempting to create IPsec site %s", createOpts.Name)

	site, err := ipsecsites.Create(context.TODO(), client, createOpts).Extract()
	if err != nil {
		return site, err
	}

	t.Logf("Successfully created IPsec site %s", createOpts.Name)

	th.AssertEquals(t, createOpts.Name, site.Name)
	th.AssertEquals(t, createOpts.VPNGatewayID, site.VPNGatewayID)
	th.AssertEquals(t, createOpts.LocalEPGroupID, site.LocalEPGroupID)
	th.AssertEquals(t, createOpts.PeerEPGroupID, site.PeerEPGroupID)

	return site, nil
}

// DeleteIPsecSite will delete an IPsec site with a specified ID. A fatal
// error will occur if the delete was not successful. This works best when
// used as a deferred function.
func DeleteIPsecSite(t *testing.T, client *gophercloud.ServiceClient, siteID string) {
	t.Logf("Attempting to delete IPsec site: %s", siteID)

	err := ipsecsites.Delete(context.TODO(), client, siteID).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete IPsec site %s: %v", siteID, err)
	}

	t.Logf("Deleted IPsec site: %s", siteID)
}

// CreateSiteConnection will create an IPsec site connection from createOpts
// with a random name and wait for it to be ACTIVE. An error will be returned
// if the connection could not be created.
func CreateSiteConnection(t *testing.T, client *gophercloud.ServiceClient, createOpts siteconnections.CreateOpts) (*siteconnections.Connection, error) {
	createOpts.Name = tools.RandomString("TESTACC-", 8)

	t.Logf("Attempting to create IPsec site connection %s", createOpts.Name)

	connection, err := siteconnections.Create(context.TODO(), client, createOpts).Extract()
	if err != nil {
		return connection, err
	}

	th.AssertEquals(t, createOpts.Name, connection.Name)
	th.AssertEquals(t, true, connection.PSKSet)

	if err := WaitForSiteConnectionStatus(client, connection.ID, siteconnections.StatusActive); err != nil {
		return connection, err
	}

	t.Logf("Successfully created IPsec site connection %s", createOpts.Name)

	return siteconnections.Get(context.TODO(), client, connection.ID).Extract()
}

// DeleteSiteConnection will delete an IPsec site connection with a specified
// ID once it has settled, and wait for it to be gone. A fatal error will
// occur if the delete was not successful. This works best when used as a
// deferred function.
func DeleteSiteConnection(t *testing.T, client *gophercloud.ServiceClient, connectionID string) {
	t.Logf("Attempting to delete IPsec site connection: %s", connectionID)

	err := tools.WaitFor(func(ctx context.Context) (bool, error) {
		connection, err := siteconnections.Get(ctx, client, connectionID).Extract()
		if err != nil {
			return false, err
		}
		return !strings.HasPrefix(connection.Status, "PENDING_"), nil
	})
	if err != nil {
		t.Fatalf("IPsec site connection %s did not settle: %v", connectionID, err)
	}

	err = siteconnections.Delete(context.TODO(), client, connectionID).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete IPsec site connection %s: %v", connectionID, err)
	}

	err = tools.WaitFor(func(ctx context.Context) (bool, error) {
		_, err := siteconnections.Get(ctx, client, connectionID).Extract()
		if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
			return true, nil
		}
		return false, err
	})
	if err != nil {
		t.Fatalf("IPsec site connection %s was not deleted: %v", connectionID, err)
	}

	t.Logf("Deleted IPsec site connection: %s", connectionID)
}

// WaitForSiteConnectionStatus will wait for an IPsec site connection to
// reach a status. It stops early when the connection goes to ERROR instead.
func WaitForSiteConnectionStatus(client *gophercloud.ServiceClient, connectionID, status string) error {
	return tools.WaitFor(func(ctx context.Context) (bool, error) {
		connection, err := siteconnections.Get(ctx, client, connectionID).Extract()
		if err != nil {
			return false, err
		}
		if connection.Status == status {
			return true, nil
		}
		if connection.Status == siteconnections.StatusError {
			return false, fmt.Errorf("IPsec site connection %s went to ERROR", connectionID)
		}
		return false, nil
	})
}
