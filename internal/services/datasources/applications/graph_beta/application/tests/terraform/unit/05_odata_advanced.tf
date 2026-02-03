data "microsoft365_graph_beta_applications_application" "odata_advanced" {
  odata_query = "appId eq '12345678-1234-1234-1234-123456789012'"
}
