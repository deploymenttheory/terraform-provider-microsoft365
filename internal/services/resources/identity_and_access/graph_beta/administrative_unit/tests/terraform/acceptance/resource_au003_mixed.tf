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

resource "microsoft365_graph_beta_users_user" "au003_test_user" {
  user_principal_name = "au003-testuser-${random_string.suffix.result}@deploymenttheory.com"
  display_name        = "AU003 Test User"
  mail_nickname       = "au003-testuser-${random_string.suffix.result}"
  account_enabled     = true
  password_profile = {
    password                           = "P@ssw0rd!${random_string.suffix.result}"
    force_change_password_next_sign_in = false
  }
  hard_delete = true
}

# ==============================================================================
# Group Dependencies
# ==============================================================================

resource "microsoft365_graph_beta_groups_group" "au003_test_group" {
  display_name     = "AU003 Test Group ${random_string.suffix.result}"
  mail_nickname    = "au003-testgroup-${random_string.suffix.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Test group for AU003 administrative unit"
  hard_delete      = true
}

# ==============================================================================
# AU003: Mixed User and Group Administrative Unit
# ==============================================================================

resource "microsoft365_graph_beta_identity_and_access_administrative_unit" "au003_mixed" {
  display_name = "acc-test-au003-mixed-${random_string.suffix.result}"
  description  = "Administrative unit for mixed user and group testing"
  visibility   = "HiddenMembership"
  hard_delete  = true
}
