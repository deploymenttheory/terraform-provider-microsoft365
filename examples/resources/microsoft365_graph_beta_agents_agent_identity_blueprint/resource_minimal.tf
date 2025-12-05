# Minimal Agent Identity Blueprint configuration
# Creates an agent identity blueprint with required fields only
resource "microsoft365_graph_beta_agents_agent_identity_blueprint" "minimal" {
  display_name = "My Agent Blueprint"

  sponsor_user_ids = ["00000000-0000-0000-0000-000000000000"]
  owner_user_ids   = ["00000000-0000-0000-0000-000000000000"]
}

