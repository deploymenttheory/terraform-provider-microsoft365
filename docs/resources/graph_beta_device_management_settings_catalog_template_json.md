---
page_title: "microsoft365_graph_beta_device_management_settings_catalog_template_json Resource - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Manages a Settings Catalog policy template in Microsoft Intune for Windows, macOS, Linux, iOS/iPadOS and Android.
---

# microsoft365_graph_beta_device_management_settings_catalog_template_json (Resource)

Manages a Settings Catalog policy template in Microsoft Intune for `Windows`, `macOS`, `Linux`, `iOS/iPadOS` and `Android`.

## Microsoft Documentation

- [deviceManagementConfigurationPolicyTemplate resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-devicemanagementconfigurationpolicytemplate?view=graph-rest-beta)
- [Create deviceManagementConfigurationPolicyTemplate](https://learn.microsoft.com/en-us/graph/api/intune-deviceconfigv2-devicemanagementconfigurationpolicytemplate-create?view=graph-rest-beta)

## API Permissions

The following API permissions are required in order to use this resource.

### Microsoft Graph

- **Application**: `DeviceManagementConfiguration.ReadWrite.All`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.14.1-alpha | Experimental | Initial release |

## Example Usage

```terraform
resource "microsoft365_graph_beta_device_management_settings_catalog_template_json" "windows_anti_virus_defender_update_controls" {
  name                           = "Windows - Defender Update controls"
  description                    = "terraform test for settings catalog templates"
  settings_catalog_template_type = "windows_anti_virus_defender_update_controls"
  role_scope_tag_ids             = ["0"]

  settings = jsonencode({
    "settings" : [
      {
        "id" : "0",
        "settingInstance" : {
          "@odata.type" : "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
          "choiceSettingValue" : {
            "value" : "device_vendor_msft_defender_configuration_engineupdateschannel_6",
            "settingValueTemplateReference" : {
              "settingValueTemplateId" : "afc8df70-7b19-4335-b200-bf4b7e098f67",
              "useTemplateDefault" : false
            },
            "children" : []
          },
          "settingInstanceTemplateReference" : {
            "settingInstanceTemplateId" : "f7e1409d-9c85-4a3f-85a6-ad05cc8ccf13"
          },
          "settingDefinitionId" : "device_vendor_msft_defender_configuration_engineupdateschannel"
        }
      },
      {
        "id" : "1",
        "settingInstance" : {
          "@odata.type" : "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
          "choiceSettingValue" : {
            "value" : "device_vendor_msft_defender_configuration_platformupdateschannel_5",
            "settingValueTemplateReference" : {
              "settingValueTemplateId" : "d3b0d61a-bdc5-4507-84d0-5f2a4a3e11a5",
              "useTemplateDefault" : false
            },
            "children" : []
          },
          "settingInstanceTemplateReference" : {
            "settingInstanceTemplateId" : "e78b3ace-75d0-4aad-b3fa-4f49390d6483"
          },
          "settingDefinitionId" : "device_vendor_msft_defender_configuration_platformupdateschannel"
        }
      },
      {
        "id" : "2",
        "settingInstance" : {
          "@odata.type" : "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
          "choiceSettingValue" : {
            "value" : "device_vendor_msft_defender_configuration_securityintelligenceupdateschannel_4",
            "settingValueTemplateReference" : {
              "settingValueTemplateId" : "41ea06bf-e94a-482a-9aaa-7fd535fb4150",
              "useTemplateDefault" : false
            },
            "children" : []
          },
          "settingInstanceTemplateReference" : {
            "settingInstanceTemplateId" : "ba273649-e186-4377-89d5-87405bc9a87c"
          },
          "settingDefinitionId" : "device_vendor_msft_defender_configuration_securityintelligenceupdateschannel"
        }
      }
    ]
  })

  assignments = {
    all_devices = false
    # all_devices_filter_type = "exclude"
    # all_devices_filter_id   = "11111111-2222-3333-4444-555555555555"

    all_users = false
    # all_users_filter_type = "include"
    # all_users_filter_id   = "11111111-2222-3333-4444-555555555555"

    include_groups = [
      {
        group_id                   = "11111111-2222-3333-4444-555555555555"
        include_groups_filter_type = "include"
        include_groups_filter_id   = "11111111-2222-3333-4444-555555555555"
      },
      {
        group_id                   = "11111111-2222-3333-4444-555555555555"
        include_groups_filter_type = "include"
        include_groups_filter_id   = "11111111-2222-3333-4444-555555555555"
      },
    ]

    exclude_group_ids = [
      "11111111-2222-3333-4444-555555555555",
      "11111111-2222-3333-4444-555555555555",
    ]
  }

  timeouts = {
    create = "1m"
    read   = "1m"
    update = "1m"
    delete = "1m"
  }
}


resource "microsoft365_graph_beta_device_management_settings_catalog_template" "windows_anti_virus_microsoft_defender_antivirus_exclusions" {
  name                           = "Windows - Defender Update anti virus exclusions"
  description                    = "terraform test for settings catalog templates"
  settings_catalog_template_type = "windows_anti_virus_microsoft_defender_antivirus_exclusions"
  role_scope_tag_ids             = ["0"]

  settings = jsonencode({
    "settings" : [
      {
        "settingInstance" : {
          "settingDefinitionId" : "device_vendor_msft_policy_config_defender_excludedextensions",
          "@odata.type" : "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance",
          "simpleSettingCollectionValue" : [
            {
              "settingValueTemplateReference" : null,
              "@odata.type" : "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
              "value" : ".dll"
            },
            {
              "settingValueTemplateReference" : null,
              "@odata.type" : "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
              "value" : ".exe"
            }
          ],
          "settingInstanceTemplateReference" : {
            "settingInstanceTemplateId" : "c203725b-17dc-427b-9470-673a2ce9cd5e"
          }
        },
        "id" : "0"
      },
      {
        "settingInstance" : {
          "settingDefinitionId" : "device_vendor_msft_policy_config_defender_excludedpaths",
          "@odata.type" : "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance",
          "simpleSettingCollectionValue" : [
            {
              "settingValueTemplateReference" : null,
              "@odata.type" : "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
              "value" : "c:\\some\\path\\1"
            },
            {
              "settingValueTemplateReference" : null,
              "@odata.type" : "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
              "value" : "c:\\some\\path\\2"
            }
          ],
          "settingInstanceTemplateReference" : {
            "settingInstanceTemplateId" : "aaf04adc-c639-464f-b4a7-152e784092e8"
          }
        },
        "id" : "1"
      },
      {
        "settingInstance" : {
          "settingDefinitionId" : "device_vendor_msft_policy_config_defender_excludedprocesses",
          "@odata.type" : "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance",
          "simpleSettingCollectionValue" : [
            {
              "settingValueTemplateReference" : null,
              "@odata.type" : "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
              "value" : "process-1"
            },
            {
              "settingValueTemplateReference" : null,
              "@odata.type" : "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
              "value" : "process-2"
            }
          ],
          "settingInstanceTemplateReference" : {
            "settingInstanceTemplateId" : "96b046ed-f138-4250-9ae0-b0772a93d16f"
          }
        },
        "id" : "2"
      }
    ]
  })

  assignments = {
    all_devices = false
    # all_devices_filter_type = "exclude"
    # all_devices_filter_id   = "11111111-2222-3333-4444-555555555555"

    all_users = false
    # all_users_filter_type = "include"
    # all_users_filter_id   = "11111111-2222-3333-4444-555555555555"

    include_groups = [
      {
        group_id                   = "11111111-2222-3333-4444-555555555555"
        include_groups_filter_type = "include"
        include_groups_filter_id   = "11111111-2222-3333-4444-555555555555"
      },
      {
        group_id                   = "11111111-2222-3333-4444-555555555555"
        include_groups_filter_type = "include"
        include_groups_filter_id   = "11111111-2222-3333-4444-555555555555"
      },
    ]

    exclude_group_ids = [
      "11111111-2222-3333-4444-555555555555",
      "11111111-2222-3333-4444-555555555555",
    ]
  }

  timeouts = {
    create = "1m"
    read   = "1m"
    update = "1m"
    delete = "1m"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Settings Catalog Policy template name
- `settings` (String) Settings Catalog Policy template settings defined as a JSON string. Please provide a valid JSON-encoded settings structure. This can either be extracted from an existing policy using the Intune gui `export JSON` functionality if supported, via a script such as this powershell script. [ExportSettingsCatalogTemplateConfigurationById](https://github.com/deploymenttheory/terraform-provider-microsoft365/blob/main/scripts/powershell/Export-IntuneSettingsCatalogTemplateConfigurationById.ps1) or created from scratch. The JSON structure should match the graph schema of the settings catalog. Please look at the terraform documentation for the settings catalog template for examples and how to correctly format the HCL.

A correctly formatted field in the HCL should begin and end like this:
```hcl
settings = jsonencode({
  "settings": [
    {
      "id": "0",
      "settingInstance": {
      }
    }
  ]
})
```

**Note:** Settings must always be provided as an array within the settings field, even when configuring a single setting.This is required because the Microsoft Graph SDK for Go always returns settings in an array format

**Note:** When configuring secret values (identified by @odata.type: "#microsoft.graph.deviceManagementConfigurationSecretSettingValue") ensure the valueState is set to "notEncrypted". The value "encryptedValueToken" is reserved for serverresponses and should not be used when creating or updating settings.

```hcl
settings = jsonencode({
  "settings": [
    {
      "id": "0",
      "settingInstance": {
        "@odata.type": "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
        "settingDefinitionId": "com.apple.loginwindow_autologinpassword",
        "settingInstanceTemplateReference": null,
        "simpleSettingValue": {
          "@odata.type": "#microsoft.graph.deviceManagementConfigurationSecretSettingValue",
          "valueState": "notEncrypted",
          "value": "your_secret_value",
          "settingValueTemplateReference": null
        }
      }
    }
  ]
})
```
- `settings_catalog_template_type` (String) Defines the intune settings catalog template type to be deployed using the settings catalog.

This value will automatically set the correct `platform` , `templateID` , `creationSource` and `technologies` values for the settings catalog policy.This value must correctly correlate to the settings defined in the `settings` field.The available options include templates for various platforms and configurations, such as macOS, Windows, and Linux. Options available are:

`Linux settings catalog templates`

`linux_anti_virus_microsoft_defender_antivirus`: This template allows you to configure Microsoft Defender for Endpoint and deploy Antivirus settings to Linux devices.

`linux_anti_virus_microsoft_defender_antivirus_exclusions`: This template allows you to manage settings for Microsoft Defender Antivirus that define Antivirus exclusions for paths, extensions and processes. Antivirus exclusion are also managed by Microsoft Defender Antivirus policy, which includes identical settings for exclusions. Settings from both templates (Antivirus and Antivirus exclusions) are subject to policy merge, and create a super set of exclusions for applicable devices and users.

`linux_endpoint_detection_and_response`: Endpoint detection and response settings for Linux devices.

`macOS settings catalog templates`

`macOS_anti_virus_microsoft_defender_antivirus`: Microsoft Defender Antivirus is the next-generation protection component of Microsoft Defender for Endpoint on Mac. Next-generation protection brings together machine learning, big-data analysis, in-depth threat resistance research, and cloud infrastructure to protect devices in your enterprise organization.

`macOS_anti_virus_microsoft_defender_antivirus_exclusions`: This template allows you to manage settings for Microsoft Defender Antivirus that define Antivirus exclusions for paths, extensions and processes. Antivirus exclusion are also managed by Microsoft Defender Antivirus policy, which includes identical settings for exclusions. Settings from both templates (Antivirus and Antivirus exclusions) are subject to policy merge, and create a super set of exclusions for applicable devices and users.

`macOS_endpoint_detection_and_response`: Endpoint detection and response settings for macOS devices.

`Security Baselines`

`security_baseline_for_windows_10_and_later_version_24H2`: The Security Baseline for Windows 10 and later represents the recommendations for configuring Windows for security conscious customers using the Microsoft full security stack. This baseline includes relevant MDM settings consistent with the security recommendations outlined in the group policy Windows security baseline. Use this baseline to tailor and adjust Microsoft-recommended policy settings within an MDM environment.

`security_baseline_for_microsoft_defender_for_endpoint_version_24H1`: he Microsoft Defender for Endpoint Security baseline for Windows 10 and newer represents the security best practices for the Microsoft security stack on devices managed by Intune (MDM). Use the baseline to tailor and adjust Microsoft-recommended policy settings.

`security_baseline_for_microsoft_edge_version_128`: The Security Baseline for Microsoft Edge represents the recommendations for configuring Microsoft Edge for security conscious customers using the Microsoft full security stack. This baseline aligns with the security recommendations for Edge security baseline for group policy. Use this baseline to configure and customize Microsoft-recommended policy settings.

`security_baseline_for_windows_365`: Windows 365 Security Baselines are a set of policy templates that you can deploy with Microsoft Intune to configure and enforce security settings for Windows 10, Windows 11, Microsoft Edge, and Microsoft Defender for Endpoint on your Cloud PCs. They are based on security best practices and real-world implementations, and they include versioning features to help you update your policies to the latest release. You can also customize the baselines to meet your specific business needs.

`security_baseline_for_microsoft_365_apps`: The Microsoft 365 Apps for enterprise security baseline provides a starting point for IT admins to evaluate and balance the security benefits with productivity needs of their users. This baseline aligns with the security recommendations for Microsoft 365 Apps for enterprise group policy security baseline. Use this baseline to configure and customize Microsoft-recommended policy settings.

`Windows settings catalog templates`

`windows_account_protection`: Account protection policies help protect user credentials by using technology such as Windows Hello for Business and Credential Guard.

`windows_anti_virus_defender_update_controls`: This template allows you to configure the gradual release rollout of Defender Updates to targeted device groups. Use a ringed approach to test, validate, and rollout updates to devices through release channels. Updates available are platform, engine, security intelligence updates. These policy types have pause, resume, manual rollback commands similar to Windows Update ring policies.

`windows_anti_virus_microsoft_defender_antivirus`: Windows Defender Antivirus is the next-generation protection component of Microsoft Defender for Endpoint. Next-generation protection brings together machine learning, big-data analysis, in-depth threat resistance research, and cloud infrastructure to protect devices in your enterprise organization.

`windows_anti_virus_microsoft_defender_antivirus_exclusions`: This template allows you to manage settings for Microsoft Defender Antivirus that define Antivirus exclusions for paths, extensions and processes. Antivirus exclusion are also managed by Microsoft Defender Antivirus policy, which includes identical settings for exclusions. Settings from both templates (Antivirus and Antivirus exclusions) are subject to policy merge, and create a super set of exclusions for applicable devices and users.

`windows_anti_virus_security_experience`: This template allows you to configure the Windows Security app is used by a number of Windows security features to provide notifications about the health and security of the machine. These include notifications about firewalls, antivirus products, Windows Defender SmartScreen, and others.

`windows_app_control_for_business`: Application control settings for Windows devices.

`windows_attack_surface_reduction_app_and_browser_isolation`:This template allows you to configure the Microsoft Defender Application Guard (Application Guard) to help prevent old and newly emerging attacks to help keep employees productive. Using MSFT's unique hardware isolation approach, their goal is to destroy the playbook that attackers use by making current attack methods obsolete.

`windows_attack_surface_reduction_attack_surface_reduction_rules`: This template allows you to configure the Attack surface reduction rules target behaviors that malware and malicious apps typically use to infect computers, including: Executable files and scripts used in Office apps or web mail that attempt to download or run files Obfuscated or otherwise suspicious scripts Behaviors that apps don't usually initiate during normal day-to-day work

`windows_attack_surface_reduction_app_device_control`:This template allows you to configure the securing removable media, and Microsoft Defender for Endpoint provides multiple monitoring and control features to help prevent threats in unauthorized peripherals from compromising your devices.

`windows_attack_surface_reduction_exploit_protection`: This template allows you to configure the protection against malware that uses exploits to infect devices and spread. Exploit protection consists of a number of mitigations that can be applied to either the operating system or individual apps.

`windows_disk_encryption_bitlocker`: This template allows you to configure the BitLocker Drive Encryption data protection features that integrates with the operating system and addresses the threats of data theft or exposure from lost, stolen, or inappropriately decommissioned computers.

`windows_disk_encryption_personal_data`: This template allows you to configure the Personal Data Encryption feature that encrypts select folders and its contents on deployed devices. Personal Data Encryption utilizes Windows Hello for Business to link data encryption keys with user credentials. This feature can minimize the number of credentials the user has to remember to gain access to content. Users will only be able to access their protected content once they've signed into Windows using Windows Hello for Business.

`windows_endpoint_detection_and_response`: Endpoint detection and response settings for Windows devices.

`windows_firewall_rules`: Firewall rules for Windows devices.

`windows_firewall_rules_config_manager`: Rules-based firewall configuration for Windows devices.

`windows_hyper-v_firewall_rules`: Hyper-V firewall rules for Windows devices.

`windows_local_admin_password_solution_(windows_LAPS)`: Windows Local Administrator Password Solution(Windows LAPS) is a Windows feature that automatically manages and backs up the password of a local administrator account on your Azure Active Directory - joined or Windows Server Active Directory - joined devices.

`windows_local_user_group_membership`: Local user group membership policies help to add, remove, or replace members of local groups on Windows devices..

`Windows Configuration Manager settings catalog templates`

`windows_(config_mgr)_anti_virus_microsoft_defender_antivirus`: Microsoft Defender Antivirus settings for Windows devices managed via Microsoft Configuration Manager.

`windows_(config_mgr)_anti_virus_windows_security_experience`: Security experience settings for Windows devices managed via Microsoft Configuration Manager.

`windows_(config_mgr)_attack_surface_reduction`: Attack surface reduction settings for Windows devices managed via Microsoft Configuration Manager.

`windows_(config_mgr)_endpoint_detection_and_response`: Endpoint detection and response settings for Windows devices managed via Microsoft Configuration Manager.

`windows_(config_mgr)_firewall`: Firewall settings for Windows devices managed via Microsoft Configuration Manager.

`windows_(config_mgr)_firewall_profile`: Profile-specific firewall configuration for Windows devices managed via Microsoft Configuration Manager.

`windows_(config_mgr)_firewall_rules`: Rules-based firewall configuration for Windows devices managed via Microsoft Configuration Manager.

`windows_(config_mgr)_attack_surface_reduction_app_and_browser_isolation`: This template allows you to configure the Microsoft Defender Application Guard (Application Guard) settings for devices managed via Microsoft Configuration Manager to help prevent old and newly emerging attacks through hardware-based isolation.

`windows_(config_mgr)_attack_surface_reduction_attack_surface_reduction_rules`: This template allows you to configure Attack Surface Reduction rules for devices managed via Microsoft Configuration Manager. These rules target behaviors commonly used by malware and malicious apps, including suspicious scripts and unusual app behaviors.

`windows_(config_mgr)_attack_surface_reduction_web_protection`: This template allows you to configure web protection settings for devices managed via Microsoft Configuration Manager, helping to protect your organization from web-based threats and malicious content.

`windows_(config_mgr)_attack_surface_reduction_exploit_protection`: This template allows you to configure exploit protection settings for devices managed via Microsoft Configuration Manager. These settings help protect against malware that uses exploits to infect devices and spread through your network.

### Optional

- `assignments` (Attributes) The assignment configuration for this Windows Settings Catalog profile. (see [below for nested schema](#nestedatt--assignments))
- `description` (String) Settings Catalog Policy template description
- `role_scope_tag_ids` (Set of String) Set of scope tag IDs for this Entity instance.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `created_date_time` (String) Creation date and time of the settings catalog policy template
- `id` (String) The unique identifier for this settings catalog policy template
- `is_assigned` (Boolean) Indicates if the policy template is assigned to any user or device scope
- `last_modified_date_time` (String) Last modification date and time of the settings catalog policy template
- `platforms` (String) Platform type for this settings catalog policy.Can be one of: `none`, `android`, `iOS`, `macOS`, `windows10X`, `windows10`, `linux`,`unknownFutureValue`, `androidEnterprise`, or `aosp`. This is automatically set based on the `settings_catalog_template_type` field.
- `settings_count` (Number) Number of settings catalog settings with the policy template. This will change over time as the resource is updated.
- `technologies` (List of String) Describes a list of technologies this settings catalog setting can be deployed with. Valid values are: `none`,`mdm`, `windows10XManagement`, `configManager`, `intuneManagementExtension`, `thirdParty`, `documentGateway`, `appleRemoteManagement`, `microsoftSense`,`exchangeOnline`, `mobileApplicationManagement`, `linuxMdm`, `enrollment`, `endpointPrivilegeManagement`, `unknownFutureValue`, `windowsOsRecovery`, and `android`. This is automatically set based on the `settings_catalog_template_type` field.

<a id="nestedatt--assignments"></a>
### Nested Schema for `assignments`

Optional:

- `all_devices` (Boolean) Specifies whether this assignment applies to all devices. When set to `true`, the assignment targets all devices in the organization.Can be used in conjuction with `all_devices_filter_type` or `all_devices_filter_id`.Can be used as an alternative to `include_groups`.Can be used in conjuction with `all_users` and `all_users_filter_type` or `all_users_filter_id`.
- `all_devices_filter_id` (String) The ID of the device group filter to apply when `all_devices` is set to `true`. This should be a valid GUID of an existing device group filter.
- `all_devices_filter_type` (String) The filter type for all devices assignment. Valid values are:
- `include`: Apply the assignment to devices that match the filter.
- `exclude`: Do not apply the assignment to devices that match the filter.
- `none`: No filter applied.
- `all_users` (Boolean) Specifies whether this assignment applies to all users. When set to `true`, the assignment targets all licensed users within the organization.Can be used in conjuction with `all_users_filter_type` or `all_users_filter_id`.Can be used as an alternative to `include_groups`.Can be used in conjuction with `all_devices` and `all_devices_filter_type` or `all_devices_filter_id`.
- `all_users_filter_id` (String) The ID of the filter to apply when `all_users` is set to `true`. This should be a valid GUID of an existing filter.
- `all_users_filter_type` (String) The filter type for all users assignment. Valid values are:
- `include`: Apply the assignment to users that match the filter.
- `exclude`: Do not apply the assignment to users that match the filter.
- `none`: No filter applied.
- `exclude_group_ids` (List of String) A list of group IDs to exclude from the assignment. These groups will not receive the assignment, even if they match other inclusion criteria.
- `include_groups` (Attributes Set) A set of entra id group Id's to include in the assignment. Each group can have its own filter type and filter ID. (see [below for nested schema](#nestedatt--assignments--include_groups))

<a id="nestedatt--assignments--include_groups"></a>
### Nested Schema for `assignments.include_groups`

Required:

- `group_id` (String) The entra ID group ID of the group to include in the assignment. This should be a valid GUID of an existing group.

Optional:

- `include_groups_filter_id` (String) The Entra ID Group ID of the filter to apply to the included group. This should be a valid GUID of an existing filter.
- `include_groups_filter_type` (String) The device group filter type for the included group. Valid values are:
- `include`: Apply the assignment to group members that match the filter.
- `exclude`: Do not apply the assignment to group members that match the filter.
- `none`: No filter applied.



<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Import

Import is supported using the following syntax:

```shell
# {resource_id}
terraform import microsoft365_graph_beta_device_and_app_management_settings_catalog_template_json.example settings-catalog-template-id
```

