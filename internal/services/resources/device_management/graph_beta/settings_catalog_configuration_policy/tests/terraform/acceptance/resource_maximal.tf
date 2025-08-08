resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "test" {
  name               = "Test Acceptance Settings Catalog Policy - Updated"
  description        = "Updated description for acceptance testing"
  platforms          = "macOS"
  technologies       = ["mdm", "appleRemoteManagement"]
  role_scope_tag_ids = ["0", "1"]

  configuration_policy = {
    settings = [
      {
        setting_instance = {
          odata_type            = "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
          setting_definition_id = "com.apple.test.sample_string"
          simple_setting_value = {
            odata_type = "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
            value      = "hello"
          }
        }
        id = "0"
      }
    ]
  }
}

