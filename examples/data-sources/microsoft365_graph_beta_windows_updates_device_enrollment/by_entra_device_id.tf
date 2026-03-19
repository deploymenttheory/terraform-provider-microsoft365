# Example: Get enrollment status for a specific device by Entra ID
data "microsoft365_graph_beta_windows_updates_device_enrollment" "by_id" {
  entra_device_id = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
}

# Output the enrollment details
output "device_enrollment_status" {
  value = {
    device_id   = data.microsoft365_graph_beta_windows_updates_device_enrollment.by_id.devices[0].id
    enrollments = data.microsoft365_graph_beta_windows_updates_device_enrollment.by_id.devices[0].enrollments
    errors      = data.microsoft365_graph_beta_windows_updates_device_enrollment.by_id.devices[0].errors
  }
}
