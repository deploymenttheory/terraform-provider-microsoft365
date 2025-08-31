---
page_title: "microsoft365_graph_beta_device_management_windows_device_compliance_policy Resource - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Manages Windows device compliance policies in Microsoft Intune using the /deviceManagement/deviceCompliancePolicies endpoint. Device compliance policies define rules and settings that devices must meet to be considered compliant with organizational security requirements. you can find out more here: 'https://learn.microsoft.com/en-us/intune/intune-service/protect/compliance-policy-create-windows'.
---

# microsoft365_graph_beta_device_management_windows_device_compliance_policy (Resource)

Manages Windows device compliance policies in Microsoft Intune using the `/deviceManagement/deviceCompliancePolicies` endpoint. Device compliance policies define rules and settings that devices must meet to be considered compliant with organizational security requirements. you can find out more here: 'https://learn.microsoft.com/en-us/intune/intune-service/protect/compliance-policy-create-windows'.

## Microsoft Documentation

- [windowsDeviceCompliancePolicy resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-windows10compliancepolicy?view=graph-rest-beta)
- [Create windowsDeviceCompliancePolicy](https://learn.microsoft.com/en-us/graph/api/intune-deviceconfig-windows10compliancepolicy-create?view=graph-rest-beta&tabs=http)
- [Update windowsDeviceCompliancePolicy](https://learn.microsoft.com/en-us/graph/api/intune-deviceconfig-windows10compliancepolicy-update?view=graph-rest-beta&tabs=http)
- [Delete windowsDeviceCompliancePolicy](https://learn.microsoft.com/en-us/graph/api/intune-deviceconfig-windows10compliancepolicy-delete?view=graph-rest-beta&tabs=http)

## API Permissions

The following API permissions are required in order to use this resource.

### Microsoft Graph

- **Application**: `DeviceManagementConfiguration.ReadWrite.All` , `DeviceManagementConfiguration.Read.All`

## Example Usage

```terraform
// Example 1: Custom Compliance Policy
resource "microsoft365_graph_beta_device_management_windows_device_compliance_policy" "custom_compliance" {
  display_name       = "tf-reg-example-windows-device-compliance-policy-custom-compliance"
  description        = "tf-reg-example-windows-device-compliance-policy-custom-compliance"
  role_scope_tag_ids = ["0"]

  # Custom compliance script
  custom_compliance_required = true
  device_compliance_policy_script = {
    device_compliance_script_id = microsoft365_graph_beta_device_management_windows_device_compliance_script.example.id
    rules_content = jsonencode({
      "Rules" : [
        {
          "SettingName" : "BiosVersion",
          "Operator" : "GreaterEquals",
          "DataType" : "Version",
          "Operand" : "2.3",
          "MoreInfoUrl" : "https://bing.com",
          "RemediationStrings" : [
            {
              "Language" : "en_US",
              "Title" : "BIOS Version needs to be upgraded to at least 2.3. Value discovered was {ActualValue}.",
              "Description" : "BIOS must be updated. Please refer to the link above"
            },
            {
              "Language" : "de_DE",
              "Title" : "BIOS-Version muss auf mindestens 2.3 aktualisiert werden. Der erkannte Wert lautet {ActualValue}.",
              "Description" : "BIOS muss aktualisiert werden. Bitte beziehen Sie sich auf den obigen Link"
            }
          ]
        },
        {
          "SettingName" : "TPMChipPresent",
          "Operator" : "IsEquals",
          "DataType" : "Boolean",
          "Operand" : true,
          "MoreInfoUrl" : "https://bing.com",
          "RemediationStrings" : [
            {
              "Language" : "en_US",
              "Title" : "TPM chip must be enabled.",
              "Description" : "TPM chip must be enabled. Please refer to the link above"
            },
            {
              "Language" : "de_DE",
              "Title" : "TPM-Chip muss aktiviert sein.",
              "Description" : "TPM-Chip muss aktiviert sein. Bitte beziehen Sie sich auf den obigen Link"
            }
          ]
        },
        {
          "SettingName" : "Manufacturer",
          "Operator" : "IsEquals",
          "DataType" : "String",
          "Operand" : "Microsoft Corporation",
          "MoreInfoUrl" : "https://bing.com",
          "RemediationStrings" : [
            {
              "Language" : "en_US",
              "Title" : "Only Microsoft devices are supported.",
              "Description" : "You are not currently using a Microsoft device."
            },
            {
              "Language" : "de_DE",
              "Title" : "Nur Microsoft-Geräte werden unterstützt.",
              "Description" : "Sie verwenden derzeit kein Microsoft-Gerät."
            }
          ]
        }
      ]
    })
  }


  scheduled_actions_for_rule = [
    {
      scheduled_action_configurations = [
        {
          action_type        = "block"
          grace_period_hours = 12
        },
        {
          action_type                  = "notification"
          grace_period_hours           = 24
          notification_template_id     = microsoft365_graph_beta_device_management_device_compliance_notification_template.basic.id
          notification_message_cc_list = [microsoft365_graph_beta_groups_group.acc_test_group_1.id]
        },
        {
          action_type        = "retire"
          grace_period_hours = 48
        },
      ]
    }
  ]

}

// Example 2: Device Health Policy
resource "microsoft365_graph_beta_device_management_windows_device_compliance_policy" "device_health" {
  display_name       = "tf-reg-windows-device-compliance-policy-device-health"
  description        = "tf-reg-windows-device-compliance-policy-device-health"
  role_scope_tag_ids = ["0"]

  # Device Health Settings
  device_health = {
    bit_locker_enabled     = true
    secure_boot_enabled    = true
    code_integrity_enabled = true
  }

  scheduled_actions_for_rule = [
    {
      scheduled_action_configurations = [
        {
          action_type        = "block"
          grace_period_hours = 12
        },
        {
          action_type                  = "notification"
          grace_period_hours           = 24
          notification_template_id     = microsoft365_graph_beta_device_management_device_compliance_notification_template.basic.id
          notification_message_cc_list = [microsoft365_graph_beta_groups_group.acc_test_group_1.id]
        },
        {
          action_type        = "retire"
          grace_period_hours = 48
        },
      ]
    }
  ]

}

// Example 3: Device Properties Policy
resource "microsoft365_graph_beta_device_management_windows_device_compliance_policy" "device_properties" {
  display_name       = "tf-reg-windows-device-compliance-policy-device-properties"
  description        = "tf-reg-windows-device-compliance-policy-device-properties"
  role_scope_tag_ids = ["0"]

  device_properties = {
    os_minimum_version        = "10.0.22631.5768"
    os_maximum_version        = "10.0.26100.9999"
    mobile_os_minimum_version = "10.0.22631.5768"
    mobile_os_maximum_version = "10.0.26100.9999"
    valid_operating_system_build_ranges = [
      {
        description     = "Windows 11 24H2"
        low_os_version  = "10.0.26100.4946"
        high_os_version = "10.0.26100.9999"
      },
      {
        description     = "Windows 11 23H2"
        low_os_version  = "10.0.22631.5768"
        high_os_version = "10.0.22631.9999"
      },

    ]
  }

  scheduled_actions_for_rule = [
    {
      scheduled_action_configurations = [
        {
          action_type        = "block"
          grace_period_hours = 12
        },
        {
          action_type                  = "notification"
          grace_period_hours           = 24
          notification_template_id     = microsoft365_graph_beta_device_management_device_compliance_notification_template.basic.id
          notification_message_cc_list = [microsoft365_graph_beta_groups_group.acc_test_group_1.id]
        },
        {
          action_type        = "retire"
          grace_period_hours = 48
        },
      ]
    }
  ]

}

// Example 4: Microsoft Defender for Endpoint Policy
resource "microsoft365_graph_beta_device_management_windows_device_compliance_policy" "microsoft_defender_for_endpoint" {
  display_name       = "tf-reg-windows-device-compliance-policy-microsoft-defender-for-endpoint"
  description        = "tf-reg-windows-device-compliance-policy-microsoft-defender-for-endpoint"
  role_scope_tag_ids = ["0"]

  # Microsoft Defender for Endpoint Settings
  microsoft_defender_for_endpoint = {
    device_threat_protection_enabled                 = true
    device_threat_protection_required_security_level = "medium"
  }

  # Scheduled Actions for Rule
  scheduled_actions_for_rule = [
    {
      scheduled_action_configurations = [
        {
          action_type        = "block"
          grace_period_hours = 12
        },
        {
          action_type                  = "notification"
          grace_period_hours           = 24
          notification_template_id     = microsoft365_graph_beta_device_management_device_compliance_notification_template.basic.id
          notification_message_cc_list = [microsoft365_graph_beta_groups_group.acc_test_group_1.id]
        },
        {
          action_type        = "retire"
          grace_period_hours = 48
        },
      ]
    }
  ]

}

// Example 5: System Security Policy
resource "microsoft365_graph_beta_device_management_windows_device_compliance_policy" "system_security" {
  display_name       = "tf-reg-windows-device-compliance-policy-system-security"
  description        = "tf-reg-windows-device-compliance-policy-system-security"
  role_scope_tag_ids = ["0"]

  # System Security Settings
  system_security = {
    active_firewall_required                         = true
    anti_spyware_required                            = true
    antivirus_required                               = true
    configuration_manager_compliance_required        = true
    defender_enabled                                 = true
    defender_version                                 = "1.0.0.0"
    device_threat_protection_enabled                 = true
    device_threat_protection_required_security_level = "medium"
    password_block_simple                            = true
    password_minimum_character_set_count             = 4
    password_required                                = true
    password_required_to_unlock_from_idle            = true
    password_required_type                           = "alphanumeric"
    rtp_enabled                                      = true
    signature_out_of_date                            = true
    storage_require_encryption                       = true
    tpm_required                                     = true
  }

  scheduled_actions_for_rule = [
    {
      scheduled_action_configurations = [
        {
          action_type        = "block"
          grace_period_hours = 12
        },
        {
          action_type                  = "notification"
          grace_period_hours           = 24
          notification_template_id     = microsoft365_graph_beta_device_management_device_compliance_notification_template.basic.id
          notification_message_cc_list = [microsoft365_graph_beta_groups_group.acc_test_group_1.id]
        },
        {
          action_type        = "retire"
          grace_period_hours = 48
        },
      ]
    }
  ]

}

// Example 6: WSL Policy with assignments
resource "microsoft365_graph_beta_device_management_windows_device_compliance_policy" "wsl_assignments" {
  display_name       = "tf-reg-windows-device-compliance-policy-wsl-assignments"
  description        = "tf-reg-windows-device-compliance-policy-wsl-assignments"
  role_scope_tag_ids = ["0"]

  wsl_distributions = [
    {
      distribution       = "Ubuntu"
      minimum_os_version = "1.0"
      maximum_os_version = "1.0"
    },
    {
      distribution       = "redhat"
      minimum_os_version = "1.0"
      maximum_os_version = "1.0"
    }
  ]

  scheduled_actions_for_rule = [
    {
      scheduled_action_configurations = [
        {
          action_type        = "block"
          grace_period_hours = 12
        },
        {
          action_type                  = "notification"
          grace_period_hours           = 24
          notification_template_id     = microsoft365_graph_beta_device_management_device_compliance_notification_template.basic.id
          notification_message_cc_list = [microsoft365_graph_beta_groups_group.acc_test_group_1.id]
        },
        {
          action_type        = "retire"
          grace_period_hours = 48
        },
      ]
    }
  ]

  # Assignments
  assignments = [
    # Optional: Assignment targeting all devices with a daily schedule
    {
      type        = "allDevicesAssignmentTarget"
      filter_id   = microsoft365_graph_beta_device_management_assignment_filter.acc_test_assignment_filter_1.id
      filter_type = "include"
    },
    # Optional: Assignment targeting all licensed users with an hourly schedule
    {
      type        = "allLicensedUsersAssignmentTarget"
      filter_id   = microsoft365_graph_beta_device_management_assignment_filter.acc_test_assignment_filter_2.id
      filter_type = "exclude"
    },
    # Optional: Assignment targeting a specific group with include filter
    {
      type        = "groupAssignmentTarget"
      group_id    = microsoft365_graph_beta_groups_group.acc_test_group_1.id
      filter_id   = microsoft365_graph_beta_device_management_assignment_filter.acc_test_assignment_filter_3.id
      filter_type = "include"

    },
    # Optional: Assignment targeting a specific group with exclude filter
    {
      type        = "groupAssignmentTarget"
      group_id    = microsoft365_graph_beta_groups_group.acc_test_group_2.id
      filter_id   = microsoft365_graph_beta_device_management_assignment_filter.acc_test_assignment_filter_4.id
      filter_type = "exclude"
    },
    # Optional: Exclusion group assignments
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_3.id
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_4.id
    },
  ]

}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `display_name` (String) The display name of the device compliance policy
- `scheduled_actions_for_rule` (Attributes List) The list of scheduled action for this rule (see [below for nested schema](#nestedatt--scheduled_actions_for_rule))

### Optional

- `assignments` (Attributes Set) Assignments for the compliance policy. Each assignment specifies the target group and schedule for script execution. (see [below for nested schema](#nestedatt--assignments))
- `custom_compliance_required` (Boolean) Indicates whether custom compliance is required
- `description` (String) Admin provided description of the Device Configuration
- `device_compliance_policy_script` (Attributes) Device compliance policy script for custom compliance. When wsl block is set, this block is computed and should not be set. (see [below for nested schema](#nestedatt--device_compliance_policy_script))
- `device_health` (Attributes) Microsoft Attestation Service evaluation settings. Use these settings to confirm that a device has protective measures enabled at boot time.Learn more here 'https://learn.microsoft.com/en-us/intune/intune-service/protect/compliance-policy-create-windows?WT.mc_id=Portal-Microsoft_Intune_DeviceSettings#device-health' (see [below for nested schema](#nestedatt--device_health))
- `device_properties` (Attributes) Device operating system version requirements and build ranges for compliance evaluation (see [below for nested schema](#nestedatt--device_properties))
- `microsoft_defender_for_endpoint` (Attributes) Microsoft Defender for Endpoint device threat protection settings (see [below for nested schema](#nestedatt--microsoft_defender_for_endpoint))
- `role_scope_tag_ids` (Set of String) Set of scope tag IDs for this Entity instance.
- `system_security` (Attributes) System security settings for device compliance including firewall, antivirus, TPM, and encryption requirements (see [below for nested schema](#nestedatt--system_security))
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `wsl_distributions` (Attributes Set) Windows Subsystem for Linux distributions configuration (see [below for nested schema](#nestedatt--wsl_distributions))

### Read-Only

- `id` (String) The id of the driver.

<a id="nestedatt--scheduled_actions_for_rule"></a>
### Nested Schema for `scheduled_actions_for_rule`

Required:

- `scheduled_action_configurations` (Attributes Set) The list of scheduled action configurations for this compliance policy (see [below for nested schema](#nestedatt--scheduled_actions_for_rule--scheduled_action_configurations))

<a id="nestedatt--scheduled_actions_for_rule--scheduled_action_configurations"></a>
### Nested Schema for `scheduled_actions_for_rule.scheduled_action_configurations`

Required:

- `action_type` (String) What action to take. Possible values are: 'noAction', 'notification', 'block', 'retire', 'wipe', 'removeResourceAccessProfiles', 'pushNotification', 'remoteLock'.
- `grace_period_hours` (Number) Number of hours to wait till the action will be enforced. Value must be between 0 and 365

Optional:

- `notification_message_cc_list` (List of String) A list of group GUIDs to specify who to CC this notification message to
- `notification_template_id` (String) What notification Message template to use



<a id="nestedatt--assignments"></a>
### Nested Schema for `assignments`

Required:

- `type` (String) Type of assignment target. Must be one of: 'allDevicesAssignmentTarget', 'allLicensedUsersAssignmentTarget', 'groupAssignmentTarget', 'exclusionGroupAssignmentTarget'.

Optional:

- `filter_id` (String) ID of the filter to apply to the assignment.
- `filter_type` (String) Type of filter to apply. Must be one of: 'include', 'exclude', or 'none'.
- `group_id` (String) The Entra ID group ID to include or exclude in the assignment. Required when type is 'groupAssignmentTarget' or 'exclusionGroupAssignmentTarget'.


<a id="nestedatt--device_compliance_policy_script"></a>
### Nested Schema for `device_compliance_policy_script`

Optional:

- `device_compliance_script_id` (String) The ID of the device compliance script
- `rules_content` (String) The base64 encoded rules content of the compliance script


<a id="nestedatt--device_health"></a>
### Nested Schema for `device_health`

Optional:

- `bit_locker_enabled` (Boolean) Windows BitLocker Drive Encryption encrypts all data stored on the Windows operating system volume. BitLocker uses the Trusted Platform Module (TPM) to help protect the Windows operating system and user data. It also helps confirm that a computer isn't tampered with, even if its left unattended, lost, or stolen. If the computer is equipped with a compatible TPM, BitLocker uses the TPM to lock the encryption keys that protect the data. As a result, the keys can't be accessed until the TPM verifies the state of the computer. Not configured (default) - This setting isn't evaluated for compliance or non-compliance. Require - The device can protect data that's stored on the drive from unauthorized access when the system is off, or hibernates.
- `code_integrity_enabled` (Boolean) Require code integrity: Code integrity is a feature that validates the integrity of a driver or system file each time it's loaded into memory.Not configured (default) - This setting isn't evaluated for compliance or non-compliance.Require - Require code integrity, which detects if an unsigned driver or system file is being loaded into the kernel. It also detects if a system file is changed by malicious software or run by a user account with administrator privileges.
- `secure_boot_enabled` (Boolean) Require Secure Boot to be enabled on the device:Not configured (default) - This setting isn't evaluated for compliance or non-compliance. Require - The system is forced to boot to a factory trusted state. The core components that are used to boot the machine must have correct cryptographic signatures that are trusted by the organization that manufactured the device. The UEFI firmware verifies the signature before it lets the machine start. If any files are tampered with, which breaks their signature, the system doesn't boot.


<a id="nestedatt--device_properties"></a>
### Nested Schema for `device_properties`

Optional:

- `mobile_os_maximum_version` (String) Enter the maximum allowed version, in the major.minor.build number. When a device is using an OS version later than the version entered, access to organization resources is blocked. The end user is asked to contact their IT administrator. The device can't access organization resources until the rule is changed to allow the OS version.
- `mobile_os_minimum_version` (String) Enter the minimum allowed version, in the major.minor.build number format. When a device has an earlier version that the OS version you enter, it's reported as noncompliant. A link with information on how to upgrade is shown. The end user can choose to upgrade their device. After they upgrade, they can access company resources.
- `os_maximum_version` (String) Maximum OS version:Enter the maximum allowed version, in the major.minor.build.revision number format. To get the correct value, open a command prompt, and type ver. The ver command returns the version in the following format: Microsoft Windows [Version 10.0.17134.1] When a device is using an OS version later than the version entered, access to organization resources is blocked. The end user is asked to contact their IT administrator. The device can't access organization resources until the rule is changed to allow the OS version.
- `os_minimum_version` (String) Minimum OS version. Enter the minimum allowed version in the major.minor.build.revision number format. To get the correct value, open a command prompt, and type ver. The ver command returns the version in the following format: Microsoft Windows [Version 10.0.17134.1] When a device has an earlier version than the OS version you enter, it's reported as noncompliant. A link with information on how to upgrade is shown. The end user can choose to upgrade their device. After they upgrade, they can access company resources.
- `valid_operating_system_build_ranges` (Attributes Set) The valid operating system build ranges on Windows devices (see [below for nested schema](#nestedatt--device_properties--valid_operating_system_build_ranges))

<a id="nestedatt--device_properties--valid_operating_system_build_ranges"></a>
### Nested Schema for `device_properties.valid_operating_system_build_ranges`

Required:

- `high_os_version` (String) The maximum allowed OS version for this build range
- `low_os_version` (String) The minimum allowed OS version for this build range

Optional:

- `description` (String) Description for this valid operating system build range



<a id="nestedatt--microsoft_defender_for_endpoint"></a>
### Nested Schema for `microsoft_defender_for_endpoint`

Optional:

- `device_threat_protection_enabled` (Boolean) Require that devices have enabled device threat protection
- `device_threat_protection_required_security_level` (String) Require Device Threat Protection minimum risk level to report noncompliance. Possible values are: unavailable, secured, low, medium, high, notSet


<a id="nestedatt--system_security"></a>
### Nested Schema for `system_security`

Optional:

- `active_firewall_required` (Boolean) Require active firewall on Windows devices
- `anti_spyware_required` (Boolean) Require any AntiSpyware solution registered with Windows Security Center to be on and monitoring
- `antivirus_required` (Boolean) Require any Antivirus solution registered with Windows Security Center to be on and monitoring
- `configuration_manager_compliance_required` (Boolean) Require device compliance from Configuration Manager: Not configured (default) - Intune doesn't check for any of the Configuration Manager settings for compliance. Require - Require all settings (configuration items) in Configuration Manager to be compliant.
- `defender_enabled` (Boolean) Require Windows Defender Antimalware on Windows devices
- `defender_version` (String) Require Windows Defender Antimalware minimum version on Windows devices
- `password_block_simple` (Boolean) Indicates whether or not to block simple password
- `password_minimum_character_set_count` (Number) The number of character sets required in the password
- `password_required` (Boolean) Require a password to unlock Windows device
- `password_required_to_unlock_from_idle` (Boolean) Require a password to unlock an idle device
- `password_required_type` (String) The required password type. Possible values are: deviceDefault, alphanumeric, numeric
- `rtp_enabled` (Boolean) Require Windows Defender Antimalware Real-Time Protection on Windows devices
- `signature_out_of_date` (Boolean) Require Windows Defender Antimalware Signature to be up to date on Windows devices
- `storage_require_encryption` (Boolean) Require encryption on windows devices
- `tpm_required` (Boolean) Require Trusted Platform Module(TPM) to be present


<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).


<a id="nestedatt--wsl_distributions"></a>
### Nested Schema for `wsl_distributions`

Required:

- `distribution` (String) The name of the WSL distribution
- `maximum_os_version` (String) The maximum OS version for the WSL distribution
- `minimum_os_version` (String) The minimum OS version for the WSL distribution

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.23.0-alpha | Experimental | Initial release |

## Import

Import is supported using the following syntax:

```shell
#!/bin/bash
terraform import microsoft365_graph_beta_device_management_windows_device_compliance_policy.example 00000000-0000-0000-0000-000000000000
```