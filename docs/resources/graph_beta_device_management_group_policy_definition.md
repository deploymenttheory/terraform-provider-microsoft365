---
page_title: "microsoft365_graph_beta_device_management_group_policy_definition Resource - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Manages a group policy definition with all its presentation values in Microsoft Intune. This resource provides a unified interface for configuring any group policy regardless of presentation types (checkboxes, textboxes, dropdowns, etc.). Values are provided as label-value pairs, and the resource automatically handles type conversion based on the policy's catalog definition. Uses the deviceManagement/groupPolicyConfigurations('{groupPolicyConfigurationId}')/updateDefinitionValues endpoint.
---

# microsoft365_graph_beta_device_management_group_policy_definition (Resource)

Manages a group policy definition with all its presentation values in Microsoft Intune. This resource provides a unified interface for configuring any group policy regardless of presentation types (checkboxes, textboxes, dropdowns, etc.). Values are provided as label-value pairs, and the resource automatically handles type conversion based on the policy's catalog definition. Uses the `deviceManagement/groupPolicyConfigurations('{groupPolicyConfigurationId}')/updateDefinitionValues` endpoint.

## Microsoft Documentation

- [Group policy definition value resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-grouppolicy-grouppolicydefinitionvalue?view=graph-rest-beta)
- [Update definition values](https://learn.microsoft.com/en-us/graph/api/intune-grouppolicy-grouppolicyconfiguration-updatedefinitionvalues?view=graph-rest-beta)
- [Group policy definition resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-grouppolicy-grouppolicydefinition?view=graph-rest-beta)
- [Group policy presentation value types](https://learn.microsoft.com/en-us/graph/api/resources/intune-grouppolicy-grouppolicypresentationvalue?view=graph-rest-beta)

## API Permissions

The following API permissions are required in order to use this resource.

### Microsoft Graph

- **Application**: `DeviceManagementConfiguration.Read.All`, `DeviceManagementConfiguration.ReadWrite.All`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.40.0-alpha | Experimental | Initial release |

## Example Usage

This resource supports all Group Policy presentation types through a unified interface. The resource automatically handles type conversion based on the presentation type defined in Microsoft's Group Policy catalog.

### Boolean (CheckBox) Presentation

```terraform
# Example: Group Policy Definition with Boolean (CheckBox) presentation values
# This example demonstrates configuring a policy with multiple boolean checkboxes

resource "microsoft365_graph_beta_device_management_group_policy_configuration" "example" {
  display_name = "Example Group Policy Configuration"
  description  = "Configuration for managing Windows Store packages"
}

resource "microsoft365_graph_beta_device_management_group_policy_definition" "boolean_example" {
  group_policy_configuration_id = microsoft365_graph_beta_device_management_group_policy_configuration.example.id
  policy_name                   = "Remove Default Microsoft Store packages from the system."
  class_type                    = "machine"
  category_path                 = "\\Windows Components\\App Package Deployment"
  enabled                       = true

  values = [
    {
      label = "Microsoft Teams"
      value = "true"
    },
    {
      label = "Paint"
      value = "false"
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
```

### TextBox Presentation

```terraform
# Example: Group Policy Definition with TextBox presentation value
# This example demonstrates configuring a policy with a text input field

resource "microsoft365_graph_beta_device_management_group_policy_configuration" "example" {
  display_name = "Example Group Policy Configuration"
  description  = "Configuration for Microsoft Edge browsing data lifetime"
}

resource "microsoft365_graph_beta_device_management_group_policy_definition" "textbox_example" {
  group_policy_configuration_id = microsoft365_graph_beta_device_management_group_policy_configuration.example.id
  policy_name                   = "Browsing Data Lifetime Settings"
  class_type                    = "machine"
  category_path                 = "\\Microsoft Edge"
  enabled                       = true

  values = [
    {
      label = "Browsing Data Lifetime Settings"
      value = "[{\"data_types\":[\"browsing_history\"],\"time_to_live_in_hours\":168}]"
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
```

### Decimal (DecimalTextBox) Presentation

```terraform
# Example: Group Policy Definition with Decimal (DecimalTextBox) presentation value
# This example demonstrates configuring a policy with a numeric input field

resource "microsoft365_graph_beta_device_management_group_policy_configuration" "example" {
  display_name = "Example Group Policy Configuration"
  description  = "Configuration for Microsoft Defender Antivirus settings"
}

resource "microsoft365_graph_beta_device_management_group_policy_definition" "decimal_example" {
  group_policy_configuration_id = microsoft365_graph_beta_device_management_group_policy_configuration.example.id
  policy_name                   = "Configure time out for detections in non-critical failed state"
  class_type                    = "machine"
  category_path                 = "\\Windows Components\\Microsoft Defender Antivirus\\Reporting"
  enabled                       = true

  values = [
    {
      label = "Configure time out for detections in non-critical failed state"
      value = "7200"
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
```

### MultiText (MultiTextBox) Presentation

```terraform
# Example: Group Policy Definition with MultiText (MultiTextBox) presentation value
# This example demonstrates configuring a policy with a multi-line text input field

resource "microsoft365_graph_beta_device_management_group_policy_configuration" "example" {
  display_name = "Example Group Policy Configuration"
  description  = "Configuration for filesystem filter drivers"
}

resource "microsoft365_graph_beta_device_management_group_policy_definition" "multitext_example" {
  group_policy_configuration_id = microsoft365_graph_beta_device_management_group_policy_configuration.example.id
  policy_name                   = "Dev drive filter attach policy"
  class_type                    = "machine"
  category_path                 = "\\System\\Filesystem"
  enabled                       = true

  values = [
    {
      label = "Filter list"
      value = "FilterDriver1\nFilterDriver2\nFilterDriver3"
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
```

### Dropdown (DropdownList) Presentation

```terraform
# Example: Group Policy Definition with Dropdown (DropdownList) presentation value
# This example demonstrates configuring a policy with a dropdown selection field

resource "microsoft365_graph_beta_device_management_group_policy_configuration" "example" {
  display_name = "Example Group Policy Configuration"
  description  = "Configuration for Internet Explorer security settings"
}

resource "microsoft365_graph_beta_device_management_group_policy_definition" "dropdown_example" {
  group_policy_configuration_id = microsoft365_graph_beta_device_management_group_policy_configuration.example.id
  policy_name                   = "Navigate windows and frames across different domains"
  class_type                    = "machine"
  category_path                 = "\\Windows Components\\Internet Explorer\\Internet Control Panel\\Security Page\\Internet Zone"
  enabled                       = true

  values = [
    {
      label = "Navigate windows and frames across different domains"
      value = "1"
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `category_path` (String) The category path of the group policy definition (e.g., '\Windows Components\App Package Deployment'). Used to identify the policy in the catalog
- `class_type` (String) The class type of the group policy definition. Must be 'user' or 'machine'
- `enabled` (Boolean) Whether this group policy definition is enabled (true) or disabled (false)
- `group_policy_configuration_id` (String) The unique identifier of the group policy configuration that contains this definition
- `policy_name` (String) The display name of the group policy definition (e.g., 'Remove Default Microsoft Store packages from the system.')

### Optional

- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `values` (Attributes Set) Set of presentation values for this group policy definition. Each value corresponds to a specific presentation (checkbox, textbox, dropdown, etc.) identified by its label. The resource automatically handles type conversion based on the presentation type in the catalog. (see [below for nested schema](#nestedatt--values))

### Read-Only

- `created_date_time` (String) The date and time when the definition value was created
- `id` (String) The unique identifier for the group policy definition value
- `last_modified_date_time` (String) The date and time when the definition value was last modified

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).


<a id="nestedatt--values"></a>
### Nested Schema for `values`

Required:

- `label` (String) The human-readable label of the presentation (e.g., 'Xbox Gaming App', 'Feedback Hub'). Must match a label from the policy's catalog definition
- `value` (String) The value for this presentation as a string. Format depends on presentation type: 'true'/'false' for checkboxes, numeric strings for decimal fields, plain text for textboxes, etc. The resource validates and converts this based on the presentation type

Read-Only:

- `id` (String) The unique identifier of the presentation template (computed from catalog)

## Important Notes

### Presentation Types and Values

This resource automatically handles different presentation types:

- **Boolean (CheckBox)**: Values must be `"true"` or `"false"` (as strings)
- **TextBox**: Any string value
- **Decimal**: Numeric values (as strings)
- **MultiText**: Multi-line text with `\n` for line breaks
- **Dropdown**: Numeric or string values corresponding to dropdown options

The resource uses the `policy_name`, `class_type`, and `category_path` fields to resolve the policy from Microsoft's Group Policy catalog. These fields must match exactly as they appear in the catalog.

### Resource Discovery

To find available policies and their metadata:

1. Use the PowerShell discovery script: `Get-GroupPolicyTemplateAndPresentationValues.ps1`
2. Query the catalog: `GET /deviceManagement/groupPolicyDefinitions`
3. Use the export script to generate HCL from existing configurations: `Export-GroupPolicyDefinitionToHCLForImport.ps1`

### Lifecycle Behavior

- **RequiresReplace**: Changes to `policy_name`, `class_type`, or `category_path` will trigger a destroy and recreate operation, as these identify a completely different policy definition.
- **In-place Update**: Changes to `enabled` or `values` are updated in place.

### Import Support

This resource supports importing using a composite ID format that enables proper multi-endpoint resolution:

- **Import ID Format**: `configurationID/definitionValueID`
- The import process automatically fetches `policy_name`, `class_type`, and `category_path` from the Microsoft Graph API
- Both simple (`configurationID`) and composite formats are supported during import

For bulk imports, use the export script to generate import commands automatically.

## Import

Import is supported using the following syntax:

```shell
#!/bin/bash

# Import a Group Policy Definition resource using the composite ID format
# The ID format is: configurationID/definitionValueID
#
# To find the composite ID:
# 1. Navigate to the Intune portal
# 2. Go to Devices > Group Policy Configurations
# 3. Select your configuration and view the definition values
# 4. Use the configuration GUID and definition value GUID

terraform import microsoft365_graph_beta_device_management_group_policy_definition.example \
  "a1b2c3d4-e5f6-g7h8-i9j0-k1l2m3n4o5p6/x1y2z3a4-b5c6-d7e8-f9g0-h1i2j3k4l5m6"

# Alternative: Use the PowerShell export script to generate import commands automatically
# pwsh Export-GroupPolicyDefinitionToHCLForImport.ps1 -TenantId "<tenant-id>" -ClientId "<client-id>" -ClientSecret "<secret>"
```

