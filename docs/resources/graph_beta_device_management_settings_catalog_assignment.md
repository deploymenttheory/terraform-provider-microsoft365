---
page_title: "microsoft365_graph_beta_device_management_settings_catalog_assignment Resource - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Manages Device Management Configuration Policy Assignments in Microsoft Intune.
---

# microsoft365_graph_beta_device_management_settings_catalog_assignment (Resource)

Manages Device Management Configuration Policy Assignments in Microsoft Intune.

## Microsoft Documentation

- [deviceManagementConfigurationPolicyAssignment resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-devicemanagementconfigurationpolicyassignment?view=graph-rest-beta)
- [Create deviceManagementConfigurationPolicyAssignment](https://learn.microsoft.com/en-us/graph/api/intune-deviceconfigv2-devicemanagementconfigurationpolicy-post-assignments?view=graph-rest-beta)

## API Permissions

The following API permissions are required in order to use this resource.

### Microsoft Graph

- **Application**: `DeviceManagementConfiguration.ReadWrite.All`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.14.1-alpha | Experimental | Initial release |

## Example Usage

```terraform
# Example 1: Assign configuration policy to all devices
resource "graph_beta_device_management_settings_catalog_assignment" "all_devices_assignment" {
  settings_catalog_id = "12345678-1234-1234-1234-123456789012"

  target {
    target_type = "allDevices"
  }
}

# Example 2: Assign configuration policy to all licensed users
resource "graph_beta_device_management_settings_catalog_assignment" "all_users_assignment" {
  settings_catalog_id = "12345678-1234-1234-1234-123456789012"

  target {
    target_type = "allLicensedUsers"
  }
}

# Example 3: Assign configuration policy to a specific Entra ID group
resource "graph_beta_device_management_settings_catalog_assignment" "group_assignment" {
  settings_catalog_id = "12345678-1234-1234-1234-123456789012"

  target {
    target_type = "groupAssignment"
    group_id    = "87654321-4321-4321-4321-210987654321"
  }
}

# Example 4: Assign configuration policy with group exclusion
resource "graph_beta_device_management_settings_catalog_assignment" "exclusion_assignment" {
  settings_catalog_id = "12345678-1234-1234-1234-123456789012"

  target {
    target_type = "exclusionGroupAssignment"
    group_id    = "87654321-4321-4321-4321-210987654321"
  }
}

# Example 5: Assign configuration policy to SCCM collection
resource "graph_beta_device_management_settings_catalog_assignment" "sccm_assignment" {
  settings_catalog_id = "12345678-1234-1234-1234-123456789012"

  target {
    target_type   = "configurationManagerCollection"
    collection_id = "SMS00000001" # Default SMS collection or use custom like "MEM12345678"
  }
}

# Example 6: Group assignment with include filter
resource "graph_beta_device_management_settings_catalog_assignment" "group_with_include_filter" {
  settings_catalog_id = "12345678-1234-1234-1234-123456789012"

  target {
    target_type                                      = "groupAssignment"
    group_id                                         = "87654321-4321-4321-4321-210987654321"
    device_and_app_management_assignment_filter_id   = "11111111-2222-3333-4444-555555555555"
    device_and_app_management_assignment_filter_type = "include"
  }
}

# Example 7: Group assignment with exclude filter
resource "graph_beta_device_management_settings_catalog_assignment" "group_with_exclude_filter" {
  settings_catalog_id = "12345678-1234-1234-1234-123456789012"

  target {
    target_type                                      = "groupAssignment"
    group_id                                         = "87654321-4321-4321-4321-210987654321"
    device_and_app_management_assignment_filter_id   = "11111111-2222-3333-4444-555555555555"
    device_and_app_management_assignment_filter_type = "exclude"
  }
}

# Example 8: Assignment from policy sets
resource "graph_beta_device_management_settings_catalog_assignment" "policy_set_assignment" {
  settings_catalog_id = "12345678-1234-1234-1234-123456789012"
  source              = "policySets"
  source_id           = "99999999-8888-7777-6666-555555555555"

  target {
    target_type = "groupAssignment"
    group_id    = "87654321-4321-4321-4321-210987654321"
  }
}

# Example 9: Assignment with custom timeouts
resource "graph_beta_device_management_settings_catalog_assignment" "assignment_with_timeouts" {
  settings_catalog_id = "12345678-1234-1234-1234-123456789012"

  target {
    target_type = "allDevices"
  }

  timeouts {
    create = "5m"
    read   = "3m"
    update = "5m"
    delete = "3m"
  }
}

# Example 10: Multiple assignments for the same policy
resource "graph_beta_device_management_settings_catalog_assignment" "primary_group" {
  settings_catalog_id = "12345678-1234-1234-1234-123456789012"

  target {
    target_type = "groupAssignment"
    group_id    = "primary-group-id-1234-1234-1234-123456789012"
  }
}

resource "graph_beta_device_management_settings_catalog_assignment" "secondary_group" {
  settings_catalog_id = "12345678-1234-1234-1234-123456789012"

  target {
    target_type = "groupAssignment"
    group_id    = "secondary-group-id-5678-5678-5678-567856785678"
  }
}

resource "graph_beta_device_management_settings_catalog_assignment" "exclude_test_group" {
  settings_catalog_id = "12345678-1234-1234-1234-123456789012"

  target {
    target_type = "exclusionGroupAssignment"
    group_id    = "test-group-id-9999-9999-9999-999999999999"
  }
}

# Example 11: Using data sources for dynamic assignment
data "azuread_group" "it_department" {
  display_name = "IT Department"
}

resource "graph_beta_device_management_settings_catalog_assignment" "it_department_assignment" {
  settings_catalog_id = "12345678-1234-1234-1234-123456789012"

  target {
    target_type = "groupAssignment"
    group_id    = data.azuread_group.it_department.object_id
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `settings_catalog_id` (String) The ID of the settings catalog (configuration policy) to attach the assignment to.
- `target` (Attributes) (see [below for nested schema](#nestedatt--target))

### Optional

- `source` (String) Represents source of assignment. Possible values are: direct, policySets. Default is direct.
- `source_id` (String) The identifier of the source of the assignment. This property is read-only.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `id` (String) The ID of the device management configuration policy assignment.

<a id="nestedatt--target"></a>
### Nested Schema for `target`

Required:

- `target_type` (String) The target group type for the configuration policy assignment. Possible values are:

- **allDevices**: Target all devices in the tenant
- **allLicensedUsers**: Target all licensed users in the tenant
- **configurationManagerCollection**: Target System Center Configuration Manager collection
- **exclusionGroupAssignment**: Target a specific Entra ID group for exclusion
- **groupAssignment**: Target a specific Entra ID group

Optional:

- `collection_id` (String) The SCCM group collection ID for the assignment target. Default collections start with 'SMS', while custom collections start with your site code (e.g., 'MEM').
- `device_and_app_management_assignment_filter_id` (String) The Id of the scope filter applied to the target assignment.
- `device_and_app_management_assignment_filter_type` (String) The type of scope filter for the target assignment. Defaults to 'none'. Possible values are:

- **include**: Only include devices or users matching the filter
- **exclude**: Exclude devices or users matching the filter
- **none**: No assignment filter applied
- `group_id` (String) The Entra ID group ID for the assignment target. Required when target_type is 'groupAssignment' or 'exclusionGroupAssignment'.


<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Import

Import is supported using the following syntax:

```shell
# {resource_id}
terraform import microsoft365_graph_beta_device_and_app_management_settings_catalog_assignment.example settings-catalog-assignment-id
```

