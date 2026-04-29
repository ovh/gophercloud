package testing

import (
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/routetables"
)

// RouteTablesListResult represents raw response for the List request.
const RouteTablesListResult = `
{
    "route_tables": [
        {
            "id": "e3f5f925-6c6c-4b1e-9c1e-7c6d4d7c3f2a",
            "name": "public",
            "description": "Public route table",
            "router_id": "a5c3a4d4-5e0f-4c2a-bb02-e2723de5b8f1",
            "routes": [
                {"destination": "10.0.0.0/16", "nexthop": null, "nexthop_type": "local"},
                {"destination": "0.0.0.0/0", "nexthop": "10.0.0.1", "nexthop_type": "ip"}
            ],
            "subnets": ["08eae331-0402-425a-923c-34f7cfe39c1b"],
            "vpc_status": "ACTIVE",
            "tenant_id": "4fd44f30292945e481c7b8a0c8908869",
            "project_id": "4fd44f30292945e481c7b8a0c8908869"
        },
        {
            "id": "b2f3c4d5-1a2b-3c4d-5e6f-7a8b9c0d1e2f",
            "name": "private",
            "description": "",
            "router_id": "a5c3a4d4-5e0f-4c2a-bb02-e2723de5b8f1",
            "routes": [],
            "subnets": [],
            "vpc_status": "PENDING_CREATE",
            "tenant_id": "4fd44f30292945e481c7b8a0c8908869",
            "project_id": "4fd44f30292945e481c7b8a0c8908869"
        }
    ]
}
`

// RouteTable1 represents the first route table of the List response.
var RouteTable1 = routetables.RouteTable{
	ID:          "e3f5f925-6c6c-4b1e-9c1e-7c6d4d7c3f2a",
	Name:        "public",
	Description: "Public route table",
	RouterID:    "a5c3a4d4-5e0f-4c2a-bb02-e2723de5b8f1",
	Routes: []routetables.Route{
		{Destination: "10.0.0.0/16", Nexthop: "", NexthopType: "local"},
		{Destination: "0.0.0.0/0", Nexthop: "10.0.0.1", NexthopType: "ip"},
	},
	Subnets:   []string{"08eae331-0402-425a-923c-34f7cfe39c1b"},
	VpcStatus: "ACTIVE",
	TenantID:  "4fd44f30292945e481c7b8a0c8908869",
	ProjectID: "4fd44f30292945e481c7b8a0c8908869",
}

// RouteTable2 represents the second route table of the List response.
var RouteTable2 = routetables.RouteTable{
	ID:          "b2f3c4d5-1a2b-3c4d-5e6f-7a8b9c0d1e2f",
	Name:        "private",
	Description: "",
	RouterID:    "a5c3a4d4-5e0f-4c2a-bb02-e2723de5b8f1",
	Routes:      []routetables.Route{},
	Subnets:     []string{},
	VpcStatus:   "PENDING_CREATE",
	TenantID:    "4fd44f30292945e481c7b8a0c8908869",
	ProjectID:   "4fd44f30292945e481c7b8a0c8908869",
}

// RouteTableGetResult represents raw response for the Get request.
const RouteTableGetResult = `
{
    "route_table": {
        "id": "e3f5f925-6c6c-4b1e-9c1e-7c6d4d7c3f2a",
        "name": "public",
        "description": "Public route table",
        "router_id": "a5c3a4d4-5e0f-4c2a-bb02-e2723de5b8f1",
        "routes": [
            {"destination": "10.0.0.0/16", "nexthop": null, "nexthop_type": "local"},
            {"destination": "0.0.0.0/0", "nexthop": "10.0.0.1", "nexthop_type": "ip"}
        ],
        "subnets": ["08eae331-0402-425a-923c-34f7cfe39c1b"],
        "vpc_status": "ACTIVE",
        "tenant_id": "4fd44f30292945e481c7b8a0c8908869",
        "project_id": "4fd44f30292945e481c7b8a0c8908869"
    }
}
`

// RouteTableCreateRequest represents raw Create request.
const RouteTableCreateRequest = `
{
    "route_table": {
        "name": "private",
        "description": "Private route table",
        "router_id": "a5c3a4d4-5e0f-4c2a-bb02-e2723de5b8f1",
        "project_id": "4fd44f30292945e481c7b8a0c8908869"
    }
}
`

// RouteTableCreateResult represents raw Create response.
const RouteTableCreateResult = `
{
    "route_table": {
        "id": "b2f3c4d5-1a2b-3c4d-5e6f-7a8b9c0d1e2f",
        "name": "private",
        "description": "Private route table",
        "router_id": "a5c3a4d4-5e0f-4c2a-bb02-e2723de5b8f1",
        "routes": [],
        "subnets": [],
        "vpc_status": "PENDING_CREATE",
        "tenant_id": "4fd44f30292945e481c7b8a0c8908869",
        "project_id": "4fd44f30292945e481c7b8a0c8908869"
    }
}
`

