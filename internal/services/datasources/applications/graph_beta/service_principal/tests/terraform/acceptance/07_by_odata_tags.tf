data "microsoft365_graph_beta_applications_service_principal" "by_odata_tags" {
  odata_query = "appId eq '00000003-0000-0000-c000-000000000000'"
}