data "microsoft365_graph_beta_applications_application" "odata_comprehensive" {
  odata_query = "tags/any(t:t eq 'MyCustomTag') and signInAudience eq 'AzureADMyOrg'"
}
