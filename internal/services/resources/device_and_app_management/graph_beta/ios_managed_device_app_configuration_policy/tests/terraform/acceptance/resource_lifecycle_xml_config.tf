data "microsoft365_graph_beta_device_and_app_management_mobile_app" "company_portal" {
  display_name    = "Microsoft Intune Company Portal"
  app_type_filter = "iosStoreApp"
}

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_and_app_management_ios_managed_device_app_configuration_policy" "custom_settings" {
  display_name         = "acc-test-ios-managed-device-app-configuration-policy-xml-encoded-${random_string.suffix.result}"
  description          = "Updated description for acceptance testing"
  targeted_mobile_apps = [data.microsoft365_graph_beta_device_and_app_management_mobile_app.company_portal.items[0].id]
  role_scope_tag_ids   = ["0"]

  encoded_setting_xml = <<-XML
    <dict>
      <key>testKey1</key>
      <string>testValue1</string>
      <key>testKey2</key>
      <integer>123</integer>
    </dict>
  XML
}
