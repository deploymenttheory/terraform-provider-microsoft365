# Maximal Password Credential configuration for unit testing
resource "microsoft365_graph_beta_applications_application_password_credential" "test_maximal" {
  application_id  = "22222222-2222-2222-2222-222222222222"
  display_name    = "unit-test-password-credential-maximal"
  start_date_time = "2026-01-01T00:00:00Z"
  end_date_time   = "2028-01-01T00:00:00Z"
}
