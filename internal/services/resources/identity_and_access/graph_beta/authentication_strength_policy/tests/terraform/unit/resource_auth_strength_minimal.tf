resource "microsoft365_graph_beta_identity_and_access_authentication_strength_policy" "auth_strength_minimal" {
  display_name = "unit-test-auth-strength-min"
  description  = "Unit test minimal authentication strength policy"

  allowed_combinations = [
    "password,sms"
  ]
}
