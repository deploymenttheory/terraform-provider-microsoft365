# Token signing certificate with API defaults for display name and end date
resource "microsoft365_graph_beta_applications_service_principal_token_signing_certificate" "test" {
  service_principal_id = "11111111-1111-1111-1111-111111111111"
}
