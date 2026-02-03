data "microsoft365_graph_beta_applications_service_principal" "by_object_id" {
  app_id = "00000003-0000-0000-c000-000000000000" # Microsoft Graph - using app_id instead since object_id varies by tenant
}
