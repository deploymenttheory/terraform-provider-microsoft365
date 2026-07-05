resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = true
  lower   = true
}

resource "microsoft365_graph_beta_identity_and_access_network_filtering_profile" "enabled" {
  name        = "acc-test-filtering-profile-enabled-${random_string.suffix.result}"
  description = "Acceptance test enabled filtering profile"
  priority    = 100
  state       = "enabled"
}
