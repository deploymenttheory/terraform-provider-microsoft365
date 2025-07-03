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

## API Permissions

The following API permissions are required in order to use this data source.

### Microsoft Graph

- **Application**: `DeviceManagementApps.Read.All`, `DeviceManagementApps.ReadWrite.All`

## Example Usage

```terraform
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `filter_type` (String) Type of filter to apply. Valid values are: `all`, `id`, `display_name`, `odata`.

### Optional

- `app_type_filter` (String) Optional filter that filters returned apps by the application type. Supported values are: `macOSPkgApp`, `macOSDmgApp`, `macOSLobApp`, `macOSMicrosoftDefenderApp`, `macOSMicrosoftEdgeApp`, `macOSOfficeSuiteApp`, `macOsVppApp`, `macOSWebClip`, `androidForWorkApp`, `androidLobApp`, `androidManagedStoreApp`, `androidManagedStoreWebApp`, `androidStoreApp`, `managedAndroidLobApp`, `managedAndroidStoreApp`, `iosiPadOSWebClip`, `iosLobApp`, `iosStoreApp`, `iosVppApp`, `managedIOSLobApp`, `managedIOSStoreApp`, `windowsAppX`, `windowsMicrosoftEdgeApp`, `windowsMobileMSI`, `windowsPhone81AppX`, `windowsPhone81AppXBundle`, `windowsPhone81StoreApp`, `windowsPhoneXAP`, `windowsStoreApp`, `windowsUniversalAppX`, `windowsWebApp`, `winGetApp`, `webApp`, `microsoftStoreForBusinessApp`, `officeSuiteApp`, `win32CatalogApp`, `win32LobApp`, `managedApp`, `managedMobileLobApp`, `mobileLobApp`.
- `filter_value` (String) Value to filter by. Not required when filter_type is 'all'.
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