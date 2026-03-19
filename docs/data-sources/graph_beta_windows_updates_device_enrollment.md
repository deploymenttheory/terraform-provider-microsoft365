---
page_title: "microsoft365_graph_beta_windows_updates_device_enrollment Data Source - terraform-provider-microsoft365"
subcategory: "Windows Updates"

description: |-
  Retrieves Windows Autopatch enrollment status for Azure AD devices using the /admin/windows/updates/updatableAssets endpoint. This data source supports multiple lookup methods: by Entra device ID, by device name, list all enrolled devices, or use custom OData queries for advanced filtering.
---

# microsoft365_graph_beta_windows_updates_device_enrollment (Data Source)

Retrieves Windows Autopatch enrollment status for Azure AD devices using the `/admin/windows/updates/updatableAssets` endpoint. This data source supports multiple lookup methods: by Entra device ID, by device name, list all enrolled devices, or use custom OData queries for advanced filtering.

## Microsoft Documentation

- [azureADDevice resource type](https://learn.microsoft.com/en-us/graph/api/resources/windowsupdates-azureaddevice?view=graph-rest-beta)
- [updatableAsset resource type](https://learn.microsoft.com/en-us/graph/api/resources/windowsupdates-updatableasset?view=graph-rest-beta)
- [List updatableAssets](https://learn.microsoft.com/en-us/graph/api/adminwindowsupdates-list-updatableassets?view=graph-rest-beta)
- [Get updatableAsset](https://learn.microsoft.com/en-us/graph/api/windowsupdates-updatableasset-get?view=graph-rest-beta)
- [updateManagementEnrollment resource type](https://learn.microsoft.com/en-us/graph/api/resources/windowsupdates-updatemanagementenrollment?view=graph-rest-beta)

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this data source:

**Required:**
- `WindowsUpdates.ReadWrite.All`

**Optional:**
- `None` `[N/A]`

## Lookup Methods

This data source supports multiple lookup methods (mutually exclusive):

1. **By Entra Device ID** - Look up a single device by its Entra ID object ID
2. **By Device Name** - Look up a single device by resolving its name to an Entra ID
3. **List All** - Retrieve all enrolled devices
4. **OData Filter** - Use custom OData filter expressions for advanced queries

### Optional Filters

When using `list_all` or `odata_filter`, you can additionally filter by:
- `update_category` - Filter devices enrolled in specific update categories (quality, feature, driver)

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.50.0-alpha | Experimental | Initial release |

## Example Usage

### Get Enrollment Status by Entra Device ID

```terraform
# Example: Get enrollment status for a specific device by Entra ID
data "microsoft365_graph_beta_windows_updates_device_enrollment" "by_id" {
  entra_device_id = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
}

# Output the enrollment details
output "device_enrollment_status" {
  value = {
    device_id   = data.microsoft365_graph_beta_windows_updates_device_enrollment.by_id.devices[0].id
    enrollments = data.microsoft365_graph_beta_windows_updates_device_enrollment.by_id.devices[0].enrollments
    errors      = data.microsoft365_graph_beta_windows_updates_device_enrollment.by_id.devices[0].errors
  }
}
```

### Get Enrollment Status by Device Name

```terraform
# Example: Get enrollment status for a device by name
data "microsoft365_graph_beta_windows_updates_device_enrollment" "by_name" {
  device_name = "DESKTOP-ABC123"
}

# Output the enrollment details
output "device_enrollment_by_name" {
  value = {
    device_id   = data.microsoft365_graph_beta_windows_updates_device_enrollment.by_name.devices[0].id
    enrollments = data.microsoft365_graph_beta_windows_updates_device_enrollment.by_name.devices[0].enrollments
  }
}
```

### List All Enrolled Devices

```terraform
# Example: List all enrolled devices
data "microsoft365_graph_beta_windows_updates_device_enrollment" "all_devices" {
  list_all = true
}

# Output count of enrolled devices
output "total_enrolled_devices" {
  value = length(data.microsoft365_graph_beta_windows_updates_device_enrollment.all_devices.devices)
}

# Output devices enrolled in quality updates
output "quality_enrolled_devices" {
  value = [
    for device in data.microsoft365_graph_beta_windows_updates_device_enrollment.all_devices.devices :
    device.id if device.enrollments.quality != null
  ]
}
```

### Filter by Update Category

```terraform
# Example: Filter enrolled devices by update category
data "microsoft365_graph_beta_windows_updates_device_enrollment" "quality_updates" {
  list_all        = true
  update_category = "quality"
}

# Output devices enrolled in quality updates
output "quality_update_devices" {
  value = [
    for device in data.microsoft365_graph_beta_windows_updates_device_enrollment.quality_updates.devices :
    {
      id         = device.id
      enrollment = device.enrollments.quality
    }
  ]
}

# Example: Filter for feature updates
data "microsoft365_graph_beta_windows_updates_device_enrollment" "feature_updates" {
  list_all        = true
  update_category = "feature"
}
```

### Use OData Filter

```terraform
# Example: Use OData filter for advanced queries
data "microsoft365_graph_beta_windows_updates_device_enrollment" "filtered_devices" {
  odata_filter = "id eq 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa' or id eq 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb'"
}

# Output the filtered devices
output "filtered_enrollment_status" {
  value = [
    for device in data.microsoft365_graph_beta_windows_updates_device_enrollment.filtered_devices.devices :
    {
      id          = device.id
      enrollments = device.enrollments
    }
  ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `device_name` (String) The device name to search for. The data source will resolve the name to an Entra device ID and then fetch enrollment status. One of `entra_device_id`, `device_name`, or `list_all` must be specified.
- `entra_device_id` (String) The Entra ID (Azure AD) device object ID to query enrollment status for. One of `entra_device_id`, `device_name`, or `list_all` must be specified.
- `list_all` (Boolean) Set to `true` to list all enrolled devices. Cannot be combined with `entra_device_id` or `device_name`. When using this option, the data source returns all devices in the `devices` attribute. Use `update_category` or `odata_filter` to narrow results.
- `odata_filter` (String) Custom OData filter query for advanced filtering when using `list_all`. Example: `id eq '12345678-1234-1234-1234-123456789012'`. Only applicable when `list_all` is `true`.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `update_category` (String) Optional filter to only return devices enrolled in a specific update category. Valid values: `feature`, `quality`, `driver`. Can be used with any lookup method to filter results.

### Read-Only

- `devices` (Attributes List) List of enrolled devices with their enrollment status. When querying by `entra_device_id` or `device_name`, this will contain a single device. When using `list_all`, this may contain multiple devices. (see [below for nested schema](#nestedatt--devices))

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).


<a id="nestedatt--devices"></a>
### Nested Schema for `devices`

Read-Only:

- `enrollments` (Attributes List) List of update management enrollments for this device. (see [below for nested schema](#nestedatt--devices--enrollments))
- `errors` (Attributes List) List of errors associated with this device's enrollment. (see [below for nested schema](#nestedatt--devices--errors))
- `id` (String) The Entra ID (Azure AD) device object ID.

<a id="nestedatt--devices--enrollments"></a>
### Nested Schema for `devices.enrollments`

Read-Only:

- `update_category` (String) The update category the device is enrolled in (feature, quality, driver).


<a id="nestedatt--devices--errors"></a>
### Nested Schema for `devices.errors`

Read-Only:

- `error_code` (String) The error code.
- `error_message` (String) The error message.
