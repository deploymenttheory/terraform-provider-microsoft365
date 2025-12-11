# Minimal configuration for update testing
resource "microsoft365_graph_beta_agents_agent_collection" "test_update" {
  display_name = "Unit Test Agent Collection Update Minimal"
  owner_ids    = ["11111111-1111-1111-1111-111111111111"]
}
