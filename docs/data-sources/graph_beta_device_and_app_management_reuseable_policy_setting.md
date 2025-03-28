---
page_title: "microsoft365_graph_beta_device_and_app_management_reuseable_policy_setting Data Source - terraform-provider-microsoft365"
subcategory: "Intune"
description: |-
  Manages a Reuseable Settings Policy using Settings Catalog in Microsoft Intune for Endpoint Privilege Management.Endpoint Privilege Management supports using reusable settings groups to manage the certificates in place of adding that certificatedirectly to an elevation rule. Like all reusable settings groups for Intune, configurations and changes made to a reusable settingsgroup are automatically passed to the policies that reference the group.
---

# microsoft365_graph_beta_device_and_app_management_reuseable_policy_setting (Data Source)

The Microsoft 365 Intune role scope tag data source provides information about a specific scope tag.

## Example Usage

```terraform
// Data Source: Reusable Policy Settings
// Basic usage: lookup by display name
data "microsoft365_graph_beta_device_and_app_management_reuseable_policy_setting" "example" {
  display_name = "epm certificate"
}

// Output to verify data source
output "reuseable_policy_settings_details" {
  value = {
    id                                     = data.microsoft365_graph_beta_device_and_app_management_reuseable_policy_settings.example.id
    display_name                           = data.microsoft365_graph_beta_device_and_app_management_reuseable_policy_settings.example.display_name
    description                            = data.microsoft365_graph_beta_device_and_app_management_reuseable_policy_settings.example.description
    settings                               = data.microsoft365_graph_beta_device_and_app_management_reuseable_policy_settings.example.settings
    created_date_time                      = data.microsoft365_graph_beta_device_and_app_management_reuseable_policy_settings.example.created_date_time
    last_modified_date_time                = data.microsoft365_graph_beta_device_and_app_management_reuseable_policy_settings.example.last_modified_date_time
    version                                = data.microsoft365_graph_beta_device_and_app_management_reuseable_policy_settings.example.version
    referencing_configuration_policies     = data.microsoft365_graph_beta_device_and_app_management_reuseable_policy_settings.example.referencing_configuration_policies
    referencing_configuration_policy_count = data.microsoft365_graph_beta_device_and_app_management_reuseable_policy_settings.example.referencing_configuration_policy_count
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `display_name` (String) The reusable setting display name supplied by user.
- `id` (String) The unique identifier for this Reuseable Settings Policy
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `created_date_time` (String) Creation date and time of the settings catalog policy
- `description` (String) Reuseable Settings Policy description
- `last_modified_date_time` (String) Last modification date and time of the settings catalog policy
- `referencing_configuration_policies` (List of String) List of configuration policies referencing this reuseable policy
- `referencing_configuration_policy_count` (Number) Number of configuration policies referencing this reuseable policy
- `settings` (String) Reuseable Settings Policy with settings catalog settings defined as a valid JSON string.
- `version` (Number) Version of the policy

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).