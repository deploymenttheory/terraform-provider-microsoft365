# Link a Private Access app policy to a custom user forwarding profile.
resource "microsoft365_graph_beta_identity_and_access_network_forwarding_profile_private_access_policy_link" "app" {
  forwarding_profile_id = "00000000-0000-0000-0000-000000000001"
  forwarding_policy_id  = "00000000-0000-0000-0000-000000000002"
  config_type           = "appPolicy"
  priority              = 120
}

# Link a Private Access app policy to a pre-logon (device) forwarding profile.
resource "microsoft365_graph_beta_identity_and_access_network_forwarding_profile_private_access_policy_link" "device" {
  forwarding_profile_id = "00000000-0000-0000-0000-000000000003"
  profile_type          = "device"
  forwarding_policy_id  = "00000000-0000-0000-0000-000000000002"
  config_type           = "appPolicy"
}
