resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = true
  lower   = true
}

resource "microsoft365_graph_beta_identity_and_access_network_filtering_profile" "test" {
  name        = "acc-test-filtering-profile-updated-${random_string.suffix.result}"
  description = "Updated acceptance test filtering profile"
  priority    = 200
  state       = "disabled"
}
