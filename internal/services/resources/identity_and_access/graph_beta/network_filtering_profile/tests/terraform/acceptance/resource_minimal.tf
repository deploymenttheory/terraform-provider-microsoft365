resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = true
  lower   = true
}

resource "microsoft365_graph_beta_identity_and_access_network_filtering_profile" "test" {
  name        = "acc-test-filtering-profile-minimal-${random_string.suffix.result}"
  description = "Acceptance test minimal filtering profile configuration"
  priority    = 100
  state       = "enabled"
}
