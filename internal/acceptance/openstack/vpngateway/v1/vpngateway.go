// Package v1 contains common functions for creating VPN gateway resources
// for use in acceptance tests. See the `*_test.go` files for example usages.
package v1

import (
	"context"
	"errors"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/vpngateway/v1/endpointgroups"
	"github.com/gophercloud/gophercloud/v2/openstack/vpngateway/v1/ikepolicies"
	"github.com/gophercloud/gophercloud/v2/openstack/vpngateway/v1/ipsecpolicies"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

// NewClient returns a VPN gateway client, skipping the test when the cloud
// has no vpn-gateway endpoint in its catalog.
func NewClient(t *testing.T) *gophercloud.ServiceClient {
	client, err := clients.NewVPNGatewayV1Client()
	var notFound *gophercloud.ErrEndpointNotFound
	if errors.As(err, &notFound) {
		t.Skip("The cloud has no vpn-gateway service")
	}
	th.AssertNoErr(t, err)
	return client
}

// CreateIKEPolicy will create an IKE policy with a random name. An error
// will be returned if the policy could not be created.
func CreateIKEPolicy(t *testing.T, client *gophercloud.ServiceClient) (*ikepolicies.Policy, error) {
	policyName := tools.RandomString("TESTACC-", 8)

	t.Logf("Attempting to create IKE policy %s", policyName)

	createOpts := ikepolicies.CreateOpts{
		Name:                policyName,
		Description:         "IKE policy",
		AuthAlgorithm:       ikepolicies.AuthAlgorithmSHA384,
		EncryptionAlgorithm: ikepolicies.EncryptionAlgorithmAES256GCM,
		PFS:                 ikepolicies.PFSGroup19,
		Lifetime: &ikepolicies.LifetimeOpts{
			Units: ikepolicies.UnitSeconds,
			Value: 7200,
		},
	}

	policy, err := ikepolicies.Create(context.TODO(), client, createOpts).Extract()
	if err != nil {
		return policy, err
	}

	t.Logf("Successfully created IKE policy %s", policyName)

	th.AssertEquals(t, policyName, policy.Name)
	th.AssertEquals(t, "sha384", policy.AuthAlgorithm)
	th.AssertEquals(t, "aes-256-gcm-16", policy.EncryptionAlgorithm)
	th.AssertEquals(t, "group19", policy.PFS)
	th.AssertEquals(t, "main", policy.Phase1NegotiationMode)
	th.AssertEquals(t, "v2", policy.IKEVersion)
	th.AssertEquals(t, 7200, policy.Lifetime.Value)

	return policy, nil
}

// DeleteIKEPolicy will delete an IKE policy with a specified ID. A fatal
// error will occur if the delete was not successful. This works best when
// used as a deferred function.
func DeleteIKEPolicy(t *testing.T, client *gophercloud.ServiceClient, policyID string) {
	t.Logf("Attempting to delete IKE policy: %s", policyID)

	err := ikepolicies.Delete(context.TODO(), client, policyID).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete IKE policy %s: %v", policyID, err)
	}

	t.Logf("Deleted IKE policy: %s", policyID)
}

// CreateIPsecPolicy will create an AES-GCM IPsec policy with a random name.
// An error will be returned if the policy could not be created.
func CreateIPsecPolicy(t *testing.T, client *gophercloud.ServiceClient) (*ipsecpolicies.Policy, error) {
	policyName := tools.RandomString("TESTACC-", 8)

	t.Logf("Attempting to create IPsec policy %s", policyName)

	createOpts := ipsecpolicies.CreateOpts{
		Name:                policyName,
		Description:         "IPsec policy",
		TransformProtocol:   ipsecpolicies.TransformProtocolESP,
		EncryptionAlgorithm: ipsecpolicies.EncryptionAlgorithmAES256GCM,
		EncapsulationMode:   ipsecpolicies.EncapsulationModeTunnel,
		PFS:                 ipsecpolicies.PFSGroup20,
		Lifetime: &ipsecpolicies.LifetimeOpts{
			Units: ipsecpolicies.UnitSeconds,
			Value: 1800,
		},
	}

	policy, err := ipsecpolicies.Create(context.TODO(), client, createOpts).Extract()
	if err != nil {
		return policy, err
	}

	t.Logf("Successfully created IPsec policy %s", policyName)

	th.AssertEquals(t, policyName, policy.Name)
	th.AssertEquals(t, "aes-256-gcm-16", policy.EncryptionAlgorithm)
	th.AssertEquals(t, "group20", policy.PFS)
	th.AssertEquals(t, 1800, policy.Lifetime.Value)
	if policy.AuthAlgorithm != nil {
		t.Fatalf("Expected no auth algorithm with AES-GCM, got %s", *policy.AuthAlgorithm)
	}

	return policy, nil
}

// DeleteIPsecPolicy will delete an IPsec policy with a specified ID. A fatal
// error will occur if the delete was not successful. This works best when
// used as a deferred function.
func DeleteIPsecPolicy(t *testing.T, client *gophercloud.ServiceClient, policyID string) {
	t.Logf("Attempting to delete IPsec policy: %s", policyID)

	err := ipsecpolicies.Delete(context.TODO(), client, policyID).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete IPsec policy %s: %v", policyID, err)
	}

	t.Logf("Deleted IPsec policy: %s", policyID)
}

// CreateEndpointGroup will create an endpoint group of the given CIDRs with
// a random name. An error will be returned if the group could not be created.
func CreateEndpointGroup(t *testing.T, client *gophercloud.ServiceClient, cidrs []string) (*endpointgroups.EndpointGroup, error) {
	groupName := tools.RandomString("TESTACC-", 8)

	t.Logf("Attempting to create endpoint group %s", groupName)

	createOpts := endpointgroups.CreateOpts{
		Name:        groupName,
		Description: "endpoint group",
		Type:        endpointgroups.TypeCIDR,
		CIDRs:       cidrs,
	}

	group, err := endpointgroups.Create(context.TODO(), client, createOpts).Extract()
	if err != nil {
		return group, err
	}

	t.Logf("Successfully created endpoint group %s", groupName)

	th.AssertEquals(t, groupName, group.Name)
	th.AssertEquals(t, "cidr", group.Type)
	th.AssertDeepEquals(t, cidrs, group.CIDRs)

	return group, nil
}

// DeleteEndpointGroup will delete an endpoint group with a specified ID. A
// fatal error will occur if the delete was not successful. This works best
// when used as a deferred function.
func DeleteEndpointGroup(t *testing.T, client *gophercloud.ServiceClient, groupID string) {
	t.Logf("Attempting to delete endpoint group: %s", groupID)

	err := endpointgroups.Delete(context.TODO(), client, groupID).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete endpoint group %s: %v", groupID, err)
	}

	t.Logf("Deleted endpoint group: %s", groupID)
}
