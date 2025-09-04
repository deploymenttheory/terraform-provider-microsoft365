---
page_title: "microsoft365_graph_beta_device_management_autopatch_groups Resource - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Manages Windows Autopatch groups using the https://services.autopatch.microsoft.com/device/v2/autopatchGroups endpoint. Autopatch groups help organize devices into logical groups for automated Windows Update deployment with customizable deployment rings and policy settings.This resource is not documented in the Microsoft Graph API documentation. This resource is experimental and may change in the future.There's currently 401 errors when trying using this resource. There appears to be a seperate unobservable authentication step between intune and  autopatch micro service that cannot be replicated in the terraform provider. Entra ID client id / secret are not sufficient to authenticate.
---

# microsoft365_graph_beta_device_management_autopatch_groups (Resource)

Manages Windows Autopatch groups using the `https://services.autopatch.microsoft.com/device/v2/autopatchGroups` endpoint. Autopatch groups help organize devices into logical groups for automated Windows Update deployment with customizable deployment rings and policy settings.This resource is not documented in the Microsoft Graph API documentation. This resource is experimental and may change in the future.There's currently 401 errors when trying using this resource. There appears to be a seperate unobservable authentication step between intune and  autopatch micro service that cannot be replicated in the terraform provider. Entra ID client id / secret are not sufficient to authenticate.

## Undocumented

This resource is not documented in the Microsoft Graph API documentation. This resource is experimental and may change in the future. There's currently 401 errors when using this resource.
This occurs for both client id and client secret authentication.

APIEndpoint:  "https://services.autopatch.microsoft.com"
ResourcePath: "/device/v2/autopatchGroups"

## API Permissions

The following API permissions are required in order to use this resource.

### Microsoft Graph

- **Application**: `TBC`

## Example Usage

