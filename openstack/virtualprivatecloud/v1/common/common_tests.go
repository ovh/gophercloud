package common

import (
	"github.com/gophercloud/gophercloud/v2"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

const TokenID = client.TokenID

func ServiceClient(fakeServer th.FakeServer) *gophercloud.ServiceClient {
	return client.ServiceClient(fakeServer)
}
