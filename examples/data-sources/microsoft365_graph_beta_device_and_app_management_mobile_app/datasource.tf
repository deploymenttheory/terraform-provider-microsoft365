# =============================================================================
# Example 1: Get all mobile apps
# =============================================================================
data "microsoft365_graph_beta_device_and_app_management_mobile_app" "all_apps" {
  filter_type = "all"
  timeouts = {
    read = "10s"
  }
}

output "all_intune_apps" {
  value       = data.microsoft365_graph_beta_device_and_app_management_mobile_app.all_apps.items != null ? data.microsoft365_graph_beta_device_and_app_management_mobile_app.all_apps.items : []
  description = "Complete list of all mobile apps in Intune"
}

# Focused output showing just names and IDs
output "all_intune_apps_summary" {
  value = data.microsoft365_graph_beta_device_and_app_management_mobile_app.all_apps.items != null ? [
    for app in data.microsoft365_graph_beta_device_and_app_management_mobile_app.all_apps.items : {
      id           = app.id
      display_name = app.display_name
      publisher    = app.publisher
      is_assigned  = app.is_assigned
    }
  ] : []
  description = "Summary of all apps with key fields"
}

# =============================================================================
# Example 2: Get mobile app by ID
# =============================================================================
data "microsoft365_graph_beta_device_and_app_management_mobile_app" "by_id" {
  filter_type  = "id"
  filter_value = "b395af0b-910f-40f9-ad74-1cb84406a20f" # Replace with actual app ID

  timeouts = {
    read = "10s"
  }
}

output "app_by_id" {
  value       = try(data.microsoft365_graph_beta_device_and_app_management_mobile_app.by_id.items[0], null)
  description = "Complete details of the specific app"
}

# =============================================================================
# Example 3: Filter apps by display name (case-insensitive contains match)
# =============================================================================
data "microsoft365_graph_beta_device_and_app_management_mobile_app" "by_display_name" {
  filter_type  = "display_name"
  filter_value = "Microsoft" # Finds all apps with "Microsoft" in the name

  timeouts = {
    read = "10s"
  }
}

output "microsoft_apps_by_name" {
  value       = data.microsoft365_graph_beta_device_and_app_management_mobile_app.by_display_name.items != null ? data.microsoft365_graph_beta_device_and_app_management_mobile_app.by_display_name.items : []
  description = "All apps with 'Microsoft' in the display name"
}

output "microsoft_apps_names" {
  value = data.microsoft365_graph_beta_device_and_app_management_mobile_app.by_display_name.items != null ? [
    for app in data.microsoft365_graph_beta_device_and_app_management_mobile_app.by_display_name.items : app.display_name
  ] : []
  description = "List of display names for Microsoft apps"
}

# =============================================================================
# Example 4: Filter apps by publisher name (case-insensitive contains match)
# NEW FEATURE - Filter by publisher
# =============================================================================
data "microsoft365_graph_beta_device_and_app_management_mobile_app" "by_publisher" {
  filter_type  = "publisher_name"
  filter_value = "Adobe" # Finds all apps published by Adobe

  timeouts = {
    read = "10s"
  }
}

output "adobe_apps" {
  value = data.microsoft365_graph_beta_device_and_app_management_mobile_app.by_publisher.items != null ? [
    for app in data.microsoft365_graph_beta_device_and_app_management_mobile_app.by_publisher.items : {
      id           = app.id
      display_name = app.display_name
      publisher    = app.publisher
      developer    = app.developer
      is_assigned  = app.is_assigned
      categories   = app.categories
    }
  ] : []
  description = "All apps from Adobe with key details"
}

# =============================================================================
# Example 5: Filter by app type - Get only Win32 LOB apps
# =============================================================================
data "microsoft365_graph_beta_device_and_app_management_mobile_app" "win32_apps" {
  filter_type     = "all"
  app_type_filter = "win32LobApp"
  
  timeouts = {
    read = "10s"
  }
}

output "all_win32_apps_summary" {
  value = data.microsoft365_graph_beta_device_and_app_management_mobile_app.win32_apps.items != null ? [
    for app in data.microsoft365_graph_beta_device_and_app_management_mobile_app.win32_apps.items : {
      id           = app.id
      display_name = app.display_name
      publisher    = app.publisher
      is_assigned  = app.is_assigned
    }
  ] : []
  description = "Summary of all Win32 LOB apps"
}

# =============================================================================
# Example 6: Combine publisher filter with app type filter
# Get all Microsoft Win32 apps
# =============================================================================
data "microsoft365_graph_beta_device_and_app_management_mobile_app" "microsoft_win32" {
  filter_type     = "publisher_name"
  filter_value    = "Microsoft"
  app_type_filter = "win32LobApp"

  timeouts = {
    read = "10s"
  }
}

