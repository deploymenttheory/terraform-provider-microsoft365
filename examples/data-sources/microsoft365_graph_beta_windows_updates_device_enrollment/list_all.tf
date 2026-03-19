# Example: List all enrolled devices
data "microsoft365_graph_beta_windows_updates_device_enrollment" "all_devices" {
  list_all = true
}

# Output count of enrolled devices
output "total_enrolled_devices" {
  value = length(data.microsoft365_graph_beta_windows_updates_device_enrollment.all_devices.devices)
}

# Output devices enrolled in quality updates
output "quality_enrolled_devices" {
  value = [
    for device in data.microsoft365_graph_beta_windows_updates_device_enrollment.all_devices.devices :
    device.id if device.enrollments.quality != null
  ]
}
