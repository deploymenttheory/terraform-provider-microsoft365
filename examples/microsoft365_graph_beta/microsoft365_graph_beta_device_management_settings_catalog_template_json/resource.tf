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