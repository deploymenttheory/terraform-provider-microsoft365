---
page_title: "microsoft365_graph_beta_users_user_license_assignment Resource - terraform-provider-microsoft365"
subcategory: "Users"
description: |-
  Manages a single Microsoft 365 license assignment for an individual user using the /users/{userId}/assignLicense endpoint. Each resource instance manages one license (SKU) for a user. To assign multiple licenses to a user, create multiple instances of this resource with different SKU IDs.
---

# microsoft365_graph_beta_users_user_license_assignment (Resource)

Manages a single Microsoft 365 license assignment for an individual user using the `/users/{userId}/assignLicense` endpoint. Each resource instance manages one license (SKU) for a user. To assign multiple licenses to a user, create multiple instances of this resource with different SKU IDs.

## Microsoft Documentation

- [user: assignLicense](https://learn.microsoft.com/en-us/graph/api/user-assignlicense?view=graph-rest-beta)
- [user resource type](https://learn.microsoft.com/en-us/graph/api/resources/user?view=graph-rest-beta)
- [subscribedSku resource type](https://learn.microsoft.com/en-us/graph/api/resources/subscribedsku?view=graph-rest-beta)

## API Permissions

The following API permissions are required in order to use this resource.

### Microsoft Graph

- **Application**: `User.ReadWrite.All`, `Directory.ReadWrite.All`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.15.0-alpha | Experimental | Initial release |
| v0.37.0-alpha | Experimental | Refactored to make resource atomic and added full test harness|

## Example Usage

```terraform
# Example 1: Assign a single Office 365 E3 license to a user with disabled service plans
resource "microsoft365_graph_beta_users_user_license_assignment" "user_e3_license" {
  user_id = "john.doe@example.com" # Can be user ID (UUID) or UPN

  sku_id = "6fd2c87f-b296-42f0-b197-1e91e994b900" # Office 365 E3

  # Optional: Disable specific service plans within this license
  disabled_plans = [
    "efb87545-963c-4e0d-99df-69c6916d9eb0" # Example: Microsoft Stream
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

# Example 2: Assign a license without disabling any service plans
resource "microsoft365_graph_beta_users_user_license_assignment" "user_powerbi_license" {
  user_id = "john.doe@example.com"
  sku_id  = "f30db892-07e9-47e9-837c-80727f46fd3d" # Power BI (free)
}

# Example 3: Assign multiple licenses to a single user
# Note: Create multiple resource instances, one per license
resource "microsoft365_graph_beta_users_user_license_assignment" "jane_e3" {
  user_id = "jane.smith@example.com"
  sku_id  = "6fd2c87f-b296-42f0-b197-1e91e994b900" # Office 365 E3
}

resource "microsoft365_graph_beta_users_user_license_assignment" "jane_ems_e5" {
  user_id = "jane.smith@example.com"
  sku_id  = "b05e124f-c7cc-45a0-a6aa-8cf78c946968" # Enterprise Mobility + Security E5

  disabled_plans = [
    "8a256a2b-b617-496d-b51b-e76466e88db0" # Microsoft Defender for Cloud Apps
  ]
}

# Example 4: Assign Office 365 E5 license with multiple disabled plans
resource "microsoft365_graph_beta_users_user_license_assignment" "alice_e5" {
  user_id = "alice.wilson@example.com"
  sku_id  = "c7df2760-2c81-4ef7-b578-5b5392b571df" # Office 365 E5

  disabled_plans = [
    "57ff2da0-773e-42df-b2af-ffb7a2317929", # Teams
    "0feaeb32-d00e-4d66-bd5a-43b5b83db82c"  # Mya
  ]
}

# Example 5: Using a data source to get user ID dynamically
data "microsoft365_graph_beta_users_user" "target_user" {
  user_principal_name = "dynamic.user@example.com"
}

resource "microsoft365_graph_beta_users_user_license_assignment" "dynamic_user_license" {
  user_id = data.microsoft365_graph_beta_users_user.target_user.id
  sku_id  = "6fd2c87f-b296-42f0-b197-1e91e994b900" # Office 365 E3
}

# Example 6: Using for_each to assign the same license to multiple users
variable "licensed_users" {
  type = set(string)
  default = [
    "user1@example.com",
    "user2@example.com",
    "user3@example.com"
  ]
}

resource "microsoft365_graph_beta_users_user_license_assignment" "bulk_e3_assignment" {
  for_each = var.licensed_users

  user_id = each.value
  sku_id  = "6fd2c87f-b296-42f0-b197-1e91e994b900" # Office 365 E3
}

# Example 7: Create user and assign license in one configuration
resource "microsoft365_graph_beta_users_user" "new_user" {
  user_principal_name = "new.employee@example.com"
  display_name        = "New Employee"
  mail_nickname       = "new.employee"
  account_enabled     = true
  usage_location      = "US"

  password_profile = {
    password                           = "TemporaryP@ssw0rd123!"
    force_change_password_next_sign_in = true
  }
}

resource "microsoft365_graph_beta_users_user_license_assignment" "new_user_license" {
  user_id = microsoft365_graph_beta_users_user.new_user.id
  sku_id  = "6fd2c87f-b296-42f0-b197-1e91e994b900" # Office 365 E3

  depends_on = [microsoft365_graph_beta_users_user.new_user]
}

# Common Microsoft 365 SKU IDs for reference:
# - Office 365 E3: 6fd2c87f-b296-42f0-b197-1e91e994b900
# - Office 365 E5: c7df2760-2c81-4ef7-b578-5b5392b571df
# - Enterprise Mobility + Security E5: b05e124f-c7cc-45a0-a6aa-8cf78c946968
# - Microsoft 365 Business Premium: cbdc14ab-d96c-4c30-b9f4-6ada7cdc1d46
# - Power BI (free): f30db892-07e9-47e9-837c-80727f46fd3d

# Note: To remove a license, simply destroy the resource:
# terraform destroy -target=microsoft365_graph_beta_users_user_license_assignment.user_e3_license
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `sku_id` (String) The unique identifier (GUID) for the license SKU to assign to the user.
- `user_id` (String) The unique identifier for the user. Can be either the object ID (UUID) or user principal name (UPN).

### Optional

- `disabled_plans` (Set of String) A collection of the unique identifiers for service plans to disable for this license.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `id` (String) The unique identifier for this license assignment resource. Format: `{user_id}_{sku_id}`.
- `user_principal_name` (String) The user principal name (UPN) of the user. This is computed and read-only.

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Important Notes

- **License Management**: This resource manages license assignments for Microsoft 365 users using the [assignLicense](https://learn.microsoft.com/en-us/graph/api/user-assignlicense?view=graph-rest-beta&tabs=http) Microsoft Graph API.
- **SKU IDs**: License SKU IDs are required to assign licenses. You can get available SKUs using the `GET /subscribedSkus` endpoint.
- **Service Plans**: Individual service plans within a license can be disabled using the `disabled_plans` attribute.
- **Atomic Operations**: License assignments are atomic - the API processes all additions and removals in a single operation.
- **Permissions**: Users must have appropriate permissions in Azure AD to assign licenses to other users.

## Common License SKU IDs

Here are some common Microsoft 365 license SKU IDs:

| License | SKU ID |
|---------|--------|
| Office 365 E1 | `18181a46-0d4e-45cd-891e-60aabd171b4e` |
| Office 365 E3 | `6fd2c87f-b296-42f0-b197-1e91e994b900` |
| Office 365 E5 | `c7df2760-2c81-4ef7-b578-5b5392b571df` |
| Microsoft 365 E3 | `05e9a617-0261-4cee-bb44-138d3ef5d965` |
| Microsoft 365 E5 | `06ebc4ee-1bb5-47dd-8120-11324bc54e06` |
| Enterprise Mobility + Security E3 | `efccb6f7-5641-4e0e-bd10-b4976e1bf68e` |
| Enterprise Mobility + Security E5 | `b05e124f-c7cc-45a0-a6aa-8cf78c946968` |

## Import

Import is supported using the following syntax:

```shell
#!/bin/bash
# Import format: {user_id}_{sku_id}

# Import using user object ID and license SKU ID
terraform import microsoft365_graph_beta_users_user_license_assignment.example "12345678-1234-1234-1234-123456789012_6fd2c87f-b296-42f0-b197-1e91e994b900"
``` 