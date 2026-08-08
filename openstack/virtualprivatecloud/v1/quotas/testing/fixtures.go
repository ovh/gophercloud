package testing

const GetResponseDefault = `
{
    "quota": {
        "project_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
        "max_vpcs": 5,
        "source": "config_default"
    }
}`

const GetResponseOverride = `
{
    "quota": {
        "project_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
        "max_vpcs": 10,
        "source": "project_override"
    }
}`

const SetRequest = `
{
    "quota": {
        "max_vpcs": 10
    }
}`

const SetResponse = `
{
    "quota": {
        "project_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
        "max_vpcs": 10
    }
}`
