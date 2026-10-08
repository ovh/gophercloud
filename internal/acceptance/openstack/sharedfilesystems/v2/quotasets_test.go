//go:build acceptance
// +build acceptance

package v2

import (
	"context"
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/projects"
	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/quotasets"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestGet(t *testing.T) {
	client, err := clients.NewSharedFileSystemV2Client()
	th.AssertNoErr(t, err)

	projectID, err := currentProjectID()
	th.AssertNoErr(t, err)

	// Get the quotaset for the current tenant
	quotaset, err := quotasets.Get(context.TODO(), client, projectID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, quotaset)
}

func TestUpdate(t *testing.T) {
	client, err := clients.NewSharedFileSystemV2Client()
	th.AssertNoErr(t, err)

	projectID, err := currentProjectID()
	th.AssertNoErr(t, err)

	// Get the quotaset for the current tenant
	_, err = quotasets.Get(context.TODO(), client, projectID).Extract()
	th.AssertNoErr(t, err)

	// Update the quotaset
	updateOpts := quotasets.UpdateOpts{
		Gigabytes: gophercloud.IntToPointer(100),
	}
	quotaset, err := quotasets.Update(context.TODO(), client, projectID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, quotaset)
}

func TestGetByShareType(t *testing.T) {
	client, err := clients.NewSharedFileSystemV2Client()
	th.AssertNoErr(t, err)

	projectID, err := currentProjectID()
	th.AssertNoErr(t, err)

	// Get the quotaset for the current tenant
	quotaset, err := quotasets.GetByShareType(context.TODO(), client, projectID, "default").Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, quotaset)
}

func TestUpdateByShareType(t *testing.T) {
	client, err := clients.NewSharedFileSystemV2Client()
	th.AssertNoErr(t, err)

	projectID, err := currentProjectID()
	th.AssertNoErr(t, err)

	// Get the quotaset for the current tenant
	_, err = quotasets.GetByShareType(context.TODO(), client, projectID, "default").Extract()
	th.AssertNoErr(t, err)

	// Update the quotaset
	updateOpts := quotasets.UpdateOpts{
		Gigabytes: gophercloud.IntToPointer(100),
	}
	quotaset, err := quotasets.UpdateByShareType(context.TODO(), client, projectID, "default", updateOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, quotaset)
}

func TestGetByUser(t *testing.T) {
	client, err := clients.NewSharedFileSystemV2Client()
	th.AssertNoErr(t, err)

	projectID, err := currentProjectID()
	th.AssertNoErr(t, err)

	// Get the quotaset for the current tenant
	quotaset, err := quotasets.GetByUser(context.TODO(), client, projectID, "admin").Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, quotaset)
}

func TestUpdateByUser(t *testing.T) {
	client, err := clients.NewSharedFileSystemV2Client()
	th.AssertNoErr(t, err)

	projectID, err := currentProjectID()
	th.AssertNoErr(t, err)

	// Get the quotaset for the current tenant
	_, err = quotasets.GetByUser(context.TODO(), client, projectID, "admin").Extract()
	th.AssertNoErr(t, err)

	// Update the quotaset
	updateOpts := quotasets.UpdateOpts{
		Gigabytes: gophercloud.IntToPointer(100),
	}
	quotaset, err := quotasets.UpdateByUser(context.TODO(), client, projectID, "admin", updateOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, quotaset)
}

// currentProjectID returns the ID of the project the acceptance credentials
// are scoped to (the v2 ServiceClient no longer carries a TenantID).
func currentProjectID() (string, error) {
	identityClient, err := clients.NewIdentityV3Client()
	if err != nil {
		return "", err
	}
	allPages, err := projects.ListAvailable(identityClient).AllPages(context.TODO())
	if err != nil {
		return "", err
	}
	allProjects, err := projects.ExtractProjects(allPages)
	if err != nil {
		return "", err
	}
	if len(allProjects) == 0 {
		return "", fmt.Errorf("unable to get project ID")
	}
	return allProjects[0].ID, nil
}
