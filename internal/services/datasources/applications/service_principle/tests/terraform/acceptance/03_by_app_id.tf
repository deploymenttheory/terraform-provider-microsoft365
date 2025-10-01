data "microsoft365_graph_beta_applications_service_principal" "by_app_id" {
  filter_type  = "app_id"
  filter_value = "00000003-0000-0000-c000-000000000000" // Microsoft Graph PowerShell
}