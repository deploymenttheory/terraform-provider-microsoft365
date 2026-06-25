resource "microsoft365_graph_beta_applications_token_lifetime_policy" "minimal" {
  display_name = "unit-test-token-lifetime-policy-min"
  description  = "Unit test minimal token lifetime policy"
  definition   = ["{\"TokenLifetimePolicy\":{\"Version\":1,\"AccessTokenLifetime\":\"01:00:00\"}}"]
}
