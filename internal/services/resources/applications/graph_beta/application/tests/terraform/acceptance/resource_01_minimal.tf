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
  display_name        = "acc-test-app-owner-${random_string.app_suffix.result}"
  user_principal_name = "acc-test-app-owner-${random_string.app_suffix.result}@deploymenttheory.com"
  mail_nickname       = "acc-test-app-owner-${random_string.app_suffix.result}"
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

# APP001: Minimal Application Configuration
# Tests application creation with only the required fields

resource "microsoft365_graph_beta_applications_application" "test_minimal" {
  display_name = "acc-test-app-minimal-${random_string.app_suffix.result}"
  description  = "Minimal acceptance test application"

  owner_user_ids = [
    microsoft365_graph_beta_users_user.dependency_owner.id
  ]

  prevent_duplicate_names = false
  hard_delete             = true
}


