---
page_title: "microsoft365_graph_beta_windows_365_cloud_pc_source_device_image Data Source - terraform-provider-microsoft365"
subcategory: "Windows 365"
description: |-
  Retrieves Cloud PC source device images available for upload and use on Cloud PCs using the /deviceManagement/virtualEndpoint/deviceImages/getSourceImages endpoint. Supports filtering by all, id, or display_name for image discovery and automation.
---

# microsoft365_graph_beta_windows_365_cloud_pc_source_device_image (Data Source)

Retrieves Cloud PC source device images available for upload and use on Cloud PCs using the `/deviceManagement/virtualEndpoint/deviceImages/getSourceImages` endpoint. Supports filtering by all, id, or display_name for image discovery and automation.

## Microsoft Documentation

- [cloudPcDeviceImage: getSourceImages](https://learn.microsoft.com/en-us/graph/api/cloudpcdeviceimage-getsourceimages?view=graph-rest-beta)

## API Permissions

The following API permissions are required in order to use this data source.

### Microsoft Graph

- **Application**: `CloudPC.Read.All`, `CloudPC.ReadWrite.All`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.18.0-alpha | Experimental | Initial release |

## Example Usage

```terraform
# Example: Retrieve all Cloud PC source device images
# filter_type = "all" returns all images. filter_value is not required.
data "microsoft365_graph_beta_windows_365_cloud_pc_source_device_image" "all" {
  filter_type = "all"
}

output "all_cloud_pc_source_device_images" {
  value = data.microsoft365_graph_beta_windows_365_cloud_pc_source_device_image.all.items
}

# Example: Retrieve a Cloud PC source device image by ID
# filter_type = "id" requires filter_value to be set to the exact image id (see output from the 'all' query above)
data "microsoft365_graph_beta_windows_365_cloud_pc_source_device_image" "by_id" {
  filter_type  = "id"
  filter_value = "<image_id>" # Replace with a real image id
}

output "cloud_pc_source_device_image_by_id" {
  value = data.microsoft365_graph_beta_windows_365_cloud_pc_source_device_image.by_id.items
}

# Example: Retrieve Cloud PC source device images by display name substring
# filter_type = "display_name" requires filter_value to be a substring to match against image display names
data "microsoft365_graph_beta_windows_365_cloud_pc_source_device_image" "by_display_name" {
  filter_type  = "display_name"
  filter_value = "<substring>" # Replace with part of the display name
}

output "cloud_pc_source_device_images_by_display_name" {
  value = data.microsoft365_graph_beta_windows_365_cloud_pc_source_device_image.by_display_name.items
}

# Valid values for filter_type:
#   - "all": Returns all images. filter_value is ignored.
#   - "id": Returns the image with the exact id specified in filter_value.
#   - "display_name": Returns images whose display_name contains the filter_value substring (case-insensitive).
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `filter_type` (String) Type of filter to apply. Valid values are: `all`, `id`, `display_name`. Use 'all' to retrieve all images, 'id' to retrieve a specific image by its unique identifier, or 'display_name' to filter by the image's display name.

### Optional

- `filter_value` (String) Value to filter by. Not required when filter_type is 'all'. For 'id', provide the image ID. For 'display_name', provide a substring to match against image display names.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `items` (Attributes List) The list of Cloud PC source device images that match the filter criteria. (see [below for nested schema](#nestedatt--items))

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

- `display_name` (String) The display name of the source device image.
- `id` (String) The unique identifier for the source device image.
- `resource_id` (String) The resource ID for the source device image.
- `subscription_display_name` (String) The display name of the subscription.
- `subscription_id` (String) The subscription ID associated with the image. 