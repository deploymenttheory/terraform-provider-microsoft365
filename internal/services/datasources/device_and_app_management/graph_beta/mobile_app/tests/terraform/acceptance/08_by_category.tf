data "microsoft365_graph_beta_device_and_app_management_mobile_app" "by_category" {
  category = "Productivity"
}

output "by_category_count" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app.by_category.items)
  description = "Number of apps in Productivity category"
}
