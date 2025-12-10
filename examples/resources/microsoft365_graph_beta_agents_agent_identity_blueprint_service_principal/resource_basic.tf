# Example: Create a service principal for an agent identity blueprint

# First, create or reference an existing agent identity blueprint
resource "microsoft365_graph_beta_agents_agent_identity_blueprint" "example" {
  display_name     = "example-agent-identity-blueprint"
  sponsor_user_ids = ["00000000-0000-0000-0000-000000000000"]
  owner_user_ids   = ["00000000-0000-0000-0000-000000000000"]
  description      = "Example agent identity blueprint for service principal creation"
  hard_delete      = true
}

# Create a service principal for the blueprint
resource "microsoft365_graph_beta_agents_agent_identity_blueprint_service_principal" "example" {
  app_id      = microsoft365_graph_beta_agents_agent_identity_blueprint.example.app_id
  hard_delete = true
}

# Output the service principal ID
output "service_principal_id" {
  value       = microsoft365_graph_beta_agents_agent_identity_blueprint_service_principal.example.id
  description = "The Object ID of the created service principal"
}
