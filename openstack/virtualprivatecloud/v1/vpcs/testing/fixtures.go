package testing

import (
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/virtualprivatecloud/v1/vpcs"
)

const ListResponse = `
{
    "vpcs": [
        {
            "id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
            "project_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
            "name": "my-vpc",
            "description": "Production VPC",
            "cidr_block": "10.0.0.0/16",
            "status": "ACTIVE",
            "created_at": "2026-03-04T12:00:00Z",
            "updated_at": "2026-03-04T12:05:00Z"
        },
        {
            "id": "b2c3d4e5-f6a7-8901-bcde-f12345678901",
            "project_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
            "name": "dev-vpc",
            "description": "",
            "cidr_block": "172.16.0.0/20",
            "status": "CREATING",
            "created_at": "2026-03-04T13:00:00Z",
            "updated_at": "2026-03-04T13:00:00Z"
        }
    ]
}`

const GetResponse = `
{
    "vpc": {
        "id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
        "project_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
        "name": "my-vpc",
        "description": "Production VPC",
        "cidr_block": "10.0.0.0/16",
        "status": "ACTIVE",
        "created_at": "2026-03-04T12:00:00Z",
        "updated_at": "2026-03-04T12:05:00Z"
    }
}`

const CreateRequest = `
{
    "vpc": {
        "name": "my-vpc",
        "cidr_block": "10.0.0.0/16",
        "description": "Production VPC"
    }
}`

const CreateResponse = `
{
    "vpc": {
        "id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
        "project_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
        "name": "my-vpc",
        "description": "Production VPC",
        "cidr_block": "10.0.0.0/16",
        "status": "CREATING",
        "created_at": "2026-03-04T12:00:00Z",
        "updated_at": "2026-03-04T12:00:00Z"
    }
}`

const UpdateRequest = `
{
    "vpc": {
        "name": "my-vpc-renamed"
    }
}`

const UpdateResponse = `
{
    "vpc": {
        "id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
        "project_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
        "name": "my-vpc-renamed",
        "description": "Production VPC",
        "cidr_block": "10.0.0.0/16",
        "status": "ACTIVE",
        "created_at": "2026-03-04T12:00:00Z",
        "updated_at": "2026-03-04T12:10:00Z"
    }
}`

const DeleteResponse = `
{
    "vpc": {
        "id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
        "project_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
        "name": "my-vpc",
        "description": "Production VPC",
        "cidr_block": "10.0.0.0/16",
        "status": "DELETING",
        "created_at": "2026-03-04T12:00:00Z",
        "updated_at": "2026-03-04T12:15:00Z"
    }
}`

var (
	createdTime1, _ = time.Parse(time.RFC3339, "2026-03-04T12:00:00Z")
	updatedTime1, _ = time.Parse(time.RFC3339, "2026-03-04T12:05:00Z")
	createdTime2, _ = time.Parse(time.RFC3339, "2026-03-04T13:00:00Z")
	updatedTime2, _ = time.Parse(time.RFC3339, "2026-03-04T13:00:00Z")
	updatedTime3, _ = time.Parse(time.RFC3339, "2026-03-04T12:10:00Z")
	updatedTime4, _ = time.Parse(time.RFC3339, "2026-03-04T12:15:00Z")
)

var Vpc1 = vpcs.Vpc{
	ID:          "f47ac10b-58cc-4372-a567-0e02b2c3d479",
	ProjectID:   "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
	Name:        "my-vpc",
	Description: "Production VPC",
	CIDRBlock:   "10.0.0.0/16",
	Status:      "ACTIVE",
	CreatedAt:   createdTime1,
	UpdatedAt:   updatedTime1,
}

var Vpc2 = vpcs.Vpc{
	ID:          "b2c3d4e5-f6a7-8901-bcde-f12345678901",
	ProjectID:   "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
	Name:        "dev-vpc",
	Description: "",
	CIDRBlock:   "172.16.0.0/20",
	Status:      "CREATING",
	CreatedAt:   createdTime2,
	UpdatedAt:   updatedTime2,
}

var ExpectedVpcSlice = []vpcs.Vpc{Vpc1, Vpc2}
