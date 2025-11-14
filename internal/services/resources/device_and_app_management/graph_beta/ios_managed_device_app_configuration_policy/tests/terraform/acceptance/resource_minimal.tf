# Data source to find specific iOS app by display name (e.g., "Company Portal")
data "microsoft365_graph_beta_device_and_app_management_mobile_app" "company_portal" {
  filter_type     = "display_name"
  filter_value    = "Microsoft Intune Company Portal"
  app_type_filter = "iosStoreApp" # Only search iOS store apps
}

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_and_app_management_ios_managed_device_app_configuration_policy" "custom_settings" {
  display_name         = "acc-test-ios-managed-device-app-configuration-policy-custom-settings-${random_string.suffix.result}"
  description          = "Updated description for acceptance testing"
  targeted_mobile_apps = [data.microsoft365_graph_beta_device_and_app_management_mobile_app.company_portal.items[0].id]
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
