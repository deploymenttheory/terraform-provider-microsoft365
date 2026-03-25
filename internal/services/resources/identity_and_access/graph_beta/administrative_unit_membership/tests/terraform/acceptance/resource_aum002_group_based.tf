# ==============================================================================
# Random Suffix for Unique Resource Names
# ==============================================================================

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Test Group
# ==============================================================================

resource "microsoft365_graph_beta_groups_group" "test_group" {
  display_name     = "acc-test-aum002-group-${random_string.suffix.result}"
  mail_nickname    = "aum002-group-${random_string.suffix.result}"
  mail_enabled     = false
  security_enabled = true
  hard_delete      = true
}

# ==============================================================================
# Administrative Unit
# ==============================================================================

resource "microsoft365_graph_beta_identity_and_access_administrative_unit" "test_au" {
  display_name = "acc-test-aum002-au-${random_string.suffix.result}"
  description  = "Administrative unit for group membership testing"
  hard_delete  = true
}

# ==============================================================================
# AUM002: Group-Based Membership
# ==============================================================================

resource "time_sleep" "wait_for_dependencies" {
  create_duration = "30s"

  depends_on = [
    microsoft365_graph_beta_groups_group.test_group,
    microsoft365_graph_beta_identity_and_access_administrative_unit.test_au,
  ]
}

resource "microsoft365_graph_beta_identity_and_access_administrative_unit_membership" "aum002_group_based" {
  administrative_unit_id = microsoft365_graph_beta_identity_and_access_administrative_unit.test_au.id
  members = [
    microsoft365_graph_beta_groups_group.test_group.id
  ]

  depends_on = [time_sleep.wait_for_dependencies]
}
