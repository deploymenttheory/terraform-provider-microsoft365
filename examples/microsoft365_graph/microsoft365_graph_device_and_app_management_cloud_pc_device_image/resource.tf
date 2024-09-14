resource "microsoft365_device_and_app_management_cloud_pc_device_image" "example" {
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

  # Optional: Define custom timeouts
  timeouts = {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "30m"
  }
}

# Resource: Create a new Cloud PC Device Image
resource "microsoft365_device_and_app_management_cloud_pc_device_image" "example" {
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

  # Optional: Define custom timeouts
  timeouts = {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "30m"
  }
}

# Data Source: Retrieve an existing Cloud PC Device Image
data "microsoft365_cloud_pc_device_image" "existing" {
  id = microsoft365_cloud_pc_device_image.example.id

  timeouts = {
    read = "30s"
  }
}

