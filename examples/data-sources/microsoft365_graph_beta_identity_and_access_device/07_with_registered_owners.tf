# Example: Get a device and its registered owners

data "microsoft365_graph_beta_identity_and_access_device" "with_owners" {
  object_id              = "00000000-0000-0000-0000-000000000000"
  list_registered_owners = true
}

output "device_name" {
  value = length(data.microsoft365_graph_beta_identity_and_access_device.with_owners.items) > 0 ? data.microsoft365_graph_beta_identity_and_access_device.with_owners.items[0].display_name : null
}

output "registered_owners" {
  value = [for owner in data.microsoft365_graph_beta_identity_and_access_device.with_owners.registered_owners : {
    id           = owner.id
    display_name = owner.display_name
    type         = owner.odata_type
  }]
}

output "owner_count" {
  value = length(data.microsoft365_graph_beta_identity_and_access_device.with_owners.registered_owners)
}
