data "microsoft365_graph_beta_device_and_app_management_mobile_app" "by_publisher" {
  filter_type  = "publisher_name"
  filter_value = "Microsoft"
}

output "by_publisher_count" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app.by_publisher.items)
  description = "Number of apps from Microsoft publisher"
}

output "by_publisher_first_app" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app.by_publisher.items) > 0 ? data.microsoft365_graph_beta_device_and_app_management_mobile_app.by_publisher.items[0].display_name : null
  description = "First app from Microsoft"
}

