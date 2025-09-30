# Example: Retrieve all Cloud PC Device Images

data "microsoft365_graph_beta_windows_365_cloud_pc_device_image" "all" {
  filter_type = "all"
}

# Output: List all device image IDs
output "all_device_image_ids" {
  value = [for image in data.microsoft365_graph_beta_windows_365_cloud_pc_device_image.all.items : image.id]
}

# Output: Show all details for the first device image (if present)
output "first_device_image_details" {
  value = data.microsoft365_graph_beta_windows_365_cloud_pc_device_image.all.items[0]
}

# Example: Retrieve a specific device image by ID
# Replace with a valid Device Image ID from your environment

data "microsoft365_graph_beta_windows_365_cloud_pc_device_image" "by_id" {
  filter_type  = "id"
  filter_value = "<device_image_id>" # Example: "MicrosoftWindowsDesktop_windows-ent-cpc_win11-22h2-ent-cpc-m365"
}

output "device_image_by_id" {
  value = data.microsoft365_graph_beta_windows_365_cloud_pc_device_image.by_id.items[0]
}

# Example: Retrieve device images by display name substring
# This will match images containing "Windows 11" in their display name

data "microsoft365_graph_beta_windows_365_cloud_pc_device_image" "by_display_name" {
  filter_type  = "display_name"
  filter_value = "Windows 11"
}

output "device_image_by_display_name" {
  value = data.microsoft365_graph_beta_windows_365_cloud_pc_device_image.by_display_name.items
}

# Output: Show available Windows 11 images with their status, expiration date, and OS version number
output "windows_11_images_summary" {
  value = [for image in data.microsoft365_graph_beta_windows_365_cloud_pc_device_image.by_display_name.items : {
    display_name      = image.display_name
    status            = image.status
    expiration_date   = image.expiration_date
    os_version_number = image.os_version_number
  }]
} 