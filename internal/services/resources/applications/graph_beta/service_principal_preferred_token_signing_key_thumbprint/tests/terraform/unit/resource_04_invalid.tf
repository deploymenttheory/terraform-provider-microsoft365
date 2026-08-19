# Invalid thumbprint format: must fail schema validation at plan time
resource "microsoft365_graph_beta_applications_service_principal_preferred_token_signing_key_thumbprint" "test" {
  service_principal_id = "11111111-1111-1111-1111-111111111111"
  thumbprint           = "aa:bb:cc:dd:ee:ff:00:11:22:33:44:55:66:77:88:99:aa:bb:cc:dd"
}
