resource "microsoft365_graph_beta_device_and_app_management_ios_managed_device_app_configuration_policy" "custom_settings" {
  display_name         = "unit-test-ios-managed-device-app-configuration-policy-custom-settings"
  description          = "Updated description for acceptance testing"
  targeted_mobile_apps = ["12345678-1234-1234-1234-123456789012"]
  role_scope_tag_ids   = ["0"]

  settings = [
    {
      app_config_key       = "testKey1"
      app_config_key_type  = "stringType"
      app_config_key_value = "testValue1"
    },
    {
      app_config_key       = "testKey2"
      app_config_key_type  = "integerType"
      app_config_key_value = "123"
    }
  ]
}
