data "microsoft365_graph_beta_applications_service_principal" "odata_advanced" {
  filter_type   = "odata"
  odata_select  = "appId,displayName,publisherName"
  odata_top     = 10
  odata_skip    = 0
}