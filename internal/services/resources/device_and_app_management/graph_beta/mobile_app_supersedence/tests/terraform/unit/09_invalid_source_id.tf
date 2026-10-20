resource "microsoft365_graph_beta_device_and_app_management_mobile_app_supersedence" "test" {
  source_id         = "not-a-guid"
  target_id         = "22222222-2222-2222-2222-222222222222"
  supersedence_type = "update"
}
