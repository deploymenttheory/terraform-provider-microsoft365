# Example: Create a password credential with custom validity period

# First, create or reference an existing agent identity blueprint
resource "microsoft365_graph_beta_agents_agent_identity_blueprint" "example" {
  display_name     = "my-agent-blueprint"
  sponsor_user_ids = ["00000000-0000-0000-0000-000000000000"]
  owner_user_ids   = ["00000000-0000-0000-0000-000000000000"]
  description      = "Agent identity blueprint for automated workflows"
  hard_delete      = true
}

# Create a password credential with specific start and end dates
resource "microsoft365_graph_beta_agents_agent_identity_blueprint_password_credential" "example" {
  blueprint_id    = microsoft365_graph_beta_agents_agent_identity_blueprint.example.id
  display_name    = "short-lived-credential"
  start_date_time = "2025-01-01T00:00:00Z"
  end_date_time   = "2025-06-30T23:59:59Z" # 6 month validity
}

# IMPORTANT: Store the secret securely - it is only available at creation time
output "client_secret" {
  value       = microsoft365_graph_beta_agents_agent_identity_blueprint_password_credential.example.secret_text
  description = "The generated client secret - store this securely!"
  sensitive   = true
}

output "credential_info" {
  value = {
    key_id          = microsoft365_graph_beta_agents_agent_identity_blueprint_password_credential.example.key_id
    hint            = microsoft365_graph_beta_agents_agent_identity_blueprint_password_credential.example.hint
    start_date_time = microsoft365_graph_beta_agents_agent_identity_blueprint_password_credential.example.start_date_time
    end_date_time   = microsoft365_graph_beta_agents_agent_identity_blueprint_password_credential.example.end_date_time
  }
  description = "Password credential metadata"
}

