resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy_json" "assignments" {
  name               = "Test All Assignment Types Settings Catalog Policy JSON"
  description        = "Test policy with all assignment types configured - JSON format"
  platforms          = "macOS"
  technologies       = ["mdm"]
  role_scope_tag_ids = [microsoft365_graph_beta_device_management_role_scope_tag.acc_test_role_scope_tag_1.id, microsoft365_graph_beta_device_management_role_scope_tag.acc_test_role_scope_tag_2.id]

  settings = jsonencode({
    "settings" = [
      {
        "id" = "0",
        "settingInstance" = {
          "@odata.type"                      = "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance",
          "settingDefinitionId"              = "com.apple.mcx.filevault2_com.apple.mcx.filevault2",
          "settingInstanceTemplateReference" = null,
          "groupSettingCollectionValue" = [
            {
              "settingValueTemplateReference" = null,
              "children" = [
                {
                  "@odata.type"                      = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
                  "settingDefinitionId"              = "com.apple.mcx.filevault2_enable",
                  "settingInstanceTemplateReference" = null,
                  "choiceSettingValue" = {
                    "settingValueTemplateReference" = null,
                    "children"                      = [],
                    "value"                         = "com.apple.mcx.filevault2_enable_1"
                  }
                }
              ]
            }
          ]
        }
      }
    ]
  })

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_1.id
    },
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_2.id
    },
    {
      type = "allLicensedUsersAssignmentTarget"
    },
    {
      type = "allDevicesAssignmentTarget"
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_3.id
    }
  ]

  timeouts = {
    create = "900s"
    read   = "900s"
    update = "900s"
    delete = "900s"
  }
}