# Example: Get enrollment status for a device by name
data "microsoft365_graph_beta_windows_updates_device_enrollment" "by_name" {
  device_name = "DESKTOP-ABC123"
}

# Output the enrollment details
output "device_enrollment_by_name" {
  value = {
    device_id   = data.microsoft365_graph_beta_windows_updates_device_enrollment.by_name.devices[0].id
    enrollments = data.microsoft365_graph_beta_windows_updates_device_enrollment.by_name.devices[0].enrollments
  }
}
