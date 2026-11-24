resource "microsoft365_graph_beta_identity_and_access_network_filtering_profile_policy_link" "test" {
  filtering_profile_id = "11111111-2222-3333-4444-555555555555"
  policy_id            = "66666666-7777-8888-9999-000000000000"
  policy_type          = "filtering"
  state                = "disabled"
  priority             = 200
  logging_state        = "disabled"
}
