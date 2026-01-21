---
page_title: "microsoft365_graph_beta_device_and_app_management_mobile_app Data Source - terraform-provider-microsoft365"
subcategory: "Device and App Management"

description: |-
  Retrieves mobile applications from Microsoft Intune using the /deviceAppManagement/mobileApps endpoint. This data source enables querying all mobile app types including Win32 LOB apps, PKG/DMG apps, store apps, and web apps with advanced filtering capabilities for application discovery and configuration planning.
---

# microsoft365_graph_beta_device_and_app_management_mobile_app (Data Source)

Retrieves mobile applications from Microsoft Intune using the `/deviceAppManagement/mobileApps` endpoint. This data source enables querying all mobile app types including Win32 LOB apps, PKG/DMG apps, store apps, and web apps with advanced filtering capabilities for application discovery and configuration planning.

## Microsoft Documentation

- [mobileApp resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-mobileapp?view=graph-rest-beta)

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this data source:

**Required:**
- `DeviceManagementApps.Read.All`
- `DeviceManagementApps.ReadWrite.All`

**Optional:**
- `None` `[N/A]`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.42.0-alpha | Experimental | Added missing version history |

## Example Usage

```terraform
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
  value       = try(data.microsoft365_graph_beta_device_and_app_management_mobile_app.recent_apps.items, [])
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
  value       = try(data.microsoft365_graph_beta_device_and_app_management_mobile_app.apps_starting_with.items, [])
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
  value       = try(data.microsoft365_graph_beta_device_and_app_management_mobile_app.apps_containing.items, [])
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `filter_type` (String) Type of filter to apply. Valid values are: `all`, `id`, `display_name`, `publisher_name`, `odata`.

### Optional

- `app_type_filter` (String) Optional filter that filters returned apps by the application type. Supported values are: `macOSPkgApp`, `macOSDmgApp`, `macOSLobApp`, `macOSMicrosoftDefenderApp`, `macOSMicrosoftEdgeApp`, `macOSOfficeSuiteApp`, `macOsVppApp`, `macOSWebClip`, `androidForWorkApp`, `androidLobApp`, `androidManagedStoreApp`, `androidManagedStoreWebApp`, `androidStoreApp`, `managedAndroidLobApp`, `managedAndroidStoreApp`, `iosiPadOSWebClip`, `iosLobApp`, `iosStoreApp`, `iosVppApp`, `managedIOSLobApp`, `managedIOSStoreApp`, `windowsAppX`, `windowsMicrosoftEdgeApp`, `windowsMobileMSI`, `windowsPhone81AppX`, `windowsPhone81AppXBundle`, `windowsPhone81StoreApp`, `windowsPhoneXAP`, `windowsStoreApp`, `windowsUniversalAppX`, `windowsWebApp`, `winGetApp`, `webApp`, `microsoftStoreForBusinessApp`, `officeSuiteApp`, `win32CatalogApp`, `win32LobApp`, `managedApp`, `managedMobileLobApp`, `mobileLobApp`.
- `filter_value` (String) Value to filter by. Not required when filter_type is 'all' or 'odata'.
- `odata_filter` (String) OData $filter parameter for filtering results. Only used when filter_type is 'odata'.
- `odata_orderby` (String) OData $orderby parameter to sort results. Only used when filter_type is 'odata'.
- `odata_select` (String) OData $select parameter to specify which fields to include. Only used when filter_type is 'odata'.
- `odata_skip` (Number) OData $skip parameter for pagination. Only used when filter_type is 'odata'.
- `odata_top` (Number) OData $top parameter to limit the number of results. Only used when filter_type is 'odata'.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `items` (Attributes List) The list of mobile apps that match the filter criteria. (see [below for nested schema](#nestedatt--items))

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).


<a id="nestedatt--items"></a>
### Nested Schema for `items`

Read-Only:

- `categories` (List of String) The list of categories for this app.
- `created_date_time` (String) The date and time the app was created.
- `dependent_app_count` (Number) The total number of dependencies the child app has.
- `description` (String) The description of the app.
- `developer` (String) The developer of the app.
- `display_name` (String) The admin provided or imported title of the app.
- `id` (String) Key of the entity.
- `information_url` (String) The more information Url.
- `is_assigned` (Boolean) The value indicating whether the app is assigned to at least one group.
- `is_featured` (Boolean) The value indicating whether the app is marked as featured by the admin.
- `last_modified_date_time` (String) The date and time the app was last modified.
- `notes` (String) Notes for the app.
- `owner` (String) The owner of the app.
- `privacy_information_url` (String) The privacy statement Url.
- `publisher` (String) The publisher of the app.
- `publishing_state` (String) The publishing state for the app. The app cannot be assigned unless the app is published. Possible values are: `notPublished`, `processing`, `published`.
- `role_scope_tag_ids` (List of String) List of scope tag ids for this mobile app.
- `superseded_app_count` (Number) The total number of apps this app is directly or indirectly superseded by.
- `superseding_app_count` (Number) The total number of apps this app directly or indirectly supersedes.
- `upload_state` (Number) The upload state.