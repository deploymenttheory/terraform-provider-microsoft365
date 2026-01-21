---
page_title: "microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile Resource - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Manages Windows Autopilot deployment profiles using the /deviceManagement/windowsAutopilotDeploymentProfiles endpoint. Autopilot deployment profiles define the out-of-box experience (OOBE) settings, device naming templates, and enrollment configurations for automated Windows device provisioning and domain joining.
---

# microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile (Resource)

Manages Windows Autopilot deployment profiles using the `/deviceManagement/windowsAutopilotDeploymentProfiles` endpoint. Autopilot deployment profiles define the out-of-box experience (OOBE) settings, device naming templates, and enrollment configurations for automated Windows device provisioning and domain joining.

## Microsoft Documentation

- [windowsAutopilotDeploymentProfile resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-enrollment-windowsautopilotdeploymentprofile?view=graph-rest-beta)

## API Permissions

The following API permissions are required in order to use this resource.

### Microsoft Graph

- **Application**: `DeviceManagementServiceConfig.ReadWrite.All`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.42.0-alpha | Experimental | Added missing version history |

## Example Usage

```terraform
# User-Driven Deployment Profile Example
resource "microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile" "user_driven" {
  display_name                                 = "user driven autopilot"
  description                                  = "user driven autopilot profile with os default locale"
  device_name_template                         = "thing-%RAND:5%" // Apply device name template max 15 characters
  locale                                       = "os-default"
  preprovisioning_allowed                      = true // Allow pre-provisioned deployment
  device_type                                  = "windowsPc"
  hardware_hash_extraction_enabled             = true
  role_scope_tag_ids                           = ["0", "1"]
  device_join_type                             = "microsoft_entra_joined"
  hybrid_azure_ad_join_skip_connectivity_check = false // always false when using microsoft_entra_joined

  out_of_box_experience_setting = {
    device_usage_type               = "singleUser"
    privacy_settings_hidden         = true       // Privacy settings
    eula_hidden                     = true       // Microsoft Software License Terms
    user_type                       = "standard" // standard or administrator
    keyboard_selection_page_skipped = true       // Automatically configure keyboard
  }

  // Optional assignments, can be either group based or all devices based
  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000001"
    },
    {
      type     = "groupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000002"
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000003"
    }
  ]
}

# User-Driven with Japanese Language and Allow Pre-provisioned Deployment
resource "microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile" "user_driven_japanese_preprovisioned_with_assignments" {
  display_name                                 = "acc_test_user_driven_japanese_preprovisioned"
  description                                  = "user driven autopilot profile with japanese locale and allow pre provisioned deployment"
  locale                                       = "ja-JP"
  preprovisioning_allowed                      = true
  device_type                                  = "windowsPc"
  hardware_hash_extraction_enabled             = true
  role_scope_tag_ids                           = ["0"]
  device_join_type                             = "microsoft_entra_hybrid_joined"
  hybrid_azure_ad_join_skip_connectivity_check = true

  out_of_box_experience_setting = {
    device_usage_type               = "singleUser"
    privacy_settings_hidden         = true
    eula_hidden                     = true
    user_type                       = "standard"
    keyboard_selection_page_skipped = true
  }

  // Optional assignments, can be either group based or all devices based
  assignments = [
    {
      type = "allDevicesAssignmentTarget"
    }
  ]
}

# Self-Deploying Deployment Profile Example
resource "microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile" "self_deploying" {
  display_name                                 = "self deploying autopilot"
  description                                  = "self deploying autopilot profile with os default locale"
  device_name_template                         = "thing-%RAND:2%"
  locale                                       = "os-default"
  preprovisioning_allowed                      = false
  device_type                                  = "windowsPc"
  hardware_hash_extraction_enabled             = true
  role_scope_tag_ids                           = ["0"]
  device_join_type                             = "microsoft_entra_joined"
  hybrid_azure_ad_join_skip_connectivity_check = false

  out_of_box_experience_setting = {
    device_usage_type               = "shared"
    privacy_settings_hidden         = true
    eula_hidden                     = true
    user_type                       = "standard"
    keyboard_selection_page_skipped = true
  }
}

# HoloLens AutopilotDeployment Profile Example
resource "microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile" "hololens" {
  display_name                                 = "hololens"
  description                                  = "hololens autopilot profile with os default locale"
  device_name_template                         = "thing-%RAND:2%"
  locale                                       = "zh-HK"
  preprovisioning_allowed                      = false
  device_type                                  = "holoLens"
  hardware_hash_extraction_enabled             = false
  role_scope_tag_ids                           = ["0"]
  device_join_type                             = "microsoft_entra_joined"
  hybrid_azure_ad_join_skip_connectivity_check = false

  out_of_box_experience_setting = {
    device_usage_type               = "shared"
    privacy_settings_hidden         = true
    eula_hidden                     = true
    user_type                       = "standard"
    keyboard_selection_page_skipped = true
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `device_join_type` (String) The type of device join to configure. Determines which Windows Autopilot deployment profile type to use. Possible values are: `microsoft_entra_joined`, `microsoft_entra_hybrid_joined`. Note: HoloLens devices must use `microsoft_entra_joined`.
- `display_name` (String) The display name of the deployment profile. Max allowed length is 200 chars. Cannot contain the following characters: ! # % ^ * ) ( - + ; ' > <
- `out_of_box_experience_setting` (Attributes) The Windows Autopilot Deployment Profile settings used by the device for the out-of-box experience. (see [below for nested schema](#nestedatt--out_of_box_experience_setting))

### Optional

- `assignments` (Attributes Set) The list of assignments for this deployment profile. (see [below for nested schema](#nestedatt--assignments))
- `description` (String) A description of the windows autopilotdeployment profile. Max allowed length is 1500 chars.
- `device_name_template` (String) The template used to name the Autopilot device. This can be a custom text and can also contain either the serial number of the device, or a randomly generated number. The total length of the text generated by the template can be no more than 15 characters. For Microsoft Entra hybrid joined type of Autopilot deployment profiles, devices are named using settings specified in Domain Join configuration.
- `device_type` (String) The Windows device type that this profile is applicable to. Possible values include `windowsPc`, `holoLens`, `surfaceHub2`, `surfaceHub2S`, `virtualMachine`, `unknownFutureValue`. The default is `windowsPc`.
- `hardware_hash_extraction_enabled` (Boolean) Select Yes to register all targeted devices to Autopilot if they are not already registered. The next time registered devices go through the Windows Out of Box Experience (OOBE), they will go through the assigned Autopilot scenario.Please note that certain Autopilot scenarios require specific minimum builds of Windows. Please make sure your device has the required minimum build to go through the scenario.Removing this profile won't remove affected devices from Autopilot. To remove a device from Autopilot, use the Windows Autopilot Devices view.Default value is FALSE.
- `hybrid_azure_ad_join_skip_connectivity_check` (Boolean) The Autopilot Hybrid Azure AD join flow will continue even if it does not establish domain controller connectivity during OOBE. This should only be set to true when using `microsoft_entra_hybrid_joined` device join type, else always false.
- `locale` (String) The locale (language) to be used when configuring the device. Possible values are: `user_select` (allows user to select language during OOBE), `os-default` (uses OS default), or specific country codes like `en-US`, `ja-JP`, `fr-FR`, etc. Default value is `os-default`.
- `management_service_app_id` (String) The Entra management service App ID which gets used during client device-based enrollment discovery.
- `preprovisioning_allowed` (Boolean) Indicates whether the user is allowed to use Windows Autopilot for pre-provisioned deployment mode during Out of Box experience (OOBE). When TRUE, indicates that Windows Autopilot for pre-provisioned deployment mode for OOBE is allowed to be used. When false, Windows Autopilot for pre-provisioned deployment mode for OOBE is not allowed. The default is FALSE.
- `role_scope_tag_ids` (Set of String) List of role scope tags for the deployment profile.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `created_date_time` (String) The date and time of when the deployment profile was created. Read-Only.
- `id` (String) The profile key.
- `last_modified_date_time` (String) The date and time of when the deployment profile was last modified. Read-Only.

<a id="nestedatt--out_of_box_experience_setting"></a>
### Nested Schema for `out_of_box_experience_setting`

Required:

- `device_usage_type` (String) The Entra join authentication type. Possible values are singleUser and shared. The default is singleUser. Possible values are: `singleUser`, `shared`, `unknownFutureValue`.
- `keyboard_selection_page_skipped` (Boolean) When TRUE, the keyboard selection page is hidden to the end user during OOBE if Language and Region are set. When FALSE, the keyboard selection page is skipped during OOBE.
- `privacy_settings_hidden` (Boolean) When TRUE, privacy settings is hidden to the end user during OOBE. When FALSE, privacy settings is shown to the end user during OOBE. Default value is FALSE.
- `user_type` (String) The type of user. Possible values are administrator and standard. Default value is administrator. Possible values are: `administrator`, `standard`, `unknownFutureValue`.

Optional:

- `eula_hidden` (Boolean) When TRUE, EULA is hidden to the end user during OOBE. When FALSE, EULA is shown to the end user during OOBE. Default value is FALSE.

Read-Only:

- `escape_link_hidden` (Boolean) When TRUE, the link that allows user to start over with a different account on company sign-in is hidden. When false, the link that allows user to start over with a different account on company sign-in is available.  This field is defaulted to TRUE for a valid api call but doesnt configure anything in the gui. This field is always required to be set to TRUE.


<a id="nestedatt--assignments"></a>
### Nested Schema for `assignments`

Required:

- `type` (String) The type of assignment target. Possible values are: `groupAssignmentTarget`, `exclusionGroupAssignmentTarget`, `allDevicesAssignmentTarget`.

Optional:

- `group_id` (String) The ID of the target group. Required when type is `groupAssignmentTarget` or `exclusionGroupAssignmentTarget`.


<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Important Notes

- **Windows Autopilot**: This resource manages Windows Autopilot deployment profiles that define the out-of-box experience (OOBE) for devices.
- **Device Configuration**: Deployment profiles control how devices are configured during the initial setup process.
- **User Experience**: Profiles can be configured to provide a customized and streamlined setup experience for end users.
- **Assignment Required**: Profiles must be assigned to Windows Autopilot device groups to take effect.
- **Profile Types**: Supports both User-Driven and Self-Deploying deployment scenarios.
- **OOBE Customization**: Configure settings like skip privacy settings, create local admin account, and join domain options.
- **Hybrid Azure AD**: Supports both cloud-only and hybrid Azure AD join scenarios.

## Import

Import is supported using the following syntax:

```shell
# {resource_id}
terraform import microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.example windows-autopilot-deployment-profile-id
```