resource "microsoft365_graph_beta_device_and_app_management_mobile_app_supersedence" "test" {
  source_id         = "11111111-1111-1111-1111-111111111111"
  target_id         = "not-a-guid"
  supersedence_type = "update"
}
