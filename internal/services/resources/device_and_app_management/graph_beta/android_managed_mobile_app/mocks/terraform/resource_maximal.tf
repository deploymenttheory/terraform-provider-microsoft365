resource "microsoft365_graph_beta_device_and_app_management_android_managed_mobile_app" "maximal" {
  managed_app_protection_id = "00000000-0000-0000-0000-000000000003"
  version                   = "1.5"
  mobile_app_identifier = {
    package_id = "com.example.complexapp"
  }
}