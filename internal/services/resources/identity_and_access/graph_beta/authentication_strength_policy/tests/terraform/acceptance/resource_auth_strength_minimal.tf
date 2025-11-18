resource "microsoft365_graph_beta_identity_and_access_authentication_strength" "auth_strength_minimal" {
  # Display name must be 30 characters or less
  display_name = "acc-test-auth-strength-min"
  description  = "Acceptance test minimal authentication strength policy"

  allowed_combinations = [
    "password,sms"
  ]
}