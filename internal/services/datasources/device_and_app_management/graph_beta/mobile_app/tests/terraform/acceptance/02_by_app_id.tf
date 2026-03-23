data "microsoft365_graph_beta_device_and_app_management_mobile_app" "all_apps" {
  list_all = true
}

data "microsoft365_graph_beta_device_and_app_management_mobile_app" "by_id" {
  app_id = data.microsoft365_graph_beta_device_and_app_management_mobile_app.all_apps.items[0].id
}

output "by_id_app_name" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app.by_id.items) > 0 ? data.microsoft365_graph_beta_device_and_app_management_mobile_app.by_id.items[0].display_name : null
  description = "Display name of the app retrieved by ID"
}
