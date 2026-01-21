---
page_title: "microsoft365_graph_beta_device_management_ios_device_compliance_policy Resource - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Manages ios device compliance policies using the /deviceManagement/deviceCompliancePolicies endpoint. This resource is used to device compliance policies define rules and settings that devices must meet to be considered compliant with organizational security requirements.
---

# microsoft365_graph_beta_device_management_ios_device_compliance_policy (Resource)

Manages ios device compliance policies using the `/deviceManagement/deviceCompliancePolicies` endpoint. This resource is used to device compliance policies define rules and settings that devices must meet to be considered compliant with organizational security requirements.

## Microsoft Documentation

- [iosDeviceCompliancePolicy resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-ioscompliancepolicy?view=graph-rest-beta)
- [Create iosDeviceCompliancePolicy](https://learn.microsoft.com/en-us/graph/api/intune-deviceconfig-ioscompliancepolicy-create?view=graph-rest-beta&tabs=http)
- [Update iosDeviceCompliancePolicy](https://learn.microsoft.com/en-us/graph/api/intune-deviceconfig-ioscompliancepolicy-update?view=graph-rest-beta&tabs=http)
- [Delete iosDeviceCompliancePolicy](https://learn.microsoft.com/en-us/graph/api/intune-devicecon fig-ioscompliancepolicy-delete?view=graph-rest-beta&tabs=http)

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this resource:

**Required:**
- `DeviceManagementConfiguration.ReadWrite.All`
- `DeviceManagementConfiguration.Read.All`

**Optional:**
- `None` `[N/A]`

## Example Usage

```terraform
# Example with minimal configuration
resource "microsoft365_graph_beta_device_management_ios_device_compliance_policy" "minimal" {
  display_name = "iOS Minimal Compliance Policy"
  description  = "Minimal iOS device compliance policy with basic security requirements"

  # Basic security requirements
  passcode_required                 = true
  security_block_jailbroken_devices = true

  # Scheduled actions for rules (required)
  scheduled_actions_for_rule = [
    {
      rule_name = "PasscodeRequired"
      scheduled_action_configurations = [
        {
          action_type        = "block"
          grace_period_hours = 0
        }
      ]
    }
  ]
}

# Example with advanced security settings
resource "microsoft365_graph_beta_device_management_ios_device_compliance_policy" "advanced" {
  display_name = "iOS Advanced Compliance Policy"
  description  = "Advanced iOS device compliance policy with strict security requirements"

  # Strict passcode requirements
  passcode_required                                    = true
  passcode_block_simple                                = true
  passcode_minimum_length                              = 8
  passcode_minimum_character_set_count                 = 3
  passcode_required_type                               = "alphanumeric"
  passcode_expiration_days                             = 30
  passcode_previous_passcode_block_count               = 5
  passcode_minutes_of_inactivity_before_lock           = 2
  passcode_minutes_of_inactivity_before_screen_timeout = 1

  # Strict OS version requirements
  os_minimum_version       = "16.0"
  os_minimum_build_version = "20A362"

  # Security settings
  security_block_jailbroken_devices                  = true
  device_threat_protection_enabled                   = true
  device_threat_protection_required_security_level   = "high"
  advanced_threat_protection_required_security_level = "secured"
  managed_email_profile_required                     = true

  # Restricted apps
  restricted_apps = [
    {
      name          = "Prohibited App"
      publisher     = "Prohibited Publisher"
      app_id        = "com.prohibited.app"
      app_store_url = "https://apps.apple.com/app/prohibited-app/id123456789"
    }
  ]

  # Scheduled actions for rules
  scheduled_actions_for_rule = [
    {
      rule_name = "PasswordRequired"
      scheduled_action_configurations = [
        {
          action_type              = "notification"
          grace_period_hours       = 0
          notification_template_id = "00000000-0000-0000-0000-000000000000"
          notification_message_cc_list = [
            "00000000-0000-0000-0000-000000000000",
            "00000000-0000-0000-0000-000000000000"
          ]
        },
        {
          action_type              = "notification"
          grace_period_hours       = 24
          notification_template_id = "00000000-0000-0000-0000-000000000000"
          notification_message_cc_list = [
            "00000000-0000-0000-0000-000000000000",
            "00000000-0000-0000-0000-000000000000"
          ]
        },
        {
          action_type              = "remoteLock"
          grace_period_hours       = 72
          notification_template_id = ""
        },
        {
          action_type              = "retire"
          grace_period_hours       = 120
          notification_template_id = ""
        },
        {
          action_type              = "block"
          grace_period_hours       = 0
          notification_template_id = ""
        }
      ]
    }
  ]

  # Assignments
  assignments = [
    # Assignment targeting all devices with an include filter
    {
      type        = "allDevicesAssignmentTarget"
      filter_id   = "00000000-0000-0000-0000-000000000000"
      filter_type = "include"
    },
    # Assignment targeting all licensed users with an exclude filter
    {
      type        = "allLicensedUsersAssignmentTarget"
      filter_id   = "00000000-0000-0000-0000-000000000000"
      filter_type = "exclude"
    },
    # Assignment targeting a specific group with include filter
    {
      type        = "groupAssignmentTarget"
      group_id    = "00000000-0000-0000-0000-000000000000"
      filter_id   = "00000000-0000-0000-0000-000000000000"
      filter_type = "include"
    },
    # Assignment targeting a specific group with exclude filter
    {
      type        = "groupAssignmentTarget"
      group_id    = "00000000-0000-0000-0000-000000000000"
      filter_id   = "00000000-0000-0000-0000-000000000000"
      filter_type = "exclude"
    },
    # Exclusion group assignments
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000000"
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000000"
    }
  ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `display_name` (String) The display name of the device compliance policy

### Optional

- `advanced_threat_protection_required_security_level` (String) Require Microsoft Defender for Endpoint minimum risk level to report noncompliance. Possible values are: unavailable, secured, low, medium, high, notSet
- `assignments` (Attributes Set) Assignments for the compliance policy. Each assignment specifies the target group and schedule for script execution. (see [below for nested schema](#nestedatt--assignments))
- `description` (String) Optional description of the resource. Maximum length is 1500 characters.
- `device_threat_protection_enabled` (Boolean) Require that devices have enabled device threat protection
- `device_threat_protection_required_security_level` (String) Require Device Threat Protection minimum risk level to report noncompliance. Possible values are: unavailable, secured, low, medium, high, notSet
- `managed_email_profile_required` (Boolean) Indicates whether or not to require a managed email profile
- `os_maximum_build_version` (String) Maximum iOS build version
- `os_maximum_version` (String) Maximum iOS version allowed.
- `os_minimum_build_version` (String) Minimum iOS build version
- `os_minimum_version` (String) Minimum iOS version allowed.
- `passcode_block_simple` (Boolean) Indicates whether or not to block simple passcodes
- `passcode_expiration_days` (Number) Number of days before the passcode expires. Valid values 1 to 65535
- `passcode_minimum_character_set_count` (Number) The number of character sets required in the password
- `passcode_minimum_length` (Number) Minimum length of passcode. Valid values 4 to 14
- `passcode_minutes_of_inactivity_before_lock` (Number) Minutes of inactivity before a passcode is required
- `passcode_minutes_of_inactivity_before_screen_timeout` (Number) Minutes of inactivity before the screen times out
- `passcode_previous_passcode_block_count` (Number) Number of previous passwords to block
- `passcode_required` (Boolean) Indicates whether or not to require a passcode
- `passcode_required_type` (String) The required passcode type. Possible values are: deviceDefault, alphanumeric, numeric
- `restricted_apps` (Attributes Set) Require the device to not have the specified apps installed. This collection can contain a maximum of 100 elements (see [below for nested schema](#nestedatt--restricted_apps))
- `role_scope_tag_ids` (Set of String) Set of scope tag IDs for this Entity instance.
- `scheduled_actions_for_rule` (Attributes List) The list of scheduled action for this rule (see [below for nested schema](#nestedatt--scheduled_actions_for_rule))
- `security_block_jailbroken_devices` (Boolean) Indicates the device should not be jailbroken. When TRUE, if the device is detected as jailbroken it will be reported non-compliant
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `id` (String) The id of the driver.

<a id="nestedatt--assignments"></a>
### Nested Schema for `assignments`

Required:

- `type` (String) Type of assignment target. Must be one of: 'allDevicesAssignmentTarget', 'allLicensedUsersAssignmentTarget', 'groupAssignmentTarget', 'exclusionGroupAssignmentTarget'.

Optional:

- `filter_id` (String) ID of the filter to apply to the assignment.
- `filter_type` (String) Type of filter to apply. Must be one of: 'include', 'exclude', or 'none'.
- `group_id` (String) The Entra ID group ID to include or exclude in the assignment. Required when type is 'groupAssignmentTarget' or 'exclusionGroupAssignmentTarget'.


<a id="nestedatt--restricted_apps"></a>
### Nested Schema for `restricted_apps`

Required:

- `name` (String) The application name

Optional:

- `app_id` (String) The application or bundle identifier of the application
- `app_store_url` (String) The Store URL of the application
- `publisher` (String) The publisher of the application


<a id="nestedatt--scheduled_actions_for_rule"></a>
### Nested Schema for `scheduled_actions_for_rule`

Required:

- `scheduled_action_configurations` (Attributes Set) The list of scheduled action configurations for this compliance policy (see [below for nested schema](#nestedatt--scheduled_actions_for_rule--scheduled_action_configurations))

Optional:

- `rule_name` (String) Name of the scheduled action rule

<a id="nestedatt--scheduled_actions_for_rule--scheduled_action_configurations"></a>
### Nested Schema for `scheduled_actions_for_rule.scheduled_action_configurations`

Required:

- `action_type` (String) What action to take. Possible values are: 'noAction', 'notification', 'block', 'retire', 'wipe', 'removeResourceAccessProfiles', 'pushNotification', 'remoteLock'.

Optional:

- `grace_period_hours` (Number) Number of hours to wait till the action will be enforced
- `notification_message_cc_list` (List of String) A list of group GUIDs to specify who to CC this notification message to
- `notification_template_id` (String) What notification Message template to use



<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.23.0-alpha | Experimental | Initial release |

## Import

Import is supported using the following syntax:

```shell
#!/bin/bash

# Import using the resource ID
terraform import microsoft365_graph_beta_device_management_ios_device_compliance_policy.example 00000000-0000-0000-0000-000000000000
```