```terraform
resource "microsoft365_graph_beta_device_management_autopatch_groups" "auto_patch_group" {
  name        = "auto-patch-group"
  description = ""

  global_user_managed_aad_groups = []

  # Deployment Groups
  deployment_groups = [
    {
      name = "auto-patch-group - Test"
      user_managed_aad_groups = [
        {
          id   = "00000000-0000-0000-0000-000000000000"
          name = "group-name-01"
        }
      ]
      deployment_group_policy_settings = {
        aad_group_name              = "auto-patch-group - Test"
        is_update_settings_modified = false
        device_configuration_setting = {
          policy_id            = "000"
          update_behavior      = "AutoInstallAndRestart"
          notification_setting = "DefaultNotifications"
          quality_deployment_settings = {
            deadline     = 1
            deferral     = 0
            grace_period = 0
          }
          feature_deployment_settings = {
            deadline = 5
            deferral = 0
          }
        }
      }
    },
    {
      name = "auto-patch-group - Ring1"
      user_managed_aad_groups = [
        {
          id   = "00000000-0000-0000-0000-000000000000"
          name = "group-name-02"
          type = 0
        }
      ]
      deployment_group_policy_settings = {
        aad_group_name              = "auto-patch-group - Ring1"
        is_update_settings_modified = false
        device_configuration_setting = {
          policy_id            = "000"
          update_behavior      = "AutoInstallAndRestart"
          notification_setting = "DefaultNotifications"
          quality_deployment_settings = {
            deadline     = 2
            deferral     = 1
            grace_period = 2
          }
          feature_deployment_settings = {
            deadline = 5
            deferral = 0
          }
        }
      }
    },
    {
      name = "auto-patch-group - Last"
      user_managed_aad_groups = [
        {
          id   = "00000000-0000-0000-0000-000000000000"
          name = "group-name-03"
        }
      ]
      deployment_group_policy_settings = {
        aad_group_name              = "auto-patch-group - Last"
        is_update_settings_modified = false
        device_configuration_setting = {
          policy_id            = "000"
          update_behavior      = "AutoInstallAndRestart"
          notification_setting = "DefaultNotifications"
          quality_deployment_settings = {
            deadline     = 3
            deferral     = 5
            grace_period = 2
          }
          feature_deployment_settings = {
            deadline = 5
            deferral = 0
          }
        }
      }
    }
  ]

  # Driver update settings
  enable_driver_update = true

  # Scope tags
  scope_tags = [0]

  # Enabled content types
  enabled_content_types = 31

  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the Autopatch group

### Optional

- `deployment_groups` (Attributes Set) The deployment groups (rings) within this Autopatch group (see [below for nested schema](#nestedatt--deployment_groups))
- `description` (String) The description of the Autopatch group
- `enable_driver_update` (Boolean) Whether driver updates are enabled
- `enabled_content_types` (Number) Enabled content types bitmask
- `global_user_managed_aad_groups` (Attributes Set) Global user-managed Azure AD groups (see [below for nested schema](#nestedatt--global_user_managed_aad_groups))
- `scope_tags` (Set of String) Set of scope tag IDs for this Autopatch group.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `distribution_type` (String) The distribution type (Mixed, AdminAssigned)
- `flow_id` (String) The flow ID for the operation
- `flow_status` (String) The flow status for the operation
- `flow_type` (String) The flow type for the operation
- `id` (String) The unique identifier for this Autopatch group
- `is_locked_by_policy` (Boolean) Whether the group is locked by policy
- `number_of_registered_devices` (Number) The number of registered devices in the group
- `read_only` (Boolean) Whether the group is read-only
- `status` (String) The status of the Autopatch group (Active, Creating, etc.)
- `tenant_id` (String) The tenant ID associated with this Autopatch group
- `type` (String) The type of the Autopatch group (Default, User)
- `umbrella_group_id` (String) The umbrella group ID
- `user_has_all_scope_tag` (Boolean) Whether the user has all scope tags

<a id="nestedatt--deployment_groups"></a>
### Nested Schema for `deployment_groups`

Required:

- `name` (String) The name of the deployment group

Optional:

- `aad_id` (String) The Azure AD group ID for this deployment group
- `deployment_group_policy_settings` (Attributes) Policy settings for this deployment group (see [below for nested schema](#nestedatt--deployment_groups--deployment_group_policy_settings))
- `distribution` (Number) Distribution percentage for this deployment group
- `user_managed_aad_groups` (Attributes Set) User-managed Azure AD groups for this deployment group (see [below for nested schema](#nestedatt--deployment_groups--user_managed_aad_groups))

Read-Only:

- `failed_prerequisite_check_count` (Number) Number of failed prerequisite checks

<a id="nestedatt--deployment_groups--deployment_group_policy_settings"></a>
### Nested Schema for `deployment_groups.deployment_group_policy_settings`

Optional:

- `aad_group_name` (String) The Azure AD group name
- `device_configuration_setting` (Attributes) Device configuration settings (see [below for nested schema](#nestedatt--deployment_groups--deployment_group_policy_settings--device_configuration_setting))
- `is_update_settings_modified` (Boolean) Whether update settings are modified

<a id="nestedatt--deployment_groups--deployment_group_policy_settings--device_configuration_setting"></a>
### Nested Schema for `deployment_groups.deployment_group_policy_settings.device_configuration_setting`

Optional:

- `feature_deployment_settings` (Attributes) Feature update deployment settings (see [below for nested schema](#nestedatt--deployment_groups--deployment_group_policy_settings--device_configuration_setting--feature_deployment_settings))
- `notification_setting` (String) Notification setting
- `policy_id` (String) The policy ID
- `quality_deployment_settings` (Attributes) Quality update deployment settings (see [below for nested schema](#nestedatt--deployment_groups--deployment_group_policy_settings--device_configuration_setting--quality_deployment_settings))
- `update_behavior` (String) Update behavior setting

<a id="nestedatt--deployment_groups--deployment_group_policy_settings--device_configuration_setting--feature_deployment_settings"></a>
### Nested Schema for `deployment_groups.deployment_group_policy_settings.device_configuration_setting.feature_deployment_settings`

Optional:

- `deadline` (Number) Deadline in days
- `deferral` (Number) Deferral in days


<a id="nestedatt--deployment_groups--deployment_group_policy_settings--device_configuration_setting--quality_deployment_settings"></a>
### Nested Schema for `deployment_groups.deployment_group_policy_settings.device_configuration_setting.quality_deployment_settings`

Optional:

- `deadline` (Number) Deadline in days
- `deferral` (Number) Deferral in days
- `grace_period` (Number) Grace period in days




<a id="nestedatt--deployment_groups--user_managed_aad_groups"></a>
### Nested Schema for `deployment_groups.user_managed_aad_groups`

Required:

- `id` (String) The ID of the Azure AD group

Optional:

- `name` (String) The name of the Azure AD group
- `type` (Number) The type of the group



<a id="nestedatt--global_user_managed_aad_groups"></a>
### Nested Schema for `global_user_managed_aad_groups`

Required:

- `id` (String) The ID of the Azure AD group
- `type` (String) The type of the group (Device, User)


<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Important Notes

- **Autopatch Groups**: This resource creates Autopatch groups that can be used to manage Windows Update deployment with customizable deployment rings and policy settings.
- **Deployment Rings**: Autopatch groups can be used to create deployment rings that can be used to manage Windows Update deployment with customizable deployment rings and policy settings.
- **Policy Settings**: Autopatch groups can be used to manage Windows Update deployment with customizable deployment rings and policy settings.


## Import

Import is supported using the following syntax:

```shell
# {resource_id}
terraform import microsoft365_graph_beta_device_management_autopatch_groups.example 00000000-0000-0000-0000-000000000000
```