resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "simple_string" {
  name         = "Test Simple String Setting - Unit"
  description  = "Testing simple string setting type"
  platforms    = "macOS"
  technologies = ["mdm"]
  
  template_reference = {
    template_id = ""
  }
  
  configuration_policy = {
    settings = [
      {
        setting_instance = {
          odata_type            = "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
          setting_definition_id = "com.apple.security.fderecoverykeyescrow_location"
          simple_setting_value = {
            odata_type = "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
            value      = "Personal recovery key location message"
          }
        }
        id = "0"
      }
    ]
  }
}