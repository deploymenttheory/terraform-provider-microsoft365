# Example: Create a password credential for an Agent Identity Blueprint

# First, create or reference an existing agent identity blueprint
resource "microsoft365_graph_beta_agents_agent_identity_blueprint" "example" {
  display_name     = "my-agent-blueprint"
  sponsor_user_ids = ["00000000-0000-0000-0000-000000000000"]
  owner_user_ids   = ["00000000-0000-0000-0000-000000000000"]
  description      = "Agent identity blueprint for automated workflows"
}

# Create a password credential for the blueprint
resource "microsoft365_graph_beta_agents_agent_identity_blueprint_password_credential" "example" {
  blueprint_id = microsoft365_graph_beta_agents_agent_identity_blueprint.example.id
  display_name = "api-access-credential"
}

# IMPORTANT: Store the secret securely - it is only available at creation time
# You can use a secrets manager or output to a secure location
output "client_secret" {
  value       = microsoft365_graph_beta_agents_agent_identity_blueprint_password_credential.example.secret_text
  description = "The generated client secret - store this securely!"
  sensitive   = true
}

output "key_id" {
  value       = microsoft365_graph_beta_agents_agent_identity_blueprint_password_credential.example.key_id
  description = "The key ID of the password credential"
}

