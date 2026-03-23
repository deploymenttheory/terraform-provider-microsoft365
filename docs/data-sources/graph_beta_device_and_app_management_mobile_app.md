---
page_title: "microsoft365_graph_beta_device_and_app_management_mobile_app Data Source - terraform-provider-microsoft365"
subcategory: "Device and App Management"

description: |-
  Retrieves mobile applications from Microsoft Intune using the /deviceAppManagement/mobileApps endpoint. Supports flexible lookup by app ID, display name, publisher, developer, category, or custom OData queries.
---

# microsoft365_graph_beta_device_and_app_management_mobile_app (Data Source)

Retrieves mobile applications from Microsoft Intune using the `/deviceAppManagement/mobileApps` endpoint. Supports flexible lookup by app ID, display name, publisher, developer, category, or custom OData queries.

## Microsoft Documentation

- [mobileApp resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-mobileapp?view=graph-rest-beta)

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this data source:

**Required:**
- `DeviceManagementApps.Read.All`

**Optional:**
- `None` `[N/A]`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.42.0-alpha | Experimental | Added missing version history |

## Example Usage

### Example 1: List All Mobile Apps

```terraform
# Get all mobile apps from Intune
data "microsoft365_graph_beta_device_and_app_management_mobile_app" "all_apps" {
  list_all = true
  timeouts = {
    read = "10s"
  }
}

output "all_intune_apps_count" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app.all_apps.items)
  description = "Total number of mobile apps in Intune"
}

output "all_intune_apps_summary" {
  value = [
    for app in data.microsoft365_graph_beta_device_and_app_management_mobile_app.all_apps.items : {
      id           = app.id
      display_name = app.display_name
      publisher    = app.publisher
      is_assigned  = app.is_assigned
    }
  ]
  description = "Summary of all apps with key fields"
}
```

### Example 2: Get Mobile App by ID

```terraform
# Get a specific mobile app by its ID
data "microsoft365_graph_beta_device_and_app_management_mobile_app" "by_id" {
  app_id = "b395af0b-910f-40f9-ad74-1cb84406a20f" # Replace with actual app ID

  timeouts = {
    read = "10s"
  }
}

output "app_by_id" {
  value       = try(data.microsoft365_graph_beta_device_and_app_management_mobile_app.by_id.items[0], null)
  description = "Complete details of the specific app"
}

output "app_by_id_name" {
  value       = try(data.microsoft365_graph_beta_device_and_app_management_mobile_app.by_id.items[0].display_name, null)
  description = "Display name of the app"
}
```

### Example 3: Filter by Display Name

```terraform
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
```

### Example 4: Filter by Publisher

```terraform
# Filter apps by publisher name (case-insensitive partial match)
# Uses server-side OData filtering for optimal performance
data "microsoft365_graph_beta_device_and_app_management_mobile_app" "by_publisher" {
  publisher = "Adobe" # Finds all apps published by Adobe

  timeouts = {
    read = "10s"
  }
}

output "adobe_apps" {
  value = [
    for app in data.microsoft365_graph_beta_device_and_app_management_mobile_app.by_publisher.items : {
      id           = app.id
      display_name = app.display_name
      publisher    = app.publisher
      developer    = app.developer
      is_assigned  = app.is_assigned
      categories   = app.categories
    }
  ]
  description = "All apps from Adobe with key details"
}
```

### Example 5: Filter by Developer

```terraform
# Filter apps by developer name (case-insensitive partial match)
# Uses server-side OData filtering for optimal performance
data "microsoft365_graph_beta_device_and_app_management_mobile_app" "by_developer" {
  developer = "Microsoft" # Finds all apps developed by Microsoft

  timeouts = {
    read = "10s"
  }
}

output "microsoft_developed_apps" {
  value = [
    for app in data.microsoft365_graph_beta_device_and_app_management_mobile_app.by_developer.items : {
      id           = app.id
      display_name = app.display_name
      developer    = app.developer
      publisher    = app.publisher
      is_assigned  = app.is_assigned
    }
  ]
  description = "All apps developed by Microsoft"
}
```

