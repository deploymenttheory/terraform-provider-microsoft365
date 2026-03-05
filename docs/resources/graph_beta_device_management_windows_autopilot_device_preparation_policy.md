---
page_title: "microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy Resource - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Manages Windows Autopilot Device Preparation Policy using the /deviceManagement/configurationPolicies endpoint. This resource is used to windows Autopilot Device Preparation is used to set up and configure new devices, getting them ready for productive use by delivering consistent configurations and enhancing the setup experience.
---

# microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy (Resource)

Manages Windows Autopilot Device Preparation Policy using the `/deviceManagement/configurationPolicies` endpoint. This resource is used to windows Autopilot Device Preparation is used to set up and configure new devices, getting them ready for productive use by delivering consistent configurations and enhancing the setup experience.

## Microsoft Documentation

- [Windows Autopilot Device Preparation](https://learn.microsoft.com/en-us/autopilot/device-preparation/overview)
- [deviceManagementConfigurationPolicy resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-devicemanagementconfigurationpolicy?view=graph-rest-beta)

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this resource:

**Required:**
- `DeviceManagementApps.Read.All`
- `DeviceManagementConfiguration.Read.All`
- `DeviceManagementConfiguration.ReadWrite.All`
- `Directory.Read.All`
- `Group.Read.All`
- `GroupMember.Read.All`

**Optional:**
- `None` `[N/A]`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.16.0-alpha | Experimental | Initial release |
| v0.47.0-alpha | Experimental | Added updates for deployment scope filters |

## Known Issues

### User-Driven Mode Authentication Limitation - 05/03/2026

There is currently a breaking bug with the Microsoft Graph API endpoint `POST /deviceManagement/configurationPolicies/{deviceManagementConfigurationPolicyId}/setEnrollmentTimeDeviceMembershipTarget` that affects **user-driven mode** configurations.

**Issue Summary:**
- The endpoint returns a `500` error when using application-only authentication (e.g service principal with client credentials)
- The same endpoint succeeds when using user-delegated authentication (interactive login)
- This affects the ability to set the device security group for user-driven policies
- **Automatic mode** configurations work correctly as they do not require this endpoint

**Impact:**
- User-driven mode policies cannot be created or updated using service principal authentication
- Automatic mode policies are not affected and work as expected

**Workaround:**
Currently, the only workaround is to use user-delegated authentication (e.g., `az login` with a service account that has the Intune Administrator role). However, this provider does not currently support Azure CLI authentication.

**Reference:**
- Test configuration: [003_scenario_user_driven_minimal.tf](https://github.com/deploymenttheory/terraform-provider-microsoft365/blob/main/internal/services/resources/device_management/graph_beta/windows_autopilot_device_preparation_policy/tests/terraform/acceptance/003_scenario_user_driven_minimal.tf)

## Example Usage

### Automatic Mode

```terraform
# ==============================================================================
# Automatic Mode - Maximal Configuration
# ==============================================================================
# This example demonstrates an automatic deployment with apps and scripts.
# Automatic mode requires minimal configuration and is ideal for shared devices.

resource "microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy" "auto_maximal" {
  name               = "Autopilot DPP - Automatic Maximal"
  description        = "Automatic mode maximal configuration with apps and scripts"
  role_scope_tag_ids = ["0"]

  deployment_settings = {
    deployment_type = "enrollment_autopilot_dpp_deploymenttype_1" # Automatic
  }

  # Allow specific apps during device preparation
  allowed_apps = [
    {
      app_id   = "12345678-1234-1234-1234-123456789012" # Replace with your WinGet app ID
      app_type = "winGetApp"
    }
  ]

  # Allow specific scripts during device preparation
  allowed_scripts = [
    "87654321-4321-4321-4321-210987654321" # Replace with your Windows Platform Script ID
  ]

  timeouts = {
    create = "180s"
    read   = "30s"
    update = "180s"
    delete = "60s"
  }
}
```

### User-Driven Mode

```terraform
# ==============================================================================
# User-Driven Mode - Maximal Configuration
# ==============================================================================
# This example demonstrates a user-driven deployment with enhanced mode features,
# OOBE settings, apps, and scripts. User-driven mode provides more control and
# customization options for the end-user experience.

resource "microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy" "ud_maximal" {
  name                  = "Autopilot DPP - User-Driven Maximal"
  description           = "User-driven mode maximal configuration with enhanced mode features"
  role_scope_tag_ids    = ["0"]
  device_security_group = "12345678-1234-1234-1234-123456789012" # Replace with your device security group ID

  deployment_settings = {
    deployment_mode = "enrollment_autopilot_dpp_deploymentmode_1" # Enhanced
    deployment_type = "enrollment_autopilot_dpp_deploymenttype_0" # User-driven
    account_type    = "enrollment_autopilot_dpp_accountype_1"     # Standard user
  }

  # Configure the out-of-box experience
  oobe_settings = {
    timeout_in_minutes   = 120
    custom_error_message = "Please contact your IT administrator for assistance with device setup."
    allow_skip           = true
    allow_diagnostics    = true
  }

  allowed_apps = [
    {
      app_id   = "12345678-1234-1234-1234-123456789012" # Replace with your WinGet app ID
      app_type = "winGetApp"
    },
    {
      app_id   = "12345678-1234-1234-1234-234567890123" # Replace with your WinGet app ID
      app_type = "win32LobApp"
    },
    {
      app_id   = "12345678-1234-1234-1234-345678901234" # Replace with your Win32CatalogApp app ID
      app_type = "win32CatalogApp"
    },
    {
      app_id   = "12345678-1234-1234-1234-456789012345" # Replace with your OfficeSuiteApp app ID
      app_type = "officeSuiteApp"
    },
  ]

  allowed_scripts = [
    "87654321-4321-4321-4321-210987654321",
    "87654321-4321-4321-4321-210987654322", 
  ]

  timeouts = {
    create = "180s"
    read   = "30s"
    update = "180s"
    delete = "60s"
  }
}
```

### User-Driven Mode - With Assignments

```terraform
# ==============================================================================
# User-Driven Mode - With Maximal Assignments
# ==============================================================================
# This example demonstrates how to configure policy assignments using both
# group-based targeting and all licensed users. Assignments determine which
# users will have this policy applied during device enrollment.

resource "microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy" "ud_max_assign" {
  name                  = "Autopilot DPP - User-Driven with Assignments"
  description           = "User-driven mode with maximal assignments demonstrating group and all licensed users targeting"
  role_scope_tag_ids    = ["0"]
  device_security_group = "12345678-1234-1234-1234-123456789012" # Replace with your device security group ID

  deployment_settings = {
    deployment_mode = "enrollment_autopilot_dpp_deploymentmode_1" # Enhanced
    deployment_type = "enrollment_autopilot_dpp_deploymenttype_0" # User-driven
    account_type    = "enrollment_autopilot_dpp_accountype_1"     # Standard user
  }

  # Configure the out-of-box experience
  oobe_settings = {
    timeout_in_minutes   = 90
    custom_error_message = "Please contact your IT administrator for assistance."
    allow_skip           = true
    allow_diagnostics    = true
  }

  # Allow specific apps during device preparation
  allowed_apps = [
    {
      app_id   = "12345678-1234-1234-1234-123456789012" # Replace with your WinGet app ID
      app_type = "winGetApp"
    }
  ]

  # Allow specific scripts during device preparation
  allowed_scripts = [
    "87654321-4321-4321-4321-210987654321" # Replace with your Windows Platform Script ID
  ]

  # Configure policy assignments
  # This demonstrates multiple assignment types:
  # - All licensed users (applies to all users with Intune licenses)
  # - Specific security groups (targeted deployment)
  assignments = [
    {
      type = "allLicensedUsersAssignmentTarget"
    },
    {
      type     = "groupAssignmentTarget"
      group_id = "11111111-1111-1111-1111-111111111111" # Replace with your group ID
    },
    {
      type     = "groupAssignmentTarget"
      group_id = "22222222-2222-2222-2222-222222222222" # Replace with your group ID
    },
    {
      type     = "groupAssignmentTarget"
      group_id = "33333333-3333-3333-3333-333333333333" # Replace with your group ID
    },
  ]

  timeouts = {
    create = "180s"
    read   = "30s"
    update = "180s"
    delete = "60s"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Policy name

### Optional

- `allowed_apps` (Attributes List) List of applications that are allowed to be installed during the Windows Autopilot Device Preparation process. Maximum of 10 items. (see [below for nested schema](#nestedatt--allowed_apps))
- `allowed_scripts` (List of String) List of script IDs that are allowed to be executed during the Windows Autopilot Device Preparation process. Maximum of 10 items.
- `assignments` (Attributes Set) Assignments for the device configuration. Each assignment specifies the target group and schedule for script execution. Supports group filters. (see [below for nested schema](#nestedatt--assignments))
- `deployment_settings` (Attributes) Deployment settings for the Windows Autopilot Device Preparation policy. The deployment_type field is required and determines the policy template. User-driven mode (deployment_type_0) requires additional fields like deployment_mode and account_type. Self-deploying mode (deployment_type_1) only requires deployment_type. (see [below for nested schema](#nestedatt--deployment_settings))
- `description` (String) Optional description of the resource. Maximum length is 1500 characters.
- `device_security_group` (String) The ID of the assigned security device group that devices will be automatically added to during the Windows Autopilot Device Preparation flow. This group must have the 'Intune Provisioning Client' service principal (AppId: f1346770-5b25-470b-88bd-d5744ab7952c) set as its owner. In some tenants, this service principal may appear as 'Intune Autopilot ConfidentialClient'. If the Intune Provisioning Client or Intune Autopilot ConfidentialClient service principal with AppId of f1346770-5b25-470b-88bd-d5744ab7952c isn't available either in the list of objects or when searching, see [Adding the Intune Provisioning Client service principal](https://learn.microsoft.com/en-us/autopilot/device-preparation/tutorial/user-driven/entra-join-device-group#adding-the-intune-provisioning-client-service-principal).
- `oobe_settings` (Attributes) Out-of-box experience settings for the Windows Autopilot Device Preparation policy. Required for user-driven mode (deployment_type_0), not applicable for self-deploying/automatic mode (deployment_type_1). (see [below for nested schema](#nestedatt--oobe_settings))
- `role_scope_tag_ids` (Set of String) Set of scope tag IDs for this Entity instance.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `created_date_time` (String) Creation date and time of the policy
- `id` (String) The unique identifier for this policy
- `is_assigned` (Boolean) Indicates if the policy is assigned to any scope
- `last_modified_date_time` (String) Last modification date and time of the policy
- `platforms` (String) The platforms this policy applies to (e.g., windows10)
- `settings_count` (Number) Number of settings with the policy. This will change over time as the resource is updated.
- `technologies` (String) The technology this policy is using (e.g., enrollment)
- `template_family` (String) The template family for this policy (e.g., enrollmentConfiguration)
- `template_id` (String) The template ID used by this policy

<a id="nestedatt--allowed_apps"></a>
### Nested Schema for `allowed_apps`

Required:

- `app_id` (String) The ID of the application.
- `app_type` (String) The type of the application. Valid values are: 'winGetApp', 'win32LobApp', 'win32CatalogApp', 'officeSuiteApp', 'windowsUniversalAppX'.


<a id="nestedatt--assignments"></a>
### Nested Schema for `assignments`

Required:

- `type` (String) Type of assignment target. Must be one of: 'allLicensedUsersAssignmentTarget', 'groupAssignmentTarget'.

Optional:

- `filter_id` (String) ID of the filter to apply to the assignment.
- `filter_type` (String) Type of filter to apply. Must be one of: 'include', 'exclude', or 'none'.
- `group_id` (String) The Entra ID group ID to include in the assignment. Required when type is 'groupAssignmentTarget'.


<a id="nestedatt--deployment_settings"></a>
### Nested Schema for `deployment_settings`

Required:

- `deployment_type` (String) The deployment type determines the policy template and available settings. Valid values are: 'enrollment_autopilot_dpp_deploymenttype_0' (User-driven mode with full configuration, device security group, and assignments) or 'enrollment_autopilot_dpp_deploymenttype_1' (Self-deploying/automatic mode with simpler configuration).

Optional:

- `account_type` (String) The account type for users in the Windows Autopilot Device Preparation policy. Required for user-driven mode (deployment_type_0). Valid values are: 'enrollment_autopilot_dpp_accountype_0' (Administrator) or 'enrollment_autopilot_dpp_accountype_1' (Standard User).
- `deployment_mode` (String) The deployment mode for the Windows Autopilot Device Preparation policy. Required for user-driven mode (deployment_type_0). Valid values are: 'enrollment_autopilot_dpp_deploymentmode_0' (Standard mode) or 'enrollment_autopilot_dpp_deploymentmode_1' (Enhanced mode).

Read-Only:

- `join_type` (String) The join type for the Windows Autopilot Device Preparation policy. Always set to 'enrollment_autopilot_dpp_jointype_0' (Entra ID joined). Hybrid join is not supported.


<a id="nestedatt--oobe_settings"></a>
### Nested Schema for `oobe_settings`

Optional:

- `allow_diagnostics` (Boolean) Whether to allow users to access diagnostics information during setup
- `allow_skip` (Boolean) Whether to allow users to skip setup after multiple failed attempts
- `custom_error_message` (String) The custom error message to display if the deployment fails. Maximum length is 1000 characters.
- `timeout_in_minutes` (Number) The timeout in minutes for the Windows Autopilot Device Preparation policy. Valid range is 15-720 minutes.


<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Important Notes

- **Device Security Group**: The specified device security group must be an assigned security device group with the "Intune Provisioning Client" service principal (AppId: f1346770-5b25-470b-88bd-d5744ab7952c) set as its owner. Devices will be automatically added to this group during the Windows Autopilot device preparation deployment process. The service principal may sometimes appear as "Intune Autopilot ConfidentialClient" in some tenants.
- **Deployment Settings**: Configure how the device will be deployed, including join type and account permissions.
- **OOBE Settings**: Control the out-of-box experience, including timeouts and error messages.
- **Application and Script Allowlists**: Specify which applications and scripts can be installed/executed during setup.
- **Assignments**: Policies can be assigned to specific user groups for targeted deployment.

## Import

Import is supported using the following syntax:

```shell
#!/bin/bash
terraform import microsoft365_graph_beta_windows_autopilot_device_preparation_policy.example 00000000-0000-0000-0000-000000000000
``` 