// RouteTableUpdateRequest represents raw Update request.
const RouteTableUpdateRequest = `
{
    "route_table": {
        "name": "renamed",
        "description": ""
    }
}
`

// RouteTableUpdateResult represents raw Update response.
const RouteTableUpdateResult = `
{
    "route_table": {
        "id": "b2f3c4d5-1a2b-3c4d-5e6f-7a8b9c0d1e2f",
        "name": "renamed",
        "description": "",
        "router_id": "a5c3a4d4-5e0f-4c2a-bb02-e2723de5b8f1",
        "routes": [],
        "subnets": [],
        "vpc_status": "PENDING_UPDATE",
        "tenant_id": "4fd44f30292945e481c7b8a0c8908869",
        "project_id": "4fd44f30292945e481c7b8a0c8908869"
    }
}
`

// AddRoutesRequest represents raw AddRoutes request.
const AddRoutesRequest = `
{
    "routes": [
        {"destination": "192.168.1.0/24", "nexthop": "10.0.0.1", "nexthop_type": "ip"},
        {"destination": "192.168.2.0/24", "nexthop_type": "blackhole"}
    ]
}
`

// AddRoutesResult represents raw AddRoutes response.
const AddRoutesResult = `
{
    "id": "b2f3c4d5-1a2b-3c4d-5e6f-7a8b9c0d1e2f",
    "name": "private",
    "description": "",
    "router_id": "a5c3a4d4-5e0f-4c2a-bb02-e2723de5b8f1",
    "routes": [
        {"destination": "192.168.1.0/24", "nexthop": "10.0.0.1", "nexthop_type": "ip"},
        {"destination": "192.168.2.0/24", "nexthop": null, "nexthop_type": "blackhole"}
    ],
    "subnets": [],
    "vpc_status": "PENDING_UPDATE",
    "tenant_id": "4fd44f30292945e481c7b8a0c8908869",
    "project_id": "4fd44f30292945e481c7b8a0c8908869"
}
`

// RemoveRoutesRequest represents raw RemoveRoutes request.
const RemoveRoutesRequest = `
{
    "routes": [
        {"destination": "192.168.1.0/24"}
    ]
}
`

// RemoveRoutesResult represents raw RemoveRoutes response.
const RemoveRoutesResult = `
{
    "id": "b2f3c4d5-1a2b-3c4d-5e6f-7a8b9c0d1e2f",
    "name": "private",
    "description": "",
    "router_id": "a5c3a4d4-5e0f-4c2a-bb02-e2723de5b8f1",
    "routes": [
        {"destination": "192.168.2.0/24", "nexthop": null, "nexthop_type": "blackhole"}
    ],
    "subnets": [],
    "vpc_status": "PENDING_UPDATE",
    "tenant_id": "4fd44f30292945e481c7b8a0c8908869",
    "project_id": "4fd44f30292945e481c7b8a0c8908869"
}
`

// AddSubnetsRequest represents raw AddSubnets request.
const AddSubnetsRequest = `
{
    "subnets": ["08eae331-0402-425a-923c-34f7cfe39c1b"],
    "routes": [
        {"destination": "10.0.1.0/24", "nexthop_type": "local"}
    ]
}
`

// AddSubnetsResult represents raw AddSubnets response.
const AddSubnetsResult = `
{
    "id": "b2f3c4d5-1a2b-3c4d-5e6f-7a8b9c0d1e2f",
    "name": "private",
    "description": "",
    "router_id": "a5c3a4d4-5e0f-4c2a-bb02-e2723de5b8f1",
    "routes": [
        {"destination": "10.0.1.0/24", "nexthop": null, "nexthop_type": "local"}
    ],
    "subnets": ["08eae331-0402-425a-923c-34f7cfe39c1b"],
    "vpc_status": "PENDING_UPDATE",
    "tenant_id": "4fd44f30292945e481c7b8a0c8908869",
    "project_id": "4fd44f30292945e481c7b8a0c8908869"
}
`

// AddSubnetsWithoutRoutesRequest represents raw AddSubnets request without
// routes.
const AddSubnetsWithoutRoutesRequest = `
{
    "subnets": ["08eae331-0402-425a-923c-34f7cfe39c1b"]
}
`

// RemoveSubnetsRequest represents raw RemoveSubnets request.
const RemoveSubnetsRequest = `
{
    "subnets": ["08eae331-0402-425a-923c-34f7cfe39c1b"],
    "routes": [
        {"destination": "10.0.1.0/24"}
    ]
}
`

// RemoveSubnetsResult represents raw RemoveSubnets response.
const RemoveSubnetsResult = `
{
    "id": "b2f3c4d5-1a2b-3c4d-5e6f-7a8b9c0d1e2f",
    "name": "private",
    "description": "",
    "router_id": "a5c3a4d4-5e0f-4c2a-bb02-e2723de5b8f1",
    "routes": [],
    "subnets": [],
    "vpc_status": "PENDING_UPDATE",
    "tenant_id": "4fd44f30292945e481c7b8a0c8908869",
    "project_id": "4fd44f30292945e481c7b8a0c8908869"
}
`
