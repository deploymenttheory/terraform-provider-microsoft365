# Unit test: Minimal Application Password Credential configuration

resource "microsoft365_graph_beta_applications_application_password_credential" "test_minimal" {
  application_id = "11111111-1111-1111-1111-111111111111"
  display_name   = "unit-test-password-credential"
}
