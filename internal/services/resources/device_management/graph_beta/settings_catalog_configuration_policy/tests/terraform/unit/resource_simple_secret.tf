resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "simple_secret" {
  name         = "Test Simple Secret Setting - Unit"
  description  = "Testing simple secret setting type with secret value"
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
          setting_definition_id = "test.secret.password"
          simple_setting_value = {
            odata_type = "#microsoft.graph.deviceManagementConfigurationSecretSettingValue"
            value      = "test-secret-password-123"
          }
        }
        id = "0"
      }
    ]
  }
}