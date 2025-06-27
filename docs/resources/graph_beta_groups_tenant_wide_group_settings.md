---
page_title: "microsoft365_graph_beta_groups_tenant_wide_group_settings Resource - terraform-provider-microsoft365"
subcategory: "Groups"
description: |-
    Manages tenant-wide directory settings for Microsoft 365 groups using the /settings endpoint. This resource enables a collection of configurations that allow admins to manage behaviors for specific Microsoft Entra objects like Microsoft 365 groups.This resource applies settings tenant-wide, enabling admins to control various aspects of group functionality.Use this resource in conjunction with the datasource 'microsoft365_graph_beta_directory_management_directory_setting_templates' to get the template_id, settings and values.
---

# microsoft365_graph_beta_groups_tenant_wide_group_settings (Resource)

Manages tenant-wide directory settings for Microsoft 365 groups using the `/settings` endpoint. This resource enables a collection of configurations that allow admins to manage behaviors for specific Microsoft Entra objects like Microsoft 365 groups.This resource applies settings tenant-wide, enabling admins to control various aspects of group functionality.Use this resource in conjunction with the datasource 'microsoft365_graph_beta_directory_management_directory_setting_templates' to get the template_id, settings and values.

## Microsoft Documentation

- [Directory setting resource type](https://learn.microsoft.com/en-us/graph/api/resources/directorysetting?view=graph-rest-beta)
- [List directory settings](https://learn.microsoft.com/en-us/graph/api/directory-list-settings?view=graph-rest-beta)
- [Create directory setting](https://learn.microsoft.com/en-us/graph/api/directory-post-settings?view=graph-rest-beta)
- [Update directory setting](https://learn.microsoft.com/en-us/graph/api/directorysetting-update?view=graph-rest-beta)
- [Delete directory setting](https://learn.microsoft.com/en-us/graph/api/directorysetting-delete?view=graph-rest-beta)

## API Permissions

The following API permissions are required in order to use this resource.

### Microsoft Graph

- **Application**: `Directory.ReadWrite.All`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.18.0-alpha | Experimental | Initial release |

## Example Usage

```terraform
# Example: Configure tenant-wide Microsoft 365 group settings
# This resource manages organization-wide policies for Microsoft 365 groups

# First, get the Group.Unified template details to see available settings
data "microsoft365_graph_beta_directory_management_directory_setting_templates" "group_unified" {
  filter_type  = "display_name"
  filter_value = "Group.Unified"
}

# Configure tenant-wide group creation and guest access policies using the Group.Unified template
resource "microsoft365_graph_beta_groups_tenant_wide_group_settings" "unified_settings" {
  # Use the template ID from the data source
  template_id = data.microsoft365_graph_beta_directory_management_directory_setting_templates.group_unified.directory_setting_templates[0].id

  # Example 1: Minimal configuration with only essential settings
  # This shows how to configure just the most commonly used settings
  values = [
    {
      name  = "EnableGroupCreation"
      value = "true" # Allow users to create Microsoft 365 groups
    },
    {
      name  = "AllowGuestsToAccessGroups"
      value = "true" # Allow guest access to groups
    },
    {
      name  = "AllowToAddGuests"
      value = "true" # Allow adding guests to groups
    }
  ]

  timeouts = {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}

# Example 2: Comprehensive configuration with naming policies and guest restrictions
resource "microsoft365_graph_beta_groups_tenant_wide_group_settings" "comprehensive_settings" {
  # Use the template ID from the data source
  template_id = data.microsoft365_graph_beta_directory_management_directory_setting_templates.group_unified.directory_setting_templates[0].id

  values = [
    # Group creation controls
    {
      name  = "EnableGroupCreation"
      value = "false" # Only specific groups can create Microsoft 365 groups
    },
    {
      name  = "GroupCreationAllowedGroupId"
      value = "00000000-0000-0000-0000-000000000000" # Replace with actual security group ID
    },

    # Naming policies
    {
      name  = "PrefixSuffixNamingRequirement"
      value = "[Marketing]-[GroupName]" # Enforce naming convention
    },
    {
      name  = "CustomBlockedWordsList"
      value = "CEO,Legal,HR" # Block specific words in group names
    },
    {
      name  = "EnableMSStandardBlockedWords"
      value = "true" # Enable Microsoft's list of blocked words
    },

    # Classification settings
    {
      name  = "ClassificationList"
      value = "Public,Internal,Confidential,Highly Confidential" # Available classifications
    },
    {
      name  = "ClassificationDescriptions"
      value = "Public:Public data,Internal:Internal data,Confidential:Confidential data,Highly Confidential:Highly Confidential data" # Descriptions
    },
    {
      name  = "DefaultClassification"
      value = "Internal" # Default classification
    },

    # Guest access settings
    {
      name  = "AllowGuestsToBeGroupOwner"
      value = "false" # Don't allow guests to be group owners
    },
    {
      name  = "AllowGuestsToAccessGroups"
      value = "true" # Allow guests to access groups
    },
    {
      name  = "AllowToAddGuests"
      value = "true" # Allow adding guests to groups
    },
    {
      name  = "GuestUsageGuidelinesUrl"
      value = "https://contoso.com/guestpolicies" # Link to guest usage guidelines
    },

    # Other settings
    {
      name  = "UsageGuidelinesUrl"
      value = "https://contoso.com/groupguidelines" # Link to general usage guidelines
    },
    {
      name  = "EnableMIPLabels"
      value = "true" # Enable sensitivity labels
    },
    {
      name  = "NewUnifiedGroupWritebackDefault"
      value = "true" # Enable group writeback to on-premises AD
    }
  ]

  timeouts = {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}

# Output the created settings ID
output "unified_settings_id" {
  value       = microsoft365_graph_beta_groups_tenant_wide_group_settings.unified_settings.id
  description = "The ID of the created tenant-wide group settings"
}

# NOTE: In a real environment, you would typically create only one tenant-wide setting per template.
# The two resources shown here are for demonstration purposes only.
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `template_id` (String) Unique identifier for the tenant-level directorySettingTemplate object that's been customized for this tenant-level settings object. The template options can be found at 'https://learn.microsoft.com/en-us/graph/group-directory-settings?tabs=http'.
- `values` (Attributes Set) Collection of name-value pairs corresponding to the name and defaultValue properties in the referenced directorySettingTemplate object. (see [below for nested schema](#nestedatt--values))

### Optional

- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `display_name` (String) Display name of this group of settings, which comes from the associated template. Read-only.
- `id` (String) The unique identifier for the tenant-wide setting. Read-only.

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

- **Tenant-wide Settings**: These settings apply to all groups in the tenant.
- **Template ID**: The `template_id` attribute determines the type of settings applied.
- **Values**: The `values` block allows specifying key-value pairs for settings.

## Import

Import is supported using the following syntax:

```shell
#!/bin/bash

# Import a tenant-wide directory setting
# Replace {setting_id} with the actual setting ID

terraform import microsoft365_graph_beta_groups_tenant_wide_group_settings.unified_settings "{setting_id}"
``` 