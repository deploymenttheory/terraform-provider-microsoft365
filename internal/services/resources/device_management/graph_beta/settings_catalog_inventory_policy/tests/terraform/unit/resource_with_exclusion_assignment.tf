resource "microsoft365_graph_beta_device_management_settings_catalog_inventory_policy" "exclusion_assignment" {
  name         = "Test Exclusion Assignment Inventory Policy - Unit"
  description  = ""
  platforms    = "windows10"
  technologies = "extensibility"

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
