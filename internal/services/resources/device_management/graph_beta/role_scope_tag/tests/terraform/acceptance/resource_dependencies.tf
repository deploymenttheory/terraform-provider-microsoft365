# Groups used for Role Scope Tag assignments testing
# These groups serve as dependencies for role scope tag assignment tests

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = true
  lower   = true
  numeric = true
}

resource "random_string" "group_suffix" {
  length  = 8
  special = false
  upper   = true
  lower   = true
  numeric = true
}

# Test Group 1 - IT Support Team
resource "microsoft365_graph_beta_groups_group" "acc_test_group_1" {
  display_name     = "acc-test-rst-it-${random_string.group_suffix.result}"
  description      = "Test group for IT support staff used in role scope tag assignments"
  mail_nickname    = "rst-it-${random_string.group_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  visibility       = "Private"

  timeouts = {
    create = "60s"
    read   = "30s"
    update = "30s"
    delete = "60s"
  }
}

# Test Group 2 - Device Management Team
resource "microsoft365_graph_beta_groups_group" "acc_test_group_2" {
  display_name     = "acc-test-rst-dm-${random_string.group_suffix.result}"
  description      = "Test group for device management staff used in role scope tag assignments"
  mail_nickname    = "rst-dm-${random_string.group_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  visibility       = "Private"

  timeouts = {
    create = "60s"
    read   = "30s"
    update = "30s"
    delete = "60s"
  }
}

# Test Group 3 - Helpdesk Team
resource "microsoft365_graph_beta_groups_group" "acc_test_group_3" {
  display_name     = "acc-test-rst-hd-${random_string.group_suffix.result}"
  description      = "Test group for helpdesk staff used in role scope tag assignments"
  mail_nickname    = "rst-hd-${random_string.group_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  visibility       = "Private"

  timeouts = {
    create = "60s"
    read   = "30s"
    update = "30s"
    delete = "60s"
  }
}

# Test Group 4 - Security Team
resource "microsoft365_graph_beta_groups_group" "acc_test_group_4" {
  display_name     = "acc-test-rst-sec-${random_string.group_suffix.result}"
  description      = "Test group for security staff used in role scope tag assignments"
  mail_nickname    = "rst-sec-${random_string.group_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  visibility       = "Private"

  timeouts = {
    create = "60s"
    read   = "30s"
    update = "30s"
    delete = "60s"
  }
}

# Outputs for easy reference in tests
output "test_group_ids" {
  description = "Group IDs for use in role scope tag assignment tests"
  value = {
    it_support_team        = microsoft365_graph_beta_groups_group.acc_test_group_1.id
    device_management_team = microsoft365_graph_beta_groups_group.acc_test_group_2.id
    helpdesk_team          = microsoft365_graph_beta_groups_group.acc_test_group_3.id
    security_team          = microsoft365_graph_beta_groups_group.acc_test_group_4.id
  }
}

# Convenience outputs for individual group IDs
output "it_support_team_id" {
  description = "IT Support Team group ID"
  value       = microsoft365_graph_beta_groups_group.acc_test_group_1.id
}

output "device_management_team_id" {
  description = "Device Management Team group ID"
  value       = microsoft365_graph_beta_groups_group.acc_test_group_2.id
}

output "helpdesk_team_id" {
  description = "Helpdesk Team group ID"
  value       = microsoft365_graph_beta_groups_group.acc_test_group_3.id
}

output "security_team_id" {
  description = "Security Team group ID"
  value       = microsoft365_graph_beta_groups_group.acc_test_group_4.id
}