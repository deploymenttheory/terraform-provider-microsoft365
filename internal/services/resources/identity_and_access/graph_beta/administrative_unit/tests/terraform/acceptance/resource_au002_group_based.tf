# ==============================================================================
# Random Suffix for Unique Resource Names
# ==============================================================================

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Group Dependencies
# ==============================================================================

resource "microsoft365_graph_beta_groups_group" "au002_test_group1" {
  display_name     = "AU002 Test Group 1 ${random_string.suffix.result}"
  mail_nickname    = "au002-testgroup1-${random_string.suffix.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Test group 1 for AU002 administrative unit"
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "au002_test_group2" {
  display_name     = "AU002 Test Group 2 ${random_string.suffix.result}"
  mail_nickname    = "au002-testgroup2-${random_string.suffix.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Test group 2 for AU002 administrative unit"
  hard_delete      = true
}

# ==============================================================================
# AU002: Group-Based Administrative Unit
# ==============================================================================

resource "microsoft365_graph_beta_identity_and_access_administrative_unit" "au002_group_based" {
  display_name = "acc-test-au002-group-based-${random_string.suffix.result}"
  description  = "Administrative unit for group-based testing"
  hard_delete  = true
}
