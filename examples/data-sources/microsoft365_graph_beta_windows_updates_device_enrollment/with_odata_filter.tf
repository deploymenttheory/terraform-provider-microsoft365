# Example: Use OData filter for advanced queries
data "microsoft365_graph_beta_windows_updates_device_enrollment" "filtered_devices" {
  odata_filter = "id eq 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa' or id eq 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb'"
}

# Output the filtered devices
output "filtered_enrollment_status" {
  value = [
    for device in data.microsoft365_graph_beta_windows_updates_device_enrollment.filtered_devices.devices :
    {
      id          = device.id
      enrollments = device.enrollments
    }
  ]
}
