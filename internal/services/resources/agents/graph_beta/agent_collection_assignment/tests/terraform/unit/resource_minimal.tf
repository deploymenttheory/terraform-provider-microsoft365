# Agent Collection Assignment Unit Test - Minimal Configuration

resource "microsoft365_graph_beta_agents_agent_collection_assignment" "test_minimal" {
  agent_instance_id   = "11111111-1111-1111-1111-111111111111"
  agent_collection_id = "22222222-2222-2222-2222-222222222222"
}
