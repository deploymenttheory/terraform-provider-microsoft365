# Example: Get a device and its registered users

data "microsoft365_graph_beta_identity_and_access_device" "with_users" {
  object_id             = "00000000-0000-0000-0000-000000000000"
  list_registered_users = true
}

output "device_name" {
  value = length(data.microsoft365_graph_beta_identity_and_access_device.with_users.items) > 0 ? data.microsoft365_graph_beta_identity_and_access_device.with_users.items[0].display_name : null
}

output "registered_users" {
  value = [for user in data.microsoft365_graph_beta_identity_and_access_device.with_users.registered_users : {
    id           = user.id
    display_name = user.display_name
    type         = user.odata_type
  }]
}

output "user_count" {
  value = length(data.microsoft365_graph_beta_identity_and_access_device.with_users.registered_users)
}
