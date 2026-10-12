resource "microsoft365_graph_beta_device_and_app_management_mobile_app_dependency" "test" {
  source_id       = "invalid"
  target_id       = "22222222-2222-2222-2222-222222222222"
  dependency_type = "detect"
}
