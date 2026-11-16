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

resource "microsoft365_graph_beta_groups_group" "scope" {
  display_name     = "acc-test-mobility-management-policy-${random_string.suffix.result}"
  mail_nickname    = "acc-test-mmp-${random_string.suffix.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Acceptance test group for the mobility management policy resource"
  hard_delete      = true
}

# ==============================================================================
# Mobility Management Policy - MAM scope restored to none
# ==============================================================================

# Restores the original none scope before the test group is deleted. Setting none removes the test group reference.
resource "microsoft365_graph_beta_identity_and_access_mobility_management_policy" "app" {
  policy_type = "app"
  policy_id   = "0000000a-0000-0000-c000-000000000000"

  applies_to = "none"
}
