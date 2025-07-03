resource "microsoft365_graph_beta_windows_365_cloud_pc_device_image" "example" {
  # This resource is read-only. Device images are created by uploading images via the Microsoft Endpoint Manager or other supported methods.
  # You can import an existing device image using its ID.
  display_name             = "My Custom Cloud PC Image"
  version                  = "1.0.0"
  # Must match the Azure image resource ID format:
  # /subscriptions/{subscription-id}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/images/{imageName}
  source_image_resource_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.Compute/images/myimage"
} 