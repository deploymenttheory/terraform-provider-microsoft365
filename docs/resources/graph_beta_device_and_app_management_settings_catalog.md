---
page_title: "microsoft365_graph_beta_device_and_app_management_settings_catalog Resource - terraform-provider-microsoft365"
subcategory: "Intune Settings Catalog"
description: |-
  Manages a Settings Catalog policy in Microsoft Intune for Windows, macOS, iOS/iPadOS and Android.
---

# microsoft365_graph_beta_device_and_app_management_settings_catalog (Resource)

Manages a Settings Catalog policy in Microsoft Intune for Windows, macOS, iOS/iPadOS and Android.

## Example Usage

```terraform
resource "microsoft365_graph_beta_device_and_app_management_settings_catalog" "test_macOS" {
  name               = "Test Settings Catalog Profile - macOS"
  description        = ""
  platforms          = "macOS" 
  technologies       = ["mdm", "appleRemoteManagement"]
  role_scope_tag_ids = ["0"]

  settings = jsonencode({

    "settingsDetails" : [
      {
        "settingInstance" : {
          "groupSettingCollectionValue" : [
            {
              "settingValueTemplateReference" : null,
              "children" : [
                {
                  "settingInstanceTemplateReference" : null,
                  "@odata.type" : "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
                  "choiceSettingValue" : {
                    "settingValueTemplateReference" : null,
                    "value" : "com.apple.mcx_disableguestaccount_true",
                    "children" : []
                  },
                  "settingDefinitionId" : "com.apple.mcx_disableguestaccount"
                },
                {
                  "settingInstanceTemplateReference" : null,
                  "@odata.type" : "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
                  "choiceSettingValue" : {
                    "settingValueTemplateReference" : null,
                    "value" : "com.apple.mcx_enableguestaccount_true",
                    "children" : []
                  },
                  "settingDefinitionId" : "com.apple.mcx_enableguestaccount"
                }
              ]
            }
          ],
          "settingInstanceTemplateReference" : null,
          "@odata.type" : "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance",
          "settingDefinitionId" : "com.apple.mcx_com.apple.mcx-accounts"
        },
        "id" : "0"
      },
      {
        "settingInstance" : {
          "groupSettingCollectionValue" : [
            {
              "settingValueTemplateReference" : null,
              "children" : [
                {
                  "settingInstanceTemplateReference" : null,
                  "@odata.type" : "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
                  "simpleSettingValue" : {
                    "settingValueTemplateReference" : null,
                    "@odata.type" : "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
                    "value" : "thing"
                  },
                  "settingDefinitionId" : "com.apple.caldav.account_caldavaccountdescription"
                },
                {
                  "settingInstanceTemplateReference" : null,
                  "@odata.type" : "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
                  "simpleSettingValue" : {
                    "settingValueTemplateReference" : null,
                    "@odata.type" : "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
                    "value" : "thing"
                  },
                  "settingDefinitionId" : "com.apple.caldav.account_caldavhostname"
                },
                {
                  "settingInstanceTemplateReference" : null,
                  "@odata.type" : "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
                  "simpleSettingValue" : {
                    "settingValueTemplateReference" : null,
                    "valueState" : "notEncrypted",
                    "@odata.type" : "#microsoft.graph.deviceManagementConfigurationSecretSettingValue",
                    "value" : "test-password"
                  },
                  "settingDefinitionId" : "com.apple.caldav.account_caldavpassword"
                },
                {
                  "settingInstanceTemplateReference" : null,
                  "@odata.type" : "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
                  "simpleSettingValue" : {
                    "settingValueTemplateReference" : null,
                    "@odata.type" : "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue",
                    "value" : 1
                  },
                  "settingDefinitionId" : "com.apple.caldav.account_caldavport"
                },
                {
                  "settingInstanceTemplateReference" : null,
                  "@odata.type" : "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
                  "simpleSettingValue" : {
                    "settingValueTemplateReference" : null,
                    "@odata.type" : "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
                    "value" : "thing"
                  },
                  "settingDefinitionId" : "com.apple.caldav.account_caldavprincipalurl"
                },
                {
                  "settingInstanceTemplateReference" : null,
                  "@odata.type" : "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
                  "choiceSettingValue" : {
                    "settingValueTemplateReference" : null,
                    "value" : "com.apple.caldav.account_caldavusessl_true",
                    "children" : []
                  },
                  "settingDefinitionId" : "com.apple.caldav.account_caldavusessl"
                },
                {
                  "settingInstanceTemplateReference" : null,
                  "@odata.type" : "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
                  "simpleSettingValue" : {
                    "settingValueTemplateReference" : null,
                    "@odata.type" : "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
                    "value" : "{{USERNAME}}"
                  },
                  "settingDefinitionId" : "com.apple.caldav.account_caldavusername"
                }
              ]
            }
          ],
          "settingInstanceTemplateReference" : null,
          "@odata.type" : "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance",
          "settingDefinitionId" : "com.apple.caldav.account_com.apple.caldav.account"
        },
        "id" : "1"
      },
      {
        "settingInstance" : {
          "groupSettingCollectionValue" : [
            {
              "settingValueTemplateReference" : null,
              "children" : [
                {
                  "settingInstanceTemplateReference" : null,
                  "@odata.type" : "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
                  "simpleSettingValue" : {
                    "settingValueTemplateReference" : null,
                    "@odata.type" : "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
                    "value" : "thing"
                  },
                  "settingDefinitionId" : "com.apple.carddav.account_carddavaccountdescription"
                },
                {
                  "settingInstanceTemplateReference" : null,
                  "@odata.type" : "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
                  "simpleSettingValue" : {
                    "settingValueTemplateReference" : null,
                    "@odata.type" : "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
                    "value" : "thing"
                  },
                  "settingDefinitionId" : "com.apple.carddav.account_carddavhostname"
                },
                {
                  "settingInstanceTemplateReference" : null,
                  "@odata.type" : "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
                  "simpleSettingValue" : {
                    "settingValueTemplateReference" : null,
                    "valueState" : "notEncrypted",
                    "@odata.type" : "#microsoft.graph.deviceManagementConfigurationSecretSettingValue",
                    "value" : "e7776185-0499-4e47-bdf5-1b3bc42ba965"
                  },
                  "settingDefinitionId" : "com.apple.carddav.account_carddavpassword"
                },
                {
                  "settingInstanceTemplateReference" : null,
                  "@odata.type" : "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
                  "simpleSettingValue" : {
                    "settingValueTemplateReference" : null,
                    "@odata.type" : "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue",
                    "value" : 1
                  },
                  "settingDefinitionId" : "com.apple.carddav.account_carddavport"
                },
                {
                  "settingInstanceTemplateReference" : null,
                  "@odata.type" : "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
                  "choiceSettingValue" : {
                    "settingValueTemplateReference" : null,
                    "value" : "com.apple.carddav.account_carddavusessl_true",
                    "children" : []
                  },
                  "settingDefinitionId" : "com.apple.carddav.account_carddavusessl"
                },
                {
                  "settingInstanceTemplateReference" : null,
                  "@odata.type" : "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
                  "simpleSettingValue" : {
                    "settingValueTemplateReference" : null,
                    "@odata.type" : "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
                    "value" : "{{USERNAME}}"
                  },
                  "settingDefinitionId" : "com.apple.carddav.account_carddavusername"
                }
              ]
            }
          ],
          "settingInstanceTemplateReference" : null,
          "@odata.type" : "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance",
          "settingDefinitionId" : "com.apple.carddav.account_com.apple.carddav.account"
        },
        "id" : "2"
      },
      {
        "settingInstance" : {
          "groupSettingCollectionValue" : [
            {
              "settingValueTemplateReference" : null,
              "children" : [
                {
                  "settingInstanceTemplateReference" : null,
                  "@odata.type" : "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
                  "simpleSettingValue" : {
                    "settingValueTemplateReference" : null,
                    "@odata.type" : "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
                    "value" : "thing"
                  },
                  "settingDefinitionId" : "com.apple.ldap.account_ldapaccountdescription"
                },
                {
                  "settingInstanceTemplateReference" : null,
                  "@odata.type" : "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
                  "simpleSettingValue" : {
                    "settingValueTemplateReference" : null,
                    "@odata.type" : "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
                    "value" : "thing"
                  },
                  "settingDefinitionId" : "com.apple.ldap.account_ldapaccounthostname"
                },
                {
                  "settingInstanceTemplateReference" : null,
                  "@odata.type" : "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
                  "simpleSettingValue" : {
                    "settingValueTemplateReference" : null,
                    "valueState" : "notEncrypted",
                    "@odata.type" : "#microsoft.graph.deviceManagementConfigurationSecretSettingValue",
                    "value" : "762b8bea-3715-449e-b4cd-abc0cb5e16ad"
                  },
                  "settingDefinitionId" : "com.apple.ldap.account_ldapaccountpassword"
                },
                {
                  "settingInstanceTemplateReference" : null,
                  "@odata.type" : "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
                  "choiceSettingValue" : {
                    "settingValueTemplateReference" : null,
                    "value" : "com.apple.ldap.account_ldapaccountusessl_true",
                    "children" : []
                  },
                  "settingDefinitionId" : "com.apple.ldap.account_ldapaccountusessl"
                },
                {
                  "settingInstanceTemplateReference" : null,
                  "@odata.type" : "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
                  "simpleSettingValue" : {
                    "settingValueTemplateReference" : null,
                    "@odata.type" : "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
                    "value" : "{{USERNAME}}"
                  },
                  "settingDefinitionId" : "com.apple.ldap.account_ldapaccountusername"
                },
                {
                  "groupSettingCollectionValue" : [
                    {
                      "settingValueTemplateReference" : null,
                      "children" : [
                        {
                          "settingInstanceTemplateReference" : null,
                          "@odata.type" : "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
                          "simpleSettingValue" : {
                            "settingValueTemplateReference" : null,
                            "@odata.type" : "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
                            "value" : "thing"
                          },
                          "settingDefinitionId" : "com.apple.ldap.account_ldapsearchsettings_item_ldapsearchsettingdescription"
                        },
                        {
                          "settingInstanceTemplateReference" : null,
                          "@odata.type" : "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
                          "choiceSettingValue" : {
                            "settingValueTemplateReference" : null,
                            "value" : "com.apple.ldap.account_ldapsearchsettings_item_ldapsearchsettingscope_2",
                            "children" : []
                          },
                          "settingDefinitionId" : "com.apple.ldap.account_ldapsearchsettings_item_ldapsearchsettingscope"
                        },
                        {
                          "settingInstanceTemplateReference" : null,
                          "@odata.type" : "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
                          "simpleSettingValue" : {
                            "settingValueTemplateReference" : null,
                            "@odata.type" : "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
                            "value" : "thing"
                          },
                          "settingDefinitionId" : "com.apple.ldap.account_ldapsearchsettings_item_ldapsearchsettingsearchbase"
                        }
                      ]
                    },
                    {
                      "settingValueTemplateReference" : null,
                      "children" : [
                        {
                          "settingInstanceTemplateReference" : null,
                          "@odata.type" : "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
                          "simpleSettingValue" : {
                            "settingValueTemplateReference" : null,
                            "@odata.type" : "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
                            "value" : "thing"
                          },
                          "settingDefinitionId" : "com.apple.ldap.account_ldapsearchsettings_item_ldapsearchsettingdescription"
                        },
                        {
                          "settingInstanceTemplateReference" : null,
                          "@odata.type" : "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
                          "choiceSettingValue" : {
                            "settingValueTemplateReference" : null,
                            "value" : "com.apple.ldap.account_ldapsearchsettings_item_ldapsearchsettingscope_2",
                            "children" : []
                          },
                          "settingDefinitionId" : "com.apple.ldap.account_ldapsearchsettings_item_ldapsearchsettingscope"
                        },
                        {
                          "settingInstanceTemplateReference" : null,
                          "@odata.type" : "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
                          "simpleSettingValue" : {
                            "settingValueTemplateReference" : null,
                            "@odata.type" : "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
                            "value" : "thing"
                          },
                          "settingDefinitionId" : "com.apple.ldap.account_ldapsearchsettings_item_ldapsearchsettingsearchbase"
                        }
                      ]
                    }
                  ],
                  "settingInstanceTemplateReference" : null,
                  "@odata.type" : "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance",
                  "settingDefinitionId" : "com.apple.ldap.account_ldapsearchsettings"
                }
              ]
            }
          ],
          "settingInstanceTemplateReference" : null,
          "@odata.type" : "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance",
          "settingDefinitionId" : "com.apple.ldap.account_com.apple.ldap.account"
        },
        "id" : "3"
      }
    ]
  })

  assignments = {
    all_devices = false
    # all_devices_filter_type = "exclude"
    # all_devices_filter_id   = "2d7956fb-e5b5-4fa3-90b2-5bee9bee7883"

    all_users = false
    # all_users_filter_type = "include"
    # all_users_filter_id   = "80f8c0a5-f3ec-4936-bcbc-420dc0ca3665"

    include_groups = [
      {
        group_id                   = "51a96cdd-4b9b-4849-b416-8c94a6d88797"
        include_groups_filter_type = "include"
        include_groups_filter_id   = "80f8c0a5-f3ec-4936-bcbc-420dc0ca3665"
      },
      {
        group_id                   = "b15228f4-9d49-41ed-9b4f-0e7c721fd9c2"
        include_groups_filter_type = "include"
        include_groups_filter_id   = "2d7956fb-e5b5-4fa3-90b2-5bee9bee7883"
      },
    ]

    exclude_group_ids = [
      "b8c661c2-fa9a-4351-af86-adc1729c343f",
      "f6ebd6ff-501e-4b3d-a00b-a2e102c3fa0f",
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

- `name` (String) Policy name
- `settings` (String) Settings Catalog Policy settings defined as a valid JSON string. Provide JSON-encoded settings structure. This can either be extracted from an existing policy using the Intune gui export to JSON, via a script such as [this PowerShell script](https://github.com/deploymenttheory/terraform-provider-microsoft365/blob/main/scripts/GetSettingsCatalogConfigurationById.ps1) or created from scratch. The JSON structure should match the graph schema of the settings catalog. Please look at the terraform documentation for the settings catalog for examples and how to correctly format the HCL.

A correctly formatted field in the HCL should begin and end like this:
```hcl
settings = jsonencode({
  "settingsDetails": [
    {
        # ... settings configuration ...
    }
  ]
})
```

Note: When setting secret values (identified by `@odata.type: "#microsoft.graph.deviceManagementConfigurationSecretSettingValue"`), ensure the `valueState` is set to `"notEncrypted"`. The value `"encryptedValueToken"` is reserved for server responses and should not be used when creating or updating settings.

### Optional

- `assignments` (Attributes) The assignment configuration for this Windows Settings Catalog profile. (see [below for nested schema](#nestedatt--assignments))
- `description` (String) Policy description
- `platforms` (String) Platforms for this policy
- `role_scope_tag_ids` (List of String) List of scope tag IDs for this Windows Settings Catalog profile.
- `technologies` (List of String) Describes a list of technologies this settings catalog setting can be deployed with. Defaults to 'mdm'.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `created_date_time` (String) Creation date and time of the settings catalog policy
- `id` (String) The unique identifier for this policy
- `is_assigned` (Boolean) Indicates if the policy is assigned to any scope
- `last_modified_date_time` (String) Last modification date and time of the settings catalog policy
- `settings_count` (Number) Number of settings catalog settings with the policy. This will change over time as the resource is updated.

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
- `include_groups` (Attributes List) A list of entra id group Id's to include in the assignment. Each group can have its own filter type and filter ID. (see [below for nested schema](#nestedatt--assignments--include_groups))

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
# Using the provider-default project ID, the import ID is:
# {resource_id}
terraform terraform import microsoft365_graph_beta_device_and_app_management_settings_catalog.example settings-catalog-id
```

