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
  display_name        = "acc-test-spa-owner-${random_string.app_suffix.result}"
  user_principal_name = "acc-test-spa-owner-${random_string.app_suffix.result}@deploymenttheory.com"
  mail_nickname       = "acc-test-spa-owner-${random_string.app_suffix.result}"
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

# APP004: Single Page Application (SPA) Configuration
# Tests SPA configuration with multiple redirect URIs

resource "microsoft365_graph_beta_applications_application" "test_spa" {
  display_name     = "acc-test-spa-${random_string.app_suffix.result}"
  description      = "SPA acceptance test application"
  sign_in_audience = "AzureADMultipleOrgs"

  identifier_uris = [
    "api://acc-test-spa-${random_string.app_suffix.result}"
  ]

  spa = {
    redirect_uris = [
      "http://localhost:3000",
      "https://acc-test-spa-${random_string.app_suffix.result}.azurestaticapps.net"
    ]
  }

  required_resource_access = []

  owner_user_ids = [
    microsoft365_graph_beta_users_user.dependency_owner.id
  ]

  prevent_duplicate_names = false
  hard_delete             = true
}


