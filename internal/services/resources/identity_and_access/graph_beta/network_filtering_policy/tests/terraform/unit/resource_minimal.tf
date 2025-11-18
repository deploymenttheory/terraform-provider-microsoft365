resource "microsoft365_graph_beta_identity_and_access_network_filtering_policy" "test" {
  name        = "Test Filtering Policy"
  description = "Test filtering policy for unit testing"
  action      = "block"
}

