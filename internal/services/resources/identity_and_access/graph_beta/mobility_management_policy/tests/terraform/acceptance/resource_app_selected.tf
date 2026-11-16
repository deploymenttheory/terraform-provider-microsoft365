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
# Mobility Management Policy - MAM scope limited to the test group
# ==============================================================================

# Scopes the Microsoft Intune MAM user scope to the test group only. The test runs only when the
# original MAM scope is none or selected, and restores the original scope and groups afterwards.
resource "microsoft365_graph_beta_identity_and_access_mobility_management_policy" "app" {
  policy_type = "app"
  policy_id   = "0000000a-0000-0000-c000-000000000000"

  applies_to      = "selected"
  included_groups = [microsoft365_graph_beta_groups_group.scope.id]
}
