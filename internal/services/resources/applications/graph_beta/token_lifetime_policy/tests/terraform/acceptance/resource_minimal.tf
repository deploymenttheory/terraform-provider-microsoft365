resource "random_uuid" "minimal_suffix" {}

resource "microsoft365_graph_beta_applications_token_lifetime_policy" "minimal" {
  display_name = "acc-test-min-${random_uuid.minimal_suffix.result}"
  description  = "Acceptance test minimal token lifetime policy"
  definition   = ["{\"TokenLifetimePolicy\":{\"Version\":1,\"AccessTokenLifetime\":\"01:00:00\"}}"]
}
