# Example: Get a device by its object ID

data "microsoft365_graph_beta_identity_and_access_device" "by_id" {
  object_id = "00000000-0000-0000-0000-000000000000"
}

output "device_display_name" {
  value = length(data.microsoft365_graph_beta_identity_and_access_device.by_id.items) > 0 ? data.microsoft365_graph_beta_identity_and_access_device.by_id.items[0].display_name : null
}

output "device_operating_system" {
  value = length(data.microsoft365_graph_beta_identity_and_access_device.by_id.items) > 0 ? data.microsoft365_graph_beta_identity_and_access_device.by_id.items[0].operating_system : null
}
