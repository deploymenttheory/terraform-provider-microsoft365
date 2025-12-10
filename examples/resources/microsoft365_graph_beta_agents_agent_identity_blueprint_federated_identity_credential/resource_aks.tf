# Example: Create a federated identity credential for Azure Kubernetes Service (AKS)

# First, create or reference an existing agent identity blueprint
resource "microsoft365_graph_beta_agents_agent_identity_blueprint" "example" {
  display_name     = "aks-workload-agent"
  sponsor_user_ids = ["00000000-0000-0000-0000-000000000000"]
  owner_user_ids   = ["00000000-0000-0000-0000-000000000000"]
  description      = "Agent identity blueprint for AKS workloads"
  hard_delete      = true
}

# Create a federated identity credential for AKS workload identity
resource "microsoft365_graph_beta_agents_agent_identity_blueprint_federated_identity_credential" "aks_workload" {
  blueprint_id = microsoft365_graph_beta_agents_agent_identity_blueprint.example.id
  name         = "aks-workload-identity"
  issuer       = "https://oidc.prod-aks.azure.com/00000000-0000-0000-0000-000000000000/"
  subject      = "system:serviceaccount:default:workload-identity-sa"
  audiences    = ["api://AzureADTokenExchange"]
  description  = "Federated identity credential for AKS workload identity"
}

# Output the credential details
output "credential_id" {
  value       = microsoft365_graph_beta_agents_agent_identity_blueprint_federated_identity_credential.aks_workload.id
  description = "The ID of the federated identity credential"
}
