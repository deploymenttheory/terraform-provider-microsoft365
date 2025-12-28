data "microsoft365_graph_beta_applications_service_principal" "odata_filter" {
  filter_type   = "odata"
  odata_filter  = "preferredSingleSignOnMode ne 'notSupported'"
  odata_count   = true
  odata_orderby = "displayName"
  odata_search  = "\"displayName:intune\""
}