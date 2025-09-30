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