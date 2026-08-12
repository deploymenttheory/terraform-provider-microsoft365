# Uppercase thumbprint: Graph normalizes to lowercase, state must keep the configured casing
resource "microsoft365_graph_beta_applications_service_principal_preferred_token_signing_key_thumbprint" "test" {
  service_principal_id = "11111111-1111-1111-1111-111111111111"
  thumbprint           = "AABBCCDDEEFF00112233445566778899AABBCCDD"
}
