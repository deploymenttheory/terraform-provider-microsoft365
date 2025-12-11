# Minimal Agent Collection Assignment configuration
# Assigns an agent instance to an agent collection
resource "microsoft365_graph_beta_agents_agent_collection_assignment" "example" {
  agent_instance_id   = "00000000-0000-0000-0000-000000000001"
  agent_collection_id = "00000000-0000-0000-0000-000000000002"
}
