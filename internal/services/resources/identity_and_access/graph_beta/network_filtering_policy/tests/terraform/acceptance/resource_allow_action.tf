resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = true
  lower   = true
}

resource "microsoft365_graph_beta_identity_and_access_network_filtering_policy" "allow" {
  name        = "acc-test-filtering-policy-allow-${random_string.suffix.result}"
  description = "Acceptance test filtering policy with allow action"
  action      = "allow"
}

