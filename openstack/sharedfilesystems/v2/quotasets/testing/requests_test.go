package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/quotasets"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestGetQuotaSet(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGetQuotaSetSuccessfully(t, fakeServer)

	actual, err := quotasets.Get(t.Context(), client.ServiceClient(fakeServer), tenantID).Extract()
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, &quotasets.QuotaSet{
		Gigabytes:             gophercloud.IntToPointer(10),
		Snapshots:             gophercloud.IntToPointer(10),
		Shares:                gophercloud.IntToPointer(10),
		SnapshotGigabytes:     gophercloud.IntToPointer(10),
		ShareNetworks:         gophercloud.IntToPointer(10),
		ShareGroups:           gophercloud.IntToPointer(10),
		ShareGroupSnapshots:   gophercloud.IntToPointer(10),
		ShareReplicas:         gophercloud.IntToPointer(10),
		ShareReplicaGigabytes: gophercloud.IntToPointer(10),
		PerShareGigabytes:     gophercloud.IntToPointer(10),
	}, actual)
}

func TestUpdateQuotaSet(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleUpdateQuotaSetSuccessfully(t, fakeServer)

	actual, err := quotasets.Update(t.Context(), client.ServiceClient(fakeServer), tenantID, quotasets.UpdateOpts{
		Gigabytes:             gophercloud.IntToPointer(100),
		Snapshots:             gophercloud.IntToPointer(100),
		Shares:                gophercloud.IntToPointer(100),
		SnapshotGigabytes:     gophercloud.IntToPointer(100),
		ShareNetworks:         gophercloud.IntToPointer(100),
		ShareGroups:           gophercloud.IntToPointer(100),
		ShareGroupSnapshots:   gophercloud.IntToPointer(100),
		ShareReplicas:         gophercloud.IntToPointer(100),
		ShareReplicaGigabytes: gophercloud.IntToPointer(100),
		PerShareGigabytes:     gophercloud.IntToPointer(100),
	}).Extract()

	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &quotasets.QuotaSet{
		Gigabytes:             gophercloud.IntToPointer(100),
		Snapshots:             gophercloud.IntToPointer(100),
		Shares:                gophercloud.IntToPointer(100),
		SnapshotGigabytes:     gophercloud.IntToPointer(100),
		ShareNetworks:         gophercloud.IntToPointer(100),
		ShareGroups:           gophercloud.IntToPointer(100),
		ShareGroupSnapshots:   gophercloud.IntToPointer(100),
		ShareReplicas:         gophercloud.IntToPointer(100),
		ShareReplicaGigabytes: gophercloud.IntToPointer(100),
		PerShareGigabytes:     gophercloud.IntToPointer(100),
	}, actual)
}

func TestGetByShareType(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGetByShareTypeSuccessfully(t, fakeServer)

	actual, err := quotasets.GetByShareType(t.Context(), client.ServiceClient(fakeServer), tenantID, ShareType).Extract()
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, &quotasets.QuotaSet{
		Gigabytes:             gophercloud.IntToPointer(10),
		Snapshots:             gophercloud.IntToPointer(10),
		Shares:                gophercloud.IntToPointer(10),
		SnapshotGigabytes:     gophercloud.IntToPointer(10),
		ShareNetworks:         gophercloud.IntToPointer(10),
		ShareGroups:           gophercloud.IntToPointer(10),
		ShareGroupSnapshots:   gophercloud.IntToPointer(10),
		ShareReplicas:         gophercloud.IntToPointer(10),
		ShareReplicaGigabytes: gophercloud.IntToPointer(10),
		PerShareGigabytes:     gophercloud.IntToPointer(10),
	}, actual)
}

