resource "random_string" "app_name" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_applications_application" "test" {
  display_name     = "acc-test-sp-preferred-thumbprint-${random_string.app_name.result}"
  sign_in_audience = "AzureADMyOrg"
}

resource "microsoft365_graph_beta_applications_service_principal" "test" {
  app_id                        = microsoft365_graph_beta_applications_application.test.app_id
  preferred_single_sign_on_mode = "saml"
}

resource "microsoft365_graph_beta_applications_service_principal_preferred_token_signing_key_thumbprint" "test" {
  service_principal_id = microsoft365_graph_beta_applications_service_principal.test.id
  thumbprint           = "aabbccddeeff00112233445566778899aabbccdd"
}
