---
page_title: "microsoft365_graph_beta_device_management_macos_device_compliance_policy Resource - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Manages macOS device compliance policies using the /deviceManagement/deviceCompliancePolicies endpoint. This resource is used to device compliance policies define rules and settings that devices must meet to be considered compliant with organizational security requirements.
---

# microsoft365_graph_beta_device_management_macos_device_compliance_policy (Resource)

Manages macOS device compliance policies using the `/deviceManagement/deviceCompliancePolicies` endpoint. This resource is used to device compliance policies define rules and settings that devices must meet to be considered compliant with organizational security requirements.

## Microsoft Documentation

- [windowsDeviceCompliancePolicy resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-macoscompliancepolicy?view=graph-rest-beta)
- [Create macosDeviceCompliancePolicy](https://learn.microsoft.com/en-us/graph/api/intune-deviceconfig-macoscompliancepolicy-create?view=graph-rest-beta&tabs=http)
- [Update macosDeviceCompliancePolicy](https://learn.microsoft.com/en-us/graph/api/intune-deviceconfig-macoscompliancepolicy-update?view=graph-rest-beta&tabs=http)
- [Delete macosDeviceCompliancePolicy](https://learn.microsoft.com/en-us/graph/api/intune-devicecon fig-macoscompliancepolicy-delete?view=graph-rest-beta&tabs=http)

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
resource "microsoft365_graph_beta_device_management_macos_device_compliance_policy" "minimal" {
  display_name = "macOS Minimal Compliance Policy"
  description  = "Minimal macOS device compliance policy with basic security requirements"

  # Basic security requirements
  password_required          = true
  storage_require_encryption = true
  firewall_enabled           = true

  # Scheduled actions for rules (required)
  scheduled_actions_for_rule = [
    {
      rule_name = "PasswordRequired"
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
resource "microsoft365_graph_beta_device_management_macos_device_compliance_policy" "advanced" {
  display_name = "macOS Advanced Compliance Policy"
  description  = "Advanced macOS device compliance policy with strict security requirements"

  # Strict password requirements
  password_required                          = true
  password_block_simple                      = true
  password_minimum_length                    = 12
  password_minimum_character_set_count       = 4
  password_required_type                     = "alphanumeric"
  password_expiration_days                   = 60
  password_previous_password_block_count     = 10
  password_minutes_of_inactivity_before_lock = 5

  # Strict OS version requirements
  os_minimum_version       = "14.0"
  os_minimum_build_version = "23A344"

  # Maximum security settings
  system_integrity_protection_enabled                = true
  device_threat_protection_enabled                   = true
  device_threat_protection_required_security_level   = "high"
  advanced_threat_protection_required_security_level = "high"
  storage_require_encryption                         = true
  gatekeeper_allowed_app_source                      = "macAppStore"

  # Strict firewall settings
  firewall_enabled             = true
  firewall_block_all_incoming  = true
  firewall_enable_stealth_mode = true

  # Scheduled actions with aggressive enforcement
  scheduled_actions_for_rule = [
    {
      rule_name = "PasswordRequired"
      scheduled_action_configurations = [
        {
          action_type              = "notification"
          grace_period_hours       = 0
          notification_template_id = "426e6351-c6ff-44d3-910d-8b937ee30bdd"
        },
        {
          action_type        = "block"
          grace_period_hours = 1
        },
        {
          action_type        = "retire"
          grace_period_hours = 24
        }
      ]
    }
  ]

  # Target specific high-security groups
  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = "high-security-macos-devices-group-id"
    }
  ]
}

resource "microsoft365_graph_beta_device_management_macos_device_compliance_policy" "comprehensive" {
  display_name = "macOS Comprehensive Compliance Policy"
  description  = "Comprehensive macOS device compliance policy with all available security settings"

  # Password requirements
  password_required                          = true
  password_block_simple                      = true
  password_minimum_length                    = 8
  password_minimum_character_set_count       = 3
  password_required_type                     = "alphanumeric"
  password_expiration_days                   = 90
  password_previous_password_block_count     = 5
  password_minutes_of_inactivity_before_lock = 15

  # OS version requirements
  os_minimum_version       = "13.0"
  os_maximum_version       = "14.0"
  os_minimum_build_version = "22A380"
  os_maximum_build_version = "23A344"

  # Security requirements
  system_integrity_protection_enabled                = true
  device_threat_protection_enabled                   = true
  device_threat_protection_required_security_level   = "medium"
  advanced_threat_protection_required_security_level = "low"
  storage_require_encryption                         = true
  gatekeeper_allowed_app_source                      = "macAppStoreAndIdentifiedDevelopers"

  # Firewall requirements
  firewall_enabled             = true
  firewall_block_all_incoming  = false
  firewall_enable_stealth_mode = true

  # Role scope tags
  role_scope_tag_ids = ["0"]

  # Scheduled actions for rules
  scheduled_actions_for_rule = [
    {
      rule_name = "PasswordRequired"
      scheduled_action_configurations = [
        {
          action_type              = "notification"
          grace_period_hours       = 0
          notification_template_id = "426e6351-c6ff-44d3-910d-8b937ee30bdd"
          notification_message_cc_list = [
            "aa856a09-cf0c-4b31-a315-cb53251e54d8",
            "a77240dc-2827-47af-8fcb-e209a67e176a"
          ]
        },
        {
          action_type              = "notification"
          grace_period_hours       = 24
          notification_template_id = "bbf43ceb-5e68-428b-8ad3-00c9efb54210"
          notification_message_cc_list = [
            "91710c72-1358-4438-b0b2-70eb32b542dd",
            "aa856a09-cf0c-4b31-a315-cb53251e54d8"
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
      filter_id   = "80f8c0a5-f3ec-4936-bcbc-420dc0ca3665"
      filter_type = "include"
    },
    # Assignment targeting all licensed users with an exclude filter
    {
      type        = "allLicensedUsersAssignmentTarget"
      filter_id   = "2d7956fb-e5b5-4fa3-90b2-5bee9bee7883"
      filter_type = "exclude"
    },
    # Assignment targeting a specific group with include filter
    {
      type        = "groupAssignmentTarget"
      group_id    = "51a96cdd-4b9b-4849-b416-8c94a6d88797"
      filter_id   = "80f8c0a5-f3ec-4936-bcbc-420dc0ca3665"
      filter_type = "include"
    },
    # Assignment targeting a specific group with exclude filter
    {
      type        = "groupAssignmentTarget"
      group_id    = "b15228f4-9d49-41ed-9b4f-0e7c721fd9c2"
      filter_id   = "2d7956fb-e5b5-4fa3-90b2-5bee9bee7883"
      filter_type = "exclude"
    },
    # Exclusion group assignments
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "b8c661c2-fa9a-4351-af86-adc1729c343f"
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "f6ebd6ff-501e-4b3d-a00b-a2e102c3fa0f"
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
- `firewall_block_all_incoming` (Boolean) Corresponds to 'Block all incoming connections' option
- `firewall_enable_stealth_mode` (Boolean) Corresponds to 'Enable stealth mode' option
- `firewall_enabled` (Boolean) Whether the firewall should be enabled or not
- `gatekeeper_allowed_app_source` (String) App source options for Gatekeeper. Possible values are: notConfigured, macAppStore, macAppStoreAndIdentifiedDevelopers, anywhere
- `os_maximum_build_version` (String) Maximum macOS build version
- `os_maximum_version` (String) Maximum macOS version allowed.
- `os_minimum_build_version` (String) Minimum macOS build version
- `os_minimum_version` (String) Minimum macOS version allowed.
- `password_block_simple` (Boolean) Indicates whether or not to block simple passwords
- `password_expiration_days` (Number) Number of days before the password expires
- `password_minimum_character_set_count` (Number) The number of character sets required in the password
- `password_minimum_length` (Number) Minimum length of password
- `password_minutes_of_inactivity_before_lock` (Number) Minutes of inactivity before a password is required
- `password_previous_password_block_count` (Number) Number of previous passwords to block
- `password_required` (Boolean) Whether or not to require a password
- `password_required_type` (String) The required password type. Possible values are: deviceDefault, alphanumeric, numeric
- `role_scope_tag_ids` (Set of String) Set of scope tag IDs for this Entity instance.
- `scheduled_actions_for_rule` (Attributes List) The list of scheduled action for this rule (see [below for nested schema](#nestedatt--scheduled_actions_for_rule))
- `storage_require_encryption` (Boolean) Require encryption on macOS devices
- `system_integrity_protection_enabled` (Boolean) Require System Integrity Protection (SIP) to be enabled
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

# Import an existing macOS Device Compliance Policy
# Replace the ID with the actual ID of your policy from Microsoft Graph API

terraform import microsoft365_graph_beta_device_management_macos_device_compliance_policy.basic 00000000-0000-0000-0000-000000000001
```