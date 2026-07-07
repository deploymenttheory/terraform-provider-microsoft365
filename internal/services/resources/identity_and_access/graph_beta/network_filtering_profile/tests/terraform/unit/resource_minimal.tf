resource "microsoft365_graph_beta_identity_and_access_network_filtering_profile" "test" {
  name        = "unit-test-filtering-profile-minimal"
  description = "Test filtering profile for unit testing"
  priority    = 100
  state       = "enabled"
}
