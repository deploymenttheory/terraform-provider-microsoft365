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