data "microsoft365_graph_beta_device_and_app_management_mobile_app" "by_developer" {
  developer = "Microsoft"
}

output "by_developer_count" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app.by_developer.items)
  description = "Number of apps from Microsoft developer"
}
