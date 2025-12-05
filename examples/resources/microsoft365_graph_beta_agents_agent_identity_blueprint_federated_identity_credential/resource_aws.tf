# Example: Create a federated identity credential for AWS

# First, create or reference an existing agent identity blueprint
resource "microsoft365_graph_beta_agents_agent_identity_blueprint" "example" {
  display_name     = "aws-workload-agent"
  sponsor_user_ids = ["00000000-0000-0000-0000-000000000000"]
  owner_user_ids   = ["00000000-0000-0000-0000-000000000000"]
  description      = "Agent identity blueprint for AWS workloads"
}

# Create a federated identity credential for AWS IAM
resource "microsoft365_graph_beta_agents_agent_identity_blueprint_federated_identity_credential" "aws_iam" {
  blueprint_id = microsoft365_graph_beta_agents_agent_identity_blueprint.example.id
  name         = "aws-iam-role"
  issuer       = "https://token.sts.amazonaws.com"
  subject      = "arn:aws:iam::123456789012:role/my-role"
  audiences    = ["api://AzureADTokenExchange"]
  description  = "Federated identity credential for AWS IAM role"
}

# Output the credential details
output "credential_id" {
  value       = microsoft365_graph_beta_agents_agent_identity_blueprint_federated_identity_credential.aws_iam.id
  description = "The ID of the federated identity credential"
}
