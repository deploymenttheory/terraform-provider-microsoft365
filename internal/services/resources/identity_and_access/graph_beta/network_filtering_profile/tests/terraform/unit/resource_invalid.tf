resource "microsoft365_graph_beta_identity_and_access_network_filtering_profile" "test" {
  name        = "Test Filtering Profile"
  description = "Test filtering profile with invalid state"
  state       = "invalid_state"
}
