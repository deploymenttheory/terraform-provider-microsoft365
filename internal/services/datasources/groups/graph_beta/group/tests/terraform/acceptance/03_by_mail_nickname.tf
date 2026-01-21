# Test 03: Look up group by mail_nickname
# This test creates a group and then looks it up using its mail_nickname

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
  display_name     = "acc-test-by-mail-nickname-${random_string.test.result}"
  mail_nickname    = "acctestbymailnickname${random_string.test.result}"
  description      = "Test group for mail nickname lookup"
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
# Data Source - Lookup by mail_nickname
# ==============================================================================

data "microsoft365_graph_beta_groups_group" "test" {
  mail_nickname = microsoft365_graph_beta_groups_group.test.mail_nickname

  depends_on = [time_sleep.wait_for_group]
}

# ==============================================================================
# Outputs
# ==============================================================================

output "group_id" {
  value = data.microsoft365_graph_beta_groups_group.test.id
}

output "mail_nickname" {
  value = data.microsoft365_graph_beta_groups_group.test.mail_nickname
}
