# Minimal Agent Collection configuration for unit testing
resource "microsoft365_graph_beta_agents_agent_collection" "test_minimal" {
  display_name = "Unit Test Agent Collection Minimal"
  owner_ids    = ["11111111-1111-1111-1111-111111111111"]
}
