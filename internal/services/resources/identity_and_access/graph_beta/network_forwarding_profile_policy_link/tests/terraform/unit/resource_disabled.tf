resource "microsoft365_graph_beta_identity_and_access_network_forwarding_profile_policy_link" "test" {
  forwarding_profile_id = "11111111-2222-3333-4444-555555555555"
  policy_link_id        = "66666666-7777-8888-9999-000000000000"
  state                 = "disabled"
}
