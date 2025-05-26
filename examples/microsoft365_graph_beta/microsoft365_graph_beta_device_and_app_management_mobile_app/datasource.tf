# Example 1: Get all intune mobile apps
data "microsoft365_graph_beta_device_and_app_management_mobile_app" "all_apps" {
  filter_type = "all"
  timeouts = {
    read = "10s"
  }
}

# Example output for all_apps
output "all_intune_apps" {
  value = data.microsoft365_graph_beta_device_and_app_management_mobile_app.all_apps.items != null ? data.microsoft365_graph_beta_device_and_app_management_mobile_app.all_apps.items : []
}

# More focused output showing just names and IDs
output "all_intune_apps_summary" {
  value = data.microsoft365_graph_beta_device_and_app_management_mobile_app.all_apps.items != null ? [
    for app in data.microsoft365_graph_beta_device_and_app_management_mobile_app.all_apps.items : {
      id          = app.id
      name        = app.display_name
      description = app.description
    }
  ] : []
}

# Example 2: Get all intune mobile apps for a specific app type
data "microsoft365_graph_beta_device_and_app_management_mobile_app" "all_winget_apps" {
  filter_type     = "all"
  app_type_filter = "win32LobApp"
  timeouts = {
    read = "10s"
  }
}

# Example output intune winget apps for all_apps
output "all_intune_winget_apps" {
  value = data.microsoft365_graph_beta_device_and_app_management_mobile_app.all_winget_apps.items != null ? data.microsoft365_graph_beta_device_and_app_management_mobile_app.all_winget_apps.items : []
}

# More focused output showing just names and IDs
output "all_intune_winget_apps_summary" {
  value = data.microsoft365_graph_beta_device_and_app_management_mobile_app.all_winget_apps.items != null ? [
    for app in data.microsoft365_graph_beta_device_and_app_management_mobile_app.all_winget_apps.items : {
      id         = app.id
      name       = app.display_name
      isAssigned = app.is_assigned
    }
  ] : []
}

# Example 3: Get a specific macOS PKG app by ID
data "microsoft365_graph_beta_device_and_app_management_mobile_app" "by_id" {
  filter_type  = "id"
  filter_value = "b395af0b-910f-40f9-ad74-1cb84406a20f" # Replace with actual app ID

  timeouts = {
    read = "10s"
  }
}

# Output for by_id
output "macos_app_by_id" {
  value = try(data.microsoft365_graph_beta_device_and_app_management_mobile_app.by_id.items[0], null)
}

# Example 4: Get all apps by display name (partial match)
data "microsoft365_graph_beta_device_and_app_management_mobile_app" "by_name" {
  filter_type  = "display_name"
  filter_value = "Fire" # This will find all apps with "Fire" in the name

  timeouts = {
    read = "10s"
  }
}

# Output for by_name
# returns the complete list of app objects that match the name filter, with all their properties (id, display_name, description, etc.). This gives you the full data for each app.
output "macos_apps_by_name" {
  value = data.microsoft365_graph_beta_device_and_app_management_mobile_app.by_name.items != null ? data.microsoft365_graph_beta_device_and_app_management_mobile_app.by_name.items : []
}

# Specific outputs for by_name
# returns just a list of the display names (as strings) from those same apps. It extracts only the display_name property from each app in the list.
output "fire_apps_names" {
  value = data.microsoft365_graph_beta_device_and_app_management_mobile_app.by_name.items != null ? [
    for app in data.microsoft365_graph_beta_device_and_app_management_mobile_app.by_name.items : app.display_name
  ] : []
}

# Example 5 filter by odata query where creation date is greater than jan 1st 2025
# Odata queries typically need a longer read window
data "microsoft365_graph_beta_device_and_app_management_mobile_app" "app_date_filter_test" {
  filter_type  = "odata"
  odata_filter = "createdDateTime gt 2025-01-01"
  odata_top    = 5

  timeouts = {
    read = "20s"
  }
}

output "app_created_after_creation_date" {
  value = try(data.microsoft365_graph_beta_device_and_app_management_mobile_app.app_date_filter_test.items, [])
}

# Example 5: filter by odata query where display name contains 'Firefox' and app type is winGet
data "microsoft365_graph_beta_device_and_app_management_mobile_app" "odata_contains_query" {
  filter_type  = "odata"
  odata_filter = "contains(displayName, 'Firefox')"

  timeouts = {
    read = "20s"
  }
}

output "odata_contains_query" {
  value = try(data.microsoft365_graph_beta_device_and_app_management_mobile_app.odata_contains_query.items, [])
}