### Example 6: Filter by Category

```terraform
# Filter apps by category name (case-insensitive partial match)
# Note: This uses local filtering as categories are expanded relationships
data "microsoft365_graph_beta_device_and_app_management_mobile_app" "by_category" {
  category = "Productivity" # Finds all apps in Productivity category

  timeouts = {
    read = "30s" # Category filtering may take longer as it fetches categories for each app
  }
}

output "productivity_apps" {
  value = [
    for app in data.microsoft365_graph_beta_device_and_app_management_mobile_app.by_category.items : {
      id           = app.id
      display_name = app.display_name
      publisher    = app.publisher
      categories   = app.categories
      is_assigned  = app.is_assigned
    }
  ]
  description = "All apps in the Productivity category"
}
```

### Example 7: Advanced OData Query

```terraform
# Advanced OData filter query
# Uses custom OData expressions for complex filtering
data "microsoft365_graph_beta_device_and_app_management_mobile_app" "odata_custom" {
  odata_query = "startswith(publisher, 'Microsoft') and isAssigned eq true"

  timeouts = {
    read = "20s"
  }
}

output "microsoft_assigned_apps" {
  value = [
    for app in data.microsoft365_graph_beta_device_and_app_management_mobile_app.odata_custom.items : {
      id           = app.id
      display_name = app.display_name
      publisher    = app.publisher
      is_assigned  = app.is_assigned
    }
  ]
  description = "Microsoft apps that are assigned"
}

# Example: Filter by creation date
data "microsoft365_graph_beta_device_and_app_management_mobile_app" "recent_apps" {
  odata_query = "createdDateTime gt 2024-01-01T00:00:00Z"

  timeouts = {
    read = "20s"
  }
}

output "recent_apps_count" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app.recent_apps.items)
  description = "Number of apps created after January 1, 2024"
}
```

### Example 8: Filter by App Type

```terraform
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
```

### Example 9: Combined Filters

```terraform
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `app_id` (String) The unique identifier of the mobile app. Conflicts with other lookup attributes.
- `app_type_filter` (String) Optional filter that filters returned apps by the application type. Supported values are: `macOSPkgApp`, `macOSDmgApp`, `macOSLobApp`, `macOSMicrosoftDefenderApp`, `macOSMicrosoftEdgeApp`, `macOSOfficeSuiteApp`, `macOsVppApp`, `macOSWebClip`, `androidForWorkApp`, `androidLobApp`, `androidManagedStoreApp`, `androidManagedStoreWebApp`, `androidStoreApp`, `managedAndroidLobApp`, `managedAndroidStoreApp`, `iosiPadOSWebClip`, `iosLobApp`, `iosStoreApp`, `iosVppApp`, `managedIOSLobApp`, `managedIOSStoreApp`, `windowsAppX`, `windowsMicrosoftEdgeApp`, `windowsMobileMSI`, `windowsPhone81AppX`, `windowsPhone81AppXBundle`, `windowsPhone81StoreApp`, `windowsPhoneXAP`, `windowsStoreApp`, `windowsUniversalAppX`, `windowsWebApp`, `winGetApp`, `webApp`, `microsoftStoreForBusinessApp`, `officeSuiteApp`, `win32CatalogApp`, `win32LobApp`, `managedApp`, `managedMobileLobApp`, `mobileLobApp`.
- `category` (String) Filter apps by category name (case-insensitive partial match). Conflicts with other lookup attributes.
- `developer` (String) Filter apps by developer name (case-insensitive partial match). Conflicts with other lookup attributes.
- `display_name` (String) Filter apps by display name (case-insensitive partial match). Conflicts with other lookup attributes.
- `list_all` (Boolean) Retrieve all mobile apps. Conflicts with specific lookup attributes.
- `odata_query` (String) Custom OData filter expression for advanced filtering. Example: `startswith(publisher, 'Microsoft') and isAssigned eq true`. Conflicts with other lookup attributes.
- `publisher` (String) Filter apps by publisher name (case-insensitive partial match). Conflicts with other lookup attributes.
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