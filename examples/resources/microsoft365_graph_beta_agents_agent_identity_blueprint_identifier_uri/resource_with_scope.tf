# Example: Configure an identifier URI with custom OAuth2 permission scope

# First, create or reference an existing agent identity blueprint
resource "microsoft365_graph_beta_agents_agent_identity_blueprint" "example" {
  display_name     = "my-agent-blueprint"
  sponsor_user_ids = ["00000000-0000-0000-0000-000000000000"]
  owner_user_ids   = ["00000000-0000-0000-0000-000000000000"]
  description      = "Agent identity blueprint for automated workflows"
  hard_delete      = true
}

# Configure the identifier URI with a custom permission scope
resource "microsoft365_graph_beta_agents_agent_identity_blueprint_identifier_uri" "example" {
  blueprint_id   = microsoft365_graph_beta_agents_agent_identity_blueprint.example.id
  identifier_uri = "api://${microsoft365_graph_beta_agents_agent_identity_blueprint.example.id}"

  scope = {
    admin_consent_description  = "Allow the application to access the agent on behalf of the signed-in user."
    admin_consent_display_name = "Access agent"
    is_enabled                 = true
    type                       = "User"
    value                      = "access_agent"
  }
}

output "identifier_uri" {
  value       = microsoft365_graph_beta_agents_agent_identity_blueprint_identifier_uri.example.identifier_uri
  description = "The configured identifier URI"
}

output "scope_id" {
  value       = microsoft365_graph_beta_agents_agent_identity_blueprint_identifier_uri.example.scope.id
  description = "The ID of the OAuth2 permission scope"
}

