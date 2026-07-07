resource "microsoft365_graph_beta_identity_and_access_network_filtering_profile" "test" {
  name        = "unit-test-filtering-profile-updated"
  description = "Updated description"
  priority    = 200
  state       = "disabled"
}
