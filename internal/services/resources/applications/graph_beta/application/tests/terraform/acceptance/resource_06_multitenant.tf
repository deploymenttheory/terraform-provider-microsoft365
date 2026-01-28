# ==============================================================================
# Random Suffix for Unique Resource Names
# ==============================================================================

resource "random_string" "app_suffix" {
  length  = 8
  special = false
  upper   = false
}


# ==============================================================================
# Dependencies - Users for owners
# ==============================================================================

resource "microsoft365_graph_beta_users_user" "dependency_owner" {
  display_name        = "acc-test-mt-owner-${random_string.app_suffix.result}"
  user_principal_name = "acc-test-mt-owner-${random_string.app_suffix.result}@deploymenttheory.com"
  mail_nickname       = "acc-test-mt-owner-${random_string.app_suffix.result}"
  account_enabled     = true
  password_profile = {
    password                           = "SecureP@ssw0rd123!"
    force_change_password_next_sign_in = false
  }
  hard_delete = true
}

# ==============================================================================
# Application
# ==============================================================================

# APP006: Multitenant Application Configuration
# Tests multitenant application with personal Microsoft account support

resource "microsoft365_graph_beta_applications_application" "test_multitenant" {
  display_name     = "acc-test-multitenant-${random_string.app_suffix.result}"
  description      = "Multitenant acceptance test application"
  sign_in_audience = "AzureADandPersonalMicrosoftAccount"

  identifier_uris = [
    "api://acc-test-mt-${random_string.app_suffix.result}"
  ]

  web = {
    home_page_url = "https://contoso.com"
    redirect_uris = [
      "https://contoso.com/signin-oidc"
    ]
    implicit_grant_settings = {
      enable_access_token_issuance = false
      enable_id_token_issuance     = true
    }
    redirect_uri_settings = []
  }

  spa = {
    redirect_uris = [
      "https://contoso.com/spa"
    ]
  }

  required_resource_access = []

  tags = [
    "multitenant",
    "acceptance-test"
  ]

  owner_user_ids = [
    microsoft365_graph_beta_users_user.dependency_owner.id
  ]

  prevent_duplicate_names = false
  hard_delete             = true
}


