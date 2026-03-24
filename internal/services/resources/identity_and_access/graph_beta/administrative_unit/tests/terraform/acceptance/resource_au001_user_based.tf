# ==============================================================================
# Random Suffix for Unique Resource Names
# ==============================================================================

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# User Dependencies
# ==============================================================================

resource "microsoft365_graph_beta_users_user" "au001_test_user1" {
  user_principal_name = "au001-testuser1-${random_string.suffix.result}@deploymenttheory.com"
  display_name        = "AU001 Test User 1"
  mail_nickname       = "au001-testuser1-${random_string.suffix.result}"
  account_enabled     = true
  password_profile = {
    password                           = "P@ssw0rd!${random_string.suffix.result}"
    force_change_password_next_sign_in = false
  }
  hard_delete = true
}

resource "microsoft365_graph_beta_users_user" "au001_test_user2" {
  user_principal_name = "au001-testuser2-${random_string.suffix.result}@deploymenttheory.com"
  display_name        = "AU001 Test User 2"
  mail_nickname       = "au001-testuser2-${random_string.suffix.result}"
  account_enabled     = true
  password_profile = {
    password                           = "P@ssw0rd!${random_string.suffix.result}"
    force_change_password_next_sign_in = false
  }
  hard_delete = true
}

# ==============================================================================
# AU001: User-Based Administrative Unit
# ==============================================================================

resource "microsoft365_graph_beta_identity_and_access_administrative_unit" "au001_user_based" {
  display_name = "acc-test-au001-user-based-${random_string.suffix.result}"
  description  = "Administrative unit for user-based testing"
  hard_delete  = true
}
