resource "microsoft365_graph_beta_identity_and_access_authentication_strength" "auth_strength_minimal" {
  display_name = "acc-test-authentication-strength-minimal"
  description  = "Acceptance test minimal authentication strength policy"

  allowed_combinations = [
    "password,sms"
  ]
}