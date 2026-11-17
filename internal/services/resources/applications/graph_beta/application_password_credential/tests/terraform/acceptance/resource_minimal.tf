# Acceptance test: Minimal Application Password Credential configuration
# Dependency chain: random_string -> application -> password_credential

resource "random_string" "test_id" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_applications_application" "dependency_app" {
  display_name = "acc-test-pwd-cred-app-${random_string.test_id.result}"
  description  = "Application for password credential acceptance test"
  hard_delete  = true
}

resource "microsoft365_graph_beta_applications_application_password_credential" "test_minimal" {
  application_id = microsoft365_graph_beta_applications_application.dependency_app.id
  display_name   = "acc-test-password-credential-${random_string.test_id.result}"
}
