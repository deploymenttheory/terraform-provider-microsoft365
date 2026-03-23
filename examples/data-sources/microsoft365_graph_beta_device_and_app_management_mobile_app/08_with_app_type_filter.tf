# Filter by app type - Get only Win32 LOB apps
data "microsoft365_graph_beta_device_and_app_management_mobile_app" "win32_apps" {
  list_all        = true
  app_type_filter = "win32LobApp"

  timeouts = {
    read = "10s"
  }
}

output "all_win32_apps_summary" {
  value = [
    for app in data.microsoft365_graph_beta_device_and_app_management_mobile_app.win32_apps.items : {
      id           = app.id
      display_name = app.display_name
      publisher    = app.publisher
      is_assigned  = app.is_assigned
    }
  ]
  description = "Summary of all Win32 LOB apps"
}

# Example: Get iOS Store apps
data "microsoft365_graph_beta_device_and_app_management_mobile_app" "ios_store_apps" {
  list_all        = true
  app_type_filter = "iosStoreApp"

  timeouts = {
    read = "10s"
  }
}

output "ios_store_apps_count" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app.ios_store_apps.items)
  description = "Number of iOS Store apps"
}

# Example: Get macOS PKG apps
data "microsoft365_graph_beta_device_and_app_management_mobile_app" "macos_pkg_apps" {
  list_all        = true
  app_type_filter = "macOSPkgApp"

  timeouts = {
    read = "10s"
  }
}

output "macos_pkg_apps" {
  value = [
    for app in data.microsoft365_graph_beta_device_and_app_management_mobile_app.macos_pkg_apps.items : {
      id           = app.id
      display_name = app.display_name
      publisher    = app.publisher
      is_assigned  = app.is_assigned
    }
  ]
  description = "All macOS PKG apps"
}
