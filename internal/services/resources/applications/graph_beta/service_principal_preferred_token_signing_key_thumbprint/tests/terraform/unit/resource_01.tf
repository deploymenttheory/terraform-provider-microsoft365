# Preferred token signing key thumbprint configuration for unit testing
resource "microsoft365_graph_beta_applications_service_principal_preferred_token_signing_key_thumbprint" "test" {
  service_principal_id = "11111111-1111-1111-1111-111111111111"
  thumbprint           = "aabbccddeeff00112233445566778899aabbccdd"
}
