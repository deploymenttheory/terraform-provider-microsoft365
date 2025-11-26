---
page_title: "microsoft365_graph_beta_groups_license_assignment Resource - terraform-provider-microsoft365"
subcategory: "Groups"
description: |-
  Manages a single Microsoft 365 license assignment for a group using the /groups/{groupId}/assignLicense endpoint. Each resource instance manages one license (SKU) for a group. To assign multiple licenses to a group, create multiple instances of this resource with different SKU IDs. License assignments automatically apply to all current and future group members.
---

# microsoft365_graph_beta_groups_license_assignment (Resource)

Manages a single Microsoft 365 license assignment for a group using the `/groups/{groupId}/assignLicense` endpoint. Each resource instance manages one license (SKU) for a group. To assign multiple licenses to a group, create multiple instances of this resource with different SKU IDs. License assignments automatically apply to all current and future group members.

## Microsoft Documentation

- [group: assignLicense](https://learn.microsoft.com/en-us/graph/api/group-assignlicense?view=graph-rest-beta)
- [group resource type](https://learn.microsoft.com/en-us/graph/api/resources/group?view=graph-rest-beta)
- [subscribedSku resource type](https://learn.microsoft.com/en-us/graph/api/resources/subscribedsku?view=graph-rest-beta)

## API Permissions

The following API permissions are required in order to use this resource.

### Microsoft Graph

- **Application**: `LicenseAssignment.ReadWrite.All`, `Group.ReadWrite.All`, `Directory.ReadWrite.All`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.15.0-alpha | Experimental | Initial release |
| v0.37.0-alpha | Preview | fixed broken create/ delete operations and added full test harness|

## Example Usage

```terraform
# Example 1: Minimal - Assign a single Office 365 E3 license to a group
resource "microsoft365_graph_beta_groups_license_assignment" "e3_license" {
  group_id = "1132b215-826f-42a9-8cfe-1643d19d17fd"
  sku_id   = "6fd2c87f-b296-42f0-b197-1e91e994b900" # Office 365 E3
}

# Example 2: Single license with disabled service plans
resource "microsoft365_graph_beta_groups_license_assignment" "e3_custom" {
  group_id = "1132b215-826f-42a9-8cfe-1643d19d17fd"
  sku_id   = "6fd2c87f-b296-42f0-b197-1e91e994b900" # Office 365 E3

  # Disable specific service plans within the license
  disabled_plans = [
    "efb87545-963c-4e0d-99df-69c6916d9eb0", # Azure Information Protection Premium P1
    "9f431833-0334-42de-a7dc-70aa40db46db"  # Microsoft Stream
  ]
}

# Example 3: Assign multiple licenses to the same group
# Each license requires a separate resource instance
resource "microsoft365_graph_beta_groups_license_assignment" "group_e3" {
  group_id = "2243c326-937g-53f0-c9df-2e68f106b901"
  sku_id   = "6fd2c87f-b296-42f0-b197-1e91e994b900" # Office 365 E3
}

resource "microsoft365_graph_beta_groups_license_assignment" "group_ems_e5" {
  group_id = "2243c326-937g-53f0-c9df-2e68f106b901"
  sku_id   = "b05e124f-c7cc-45a0-a6aa-8cf78c946968" # Enterprise Mobility + Security E5

  disabled_plans = [
    "113feb6c-3fe4-4440-bddc-54d774bf0318", # Exchange Foundation
    "14ab5db5-e6c4-4b20-b4bc-13e36fd2227f"  # Intune for Education
  ]
}

resource "microsoft365_graph_beta_groups_license_assignment" "group_power_bi" {
  group_id = "2243c326-937g-53f0-c9df-2e68f106b901"
  sku_id   = "f30db892-07e9-47e9-837c-80727f46fd3d" # Power BI Free
}

# Example 4: Using a data source to get group ID dynamically
data "microsoft365_graph_beta_groups_group" "sales_team" {
  display_name = "Sales Team"
}

resource "microsoft365_graph_beta_groups_license_assignment" "sales_e3" {
  group_id = data.microsoft365_graph_beta_groups_group.sales_team.id
  sku_id   = "6fd2c87f-b296-42f0-b197-1e91e994b900" # Office 365 E3
}

# Example 5: Assign license to a newly created group
resource "microsoft365_graph_beta_groups_group" "engineering" {
  display_name     = "Engineering Team"
  mail_nickname    = "engineering"
  mail_enabled     = false
  security_enabled = true
}

resource "microsoft365_graph_beta_groups_license_assignment" "engineering_e5" {
  group_id = microsoft365_graph_beta_groups_group.engineering.id
  sku_id   = "c7df2760-2c81-4ef7-b578-5b5392b571df" # Office 365 E5

  disabled_plans = [
    "57ff2da0-773e-42df-b2af-ffb7a2317929" # Teams
  ]
}

# Example 6: With custom timeouts
resource "microsoft365_graph_beta_groups_license_assignment" "custom_timeout" {
  group_id = "4465e548-159i-75h2-e1fg-4g80h328d123"
  sku_id   = "c7df2760-2c81-4ef7-b578-5b5392b571df" # Office 365 E5

  timeouts {
    create = "300s"
    read   = "180s"
    update = "300s"
    delete = "300s"
  }
}

# Example 7: Department-based license assignments
# Marketing department gets Office 365 E3
resource "microsoft365_graph_beta_groups_license_assignment" "marketing_e3" {
  group_id = "marketing-group-uuid-here"
  sku_id   = "6fd2c87f-b296-42f0-b197-1e91e994b900" # Office 365 E3
}

# Engineering department gets Office 365 E5 with full features
resource "microsoft365_graph_beta_groups_license_assignment" "engineering_e5_full" {
  group_id = "engineering-group-uuid-here"
  sku_id   = "c7df2760-2c81-4ef7-b578-5b5392b571df" # Office 365 E5

  # No disabled plans - all features enabled
}

# Sales department gets E3 with some features disabled
resource "microsoft365_graph_beta_groups_license_assignment" "sales_e3_limited" {
  group_id = "sales-group-uuid-here"
  sku_id   = "6fd2c87f-b296-42f0-b197-1e91e994b900" # Office 365 E3

  disabled_plans = [
    "efb87545-963c-4e0d-99df-69c6916d9eb0", # Azure Information Protection
    "9f431833-0334-42de-a7dc-70aa40db46db", # Microsoft Stream
    "b737dad2-2f6c-4c65-90e3-ca563267e8b9"  # Yammer Enterprise
  ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `group_id` (String) The unique identifier (UUID) for the group.
- `sku_id` (String) The unique identifier (GUID) for the license SKU to assign to the group.

### Optional

- `disabled_plans` (Set of String) A collection of the unique identifiers for service plans to disable for this license.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `display_name` (String) The display name of the group. This is computed and read-only.
- `id` (String) The unique identifier for this license assignment resource. Format: `{group_id}_{sku_id}`.

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Important Notes

- **Group-Based Licensing**: This resource manages license assignments for Microsoft 365 groups using the [group: assignLicense](https://learn.microsoft.com/en-us/graph/api/group-assignlicense?view=graph-rest-beta&tabs=http) Microsoft Graph API.
- **Transitive Assignment**: Licenses assigned to groups are automatically assigned to all users in the group.
- **SKU IDs**: License SKU IDs are required to assign licenses. You can get available SKUs using the `GET /subscribedSkus` endpoint.
- **Service Plans**: Individual service plans within a license can be disabled using the `disabled_plans` attribute.
- **Atomic Operations**: License assignments are atomic - the API processes all additions and removals in a single operation.
- **Permissions**: Users must have appropriate permissions in Azure AD to assign licenses to groups.
- **Response**: The API returns a `202 Accepted` response indicating the operation was accepted for processing.

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
# Import using group object ID
terraform import microsoft365_graph_beta_group_license_assignment.example 00000000-0000-0000-0000-000000000000
``` 