func TestUpdateByShareType(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleUpdateByShareTypeSuccessfully(t, fakeServer)

	actual, err := quotasets.UpdateByShareType(t.Context(), client.ServiceClient(fakeServer), tenantID, ShareType, quotasets.UpdateOpts{
		Gigabytes:             gophercloud.IntToPointer(100),
		Snapshots:             gophercloud.IntToPointer(100),
		Shares:                gophercloud.IntToPointer(100),
		SnapshotGigabytes:     gophercloud.IntToPointer(100),
		ShareNetworks:         gophercloud.IntToPointer(100),
		ShareGroups:           gophercloud.IntToPointer(100),
		ShareGroupSnapshots:   gophercloud.IntToPointer(100),
		ShareReplicas:         gophercloud.IntToPointer(100),
		ShareReplicaGigabytes: gophercloud.IntToPointer(100),
		PerShareGigabytes:     gophercloud.IntToPointer(100),
	}).Extract()
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, &quotasets.QuotaSet{
		Gigabytes:             gophercloud.IntToPointer(100),
		Snapshots:             gophercloud.IntToPointer(100),
		Shares:                gophercloud.IntToPointer(100),
		SnapshotGigabytes:     gophercloud.IntToPointer(100),
		ShareNetworks:         gophercloud.IntToPointer(100),
		ShareGroups:           gophercloud.IntToPointer(100),
		ShareGroupSnapshots:   gophercloud.IntToPointer(100),
		ShareReplicas:         gophercloud.IntToPointer(100),
		ShareReplicaGigabytes: gophercloud.IntToPointer(100),
		PerShareGigabytes:     gophercloud.IntToPointer(100),
	}, actual)
}

func TestGetByUser(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGetByUserSuccessfully(t, fakeServer)

	actual, err := quotasets.GetByUser(t.Context(), client.ServiceClient(fakeServer), tenantID, userID).Extract()
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, &quotasets.QuotaSet{
		Gigabytes:             gophercloud.IntToPointer(10),
		Snapshots:             gophercloud.IntToPointer(10),
		Shares:                gophercloud.IntToPointer(10),
		SnapshotGigabytes:     gophercloud.IntToPointer(10),
		ShareNetworks:         gophercloud.IntToPointer(10),
		ShareGroups:           gophercloud.IntToPointer(10),
		ShareGroupSnapshots:   gophercloud.IntToPointer(10),
		ShareReplicas:         gophercloud.IntToPointer(10),
		ShareReplicaGigabytes: gophercloud.IntToPointer(10),
		PerShareGigabytes:     gophercloud.IntToPointer(10),
	}, actual)
}

func TestUpdateByUser(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleUpdateByUserSuccessfully(t, fakeServer)

	actual, err := quotasets.UpdateByUser(t.Context(), client.ServiceClient(fakeServer), tenantID, userID, quotasets.UpdateOpts{
		Gigabytes:             gophercloud.IntToPointer(100),
		Snapshots:             gophercloud.IntToPointer(100),
		Shares:                gophercloud.IntToPointer(100),
		SnapshotGigabytes:     gophercloud.IntToPointer(100),
		ShareNetworks:         gophercloud.IntToPointer(100),
		ShareGroups:           gophercloud.IntToPointer(100),
		ShareGroupSnapshots:   gophercloud.IntToPointer(100),
		ShareReplicas:         gophercloud.IntToPointer(100),
		ShareReplicaGigabytes: gophercloud.IntToPointer(100),
		PerShareGigabytes:     gophercloud.IntToPointer(100),
	}).Extract()
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, &quotasets.QuotaSet{
		Gigabytes:             gophercloud.IntToPointer(100),
		Snapshots:             gophercloud.IntToPointer(100),
		Shares:                gophercloud.IntToPointer(100),
		SnapshotGigabytes:     gophercloud.IntToPointer(100),
		ShareNetworks:         gophercloud.IntToPointer(100),
		ShareGroups:           gophercloud.IntToPointer(100),
		ShareGroupSnapshots:   gophercloud.IntToPointer(100),
		ShareReplicas:         gophercloud.IntToPointer(100),
		ShareReplicaGigabytes: gophercloud.IntToPointer(100),
		PerShareGigabytes:     gophercloud.IntToPointer(100),
	}, actual)
}
