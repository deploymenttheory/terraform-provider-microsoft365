resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy_json" "edge_extensions_macos" {
  name               = "Test Edge Extensions macOS Policy"
  description        = "Test policy for Edge Extensions on macOS with multiple collections"
  platforms          = "macOS"
  technologies       = ["mdm", "appleRemoteManagement"]
  role_scope_tag_ids = ["0"]

  settings = jsonencode({
    settings = [
      {
        id = "0"
        settingInstance = {
          "@odata.type"                        = "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance"
          settingDefinitionId                  = "com.apple.managedclient.preferences_extensioninstallallowlist"
          settingInstanceTemplateReference     = null
          simpleSettingCollectionValue = [
            {
              "@odata.type"                     = "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
              value                             = "odfafepnkmbhccpbejgmiehpchacaeak"
              settingValueTemplateReference     = null
            }
          ]
        }
      },
      {
        id = "1"
        settingInstance = {
          "@odata.type"                        = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
          settingDefinitionId                  = "com.apple.managedclient.preferences_blockexternalextensions"
          settingInstanceTemplateReference     = null
          choiceSettingValue = {
            value                            = "com.apple.managedclient.preferences_blockexternalextensions_true"
            settingValueTemplateReference    = null
            children                         = []
          }
        }
      },
      {
        id = "2"
        settingInstance = {
          "@odata.type"                        = "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance"
          settingDefinitionId                  = "com.apple.managedclient.preferences_extensioninstallforcelist"
          settingInstanceTemplateReference     = null
          simpleSettingCollectionValue = [
            {
              "@odata.type"                     = "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
              value                             = "nkbndigcebkoaejohleckhekfmcecfja"
              settingValueTemplateReference     = null
            },
            {
              "@odata.type"                     = "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
              value                             = "ofefcgjbeghpigppfmkologfjadafddi"
              settingValueTemplateReference     = null
            }
          ]
        }
      },
      {
        id = "3"
        settingInstance = {
          "@odata.type"                        = "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance"
          settingDefinitionId                  = "com.apple.managedclient.preferences_extensioninstallblocklist"
          settingInstanceTemplateReference     = null
          simpleSettingCollectionValue = [
            {
              "@odata.type"                     = "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
              value                             = "*"
              settingValueTemplateReference     = null
            }
          ]
        }
      }
    ]
  })
}
