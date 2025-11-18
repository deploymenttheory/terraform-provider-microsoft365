resource "microsoft365_graph_beta_identity_and_access_authentication_strength" "auth_strength_minimal" {
  # Display name must be 30 characters or less
  display_name = "unit-test-auth-strength-min"
  description  = "Unit test minimal authentication strength policy"

  allowed_combinations = [
    "password,sms"
  ]
}