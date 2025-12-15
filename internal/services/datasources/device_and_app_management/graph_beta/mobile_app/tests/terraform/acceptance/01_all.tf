data "microsoft365_graph_beta_device_and_app_management_mobile_app" "all" {
  filter_type = "all"
}

output "all_apps_count" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app.all.items)
  description = "Total number of mobile apps"
}

output "first_app_id" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app.all.items) > 0 ? data.microsoft365_graph_beta_device_and_app_management_mobile_app.all.items[0].id : null
  description = "ID of the first app"
}

output "first_app_display_name" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app.all.items) > 0 ? data.microsoft365_graph_beta_device_and_app_management_mobile_app.all.items[0].display_name : null
  description = "Display name of the first app"
}

output "first_app_publisher" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app.all.items) > 0 ? data.microsoft365_graph_beta_device_and_app_management_mobile_app.all.items[0].publisher : null
  description = "Publisher of the first app"
}

