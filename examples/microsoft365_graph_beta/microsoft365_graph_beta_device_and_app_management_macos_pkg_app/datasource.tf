# Example 1: Get all macOS PKG apps
data "microsoft365_graph_beta_device_and_app_management_macos_pkg_app" "all_apps" {
  filter_type = "all"

  timeouts {
    read = "5m"
  }
}

# Example output for all_apps
output "all_macos_apps" {
  value = data.microsoft365_graph_beta_device_and_app_management_macos_pkg_app.all_apps.items
}

# More focused output showing just names and IDs
output "all_macos_apps_summary" {
  value = [
    for app in data.microsoft365_graph_beta_device_and_app_management_macos_pkg_app.all_apps.items : {
      id          = app.id
      name        = app.display_name
      description = app.description
    }
  ]
}

# Example 2: Get a specific macOS PKG app by ID
data "microsoft365_graph_beta_device_and_app_management_macos_pkg_app" "by_id" {
  filter_type  = "id"
  filter_value = "00000000-0000-0000-0000-000000000001" # Replace with actual app ID

  timeouts {
    read = "5m"
  }
}

# Output for by_id
output "macos_app_by_id" {
  value = length(data.microsoft365_graph_beta_device_and_app_management_macos_pkg_app.by_id.items) > 0 ? data.microsoft365_graph_beta_device_and_app_management_macos_pkg_app.by_id.items[0] : null
}

# Example 3: Get macOS PKG apps by display name (partial match)
data "microsoft365_graph_beta_device_and_app_management_macos_pkg_app" "by_name" {
  filter_type  = "display_name"
  filter_value = "Adobe" # This will find all apps with "Adobe" in the name

  timeouts {
    read = "5m"
  }
}

# Output for by_name
output "macos_apps_by_name" {
  value = data.microsoft365_graph_beta_device_and_app_management_macos_pkg_app.by_name.items
}

# Specific outputs for by_name
output "adobe_apps_names" {
  value = [
    for app in data.microsoft365_graph_beta_device_and_app_management_macos_pkg_app.by_name.items : app.display_name
  ]
}

# Example 4: Find a specific app by exact display name (demonstrating how to access a specific result)
data "microsoft365_graph_beta_device_and_app_management_macos_pkg_app" "exact_app" {
  filter_type  = "display_name"
  filter_value = "Microsoft OneDrive" # Replace with exact app name
}

# Output showing how to access a specific app when you expect exactly one result
output "onedrive_app" {
  value = length(data.microsoft365_graph_beta_device_and_app_management_macos_pkg_app.exact_app.items) > 0 ? data.microsoft365_graph_beta_device_and_app_management_macos_pkg_app.exact_app.items[0] : null
}

output "onedrive_description" {
  value = length(data.microsoft365_graph_beta_device_and_app_management_macos_pkg_app.exact_app.items) > 0 ? data.microsoft365_graph_beta_device_and_app_management_macos_pkg_app.exact_app.items[0].description : null
}