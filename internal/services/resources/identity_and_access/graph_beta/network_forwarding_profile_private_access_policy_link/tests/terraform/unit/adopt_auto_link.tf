resource "microsoft365_graph_beta_identity_and_access_network_forwarding_profile_private_access_policy_link" "test" {
  forwarding_profile_id = "11111111-1111-1111-1111-111111111111"
  forwarding_policy_id  = "44444444-4444-4444-4444-444444444444"
  config_type           = "appPolicy"
}
