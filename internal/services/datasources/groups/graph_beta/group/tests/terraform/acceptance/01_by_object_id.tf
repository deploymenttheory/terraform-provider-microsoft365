# Test 01: Look up group by object_id
# This test creates a security group and then looks it up using its object_id

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
  display_name     = "acc-test-by-object-id-${random_string.test.result}"
  mail_nickname    = "acctestbyobjectid${random_string.test.result}"
  description      = "Test group for data source acceptance test"
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
# Data Source - Lookup by object_id
# ==============================================================================

data "microsoft365_graph_beta_groups_group" "test" {
  object_id = microsoft365_graph_beta_groups_group.test.id

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
