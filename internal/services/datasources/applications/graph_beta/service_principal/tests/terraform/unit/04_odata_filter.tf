data "microsoft365_graph_beta_applications_service_principal" "odata_filter" {
  odata_query = "preferredSingleSignOnMode ne 'notSupported' and displayName eq 'Microsoft Intune'"
}