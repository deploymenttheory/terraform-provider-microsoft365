resource "random_string" "test_id" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_applications_application" "dependency_app" {
  display_name = "acc-test-pwd-cred-maximal-${random_string.test_id.result}"
  description  = "Application for password credential maximal test"
  hard_delete  = true
}

resource "microsoft365_graph_beta_applications_application_password_credential" "test_maximal" {
  application_id  = microsoft365_graph_beta_applications_application.dependency_app.id
  display_name    = "acc-test-password-credential-maximal-${random_string.test_id.result}"
  start_date_time = "2026-01-01T00:00:00Z"
  end_date_time   = "2028-01-01T00:00:00Z"
}
