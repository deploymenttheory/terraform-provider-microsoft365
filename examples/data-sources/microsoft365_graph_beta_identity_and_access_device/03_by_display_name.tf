# Example: Get devices by display name

data "microsoft365_graph_beta_identity_and_access_device" "by_name" {
  display_name = "DESKTOP-ABC123"
}

output "device_id" {
  value = length(data.microsoft365_graph_beta_identity_and_access_device.by_name.items) > 0 ? data.microsoft365_graph_beta_identity_and_access_device.by_name.items[0].device_id : null
}

output "device_trust_type" {
  value = length(data.microsoft365_graph_beta_identity_and_access_device.by_name.items) > 0 ? data.microsoft365_graph_beta_identity_and_access_device.by_name.items[0].trust_type : null
}
