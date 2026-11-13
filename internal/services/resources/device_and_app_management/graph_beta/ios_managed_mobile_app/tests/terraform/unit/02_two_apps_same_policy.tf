resource "microsoft365_graph_beta_device_and_app_management_ios_managed_mobile_app" "first" {
  managed_app_protection_id = "T_00000000-0000-0000-0000-000000000001"
  mobile_app_identifier = {
    bundle_id = "com.example.first"
  }
}

resource "microsoft365_graph_beta_device_and_app_management_ios_managed_mobile_app" "second" {
  managed_app_protection_id = "T_00000000-0000-0000-0000-000000000001"
  mobile_app_identifier = {
    bundle_id = "com.example.second"
  }
}
