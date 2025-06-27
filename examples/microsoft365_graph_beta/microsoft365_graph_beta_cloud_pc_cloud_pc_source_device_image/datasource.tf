# Example: Retrieve all Cloud PC source device images
# filter_type = "all" returns all images. filter_value is not required.
data "microsoft365_graph_beta_cloud_pc_cloud_pc_source_device_image" "all" {
  filter_type = "all"
}

output "all_cloud_pc_source_device_images" {
  value = data.microsoft365_graph_beta_cloud_pc_cloud_pc_source_device_image.all.items
}

# Example: Retrieve a Cloud PC source device image by ID
# filter_type = "id" requires filter_value to be set to the exact image id (see output from the 'all' query above)
data "microsoft365_graph_beta_cloud_pc_cloud_pc_source_device_image" "by_id" {
  filter_type  = "id"
  filter_value = "<image_id>" # Replace with a real image id
}

output "cloud_pc_source_device_image_by_id" {
  value = data.microsoft365_graph_beta_cloud_pc_cloud_pc_source_device_image.by_id.items
}

# Example: Retrieve Cloud PC source device images by display name substring
# filter_type = "display_name" requires filter_value to be a substring to match against image display names
data "microsoft365_graph_beta_cloud_pc_cloud_pc_source_device_image" "by_display_name" {
  filter_type  = "display_name"
  filter_value = "<substring>" # Replace with part of the display name
}

output "cloud_pc_source_device_images_by_display_name" {
  value = data.microsoft365_graph_beta_cloud_pc_cloud_pc_source_device_image.by_display_name.items
}

# Valid values for filter_type:
#   - "all": Returns all images. filter_value is ignored.
#   - "id": Returns the image with the exact id specified in filter_value.
#   - "display_name": Returns images whose display_name contains the filter_value substring (case-insensitive). 