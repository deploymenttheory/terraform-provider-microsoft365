---
page_title: "microsoft365_graph_beta_windows_updates_applicable_content Data Source - terraform-provider-microsoft365"
subcategory: "Windows Updates"
description: |-
  Retrieves applicable content (driver and firmware updates) for a deployment audience using the /admin/windows/updates/deploymentAudiences/{audienceId}/applicableContent endpoint. This data source shows which updates are applicable to devices in a deployment audience, along with which devices match each update. Supports filtering by catalog entry type, driver class, manufacturer, and specific devices.
---

# microsoft365_graph_beta_windows_updates_applicable_content (Data Source)

Retrieves applicable content (driver, firmware, quality, and feature updates) for a Windows Autopatch deployment audience. This data source shows which updates are applicable to devices in the audience and which devices match each update.

The data source supports filtering by catalog entry type, driver class, manufacturer, and specific devices to help you identify relevant updates for your deployment scenarios.

## Microsoft Documentation

- [Applicable Content Resource Type](https://learn.microsoft.com/en-us/graph/api/resources/windowsupdates-applicablecontent?view=graph-rest-beta)
- [List Applicable Content](https://learn.microsoft.com/en-us/graph/api/windowsupdates-deploymentaudience-list-applicablecontent?view=graph-rest-beta)

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this data source:

**Required:**
- `WindowsUpdates.ReadWrite.All`

**Optional:**
- `None` `[N/A]`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.50.0-alpha | Experimental | Initial release |
| v0.51.0-alpha | Experimental | Added filtering capabilities: `catalog_entry_type`, `driver_class`, `manufacturer`, `device_id`, `include_no_matches`, `odata_filter` |

## Example Usage

### Example 1: Get All Applicable Content

```terraform
# Example: Get all applicable content for a deployment audience

data "microsoft365_graph_beta_windows_updates_applicable_content" "all" {
  audience_id = "12345678-1234-1234-1234-123456789012"
}

# Output the count of applicable content
output "total_applicable_content" {
  value = length(data.microsoft365_graph_beta_windows_updates_applicable_content.all.applicable_content)
}

# Output details of the first applicable content entry
output "first_content" {
  value = length(data.microsoft365_graph_beta_windows_updates_applicable_content.all.applicable_content) > 0 ? {
    catalog_entry_id = data.microsoft365_graph_beta_windows_updates_applicable_content.all.applicable_content[0].catalog_entry_id
    display_name     = data.microsoft365_graph_beta_windows_updates_applicable_content.all.applicable_content[0].catalog_entry.display_name
    matched_devices  = length(data.microsoft365_graph_beta_windows_updates_applicable_content.all.applicable_content[0].matched_devices)
  } : null
}
```

### Example 2: Filter by Driver Updates

```terraform
# Example: Get only driver updates for a deployment audience

data "microsoft365_graph_beta_windows_updates_applicable_content" "drivers" {
  audience_id        = "12345678-1234-1234-1234-123456789012"
  catalog_entry_type = "driver"
}

# Output driver count
output "driver_count" {
  value = length(data.microsoft365_graph_beta_windows_updates_applicable_content.drivers.applicable_content)
}

# Output driver manufacturers
output "driver_manufacturers" {
  value = distinct([
    for content in data.microsoft365_graph_beta_windows_updates_applicable_content.drivers.applicable_content :
    content.catalog_entry.manufacturer
    if content.catalog_entry.manufacturer != null
  ])
}
```

### Example 3: Filter by Display Drivers

```terraform
# Example: Get only display driver updates

data "microsoft365_graph_beta_windows_updates_applicable_content" "display_drivers" {
  audience_id        = "12345678-1234-1234-1234-123456789012"
  catalog_entry_type = "driver"
  driver_class       = "Display"
}

# Output display driver details
output "display_drivers" {
  value = [
    for content in data.microsoft365_graph_beta_windows_updates_applicable_content.display_drivers.applicable_content : {
      display_name = content.catalog_entry.display_name
      manufacturer = content.catalog_entry.manufacturer
      version      = content.catalog_entry.version
      provider     = content.catalog_entry.provider
      device_count = length(content.matched_devices)
    }
  ]
}
```

### Example 4: Filter by Manufacturer

```terraform
# Example: Get updates from a specific manufacturer

data "microsoft365_graph_beta_windows_updates_applicable_content" "intel_updates" {
  audience_id  = "12345678-1234-1234-1234-123456789012"
  manufacturer = "Intel"
}

# Output Intel update details
output "intel_updates" {
  value = [
    for content in data.microsoft365_graph_beta_windows_updates_applicable_content.intel_updates.applicable_content : {
      display_name = content.catalog_entry.display_name
      driver_class = content.catalog_entry.driver_class
      version      = content.catalog_entry.version
      release_date = content.catalog_entry.release_date_time
    }
  ]
}
```

### Example 5: Get Content for Specific Device

```terraform
# Example: Get applicable content for a specific device

data "microsoft365_graph_beta_windows_updates_applicable_content" "device_updates" {
  audience_id = "12345678-1234-1234-1234-123456789012"
  device_id   = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
}

# Output updates applicable to this specific device
output "device_applicable_updates" {
  value = [
    for content in data.microsoft365_graph_beta_windows_updates_applicable_content.device_updates.applicable_content : {
      catalog_entry_id = content.catalog_entry_id
      display_name     = content.catalog_entry.display_name
      driver_class     = content.catalog_entry.driver_class
      manufacturer     = content.catalog_entry.manufacturer
    }
  ]
}
```

### Example 6: Get Only Content with Matched Devices

```terraform
# Example: Get only content that has matched devices

data "microsoft365_graph_beta_windows_updates_applicable_content" "with_matches" {
  audience_id        = "12345678-1234-1234-1234-123456789012"
  include_no_matches = false
}

# Output content with device matches
output "content_with_matches" {
  value = [
    for content in data.microsoft365_graph_beta_windows_updates_applicable_content.with_matches.applicable_content : {
      display_name    = content.catalog_entry.display_name
      matched_devices = length(content.matched_devices)
      device_ids      = [for device in content.matched_devices : device.device_id]
    }
  ]
}

# Output total devices that have applicable content
output "total_devices_with_updates" {
  value = length(distinct(flatten([
    for content in data.microsoft365_graph_beta_windows_updates_applicable_content.with_matches.applicable_content : [
      for device in content.matched_devices : device.device_id
    ]
  ])))
}
```

### Example 7: Filter by Quality Updates

```terraform
# Example: Get only quality/security updates

data "microsoft365_graph_beta_windows_updates_applicable_content" "quality_updates" {
  audience_id        = "12345678-1234-1234-1234-123456789012"
  catalog_entry_type = "quality"
}

# Output quality update details
output "quality_updates" {
  value = [
    for content in data.microsoft365_graph_beta_windows_updates_applicable_content.quality_updates.applicable_content : {
      display_name         = content.catalog_entry.display_name
      release_date         = content.catalog_entry.release_date_time
      deployable_until     = content.catalog_entry.deployable_until_date_time
      matched_device_count = length(content.matched_devices)
    }
  ]
}
```

## Filtering Options

This data source supports multiple filtering options that can be combined to narrow down results:

### By Catalog Entry Type
- `catalog_entry_type` - Filter by update type: `driver`, `quality`, or `feature`
- Useful for separating driver updates from Windows updates

### By Driver Class
- `driver_class` - Filter driver updates by class (e.g., `Display`, `Network`, `Storage`)
- Only applicable when filtering for driver updates
- Common values: `Display`, `Network`, `Storage`, `Audio`, `Bluetooth`, `System`, `HIDClass`

### By Manufacturer
- `manufacturer` - Filter by hardware manufacturer
- Examples: `Intel`, `NVIDIA`, `AMD`, `Microsoft`, `Realtek`, `Qualcomm`
- Works for both driver and firmware updates

### By Device
- `device_id` - Filter to show only content applicable to a specific Azure AD device
- Useful for troubleshooting or checking what updates are available for one device

### By Match Status
- `include_no_matches` - Control whether to include content with no matched devices
- Default: `true` (include all content)
- Set to `false` to only see content that has at least one matched device

### Advanced Filtering
- `odata_filter` - Custom OData filter expression for advanced scenarios
- Applied client-side after retrieving the applicable content
- Example: `catalogEntry/displayName contains 'NVIDIA'`

## Important Notes

1. **Audience ID is Required**: The `audience_id` parameter is mandatory as applicable content is scoped to a specific deployment audience.

2. **Client-Side Filtering**: All filtering (except the initial API query) is performed client-side. This means the data source retrieves all applicable content for the audience and then filters the results based on your criteria.

3. **Matched Devices**: The `matched_devices` list shows which devices in the audience are compatible with each update. This helps identify deployment targets.

4. **Catalog Entry Details**: Each applicable content entry includes detailed information about the update from the catalog, including version, release date, manufacturer, and more.

5. **Empty Results**: If no content matches your filters, the `applicable_content` list will be empty. This is normal and indicates no updates meet your criteria.

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `audience_id` (String) The ID of the deployment audience to query for applicable content. This is required as applicable content is scoped to a specific audience.

### Optional

- `catalog_entry_type` (String) Optional filter to only return content of a specific type. Valid values: `driver`, `quality`, `feature`. When specified, only catalog entries of this type will be returned.
- `device_id` (String) Optional Azure AD device ID to filter results. When specified, only shows applicable content that matches this specific device.
- `driver_class` (String) Optional filter to only return driver updates of a specific class. Examples: `Display`, `Network`, `Storage`, `Audio`, `Bluetooth`. Only applicable when filtering for driver updates.
- `include_no_matches` (Boolean) Whether to include content with no matched devices. Defaults to `true`. Set to `false` to only return content that has at least one matched device.
- `manufacturer` (String) Optional filter to only return updates from a specific manufacturer. Examples: `Intel`, `NVIDIA`, `AMD`, `Microsoft`, `Realtek`.
- `odata_filter` (String) Optional custom OData filter expression for advanced filtering. This is applied client-side after retrieving the applicable content. Example: `catalogEntry/displayName contains 'NVIDIA'`.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `applicable_content` (Attributes List) List of applicable content entries (drivers/firmware) for the audience. (see [below for nested schema](#nestedatt--applicable_content))

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).


<a id="nestedatt--applicable_content"></a>
### Nested Schema for `applicable_content`

Read-Only:

- `catalog_entry` (Attributes) Details about the driver update catalog entry. (see [below for nested schema](#nestedatt--applicable_content--catalog_entry))
- `catalog_entry_id` (String) The ID of the catalog entry for this applicable content.
- `matched_devices` (Attributes List) List of devices that match this driver update. (see [below for nested schema](#nestedatt--applicable_content--matched_devices))

<a id="nestedatt--applicable_content--catalog_entry"></a>
### Nested Schema for `applicable_content.catalog_entry`

Read-Only:

- `deployable_until_date_time` (String) The date and time until which the driver can be deployed, in RFC3339 format.
- `description` (String) Description of the driver update.
- `display_name` (String) The display name of the driver update.
- `driver_class` (String) The class of the driver, e.g., 'Display', 'Network'.
- `id` (String) The unique identifier for the driver update catalog entry.
- `manufacturer` (String) The manufacturer of the driver.
- `provider` (String) The provider of the driver update.
- `release_date_time` (String) The release date and time in RFC3339 format.
- `version` (String) The version of the driver.
- `version_date_time` (String) The version date and time in RFC3339 format.


<a id="nestedatt--applicable_content--matched_devices"></a>
### Nested Schema for `applicable_content.matched_devices`

Read-Only:

- `device_id` (String) The Azure AD device ID.
- `recommended_by` (List of String) List of entities recommending this driver, e.g., ['Microsoft', 'Contoso'].
