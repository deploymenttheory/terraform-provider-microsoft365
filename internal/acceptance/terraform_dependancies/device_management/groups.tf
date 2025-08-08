# Groups used for acceptance testing
# These groups serve as dependencies.

resource "random_string" "group_suffix" {
  length  = 8
  special = false
  upper   = false
}

# Test Group 1 - IT Support Team
resource "microsoft365_graph_beta_groups_group" "acc_test_group_1" {
  display_name     = "Acceptance Test Dependency - IT Support Team - ${random_string.group_suffix.result}"
  description      = "Test group for IT support staff used in terms and conditions assignments"
  mail_nickname    = "it-support-acc-test-${random_string.group_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  visibility       = "Private"

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

# Test Group 2 - Device Management Team
resource "microsoft365_graph_beta_groups_group" "acc_test_group_2" {
  display_name     = "Acceptance Test Dependency - Device Management Team - ${random_string.group_suffix.result}"
  description      = "Test group for device management staff used in terms and conditions assignments"
  mail_nickname    = "device-mgmt-acc-test-${random_string.group_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  visibility       = "Private"

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

# Test Group 3 - Helpdesk Team
resource "microsoft365_graph_beta_groups_group" "acc_test_group_3" {
  display_name     = "Acceptance Test Dependency - Helpdesk Team - ${random_string.group_suffix.result}"
  description      = "Test group for helpdesk staff used in terms and conditions assignments"
  mail_nickname    = "helpdesk-acc-test-${random_string.group_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  visibility       = "Private"

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

# Test Group 4 - Security Team
resource "microsoft365_graph_beta_groups_group" "acc_test_group_4" {
  display_name     = "Acceptance Test Dependency - Security Team - ${random_string.group_suffix.result}"
  description      = "Test group for security staff used in terms and conditions assignments"
  mail_nickname    = "security-acc-test-${random_string.group_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  visibility       = "Private"

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

# Test Group 5 - DevOps Team
resource "microsoft365_graph_beta_groups_group" "acc_test_group_5" {
  display_name     = "Acceptance Test Dependency - DevOps Team - ${random_string.group_suffix.result}"
  description      = "Test group for devops staff used in terms and conditions assignments"
  mail_nickname    = "devops-acc-test-${random_string.group_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  visibility       = "Private"

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}


# Outputs for easy reference in tests
output "test_group_ids" {
  description = "Group IDs for use in terms and conditions assignment tests"
  value = {
    it_support_team        = microsoft365_graph_beta_groups_group.acc_test_group_1.id
    device_management_team = microsoft365_graph_beta_groups_group.acc_test_group_2.id
    helpdesk_team          = microsoft365_graph_beta_groups_group.acc_test_group_3.id
    security_team          = microsoft365_graph_beta_groups_group.acc_test_group_4.id
    devops_team            = microsoft365_graph_beta_groups_group.acc_test_group_5.id
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

output "devops_team_id" {
  description = "DevOps Team group ID"
  value       = microsoft365_graph_beta_groups_group.acc_test_group_5.id
}