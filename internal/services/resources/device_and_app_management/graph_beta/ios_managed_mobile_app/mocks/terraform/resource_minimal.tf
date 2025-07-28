resource "microsoft365_graph_beta_device_and_app_management_ios_managed_mobile_app" "minimal" {
  managed_app_protection_id = "00000000-0000-0000-0000-000000000002"
  mobile_app_identifier = {
    bundle_id = "com.example.testapp"
  }
}