data "microsoft365_graph_beta_applications_service_principal" "odata_filter" {
  filter_type   = "odata"
  odata_filter  = "startsWith(displayName,'Microsoft')"
  odata_count   = true
  odata_orderby = "displayName"
}