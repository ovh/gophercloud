//go:build acceptance || vpngateway || endpointgroups

package v1

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/vpngateway/v1/endpointgroups"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestEndpointGroupsList(t *testing.T) {
	client := NewClient(t)

	allPages, err := endpointgroups.List(client, nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allGroups, err := endpointgroups.ExtractEndpointGroups(allPages)
	th.AssertNoErr(t, err)

	for _, group := range allGroups {
		tools.PrintResource(t, group)
	}
}

func TestEndpointGroupsCRUD(t *testing.T) {
	client := NewClient(t)

	group, err := CreateEndpointGroup(t, client, []string{"10.2.0.0/24", "fd00:2::/64"})
	th.AssertNoErr(t, err)
	defer DeleteEndpointGroup(t, client, group.ID)

	tools.PrintResource(t, group)

	newGroup, err := endpointgroups.Get(context.TODO(), client, group.ID).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, *group, *newGroup)

	listOpts := endpointgroups.ListOpts{
		Name:  group.Name,
		Limit: 1,
	}
	allPages, err := endpointgroups.List(client, listOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	allGroups, err := endpointgroups.ExtractEndpointGroups(allPages)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 1, len(allGroups))
	th.AssertEquals(t, group.ID, allGroups[0].ID)

	name := tools.RandomString("TESTACC-", 8)
	description := tools.RandomString("TESTACC-desc-", 8)
	updateOpts := endpointgroups.UpdateOpts{
		Name:        &name,
		Description: &description,
		CIDRs:       []string{"10.4.0.0/16"},
	}

	updatedGroup, err := endpointgroups.Update(context.TODO(), client, group.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, updatedGroup)

	th.AssertEquals(t, name, updatedGroup.Name)
	th.AssertEquals(t, description, updatedGroup.Description)
	th.AssertDeepEquals(t, []string{"10.4.0.0/16"}, updatedGroup.CIDRs)
}
