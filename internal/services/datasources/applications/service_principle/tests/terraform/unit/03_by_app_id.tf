data "microsoft365_graph_beta_applications_service_principal" "by_app_id" {
  filter_type  = "app_id"
  filter_value = "63e61dc2-f593-4a6f-92b9-92e4d2c03d4f"
}