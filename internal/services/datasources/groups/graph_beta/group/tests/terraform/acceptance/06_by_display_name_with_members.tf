# Test 06: Look up group by display_name and verify members
# This test creates a security group with 2 member users and verifies membership

# ==============================================================================
# Random Suffix for Unique Resource Names
# ==============================================================================

resource "random_string" "test" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# User Resources (Members)
# ==============================================================================

resource "microsoft365_graph_beta_users_user" "member1" {
  user_principal_name = "acc-test-member1-${random_string.test.result}@deploymenttheory.com"
  display_name        = "Test Member 1 ${random_string.test.result}"
  mail_nickname       = "acctestmember1${random_string.test.result}"
  account_enabled     = true
  hard_delete         = true
  
  password_profile = {
    password                           = "P@ssw0rd${random_string.test.result}!"
    force_change_password_next_sign_in = false
  }
}

resource "microsoft365_graph_beta_users_user" "member2" {
  user_principal_name = "acc-test-member2-${random_string.test.result}@deploymenttheory.com"
  display_name        = "Test Member 2 ${random_string.test.result}"
  mail_nickname       = "acctestmember2${random_string.test.result}"
  account_enabled     = true
  hard_delete         = true
  
  password_profile = {
    password                           = "P@ssw0rd${random_string.test.result}!"
    force_change_password_next_sign_in = false
  }
}

# ==============================================================================
# Group Resource with Members
# ==============================================================================

resource "microsoft365_graph_beta_groups_group" "test" {
  display_name     = "acc-test-by-display-name-members-${random_string.test.result}"
  mail_nickname    = "acctestbydisplaynamemembers${random_string.test.result}"
  description      = "Test security group with members for data source acceptance test"
  mail_enabled     = false
  security_enabled = true
  hard_delete      = true

  group_members = [
    microsoft365_graph_beta_users_user.member1.id,
    microsoft365_graph_beta_users_user.member2.id
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

output "members_count" {
  value = length(data.microsoft365_graph_beta_groups_group.test.members)
}

output "members" {
  value = data.microsoft365_graph_beta_groups_group.test.members
}
