resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = true
  lower   = true
}

resource "microsoft365_graph_beta_identity_and_access_network_filtering_profile" "minimal" {
  name     = "acc-test-filtering-profile-minimal-nodesc-${random_string.suffix.result}"
  priority = 300
  state    = "enabled"
}
