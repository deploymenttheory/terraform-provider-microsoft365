data "microsoft365_graph_beta_device_and_app_management_mobile_app_relationship" "test" {
  source_id   = "11111111-1111-1111-1111-111111111111"
  odata_query = "isof('microsoft.graph.mobileAppSupersedence')"
  target_type = "child"
}
