---
page_title: "microsoft365_graph_beta_device_management_endpoint_privilege_management_json Resource - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Manages Endpoint Privilege Management policies using the /deviceManagement/configurationPolicies endpoint. EPM policies control elevation settings and rules for Windows devices, allowing administrators to grant temporary administrative privileges to standard users for specific applications or processes without compromising overall security posture.
---

# microsoft365_graph_beta_device_management_endpoint_privilege_management_json (Resource)

Manages Endpoint Privilege Management policies using the `/deviceManagement/configurationPolicies` endpoint. EPM policies control elevation settings and rules for Windows devices, allowing administrators to grant temporary administrative privileges to standard users for specific applications or processes without compromising overall security posture.

## Microsoft Documentation

- [deviceManagementConfigurationPolicy resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-devicemanagementconfigurationpolicy?view=graph-rest-beta)
- [Create deviceManagementConfigurationPolicy](https://learn.microsoft.com/en-us/graph/api/intune-deviceconfigv2-devicemanagementconfigurationpolicy-create?view=graph-rest-beta)

## API Permissions

The following API permissions are required in order to use this resource.

### Microsoft Graph

- **Application**: `DeviceManagementConfiguration.ReadWrite.All`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.14.1-alpha | Experimental | Initial release |
| v0.25.0-alpha | Testing | Updated to use new settings catalog assignment schema |

## Example Usage

```terraform
# epm elevation settings policy example

resource "microsoft365_graph_beta_device_management_endpoint_privilege_management_json" "epm_elevation_settings_policy" {
  name                           = "EPM Base Elevation settings policy"
  description                    = "Elevation settings policy"
  role_scope_tag_ids             = ["0"]
  settings_catalog_template_type = "elevation_settings_policy"

  settings = jsonencode({

    "settings" : [{
      "id" : "0",
      "settingInstance" : {
        "@odata.type" : "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
        "choiceSettingValue" : {
          "children" : [
            {
              "@odata.type" : "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
              "choiceSettingValue" : {
                "children" : [
                  {
                    "@odata.type" : "#microsoft.graph.deviceManagementConfigurationChoiceSettingCollectionInstance",
                    "choiceSettingCollectionValue" : [
                      {
                        "children" : [],
                        "settingValueTemplateReference" : null,
                    "value" : "device_vendor_msft_policy_privilegemanagement_elevationclientsettings_defaultelevationresponse_validation_0" }],
                    "settingDefinitionId" : "device_vendor_msft_policy_privilegemanagement_elevationclientsettings_defaultelevationresponse_validation",
                    "settingInstanceTemplateReference" : null
                  }
                ],
                "settingValueTemplateReference" : null,
                "value" : "device_vendor_msft_policy_elevationclientsettings_defaultelevationresponse_1"
              }, "settingDefinitionId" : "device_vendor_msft_policy_elevationclientsettings_defaultelevationresponse",
              "settingInstanceTemplateReference" : null
            },
            {
              "@odata.type" : "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance", "choiceSettingValue" : {
                "children" : [],
                "settingValueTemplateReference" : null,
                "value" : "device_vendor_msft_policy_elevationclientsettings_allowelevationdetection_1"
              },
              "settingDefinitionId" : "device_vendor_msft_policy_elevationclientsettings_allowelevationdetection",
              "settingInstanceTemplateReference" : null
            },
            {
              "@odata.type" : "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance", "choiceSettingValue" : {
                "children" : [
                  {
                    "@odata.type" : "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance", "choiceSettingValue" : {
                      "children" : [],
                      "settingValueTemplateReference" : null,
                      "value" : "device_vendor_msft_policy_elevationclientsettings_reportingscope_2"
                    },
                    "settingDefinitionId" : "device_vendor_msft_policy_elevationclientsettings_reportingscope",
                    "settingInstanceTemplateReference" : null
                  }
                ],
                "settingValueTemplateReference" : null,
                "value" : "device_vendor_msft_policy_elevationclientsettings_senddata_1"
              },
              "settingDefinitionId" : "device_vendor_msft_policy_elevationclientsettings_senddata",
              "settingInstanceTemplateReference" : null
            }
          ],
          "settingValueTemplateReference" : {
            "settingValueTemplateId" : "a13cc55c-307a-4962-aaec-20b832bf75c7",
            "useTemplateDefault" : false
          },
          "value" : "device_vendor_msft_policy_elevationclientsettings_enableepm_1"
        }, "settingDefinitionId" : "device_vendor_msft_policy_elevationclientsettings_enableepm",
        "settingInstanceTemplateReference" : {
          "settingInstanceTemplateId" : "58a79a4b-ba9b-4923-a7a5-6dc1a9f638a4"
        }
      }
    }]

  })


  assignments = [
    # Optional: Assignment targeting all devices with inlcude filter
    {
      type        = "allDevicesAssignmentTarget"
      filter_id   = "00000000-0000-0000-0000-000000000000"
      filter_type = "include"
    },
    # Optional: Assignment targeting all licensed users with exclude filter
    {
      type        = "allLicensedUsersAssignmentTarget"
      filter_id   = "00000000-0000-0000-0000-000000000000"
      filter_type = "exclude"
    },
    # Optional: Assignment targeting a specific group with include filter
    {
      type        = "groupAssignmentTarget"
      group_id    = "00000000-0000-0000-0000-000000000000"
      filter_id   = "00000000-0000-0000-0000-000000000000"
      filter_type = "include"
    },
    # Optional: Assignment targeting a specific group with exclude filter
    {
      type        = "groupAssignmentTarget"
      group_id    = "00000000-0000-0000-0000-000000000000"
      filter_id   = "00000000-0000-0000-0000-000000000000"
      filter_type = "exclude"
    },
    # Optional: Assignment targeting a specific group with include filter
    {
      type        = "groupAssignmentTarget"
      group_id    = "00000000-0000-0000-0000-000000000000"
      filter_id   = "00000000-0000-0000-0000-000000000000"
      filter_type = "exclude"
    },
    # Optional: Exclusion group assignments
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000000"
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000000"
    },
  ]

  timeouts = {
    create = "1m"
    read   = "1m"
    update = "1m"
    delete = "1m"
  }
}

# epm elevation rules policy example

resource "microsoft365_graph_beta_device_management_endpoint_privilege_management_json" "epm_elevation_rules_policy" {
  name                           = "EPM Elevation rules policy"
  description                    = "Elevation rules policy"
  role_scope_tag_ids             = ["0"]
  settings_catalog_template_type = "elevation_rules_policy"

  settings = jsonencode({
    "settings" : [
      {
        "settingInstance" : {
          "settingDefinitionId" : "device_vendor_msft_policy_privilegemanagement_elevationrules_{elevationrulename}",
          "settingInstanceTemplateReference" : {
            "settingInstanceTemplateId" : "ee3d2e5f-6b3d-4cb1-af9b-37b02d3dbae2"
          },
          "@odata.type" : "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance",
          "groupSettingCollectionValue" : [
            {
              "settingValueTemplateReference" : null,
              "children" : [
                {
                  "settingDefinitionId" : "device_vendor_msft_policy_privilegemanagement_elevationrules_{elevationrulename}_appliesto",
                  "choiceSettingValue" : {
                    "value" : "device_vendor_msft_policy_privilegemanagement_elevationrules_{elevationrulename}_allusers",
                    "settingValueTemplateReference" : {
                      "settingValueTemplateId" : "2ec26569-c08f-434c-af3d-a50ac4a1ce26",
                      "useTemplateDefault" : false
                    },
                    "children" : []
                  },
                  "settingInstanceTemplateReference" : {
                    "settingInstanceTemplateId" : "0cde1c42-c701-44b1-94b7-438dd4536128"
                  },
                  "@odata.type" : "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
                },
                {
                  "settingDefinitionId" : "device_vendor_msft_policy_privilegemanagement_elevationrules_{elevationrulename}_filehash",
                  "simpleSettingValue" : {
                    "value" : "d5774b403ae04414c6c8e8eb2bc98fc55b1677684f8cee8a4b1c509e55e3d5c1",
                    "@odata.type" : "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
                    "settingValueTemplateReference" : {
                      "settingValueTemplateId" : "1adcc6f7-9fa4-4ce3-8941-2ce22cf5e404",
                      "useTemplateDefault" : false
                    }
                  },
                  "settingInstanceTemplateReference" : {
                    "settingInstanceTemplateId" : "e4436e2c-1584-4fba-8e38-78737cbbbfdf"
                  },
                  "@odata.type" : "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
                },
                {
                  "settingDefinitionId" : "device_vendor_msft_policy_privilegemanagement_elevationrules_{elevationrulename}_ruletype",
                  "choiceSettingValue" : {
                    "value" : "device_vendor_msft_policy_privilegemanagement_elevationrules_{elevationrulename}_self",
                    "settingValueTemplateReference" : {
                      "settingValueTemplateId" : "cb2ea689-ebc3-42ea-a7a4-c704bb13e3ad",
                      "useTemplateDefault" : false
                    },
                    "children" : [
                      {
                        "settingDefinitionId" : "device_vendor_msft_policy_privilegemanagement_elevationrules_{elevationrulename}_ruletype_validation",
                        "choiceSettingCollectionValue" : [
                          {
                            "value" : "device_vendor_msft_policy_privilegemanagement_elevationrules_{elevationrulename}_ruletype_validation_0",
                            "settingValueTemplateReference" : null,
                            "children" : []
                          },
                          {
                            "value" : "device_vendor_msft_policy_privilegemanagement_elevationrules_{elevationrulename}_ruletype_validation_1",
                            "settingValueTemplateReference" : null,
                            "children" : []
                          }
                        ],
                        "settingInstanceTemplateReference" : null,
                        "@odata.type" : "#microsoft.graph.deviceManagementConfigurationChoiceSettingCollectionInstance"
                      }
                    ]
                  },
                  "settingInstanceTemplateReference" : {
                    "settingInstanceTemplateId" : "bc5a31ac-95b5-4ec6-be1f-50a384bb165f"
                  },
                  "@odata.type" : "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
                },
                {
                  "settingDefinitionId" : "device_vendor_msft_policy_privilegemanagement_elevationrules_{elevationrulename}_childprocessbehavior",
                  "choiceSettingValue" : {
                    "value" : "device_vendor_msft_policy_privilegemanagement_elevationrules_{elevationrulename}_allowrunelevatedrulerequired",
                    "settingValueTemplateReference" : null,
                    "children" : []
                  },
                  "settingInstanceTemplateReference" : null,
                  "@odata.type" : "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
                },
                {
                  "settingDefinitionId" : "device_vendor_msft_policy_privilegemanagement_elevationrules_{elevationrulename}_filename",
                  "simpleSettingValue" : {
                    "value" : "test.exe",
                    "@odata.type" : "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
                    "settingValueTemplateReference" : {
                      "settingValueTemplateId" : "a165327c-f0e5-4c7d-9af1-d856b02191f7",
                      "useTemplateDefault" : false
                    }
                  },
                  "settingInstanceTemplateReference" : {
                    "settingInstanceTemplateId" : "0c1ceb2b-bbd4-46d4-9ba5-9ee7abe1f094"
                  },
                  "@odata.type" : "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
                },
                {
                  "settingDefinitionId" : "device_vendor_msft_policy_privilegemanagement_elevationrules_{elevationrulename}_name",
                  "simpleSettingValue" : {
                    "value" : "test",
                    "@odata.type" : "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
                    "settingValueTemplateReference" : {
                      "settingValueTemplateId" : "03f003e5-43ef-4e7e-bf30-57f00781fdcc",
                      "useTemplateDefault" : false
                    }
                  },
                  "settingInstanceTemplateReference" : {
                    "settingInstanceTemplateId" : "fdabfcf9-afa4-4dbf-a4ef-d5c1549065e1"
                  },
                  "@odata.type" : "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
                },
                {
                  "settingDefinitionId" : "device_vendor_msft_policy_privilegemanagement_elevationrules_{elevationrulename}_filepath",
                  "simpleSettingValue" : {
                    "value" : "c:\\path",
                    "@odata.type" : "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
                    "settingValueTemplateReference" : {
                      "settingValueTemplateId" : "f011bcfc-03cd-4b28-a1f4-305278d7a030",
                      "useTemplateDefault" : false
                    }
                  },
                  "settingInstanceTemplateReference" : {
                    "settingInstanceTemplateId" : "c3b7fda4-db6a-421d-bf04-d485e9d0cfb1"
                  },
                  "@odata.type" : "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
                }
              ]
            }
          ]
        },
        "id" : "0"
      }
    ],

  })


  assignments = [
    # Optional: Assignment targeting all devices with inlcude filter
    {
      type        = "allDevicesAssignmentTarget"
      filter_id   = "00000000-0000-0000-0000-000000000000"
      filter_type = "include"
    },
    # Optional: Assignment targeting all licensed users with exclude filter
    {
      type        = "allLicensedUsersAssignmentTarget"
      filter_id   = "00000000-0000-0000-0000-000000000000"
      filter_type = "exclude"
    },
    # Optional: Assignment targeting a specific group with include filter
    {
      type        = "groupAssignmentTarget"
      group_id    = "00000000-0000-0000-0000-000000000000"
      filter_id   = "00000000-0000-0000-0000-000000000000"
      filter_type = "include"
    },
    # Optional: Assignment targeting a specific group with exclude filter
    {
      type        = "groupAssignmentTarget"
      group_id    = "00000000-0000-0000-0000-000000000000"
      filter_id   = "00000000-0000-0000-0000-000000000000"
      filter_type = "exclude"
    },
    # Optional: Assignment targeting a specific group with include filter
    {
      type        = "groupAssignmentTarget"
      group_id    = "00000000-0000-0000-0000-000000000000"
      filter_id   = "00000000-0000-0000-0000-000000000000"
      filter_type = "exclude"
    },
    # Optional: Exclusion group assignments
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000000"
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000000"
    },
  ]

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

- `name` (String) Policy name
- `settings` (String) Endpoint Privilege Management Policy settings defined as a JSON string. Please provide a valid JSON-encoded settings structure. This can either be extracted from an existing policy using the Intune gui `export JSON` functionality if supported, via a script such as this powershell script. [ExportSettingsCatalogConfigurationById](https://github.com/deploymenttheory/terraform-provider-microsoft365/blob/main/scripts/ExportSettingsCatalogConfigurationById.ps1) or created from scratch. The JSON structure should match the graph schema of the settings catalog. Please look at the terraform documentation for the Endpoint Privilege Management Policy for examples and how to correctly format the HCL.

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
- `settings_catalog_template_type` (String) Defines which Endpoint Privilege Management Policy type with settings catalog setting will be deployed. Options available are `elevation_settings_policy` or `elevation_rules_policy`.

### Optional

- `assignments` (Attributes Set) Assignments for the device configuration. Each assignment specifies the target group and schedule for script execution. Supports group filters. (see [below for nested schema](#nestedatt--assignments))
- `description` (String) Endpoint Privilege Management Policy description
- `role_scope_tag_ids` (Set of String) Set of scope tag IDs for this Settings Catalog template profile.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `created_date_time` (String) Creation date and time of the settings catalog policy
- `id` (String) The unique identifier for this Endpoint Privilege Management Policy
- `is_assigned` (Boolean) Indicates if the policy is assigned to any scope
- `last_modified_date_time` (String) Last modification date and time of the settings catalog policy
- `platforms` (String) Platform type for this Endpoint Privilege Management Policy.Will always be set to `windows10`, as EPM currently only supports windows device types.Defaults to windows10.
- `settings_count` (Number) Number of settings catalog settings with the policy. This will change over time as the resource is updated.
- `technologies` (List of String) Describes a list of technologies this Endpoint Privilege Management Policy with settings catalog setting will be deployed with.Defaults to `mdm`, `endpointPrivilegeManagement`.

<a id="nestedatt--assignments"></a>
### Nested Schema for `assignments`

Required:

- `type` (String) Type of assignment target. Must be one of: 'allDevicesAssignmentTarget', 'allLicensedUsersAssignmentTarget', 'groupAssignmentTarget', 'exclusionGroupAssignmentTarget'.

Optional:

- `filter_id` (String) ID of the filter to apply to the assignment.
- `filter_type` (String) Type of filter to apply. Must be one of: 'include', 'exclude', or 'none'.
- `group_id` (String) The Entra ID group ID to include or exclude in the assignment. Required when type is 'groupAssignmentTarget' or 'exclusionGroupAssignmentTarget'.


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
terraform import microsoft365_graph_beta_device_and_app_management_endpoint_privilege_management.example epm-policy-id
```
