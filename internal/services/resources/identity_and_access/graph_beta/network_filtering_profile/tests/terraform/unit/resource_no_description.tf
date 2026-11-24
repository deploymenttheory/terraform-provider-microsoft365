resource "microsoft365_graph_beta_identity_and_access_network_filtering_profile" "test" {
  name     = "unit-test-filtering-profile-no-description"
  priority = 300
  state    = "enabled"
}
