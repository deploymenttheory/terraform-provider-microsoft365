# Example: Get a device and its group memberships

data "microsoft365_graph_beta_identity_and_access_device" "with_groups" {
  object_id      = "00000000-0000-0000-0000-000000000000"
  list_member_of = true
}

output "device_name" {
  value = length(data.microsoft365_graph_beta_identity_and_access_device.with_groups.items) > 0 ? data.microsoft365_graph_beta_identity_and_access_device.with_groups.items[0].display_name : null
}

output "group_memberships" {
  value = [for group in data.microsoft365_graph_beta_identity_and_access_device.with_groups.member_of : {
    id           = group.id
    display_name = group.display_name
    type         = group.odata_type
  }]
}

output "group_count" {
  value = length(data.microsoft365_graph_beta_identity_and_access_device.with_groups.member_of)
}
