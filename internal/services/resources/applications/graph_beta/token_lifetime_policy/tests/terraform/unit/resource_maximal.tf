resource "microsoft365_graph_beta_applications_token_lifetime_policy" "maximal" {
  display_name            = "unit-test-token-lifetime-policy-max"
  is_organization_default = true
  definition              = ["{\"TokenLifetimePolicy\":{\"Version\":1,\"AccessTokenLifetime\":\"02:00:00\"}}"]
}
