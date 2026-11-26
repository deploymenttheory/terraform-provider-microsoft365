resource "random_uuid" "maximal_suffix" {}

resource "microsoft365_graph_beta_applications_token_lifetime_policy" "maximal" {
  display_name            = "acc-test-max-${random_uuid.maximal_suffix.result}"
  is_organization_default = false
  definition              = ["{\"TokenLifetimePolicy\":{\"Version\":1,\"AccessTokenLifetime\":\"02:00:00\"}}"]
}
