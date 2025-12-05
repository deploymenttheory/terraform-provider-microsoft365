# Example: Configure an identifier URI for an Agent Identity Blueprint

# First, create or reference an existing agent identity blueprint
resource "microsoft365_graph_beta_agents_agent_identity_blueprint" "example" {
  display_name     = "my-agent-blueprint"
  sponsor_user_ids = ["00000000-0000-0000-0000-000000000000"]
  owner_user_ids   = ["00000000-0000-0000-0000-000000000000"]
  description      = "Agent identity blueprint for automated workflows"
}

# Configure the identifier URI using the blueprint's ID
resource "microsoft365_graph_beta_agents_agent_identity_blueprint_identifier_uri" "example" {
  blueprint_id   = microsoft365_graph_beta_agents_agent_identity_blueprint.example.id
  identifier_uri = "api://${microsoft365_graph_beta_agents_agent_identity_blueprint.example.id}"
}

output "identifier_uri" {
  value       = microsoft365_graph_beta_agents_agent_identity_blueprint_identifier_uri.example.identifier_uri
  description = "The configured identifier URI"
}

