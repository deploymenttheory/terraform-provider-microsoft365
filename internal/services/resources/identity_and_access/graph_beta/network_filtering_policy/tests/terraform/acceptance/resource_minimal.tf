resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = true
  lower   = true
}

resource "microsoft365_graph_beta_identity_and_access_network_filtering_policy" "test" {
  name        = "acc-test-filtering-policy-minimal-${random_string.suffix.result}"
  description = "Acceptance test minimal filtering policy configuration"
  action      = "block"
}

