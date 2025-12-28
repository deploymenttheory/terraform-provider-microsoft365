data "microsoft365_graph_beta_applications_service_principal" "by_display_name" {
  filter_type  = "display_name"
  filter_value = "Microsoft Graph"
}