data "microsoft365_graph_beta_applications_service_principal" "odata_comprehensive" {
  filter_type   = "odata"
  odata_filter  = "preferredSingleSignOnMode ne 'notSupported'"
  odata_count   = true
  odata_orderby = "displayName"
  odata_search  = "\"displayName:intune\""
  odata_select  = "id,appId,displayName,publisherName,servicePrincipalType"
  odata_top     = 5
  odata_skip    = 0
}