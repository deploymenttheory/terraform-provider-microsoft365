# Example: List all devices in the tenant

data "microsoft365_graph_beta_identity_and_access_device" "all" {
  list_all = true
}

# Output the first device's display name
output "first_device_name" {
  value = length(data.microsoft365_graph_beta_identity_and_access_device.all.items) > 0 ? data.microsoft365_graph_beta_identity_and_access_device.all.items[0].display_name : "No devices found"
}

# Output the total count of devices
output "device_count" {
  value = length(data.microsoft365_graph_beta_identity_and_access_device.all.items)
}
