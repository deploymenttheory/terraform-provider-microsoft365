resource "microsoft365_graph_beta_identity_and_access_network_filtering_policy" "test" {
  name        = "unit-test-filtering-policy-updated"
  description = "Updated description"
  action      = "allow"
}

