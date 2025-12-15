data "microsoft365_graph_beta_device_and_app_management_mobile_app" "by_display_name" {
  filter_type  = "display_name"
  filter_value = "Microsoft"
}

output "by_display_name_count" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app.by_display_name.items)
  description = "Number of apps with 'Microsoft' in display name"
}

