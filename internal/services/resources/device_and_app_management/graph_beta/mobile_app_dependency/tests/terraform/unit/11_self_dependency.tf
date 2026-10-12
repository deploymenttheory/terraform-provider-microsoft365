resource "microsoft365_graph_beta_device_and_app_management_mobile_app_dependency" "test" {
  source_id       = "11111111-1111-1111-1111-111111111111"
  target_id       = "11111111-1111-1111-1111-111111111111"
  dependency_type = "detect"
}
