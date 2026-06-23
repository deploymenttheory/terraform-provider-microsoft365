# Acceptance test: Minimal Application Federated Identity Credential configuration
# Full dependency chain: random_string -> application -> federated_identity_credential

resource "random_string" "test_id" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_applications_application" "test_app_minimal" {
  display_name = "acc-test-app-fic-${random_string.test_id.result}"
  description  = "Application for federated identity credential acceptance test"
  hard_delete  = true
}

resource "time_sleep" "wait_for_app_minimal" {
  depends_on      = [microsoft365_graph_beta_applications_application.test_app_minimal]
  create_duration = "30s"
}

# Federated credential scenario - GitHub Actions deploying Azure resources
resource "microsoft365_graph_beta_applications_application_federated_identity_credential" "test_minimal" {
  application_id = microsoft365_graph_beta_applications_application.test_app_minimal.id
  name           = "acc-test-fic-minimal-${random_string.test_id.result}"
  description    = "Federated credential scenario - GitHub Actions deploying Azure resources"
  issuer         = "https://token.actions.githubusercontent.com"
  subject        = "repo:deploymenttheory/test-repo-${random_string.test_id.result}:environment:Production"
  audiences      = ["api://AzureADTokenExchange"]
  
  depends_on = [time_sleep.wait_for_app_minimal]
}
