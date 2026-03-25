# ==============================================================================
# Random Suffix for Unique Resource Names
# ==============================================================================

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Test Users
# ==============================================================================

resource "microsoft365_graph_beta_users_user" "test_user1" {
  user_principal_name = "aum001-user1-${random_string.suffix.result}@deploymenttheory.com"
  display_name        = "AUM001 Test User 1"
  mail_nickname       = "aum001-user1-${random_string.suffix.result}"
  account_enabled     = true
  password_profile = {
    password                           = "P@ssw0rd!${random_string.suffix.result}"
    force_change_password_next_sign_in = false
  }
  hard_delete = true
}

resource "microsoft365_graph_beta_users_user" "test_user2" {
  user_principal_name = "aum001-user2-${random_string.suffix.result}@deploymenttheory.com"
  display_name        = "AUM001 Test User 2"
  mail_nickname       = "aum001-user2-${random_string.suffix.result}"
  account_enabled     = true
  password_profile = {
    password                           = "P@ssw0rd!${random_string.suffix.result}"
    force_change_password_next_sign_in = false
  }
  hard_delete = true
}

# ==============================================================================
# Administrative Unit
# ==============================================================================

resource "microsoft365_graph_beta_identity_and_access_administrative_unit" "test_au" {
  display_name = "acc-test-aum001-au-${random_string.suffix.result}"
  description  = "Administrative unit for membership testing"
  hard_delete  = true
}

# ==============================================================================
# AUM001: User-Based Membership
# ==============================================================================

resource "time_sleep" "wait_for_dependencies" {
  create_duration = "30s"

  depends_on = [
    microsoft365_graph_beta_users_user.test_user1,
    microsoft365_graph_beta_users_user.test_user2,
    microsoft365_graph_beta_identity_and_access_administrative_unit.test_au,
  ]
}

resource "microsoft365_graph_beta_identity_and_access_administrative_unit_membership" "aum001_user_based" {
  administrative_unit_id = microsoft365_graph_beta_identity_and_access_administrative_unit.test_au.id
  members = [
    microsoft365_graph_beta_users_user.test_user1.id,
    microsoft365_graph_beta_users_user.test_user2.id
  ]

  depends_on = [time_sleep.wait_for_dependencies]
}
