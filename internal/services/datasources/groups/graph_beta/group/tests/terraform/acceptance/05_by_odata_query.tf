# Test 05: Look up group using custom OData query
# This test creates a security group and then looks it up using a custom OData filter

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
  display_name     = "acc-test-by-odata-query-${random_string.test.result}"
  mail_nickname    = "acctestbyodataquery${random_string.test.result}"
  description      = "Test group for OData query lookup"
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
# Data Source - Lookup using custom OData query
# ==============================================================================

data "microsoft365_graph_beta_groups_group" "test" {
  odata_query = "displayName eq '${microsoft365_graph_beta_groups_group.test.display_name}' and securityEnabled eq true"

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
