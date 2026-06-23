# Acceptance test: Maximal Application Federated Identity Credential configuration
# Full dependency chain: random_string -> application -> federated_identity_credential

resource "random_string" "test_id_maximal" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_applications_application" "test_app_maximal" {
  display_name = "acc-test-app-fic-maximal-${random_string.test_id_maximal.result}"
  description  = "Application for federated identity credential maximal acceptance test"
  hard_delete  = true
}

resource "time_sleep" "wait_for_app_maximal" {
  depends_on      = [microsoft365_graph_beta_applications_application.test_app_maximal]
  create_duration = "30s"
}

# Federated credential scenario - GitHub Actions deploying Azure resources with all optional fields
resource "microsoft365_graph_beta_applications_application_federated_identity_credential" "test_maximal" {
  application_id = microsoft365_graph_beta_applications_application.test_app_maximal.id
  name           = "acc-test-fic-maximal-${random_string.test_id_maximal.result}"
  description    = "Federated credential scenario - GitHub Actions with all optional fields configured"
  issuer         = "https://token.actions.githubusercontent.com"
  subject        = "repo:deploymenttheory/test-repo-${random_string.test_id_maximal.result}:environment:Production"
  audiences      = ["api://AzureADTokenExchange"]

  depends_on = [time_sleep.wait_for_app_maximal]
}
