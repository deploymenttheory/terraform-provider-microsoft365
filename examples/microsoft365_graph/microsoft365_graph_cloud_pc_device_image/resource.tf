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

# Resource: Create a new Cloud PC Device Image
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

# Data Source: Retrieve an existing Cloud PC Device Image
data "microsoft365_cloud_pc_device_image" "existing" {
  id = microsoft365_cloud_pc_device_image.example.id

  timeouts = {
    read = "30s"
  }
}

# Outputs for the Resource
output "cloud_pc_device_image_id" {
  value       = microsoft365_cloud_pc_device_image.example.id
  description = "The ID of the created Cloud PC Device Image"
}

output "cloud_pc_device_image_display_name" {
  value       = microsoft365_cloud_pc_device_image.example.display_name
  description = "The display name of the created Cloud PC Device Image"
}

output "cloud_pc_device_image_version" {
  value       = microsoft365_cloud_pc_device_image.example.version
  description = "The version of the created Cloud PC Device Image"
}

output "cloud_pc_device_image_source_id" {
  value       = microsoft365_cloud_pc_device_image.example.source_image_resource_id
  description = "The source image resource ID of the created Cloud PC Device Image"
}

output "existing_cloud_pc_device_image_display_name" {
  value       = data.microsoft365_cloud_pc_device_image.existing.display_name
  description = "The display name of the existing Cloud PC Device Image"
}

output "existing_cloud_pc_device_image_error_code" {
  value       = data.microsoft365_cloud_pc_device_image.existing.error_code
  description = "The error code of the existing Cloud PC Device Image, if any"
}

output "existing_cloud_pc_device_image_expiration_date" {
  value       = data.microsoft365_cloud_pc_device_image.existing.expiration_date
  description = "The expiration date of the existing Cloud PC Device Image"
}

output "existing_cloud_pc_device_image_last_modified" {
  value       = data.microsoft365_cloud_pc_device_image.existing.last_modified_date_time
  description = "The last modified date and time of the existing Cloud PC Device Image"
}

output "existing_cloud_pc_device_image_os" {
  value       = data.microsoft365_cloud_pc_device_image.existing.operating_system
  description = "The operating system of the existing Cloud PC Device Image"
}

output "existing_cloud_pc_device_image_os_build" {
  value       = data.microsoft365_cloud_pc_device_image.existing.os_build_number
  description = "The OS build number of the existing Cloud PC Device Image"
}

output "existing_cloud_pc_device_image_os_status" {
  value       = data.microsoft365_cloud_pc_device_image.existing.os_status
  description = "The OS status of the existing Cloud PC Device Image"
}

output "existing_cloud_pc_device_image_source_id" {
  value       = data.microsoft365_cloud_pc_device_image.existing.source_image_resource_id
  description = "The source image resource ID of the existing Cloud PC Device Image"
}

output "existing_cloud_pc_device_image_status" {
  value       = data.microsoft365_cloud_pc_device_image.existing.status
  description = "The status of the existing Cloud PC Device Image"
}

output "existing_cloud_pc_device_image_version" {
  value       = data.microsoft365_cloud_pc_device_image.existing.version
  description = "The version of the existing Cloud PC Device Image"
}