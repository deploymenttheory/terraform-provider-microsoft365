data "microsoft365_graph_beta_applications_application" "odata_filter" {
  odata_query = "displayName eq 'Test Application' and signInAudience eq 'AzureADMyOrg'"
}
