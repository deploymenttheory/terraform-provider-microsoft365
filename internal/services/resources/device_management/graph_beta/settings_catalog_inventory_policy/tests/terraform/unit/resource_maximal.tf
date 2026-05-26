resource "microsoft365_graph_beta_device_management_settings_catalog_inventory_policy" "maximal" {
  name               = "Test Maximal Inventory Policy - Unit"
  description        = "Comprehensive inventory policy with full settings and all assignment types"
  platforms          = "windows10"
  technologies       = "extensibility"
  role_scope_tag_ids = ["0", "1"]

  configuration_policy = {
    settings = [
      {
        setting_instance = {
          odata_type            = "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance"
          setting_definition_id = "windows_applicationinventory_applicationproperties_applicationkey"
          group_setting_collection_value = [
            {
              children = [
                {
                  odata_type            = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
                  setting_definition_id = "windows_applicationinventory_applicationproperties_applicationkey_appname"
                  choice_setting_value = {
                    children = []
                    value    = "windows_applicationinventory_applicationproperties_applicationkey_appname_24"
                  }
                },
                {
                  odata_type            = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
                  setting_definition_id = "windows_applicationinventory_applicationproperties_applicationkey_appversion"
                  choice_setting_value = {
                    children = []
                    value    = "windows_applicationinventory_applicationproperties_applicationkey_appversion_24"
                  }
                },
                {
                  odata_type            = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
                  setting_definition_id = "windows_applicationinventory_applicationproperties_applicationkey_architectures"
                  choice_setting_value = {
                    children = []
                    value    = "windows_applicationinventory_applicationproperties_applicationkey_architectures_24"
                  }
                },
                {
                  odata_type            = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
                  setting_definition_id = "windows_applicationinventory_applicationproperties_applicationkey_publisher"
                  choice_setting_value = {
                    children = []
                    value    = "windows_applicationinventory_applicationproperties_applicationkey_publisher_24"
                  }
                },
                {
                  odata_type            = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
                  setting_definition_id = "windows_applicationinventory_applicationproperties_applicationkey_installdatetime"
                  choice_setting_value = {
                    children = []
                    value    = "windows_applicationinventory_applicationproperties_applicationkey_installdatetime_24"
                  }
                },
                {
                  odata_type            = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
                  setting_definition_id = "windows_applicationinventory_applicationproperties_applicationkey_installlocation"
                  choice_setting_value = {
                    children = []
                    value    = "windows_applicationinventory_applicationproperties_applicationkey_installlocation_24"
                  }
                }
              ]
            }
          ]
        }
        id = "0"
      }
    ]
  }

  assignments = [
    {
      type        = "groupAssignmentTarget"
      group_id    = "11111111-1111-1111-1111-111111111111"
      filter_type = "include"
      filter_id   = "22222222-2222-2222-2222-222222222222"
    },
    {
      type        = "groupAssignmentTarget"
      group_id    = "33333333-3333-3333-3333-333333333333"
      filter_type = "include"
      filter_id   = "44444444-4444-4444-4444-444444444444"
    },
    {
      type        = "allLicensedUsersAssignmentTarget"
      filter_type = "include"
      filter_id   = "55555555-5555-5555-5555-555555555555"
    },
    {
      type        = "allDevicesAssignmentTarget"
      filter_type = "include"
      filter_id   = "66666666-6666-6666-6666-666666666666"
    },
    {
      type        = "exclusionGroupAssignmentTarget"
      group_id    = "77777777-7777-7777-7777-777777777777"
      filter_type = "include"
      filter_id   = "88888888-8888-8888-8888-888888888888"
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
