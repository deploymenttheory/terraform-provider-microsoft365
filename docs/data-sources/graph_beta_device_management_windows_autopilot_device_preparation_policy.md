---
page_title: "microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy Data Source - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Retrieves Windows Autopilot device preparation policies from Microsoft Intune using the /deviceManagement/configurationPolicies endpoint. Device preparation policies are settings catalog policies created from the Autopilot device preparation templates, so results are always restricted to policies whose templateReference matches one of those templates. Supports flexible lookup by policy ID, name, custom OData queries, or listing all policies.
---

# microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy (Data Source)

Retrieves Windows Autopilot device preparation policies from Microsoft Intune using the `/deviceManagement/configurationPolicies` endpoint. Device preparation policies are settings catalog policies created from the Autopilot device preparation templates, so results are always restricted to policies whose `templateReference` matches one of those templates. Supports flexible lookup by policy ID, name, custom OData queries, or listing all policies.

## Microsoft Documentation

- [deviceManagementConfigurationPolicy resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-devicemanagementconfigurationpolicy?view=graph-rest-beta)

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this data source:

**Required:**
- `DeviceManagementConfiguration.Read.All`

**Optional:**
- `None` `[N/A]`

## Example Usage

```terraform
# Look up a Windows Autopilot device preparation policy by its Intune object id
data "microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy" "by_id" {
  policy_id = "00000000-0000-0000-0000-000000000000"
}

# Look up a Windows Autopilot device preparation policy by its name
data "microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy" "by_name" {
  name = "Autopilot Device Preparation - Cloud PC"
}

# List every Windows Autopilot device preparation policy in the tenant
data "microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy" "all" {
  list_all = true
}

# Use a custom OData query for advanced filtering. Results are always restricted to
# Autopilot device preparation policies.
data "microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy" "assigned" {
  odata_query = "isAssigned eq true"
}

# Look up a policy and also retrieve its assignments
data "microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy" "with_assignments" {
  name             = "Autopilot Device Preparation - Cloud PC"
  list_assignments = true
}

# Use the looked up policy id when provisioning Windows 365 Cloud PCs. This allows the
# provisioning policy to be managed in a different Terraform state than the device
# preparation policy itself.
resource "microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy" "example" {
  display_name             = "Cloud PC - Autopilot device preparation"
  description              = "Provisioning policy using an existing device preparation policy"
  image_id                 = "microsoftwindowsdesktop_windows-ent-cpc_win11-24H2-ent-cpc-m365"
  image_type               = "gallery"
  provisioning_type        = "dedicated"
  cloud_pc_naming_template = "CPC-%USERNAME:4%-%RAND:5%"

  domain_join_configurations = [
    {
      domain_join_type = "azureADJoin"
      region_group     = "europeUnion"
      region_name      = "automatic"
    }
  ]

  autopilot_configuration = {
    device_preparation_profile_id = data.microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.by_name.items[0].id
  }
}

output "device_preparation_policy_id" {
  value       = data.microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.by_name.items[0].id
  description = "ID of the existing Windows Autopilot device preparation policy"
}

output "device_preparation_policy_deployment_mode" {
  value       = data.microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.by_id.items[0].template_reference.deployment_mode
  description = "Deployment mode of the existing Windows Autopilot device preparation policy"
}

output "all_device_preparation_policy_names" {
  value       = [for policy in data.microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.all.items : policy.name]
  description = "Names of every Windows Autopilot device preparation policy in the tenant"
}

output "device_preparation_policy_assigned_groups" {
  value       = [for assignment in data.microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.with_assignments.assignments : assignment.group_id]
  description = "Groups the Windows Autopilot device preparation policy is assigned to"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `list_all` (Boolean) Retrieve all Windows Autopilot device preparation policies in the tenant. Conflicts with specific lookup attributes.
- `list_assignments` (Boolean) When true, retrieves the assignments of the matched policy. When several policies match, the assignments of the first policy are returned.
- `name` (String) The name of the device preparation policy. The lookup is an exact match on the policy name. Conflicts with other lookup attributes.
- `odata_query` (String) Custom OData filter expression for advanced queries (e.g., `name eq 'Autopilot policy' and isAssigned eq true`). Results are still restricted to device preparation policies. Conflicts with specific lookup attributes.
- `policy_id` (String) The unique identifier of the device preparation policy. Conflicts with other lookup attributes.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `assignments` (Attributes List) Assignments of the matched device preparation policy. Only populated when `list_assignments` is true. (see [below for nested schema](#nestedatt--assignments))
- `id` (String) The unique identifier for the data source. This is a placeholder attribute required by Terraform.
- `items` (Attributes List) List of Windows Autopilot device preparation policies matching the query criteria. (see [below for nested schema](#nestedatt--items))

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).


<a id="nestedatt--assignments"></a>
### Nested Schema for `assignments`

Read-Only:

- `filter_id` (String) The assignment filter applied to the assignment, if any.
- `filter_type` (String) The assignment filter mode. One of `include`, `exclude` or `none`.
- `group_id` (String) The group the policy is assigned to, when the target is a group assignment target.
- `id` (String) The unique identifier of the assignment.
- `type` (String) The assignment target type, for example `groupAssignmentTarget`.


<a id="nestedatt--items"></a>
### Nested Schema for `items`

Read-Only:

- `created_date_time` (String) The date and time the policy was created.
- `creation_source` (String) The source that created the policy.
- `description` (String) The description of the device preparation policy.
- `disable_entra_group_policy_assignment` (Boolean) Indicates whether Microsoft Entra group policy assignment is disabled for the policy.
- `id` (String) The unique identifier of the device preparation policy. This is the value required by `autopilot_configuration.device_preparation_profile_id` on a Windows 365 Cloud PC provisioning policy.
- `is_assigned` (Boolean) Indicates whether the policy is assigned to any group.
- `last_modified_date_time` (String) The date and time the policy was last modified.
- `name` (String) The name of the device preparation policy.
- `platforms` (String) The platforms the policy applies to.
- `priority` (Number) The priority of the policy, taken from its priority metadata.
- `role_scope_tag_ids` (List of String) The list of scope tag IDs applied to the policy.
- `setting_count` (Number) The number of settings configured on the policy.
- `technologies` (String) The technologies used to deliver the policy.
- `template_reference` (Attributes) The settings catalog template the policy was created from. (see [below for nested schema](#nestedatt--items--template_reference))

<a id="nestedatt--items--template_reference"></a>
### Nested Schema for `items.template_reference`

Read-Only:

- `deployment_mode` (String) The Autopilot device preparation deployment mode derived from the template ID. Either `automatic` or `user_driven`.
- `template_display_name` (String) The display name of the template.
- `template_display_version` (String) The display version of the template.
- `template_family` (String) The template family. Always `enrollmentConfiguration` for device preparation policies.
- `template_id` (String) The identifier of the template, including its version suffix.
