# Maximal Agent Collection configuration for unit testing
resource "microsoft365_graph_beta_agents_agent_collection" "test_maximal" {
  display_name = "Unit Test Agent Collection Maximal"
  owner_ids = [
    "11111111-1111-1111-1111-111111111111",
    "22222222-2222-2222-2222-222222222222"
  ]
  description       = "A comprehensive test agent collection with all fields configured"
  managed_by        = "33333333-3333-3333-3333-333333333333"
  originating_store = "Terraform"
}
