resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = true
  lower   = true
}

resource "microsoft365_graph_beta_identity_and_access_network_filtering_policy" "minimal" {
  name   = "acc-test-filtering-policy-minimal-nodesc-${random_string.suffix.result}"
  action = "block"
}

