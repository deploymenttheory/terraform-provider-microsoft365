# Return all service principals
data "microsoft365_graph_beta_applications_service_principal" "all" {
  filter_type = "all"
}

# Return a service principal by display name
data "microsoft365_graph_beta_applications_service_principal" "by_display_name" {
  filter_type  = "display_name"
  filter_value = "Microsoft Graph"
}

# Return a service principal by app id
data "microsoft365_graph_beta_applications_service_principal" "by_app_id" {
  filter_type  = "app_id"
  filter_value = "00000003-0000-0000-c000-000000000000" // Microsoft Graph PowerShell
}

# Return a service principal by odata filter
data "microsoft365_graph_beta_applications_service_principal" "odata_advanced" {
  filter_type  = "odata"
  odata_filter = "startsWith(displayName,'Microsoft')"
  odata_select = "id,appId,displayName,publisherName"
  odata_top    = 10
  odata_skip   = 0
}

# Return a service principal by odata filter with comprehensive options
data "microsoft365_graph_beta_applications_service_principal" "odata_comprehensive" {
  filter_type   = "odata"
  odata_filter  = "startsWith(displayName,'Microsoft')"
  odata_count   = true
  odata_orderby = "displayName"
  odata_search  = "\"displayName:Graph\""
  odata_select  = "id,appId,displayName,publisherName,servicePrincipalType"
  odata_top     = 5
  odata_skip    = 0
}

# Return a service principal by odata filter with matching display name
data "microsoft365_graph_beta_applications_service_principal" "odata_matching_display_name" {
  filter_type  = "odata"
  odata_search = "\"displayName:Intune\""
  odata_count  = true
}

