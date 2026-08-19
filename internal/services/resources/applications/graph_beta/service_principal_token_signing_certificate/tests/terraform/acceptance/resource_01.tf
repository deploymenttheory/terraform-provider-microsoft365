resource "random_string" "app_name" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_applications_application" "test" {
  display_name     = "acc-test-sp-token-signing-cert-${random_string.app_name.result}"
  sign_in_audience = "AzureADMyOrg"
}

resource "microsoft365_graph_beta_applications_service_principal" "test" {
  app_id                        = microsoft365_graph_beta_applications_application.test.app_id
  preferred_single_sign_on_mode = "saml"
}

resource "microsoft365_graph_beta_applications_service_principal_token_signing_certificate" "test" {
  service_principal_id = microsoft365_graph_beta_applications_service_principal.test.id
  display_name         = "CN=Acceptance Test SAML Signing"
  end_date_time        = "2028-01-01T14:59:59Z"

  lifecycle {
    create_before_destroy = true
  }
}
