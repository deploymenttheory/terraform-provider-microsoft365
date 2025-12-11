# Maximal configuration for update testing
# Note: originating_store is not included as it requires resource replacement
resource "microsoft365_graph_beta_agents_agent_collection" "test_update" {
  display_name = "Unit Test Agent Collection Update Maximal"
  owner_ids = [
    "11111111-1111-1111-1111-111111111111",
    "22222222-2222-2222-2222-222222222222"
  ]
  description = "Updated agent collection with all fields configured"
  managed_by  = "33333333-3333-3333-3333-333333333333"
}
