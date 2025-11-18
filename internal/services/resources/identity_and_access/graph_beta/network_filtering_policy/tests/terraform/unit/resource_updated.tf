resource "microsoft365_graph_beta_identity_and_access_network_filtering_policy" "test" {
  name        = "Updated Filtering Policy"
  description = "Updated description"
  action      = "allow"
}

