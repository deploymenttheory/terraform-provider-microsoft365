# Token signing certificate configuration for unit testing
resource "microsoft365_graph_beta_applications_service_principal_token_signing_certificate" "test" {
  service_principal_id = "11111111-1111-1111-1111-111111111111"
  display_name         = "CN=Test SAML Signing"
  end_date_time        = "2029-01-01T00:00:00Z"
}
