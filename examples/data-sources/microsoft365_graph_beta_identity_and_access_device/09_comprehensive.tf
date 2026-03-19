# Example: Get a device with all related information

data "microsoft365_graph_beta_identity_and_access_device" "comprehensive" {
  object_id              = "00000000-0000-0000-0000-000000000000"
  list_member_of         = true
  list_registered_owners = true
  list_registered_users  = true
}

# Device information
output "device_info" {
  value = length(data.microsoft365_graph_beta_identity_and_access_device.comprehensive.items) > 0 ? {
    id                       = data.microsoft365_graph_beta_identity_and_access_device.comprehensive.items[0].id
    display_name             = data.microsoft365_graph_beta_identity_and_access_device.comprehensive.items[0].display_name
    device_id                = data.microsoft365_graph_beta_identity_and_access_device.comprehensive.items[0].device_id
    operating_system         = data.microsoft365_graph_beta_identity_and_access_device.comprehensive.items[0].operating_system
    operating_system_version = data.microsoft365_graph_beta_identity_and_access_device.comprehensive.items[0].operating_system_version
    is_compliant             = data.microsoft365_graph_beta_identity_and_access_device.comprehensive.items[0].is_compliant
    is_managed               = data.microsoft365_graph_beta_identity_and_access_device.comprehensive.items[0].is_managed
    trust_type               = data.microsoft365_graph_beta_identity_and_access_device.comprehensive.items[0].trust_type
    account_enabled          = data.microsoft365_graph_beta_identity_and_access_device.comprehensive.items[0].account_enabled
  } : null
}

# Group memberships
output "group_memberships" {
  value = length(data.microsoft365_graph_beta_identity_and_access_device.comprehensive.member_of)
}

# Registered owners
output "registered_owners" {
  value = length(data.microsoft365_graph_beta_identity_and_access_device.comprehensive.registered_owners)
}

# Registered users
output "registered_users" {
  value = length(data.microsoft365_graph_beta_identity_and_access_device.comprehensive.registered_users)
}
