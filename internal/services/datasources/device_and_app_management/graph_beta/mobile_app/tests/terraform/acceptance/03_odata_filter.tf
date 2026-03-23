data "microsoft365_graph_beta_device_and_app_management_mobile_app" "odata_filter" {
  odata_query = "startswith(publisher, 'Microsoft')"
}

output "odata_filter_count" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app.odata_filter.items)
  description = "Number of apps matching OData filter"
}
