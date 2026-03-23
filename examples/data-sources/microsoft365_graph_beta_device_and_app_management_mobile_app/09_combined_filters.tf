# Combine publisher filter with app type filter
# Get all Microsoft Win32 apps
data "microsoft365_graph_beta_device_and_app_management_mobile_app" "microsoft_win32" {
  publisher       = "Microsoft"
  app_type_filter = "win32LobApp"

  timeouts = {
    read = "10s"
  }
}

output "microsoft_win32_apps" {
  value = [
    for app in data.microsoft365_graph_beta_device_and_app_management_mobile_app.microsoft_win32.items : {
      id           = app.id
      display_name = app.display_name
      publisher    = app.publisher
      description  = app.description
      is_assigned  = app.is_assigned
      is_featured  = app.is_featured
    }
  ]
  description = "All Microsoft Win32 LOB apps"
}

# Example: Get assigned apps using local filtering
data "microsoft365_graph_beta_device_and_app_management_mobile_app" "all_for_filtering" {
  list_all = true
  timeouts = {
    read = "10s"
  }
}

locals {
  assigned_apps = [
    for app in data.microsoft365_graph_beta_device_and_app_management_mobile_app.all_for_filtering.items : app
    if app.is_assigned == true
  ]
}

output "assigned_apps_only" {
  value       = local.assigned_apps
  description = "Only apps that are assigned to groups"
}

output "assigned_apps_count" {
  value       = length(local.assigned_apps)
  description = "Number of assigned apps"
}
