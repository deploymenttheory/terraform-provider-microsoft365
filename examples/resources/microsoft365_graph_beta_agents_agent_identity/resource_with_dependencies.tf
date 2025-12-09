# Example: Agent Identity with Full Dependency Chain
#
# This example shows the complete setup including:
# - A user to act as sponsor and owner
# - An Agent Identity Blueprint
# - The Blueprint's Service Principal
# - The Agent Identity itself

# Create a user to be the sponsor
resource "microsoft365_graph_beta_users_user" "agent_sponsor" {
  display_name        = "Agent Sponsor User"
  user_principal_name = "agent-sponsor@yourdomain.com"
  mail_nickname       = "agent-sponsor"
  account_enabled     = true
  password_profile = {
    password                           = "SecureP@ssw0rd123!"
    force_change_password_next_sign_in = true
  }
}

# Create a user to be the owner
resource "microsoft365_graph_beta_users_user" "agent_owner" {
  display_name        = "Agent Owner User"
  user_principal_name = "agent-owner@yourdomain.com"
  mail_nickname       = "agent-owner"
  account_enabled     = true
  password_profile = {
    password                           = "SecureP@ssw0rd123!"
    force_change_password_next_sign_in = true
  }
}

# Create an agent identity blueprint
resource "microsoft365_graph_beta_agents_agent_identity_blueprint" "example" {
  display_name     = "Customer Service Agent Blueprint"
  description      = "Blueprint for customer service AI agents"
  sponsor_user_ids = [microsoft365_graph_beta_users_user.agent_sponsor.id]
  owner_user_ids   = [microsoft365_graph_beta_users_user.agent_owner.id]
  tags             = ["customer-service", "production"]
}

# Create the service principal for the blueprint (required before creating agent identities)
resource "microsoft365_graph_beta_agents_agent_identity_blueprint_service_principal" "example" {
  app_id = microsoft365_graph_beta_agents_agent_identity_blueprint.example.app_id
}

# Create an agent identity from the blueprint
resource "microsoft365_graph_beta_agents_agent_identity" "example" {
  display_name                = "Customer Service Agent 01"
  agent_identity_blueprint_id = microsoft365_graph_beta_agents_agent_identity_blueprint.example.app_id
  account_enabled             = true
  sponsor_ids                 = [microsoft365_graph_beta_users_user.agent_sponsor.id]
  owner_ids                   = [microsoft365_graph_beta_users_user.agent_owner.id]
  tags                        = ["customer-service", "agent-instance"]

  depends_on = [
    microsoft365_graph_beta_agents_agent_identity_blueprint_service_principal.example
  ]
}

# Outputs
output "agent_identity_id" {
  description = "The ID of the created agent identity"
  value       = microsoft365_graph_beta_agents_agent_identity.example.id
}

output "agent_identity_display_name" {
  description = "The display name of the agent identity"
  value       = microsoft365_graph_beta_agents_agent_identity.example.display_name
}

output "agent_identity_service_principal_type" {
  description = "The service principal type of the agent identity"
  value       = microsoft365_graph_beta_agents_agent_identity.example.service_principal_type
}

output "blueprint_app_id" {
  description = "The app ID of the agent identity blueprint"
  value       = microsoft365_graph_beta_agents_agent_identity_blueprint.example.app_id
}

