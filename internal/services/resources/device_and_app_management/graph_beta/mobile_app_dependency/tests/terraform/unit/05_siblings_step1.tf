resource "microsoft365_graph_beta_device_and_app_management_mobile_app_dependency" "test" {
  source_id       = "11111111-1111-1111-1111-111111111111"
  target_id       = "22222222-2222-2222-2222-222222222222"
  dependency_type = "detect"
}

resource "microsoft365_graph_beta_device_and_app_management_mobile_app_dependency" "sibling" {
  source_id       = "11111111-1111-1111-1111-111111111111"
  target_id       = "33333333-3333-3333-3333-333333333333"
  dependency_type = "autoInstall"
}
