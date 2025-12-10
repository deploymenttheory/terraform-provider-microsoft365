# Example: Create a federated identity credential for GitHub Actions

# First, create or reference an existing agent identity blueprint
resource "microsoft365_graph_beta_agents_agent_identity_blueprint" "example" {
  display_name     = "github-actions-agent"
  sponsor_user_ids = ["00000000-0000-0000-0000-000000000000"]
  owner_user_ids   = ["00000000-0000-0000-0000-000000000000"]
  description      = "Agent identity blueprint for GitHub Actions workflows"
  hard_delete      = true
}

# Create a federated identity credential for GitHub Actions
resource "microsoft365_graph_beta_agents_agent_identity_blueprint_federated_identity_credential" "github_actions" {
  blueprint_id = microsoft365_graph_beta_agents_agent_identity_blueprint.example.id
  name         = "github-actions-production"
  issuer       = "https://token.actions.githubusercontent.com"
  subject      = "repo:my-org/my-repo:environment:Production"
  audiences    = ["api://AzureADTokenExchange"]
  description  = "Federated identity credential for GitHub Actions in production environment"
}

# Output the credential details
output "credential_id" {
  value       = microsoft365_graph_beta_agents_agent_identity_blueprint_federated_identity_credential.github_actions.id
  description = "The ID of the federated identity credential"
}