output "microsoft_win32_apps" {
  value = data.microsoft365_graph_beta_device_and_app_management_mobile_app.microsoft_win32.items != null ? [
    for app in data.microsoft365_graph_beta_device_and_app_management_mobile_app.microsoft_win32.items : {
      id           = app.id
      display_name = app.display_name
      publisher    = app.publisher
      description  = app.description
      is_assigned  = app.is_assigned
      is_featured  = app.is_featured
    }
  ] : []
  description = "All Microsoft Win32 LOB apps"
}

# =============================================================================
# Example 7: Filter macOS PKG apps
# =============================================================================
data "microsoft365_graph_beta_device_and_app_management_mobile_app" "macos_pkg_apps" {
  filter_type     = "all"
  app_type_filter = "macOSPkgApp"

  timeouts = {
    read = "10s"
  }
}

output "macos_pkg_apps" {
  value = data.microsoft365_graph_beta_device_and_app_management_mobile_app.macos_pkg_apps.items != null ? [
    for app in data.microsoft365_graph_beta_device_and_app_management_mobile_app.macos_pkg_apps.items : {
      id           = app.id
      display_name = app.display_name
      publisher    = app.publisher
      is_assigned  = app.is_assigned
    }
  ] : []
  description = "All macOS PKG apps"
}

# =============================================================================
# Example 8: Advanced OData filter - Apps created after a specific date
# =============================================================================
data "microsoft365_graph_beta_device_and_app_management_mobile_app" "recent_apps" {
  filter_type  = "odata"
  odata_filter = "createdDateTime gt 2024-01-01"
  odata_top    = 10

  timeouts = {
    read = "20s" # OData queries may take longer
  }
}

output "recent_apps" {
  value = try(data.microsoft365_graph_beta_device_and_app_management_mobile_app.recent_apps.items, [])
  description = "Apps created after January 1, 2024 (limited to 10)"
}

# =============================================================================
# Example 9: Advanced OData filter - Using startswith
# =============================================================================
data "microsoft365_graph_beta_device_and_app_management_mobile_app" "apps_starting_with" {
  filter_type  = "odata"
  odata_filter = "startswith(displayName, 'Fire')"

  timeouts = {
    read = "20s"
  }
}

output "apps_starting_with_fire" {
  value = try(data.microsoft365_graph_beta_device_and_app_management_mobile_app.apps_starting_with.items, [])
  description = "Apps with display name starting with 'Fire'"
}

# =============================================================================
# Example 10: Advanced OData filter - Using contains
# =============================================================================
data "microsoft365_graph_beta_device_and_app_management_mobile_app" "apps_containing" {
  filter_type  = "odata"
  odata_filter = "contains(displayName, 'Office')"

  timeouts = {
    read = "20s"
  }
}

output "apps_containing_office" {
  value = try(data.microsoft365_graph_beta_device_and_app_management_mobile_app.apps_containing.items, [])
  description = "Apps with 'Office' anywhere in the display name"
}

# =============================================================================
# Example 11: Get iOS store apps
# =============================================================================
data "microsoft365_graph_beta_device_and_app_management_mobile_app" "ios_store_apps" {
  filter_type     = "all"
  app_type_filter = "iosStoreApp"

  timeouts = {
    read = "10s"
  }
}

output "ios_store_apps" {
  value = data.microsoft365_graph_beta_device_and_app_management_mobile_app.ios_store_apps.items != null ? [
    for app in data.microsoft365_graph_beta_device_and_app_management_mobile_app.ios_store_apps.items : {
      id           = app.id
      display_name = app.display_name
      publisher    = app.publisher
      is_assigned  = app.is_assigned
      categories   = app.categories
    }
  ] : []
  description = "All iOS Store apps"
}

# =============================================================================
# Example 12: Get only assigned apps
# Using local filtering on the results
# =============================================================================
locals {
  assigned_apps = data.microsoft365_graph_beta_device_and_app_management_mobile_app.all_apps.items != null ? [
    for app in data.microsoft365_graph_beta_device_and_app_management_mobile_app.all_apps.items : app
    if app.is_assigned == true
  ] : []
}

output "assigned_apps_only" {
  value       = local.assigned_apps
  description = "Only apps that are assigned to groups"
}

output "assigned_apps_count" {
  value       = length(local.assigned_apps)
  description = "Number of assigned apps"
}

# =============================================================================
# Example 13: Complex local filtering
# Get featured Microsoft apps
# =============================================================================
locals {
  featured_microsoft_apps = data.microsoft365_graph_beta_device_and_app_management_mobile_app.by_publisher.items != null ? [
    for app in data.microsoft365_graph_beta_device_and_app_management_mobile_app.by_publisher.items : app
    if app.is_featured == true && contains(lower(app.publisher), "microsoft")
  ] : []
}

output "featured_microsoft_apps" {
  value = [
    for app in local.featured_microsoft_apps : {
      display_name = app.display_name
      publisher    = app.publisher
      is_assigned  = app.is_assigned
      categories   = app.categories
    }
  ]
  description = "Featured apps from Microsoft"
}
