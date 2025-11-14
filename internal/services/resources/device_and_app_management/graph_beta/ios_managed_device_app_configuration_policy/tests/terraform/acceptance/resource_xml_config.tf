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

resource "microsoft365_graph_beta_device_and_app_management_ios_managed_device_app_configuration_policy" "xml_encoded" {
  display_name         = "acc-test-ios-managed-device-app-configuration-policy-xml-encoded-${random_string.suffix.result}"
  description          = "Updated description for acceptance testing"
  targeted_mobile_apps = [data.microsoft365_graph_beta_device_and_app_management_mobile_app.company_portal.items[0].id]
  role_scope_tag_ids   = ["0"]

  encoded_setting_xml = <<-XML
    <dict>
      <key>metadata</key>
      <dict>
          <key>version</key>
          <string>1.0</string>
          <key>created</key>
          <string>2025-10-14</string>
          <key>author</key>
          <string>System</string>
      </dict>

      <key>items</key>
      <array>
          <dict>
              <key>id</key>
              <string>001</string>
              <key>category</key>
              <string>electronics</string>
              <key>name</key>
              <string>Wireless Mouse</string>
              <key>description</key>
              <string>Ergonomic wireless mouse with USB receiver</string>
              <key>price</key>
              <real>29.99</real>
              <key>stock</key>
              <integer>150</integer>
              <key>specifications</key>
              <dict>
                  <key>battery</key>
                  <string>AA batteries</string>
                  <key>range</key>
                  <string>10 meters</string>
                  <key>color</key>
                  <string>Black</string>
              </dict>
          </dict>

          <dict>
              <key>id</key>
              <string>002</string>
              <key>category</key>
              <string>books</string>
              <key>name</key>
              <string>The Art of Programming</string>
              <key>description</key>
              <string>A comprehensive guide to software development</string>
              <key>price</key>
              <real>49.99</real>
              <key>stock</key>
              <integer>75</integer>
              <key>specifications</key>
              <dict>
                  <key>pages</key>
                  <integer>500</integer>
                  <key>isbn</key>
                  <string>978-1234567890</string>
                  <key>format</key>
                  <string>Hardcover</string>
              </dict>
          </dict>
      </array>
    </dict>
  XML

}
