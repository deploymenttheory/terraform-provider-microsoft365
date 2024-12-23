---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "microsoft365_graph_beta_device_and_app_management_windows_platform_script Data Source - terraform-provider-microsoft365"
subcategory: ""
description: |-
  Retrieves information about a windows platform script.
---

# microsoft365_graph_beta_device_and_app_management_windows_platform_script (Data Source)

Retrieves information about a windows platform script.



<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `assignments` (Attributes) The assignment configuration for this Windows Settings Catalog profile. (see [below for nested schema](#nestedatt--assignments))
- `display_name` (String) Name of the windows platform script.
- `id` (String) Unique identifier for the windows platform script.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `description` (String) Description of the windows platform script.
- `enforce_signature_check` (Boolean) Indicate whether the script signature needs be checked.
- `file_name` (String) Script file name.
- `role_scope_tag_ids` (List of String) List of Scope Tag IDs for this PowerShellScript instance.
- `run_as_32_bit` (Boolean) A value indicating whether the PowerShell script should run as 32-bit.
- `run_as_account` (String) Indicates the type of execution context.
- `script_content` (String, Sensitive) The script content.

<a id="nestedatt--assignments"></a>
### Nested Schema for `assignments`

Optional:

- `all_devices` (Boolean) Specifies whether this assignment applies to all devices. When set to `true`, the assignment targets all devices in the organization.Can be used in conjuction with `all_users`.Can be used as an alternative to `include_groups`.Can be used in conjuction with `all_users` and `exclude_group_ids`.
- `all_users` (Boolean) Specifies whether this assignment applies to all users. When set to `true`, the assignment targets all licensed users within the organization.Can be used in conjuction with `all_devices`.Can be used as an alternative to `include_groups`.Can be used in conjuction with `all_devices` and `exclude_group_ids`.
- `exclude_group_ids` (Set of String) A set of group IDs to exclude from the assignment. These groups will not receive the assignment, even if they match other inclusion criteria.
- `include_group_ids` (Set of String) A set of entra id group Id's to include in the assignment.


<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
