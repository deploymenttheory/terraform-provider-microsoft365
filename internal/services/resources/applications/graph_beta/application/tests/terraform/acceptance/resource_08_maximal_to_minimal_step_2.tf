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

resource "microsoft365_graph_beta_users_user" "dependency_owner_1" {
  display_name        = "acc-test-app-owner1-${random_string.app_suffix.result}"
  user_principal_name = "acc-test-app-owner1-${random_string.app_suffix.result}@deploymenttheory.com"
  mail_nickname       = "acc-test-app-owner1-${random_string.app_suffix.result}"
  account_enabled     = true
  password_profile = {
    password                           = "SecureP@ssw0rd123!"
    force_change_password_next_sign_in = false
  }
  hard_delete = true
}

# ==============================================================================
# Time Sleep for Dependency Consistency
# ==============================================================================

resource "time_sleep" "wait_for_dependencies" {
  create_duration = "30s"

  depends_on = [
    microsoft365_graph_beta_users_user.dependency_owner_1
  ]
}

# ==============================================================================
# Application
# ==============================================================================

# APP008 Step 2: Minimal Application Configuration
# Tests application update from maximal to minimal configuration

resource "microsoft365_graph_beta_applications_application" "test_maximal_to_minimal" {
  display_name = "acc-test-app-max-to-min-${random_string.app_suffix.result}"
  description  = "Maximal to minimal test application - step 2"

  owner_user_ids = [
    microsoft365_graph_beta_users_user.dependency_owner_1.id
  ]

  prevent_duplicate_names = false
  hard_delete             = true

  depends_on = [
    time_sleep.wait_for_dependencies
  ]
}


