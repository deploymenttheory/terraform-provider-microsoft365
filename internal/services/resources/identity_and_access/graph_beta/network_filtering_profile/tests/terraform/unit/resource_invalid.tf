resource "microsoft365_graph_beta_identity_and_access_network_filtering_profile" "test" {
  name        = "Test Filtering Profile"
  description = "Test filtering profile with invalid state"
  priority    = 100
  state       = "invalid_state"
}
