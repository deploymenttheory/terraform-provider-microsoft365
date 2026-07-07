resource "microsoft365_graph_beta_applications_service_principal_token_lifetime_policy_assignment" "basic" {
  service_principal_id     = var.service_principal_id
  token_lifetime_policy_id = var.token_lifetime_policy_id
}
