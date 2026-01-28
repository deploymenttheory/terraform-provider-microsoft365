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
  display_name        = "acc-test-web-owner-${random_string.app_suffix.result}"
  user_principal_name = "acc-test-web-owner-${random_string.app_suffix.result}@deploymenttheory.com"
  mail_nickname       = "acc-test-web-owner-${random_string.app_suffix.result}"
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

# APP003: Web Application Configuration
# Tests typical web application configuration with OIDC authentication

resource "microsoft365_graph_beta_applications_application" "test_web_app" {
  display_name     = "acc-test-web-app-${random_string.app_suffix.result}"
  description      = "Web application acceptance test"
  sign_in_audience = "AzureADMyOrg"

  identifier_uris = [
    "https://acc-test-web-${random_string.app_suffix.result}.azurewebsites.net"
  ]

  web = {
    home_page_url = "https://acc-test-web-${random_string.app_suffix.result}.azurewebsites.net"
    logout_url    = "https://acc-test-web-${random_string.app_suffix.result}.azurewebsites.net/signout"
    redirect_uris = [
      "https://acc-test-web-${random_string.app_suffix.result}.azurewebsites.net/signin-oidc"
    ]
    implicit_grant_settings = {
      enable_access_token_issuance = false
      enable_id_token_issuance     = true
    }
    redirect_uri_settings = []
  }

  required_resource_access = []

  owner_user_ids = [
    microsoft365_graph_beta_users_user.dependency_owner.id
  ]

  prevent_duplicate_names = false
  hard_delete             = true
}


