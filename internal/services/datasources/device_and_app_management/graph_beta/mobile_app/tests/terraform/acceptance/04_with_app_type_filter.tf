data "microsoft365_graph_beta_device_and_app_management_mobile_app" "win32_apps" {
  list_all        = true
  app_type_filter = "win32LobApp"
}

output "win32_apps_count" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app.win32_apps.items)
  description = "Number of Win32 LOB apps"
}

output "win32_apps_first" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app.win32_apps.items) > 0 ? data.microsoft365_graph_beta_device_and_app_management_mobile_app.win32_apps.items[0].display_name : null
  description = "First Win32 LOB app"
}
