resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "choice" {
  name         = "Test Choice Setting - Unit"
  description  = "Testing choice setting type from Edge security settings"
  platforms    = "macOS"
  technologies = ["mdm"]

  template_reference = {
    template_id = ""
  }

  configuration_policy = {
    settings = [
      {
        setting_instance = {
          odata_type            = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
          setting_definition_id = "com.apple.managedclient.preferences_smartscreenenabled"
          choice_setting_value = {
            children = []
            value    = "com.apple.managedclient.preferences_smartscreenenabled_true"
          }
        }
        id = "0"
      }
    ]
  }
}