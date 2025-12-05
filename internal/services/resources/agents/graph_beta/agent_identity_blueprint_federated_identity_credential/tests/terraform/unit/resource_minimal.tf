# Minimal Federated Identity Credential configuration for unit testing
resource "microsoft365_graph_beta_agents_agent_identity_blueprint_federated_identity_credential" "test_minimal" {
  blueprint_id = "11111111-1111-1111-1111-111111111111"
  name         = "unit-test-fic-minimal"
  issuer       = "https://token.actions.githubusercontent.com"
  subject      = "repo:octo-org/octo-repo:environment:Production"
  audiences    = ["api://AzureADTokenExchange"]
}
