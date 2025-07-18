---
page_title: "microsoft365_graph_beta_windows_365_cloud_pc_gallery_image Data Source - terraform-provider-microsoft365"
subcategory: "Windows 365"
description: |-
  Retrieves Cloud PC Gallery Images from Microsoft Intune using the /deviceManagement/virtualEndpoint/galleryImages endpoint. Supports filtering by all, id, or display_name for image discovery and selection.
---

# microsoft365_graph_beta_windows_365_cloud_pc_gallery_image (Data Source)

Retrieves Cloud PC Gallery Images from Microsoft Intune using the `/deviceManagement/virtualEndpoint/galleryImages` endpoint. Supports filtering by all, id, or display_name for image discovery and selection.

Gallery images represent the available OS images that can be used to provision Cloud PCs. This data source allows you to discover available images, their support status, and key details like OS version and size.

## Microsoft Documentation

- [cloudPcGalleryImage resource type](https://learn.microsoft.com/en-us/graph/api/resources/cloudpcgalleryimage?view=graph-rest-beta)

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
# Example: Retrieve all Cloud PC Gallery Images

data "microsoft365_graph_beta_windows_365_cloud_pc_gallery_image" "all" {
  filter_type = "all"
}

# Output: List all gallery image IDs
output "all_gallery_image_ids" {
  value = [for image in data.microsoft365_graph_beta_windows_365_cloud_pc_gallery_image.all.items : image.id]
}

# Output: Show all details for the first gallery image (if present)
output "first_gallery_image_details" {
  value = data.microsoft365_graph_beta_windows_365_cloud_pc_gallery_image.all.items[0]
}

# Example: Retrieve a specific gallery image by ID
data "microsoft365_graph_beta_windows_365_cloud_pc_gallery_image" "by_id" {
  filter_type  = "id"
  filter_value = "MicrosoftWindowsDesktop_windows-ent-cpc_win11-22h2-ent-cpc-m365" # Example ID format
}

output "gallery_image_by_id" {
  value = data.microsoft365_graph_beta_windows_365_cloud_pc_gallery_image.by_id.items[0]
}

# Example: Retrieve gallery images by display name substring
data "microsoft365_graph_beta_windows_365_cloud_pc_gallery_image" "by_display_name" {
  filter_type  = "display_name"
  filter_value = "Windows 11" # This will match images containing "Windows 11" in their name
}

output "gallery_images_by_display_name" {
  value = data.microsoft365_graph_beta_windows_365_cloud_pc_gallery_image.by_display_name.items
}

# Example: Show available Windows 11 images with their status and dates
output "windows_11_images_status" {
  value = [for image in data.microsoft365_graph_beta_windows_365_cloud_pc_gallery_image.by_display_name.items : {
    display_name    = image.display_name
    status          = image.status
    start_date      = image.start_date
    end_date        = image.end_date
    expiration_date = image.expiration_date
    size_in_gb      = image.size_in_gb
  }]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `filter_type` (String) Type of filter to apply. Valid values are: `all`, `id`, `display_name`. Use 'all' to retrieve all gallery images, 'id' to retrieve a specific image by its unique identifier, or 'display_name' to filter by the image's display name.

### Optional

- `filter_value` (String) Value to filter by. Not required when filter_type is 'all'. For 'id', provide the gallery image ID. For 'display_name', provide a substring to match against image display names.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `items` (Attributes List) The list of Cloud PC gallery images that match the filter criteria. (see [below for nested schema](#nestedatt--items))

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

- `display_name` (String) The display name of this gallery image. For example, Windows 11 Enterprise + Microsoft 365 Apps 22H2.
- `end_date` (String) The date when the status of image becomes supportedWithWarning.
- `expiration_date` (String) The date when the image is no longer available. Users are unable to provision new Cloud PCs if the current time is later than this date.
- `id` (String) The unique identifier (ID) of the gallery image resource on Cloud PC. The ID format is {publisherName_offerName_skuName}.
- `offer_name` (String) The offer name of this gallery image that is passed to ARM to retrieve the image resource.
- `os_version_number` (String) The operating system version of this gallery image. For example, 10.0.22000.296.
- `publisher_name` (String) The publisher name of this gallery image that is passed to ARM to retrieve the image resource.
- `size_in_gb` (Number) Indicates the size of this image in gigabytes. For example, 64.
- `sku_name` (String) The SKU name of this image that is passed to ARM to retrieve the image resource.
- `start_date` (String) The date when the Cloud PC image is available for provisioning new Cloud PCs.
- `status` (String) The status of the gallery image on the Cloud PC. Possible values are: supported, supportedWithWarning, notSupported, unknownFutureValue. 