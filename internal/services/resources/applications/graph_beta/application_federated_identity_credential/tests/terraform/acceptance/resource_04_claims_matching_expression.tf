# Acceptance test: Federated Identity Credential using a flexible-FIC claims matching expression
# Full dependency chain: random_string -> application -> federated_identity_credential
# claims_matching_expression is mutually exclusive with subject, so subject is omitted.

resource "random_string" "test_id_claims" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_applications_application" "test_app_claims" {
  display_name = "acc-test-app-fic-claims-${random_string.test_id_claims.result}"
  description  = "Application for federated identity credential claims matching acceptance test"
  hard_delete  = true
}

resource "time_sleep" "wait_for_app_claims" {
  depends_on      = [microsoft365_graph_beta_applications_application.test_app_claims]
  create_duration = "30s"
}

# Federated credential scenario - GitHub Actions matched via a claims matching expression
resource "microsoft365_graph_beta_applications_application_federated_identity_credential" "test_claims_matching" {
  application_id             = microsoft365_graph_beta_applications_application.test_app_claims.id
  name                       = "acc-test-fic-claims-${random_string.test_id_claims.result}"
  description                = "Federated credential scenario - flexible FIC claims matching expression"
  issuer                     = "https://token.actions.githubusercontent.com"
  claims_matching_expression = "claims['sub'] matches 'repo:deploymenttheory/*'"
  audiences                  = ["api://AzureADTokenExchange"]

  depends_on = [time_sleep.wait_for_app_claims]
}
