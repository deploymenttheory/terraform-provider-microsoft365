# Acceptance test: Service Principal - Maximal
# Full dependency chain: random_string -> application -> service_principal

resource "random_string" "test_id_max" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_applications_application" "test_max" {
  display_name = "acc-test-sp-maximal-${random_string.test_id_max.result}"
  description  = "Application for service principal maximal acceptance test"
  hard_delete  = true
}

resource "time_sleep" "wait_for_app_max" {
  depends_on      = [microsoft365_graph_beta_applications_application.test_max]
  create_duration = "15s"
}

resource "microsoft365_graph_beta_applications_service_principal" "test_maximal" {
  app_id                        = microsoft365_graph_beta_applications_application.test_max.app_id
  account_enabled               = true
  app_role_assignment_required  = true
  description                   = "Maximal service principal configuration for testing"
  login_url                     = "https://login.example.com"
  notes                         = "Service principal for maximal acceptance testing"
  notification_email_addresses  = ["admin@example.com", "security@example.com"]
  preferred_single_sign_on_mode = "saml"
  tags                          = ["HideApp", "WindowsAzureActiveDirectoryIntegratedApp"]
  alternative_names             = ["isExplicit=True", "acc-test-sp-maximal"]
  hard_delete                   = true

  # token_encryption_key_id is not exercised here: Graph rejects it unless the
  # referenced key credential already exists on the service principal.
  saml_single_sign_on_settings = {
    relay_state = "https://example.com/relay"
  }

  depends_on = [time_sleep.wait_for_app_max]
}
