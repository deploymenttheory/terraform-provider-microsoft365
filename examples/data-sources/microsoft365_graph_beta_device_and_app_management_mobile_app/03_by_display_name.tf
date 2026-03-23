# Filter apps by display name (case-insensitive partial match)
# Uses server-side OData filtering for optimal performance
data "microsoft365_graph_beta_device_and_app_management_mobile_app" "by_display_name" {
  display_name = "Microsoft" # Finds all apps with "Microsoft" in the name

  timeouts = {
    read = "10s"
  }
}

output "microsoft_apps_by_name" {
  value       = data.microsoft365_graph_beta_device_and_app_management_mobile_app.by_display_name.items
  description = "All apps with 'Microsoft' in the display name"
}

output "microsoft_apps_names" {
  value = [
    for app in data.microsoft365_graph_beta_device_and_app_management_mobile_app.by_display_name.items : app.display_name
  ]
  description = "List of display names for Microsoft apps"
}
