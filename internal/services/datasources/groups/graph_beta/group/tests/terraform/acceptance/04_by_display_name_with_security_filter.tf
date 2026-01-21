# Test 04: Look up group by display_name with security_enabled filter
# This test creates a security group and then looks it up using display name with security filter

# ==============================================================================
# Random Suffix for Unique Resource Names
# ==============================================================================

resource "random_string" "test" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Group Resource
# ==============================================================================

resource "microsoft365_graph_beta_groups_group" "test" {
  display_name     = "acc-test-by-display-name-security-${random_string.test.result}"
  mail_nickname    = "acctestbydisplaysec${random_string.test.result}"
  description      = "Test security group for data source acceptance test with security filter"
  mail_enabled     = false
  security_enabled = true
  hard_delete      = true

}

# ==============================================================================
# Wait for Group Propagation
# ==============================================================================

resource "time_sleep" "wait_for_group" {
  depends_on      = [microsoft365_graph_beta_groups_group.test]
  create_duration = "30s"
}

# ==============================================================================
# Data Source - Lookup by display_name with security_enabled filter
# ==============================================================================

data "microsoft365_graph_beta_groups_group" "test" {
  display_name     = microsoft365_graph_beta_groups_group.test.display_name
  security_enabled = true

  depends_on = [time_sleep.wait_for_group]
}

# ==============================================================================
# Outputs
# ==============================================================================

output "group_id" {
  value = data.microsoft365_graph_beta_groups_group.test.id
}

output "display_name" {
  value = data.microsoft365_graph_beta_groups_group.test.display_name
}

output "security_enabled" {
  value = data.microsoft365_graph_beta_groups_group.test.security_enabled
}
