# Agent User with Full Dependency Chain
# This example demonstrates the complete resource hierarchy:
# User (sponsor/owner) -> Agent Identity Blueprint -> Service Principal -> Agent Identity -> Agent User

########################################################################################
# Look up existing user to use as sponsor/owner
########################################################################################

data "microsoft365_graph_beta_users_user_by_filter" "sponsor" {
  filter_type  = "display_name"
  filter_value = "Admin User"
}

########################################################################################
# Agent Identity Blueprint
# The blueprint defines the template for agent identities
########################################################################################

resource "microsoft365_graph_beta_agents_agent_identity_blueprint" "example" {
  display_name     = "Example Agent Blueprint"
  description      = "Blueprint for example agent identity"
  sponsor_user_ids = [data.microsoft365_graph_beta_users_user_by_filter.sponsor.id]
  owner_user_ids   = [data.microsoft365_graph_beta_users_user_by_filter.sponsor.id]
  hard_delete      = true
}

########################################################################################
# Agent Identity Blueprint Service Principal
# Required before creating agent identities from the blueprint
########################################################################################

resource "microsoft365_graph_beta_agents_agent_identity_blueprint_service_principal" "example" {
  app_id      = microsoft365_graph_beta_agents_agent_identity_blueprint.example.app_id
  hard_delete = true
}

########################################################################################
# Agent Identity
# The parent identity that the agent user will be associated with
########################################################################################

resource "microsoft365_graph_beta_agents_agent_identity" "example" {
  display_name                = "Example Agent Identity"
  agent_identity_blueprint_id = microsoft365_graph_beta_agents_agent_identity_blueprint.example.app_id
  account_enabled             = true
  sponsor_ids                 = [data.microsoft365_graph_beta_users_user_by_filter.sponsor.id]
  owner_ids                   = [data.microsoft365_graph_beta_users_user_by_filter.sponsor.id]
  hard_delete                 = true

  depends_on = [
    microsoft365_graph_beta_agents_agent_identity_blueprint_service_principal.example
  ]
}

########################################################################################
# Agent User
# The user identity that authenticates through the parent agent identity
########################################################################################

resource "microsoft365_graph_beta_agents_agent_user" "example" {
  display_name        = "Example Agent User"
  agent_identity_id   = microsoft365_graph_beta_agents_agent_identity.example.id
  account_enabled     = true
  user_principal_name = "agent-user@${var.domain}"
  mail_nickname       = "agent-user"
  sponsor_ids         = [data.microsoft365_graph_beta_users_user_by_filter.sponsor.id]
  hard_delete         = true

  # Optional fields
  given_name     = "Agent"
  surname        = "User"
  job_title      = "AI Agent"
  department     = "Engineering"
  usage_location = "US"

  depends_on = [
    microsoft365_graph_beta_agents_agent_identity.example
  ]
}

########################################################################################
# Variables
########################################################################################

variable "domain" {
  description = "The verified domain for the tenant (e.g., contoso.com)"
  type        = string
}

