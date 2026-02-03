data "microsoft365_graph_beta_applications_service_principal" "odata_advanced" {
  odata_query = "displayName eq 'Microsoft Graph' and accountEnabled eq true"
}