resource "microsoft365_graph_beta_device_and_app_management_mobile_app_supersedence" "test" {
  source_id         = "11111111-1111-1111-1111-111111111111"
  target_id         = "22222222-2222-2222-2222-222222222222"
  supersedence_type = "update"
}

resource "microsoft365_graph_beta_device_and_app_management_mobile_app_supersedence" "sibling" {
  source_id         = "11111111-1111-1111-1111-111111111111"
  target_id         = "33333333-3333-3333-3333-333333333333"
  supersedence_type = "replace"
}
