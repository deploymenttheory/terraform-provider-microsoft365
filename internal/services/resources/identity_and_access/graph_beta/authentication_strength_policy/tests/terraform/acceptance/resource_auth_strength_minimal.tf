resource "random_uuid" "auth_strength_minimal" {}

resource "microsoft365_graph_beta_identity_and_access_authentication_strength_policy" "auth_strength_minimal" {
  display_name = substr("acc-test-min-${random_uuid.auth_strength_minimal.result}", 0, 30)
  description  = "Acceptance test minimal authentication strength policy"

  allowed_combinations = [
    "password,sms"
  ]
}
