resource "microsoft365_graph_beta_applications_application" "example" {
  display_name = "my-github-deployment-app"
  description  = "Application for GitHub Actions deployments"
}

# Federated credential for GitHub Actions to deploy to Azure
resource "microsoft365_graph_beta_applications_application_federated_identity_credential" "github_actions" {
  application_id = microsoft365_graph_beta_applications_application.example.id
  name           = "github-actions-production"
  description    = "GitHub Actions deploying to Production environment"
  issuer         = "https://token.actions.githubusercontent.com"
  subject        = "repo:myorg/myrepo:environment:Production"
  audiences      = ["api://AzureADTokenExchange"]
}
