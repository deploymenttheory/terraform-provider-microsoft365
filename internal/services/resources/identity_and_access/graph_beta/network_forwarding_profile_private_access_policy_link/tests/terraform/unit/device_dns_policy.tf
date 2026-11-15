resource "microsoft365_graph_beta_identity_and_access_network_forwarding_profile_private_access_policy_link" "test" {
  forwarding_profile_id = "22222222-2222-2222-2222-222222222222"
  profile_type          = "device"
  forwarding_policy_id  = "55555555-5555-5555-5555-555555555555"
  config_type           = "dnsPolicy"
  state                 = "disabled"
}
