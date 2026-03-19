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
