# Minimal Agent Collection configuration
# Creates an agent collection with required fields only
resource "microsoft365_graph_beta_agents_agent_collection" "minimal" {
  display_name = "My Agent Collection"
  owner_ids    = ["00000000-0000-0000-0000-000000000000"]
}
