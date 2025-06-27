---
page_title: "microsoft365_graph_beta_groups_group_settings Resource - terraform-provider-microsoft365"
subcategory: "Groups"
description: |-
    Manages group-specific directory settings for Microsoft 365 groups using the /groups/{group-id}/settings endpoint.This resource enables configuration of group-level settings such as guest access permissions and other group-specific policies that override tenant-wide defaults.Use this resource in conjunction with the datasource 'microsoft365_graph_beta_directory_management_directory_setting_templates' to get the template_id, settings and values.Use this resource in conjection with the resource 'microsoft365_graph_beta_groups_group' to get the group_id.
---

# microsoft365_graph_beta_groups_group_settings (Resource)

Manages group-specific directory settings for Microsoft 365 groups using the `/groups/{group-id}/settings` endpoint.This resource enables configuration of group-level settings such as guest access permissions and other group-specific policies that override tenant-wide defaults.Use this resource in conjunction with the datasource 'microsoft365_graph_beta_directory_management_directory_setting_templates' to get the template_id, settings and values.Use this resource in conjection with the resource 'microsoft365_graph_beta_groups_group' to get the group_id.

## Microsoft Documentation

- [Directory setting resource type](https://learn.microsoft.com/en-us/graph/api/resources/directorysetting?view=graph-rest-beta)
- [List group settings](https://learn.microsoft.com/en-us/graph/api/group-list-settings?view=graph-rest-beta)
- [Create group setting](https://learn.microsoft.com/en-us/graph/api/group-post-settings?view=graph-rest-beta)
- [Update group setting](https://learn.microsoft.com/en-us/graph/api/directorysetting-update?view=graph-rest-beta)
- [Delete group setting](https://learn.microsoft.com/en-us/graph/api/directorysetting-delete?view=graph-rest-beta)

## API Permissions

The following API permissions are required in order to use this resource.

### Microsoft Graph

- **Application**: `Group.ReadWrite.All`, `Directory.ReadWrite.All`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.18.0-alpha | Experimental | Initial release |

## Example Usage

```terraform
# Example: Configure group-specific settings for Microsoft 365 groups
# This resource manages directory settings for a specific Microsoft 365 group

# Get the Group.Unified.Guest template details to see available settings
data "microsoft365_graph_beta_directory_management_directory_setting_templates" "group_unified_guest" {
  filter_type  = "display_name"
  filter_value = "Group.Unified.Guest"
}

# Data source to get the group ID
data "microsoft365_graph_beta_groups_group" "example" {
  display_name = "Marketing Team"
}

# Example 1: Configure group-specific guest access settings
# This overrides tenant-wide guest settings for this specific group
resource "microsoft365_graph_beta_groups_group_settings" "guest_settings" {
  group_id    = data.microsoft365_graph_beta_groups_group.example.id
  template_id = data.microsoft365_graph_beta_directory_management_directory_setting_templates.group_unified_guest.directory_setting_templates[0].id

  values = [
    {
      name  = "AllowToAddGuests"
      value = "false" # Disable guest access for this specific group
    }
  ]

  timeouts = {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}

# Get the Group.Unified template details
data "microsoft365_graph_beta_directory_management_directory_setting_templates" "group_unified" {
  filter_type  = "display_name"
  filter_value = "Group.Unified"
}

# Example 2: Configure group-specific unified settings
# This shows how to override other group settings for a specific group
resource "microsoft365_graph_beta_groups_group_settings" "unified_settings" {
  group_id    = data.microsoft365_graph_beta_groups_group.example.id
  template_id = data.microsoft365_graph_beta_directory_management_directory_setting_templates.group_unified.directory_setting_templates[0].id

  values = [
    {
      name  = "ClassificationList"
      value = "Confidential,Secret,Top Secret" # Custom classifications for this group
    },
    {
      name  = "DefaultClassification"
      value = "Confidential" # Default classification for this group
    },
    {
      name  = "UsageGuidelinesUrl"
      value = "https://contoso.com/marketing-group-guidelines" # Group-specific guidelines
    }
  ]

  timeouts = {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}

# Output the created settings IDs
output "guest_settings_id" {
  value       = microsoft365_graph_beta_groups_group_settings.guest_settings.id
  description = "The ID of the created group-specific guest settings"
}

output "unified_settings_id" {
  value       = microsoft365_graph_beta_groups_group_settings.unified_settings.id
  description = "The ID of the created group-specific unified settings"
}

# NOTE: In a real environment, you would typically create only one setting per template per group.
# The two resources shown here are for demonstration purposes only.
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `group_id` (String) The unique identifier of the group for which the settings apply.
- `template_id` (String) Unique identifier for the tenant-level directorySettingTemplate object that's been customized for this group-level settings object. The template named 'Group.Unified.Guest' can be used to configure group-specific settings.
- `values` (Attributes Set) Collection of name-value pairs corresponding to the name and defaultValue properties in the referenced directorySettingTemplate object. (see [below for nested schema](#nestedatt--values))

### Optional

- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `display_name` (String) Display name of this group of settings, which comes from the associated template. Read-only.
- `id` (String) The unique identifier for the group setting. Read-only.

<a id="nestedatt--values"></a>
### Nested Schema for `values`

Required:

- `name` (String) Name of the setting from the referenced directorySettingTemplate.
- `value` (String) Value of the setting.


<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Important Notes

- **Group Settings**: These settings apply to a specific group.
- **Template ID**: The `template_id` attribute determines the type of settings applied.
- **Values**: The `values` block allows specifying key-value pairs for settings.

## Import

Import is supported using the following syntax:

```shell
#!/bin/bash

# Import a group-specific directory setting
# Format: {group_id}/{setting_id}
# Replace {group_id} with the actual group ID and {setting_id} with the setting ID

terraform import microsoft365_graph_beta_groups_group_settings.guest_settings "{group_id}/{setting_id}"
``` 