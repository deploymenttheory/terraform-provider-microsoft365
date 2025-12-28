data "microsoft365_graph_beta_applications_service_principal" "odata_search_only" {
  filter_type  = "odata"
  odata_search = "\"displayName:Intune\""
  odata_count  = true
}