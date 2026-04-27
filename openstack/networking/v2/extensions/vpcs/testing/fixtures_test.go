package testing

import (
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/vpcs"
)

// VpcsListResult represents raw response for the List request.
const VpcsListResult = `
{
    "vpcs": [
        {
            "id": "6f1a5c3e-2b4d-4c8e-9f0a-1b2c3d4e5f60",
            "name": "vpc-prod",
            "description": "Production VPC",
            "cidr_block": "10.0.0.0/16",
            "vpc_status": "ACTIVE",
            "tenant_id": "4fd44f30292945e481c7b8a0c8908869",
            "project_id": "4fd44f30292945e481c7b8a0c8908869"
        },
        {
            "id": "a3b2c1d0-9e8f-4a7b-8c6d-5e4f3a2b1c0d",
            "name": "vpc-dev",
            "description": "",
            "cidr_block": "172.16.0.0/20",
            "vpc_status": "PENDING_CREATE",
            "tenant_id": "4fd44f30292945e481c7b8a0c8908869",
            "project_id": "4fd44f30292945e481c7b8a0c8908869"
        }
    ]
}
`

// Vpc1 represents the first VPC of the List response.
var Vpc1 = vpcs.Vpc{
	ID:          "6f1a5c3e-2b4d-4c8e-9f0a-1b2c3d4e5f60",
	Name:        "vpc-prod",
	Description: "Production VPC",
	CIDRBlock:   "10.0.0.0/16",
	VpcStatus:   "ACTIVE",
	TenantID:    "4fd44f30292945e481c7b8a0c8908869",
	ProjectID:   "4fd44f30292945e481c7b8a0c8908869",
}

// Vpc2 represents the second VPC of the List response.
var Vpc2 = vpcs.Vpc{
	ID:          "a3b2c1d0-9e8f-4a7b-8c6d-5e4f3a2b1c0d",
	Name:        "vpc-dev",
	Description: "",
	CIDRBlock:   "172.16.0.0/20",
	VpcStatus:   "PENDING_CREATE",
	TenantID:    "4fd44f30292945e481c7b8a0c8908869",
	ProjectID:   "4fd44f30292945e481c7b8a0c8908869",
}

// VpcGetResult represents raw response for the Get request.
const VpcGetResult = `
{
    "vpc": {
        "id": "6f1a5c3e-2b4d-4c8e-9f0a-1b2c3d4e5f60",
        "name": "vpc-prod",
        "description": "Production VPC",
        "cidr_block": "10.0.0.0/16",
        "vpc_status": "ACTIVE",
        "tenant_id": "4fd44f30292945e481c7b8a0c8908869",
        "project_id": "4fd44f30292945e481c7b8a0c8908869"
    }
}
`

// VpcCreateRequest represents raw Create request.
const VpcCreateRequest = `
{
    "vpc": {
        "name": "vpc-prod",
        "description": "Production VPC",
        "cidr_block": "10.0.0.0/16",
        "project_id": "4fd44f30292945e481c7b8a0c8908869"
    }
}
`

// VpcCreateResult represents raw Create response.
const VpcCreateResult = `
{
    "vpc": {
        "id": "6f1a5c3e-2b4d-4c8e-9f0a-1b2c3d4e5f60",
        "name": "vpc-prod",
        "description": "Production VPC",
        "cidr_block": "10.0.0.0/16",
        "vpc_status": "PENDING_CREATE",
        "tenant_id": "4fd44f30292945e481c7b8a0c8908869",
        "project_id": "4fd44f30292945e481c7b8a0c8908869"
    }
}
`

// VpcUpdateRequest represents raw Update request.
const VpcUpdateRequest = `
{
    "vpc": {
        "name": "vpc-renamed",
        "description": ""
    }
}
`

// VpcUpdateResult represents raw Update response.
const VpcUpdateResult = `
{
    "vpc": {
        "id": "6f1a5c3e-2b4d-4c8e-9f0a-1b2c3d4e5f60",
        "name": "vpc-renamed",
        "description": "",
        "cidr_block": "10.0.0.0/16",
        "vpc_status": "PENDING_UPDATE",
        "tenant_id": "4fd44f30292945e481c7b8a0c8908869",
        "project_id": "4fd44f30292945e481c7b8a0c8908869"
    }
}
`
