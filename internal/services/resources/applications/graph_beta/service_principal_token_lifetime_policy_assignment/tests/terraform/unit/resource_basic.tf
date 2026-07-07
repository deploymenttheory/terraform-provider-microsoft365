resource "microsoft365_graph_beta_applications_service_principal_token_lifetime_policy_assignment" "basic" {
  service_principal_id     = "00000000-0000-0000-0000-000000000020"
  token_lifetime_policy_id = "00000000-0000-0000-0000-000000000010"
}
