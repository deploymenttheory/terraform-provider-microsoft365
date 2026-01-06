---
page_title: "microsoft365_graph_beta_device_management_android_device_owner_compliance_policy Resource - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Manages Android Device Owner compliance policies in Microsoft Intune using the /deviceManagement/deviceCompliancePolicies endpoint. Device compliance policies define rules and settings that devices must meet to be considered compliant with organizational security requirements.
---

# microsoft365_graph_beta_device_management_android_device_owner_compliance_policy (Resource)

Manages Android Device Owner compliance policies in Microsoft Intune using the `/deviceManagement/deviceCompliancePolicies` endpoint. Device compliance policies define rules and settings that devices must meet to be considered compliant with organizational security requirements.

## Microsoft Documentation

- [androidDeviceOwnerCompliancePolicy resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-androiddeviceownercompliancepolicy?view=graph-rest-beta)
- [Create androidDeviceOwnerCompliancePolicy](https://learn.microsoft.com/en-us/graph/api/intune-deviceconfig-androiddeviceownercompliancepolicy-create?view=graph-rest-beta&tabs=http)
- [Update androidDeviceOwnerCompliancePolicy](https://learn.microsoft.com/en-us/graph/api/intune-deviceconfig-androiddeviceownercompliancepolicy-update?view=graph-rest-beta&tabs=http)
- [Delete androidDeviceOwnerCompliancePolicy](https://learn.microsoft.com/en-us/graph/api/intune-deviceconfig-androiddeviceownercompliancepolicy-delete?view=graph-rest-beta&tabs=http)

## API Permissions

The following API permissions are required in order to use this resource.

### Microsoft Graph

- **Application**: `DeviceManagementConfiguration.ReadWrite.All` , `DeviceManagementConfiguration.Read.All`

## Example Usage

```terraform
# Example with minimal configuration
resource "microsoft365_graph_beta_device_management_android_device_owner_compliance_policy" "minimal" {
  display_name = "Android Device Owner Minimal Compliance Policy"
  description  = "Minimal Android device owner compliance policy with basic security requirements"

  # Basic password requirements
  password_required       = true
  password_minimum_length = 6

  # Security settings
  security_block_jailbroken_devices = true
  storage_require_encryption        = true

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

# Example with comprehensive security settings
resource "microsoft365_graph_beta_device_management_android_device_owner_compliance_policy" "comprehensive" {
  display_name = "Android Device Owner Comprehensive Compliance Policy"
  description  = "Comprehensive Android device owner compliance policy with advanced security requirements"

  # Threat protection settings
  device_threat_protection_enabled                   = true
  device_threat_protection_required_security_level   = "medium"
  advanced_threat_protection_required_security_level = "high"

  # Security settings
  security_block_jailbroken_devices                        = true
  security_require_safety_net_attestation_basic_integrity  = true
  security_require_safety_net_attestation_certified_device = true
  security_require_intune_app_integrity                    = true
  require_no_pending_system_updates                        = true
  security_required_android_safety_net_evaluation_type     = "hardwareBacked"

  # OS version requirements
  os_minimum_version               = "14"
  os_maximum_version               = "15"
  min_android_security_patch_level = "February 1, 2025"

  # Comprehensive password requirements
  password_required                          = true
  password_minimum_length                    = 12
  password_minimum_letter_characters         = 2
  password_minimum_lower_case_characters     = 1
  password_minimum_upper_case_characters     = 1
  password_minimum_numeric_characters        = 2
  password_minimum_symbol_characters         = 1
  password_minimum_non_letter_characters     = 3
  password_required_type                     = "alphanumericWithSymbols"
  password_minutes_of_inactivity_before_lock = 5
  password_expiration_days                   = 90
  password_previous_password_count_to_block  = 5

  # Storage settings
  storage_require_encryption = true

  # Role scope tags
  role_scope_tag_ids = ["0", "1"]

  # Scheduled actions for rules
  scheduled_actions_for_rule = [
    {
      rule_name = "PasswordRequired"
      scheduled_action_configurations = [
        {
          action_type              = "notification"
          grace_period_hours       = 24
          notification_template_id = "00000000-0000-0000-0000-000000000000"
          notification_message_cc_list = [
            "00000000-0000-0000-0000-000000000000"
          ]
        },
        {
          action_type        = "block"
          grace_period_hours = 72
        }
      ]
    }
  ]

  # Assignments
  assignments = [
    # Assignment targeting all devices
    {
      type        = "allDevicesAssignmentTarget"
      filter_id   = "00000000-0000-0000-0000-000000000000"
      filter_type = "none"
    },
    # Assignment targeting a specific group
    {
      type        = "groupAssignmentTarget"
      group_id    = "00000000-0000-0000-0000-000000000000"
      filter_id   = "00000000-0000-0000-0000-000000000000"
      filter_type = "none"
    },
    # Exclusion group assignment
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000000"
    }
  ]
}

# Example with moderate security settings for enterprise use
resource "microsoft365_graph_beta_device_management_android_device_owner_compliance_policy" "enterprise" {
  display_name = "Android Device Owner Enterprise Compliance Policy"
  description  = "Enterprise Android device owner compliance policy balancing security and usability"

  # Threat protection settings
  device_threat_protection_enabled                 = true
  device_threat_protection_required_security_level = "low"

  # Security settings
  security_block_jailbroken_devices                       = true
  security_require_safety_net_attestation_basic_integrity = true
  security_require_intune_app_integrity                   = true
  security_required_android_safety_net_evaluation_type    = "basic"

  # OS version requirements - allowing broader range for compatibility
  os_minimum_version               = "13"
  min_android_security_patch_level = "January 1, 2024"

  # Balanced password requirements
  password_required                          = true
  password_minimum_length                    = 8
  password_minimum_letter_characters         = 1
  password_minimum_numeric_characters        = 1
  password_required_type                     = "alphanumeric"
  password_minutes_of_inactivity_before_lock = 15
  password_expiration_days                   = 180
  password_previous_password_count_to_block  = 3

  # Storage encryption required
  storage_require_encryption = true

  # Scheduled actions with grace periods for user adaptation
  scheduled_actions_for_rule = [
    {
      rule_name = "PasswordRequired"
      scheduled_action_configurations = [
        {
          action_type        = "notification"
          grace_period_hours = 48
        },
        {
          action_type        = "block"
          grace_period_hours = 168 # 7 days
        }
      ]
    }
  ]

  # Assignment to enterprise device group
  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000000" # Replace with actual enterprise device group ID
    }
  ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `display_name` (String) Admin provided name of the device configuration. Inherited from deviceCompliancePolicy

### Optional

- `advanced_threat_protection_required_security_level` (String) Indicates the Microsoft Defender for Endpoint (also referred to Microsoft Defender Advanced Threat Protection (MDATP)) minimum risk level to report noncompliance. Possible values are: unavailable, secured, low, medium, high, notSet.
- `assignments` (Attributes Set) Assignments for the compliance policy. Each assignment specifies the target group and schedule for script execution. (see [below for nested schema](#nestedatt--assignments))
- `description` (String) Optional description of the resource. Maximum length is 1500 characters.
- `device_threat_protection_enabled` (Boolean) Indicates whether the policy requires devices have device threat protection enabled. When TRUE, threat protection is enabled. When FALSE, threat protection is not enabled. Default is FALSE.
- `device_threat_protection_required_security_level` (String) Indicates the minimum mobile threat protection risk level to that results in Intune reporting device noncompliance. Possible values are: unavailable, secured, low, medium, high, notSet.
- `min_android_security_patch_level` (String) Indicates the minimum Android security patch level required to mark the device as compliant. Must be a valid date format (YYYY-MM-DD). Example: 2026-10-01, 2026-10-31 etc.
- `os_maximum_version` (String) Indicates the maximum Android version required to mark the device as compliant. For example: '15'
- `os_minimum_version` (String) Indicates the minimum Android version required to mark the device as compliant. For example: '14'
- `password_expiration_days` (Number) Indicates the number of days before the password expires. Valid values 1 to 365.
- `password_minimum_length` (Number) Indicates the minimum password length required to mark the device as compliant. Valid values are 4 to 16, inclusive. Valid values 4 to 16
- `password_minimum_letter_characters` (Number) Indicates the minimum number of letter characters required for device password for the device to be marked compliant. Valid values 1 to 16.
- `password_minimum_lower_case_characters` (Number) Indicates the minimum number of lower case characters required for device password for the device to be marked compliant. Valid values 1 to 16.
- `password_minimum_non_letter_characters` (Number) Indicates the minimum number of non-letter characters required for device password for the device to be marked compliant. Valid values 1 to 16.
- `password_minimum_numeric_characters` (Number) Indicates the minimum number of numeric characters required for device password for the device to be marked compliant. Valid values 1 to 16.
- `password_minimum_symbol_characters` (Number) Indicates the minimum number of symbol characters required for device password for the device to be marked compliant. Valid values 1 to 16.
- `password_minimum_upper_case_characters` (Number) Indicates the minimum number of upper case letter characters required for device password for the device to be marked compliant. Valid values 1 to 16.
- `password_minutes_of_inactivity_before_lock` (Number) Indicates the number of minutes of inactivity before a password is required.
- `password_previous_password_count_to_block` (Number) Indicates the number of previous passwords to block. Valid values 1 to 24.
- `password_required` (Boolean) Indicates whether a password is required to unlock the device. When TRUE, there must be a password set that unlocks the device for the device to be marked as compliant. When FALSE, a device is marked as compliant whether or not a password is set as required to unlock the device. Default is FALSE.
- `password_required_type` (String) Indicates the password complexity requirement for the device to be marked compliant. Possible values are: deviceDefault, required, numeric, numericComplex, alphabetic, alphanumeric, alphanumericWithSymbols, lowSecurityBiometric, customPassword.
- `require_no_pending_system_updates` (Boolean) Indicates whether the device has pending security or OS updates and sets the compliance state accordingly. When TRUE, checks if there are any pending system updates on each check in and if there are any pending security or OS version updates (System Updates), the device will be reported as non-compliant. If set to FALSE, then checks for any pending security or OS version updates (System Updates) are done without impact to device compliance state. Default is FALSE.
- `role_scope_tag_ids` (Set of String) List of Scope Tags for this Entity instance. Inherited from deviceCompliancePolicy
- `scheduled_actions_for_rule` (Attributes List) The list of scheduled action for this rule (see [below for nested schema](#nestedatt--scheduled_actions_for_rule))
- `security_block_jailbroken_devices` (Boolean) Indicates the device should not be rooted. When TRUE, if the device is detected as rooted it will be reported non-compliant. When FALSE, the device is not reported as non-compliant regardless of device rooted state. Default is FALSE.
- `security_require_intune_app_integrity` (Boolean) Indicates whether Intune application integrity is required to mark the device as compliant. When TRUE, Intune checks that the Intune app installed on fully managed, dedicated, or corporate-owned work profile Android Enterprise enrolled devices, is the one provided by Microsoft from the Managed Google Play store. If the check fails, the device will be reported as non-compliant. Default is FALSE.
- `security_require_safety_net_attestation_basic_integrity` (Boolean) Indicates whether the compliance check will validate the Google Play Integrity check. When TRUE, the Google Play integrity basic check must pass to consider the device compliant. When FALSE, the Google Play integrity basic check can pass or fail and the device will be considered compliant. Default is FALSE.
- `security_require_safety_net_attestation_certified_device` (Boolean) Indicates whether the compliance check will validate the Google Play Integrity check. When TRUE, the Google Play integrity device check must pass to consider the device compliant. When FALSE, the Google Play integrity device check can pass or fail and the device will be considered compliant. Default is FALSE.
- `security_required_android_safety_net_evaluation_type` (String) Indicates the types of measurements and reference data used to evaluate the device SafetyNet evaluation. Evaluation is completed on the device to assess device integrity based on checks defined by Android and built into the device hardware, for example, compromised OS version or root detection. Possible values are: basic, hardwareBacked, with default value of basic.
- `storage_require_encryption` (Boolean) Indicates whether encryption on Android devices is required to mark the device as compliant.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `id` (String) Key of the entity. Inherited from deviceCompliancePolicy

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

# Import using the resource ID
terraform import microsoft365_graph_beta_device_management_android_device_owner_compliance_policy.example 00000000-0000-0000-0000-000000000000
```