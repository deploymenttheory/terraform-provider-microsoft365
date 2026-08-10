# Service Principal configuration for unit testing - Maximal
resource "microsoft365_graph_beta_applications_service_principal" "test_maximal" {
  app_id                        = "22222222-2222-2222-2222-222222222222"
  account_enabled               = true
  app_role_assignment_required  = true
  description                   = "Maximal service principal configuration for testing"
  login_url                     = "https://login.example.com"
  notes                         = "Service principal for maximal unit testing"
  notification_email_addresses  = ["admin@example.com", "security@example.com"]
  preferred_single_sign_on_mode = "saml"
  tags                          = ["HideApp", "WindowsAzureActiveDirectoryIntegratedApp"]
  alternative_names             = ["isExplicit=True", "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test"]
  # References the key credential seeded by the mock; Graph rejects unknown key IDs
  token_encryption_key_id       = "dddddddd-1111-2222-3333-444444444444"
  hard_delete                   = true

  saml_single_sign_on_settings = {
    relay_state = "https://example.com/relay"
  }
}
