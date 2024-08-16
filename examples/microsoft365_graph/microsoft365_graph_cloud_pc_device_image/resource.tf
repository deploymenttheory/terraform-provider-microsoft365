resource "microsoft365_cloud_pc_device_image" "example" {
  display_name             = "My Custom Device Image"
  source_image_resource_id = "/subscriptions/your-subscription-id/resourceGroups/your-resource-group/providers/Microsoft.Compute/images/your-image-name"
  version                  = "1.0.0"

  lifecycle {
    ignore_changes = [
      error_code,
      expiration_date,
      last_modified_date_time,
      operating_system,
      os_build_number,
      os_status,
      status
    ]
  }
}

output "cloud_pc_device_image_id" {
  value = microsoft365_cloud_pc_device_image.example.id
}