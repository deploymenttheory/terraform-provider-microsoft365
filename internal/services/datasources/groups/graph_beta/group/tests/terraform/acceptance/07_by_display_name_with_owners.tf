# Test 07: Look up group by display_name and verify owners
# This test creates a security group with 1 owner user and verifies ownership

# ==============================================================================
# Random Suffix for Unique Resource Names
# ==============================================================================

resource "random_string" "test" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# User Resource (Owner)
# ==============================================================================

resource "microsoft365_graph_beta_users_user" "owner" {
  user_principal_name = "acc-test-owner-${random_string.test.result}@deploymenttheory.com"
  display_name        = "Test Owner ${random_string.test.result}"
  mail_nickname       = "acctestowner${random_string.test.result}"
  account_enabled     = true
  hard_delete         = true

  password_profile = {
    password                           = "P@ssw0rd${random_string.test.result}!"
    force_change_password_next_sign_in = false
  }
}

# ==============================================================================
# Group Resource with Owner
# ==============================================================================

resource "microsoft365_graph_beta_groups_group" "test" {
  display_name     = "acc-test-by-display-name-owners-${random_string.test.result}"
  mail_nickname    = "acctestbydisplaynameowners${random_string.test.result}"
  description      = "Test security group with owner for data source acceptance test"
  mail_enabled     = false
  security_enabled = true
  hard_delete      = true

  group_owners = [
    microsoft365_graph_beta_users_user.owner.id
  ]
}

# ==============================================================================
# Wait for Group Propagation
# ==============================================================================

resource "time_sleep" "wait_for_group" {
  depends_on      = [microsoft365_graph_beta_groups_group.test]
  create_duration = "30s"
}

# ==============================================================================
# Data Source - Lookup by display_name with security filter
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

output "owners_count" {
  value = length(data.microsoft365_graph_beta_groups_group.test.owners)
}

output "owners" {
  value = data.microsoft365_graph_beta_groups_group.test.owners
}
