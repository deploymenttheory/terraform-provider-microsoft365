# Example: Get devices by device ID (Azure Device Registration Service ID)

data "microsoft365_graph_beta_identity_and_access_device" "by_device_id" {
  device_id = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
}

output "device_display_name" {
  value = length(data.microsoft365_graph_beta_identity_and_access_device.by_device_id.items) > 0 ? data.microsoft365_graph_beta_identity_and_access_device.by_device_id.items[0].display_name : null
}

output "device_is_compliant" {
  value = length(data.microsoft365_graph_beta_identity_and_access_device.by_device_id.items) > 0 ? data.microsoft365_graph_beta_identity_and_access_device.by_device_id.items[0].is_compliant : null
}
