resource "microsoft365_graph_beta_identity_and_access_network_forwarding_profile_private_access_policy_link" "test" {
  forwarding_profile_id = "11111111-1111-1111-1111-111111111111"
  profile_type          = "tenant"
  forwarding_policy_id  = "33333333-3333-3333-3333-333333333333"
  config_type           = "appPolicy"
}
