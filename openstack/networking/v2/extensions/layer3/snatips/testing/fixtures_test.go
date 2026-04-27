package testing

import (
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/snatips"
)

// SnatIPsListResult represents raw response for the List request.
const SnatIPsListResult = `
{
    "snat_ips": [
        {
            "id": "2c4e6a8b-1d3f-4a5c-8e7b-9d0f1a2b3c4d",
            "network_id": "8f9a1c2b-3d4e-4f5a-9b6c-7d8e9f0a1b2c",
            "availability_zone": "eu-west-par-a",
            "ip_address": "203.0.113.10",
            "description": "NAT for production",
            "vpc_id": "6f1a5c3e-2b4d-4c8e-9f0a-1b2c3d4e5f60",
            "status": "ACTIVE",
            "vpc_status": "ACTIVE",
            "tenant_id": "4fd44f30292945e481c7b8a0c8908869",
            "project_id": "4fd44f30292945e481c7b8a0c8908869"
        },
        {
            "id": "7e9f1a3b-5c6d-4e8f-a0b1-c2d3e4f5a6b7",
            "network_id": "8f9a1c2b-3d4e-4f5a-9b6c-7d8e9f0a1b2c",
            "availability_zone": "eu-west-par-b",
            "ip_address": "203.0.113.11",
            "description": "",
            "vpc_id": null,
            "status": "DOWN",
            "vpc_status": "ACTIVE",
            "tenant_id": "4fd44f30292945e481c7b8a0c8908869",
            "project_id": "4fd44f30292945e481c7b8a0c8908869"
        }
    ]
}
`

// SnatIP1 represents the first SNAT IP of the List response.
var SnatIP1 = snatips.SnatIP{
	ID:               "2c4e6a8b-1d3f-4a5c-8e7b-9d0f1a2b3c4d",
	NetworkID:        "8f9a1c2b-3d4e-4f5a-9b6c-7d8e9f0a1b2c",
	AvailabilityZone: "eu-west-par-a",
	IPAddress:        "203.0.113.10",
	Description:      "NAT for production",
	VpcID:            "6f1a5c3e-2b4d-4c8e-9f0a-1b2c3d4e5f60",
	Status:           "ACTIVE",
	VpcStatus:        "ACTIVE",
	TenantID:         "4fd44f30292945e481c7b8a0c8908869",
	ProjectID:        "4fd44f30292945e481c7b8a0c8908869",
}

// SnatIP2 represents the second SNAT IP of the List response.
var SnatIP2 = snatips.SnatIP{
	ID:               "7e9f1a3b-5c6d-4e8f-a0b1-c2d3e4f5a6b7",
	NetworkID:        "8f9a1c2b-3d4e-4f5a-9b6c-7d8e9f0a1b2c",
	AvailabilityZone: "eu-west-par-b",
	IPAddress:        "203.0.113.11",
	Description:      "",
	VpcID:            "",
	Status:           "DOWN",
	VpcStatus:        "ACTIVE",
	TenantID:         "4fd44f30292945e481c7b8a0c8908869",
	ProjectID:        "4fd44f30292945e481c7b8a0c8908869",
}

// SnatIPGetResult represents raw response for the Get request.
const SnatIPGetResult = `
{
    "snat_ip": {
        "id": "2c4e6a8b-1d3f-4a5c-8e7b-9d0f1a2b3c4d",
        "network_id": "8f9a1c2b-3d4e-4f5a-9b6c-7d8e9f0a1b2c",
        "availability_zone": "eu-west-par-a",
        "ip_address": "203.0.113.10",
        "description": "NAT for production",
        "vpc_id": "6f1a5c3e-2b4d-4c8e-9f0a-1b2c3d4e5f60",
        "status": "ACTIVE",
        "vpc_status": "ACTIVE",
        "tenant_id": "4fd44f30292945e481c7b8a0c8908869",
        "project_id": "4fd44f30292945e481c7b8a0c8908869"
    }
}
`

// SnatIPCreateRequest represents raw Create request.
const SnatIPCreateRequest = `
{
    "snat_ip": {
        "network_id": "8f9a1c2b-3d4e-4f5a-9b6c-7d8e9f0a1b2c",
        "availability_zone": "eu-west-par-a",
        "ip_address": "203.0.113.10",
        "description": "NAT for production",
        "project_id": "4fd44f30292945e481c7b8a0c8908869"
    }
}
`

// SnatIPCreateResult represents raw Create response.
const SnatIPCreateResult = `
{
    "snat_ip": {
        "id": "2c4e6a8b-1d3f-4a5c-8e7b-9d0f1a2b3c4d",
        "network_id": "8f9a1c2b-3d4e-4f5a-9b6c-7d8e9f0a1b2c",
        "availability_zone": "eu-west-par-a",
        "ip_address": "203.0.113.10",
        "description": "NAT for production",
        "vpc_id": null,
        "status": "DOWN",
        "vpc_status": "PENDING_CREATE",
        "tenant_id": "4fd44f30292945e481c7b8a0c8908869",
        "project_id": "4fd44f30292945e481c7b8a0c8908869"
    }
}
`

// SnatIPUpdateRequest represents raw Update request.
const SnatIPUpdateRequest = `
{
    "snat_ip": {
        "description": "NAT for staging"
    }
}
`

// SnatIPUpdateResult represents raw Update response.
const SnatIPUpdateResult = `
{
    "snat_ip": {
        "id": "2c4e6a8b-1d3f-4a5c-8e7b-9d0f1a2b3c4d",
        "network_id": "8f9a1c2b-3d4e-4f5a-9b6c-7d8e9f0a1b2c",
        "availability_zone": "eu-west-par-a",
        "ip_address": "203.0.113.10",
        "description": "NAT for staging",
        "vpc_id": "6f1a5c3e-2b4d-4c8e-9f0a-1b2c3d4e5f60",
        "status": "ACTIVE",
        "vpc_status": "ACTIVE",
        "tenant_id": "4fd44f30292945e481c7b8a0c8908869",
        "project_id": "4fd44f30292945e481c7b8a0c8908869"
    }
}
`
